package zrenderprovider

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

//go:embed shaders/fragment.glsl
var paletteFragmentShader string

//go:embed shaders/vertex.glsl
var standardVertexShader string

type Shader struct {
	id uint32
	UniformModel, UniformProjection,
	UniformImage, UniformPaletteTex,
	UniformUsePalette, UniformPaletteOffset int32
}

func InitShaders() Shader {
	shader := NewProgram(standardVertexShader, paletteFragmentShader)

	shader.UniformModel = GetShaderLocation(shader, "model")
	shader.UniformProjection = GetShaderLocation(shader, "projection")
	shader.UniformImage = GetShaderLocation(shader, "image")
	shader.UniformPaletteTex = GetShaderLocation(shader, "paletteTex")
	shader.UniformPaletteOffset = GetShaderLocation(shader, "paletteOffset")
	shader.UniformUsePalette = GetShaderLocation(shader, "usePalette")

	return shader
}

func GetShaderLocation(program Shader, location string) int32 {
	location += string(rune(0))
	string := make([]uint8, len(location))
	copy(string, location)
	return gl.GetUniformLocation(program.id, &string[0])
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

		debugPrint(5, fmt.Sprintf("Failed to link shader %s", msg))

		return Shader{id: 0}
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	s := Shader{id: program}
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

		debugPrint(5, fmt.Sprintf("Failed to compile shader %s", msg))

		return 0
	}

	return shader
}
