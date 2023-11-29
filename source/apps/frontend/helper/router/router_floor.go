package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Floor(app *fiber.App) {
	Floor := new(controller.Floor)
	app.Get(config.URIFloor, Floor.Index)
	app.Post(config.URIFloor, Floor.Filter)
	app.Get(config.URIFloorAdd, Floor.Add)
	app.Post(config.URIFloorAdd, Floor.AddPost)
	app.Get(config.URIFloorEdit, Floor.Edit)
	app.Post(config.URIFloorEdit, Floor.EditPost)
	app.Post(config.URIFloorDel, Floor.Delete)
	app.Post(config.URIFloorCollapse, Floor.Collapse)
}
