package renderer

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/rs/zerolog/log"
	"unsafe"
)

func debugCb(
	source uint32,
	gltype uint32,
	id uint32,
	severity uint32,
	length int32,
	message string,
	userParam unsafe.Pointer) {

	var _source string
	var _type string
	var _severity string

	switch source {
	case gl.DEBUG_SOURCE_API:
		_source = "API"
	case gl.DEBUG_SOURCE_WINDOW_SYSTEM:
		_source = "Window System"
	case gl.DEBUG_SOURCE_SHADER_COMPILER:
		_source = "Shader Compiler"
	case gl.DEBUG_SOURCE_THIRD_PARTY:
		_source = "Third Party"
	case gl.DEBUG_SOURCE_APPLICATION:
		_source = "Application"
	case gl.DEBUG_SOURCE_OTHER:
		fallthrough
	default:
		_source = "Unknown"
	}

	switch severity {
	case gl.DEBUG_SEVERITY_HIGH:
		_severity = "High"
	case gl.DEBUG_SEVERITY_MEDIUM:
		_severity = "Medium"
	case gl.DEBUG_SEVERITY_LOW:
		_severity = "Low"
	case gl.DEBUG_SEVERITY_NOTIFICATION:
		_severity = "Notification"
	default:
		_severity = "Unknown"

	}

	switch gltype {
	case gl.DEBUG_TYPE_ERROR:
		_type = "Error"
		log.Error().Msgf("Renderer: %s %s %s %s ", _source, _severity, _type, message)
	case gl.DEBUG_TYPE_DEPRECATED_BEHAVIOR:
		_type = "Deprecated"
		log.Warn().Msgf("Renderer: %s %s %s %s ", _source, _severity, _type, message)
	case gl.DEBUG_TYPE_UNDEFINED_BEHAVIOR:
		_type = "UB"
		log.Warn().Msgf("Renderer: %s %s %s %s ", _source, _severity, _type, message)
	case gl.DEBUG_TYPE_PORTABILITY:
		_type = "Portability"
		log.Warn().Msgf("Renderer: %s %s %s %s ", _source, _severity, _type, message)
	case gl.DEBUG_TYPE_PERFORMANCE:
		_type = "Perf"
		log.Info().Msgf("Renderer: %s %s %s %s ", _source, _severity, _type, message)
	case gl.DEBUG_TYPE_OTHER:
		_type = "Other"
		log.Info().Msgf("Renderer: %s %s %s %s ", _source, _severity, _type, message)
	case gl.DEBUG_TYPE_MARKER:
		_type = "Marker"
		log.Info().Msgf("Renderer: %s %s %s %s ", _source, _severity, _type, message)
	default:
		_type = "Unknown"
	}

}
