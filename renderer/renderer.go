package renderer

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/rs/zerolog/log"
)

type BlendMode int

const (
	BlendModeNone BlendMode = iota
	BlendModeAlpha
	BlendModeAdditive
	BlendModeMultiplied
	BlendModeAddColors
	BlendModeSubtractColors
)

func SetBlendMode(mode BlendMode) {
	switch mode {
	case BlendModeNone:
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		gl.BlendEquation(gl.FUNC_ADD)
	case BlendModeAdditive:
		gl.BlendFunc(gl.ONE, gl.ONE)
		gl.BlendEquation(gl.FUNC_ADD)
	case BlendModeMultiplied:
		gl.BlendFunc(gl.DST_COLOR, gl.ONE_MINUS_SRC_ALPHA)
		gl.BlendEquation(gl.FUNC_ADD)
	default:
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
		gl.BlendEquation(gl.FUNC_ADD)
	}
}

type Renderer struct {
	window *glfw.Window
}

var (
	projection = mgl32.Ortho(0.0, 800.0, 600.0, 0.0, -1.0, 1.0)
	quadVAO    uint32
	vbo        uint32

	verts = []float32{
		0.0, 1.0, 0.0, 1.0,
		1.0, 0.0, 1.0, 0.0,
		0.0, 0.0, 0.0, 0.0,

		0.0, 1.0, 0.0, 1.0,
		1.0, 1.0, 1.0, 1.0,
		1.0, 0.0, 1.0, 0.0,
	}

	Width, Height int
)

func Init() *glfw.Window {
	log.Info().Msgf("Initialising OpenGL")

	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	window, err := glfw.CreateWindow(800, 600, "Abyss Engine", nil, nil)
	if err != nil {
		panic(err)
	}

	window.SetAspectRatio(4, 3)
	window.SetFramebufferSizeCallback(framebufferSizeCallback)
	window.SetCursorPosCallback(cursorPosCallback)
	window.SetMouseButtonCallback(mouseButtonCallback)

	window.MakeContextCurrent()
	window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	gl.DebugMessageCallback(debugCb, nil)
	gl.Enable(gl.DEBUG_OUTPUT)

	//glfw.SwapInterval(0)

	gl.GenVertexArrays(1, &quadVAO)
	gl.GenBuffers(1, &vbo)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)
	gl.BindVertexArray(quadVAO)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 4*4, nil)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	gl.ClearColor(0, 0, 0, 255)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	return window
}

func framebufferSizeCallback(w *glfw.Window, width int, height int) {
	//rX := float32(width) / 800.0
	//rY := float32(height) / 600.0
	projection = mgl32.Ortho(0, float32(width), float32(height), 0, -1.0, 1.0)
	gl.Viewport(0, 0, int32(width), int32(height))

	//var r float32
	//if rX < rY {
	//	r = rX
	//} else {
	//	r = rY
	//}

	//vW := int(r * 800)
	//vH := int(r * 600)

	//vX := int((float32(width)-800.0*r) / 2)
	//vY := int((float32(width)-600.0*r) / 2)

	//gl.Viewport(int32(vX), int32(vY), int32(vW), int32(vH))

	Width = width
	Height = height
}

func GetScreenWidth() int {
	return Width
}

func GetScreenHeight() int {
	return Height
}

func Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func GetShaderLocation(program Shader, location string) int32 {
	location += string(rune(0))
	string := make([]uint8, len(location))
	copy(string, location)
	return gl.GetUniformLocation(program.id, &string[0])
}
