package unibuild

import (
	"github.com/andycai/unitool/core"
)

var app *core.App

type uniBuildModule struct {
}

func (m *uniBuildModule) Awake(a *core.App) error {
	app = a
	return nil
}

func (m *uniBuildModule) Start() error {
	return nil
}

func (m *uniBuildModule) AddPublicRouters() error {
	// public
	app.RouterPublicApi.Post("/unibuild/res", buildResources)
	app.RouterPublicApi.Post("/unibuild/app", buildApp)

	return nil
}

func (m *uniBuildModule) AddAuthRouters() error {
	return nil
}

func init() {
	core.RegisterModule(&uniBuildModule{})
}
