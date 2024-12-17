package modules

// 新增的模块必须在这里进行导入，不然模块 init 方法不会执行
import (
	_ "github.com/andycai/unitool/modules/adminlog"
	_ "github.com/andycai/unitool/modules/browse"
	_ "github.com/andycai/unitool/modules/citask"
	_ "github.com/andycai/unitool/modules/gamelog"
	_ "github.com/andycai/unitool/modules/login"
	_ "github.com/andycai/unitool/modules/menu"
	_ "github.com/andycai/unitool/modules/permission"
	_ "github.com/andycai/unitool/modules/role"
	_ "github.com/andycai/unitool/modules/serverconf"
	_ "github.com/andycai/unitool/modules/shell"
	_ "github.com/andycai/unitool/modules/stats"
	_ "github.com/andycai/unitool/modules/unibuild"
	_ "github.com/andycai/unitool/modules/user"
	// _ "github.com/andycai/unitool/modules/note"
)
