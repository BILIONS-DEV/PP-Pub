package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Support(app *fiber.App) {
	Support := new(controller.Support)
	app.Get(config.URISUPPORT, Support.Index)
	app.Get(config.URISUPPORTPRODUCTDESCRIOPTION, Support.ProductDescription)
	app.Get(config.URISUPPORTTICKETSNEW, Support.TichketsNew)
}
