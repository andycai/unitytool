package admin

// 新增的模块必须在这里进行导入，不然模块 init 方法不会执行
import (
	_ "github.com/andycai/unitool/admin/adminlog"
	_ "github.com/andycai/unitool/admin/browse"
	_ "github.com/andycai/unitool/admin/citask"
	_ "github.com/andycai/unitool/admin/gamelog"
	_ "github.com/andycai/unitool/admin/login"
	_ "github.com/andycai/unitool/admin/menu"
	_ "github.com/andycai/unitool/admin/permission"
	_ "github.com/andycai/unitool/admin/role"
	_ "github.com/andycai/unitool/admin/serverconf"
	_ "github.com/andycai/unitool/admin/shell"
	_ "github.com/andycai/unitool/admin/stats"
	_ "github.com/andycai/unitool/admin/unibuild"
	_ "github.com/andycai/unitool/admin/user"
)
