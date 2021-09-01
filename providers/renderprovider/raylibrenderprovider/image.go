package raylibrenderprovider

import rl "github.com/gen2brain/raylib-go/raylib"

type Image struct {
	src *rl.Image
}

func (i *Image) Width() int {
	return int(i.src.Width)
}

func (i *Image) Height() int {
	return int(i.src.Height)
}
