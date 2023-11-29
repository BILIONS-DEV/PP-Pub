package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Report(app *fiber.App) {
	Report := new(controller.Report)
	app.Get(config.URIReport, Report.Index)
	app.Post(config.URIReport, Report.Filter)
}
