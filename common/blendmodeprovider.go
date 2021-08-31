package common

import (
	"errors"
	"strings"
	ren "github.com/OpenDiablo2/AbyssEngine/renderer"
)

type BlendModeProvider interface {
	SetBlendMode(mode ren.BlendMode)
}

func BlendModeToString(mode ren.BlendMode) string {
	switch mode {
	case ren.BlendModeNone:
		return ""
	case ren.BlendModeAlpha:
		return "alpha"
	case ren.BlendModeAdditive:
		return "add"
	case ren.BlendModeMultiplied:
		return "multiply"
	case ren.BlendModeAddColors:
		return "addcolors"
	case ren.BlendModeSubtractColors:
		return "subcolors"
	default:
		return ""
	}
}

func StringToBlendMode(mode string) (ren.BlendMode, error) {
	switch strings.ToLower(mode) {
	case "":
		return ren.BlendModeNone, nil
	case "alpha":
		return ren.BlendModeAlpha, nil
	case "add":
		return ren.BlendModeAdditive, nil
	case "multiply":
		return ren.BlendModeMultiplied, nil
	case "addcolors":
		return ren.BlendModeAddColors, nil
	case "subcolors":
		return ren.BlendModeSubtractColors, nil
	default:
		return -1, errors.New("invalid blend mode")
	}
}
