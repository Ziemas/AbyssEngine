package raylibrenderprovider

import rl "github.com/gen2brain/raylib-go/raylib"

type Texture struct {
	tex             rl.Texture2D
	rtex            rl.RenderTexture2D
	isRenderTexture bool
}

func (t *Texture) ID() int {
	if t.isRenderTexture {
		return int(t.rtex.ID)
	}
	return int(t.tex.ID)
}

func (t *Texture) Width() int {
	if t.isRenderTexture {
		return int(t.rtex.Texture.Width)
	}
	return int(t.tex.Width)
}

func (t *Texture) Height() int {
	if t.isRenderTexture {
		return int(t.rtex.Texture.Height)
	}
	return int(t.tex.Height)
}
