package renderprovider

import (
	"image"
	"image/color"
	"io"
)

type WindowFlag uint32

const (
	WindowFlagVSync WindowFlag = 1 << iota
	WindowFlagFullScreen
	WindowFlagResizable
)

type FileType uint8

const (
	FileTypePng FileType = iota
)

type BlendMode int

const (
	BlendModeNone BlendMode = iota
	BlendModeAlpha
	BlendModeAdditive
	BlendModeMultiplied
	BlendModeAddColors
	BlendModeSubtractColors
)

type ImageColorMode int

const (
	ImageColorModeRGBA8Uncompressed ImageColorMode = iota
	ImageColorModeGrayscale
)

type RenderProvider interface {
	CreateWindow(width, height int, title string)
	SetLoggerCallback(callback func(logLevel int, s string))
	SetLoggerLevel(level int)
	SetWindowFlags(flags WindowFlag)
	SetWindowMinimumSize(width, height int)
	LoadImage(file io.Reader, fileType FileType) (Image, error)
	LoadTextureFromImage(img Image) (Texture, error)
	CreateRenderTexture(width, height int) (Texture, error)
	SetWindowIcon(image Image)
	FreeImage(image Image) error
	FreeTexture(texture Texture) error
	DrawTexture(texture Texture, x, y int, palette string, paletteOffset int) error
	DrawTextureEx(texture Texture, srcRect, destRec image.Rectangle, palette string, paletteOffset int) error
	DrawFontTexture(texture Texture, x, y int, palette string, color int) error
	NewImage(reader io.Reader, width, height int, imageColorMode ImageColorMode) (Image, error)
	BeginBlendMode(blendMode BlendMode)
	EndBlendMode()
	LoadPalette(name string, paletteStream io.Reader) error
	GetWindowSize() (width, height int)
	BeginDrawing()
	EndDrawing()
	IsRunning() bool
	LoadFontTTF(stream io.Reader, fontSize int) (Font, error)
	FreeFont(font Font) error
	DrawText(font Font, x, y int, color color.Color, text string) error
	ClearScreen(color color.Color)
	GetFPS() float32
	SetTargetFPS(fps int)
	GetFrameTime() float32
	CloseWindow()
	GetMousePosition() (x, y int)
	SetMouseVisible(visible bool)
	IsMouseButtonPressed(button MouseButton) bool
	BeginTextureMode(texture Texture) error
	EndTextureMode()
}
