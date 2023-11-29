package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func AdsTxt(app *fiber.App) {
	AdsTxt := new(controller.AdsTxt)
	app.Get("/jsx/test", AdsTxt.Test)
	app.Get(config.URIAdsTxt, AdsTxt.Index)
	app.Post(config.URIAdsTxt, AdsTxt.Filter)
	app.Get(config.URIAdsTxtDetail, AdsTxt.Detail)
	app.Post(config.URIAdsTxtDetail, AdsTxt.SaveAdsTxt)
	app.Post(config.URIAdsTxtScan, AdsTxt.Scan)
	app.Post(config.URIAdsTxtLoad, AdsTxt.Load)
}
