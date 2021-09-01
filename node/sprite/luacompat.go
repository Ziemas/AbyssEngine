package sprite

import (
	"errors"
	"fmt"
	"github.com/OpenDiablo2/AbyssEngine/providers/renderprovider"
	"strings"

	"github.com/OpenDiablo2/AbyssEngine/common"
	lua "github.com/yuin/gopher-lua"
)

var luaTypeExportName = "sprite"
var LuaTypeExport = common.LuaTypeExport{
	Name: luaTypeExportName,
	//ConstructorFunc: newLuaEntity,
	Methods: map[string]lua.LGFunction{
		"node":                   luaGetNode,
		"cellSize":               luaGetSetCellSize,
		"active":                 luaGetSetActive,
		"bottomOrigin":           luaGetSetBottomOrigin,
		"visible":                luaGetSetVisible,
		"position":               luaGetSetPosition,
		"currentSequence":        luaGetSetCurrentSequence,
		"currentFrame":           luaGetSetCurrentFrame,
		"sequenceCount":          luaGetSequenceCount,
		"frameCount":             luaGetFrameCount,
		"destroy":                luaDestroy,
		"mouseButtonDownHandler": luaGetSetMouseButtonDownHandler,
		"mouseButtonUpHandler":   luaGetSetMouseButtonUpHandler,
		"mouseOverHandler":       luaGetSetMouseOverHandler,
		"mouseLeaveHandler":      luaGetSetMouseLeaveHandler,
		"playForward":            luaGetPlayForward,
		"blendMode":              luaGetSetBlendMode,
	},
}

func BlendModeToString(mode renderprovider.BlendMode) string {
	switch mode {
	case renderprovider.BlendModeNone:
		return ""
	case renderprovider.BlendModeAlpha:
		return "alpha"
	case renderprovider.BlendModeAdditive:
		return "add"
	case renderprovider.BlendModeMultiplied:
		return "multiply"
	case renderprovider.BlendModeAddColors:
		return "addcolors"
	case renderprovider.BlendModeSubtractColors:
		return "subcolors"
	default:
		return ""
	}
}

func StringToBlendMode(mode string) (renderprovider.BlendMode, error) {
	switch strings.ToLower(mode) {
	case "":
		return renderprovider.BlendModeNone, nil
	case "alpha":
		return renderprovider.BlendModeAlpha, nil
	case "add":
		return renderprovider.BlendModeAdditive, nil
	case "multiply":
		return renderprovider.BlendModeMultiplied, nil
	case "addcolors":
		return renderprovider.BlendModeAddColors, nil
	case "subcolors":
		return renderprovider.BlendModeSubtractColors, nil
	default:
		return -1, errors.New("invalid blend mode")
	}
}

func luaGetSetBottomOrigin(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LBool(sprite.bottomOrigin))
		return 1
	}

	newValue := l.CheckBool(2)
	sprite.bottomOrigin = newValue

	return 0
}

func luaGetSetBlendMode(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LString(BlendModeToString(sprite.blendMode)))
		return 1
	}

	newMode, err := StringToBlendMode(l.CheckString(2))

	if err != nil {
		l.RaiseError(err.Error())
		return 0
	}

	sprite.blendMode = newMode

	return 0
}

func luaGetPlayForward(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	sprite.PlayForward()

	return 0
}

func luaGetSetMouseOverHandler(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(l.NewFunction(func(l *lua.LState) int {
			sprite.OnMouseOver()
			return 0
		}))

		return 1
	}

	luaFunc := l.CheckFunction(2)
	sprite.OnMouseOver = func() {
		go func() {
			if err := l.CallByParam(lua.P{
				Fn:      luaFunc,
				NRet:    1,
				Protect: true,
			}, sprite.ToLua(l)); err != nil {
				panic(err)
			}
		}()
	}

	return 0
}

func luaGetSetMouseLeaveHandler(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(l.NewFunction(func(l *lua.LState) int {
			sprite.OnMouseLeave()
			return 0
		}))

		return 1
	}

	luaFunc := l.CheckFunction(2)
	sprite.OnMouseLeave = func() {
		go func() {
			if err := l.CallByParam(lua.P{
				Fn:      luaFunc,
				NRet:    1,
				Protect: true,
			}, sprite.ToLua(l)); err != nil {
				panic(err)
			}
		}()
	}

	return 0
}

func luaDestroy(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	sprite.Destroy()

	return 0
}

func luaGetSetMouseButtonUpHandler(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(l.NewFunction(func(l *lua.LState) int {
			sprite.OnMouseButtonUp()
			return 0
		}))

		return 1
	}

	luaFunc := l.CheckFunction(2)
	sprite.OnMouseButtonUp = func() {
		go func() {
			if err := l.CallByParam(lua.P{
				Fn:      luaFunc,
				NRet:    1,
				Protect: true,
			}, sprite.ToLua(l)); err != nil {
				panic(err)
			}
		}()
	}

	return 0
}

func luaGetFrameCount(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	l.Push(lua.LNumber(sprite.Sequences.FrameCount(sprite.CurrentSequence())))
	return 1
}

func luaGetSequenceCount(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	l.Push(lua.LNumber(sprite.Sequences.SequenceCount()))
	return 1
}

func luaGetSetCurrentFrame(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(sprite.CurrentFrame))
		return 1
	}

	newFrame := l.CheckInt(2)

	if (newFrame < 0) || (newFrame >= sprite.Sequences.FrameCount(sprite.CurrentSequence())) {
		l.RaiseError("frame index out of bounds")
		return 0
	}

	sprite.CurrentFrame = newFrame

	return 0
}

func luaGetSetCurrentSequence(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(sprite.CurrentSequence()))
		return 1
	}

	newSequence := l.CheckInt(2)

	if (newSequence < 0) || (newSequence >= sprite.Sequences.SequenceCount()) {
		l.RaiseError("sequence index out of bounds")
		return 0
	}

	sprite.SetSequence(newSequence)
	sprite.CurrentFrame = 0

	return 0
}

func luaGetNode(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	l.Push(sprite.Node.ToLua(l))

	return 1
}

func (s *Sprite) ToLua(l *lua.LState) *lua.LUserData {
	result := l.NewUserData()
	result.Value = s

	l.SetMetatable(result, l.GetTypeMetatable(luaTypeExportName))

	return result
}

func FromLua(ud *lua.LUserData) (*Sprite, error) {
	v, ok := ud.Value.(*Sprite)

	if !ok {
		return nil, fmt.Errorf("failed to convert")
	}

	return v, nil
}

func luaGetSetPosition(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(sprite.X))
		l.Push(lua.LNumber(sprite.Y))
		return 2
	}

	posX := l.ToNumber(2)
	posY := l.ToNumber(3)

	sprite.X = int(posX)
	sprite.Y = int(posY)

	return 0
}

func luaGetSetCellSize(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(sprite.CellSizeX))
		l.Push(lua.LNumber(sprite.CellSizeY))
		return 2
	}

	sizeX := l.ToNumber(2)
	sizeY := l.ToNumber(3)

	sprite.CellSizeX = int(sizeX)
	sprite.CellSizeY = int(sizeY)

	return 0
}

func luaGetSetMouseButtonDownHandler(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(l.NewFunction(func(l *lua.LState) int {
			sprite.OnMouseButtonDown()
			return 0
		}))

		return 1
	}

	luaFunc := l.CheckFunction(2)
	sprite.OnMouseButtonDown = func() {
		go func() {
			if err := l.CallByParam(lua.P{
				Fn:      luaFunc,
				NRet:    1,
				Protect: true,
			}, sprite.ToLua(l)); err != nil {
				panic(err)
			}
		}()
	}

	return 0
}

func luaGetSetVisible(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LBool(sprite.Visible))
		return 1
	}

	newValue := l.CheckBool(2)
	sprite.Visible = newValue

	return 0
}

func luaGetSetActive(l *lua.LState) int {
	sprite, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LBool(sprite.Active))
		return 1
	}

	newValue := l.CheckBool(2)
	sprite.Active = newValue

	return 0
}
