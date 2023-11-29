package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Blocking(app *fiber.App) {
	Blocking := new(controller.Blocking)
	app.Get(config.URIBlocking, Blocking.Index)
	app.Post(config.URIBlocking, Blocking.Filter)
	app.Get(config.URIBlockingAdd, Blocking.Add)
	app.Post(config.URIBlockingAdd, Blocking.AddPost)
	app.Get(config.URIBlockingEdit, Blocking.Edit)
	app.Post(config.URIBlockingEdit, Blocking.EditPost)
	app.Get(config.URIBlockingLoadInventory, Blocking.LoadInventory)
	app.Post(config.URIBlockingDelete, Blocking.Delete)
	app.Post(config.URIBlockingValidateDomain, Blocking.ValidateDomain)
}
