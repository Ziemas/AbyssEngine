package engine

import (
	"image/color"
	"io/ioutil"

	"github.com/OpenDiablo2/AbyssEngine/common"
	pl2 "github.com/OpenDiablo2/pl2/pkg"
)

func (e *Engine) loadPalette(name string, path string) error {
	if common.PaletteTexture == nil {
		common.PaletteTexture = make(map[string]*common.PalTex)
	}

	paletteStream, err := e.loader.Load(path)

	if err != nil {
		return err
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

	colors = append(colors, transformToSlice(pal.BasePalette, pal.SelectedUnitShift)...)

	common.PaletteTextShiftOffset = len(colors)/(256*4)

	for idx := range pal.TextColorShifts {
		colors = append(colors, transformToSlice(pal.BasePalette, pal.TextColorShifts[idx])...)
	}

	for idx := range pal.HueVariations {
		colors = append(colors, transformToSlice(pal.BasePalette, pal.HueVariations[idx])...)
	}

	colors = append(colors, transformToSlice(pal.BasePalette, pal.RedTones)...)
	colors = append(colors, transformToSlice(pal.BasePalette, pal.BlueTones)...)
	colors = append(colors, transformToSlice(pal.BasePalette, pal.GreenTones)...)

	common.PaletteTransformsCount = len(colors)/(256*4)

	tex := &common.PalTex{}

	tex.Data = colors
	tex.Init = false

	common.PaletteTexture[name] = tex

	return nil
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
