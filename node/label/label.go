package label

import (
	bytes2 "bytes"
	"errors"
	"image"
	"io"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/OpenDiablo2/AbyssEngine/common"
	"github.com/OpenDiablo2/AbyssEngine/node"
	dc6 "github.com/OpenDiablo2/dc6/pkg"
	tblfont "github.com/OpenDiablo2/tbl_font/pkg"
)

type Label struct {
	*node.Node

	initialized bool
	hasTexture  bool
	texture     rl.Texture2D
	FontTable   *tblfont.FontTable
	FontGfx     common.SequenceProvider
	Palette     string
	Caption     string
}

func New(loaderProvider common.LoaderProvider, fontPath, palette string) (*Label, error) {
	result := &Label{
		Node:        node.New(),
		initialized: false,
	}

	_, ok := common.PaletteTexture[palette]
	if !ok {
		return nil, errors.New("sprite loaded with non-existent palette")
	}
	result.Palette = palette

	fontTableStream, err := loaderProvider.Load(fontPath + ".tbl")
	defer fontTableStream.Close()

	if err != nil {
		return nil, err
	}

	// hack: mpq block stream is bugged
	fontTableData, _ := io.ReadAll(fontTableStream)
	fontTable, err := tblfont.Load(bytes2.NewReader(fontTableData))

	if err != nil {
		return nil, err
	}

	result.FontTable = fontTable

	fontSpriteStream, err := loaderProvider.Load(fontPath + ".dc6")
	defer fontSpriteStream.Close()

	if err != nil {
		return nil, err
	}

	// hack: mpq block stream is bugged
	fontSpriteData, _ := io.ReadAll(fontSpriteStream)
	fontSprite, err := dc6.FromBytes(fontSpriteData)

	if err != nil {
		return nil, err
	}

	result.FontGfx = &common.DC6SequenceProvider{Sequences: fontSprite.Directions}

	result.RenderCallback = result.render
	result.UpdateCallback = result.update

	return result, nil
}

func (l *Label) render() {
	if !l.initialized || len(l.Caption) == 0 {
		return
	}

	tex := common.PaletteTexture[l.Palette]
	if !tex.Init {
		img := rl.NewImage(tex.Data, 256, 1, 1, rl.UncompressedR8g8b8a8)
		tex.Texture = rl.LoadTextureFromImage(img)

		tex.Init = true
	}

	posX, posY := l.GetPosition()

	rl.BeginShaderMode(common.PaletteShader)
	rl.SetShaderValueTexture(common.PaletteShader, common.PaletteShaderLoc, tex.Texture)
	rl.DrawTexture(l.texture, int32(posX), int32(posY), rl.White)
	rl.EndShaderMode()

}

func (l *Label) update() {
	if !l.initialized && len(l.Caption) > 0 {
		l.initialized = true
		l.initializeTexture()
	}
}

func (l *Label) initializeTexture() {
	width := 0
	height := 0

	charOffsets := make([]image.Point, len(l.Caption))

	for idx := range l.Caption {
		charOffsets[idx] = image.Point{X: width, Y: 0}
		glyph := l.FontTable.Glyphs[rune(l.Caption[idx])]
		width += glyph.Width()
		gHeight := glyph.Height()
		if gHeight > height {
			height = gHeight
		}
	}

	pixels := make([]byte, width*height)

	for idx := range l.Caption {
		glyph := l.FontTable.Glyphs[rune(l.Caption[idx])]
		frameIdx := glyph.FrameIndex()
		glyphWidth := glyph.Width()
		glyphHeight := glyph.Height()

		if glyphWidth == 0 || glyphHeight == 0 {
			continue
		}

		glyphOriginY := (l.FontGfx.FrameHeight(0, frameIdx) - glyphHeight) - 1

		for y := 0; y < glyphHeight; y++ {
			for x := 0; x < glyphWidth; x++ {
				c := l.FontGfx.GetColorIndexAt(0, frameIdx, x, y+glyphOriginY)
				idx := (charOffsets[idx].X + x) + ((charOffsets[idx].Y + y) * width)
				pixels[idx] = c
			}
		}
	}

	img := rl.NewImage(pixels, int32(width), int32(height), 1, rl.UncompressedGrayscale)

	if !l.hasTexture {
		l.hasTexture = true
	} else {
		rl.UnloadTexture(l.texture)
	}

	l.texture = rl.LoadTextureFromImage(img)
}
