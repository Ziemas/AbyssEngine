package raylibrenderprovider

import (
	"github.com/OpenDiablo2/AbyssEngine/providers/renderprovider"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var mouseButtonLookup = map[renderprovider.MouseButton]int32{
	renderprovider.MouseButtonLeft:   rl.MouseLeftButton,
	renderprovider.MouseButtonRight:  rl.MouseRightButton,
	renderprovider.MouseButtonMiddle: rl.MouseMiddleButton,
}
