package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Target(app *fiber.App) {
	Target := new(controller.Target)
	app.Get(config.URITargetLoadInventory, Target.LoadMoreData)
	app.Get(config.URITargetLoadSelected, Target.LoadSelected)
	app.Post(config.URITargetFilterAdTAg, Target.FilterAdTag)
	app.Get("/target/loadPlaylist", Target.LoadMoreData)
	app.Get("/target/loadsize", Target.LoadSize)
	app.Get("/target/test", Target.Test)
}
