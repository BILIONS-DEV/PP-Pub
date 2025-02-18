package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Api(app *fiber.App) {
	api := new(controller.Api)
	app.Get(config.URIAPIGetInfoAccount, api.GetInfoAccount)
}
