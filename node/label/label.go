package label

import (
	bytes "bytes"
	"errors"
	"github.com/OpenDiablo2/AbyssEngine/providers/renderprovider"
	"image"
	"io"
	"strings"

	"github.com/OpenDiablo2/AbyssEngine/common"
	"github.com/OpenDiablo2/AbyssEngine/node"
	dc6 "github.com/OpenDiablo2/dc6/pkg"
	tblfont "github.com/OpenDiablo2/tbl_font/pkg"
)

type LabelAlign int

const (
	LabelAlignStart LabelAlign = iota
	LabelAlignCenter
	LabelAlignEnd
)

func (l LabelAlign) ToString() string {
	switch l {
	case LabelAlignStart:
		return "start"
	case LabelAlignCenter:
		return "center"
	case LabelAlignEnd:
		return "end"
	}

	return "start"
}

func StringToLabelAlign(s string) (LabelAlign, error) {
	switch strings.ToLower(s) {
	case "start":
		return LabelAlignStart, nil
	case "center":
		return LabelAlignCenter, nil
	case "end":
		return LabelAlignEnd, nil
	}

	return LabelAlignStart, errors.New("unknown alignment value")
}

type Label struct {
	*node.Node

	initialized    bool
	hasTexture     bool
	renderProvider renderprovider.RenderProvider
	texture        renderprovider.Texture
	FontTable      *tblfont.FontTable
	FontGfx        common.SequenceProvider
	Palette        string
	Caption        string
	BlendMode      renderprovider.BlendMode
	color          int
	HAlign         LabelAlign
	VAlign         LabelAlign
}

func New(
	loaderProvider common.LoaderProvider,
	renderProvider renderprovider.RenderProvider,
	fontPath, palette string) (*Label, error) {
	result := &Label{
		Node:           node.New(),
		renderProvider: renderProvider,
		initialized:    false,
		HAlign:         LabelAlignStart,
		VAlign:         LabelAlignStart,
		BlendMode:      renderprovider.BlendModeNone,
		color:          7,
	}

	result.Palette = palette

	fontTableStream, err := loaderProvider.Load(fontPath + ".tbl")
	defer func() { _ = fontTableStream.Close() }()

	if err != nil {
		return nil, err
	}

	// hack: mpq block stream is bugged
	fontTableData, _ := io.ReadAll(fontTableStream)
	fontTable, err := tblfont.Load(bytes.NewReader(fontTableData))

	if err != nil {
		return nil, err
	}

	result.FontTable = fontTable

	fontSpriteStream, err := loaderProvider.Load(fontPath + ".dc6")
	defer func() { _ = fontSpriteStream.Close() }()

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

	posX, posY := l.GetPosition()

	switch l.HAlign {
	case LabelAlignCenter:
		posX -= l.texture.Width() / 2
	case LabelAlignEnd:
		posX -= l.texture.Width()
	}

	switch l.VAlign {
	case LabelAlignCenter:
		posY -= l.texture.Height() / 2
	case LabelAlignEnd:
		posY -= l.texture.Height()
	}

	l.renderProvider.BeginBlendMode(l.BlendMode)

	_ = l.renderProvider.DrawFontTexture(l.texture, posX, posY, l.Palette, l.color)
	l.renderProvider.EndBlendMode()

}

func (l *Label) update(elapsed float64) {
	if !l.initialized && len(l.Caption) > 0 {
		l.initialized = true
		l.initializeTexture()
	}
}

func (l *Label) getTextMetrics() (width, height int) {
	var (
		lineWidth  int
		lineHeight int
	)

	for _, c := range l.Caption {
		if c == '\n' {
			width = common.MaxInt(width, lineWidth)
			height += lineHeight
			lineWidth = 0
			lineHeight = 0
		} else {
			glyph := l.FontTable.Glyphs[c]
			lineWidth += glyph.Width()
			lineHeight = common.MaxInt(lineHeight, glyph.Height())
		}
	}

	width = common.MaxInt(width, lineWidth)
	height += lineHeight

	return width, height
}

func (l *Label) initializeTexture() {
	charOffsets := make([]image.Point, len(l.Caption))
	lineHeights := make([]int, 0)
	tw := 0
	th := 0
	width := 0
	height := 0
	lineHeight := 0
	for idx, c := range l.Caption {
		charOffsets[idx] = image.Point{X: tw, Y: th}
		glyph := l.FontTable.Glyphs[c]
		glyphWidth := l.FontTable.Glyphs[c].Width()
		lineHeight = common.MaxInt(lineHeight, l.FontGfx.FrameHeight(0, glyph.FrameIndex(), 1, 1))
		width = common.MaxInt(width, glyphWidth+tw)
		if l.Caption[idx] == '\n' {
			height += lineHeight
			lineHeights = append(lineHeights, lineHeight)
			tw = 0
			th += lineHeight
			lineHeight = 0
			continue
		}
		tw += glyphWidth
	}
	lineHeights = append(lineHeights, lineHeight)
	height += lineHeight

	pixels := make([]byte, width*height)

	curLine := 0
	for idx := range l.Caption {
		if l.Caption[idx] == '\n' {
			curLine++
			continue
		}
		glyph := l.FontTable.Glyphs[rune(l.Caption[idx])]
		frameIdx := glyph.FrameIndex()
		glyphWidth := glyph.Width()
		glyphHeight := glyph.Height()

		if glyphWidth == 0 || glyphHeight == 0 {
			continue
		}
		for y := 0; y < l.FontGfx.FrameHeight(0, frameIdx, 1, 1); y++ {
			for x := 0; x < l.FontGfx.FrameWidth(0, frameIdx, 1); x++ {
				if x >= glyphWidth {
					break
				}
				c := l.FontGfx.GetColorIndexAt(0, frameIdx, x, y)
				tx := charOffsets[idx].X + x
				ty := charOffsets[idx].Y + y

				if tx < 0 || tx >= width || ty < 0 || ty >= height {
					continue
				}

				idx := tx + (ty * width)
				pixels[idx] = c
			}
		}
	}

	img, _ := l.renderProvider.NewImage(bytes.NewReader(pixels), width, height, renderprovider.ImageColorModeGrayscale)

	if !l.hasTexture {
		l.hasTexture = true
	} else {
		_ = l.renderProvider.FreeTexture(l.texture)
	}

	l.texture, _ = l.renderProvider.LoadTextureFromImage(img)
}
