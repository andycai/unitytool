package unibuild

import (
	"github.com/andycai/unitool/core"
)

var app *core.App

type uniBuildModule struct {
}

func (u *uniBuildModule) Init(a *core.App) error {
	app = a
	return nil
}

func (u *uniBuildModule) InitDB() error {
	// 数据迁移
	return nil
}

func (u *uniBuildModule) InitData() error {
	// 初始化数据
	return nil
}

func (u *uniBuildModule) InitRouter() error {
	// public
	app.RouterPublic.Post("/api/unibuild/res", buildResources)
	app.RouterPublic.Post("/api/unibuild/app", buildApp)

	// admin

	// api

	return nil
}

func init() {
	core.RegisterModules(&uniBuildModule{})
}
