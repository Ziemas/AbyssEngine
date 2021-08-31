package common

import (
	ren "github.com/OpenDiablo2/AbyssEngine/renderer"
)

type PalTex struct {
	Texture ren.Texture
	Data    []byte
	Init    bool
}

//TODO: Yeah yeah, move this out
var (
	StandardShader         ren.Shader
	PaletteShader          ren.Shader
	PaletteShaderLoc       int32
	PaletteShaderOffsetLoc int32
	PaletteTexture         map[string]*PalTex
	PaletteTextShiftOffset int
	PaletteTransformsCount int
)
