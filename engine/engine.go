package engine

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"image"
	"math"
	"runtime"

	"github.com/OpenDiablo2/AbyssEngine/common"
	"github.com/OpenDiablo2/AbyssEngine/loader"
	"github.com/OpenDiablo2/AbyssEngine/loader/filesystemloader"
	"github.com/OpenDiablo2/AbyssEngine/media"
	"github.com/OpenDiablo2/AbyssEngine/node"
	"github.com/OpenDiablo2/AbyssEngine/node/sprite"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/rs/zerolog/log"
)

// Engine represents the main game engine
type Engine struct {
	config           Configuration
	loader           *loader.Loader
	renderSurface    rl.RenderTexture2D
	systemFont       rl.Font
	bootLogo         rl.Texture2D
	bootLoadText     string
	shutdown         bool
	engineMode       EngineMode
	cursorSprite     *sprite.Sprite
	rootNode         *node.Node
	cursorX          int
	cursorY          int
	toggleFullscreen bool
	isFullscreen     bool
	cursorOffset     image.Point
	luaState         *lua.LState
	currentBlendMode common.BlendMode
}

func (e *Engine) SetBlendMode(mode common.BlendMode) {
	if e.currentBlendMode == mode {
		return
	}

	if e.currentBlendMode != common.BlendModeNone {
		rl.EndBlendMode()
	}

	e.currentBlendMode = mode

	if mode == common.BlendModeNone {
		return
	}

	rlBlendMode := common.BlendModeLookup[mode]
	rl.BeginBlendMode(rlBlendMode)
}

func (e *Engine) GetMousePosition() (X, Y int) {
	return e.cursorX, e.cursorY
}

func (e *Engine) GetLanguageCode() string {
	// TODO: Make this dynamic
	return "eng"
}

func (e *Engine) GetLanguageFontCode() string {
	// TODO: Make this dynamic
	return "latin"
}

// New creates a new instance of the engine
func New(config Configuration) *Engine {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(800, 600, "Abyss Engine")
	rl.SetWindowMinSize(800, 600)

	windowIcon := rl.LoadImageFromMemory(".png", media.AbyssIcon, int32(len(media.BootLogo)))
	rl.SetWindowIcon(*windowIcon)

	result := &Engine{
		shutdown:         false,
		toggleFullscreen: false,
		isFullscreen:     false,
		config:           config,
		engineMode:       EngineModeBoot,
		renderSurface:    rl.LoadRenderTexture(800, 600),
		systemFont:       rl.LoadFontFromMemory(".ttf", media.FontDiabloHeavy, int32(len(media.FontDiabloHeavy)), 18, nil, 0),
		rootNode:         node.New(),
		cursorOffset:     image.Point{},
		currentBlendMode: common.BlendModeNone,
	}

	result.loader = loader.New(result)
	result.loader.AddProvider(filesystemloader.New(config.RootPath))

	logo := rl.LoadImageFromMemory(".png", media.BootLogo, int32(len(media.BootLogo)))
	result.bootLogo = rl.LoadTextureFromImage(logo)
	rl.UnloadImage(logo)

	rl.GenTextureMipmaps(&result.systemFont.Texture)
	rl.SetTextureFilter(result.systemFont.Texture, rl.FilterAnisotropic16x)

	common.PaletteShader = rl.LoadShaderFromMemory(media.StandardVertexShader, media.PaletteFragmentShader)
	common.PaletteShaderLoc = rl.GetShaderLocation(common.PaletteShader, "palette")
	common.PaletteShaderOffsetLoc = rl.GetShaderLocation(common.PaletteShader, "paletteOffset")
	return result
}

// Destroy finalizes the instance of the engine
func (e *Engine) Destroy() {
	rl.UnloadTexture(e.bootLogo)
	rl.UnloadFont(e.systemFont)
}

// Run runs the engine
func (e *Engine) Run() {
	e.bootstrapScripts()

	for !rl.WindowShouldClose() {
		if e.shutdown {
			break
		}

		rl.BeginDrawing()
		rl.BeginTextureMode(e.renderSurface)
		rl.ClearBackground(rl.Black)

		switch e.engineMode {
		case EngineModeBoot:
			e.showBootSplash()
		case EngineModeGame:
			rl.BeginShaderMode(common.PaletteShader)
			e.currentBlendMode = common.BlendModeNone
			e.showGame()
			e.SetBlendMode(common.BlendModeNone)
			rl.EndShaderMode()
		}

		rl.EndTextureMode()
		e.drawMainSurface()
		rl.EndDrawing()

		if e.engineMode == EngineModeGame {
			e.updateGame(float64(rl.GetFrameTime()))
		}

		if e.toggleFullscreen {
			e.toggleFullscreen = false
			if e.isFullscreen {
				rl.ClearWindowState(rl.FlagFullscreenMode)
			} else {
				rl.SetWindowState(rl.FlagFullscreenMode)
			}
			e.isFullscreen = !e.isFullscreen
		}
	}

	e.luaState.Close()
	rl.CloseWindow()
}

func (e *Engine) showGame() {
	e.rootNode.Render()
	if e.cursorSprite != nil {
		e.cursorSprite.Render()
	}
}

func (e *Engine) updateGame(elapsed float64) {
	e.rootNode.Update(elapsed)
	if e.cursorSprite != nil {
		scale := float32(math.Min(float64(rl.GetScreenWidth())/800.0, float64(rl.GetScreenHeight())/600.0))
		xOrigin := (float32(rl.GetScreenWidth()) - (800.0 * scale)) * 0.5
		yOrigin := (float32(rl.GetScreenHeight()) - (600.0 * scale)) * 0.5

		e.cursorX = int((float32(rl.GetMouseX()) - xOrigin) * (1.0 / scale))
		e.cursorSprite.X = e.cursorX + e.cursorOffset.X

		e.cursorY = int((float32(rl.GetMouseY()) - yOrigin) * (1.0 / scale))
		e.cursorSprite.Y = e.cursorY + e.cursorOffset.Y

		e.cursorSprite.Update(elapsed)
	}
}

func (e *Engine) showBootSplash() {
	rl.DrawTexture(e.bootLogo, int32(rl.GetScreenWidth()/3)-(e.bootLogo.Width/2),
		int32(rl.GetScreenHeight()/2)-(e.bootLogo.Height/2), rl.White)

	textX := float32(rl.GetScreenWidth()) / 2
	textY := float32(rl.GetScreenHeight()/2) - 20

	rl.DrawTextEx(e.systemFont, "Abyss Engine", rl.Vector2{X: textX, Y: textY}, 18, 0, rl.White)
	rl.DrawTextEx(e.systemFont, "Local Build", rl.Vector2{X: textX, Y: textY + 16}, 18, 0, rl.Gray)
	rl.DrawTextEx(e.systemFont, e.bootLoadText,
		rl.Vector2{X: float32(rl.GetScreenWidth() / 4), Y: float32(rl.GetScreenWidth()/4) * 2.5}, 18, 0, rl.Beige)
}

func (e *Engine) drawMainSurface() {
	rl.ClearBackground(rl.Black)
	scale := float32(math.Min(float64(rl.GetScreenWidth())/800.0, float64(rl.GetScreenHeight())/600.0))

	rl.DrawTexturePro(e.renderSurface.Texture,
		rl.Rectangle{Width: float32(e.renderSurface.Texture.Width), Height: float32(-e.renderSurface.Texture.Height)},
		rl.Rectangle{
			X:      (float32(rl.GetScreenWidth()) - (800.0 * scale)) * 0.5,
			Y:      (float32(rl.GetScreenHeight()) - (600.0 * scale)) * 0.5,
			Width:  800.0 * scale,
			Height: 600.0 * scale},
		rl.Vector2{}, 0.0, rl.White)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	rl.DrawTextEx(e.systemFont, fmt.Sprintf("FPS: %d", int(rl.GetFPS())), rl.Vector2{X: 5, Y: 5}, 18, 0, rl.White)
	rl.DrawTextEx(e.systemFont, fmt.Sprintf("GC: %d (%%%d)", int(memStats.NumGC), int(memStats.GCCPUFraction*100)), rl.Vector2{X: 5, Y: 21}, 18, 0, rl.White)
	rl.DrawTextEx(e.systemFont, fmt.Sprintf("Alloc: %0.2fMB (%0.2fMB)", float32(memStats.Alloc)/1024/1024, float32(memStats.Sys)/1024/1024), rl.Vector2{X: 5, Y: 37}, 18, 0, rl.White)

}

func (e *Engine) panic(msg string) {
	// TODO: This should be a UI screen
	log.Fatal().Msg(msg)
	rl.CloseWindow()
}
