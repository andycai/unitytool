package core

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var dbMap = map[string]func([]*gorm.DB){}

var moduleMap = map[string]func(){}

var routerPublicNoCheckMap = map[string]func(fiber.Router){}
var routerRootCheckMap = map[string]func(fiber.Router){}
var routerAPINoCheckMap = map[string]func(fiber.Router){}
var routerAPICheckMap = map[string]func(fiber.Router){}
var routerAdminCheckMap = map[string]func(fiber.Router){}

func InitModules() {
	for _, f := range moduleMap {
		f()
	}
}

func RegisterModule(moduleType string, f func()) {
	if _, ok := moduleMap[moduleType]; ok {
		panic("duplicate module type: " + moduleType)
	}
	moduleMap[moduleType] = f
}

func RegisterDatabase(dbType string, f func([]*gorm.DB)) {
	if _, ok := dbMap[dbType]; ok {
		panic("duplicate db type: " + dbType)
	}
	dbMap[dbType] = f
}

func RegisterPublicRouter(routerType string, f func(fiber.Router)) {
	if _, ok := routerPublicNoCheckMap[routerType]; ok {
		panic("duplicate router type: " + routerType)
	}
	routerPublicNoCheckMap[routerType] = f
}

func RegisterRootCheckRouter(routerType string, f func(fiber.Router)) {
	if _, ok := routerRootCheckMap[routerType]; ok {
		panic("duplicate router type: " + routerType)
	}
	routerRootCheckMap[routerType] = f
}

func RegisterAPINoCheckRouter(routerType string, f func(fiber.Router)) {
	if _, ok := routerAPINoCheckMap[routerType]; ok {
		panic("duplicate router type: " + routerType)
	}
	routerAPINoCheckMap[routerType] = f
}

func RegisterAPICheckRouter(routerType string, f func(fiber.Router)) {
	if _, ok := routerAPICheckMap[routerType]; ok {
		panic("duplicate router type: " + routerType)
	}
	routerAPICheckMap[routerType] = f
}

func RegisterAdminCheckRouter(routerType string, f func(fiber.Router)) {
	if _, ok := routerAdminCheckMap[routerType]; ok {
		panic("duplicate router type: " + routerType)
	}
	routerAdminCheckMap[routerType] = f
}
