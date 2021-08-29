package engine

import (
	"io/ioutil"

	"github.com/OpenDiablo2/AbyssEngine/common"
	datPalette "github.com/OpenDiablo2/dat_palette/pkg"
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

	paletteData, err := datPalette.Decode(paletteBytes)

	if err != nil {
		return err
	}

	colors := make([]byte, 256*4)

	for i := 0; i < 256; i++ {
		if i >= len(paletteData) {
			break
		}

		offset := i * 4
		r, g, b, _ := paletteData[i].RGBA()
		colors[offset] = uint8(r >> 8)
		colors[offset+1] = uint8(g >> 8)
		colors[offset+2] = uint8(b >> 8)
		colors[offset+3] = 255
	}

	colors[3] = 0

	tex := &common.PalTex{}

	tex.Data = colors
	tex.Init = false

	common.PaletteTexture[name] = tex

	return nil
}
