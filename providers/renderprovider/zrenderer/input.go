package zrenderprovider

import "github.com/go-gl/glfw/v3.3/glfw"

var (
	cursorX, cursorY                float64
	leftDown, middleDown, rightDown bool
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

	if button == glfw.MouseButtonMiddle {
		if action == glfw.Press {
			middleDown = true
		} else {
			middleDown = false
		}
	}
}
