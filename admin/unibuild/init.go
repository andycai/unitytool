package unibuild

import (
	"github.com/andycai/unitool/core"
)

var app *core.App

type uniBuildModule struct {
}

func (m *uniBuildModule) Init(a *core.App) error {
	app = a
	return nil
}

func (m *uniBuildModule) InitDB() error {
	// 数据迁移
	return nil
}

func (m *uniBuildModule) InitModule() error {
	// public
	app.RouterPublic.Post("/api/unibuild/res", buildResources)
	app.RouterPublic.Post("/api/unibuild/app", buildApp)

	// admin

	// api

	return nil
}

func init() {
	core.RegisterModule(&uniBuildModule{})
}
