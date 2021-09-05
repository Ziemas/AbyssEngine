package zrenderprovider

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"

	rp "github.com/OpenDiablo2/AbyssEngine/providers/renderprovider"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var ImageColorModeLookup = map[rp.ImageColorMode]PixelFormat{
	rp.ImageColorModeRGBA8Uncompressed: PixelFmtRGBA8,
	rp.ImageColorModeGrayscale:         PixelFmtGrayscale,
}

type ZRenderProvider struct {
	wflags rp.WindowFlag
	window *glfw.Window
	renderer *Renderer
	currentTime float64
	lastFrameTime float64
}

func New() *ZRenderProvider {
	result := &ZRenderProvider{}

	return result
}

func (z *ZRenderProvider) CreateWindow(width, height int, title string) {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	if z.wflags & rp.WindowFlagResizable != 0 {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	} else {
		glfw.WindowHint(glfw.Resizable, glfw.False)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// TODO fullscreen
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	z.window = window

	fbw = width
	fbh = height

	window.SetFramebufferSizeCallback(framebufferSizeCallback)
	window.SetCursorPosCallback(cursorPosCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)

	window.MakeContextCurrent()

	if z.wflags & rp.WindowFlagVSync != 0 {
		glfw.SwapInterval(0)
	} else {
		glfw.SwapInterval(1)
	}

	z.renderer = NewRenderer()

	z.currentTime = glfw.GetTime()
}

func (z *ZRenderProvider) SetWindowFlags(flags rp.WindowFlag) {
	z.wflags = flags
}

func (z *ZRenderProvider) SetWindowMinimumSize(width, height int) {
	z.window.SetSizeLimits(width, height, glfw.DontCare, glfw.DontCare)
}

func (z *ZRenderProvider) LoadImage(file io.Reader, fileType rp.FileType) (rp.Image, error) {
	if fileType != rp.FileTypePng {
		panic("unimplemented filetype")
	}

	png, err := png.Decode(file)
	data := image.NewRGBA(png.Bounds())
	draw.Draw(data, png.Bounds(), png, png.Bounds().Min, draw.Src)

	img := NewImage(data.Pix, png.Bounds().Size().X, png.Bounds().Size().Y, PixelFmtRGBA8)

	return img, err
}

func (z *ZRenderProvider) LoadTextureFromImage(img rp.Image) (rp.Texture, error) {
	image := img.(*zImage)
	return NewTextureFromImage(image)
}

func (z *ZRenderProvider) CreateRenderTexture(width, height int) (rp.Texture, error) {
	return NewRenderTexture(width, height)
}

func (z *ZRenderProvider) SetWindowIcon(image rp.Image) {
}

func (z *ZRenderProvider) FreeImage(image rp.Image) error {
	// Go object, will GC
	return nil
}

func (z *ZRenderProvider) FreeTexture(texture rp.Texture) error {
	return nil
}

func (z *ZRenderProvider) DrawTexture(texture rp.Texture, x, y int, palette string, paletteOffset int) error {
	z.renderer.DrawTexture(texture.(*zTexture), x, y, palette, paletteOffset)
	return nil
}

func (z *ZRenderProvider) DrawTextureEx(texture rp.Texture, srcRect, destRec image.Rectangle, palette string, paletteOffset int) error {
	z.renderer.DrawTextureEX(texture.(*zTexture), srcRect, destRec, palette, paletteOffset)
	return nil
}

func (z *ZRenderProvider) DrawFontTexture(texture rp.Texture, x, y int, palette string, color int) error {
	z.renderer.DrawTexture(texture.(*zTexture), x, y, palette, PaletteTextShiftOffset+color)
	return nil
}

func (z *ZRenderProvider) NewImage(reader io.Reader, width, height int, imageColorMode rp.ImageColorMode) (rp.Image, error) {
	data, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	result := NewImage(data, width, height, ImageColorModeLookup[imageColorMode])

	return result, nil
}

func (z *ZRenderProvider) BeginBlendMode(blendMode rp.BlendMode) {
}

func (z *ZRenderProvider) EndBlendMode() {
}

func (z *ZRenderProvider) LoadPalette(name string, paletteStream io.Reader) error {
	pal, err := NewPalette(paletteStream)
	if err != nil {
		return err
	}

	z.renderer.LoadPalette(name, pal)

	return nil
}

func (z *ZRenderProvider) GetWindowSize() (width, height int) {
	return z.window.GetFramebufferSize()
}

func (z *ZRenderProvider) BeginDrawing() {
}

func (z *ZRenderProvider) EndDrawing() {
	oldtime := z.currentTime
	z.currentTime = glfw.GetTime()
	z.lastFrameTime = z.currentTime - oldtime

	z.window.SwapBuffers()
	glfw.PollEvents()
}

func (z *ZRenderProvider) IsRunning() bool {
	return !z.window.ShouldClose()
}

func (z *ZRenderProvider) LoadFontTTF(stream io.Reader, fontSize int) (rp.Font, error) {
	return nil, nil
}

func (z *ZRenderProvider) FreeFont(font rp.Font) error {
	return nil
}

func (z *ZRenderProvider) DrawText(font rp.Font, x, y int, color color.Color, text string) error {
	return nil
}

func (z *ZRenderProvider) ClearScreen(color color.Color) {
	z.renderer.Clear() // TODO Color
}

func (z *ZRenderProvider) GetFPS() float32 {
	// TODO
	return 60.0
}

func (z *ZRenderProvider) SetTargetFPS(fps int) {
	// TODO?
}

func (z *ZRenderProvider) GetFrameTime() float32 {
	return float32(z.lastFrameTime)
}

func (z *ZRenderProvider) CloseWindow() {
	z.window.Destroy()
}

func (z *ZRenderProvider) GetMousePosition() (x, y int) {
	return int(cursorX), int(cursorY)
}

func (z *ZRenderProvider) SetMouseVisible(visible bool) {
	if visible == true {
		z.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	} else {
		z.window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	}
}

func (z *ZRenderProvider) IsMouseButtonPressed(button rp.MouseButton) bool {
	switch button {
	case rp.MouseButtonLeft:
		return leftDown
	case rp.MouseButtonRight:
		return rightDown
	case rp.MouseButtonMiddle:
		return middleDown
	}

	return false
}

func (z *ZRenderProvider) BeginTextureMode(texture rp.Texture) error {
	tex := texture.(*zTexture)
	z.renderer.BindRenderTexture(tex)

	return nil
}

func (z *ZRenderProvider) EndTextureMode() {
	z.renderer.UnbindRenderTexture()
}

func (z *ZRenderProvider) SetLoggerCallback(callback func(logLevel int, s string)) {
	debugCallback = callback
}

func (z *ZRenderProvider) SetLoggerLevel(level int) {
	maxLogLevel = level
}
