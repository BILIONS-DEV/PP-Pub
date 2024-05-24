package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func AdsTxt(app *fiber.App) {
	// get url
	adsTxt := new(controller.AdsTxt)
	app.Get("/jsx/test", adsTxt.Test)
	app.Get(config.URIAdsTxt, adsTxt.Index)
	app.Post(config.URIAdsTxt, adsTxt.Filter)
	app.Get(config.URIAdsTxtDetail, adsTxt.Detail)
	app.Post(config.URIAdsTxtDetail, adsTxt.SaveAdsTxt)
	app.Post(config.URIAdsTxtScan, adsTxt.Scan)
	app.Post(config.URIAdsTxtLoad, adsTxt.Load)
}
