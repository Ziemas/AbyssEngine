package button

import (
	"fmt"

	"github.com/OpenDiablo2/AbyssEngine/common"
	lua "github.com/yuin/gopher-lua"
)

var luaTypeExportName = "buttonlayout"
var LuaTypeExport = common.LuaTypeExport{
	Name: luaTypeExportName,
	//ConstructorFunc: newLuaEntity,
	Methods: map[string]lua.LGFunction{},
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
