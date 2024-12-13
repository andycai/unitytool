package core

type Module interface {
	Init(*App) error
	InitDB() error
	InitModule() error
}

var modules []Module

func RegisterModule(module Module) {
	modules = append(modules, module)
}

func InitModules(app *App) {
	for _, module := range modules {
		module.Init(app)
	}

	// 初始化数据库和数据
	for _, module := range modules {
		module.InitDB()
	}

	// 初始化模块
	for _, module := range modules {
		module.InitModule()
	}
}
