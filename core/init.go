package core

type Module interface {
	Init(*App) error
	InitDB() error
	InitData() error
	InitRouter() error
}

var modules []Module

func RegisterModules(module Module) {
	modules = append(modules, module)
}

func InitModule(app *App) {
	for _, module := range modules {
		module.Init(app)
		module.InitDB()
		module.InitData()
		module.InitRouter()
	}
}
