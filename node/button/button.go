package button

import (
	"github.com/OpenDiablo2/AbyssEngine/common"
	"github.com/OpenDiablo2/AbyssEngine/node"
	"github.com/OpenDiablo2/AbyssEngine/node/button/buttonlayout"
	"github.com/OpenDiablo2/AbyssEngine/node/label"
	"github.com/OpenDiablo2/AbyssEngine/node/sprite"
	"github.com/OpenDiablo2/AbyssEngine/providers/renderprovider"
)

const (
	buttonStatePressed = iota + 1
	buttonStateToggled
	buttonStatePressedToggled
)

type Button struct {
	*node.Node

	renderProvider renderprovider.RenderProvider
	buttonLayout   buttonlayout.ButtonLayout
	enabled        bool
	pressed        bool
	toggled        bool
	labelOffset    bool
	mouseOver      bool
	onClick        func()
	sprite         *sprite.Sprite
	label          *label.Label
	text           string
}

func New(loaderProvider common.LoaderProvider, renderProvider renderprovider.RenderProvider, mousePositionProvider common.MousePositionProvider,
	buttonLayout buttonlayout.ButtonLayout) (*Button, error) {
	result := &Button{
		Node:           node.New(),
		buttonLayout:   buttonLayout,
		renderProvider: renderProvider,
		enabled:        true,
		pressed:        false,
		toggled:        false,
	}

	result.RenderCallback = result.render
	result.UpdateCallback = result.update

	var err error

	result.sprite, err = sprite.New(loaderProvider, mousePositionProvider, renderProvider,
		buttonLayout.ResourceName, buttonLayout.PaletteName)

	if err != nil {
		return nil, err
	}

	result.label, err = label.New(loaderProvider, renderProvider, buttonLayout.FontPath, buttonLayout.PaletteName)

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
	result.label.BlendMode = renderprovider.BlendModeMultiplied

	err = result.sprite.Node.AddChild(result.label.Node)

	if err != nil {
		return nil, err
	}

	result.sprite.CellSizeX = buttonLayout.XSegments
	result.sprite.CellSizeY = buttonLayout.YSegments
	result.sprite.OnMouseButtonDown = func() { result.onPressed() }
	result.sprite.OnMouseButtonUp = func() { result.onReleased() }
	result.sprite.OnMouseOver = func() { result.onHover() }
	result.sprite.OnMouseLeave = func() { result.onLeave() }

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
		b.sprite.Render()
	}
}

func (b *Button) update(elapsed float64) {
	if !b.Active {
		return
	}

	if b.pressed && !b.mouseOver && !b.renderProvider.IsMouseButtonPressed(renderprovider.MouseButtonLeft) {
		b.pressed = false
	}

	pressed := b.pressed && b.mouseOver

	if b.buttonLayout.HasImage {
		if !b.enabled {
			if b.toggled {
				b.sprite.CurrentFrame = b.buttonLayout.BaseFrame + buttonStateToggled
			} else {
				b.sprite.CurrentFrame = b.buttonLayout.DisabledFrame
			}
		} else if b.toggled && pressed {
			b.sprite.CurrentFrame = b.buttonLayout.BaseFrame + (buttonStatePressedToggled * b.sprite.CellSizeX)
		} else if pressed && b.buttonLayout.AllowFrameChange {
			b.sprite.CurrentFrame = b.buttonLayout.BaseFrame + (buttonStatePressed * b.sprite.CellSizeX)
		} else if b.toggled {
			b.sprite.CurrentFrame = b.buttonLayout.BaseFrame + (buttonStateToggled * b.sprite.CellSizeX)
		} else {
			b.sprite.CurrentFrame = b.buttonLayout.BaseFrame
		}

		if b.labelOffset && !pressed {
			b.label.X += 2
			b.label.Y -= 2
			b.labelOffset = false
		} else if !b.labelOffset && pressed {
			b.label.X -= 2
			b.label.Y += 2
			b.labelOffset = true
		}
	}

	b.sprite.Update(elapsed)
}

func (b *Button) onPressed() {
	if !b.enabled || !b.Active {
		return
	}

	b.pressed = true
}

func (b *Button) onReleased() {
	if !b.enabled || !b.Active {
		return
	}

	if b.pressed {
		b.pressed = false
		if b.onClick != nil {
			b.onClick()
		}
	}
}

func (b *Button) onHover() {
	if !b.enabled || !b.Active {
		return
	}

	b.mouseOver = true
}

func (b *Button) onLeave() {
	if !b.enabled || !b.Active {
		return
	}

	b.mouseOver = false
}
