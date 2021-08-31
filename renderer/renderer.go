package renderer

import (
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/rs/zerolog/log"
)

type Renderer struct {
	window *glfw.Window
}

var (
	projection = mgl32.Ortho(0.0, 800.0, 600.0, 0.0, -1.0, 1.0)
	quadVAO    uint32
	vbo        uint32

	// Our private uniforms for the renderer
	curModelUni uint32
	curProjUni  uint32

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

//func New(window *glfw.Window) *Renderer {
//	ren := Renderer{
//		window: window,
//	}
//
//	err := gl.Init()
//	if err != nil {
//		panic(err)
//	}
//
//	return &ren
//}

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

	window.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	gl.DebugMessageCallback(debugCb, nil)
	gl.Enable(gl.DEBUG_OUTPUT)

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

func BeginShader(s Shader) {
	gl.UseProgram(s.id)

	// Track these so we can use them without mucking about with useprogram
	curModelUni = s.modelUni
	curProjUni = s.projectionUni

	gl.UniformMatrix4fv(int32(curProjUni), 1, false, &projection[0])
}

func EndShader() {
	gl.UseProgram(0)

}

type Shader struct {
	id            uint32
	modelUni      uint32
	projectionUni uint32
}

func NewProgram(vertexShaderSrc, fragmentShaderSrc string) Shader {
	vertexShader := compileShader(vertexShaderSrc, gl.VERTEX_SHADER)
	fragmentShader := compileShader(fragmentShaderSrc, gl.FRAGMENT_SHADER)

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		msg := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(msg))

		log.Error().Msgf("Failed to compile shader %s", msg)

		return Shader{id: 0}
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	s := Shader{id: program}
	s.modelUni = uint32(GetShaderLocation(s, "model"))
	s.projectionUni = uint32(GetShaderLocation(s, "projection"))
	return s
}

func compileShader(source string, sType uint32) uint32 {
	shader := gl.CreateShader(sType)

	csource, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csource, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		msg := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(msg))

		log.Error().Msgf("Failed to compile shader %s", msg)

		return 0
	}

	return shader
}

func SetShaderValuei(location int32, value int32) {
	gl.Uniform1i(location, value)
}

func SetShaderValueF(location int32, value float32) {
	gl.Uniform1f(location, value)
}
