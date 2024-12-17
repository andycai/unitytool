package note

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/enum"
)

var app *core.App

type noteModule struct {
}

func init() {
	core.RegisterModule(&noteModule{}, enum.ModulePriorityNote)
}

func (m *noteModule) Awake(a *core.App) error {
	app = a
	// 数据迁移
	if err := autoMigrate(); err != nil {
		return err
	}

	// 初始化数据
	return initData()
}

func (m *noteModule) Start() error {
	return nil
}

func (m *noteModule) AddPublicRouters() error {
	// 公开API路由
	app.RouterPublicApi.Get("/notes/public", handlePublicNotes)
	app.RouterPublicApi.Get("/notes/public/:id", handlePublicNoteDetail)
	app.RouterPublicApi.Get("/notes/categories/public", handlePublicCategories)

	return nil
}

func (m *noteModule) AddAuthRouters() error {
	// 管理后台路由
	app.RouterAdmin.Get("/notes", app.HasPermission("note:list"), handleNoteList)

	// API路由
	app.RouterApi.Get("/notes/tree", app.HasPermission("note:list"), handleNoteTree)
	app.RouterApi.Get("/notes/:id", app.HasPermission("note:list"), handleNoteDetail)
	app.RouterApi.Post("/notes", app.HasPermission("note:create"), handleNoteCreate)
	app.RouterApi.Put("/notes/:id", app.HasPermission("note:update"), handleNoteUpdate)
	app.RouterApi.Delete("/notes/:id", app.HasPermission("note:delete"), handleNoteDelete)

	// 分类操作
	app.RouterApi.Get("/notes/categories", app.HasPermission("note:category:list"), handleCategoryList)
	app.RouterApi.Post("/notes/categories", app.HasPermission("note:category:create"), handleCategoryCreate)
	app.RouterApi.Put("/notes/categories/:id", app.HasPermission("note:category:update"), handleCategoryUpdate)
	app.RouterApi.Delete("/notes/categories/:id", app.HasPermission("note:category:delete"), handleCategoryDelete)

	return nil
}
