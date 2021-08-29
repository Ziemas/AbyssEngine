package button

import (
	"github.com/OpenDiablo2/AbyssEngine/common"
	"github.com/OpenDiablo2/AbyssEngine/node"
	"github.com/OpenDiablo2/AbyssEngine/node/button/buttonlayout"
	"github.com/OpenDiablo2/AbyssEngine/node/label"
	"github.com/OpenDiablo2/AbyssEngine/node/sprite"
	rl "github.com/gen2brain/raylib-go/raylib"
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
	label        *label.Label
	text         string
}

func New(loaderProvider common.LoaderProvider, mousePositionProvider common.MousePositionProvider,
	buttonLayout buttonlayout.ButtonLayout) (*Button, error) {
	result := &Button{
		Node:         node.New(),
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

	result.label, err = label.New(loaderProvider, buttonLayout.FontPath, buttonLayout.PaletteName)

	if err != nil {
		return nil, err
	}

	if buttonLayout.FixedWidth >= 0 {
		result.label.X = buttonLayout.FixedWidth / 2
	} else {
		width := result.sprite.Sequences.FrameWidth(result.sprite.CurrentSequence(), result.sprite.CurrentFrame, result.sprite.CellSizeX)
		result.label.X = width / 2
	}

	if buttonLayout.FixedHeight >= 0 {
		result.label.Y = buttonLayout.FixedHeight / 2
	} else {
		height := result.sprite.Sequences.FrameHeight(result.sprite.CurrentSequence(),
			result.sprite.CurrentFrame, result.sprite.CellSizeX, result.sprite.CellSizeY)

		result.label.Y = height / 2
	}

	result.label.X += buttonLayout.TextOffsetX
	result.label.Y += buttonLayout.TextOffsetY
	result.label.HAlign = label.LabelAlignCenter
	result.label.VAlign = label.LabelAlignCenter
	result.label.BlendMode = rl.BlendMultiplied

	err = result.sprite.Node.AddChild(result.label.Node)

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

func (b *Button) SetText(newText string) {
	if b.text == newText {
		return
	}

	b.text = newText
	b.label.Caption = newText
}

func (b *Button) render() {
	if !b.Visible || !b.Active {
		return
	}

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
	if !b.Active {
		return
	}

	b.sprite.Update(elapsed)
}
