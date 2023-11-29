package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func AbTesting(app *fiber.App) {
	AbTesting := new(controller.AbTesting)
	app.Get(config.URIAbTesting, AbTesting.Index)
	app.Post(config.URIAbTesting, AbTesting.Filter)
	app.Get(config.URIAbTestingAdd, AbTesting.Add)
	app.Post(config.URIAbTestingAdd, AbTesting.AddPost)
	app.Get(config.URIAbTestingEdit, AbTesting.Edit)
	app.Post(config.URIAbTestingEdit, AbTesting.EditPost)
	app.Post(config.URIAbTestingDel, AbTesting.Delete)
	app.Post(config.URIAbTestingCollapse, AbTesting.Collapse)
}
