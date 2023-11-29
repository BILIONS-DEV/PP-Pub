package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Notification(app *fiber.App) {
	Notification := new(controller.Notification)
	app.Get(config.URINOTIFICATION, Notification.Index)
	app.Get(config.URINOTIFICATIONREADALL, Notification.ReadAll)
	app.Get(config.URINOTIFICATIONREAD, Notification.Read)
}
