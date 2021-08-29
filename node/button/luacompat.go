package button

import (
	"fmt"

	"github.com/OpenDiablo2/AbyssEngine/common"
	lua "github.com/yuin/gopher-lua"
)

var luaTypeExportName = "button"
var LuaTypeExport = common.LuaTypeExport{
	Name: luaTypeExportName,
	//ConstructorFunc: newLuaEntity,
	Methods: map[string]lua.LGFunction{
		"node":    luaGetNode,
		"active":  luaGetSetActive,
		"enabled": luaGetSetEnabled,
		"pressed": luaGetSetPressed,
		"toggled": luaGetSetToggled,
	},
}

func luaGetSetToggled(l *lua.LState) int {
	button, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LBool(button.toggled))
		return 1
	}

	button.toggled = l.CheckBool(2)
	return 0
}

func luaGetSetPressed(l *lua.LState) int {
	button, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LBool(button.pressed))
		return 1
	}

	button.pressed = l.CheckBool(2)
	return 0
}

func luaGetSetEnabled(l *lua.LState) int {
	button, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LBool(button.enabled))
		return 1
	}

	button.enabled = l.CheckBool(2)
	return 0
}

func luaGetSetActive(l *lua.LState) int {
	button, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LBool(button.Active))
		return 1
	}

	button.Active = l.CheckBool(2)
	return 0
}

func (b *Button) ToLua(ls *lua.LState) *lua.LUserData {
	result := ls.NewUserData()
	result.Value = b

	ls.SetMetatable(result, ls.GetTypeMetatable(luaTypeExportName))

	return result
}

func FromLua(ud *lua.LUserData) (*Button, error) {
	v, ok := ud.Value.(*Button)

	if !ok {
		return nil, fmt.Errorf("failed to convert")
	}

	return v, nil
}

func luaGetNode(l *lua.LState) int {
	button, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	l.Push(button.Node.ToLua(l))

	return 1
}
