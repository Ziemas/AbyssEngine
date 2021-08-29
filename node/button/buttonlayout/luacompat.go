package buttonlayout

import (
	"fmt"
	"github.com/OpenDiablo2/AbyssEngine/common"
	lua "github.com/yuin/gopher-lua"
)

var luaTypeExportName = "buttonlayout"
var LuaTypeExport = common.LuaTypeExport{
	Name:            luaTypeExportName,
	ConstructorFunc: newLuaButtonLayout,
	Methods: map[string]lua.LGFunction{
		"resourceName":     luaGetSetResourceName,
		"paletteName":      luaGetSetPaletteName,
		"fontPath":         luaGetSetFontPath,
		"xSegments":        luaGetSetXSegments,
		"ySegments":        luaGetSetYSegments,
		"baseFrame":        luaGetSetBaseFrame,
		"disabledFrame":    luaGetSetDisabledFrame,
		"disabledColor":    luaGetSetDisabledColor,
		"textOffset":       luaGetSetTextOffset,
		"fixedWidth":       luaGetSetFixedWidth,
		"fixedHeight":      luaGetSetFixedHeight,
		"labelColor":       luaGetSetLabelColor,
		"toggleable":       luaGetSetToggleable,
		"allowFrameChange": luaGetSetAllowFrameChange,
		"hasImage":         luaGetSetHasImage,
		"tooltip":          luaGetSetTooltip,
		"tooltipXOffset":   luaGetSetTooltipXOffset,
		"tooltipYOffset":   luaGetSetTooltipYOffset,
	},
}

func luaGetSetResourceName(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LString(buttonLayout.ResourceName))
		return 1
	}

	newValue := l.CheckString(2)

	buttonLayout.ResourceName = newValue

	return 0
}

func luaGetSetPaletteName(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LString(buttonLayout.PaletteName))
		return 1
	}

	newValue := l.CheckString(2)

	buttonLayout.PaletteName = newValue

	return 0
}

func luaGetSetFontPath(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LString(buttonLayout.FontPath))
		return 1
	}

	newValue := l.CheckString(2)

	buttonLayout.FontPath = newValue

	return 0
}

func luaGetSetXSegments(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(buttonLayout.XSegments))
		return 1
	}

	newValue := l.CheckInt(2)

	buttonLayout.XSegments = newValue

	return 0
}

func luaGetSetYSegments(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(buttonLayout.YSegments))
		return 1
	}

	newValue := l.CheckInt(2)

	buttonLayout.YSegments = newValue

	return 0
}

func luaGetSetBaseFrame(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(buttonLayout.BaseFrame))
		return 1
	}

	newValue := l.CheckInt(2)

	buttonLayout.BaseFrame = newValue

	return 0
}

func luaGetSetDisabledFrame(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(buttonLayout.DisabledFrame))
		return 1
	}

	newValue := l.CheckInt(2)

	buttonLayout.DisabledFrame = newValue

	return 0
}

func luaGetSetDisabledColor(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(buttonLayout.DisabledColor))
		return 1
	}

	newValue := l.CheckInt64(2)

	buttonLayout.DisabledColor = uint32(newValue)

	return 0
}

func luaGetSetTextOffset(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(buttonLayout.TextOffsetX))
		l.Push(lua.LNumber(buttonLayout.TextOffsetY))
		return 1
	}

	newValueX := l.CheckInt(2)
	newValueY := l.CheckInt(3)

	buttonLayout.TextOffsetX = newValueX
	buttonLayout.TextOffsetY = newValueY

	return 0
}

func luaGetSetFixedWidth(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(buttonLayout.FixedWidth))
		return 1
	}

	newValue := l.CheckInt(2)

	buttonLayout.FixedWidth = newValue

	return 0
}

func luaGetSetFixedHeight(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(buttonLayout.FixedHeight))
		return 1
	}

	newValue := l.CheckInt(2)

	buttonLayout.FixedHeight = newValue

	return 0
}

func luaGetSetLabelColor(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(buttonLayout.LabelColor))
		return 1
	}

	newValue := l.CheckInt64(2)

	buttonLayout.LabelColor = uint32(newValue)

	return 0
}

func luaGetSetToggleable(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LBool(buttonLayout.Toggleable))
		return 1
	}

	newValue := l.CheckBool(2)

	buttonLayout.Toggleable = newValue

	return 0
}

func luaGetSetAllowFrameChange(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LBool(buttonLayout.AllowFrameChange))
		return 1
	}

	newValue := l.CheckBool(2)

	buttonLayout.AllowFrameChange = newValue

	return 0
}

func luaGetSetHasImage(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LBool(buttonLayout.HasImage))
		return 1
	}

	newValue := l.CheckBool(2)

	buttonLayout.HasImage = newValue

	return 0
}

func luaGetSetTooltip(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(buttonLayout.Tooltip))
		return 1
	}

	newValue := l.CheckInt(2)

	buttonLayout.Tooltip = newValue

	return 0
}

func luaGetSetTooltipXOffset(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(buttonLayout.TooltipXOffset))
		return 1
	}

	newValue := l.CheckInt(2)

	buttonLayout.TooltipXOffset = newValue

	return 0
}

func luaGetSetTooltipYOffset(l *lua.LState) int {
	buttonLayout, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(buttonLayout.TooltipYOffset))
		return 1
	}

	newValue := l.CheckInt(2)

	buttonLayout.TooltipYOffset = newValue

	return 0
}

func newLuaButtonLayout(l *lua.LState) int {
	result := &ButtonLayout{}
	userData := l.NewUserData()
	userData.Value = result

	l.SetMetatable(userData, l.GetTypeMetatable(luaTypeExportName))
	l.Push(userData)
	return 1
}

func (l *ButtonLayout) ToLua(ls *lua.LState) *lua.LUserData {
	result := ls.NewUserData()
	result.Value = l

	ls.SetMetatable(result, ls.GetTypeMetatable(luaTypeExportName))

	return result
}

func FromLua(ud *lua.LUserData) (*ButtonLayout, error) {
	v, ok := ud.Value.(*ButtonLayout)

	if !ok {
		return nil, fmt.Errorf("failed to convert")
	}

	return v, nil
}
