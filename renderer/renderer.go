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


func QueueDraw(drawable Renderable) {
	drawList.PushBack(drawable)
}

func framebufferSizeCallback(w *glfw.Window, width int, height int) {
	ratio := 4.0 / 3.0
	vW := width
	vH := int(float64(width) / ratio)

	if (vH > height) {
		vH = height
		vW = int(float64(height) * ratio)
	}

	vX := (width - vW) / 2
	vY := (height - vH) / 2

	gl.Viewport(int32(vX), int32(vY), int32(vW), int32(vH))

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
