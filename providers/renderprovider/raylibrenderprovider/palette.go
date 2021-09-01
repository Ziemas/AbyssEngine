package raylibrenderprovider

import (
	pl2 "github.com/OpenDiablo2/pl2/pkg"
	rl "github.com/gen2brain/raylib-go/raylib"
	"image/color"
)

type PalTex struct {
	Texture rl.Texture2D
	Data    []byte
	Init    bool
}

//TODO: Yeah yeah, move this out
var (
	PaletteShader          rl.Shader
	PaletteShaderLoc       int32
	PaletteShaderOffsetLoc int32
	PaletteTexture         map[string]*PalTex
	PaletteTextShiftOffset int
	PaletteTransformsCount int
)

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
