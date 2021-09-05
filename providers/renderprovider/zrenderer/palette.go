package zrenderprovider

import (
	"image/color"
	"io"
	"io/ioutil"

	pl2 "github.com/OpenDiablo2/pl2/pkg"
)

var (
	PaletteTextShiftOffset int
	PaletteTransformsCount int
)

type Palette struct {
	Texture *zTexture
	Data    []byte
	Init    bool
}

func NewPalette(paletteStream io.Reader) (*Palette, error) {
	paletteBytes, err := ioutil.ReadAll(paletteStream)

	if err != nil {
		return nil, err
	}

	pal, err := pl2.FromBytes(paletteBytes)

	if err != nil {
		return nil, err
	}

	colors := make([]uint8, 0)
	colors = append(colors, palToSlice(pal.BasePalette)...)

	PaletteTextShiftOffset = len(colors) / (256 * 4)

	for idx := range pal.TextColorShifts {
		colors = append(colors, transformToSlice(pal.BasePalette, pal.TextColorShifts[idx])...)
	}

	PaletteTransformsCount = len(colors) / (256 * 4)

	tex := &Palette{}

	tex.Data = colors
	tex.Init = false

	return tex, nil
}

func transformToSlice(palette color.Palette, transform pl2.Transform) []uint8 {
	colors := make([]uint8, 256*4)
	for i := 0; i < 256; i++ {
		offset := i * 4
		r, g, b, _ := palette[transform[i]].RGBA()

		colors[offset] = uint8(r >> 8)
		colors[offset+1] = uint8(g >> 8)
		colors[offset+2] = uint8(b >> 8)
		colors[offset+3] = 255
	}

	colors[3] = 0

	return colors
}

func palToSlice(color color.Palette) []uint8 {
	colors := make([]uint8, 256*4)
	for i := 0; i < 256; i++ {
		if i >= len(color) {
			break
		}

		offset := i * 4
		r, g, b, _ := color[i].RGBA()
		colors[offset] = uint8(r >> 8)
		colors[offset+1] = uint8(g >> 8)
		colors[offset+2] = uint8(b >> 8)
		colors[offset+3] = 255
	}

	colors[3] = 0

	return colors
}
