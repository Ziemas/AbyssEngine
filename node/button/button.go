package button

import (
	"github.com/OpenDiablo2/AbyssEngine/common"
	"github.com/OpenDiablo2/AbyssEngine/loader"
	"github.com/OpenDiablo2/AbyssEngine/node"
	"github.com/OpenDiablo2/AbyssEngine/node/button/buttonlayout"
	"github.com/OpenDiablo2/AbyssEngine/node/sprite"
)

const (
	buttonStatePressed = iota + 1
	buttonStateToggled
	buttonStatePressedToggled
)

type Button struct {
	*node.Node

	buttonLayout buttonlayout.ButtonLayout
	enabled      bool
	pressed      bool
	toggled      bool
	onClick      func()
	sprite       *sprite.Sprite
}

func New(loaderProvider loader.LoaderProvider, mousePositionProvider common.MousePositionProvider,
	buttonLayout buttonlayout.ButtonLayout) (*Button, error) {
	result := &Button{
		buttonLayout: buttonLayout,
		enabled:      true,
		pressed:      false,
		toggled:      false,
	}

	result.RenderCallback = result.render
	result.UpdateCallback = result.update

	var err error

	result.sprite, err = sprite.New(loaderProvider, mousePositionProvider,
		buttonLayout.ResourceName, buttonLayout.PaletteName)

	if err != nil {
		return nil, err
	}

	result.sprite.CellSizeX = buttonLayout.XSegments
	result.sprite.CellSizeY = buttonLayout.YSegments
	err = result.AddChild(result.sprite.Node)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (b *Button) render() {

	if b.buttonLayout.HasImage {

		if !b.enabled {
			if b.toggled {
				b.sprite.CurrentFrame = b.buttonLayout.BaseFrame + buttonStateToggled
			} else {
				b.sprite.CurrentFrame = b.buttonLayout.DisabledFrame
			}
		} else if b.toggled && b.pressed {
			b.sprite.CurrentFrame = b.buttonLayout.BaseFrame + buttonStatePressedToggled
		} else if b.pressed && b.buttonLayout.AllowFrameChange {
			b.sprite.CurrentFrame = b.buttonLayout.BaseFrame + buttonStatePressed
		} else if b.toggled {
			b.sprite.CurrentFrame = b.buttonLayout.BaseFrame + buttonStateToggled
		} else {
			b.sprite.CurrentFrame = b.buttonLayout.BaseFrame
		}

		b.sprite.Render()
	}
}

func (b *Button) update(elapsed float64) {

}
