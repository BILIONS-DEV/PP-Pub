package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Video(app *fiber.App) {
	Video := new(controller.Video)
	app.Get(config.URILinkVideo, Video.Index)
}
