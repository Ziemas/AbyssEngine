package renderer

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/rs/zerolog/log"
	"strings"
)

const (
	UniformModelLoc      = 1
	UniformProjectionLoc = 2
	UniformImageLoc      = 3

	UniformPaletteTexLoc    = 4
	UniformPaletteOffsetLoc = 5

	TextureImageBind   = 0
	TexturePaletteBind = 1
)

type Shader struct {
	id uint32
}

func BeginShader(s Shader) {
	gl.UseProgram(s.id)
	gl.UniformMatrix4fv(int32(UniformProjectionLoc), 1, false, &projection[0])
}

func EndShader() {
	gl.UseProgram(0)

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
