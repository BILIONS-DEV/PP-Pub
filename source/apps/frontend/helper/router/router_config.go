package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Config(app *fiber.App) {
	Config := new(controller.Config)
	app.Get(config.URIConfig, Config.IndexConfig)
	app.Post(config.URIConfigSave, Config.ConfigSave)
}
