package raylibrenderprovider

import (
	"errors"
	"image"
	"image/color"
	"io"
	"io/ioutil"

	"github.com/OpenDiablo2/AbyssEngine/media"
	renderprovider "github.com/OpenDiablo2/AbyssEngine/providers/renderprovider"
	pl2 "github.com/OpenDiablo2/pl2/pkg"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var BlendModeLookup = map[renderprovider.BlendMode]rl.BlendMode{
	renderprovider.BlendModeAlpha:          rl.BlendAlpha,
	renderprovider.BlendModeAdditive:       rl.BlendAdditive,
	renderprovider.BlendModeMultiplied:     rl.BlendMultiplied,
	renderprovider.BlendModeAddColors:      rl.BlendAddColors,
	renderprovider.BlendModeSubtractColors: rl.BlendSubtractColors,
}

var ImageColorModeLookup = map[renderprovider.ImageColorMode]rl.PixelFormat{
	renderprovider.ImageColorModeRGBA8Uncompressed: rl.UncompressedR8g8b8a8,
	renderprovider.ImageColorModeGrayscale:         rl.UncompressedGrayscale,
}

type RaylibRenderProvider struct {
}

func (r *RaylibRenderProvider) BeginTextureMode(texture renderprovider.Texture) error {
	sourceTexture, ok := texture.(*Texture)

	if !ok {
		return errors.New("invalid texture")
	}

	if !sourceTexture.isRenderTexture {
		return errors.New("not a render texture")
	}

	rl.BeginTextureMode(sourceTexture.rtex)
	return nil
}

func (r *RaylibRenderProvider) EndTextureMode() {
	rl.EndTextureMode()
}

func (r *RaylibRenderProvider) CreateRenderTexture(width, height int) (renderprovider.Texture, error) {
	texture := rl.LoadRenderTexture(int32(width), int32(height))

	return &Texture{
		isRenderTexture: true,
		rtex:            texture,
	}, nil
}

func (r *RaylibRenderProvider) DrawTextureEx(texture renderprovider.Texture, srcRect, destRec image.Rectangle, palette string, paletteOffset int) error {
	sourceTexture, ok := texture.(*Texture)

	if !ok {
		return errors.New("invalid texture")
	}

	if palette != "" {
		tex := PaletteTexture[palette]

		if !tex.Init {
			img := rl.NewImage(tex.Data, 256, int32(PaletteTransformsCount), 1, rl.UncompressedR8g8b8a8)
			tex.Texture = rl.LoadTextureFromImage(img)

			tex.Init = true
		}

		rl.BeginShaderMode(PaletteShader)
		rl.SetShaderValueTexture(PaletteShader, PaletteShaderLoc, tex.Texture)
		rl.SetShaderValue(PaletteShader, PaletteShaderOffsetLoc, []float32{float32(paletteOffset)}, rl.ShaderUniformFloat)
	}

	rlSrc := rl.Rectangle{
		X:      float32(srcRect.Min.X),
		Y:      float32(srcRect.Min.Y),
		Width:  float32(srcRect.Size().X),
		Height: -float32(srcRect.Size().Y),
	}

	rlDest := rl.Rectangle{
		X:      float32(destRec.Min.X),
		Y:      float32(destRec.Min.Y),
		Width:  float32(destRec.Size().X),
		Height: float32(destRec.Size().Y),
	}

	rl.DrawTexturePro(sourceTexture.rtex.Texture, rlSrc, rlDest, rl.Vector2{}, 0.0, rl.White)

	if palette != "" {
		rl.EndShaderMode()
	}

	return nil
}

func (r *RaylibRenderProvider) IsMouseButtonPressed(button renderprovider.MouseButton) bool {
	rlButton, ok := mouseButtonLookup[button]

	if !ok {
		return false
	}

	return rl.IsMouseButtonDown(rlButton)
}

func (r *RaylibRenderProvider) SetTargetFPS(fps int) {
	rl.SetTargetFPS(int32(fps))
}

func (r *RaylibRenderProvider) SetMouseVisible(visible bool) {
	if visible {
		rl.ShowCursor()
		return
	}

	rl.HideCursor()
}

func (r *RaylibRenderProvider) GetMousePosition() (x, y int) {
	pos := rl.GetMousePosition()

	return int(pos.X), int(pos.Y)
}

func (r *RaylibRenderProvider) GetFrameTime() float32 {
	return rl.GetFrameTime()
}

func (r *RaylibRenderProvider) CloseWindow() {
	rl.CloseWindow()
}

func (r *RaylibRenderProvider) GetFPS() float32 {
	return rl.GetFPS()
}

func (r *RaylibRenderProvider) ClearScreen(color color.Color) {
	rr, g, b, a := color.RGBA()
	rl.ClearBackground(rl.Color{R: uint8(rr), G: uint8(g), B: uint8(b), A: uint8(a)})
}

func (r *RaylibRenderProvider) DrawText(font renderprovider.Font, x, y int, color color.Color, text string) error {
	fnt, ok := font.(*Font)

	if !ok {
		return errors.New("not a font")
	}

	rr, g, b, a := color.RGBA()
	rl.DrawTextEx(fnt.fnt, text, rl.Vector2{X: float32(x), Y: float32(y)}, float32(fnt.size), 0, rl.Color{
		R: uint8(rr), G: uint8(g), B: uint8(b), A: uint8(a),
	})

	return nil
}

func (r *RaylibRenderProvider) FreeFont(font renderprovider.Font) error {
	fnt, ok := font.(*Font)

	if !ok {
		return errors.New("not a font")
	}

	rl.UnloadFont(fnt.fnt)

	return nil
}

func (r *RaylibRenderProvider) LoadFontTTF(stream io.Reader, fontSize int) (renderprovider.Font, error) {
	bytes, err := io.ReadAll(stream)

	if err != nil {
		return nil, err
	}

	font := rl.LoadFontFromMemory(".ttf", bytes, int32(len(bytes)), int32(fontSize), nil, 0)
	rl.GenTextureMipmaps(&font.Texture)
	rl.SetTextureFilter(font.Texture, rl.FilterAnisotropic16x)

	return &Font{fnt: font, size: fontSize}, nil
}

func (r *RaylibRenderProvider) IsRunning() bool {
	return !rl.WindowShouldClose()
}

func (r *RaylibRenderProvider) BeginDrawing() {
	rl.BeginDrawing()
}

func (r *RaylibRenderProvider) EndDrawing() {
	rl.EndDrawing()
}

func (r *RaylibRenderProvider) GetWindowSize() (width, height int) {
	return rl.GetScreenWidth(), rl.GetScreenHeight()
}

func (r *RaylibRenderProvider) BeginBlendMode(blendMode renderprovider.BlendMode) {
	rl.BeginBlendMode(BlendModeLookup[blendMode])
}

func (r *RaylibRenderProvider) EndBlendMode() {
	rl.EndBlendMode()
}

func (r *RaylibRenderProvider) NewImage(reader io.Reader, width, height int, imageColorMode renderprovider.ImageColorMode) (renderprovider.Image, error) {
	data, err := ioutil.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	result := rl.NewImage(data, int32(width), int32(height), 1, ImageColorModeLookup[imageColorMode])

	return &Image{src: result}, nil
}

func (r *RaylibRenderProvider) DrawTexture(texture renderprovider.Texture, x, y int, palette string, paletteOffset int) error {
	sourceTexture, ok := texture.(*Texture)

	if !ok {
		return errors.New("invalid texture")
	}

	if palette != "" {
		tex := PaletteTexture[palette]

		if !tex.Init {
			img := rl.NewImage(tex.Data, 256, int32(PaletteTransformsCount), 1, rl.UncompressedR8g8b8a8)
			tex.Texture = rl.LoadTextureFromImage(img)

			tex.Init = true
		}

		rl.BeginShaderMode(PaletteShader)
		rl.SetShaderValueTexture(PaletteShader, PaletteShaderLoc, tex.Texture)
		rl.SetShaderValue(PaletteShader, PaletteShaderOffsetLoc, []float32{float32(paletteOffset)}, rl.ShaderUniformFloat)
	}
	rl.DrawTexture(sourceTexture.tex, int32(x), int32(y), rl.White)
	if palette != "" {
		rl.EndShaderMode()
	}

	return nil
}

func (r *RaylibRenderProvider) DrawFontTexture(texture renderprovider.Texture, x, y int, palette string, color int) error {
	sourceTexture, ok := texture.(*Texture)

	if !ok {
		return errors.New("invalid texture")
	}

	tex := PaletteTexture[palette]

	if !tex.Init {
		img := rl.NewImage(tex.Data, 256, int32(PaletteTransformsCount), 1, rl.UncompressedR8g8b8a8)
		tex.Texture = rl.LoadTextureFromImage(img)

		tex.Init = true
	}

	rl.BeginShaderMode(PaletteShader)
	rl.SetShaderValueTexture(PaletteShader, PaletteShaderLoc, tex.Texture)
	rl.SetShaderValue(PaletteShader, PaletteShaderOffsetLoc, []float32{float32(color+PaletteTextShiftOffset) / float32(PaletteTransformsCount-1)}, rl.ShaderUniformFloat)
	rl.DrawTexture(sourceTexture.tex, int32(x), int32(y), rl.White)
	rl.EndShaderMode()

	return nil
}

func (r *RaylibRenderProvider) FreeImage(image renderprovider.Image) error {
	img, ok := image.(*Image)

	if !ok {
		return errors.New("invalid image")
	}

	rl.UnloadImage(img.src)
	return nil
}

func (r *RaylibRenderProvider) FreeTexture(texture renderprovider.Texture) error {
	tex, ok := texture.(*Texture)

	if !ok {
		return errors.New("invalid texture")
	}

	rl.UnloadTexture(tex.tex)
	return nil
}

func (r *RaylibRenderProvider) LoadTextureFromImage(image renderprovider.Image) (renderprovider.Texture, error) {
	img, ok := image.(*Image)

	if !ok {
		return nil, errors.New("invalid image")
	}

	result := rl.LoadTextureFromImage(img.src)

	return &Texture{tex: result}, nil
}

func (r *RaylibRenderProvider) SetWindowIcon(image renderprovider.Image) {
	img, ok := image.(*Image)

	if !ok {
		return
	}

	rl.SetWindowIcon(*img.src)
}

var fileTypeMap = map[renderprovider.FileType]string{
	renderprovider.FileTypePng: ".png",
}

func (r *RaylibRenderProvider) LoadImage(file io.Reader, fileType renderprovider.FileType) (renderprovider.Image, error) {
	fileExt, ok := fileTypeMap[fileType]

	if !ok {
		return nil, errors.New("unsupported file type")
	}

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		return nil, err
	}

	result := rl.LoadImageFromMemory(fileExt, bytes, int32(len(bytes)))

	return &Image{src: result}, nil
}

func (r *RaylibRenderProvider) SetWindowMinimumSize(width, height int) {
	rl.SetWindowMinSize(width, height)
}

func (r *RaylibRenderProvider) SetWindowFlags(flags renderprovider.WindowFlag) {
	var rlFlags byte = 0

	if flags&renderprovider.WindowFlagResizable > 0 {
		rlFlags |= rl.FlagWindowResizable
	}

	if flags&renderprovider.WindowFlagFullScreen > 0 {
		rlFlags |= rl.FlagFullscreenMode
	}

	if flags&renderprovider.WindowFlagVSync > 0 {
		rlFlags |= rl.FlagVsyncHint
	}

	rl.SetConfigFlags(rlFlags)
}

func (r *RaylibRenderProvider) SetLoggerLevel(level int) {
	rl.SetTraceLog(level)
}

func (r *RaylibRenderProvider) SetLoggerCallback(callback func(logLevel int, s string)) {
	rl.SetTraceLogCallback(callback)
}

func (r *RaylibRenderProvider) CreateWindow(width, height int, title string) {
	rl.InitWindow(int32(width), int32(height), title)

	PaletteShader = rl.LoadShaderFromMemory(media.StandardVertexShader, media.PaletteFragmentShader)
	PaletteShaderLoc = rl.GetShaderLocation(PaletteShader, "palette")
	PaletteShaderOffsetLoc = rl.GetShaderLocation(PaletteShader, "paletteOffset")

}

func (r *RaylibRenderProvider) LoadPalette(name string, paletteStream io.Reader) error {
	if PaletteTexture == nil {
		PaletteTexture = make(map[string]*PalTex)
	}

	paletteBytes, err := ioutil.ReadAll(paletteStream)

	if err != nil {
		return err
	}

	pal, err := pl2.FromBytes(paletteBytes)

	if err != nil {
		return err
	}

	colors := make([]uint8, 0)
	colors = append(colors, palToSlice(pal.BasePalette)...)

	PaletteTextShiftOffset = len(colors) / (256 * 4)

	for idx := range pal.TextColorShifts {
		colors = append(colors, transformToSlice(pal.BasePalette, pal.TextColorShifts[idx])...)
	}

	PaletteTransformsCount = len(colors) / (256 * 4)

	tex := &PalTex{}

	tex.Data = colors
	tex.Init = false

	PaletteTexture[name] = tex

	return nil
}

func New() *RaylibRenderProvider {
	result := &RaylibRenderProvider{}

	return result
}
