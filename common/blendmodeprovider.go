package common

import (
	"errors"
	"strings"
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

//var BlendModeLookup = map[BlendMode]rl.BlendMode{
//	BlendModeAlpha:          rl.BlendAlpha,
//	BlendModeAdditive:       rl.BlendAdditive,
//	BlendModeMultiplied:     rl.BlendMultiplied,
//	BlendModeAddColors:      rl.BlendAddColors,
//	BlendModeSubtractColors: rl.BlendSubtractColors,
//}

type BlendModeProvider interface {
	SetBlendMode(mode BlendMode)
}

func BlendModeToString(mode BlendMode) string {
	switch mode {
	case BlendModeNone:
		return ""
	case BlendModeAlpha:
		return "alpha"
	case BlendModeAdditive:
		return "add"
	case BlendModeMultiplied:
		return "multiply"
	case BlendModeAddColors:
		return "addcolors"
	case BlendModeSubtractColors:
		return "subcolors"
	default:
		return ""
	}
}

func StringToBlendMode(mode string) (BlendMode, error) {
	switch strings.ToLower(mode) {
	case "":
		return BlendModeNone, nil
	case "alpha":
		return BlendModeAlpha, nil
	case "add":
		return BlendModeAdditive, nil
	case "multiply":
		return BlendModeMultiplied, nil
	case "addcolors":
		return BlendModeAddColors, nil
	case "subcolors":
		return BlendModeSubtractColors, nil
	default:
		return -1, errors.New("invalid blend mode")
	}
}
