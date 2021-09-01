package engine

import (
	"bytes"
	"fmt"
	"github.com/OpenDiablo2/AbyssEngine/loader"
	"github.com/OpenDiablo2/AbyssEngine/loader/filesystemloader"
	"github.com/OpenDiablo2/AbyssEngine/media"
	"github.com/OpenDiablo2/AbyssEngine/node"
	"github.com/OpenDiablo2/AbyssEngine/node/sprite"
	"github.com/OpenDiablo2/AbyssEngine/providers/renderprovider"
	"github.com/rs/zerolog/log"
	lua "github.com/yuin/gopher-lua"
	"image"
	"image/color"
	"math"
	"runtime"
)

// Engine represents the main game engine
type Engine struct {
	config           Configuration
	renderProvider   renderprovider.RenderProvider
	loader           *loader.Loader
	renderSurface    renderprovider.Texture
	systemFont       renderprovider.Font
	bootLogo         renderprovider.Texture
	bootLoadText     string
	shutdown         bool
	engineMode       EngineMode
	cursorSprite     *sprite.Sprite
	rootNode         *node.Node
	cursorX          int
	cursorY          int
	scale            float32
	xOrigin          float32
	yOrigin          float32
	lastScreenWidth  int
	lastScreenHeight int
	cursorOffset     image.Point
	luaState         *lua.LState
	currentBlendMode renderprovider.BlendMode
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
func New(config Configuration, renderProvider renderprovider.RenderProvider) *Engine {
	renderProvider.SetWindowFlags(renderprovider.WindowFlagResizable)
	renderProvider.CreateWindow(800, 600, "Abyss Engine")
	renderProvider.SetWindowMinimumSize(800, 600)

	windowIcon, _ := renderProvider.LoadImage(bytes.NewReader(media.AbyssIcon), renderprovider.FileTypePng)
	renderProvider.SetWindowIcon(windowIcon)

	result := &Engine{
		renderProvider:   renderProvider,
		shutdown:         false,
		config:           config,
		engineMode:       EngineModeBoot,
		rootNode:         node.New(),
		cursorOffset:     image.Point{},
		currentBlendMode: renderprovider.BlendModeNone,
	}

	result.renderSurface, _ = renderProvider.CreateRenderTexture(800, 600)
	result.loader = loader.New(result)
	result.loader.AddProvider(filesystemloader.New(config.RootPath))

	logo, _ := renderProvider.LoadImage(bytes.NewReader(media.BootLogo), renderprovider.FileTypePng)
	result.bootLogo, _ = renderProvider.LoadTextureFromImage(logo)
	_ = renderProvider.FreeImage(logo)

	result.systemFont, _ = renderProvider.LoadFontTTF(bytes.NewReader(media.FontDiabloHeavy), 18)

	return result
}

// Destroy finalizes the instance of the engine
func (e *Engine) Destroy() {
	_ = e.renderProvider.FreeTexture(e.bootLogo)
	_ = e.renderProvider.FreeFont(e.systemFont)
}

// Run runs the engine
func (e *Engine) Run() {
	e.bootstrapScripts()

	for e.renderProvider.IsRunning() {
		if e.shutdown {
			break
		}

		newScreenWidth, newScreenHeight := e.renderProvider.GetWindowSize()

		if (newScreenWidth != e.lastScreenWidth) || (newScreenHeight != e.lastScreenHeight) {
			e.scale = float32(math.Min(float64(newScreenWidth)/800.0, float64(newScreenHeight)/600.0))
			e.xOrigin = (float32(newScreenWidth) - (800.0 * e.scale)) * 0.5
			e.yOrigin = (float32(newScreenHeight) - (600.0 * e.scale)) * 0.5
		}

		mousePosX, mousePosY := e.renderProvider.GetMousePosition()
		e.cursorX = int((float32(mousePosX) - e.xOrigin) * (1.0 / e.scale))
		e.cursorY = int((float32(mousePosY) - e.yOrigin) * (1.0 / e.scale))

		_ = e.renderProvider.BeginTextureMode(e.renderSurface)
		e.renderProvider.ClearScreen(color.Black)
		switch e.engineMode {
		case EngineModeBoot:
			e.showBootSplash()
		case EngineModeGame:
			e.showGame()
		}
		e.renderProvider.EndTextureMode()

		e.renderProvider.BeginDrawing()
		e.drawMainSurface()
		e.renderProvider.EndDrawing()

		if e.engineMode == EngineModeGame {
			e.updateGame(float64(e.renderProvider.GetFrameTime()))
		}

	}

	e.luaState.Close()
	e.renderProvider.CloseWindow()
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
		e.cursorSprite.X = e.cursorX + e.cursorOffset.X
		e.cursorSprite.Y = e.cursorY + e.cursorOffset.Y
		e.cursorSprite.Update(elapsed)
	}
}

func (e *Engine) showBootSplash() {
	screenWidth, screenHeight := e.renderProvider.GetWindowSize()

	_ = e.renderProvider.DrawTexture(e.bootLogo,
		(screenWidth/3)-(e.bootLogo.Width()/2),
		(screenHeight/2)-(e.bootLogo.Height()/2), "", 0)

	textX := float32(screenWidth) / 2
	textY := float32(screenHeight/2) - 20

	clrGray := color.RGBA{R: 0x7C, G: 0x7C, B: 0x7C, A: 0xFF}
	clrBeige := color.RGBA{R: 211, G: 176, B: 131, A: 255}

	_ = e.renderProvider.DrawText(e.systemFont, int(textX), int(textY), color.White, "Abyss Engine")
	_ = e.renderProvider.DrawText(e.systemFont, int(textX), int(textY+16), clrGray, "Local Build")

	blX := screenWidth / 4
	blY := int(float32(screenHeight/4) * 2.5)
	_ = e.renderProvider.DrawText(e.systemFont, blX, blY, clrBeige, e.bootLoadText)
}

func (e *Engine) drawMainSurface() {
	e.renderProvider.ClearScreen(color.Black)
	screenWidth, screenHeight := e.renderProvider.GetWindowSize()

	x := (float32(screenWidth) - (800.0 * e.scale)) * 0.5
	y := (float32(screenHeight) - (600.0 * e.scale)) * 0.5
	width := 800.0 * e.scale
	height := 600.0 * e.scale

	//rl.DrawTexturePro(e.renderSurface.Texture,
	//	rl.Rectangle{Width: float32(e.renderSurface.Texture.Width), Height: float32(-e.renderSurface.Texture.Height)},
	//	rl.Rectangle{
	//		X:      (float32(rl.GetScreenWidth()) - (800.0 * scale)) * 0.5,
	//		Y:      (float32(rl.GetScreenHeight()) - (600.0 * scale)) * 0.5,
	//		Width:  800.0 * scale,
	//		Height: 600.0 * scale},
	//	rl.Vector2{}, 0.0, rl.White)

	_ = e.renderProvider.DrawTextureEx(e.renderSurface,
		image.Rectangle{
			Max: image.Point{X: e.renderSurface.Width(), Y: e.renderSurface.Height()},
		},
		image.Rectangle{
			Min: image.Point{X: int(x), Y: int(y)},
			Max: image.Point{X: int(x + width), Y: int(y + height)},
		}, "", 0)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	fps := int(e.renderProvider.GetFPS())
	_ = e.renderProvider.DrawText(e.systemFont, 5, 5, color.White, fmt.Sprintf("FPS: %d", fps))
	_ = e.renderProvider.DrawText(e.systemFont, 5, 21, color.White, fmt.Sprintf("GC: %d (%%%d)", int(memStats.NumGC), int(memStats.GCCPUFraction*100)))
	_ = e.renderProvider.DrawText(e.systemFont, 5, 37, color.White, fmt.Sprintf("Alloc: %0.2fMB (%0.2fMB)", float32(memStats.Alloc)/1024/1024, float32(memStats.Sys)/1024/1024))
}

func (e *Engine) panic(msg string) {
	// TODO: This should be a UI screen
	log.Fatal().Msg(msg)
	e.renderProvider.CloseWindow()
}
