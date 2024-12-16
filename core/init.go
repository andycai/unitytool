package core

type Module interface {
	Awake(*App) error
	Start() error
	AddPublicRouters() error
	AddAuthRouters() error
}

var modules []Module

// RegisterModule 注册模块
func RegisterModule(module Module) {
	modules = append(modules, module)
}

// InitPublicRouters 初始化公共路由
func InitPublicRouters() {
	for _, module := range modules {
		module.AddPublicRouters()
	}
}

// InitAuthRouters 初始化管理员路由
func InitAuthRouters() {
	for _, module := range modules {
		module.AddAuthRouters()
	}
}

// AwakeModules 模块初始化
func AwakeModules(app *App) {
	// 模块初始化
	for _, module := range modules {
		module.Awake(app)
	}

	// 模块启动
	for _, module := range modules {
		module.Start()
	}
}
