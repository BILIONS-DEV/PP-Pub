package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func AdBlock(app *fiber.App) {
	AdBlock := new(controller.AdBlock)
	app.Get(config.URIAdBlockAnalytics, AdBlock.Analytics)
	app.Post(config.URIAdBlockAnalytics, AdBlock.AnalyticsFilter)
	app.Get(config.URIAdBlockAlertGenerator, AdBlock.AdblockAlertGenerator)
}
