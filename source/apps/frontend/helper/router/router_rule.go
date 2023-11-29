package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Rule(app *fiber.App) {
	Rule := new(controller.Rule)
	app.Get(config.URIRule, Rule.Index)
	app.Post(config.URIRule, Rule.Filter)
}
