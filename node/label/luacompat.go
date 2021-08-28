package label

import (
	"fmt"

	"github.com/OpenDiablo2/AbyssEngine/common"
	lua "github.com/yuin/gopher-lua"
)

var luaTypeExportName = "label"
var LuaTypeExport = common.LuaTypeExport{
	Name: luaTypeExportName,
	//ConstructorFunc: newLuaEntity,
	Methods: map[string]lua.LGFunction{
		"node":     luaGetNode,
		"caption":  luaGetSetCaption,
		"position": luaGetSetPosition,
	},
}

func luaGetSetPosition(l *lua.LState) int {
	label, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LNumber(label.X))
		l.Push(lua.LNumber(label.Y))
		return 2
	}

	posX := l.ToNumber(2)
	posY := l.ToNumber(3)

	label.X = int(posX)
	label.Y = int(posY)

	return 0
}

func luaGetSetCaption(l *lua.LState) int {
	label, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	if l.GetTop() == 1 {
		l.Push(lua.LString(label.Caption))
		return 1
	}

	newCaption := l.CheckString(2)

	if label.Caption == newCaption {
		return 0
	}

	label.Caption = newCaption
	label.initialized = false

	return 0
}

func luaGetNode(l *lua.LState) int {
	label, err := FromLua(l.ToUserData(1))

	if err != nil {
		l.RaiseError("failed to convert")
		return 0
	}

	l.Push(label.Node.ToLua(l))

	return 1
}

func (l *Label) ToLua(ls *lua.LState) *lua.LUserData {
	result := ls.NewUserData()
	result.Value = l

	ls.SetMetatable(result, ls.GetTypeMetatable(luaTypeExportName))

	return result
}

func FromLua(ud *lua.LUserData) (*Label, error) {
	v, ok := ud.Value.(*Label)

	if !ok {
		return nil, fmt.Errorf("failed to convert")
	}

	return v, nil
}
