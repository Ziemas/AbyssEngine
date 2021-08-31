package engine

import (
	"bytes"
	"image"
	"image/png"
	"runtime"

	lua "github.com/yuin/gopher-lua"

	"github.com/OpenDiablo2/AbyssEngine/common"
	"github.com/OpenDiablo2/AbyssEngine/loader"
	"github.com/OpenDiablo2/AbyssEngine/loader/filesystemloader"
	"github.com/OpenDiablo2/AbyssEngine/media"
	"github.com/OpenDiablo2/AbyssEngine/node"
	"github.com/OpenDiablo2/AbyssEngine/node/sprite"
	ren "github.com/OpenDiablo2/AbyssEngine/renderer"
	"github.com/rs/zerolog/log"

	"github.com/go-gl/glfw/v3.3/glfw"
)

// Engine represents the main game engine
type Engine struct {
	config Configuration
	loader *loader.Loader
	window *glfw.Window
	//renderer *ren.Renderer
	//renderSurface    rl.RenderTexture2D
	//systemFont       rl.Font
	//bootLogo         rl.Texture2D
	bootLogo         ren.Texture
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
	currentBlendMode ren.BlendMode
	prevTime         float64
}

func (e *Engine) SetBlendMode(mode ren.BlendMode) {
	if e.currentBlendMode == mode {
		return
	}

	ren.SetBlendMode(mode)

	if e.currentBlendMode != ren.BlendModeNone {
		//rl.EndBlendMode()
	}

	e.currentBlendMode = mode

	if mode == ren.BlendModeNone {
		return
	}

	//rlBlendMode := common.BlendModeLookup[mode]
	//rl.BeginBlendMode(rlBlendMode)
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
	runtime.LockOSThread()

	window := ren.Init()

	icon, _ := png.Decode(bytes.NewReader(media.AbyssIcon))
	window.SetIcon([]image.Image{icon, icon})

	logo, _ := png.Decode(bytes.NewReader(media.BootLogo))
	logoTex := ren.NewTextureRGBA(logo)

	result := &Engine{
		shutdown:         false,
		toggleFullscreen: false,
		isFullscreen:     false,
		config:           config,
		engineMode:       EngineModeBoot,
		window:           window,
		bootLogo:         logoTex,
		//renderer:         renderer,
		//renderSurface:    rl.LoadRenderTexture(800, 600),
		//systemFont:       rl.LoadFontFromMemory(".ttf", media.FontDiabloHeavy, int32(len(media.FontDiabloHeavy)), 18, nil, 0),
		rootNode:         node.New(),
		cursorOffset:     image.Point{},
		currentBlendMode: ren.BlendModeNone,
	}

	result.loader = loader.New(result)
	result.loader.AddProvider(filesystemloader.New(config.RootPath))
	//logo := rl.LoadImageFromMemory(".png", media.BootLogo, int32(len(media.BootLogo)))
	//result.bootLogo = rl.LoadTextureFromImage(logo)
	//rl.UnloadImage(logo)

	//rl.GenTextureMipmaps(&result.systemFont.Texture)
	//rl.SetTextureFilter(result.systemFont.Texture, rl.FilterAnisotropic16x)

	common.StandardShader = ren.NewProgram(media.StandardVertexShader, media.StandardFragmentShader)
	common.PaletteShader = ren.NewProgram(media.StandardVertexShader, media.PaletteFragmentShader)
	common.PaletteShaderLoc = ren.GetShaderLocation(common.PaletteShader, "paletteTex")
	common.PaletteShaderOffsetLoc = ren.GetShaderLocation(common.PaletteShader, "paletteOffset")
	return result
}

// Destroy finalizes the instance of the engine
func (e *Engine) Destroy() {
	//rl.UnloadTexture(e.bootLogo)
	//rl.UnloadFont(e.systemFont)
}

// Run runs the engine
func (e *Engine) Run() {
	e.bootstrapScripts()

	for !e.window.ShouldClose() {
		if e.shutdown {
			break
		}

		//rl.BeginDrawing()
		//rl.BeginTextureMode(e.renderSurface)
		//rl.ClearBackground(rl.Black)
		ren.Clear()

		switch e.engineMode {
		case EngineModeBoot:
			ren.BeginShader(common.StandardShader)
			e.showBootSplash()
			ren.EndShader()
		case EngineModeGame:
			ren.BeginShader(common.PaletteShader)
			e.currentBlendMode = ren.BlendModeNone
			e.showGame()
			e.SetBlendMode(ren.BlendModeNone)
			ren.EndShader()
		}

		//rl.EndTextureMode()
		e.drawMainSurface()
		e.window.SwapBuffers()
		glfw.PollEvents()
		//rl.EndDrawing()

		time := glfw.GetTime()
		delta := time - e.prevTime
		e.prevTime = time

		if e.engineMode == EngineModeGame {
			e.updateGame(delta)
		}

		if e.toggleFullscreen {
			e.toggleFullscreen = false
			if e.isFullscreen {
				//rl.ClearWindowState(rl.FlagFullscreenMode)
			} else {
				//rl.SetWindowState(rl.FlagFullscreenMode)
			}
			e.isFullscreen = !e.isFullscreen
		}
	}

	e.luaState.Close()
	//rl.CloseWindow()
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
		scale := float32(math.Min(float64(ren.GetScreenWidth())/800.0, float64(ren.GetScreenHeight())/600.0))
		xOrigin := (float32(ren.GetScreenWidth()) - (800.0 * scale)) * 0.5
		yOrigin := (float32(ren.GetScreenHeight()) - (600.0 * scale)) * 0.5

		e.cursorX = int((float32(ren.GetMouseX()) - xOrigin) * (1.0 / scale))
		e.cursorSprite.X = e.cursorX + e.cursorOffset.X

		e.cursorY = int((float32(ren.GetMouseY()) - yOrigin) * (1.0 / scale))
		e.cursorSprite.Y = e.cursorY + e.cursorOffset.Y

		e.cursorSprite.Update(elapsed)
	}
}

func (e *Engine) showBootSplash() {
	//rl.DrawTexture(e.bootLogo, int32(rl.GetScreenWidth()/3)-(e.bootLogo.Width/2),
	//	int32(rl.GetScreenHeight()/2)-(e.bootLogo.Height/2), rl.White)

	ren.DrawTexture(e.bootLogo, (ren.GetScreenWidth()/3)-(e.bootLogo.Width/2),
		(ren.GetScreenHeight()/2)-(e.bootLogo.Height/2))
	//ren.DrawTexture(e.bootLogo, (ren.GetScreenWidth()/2)+(e.bootLogo.Width), ren.GetScreenHeight()/2+(e.bootLogo.Height))

	//textX := float32(rl.GetScreenWidth()) / 2
	//textY := float32(rl.GetScreenHeight()/2) - 20

	//rl.DrawTextEx(e.systemFont, "Abyss Engine", rl.Vector2{X: textX, Y: textY}, 18, 0, rl.White)
	//rl.DrawTextEx(e.systemFont, "Local Build", rl.Vector2{X: textX, Y: textY + 16}, 18, 0, rl.Gray)
	//rl.DrawTextEx(e.systemFont, e.bootLoadText,
	//	rl.Vector2{X: float32(rl.GetScreenWidth() / 4), Y: float32(rl.GetScreenWidth()/4) * 2.5}, 18, 0, rl.Beige)
}

func (e *Engine) drawMainSurface() {
	//rl.ClearBackground(rl.Black)
	//scale := float32(math.Min(float64(rl.GetScreenWidth())/800.0, float64(rl.GetScreenHeight())/600.0))

	//rl.DrawTexturePro(e.renderSurface.Texture,
	//	rl.Rectangle{Width: float32(e.renderSurface.Texture.Width), Height: float32(-e.renderSurface.Texture.Height)},
	//	rl.Rectangle{
	//		X:      (float32(rl.GetScreenWidth()) - (800.0 * scale)) * 0.5,
	//		Y:      (float32(rl.GetScreenHeight()) - (600.0 * scale)) * 0.5,
	//		Width:  800.0 * scale,
	//		Height: 600.0 * scale},
	//	rl.Vector2{}, 0.0, rl.White)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	//rl.DrawTextEx(e.systemFont, fmt.Sprintf("FPS: %d", int(rl.GetFPS())), rl.Vector2{X: 5, Y: 5}, 18, 0, rl.White)
	//rl.DrawTextEx(e.systemFont, fmt.Sprintf("GC: %d (%%%d)", int(memStats.NumGC), int(memStats.GCCPUFraction*100)), rl.Vector2{X: 5, Y: 21}, 18, 0, rl.White)
	//rl.DrawTextEx(e.systemFont, fmt.Sprintf("Alloc: %0.2fMB (%0.2fMB)", float32(memStats.Alloc)/1024/1024, float32(memStats.Sys)/1024/1024), rl.Vector2{X: 5, Y: 37}, 18, 0, rl.White)

}

func (e *Engine) panic(msg string) {
	// TODO: This should be a UI screen
	log.Fatal().Msg(msg)
	//rl.CloseWindow()
}
