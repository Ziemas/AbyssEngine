package zrenderprovider

import (
	"image"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/pkg/errors"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// TODO use a texture atlas

type zChar struct {
	Size    image.Rectangle
	Bearing image.Point
	Advance fixed.Int26_6
	Texture *zTexture
}

type zFont struct {
	characters map[rune]zChar
	face       font.Face
	size       int
}

func NewFont(fontBytes []byte, size int) (*zFont, error) {
	result := &zFont{
		characters: map[rune]zChar{},
		size:       size,
	}

	fontData, err := opentype.Parse(fontBytes)
	if err != nil {
		debugPrint(4, err.Error())
		return nil, err
	}

	opts := &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	}

	face, err := opentype.NewFace(fontData, opts)
	if err != nil {
		return nil, err
	}

	result.face = face

	for c := 0; c < 128; c++ {
		dr, mask, _, advance, ok := face.Glyph(fixed.Point26_6{}, rune(c))
		if !ok {
			return nil, errors.New("Failed to read glyph from font")
		}

		data := image.NewRGBA(dr)
		draw.Draw(data, dr, mask, image.ZP, draw.Src)

		character := zChar{
			Size:    dr,
			Bearing: dr.Min,
			Advance: advance,
		}

		if dr.Bounds().Dx() > 0 && dr.Bounds().Dy() > 0 {
			tex, err := NewTexture(gl.Ptr(data.Pix), data.Bounds().Dx(), data.Bounds().Dy(), PixelFmtRGBA8)
			if err != nil {
				return nil, err
			}

			character.Texture = tex
		}

		result.characters[rune(c)] = character
	}

	return result, nil
}

func (f *zFont) RenderString(x, y int, text string, r *Renderer) {
	var xOffset fixed.Int26_6
	prev := rune(-1)
	for _, c := range text {
		if prev > -1 {
			xOffset += f.face.Kern(prev, c)
		}

		char := f.characters[c]

		if char.Texture != nil {
			dstX := x + xOffset.Round() + char.Bearing.X
			dstY := (y + f.size) + char.Bearing.Y
			r.DrawTexture(char.Texture, dstX, dstY, "", 0)
		}

		xOffset += char.Advance
	}

}
