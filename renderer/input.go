package renderer

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	cursorX, cursorY    float64
	leftDown, rightDown bool
)

func cursorPosCallback(w *glfw.Window, x, y float64) {
	cursorX = x
	cursorY = y
}

func mouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft {
		if action == glfw.Press {
			leftDown = true
		} else {
			leftDown = false
		}
	}

	if button == glfw.MouseButtonRight {
		if action == glfw.Press {
			rightDown = true
		} else {
			rightDown = false
		}
	}
}

func GetMouseX() float64 {
	return cursorX
}

func GetMouseY() float64 {
	return cursorY
}

type MouseButton int
const (
	MouseLeftButton MouseButton = iota
	MouseRightButton
)

func IsMouseButtonDown(button MouseButton) bool {
	switch button {
	case MouseLeftButton:
		return leftDown
	case MouseRightButton:
		return rightDown
	default:
		return false
	}
}
