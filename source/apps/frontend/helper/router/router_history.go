package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/controller"
	"source/apps/frontend/config"
)

func History(app *fiber.App) {
	History := new(controller.History)
	app.Get(config.URIHistory, History.Index)
	app.Post(config.URIObjectByPage, History.ObjectByPage)
	// app.Post(config.URIHistory, History.Filter)
	app.Post(config.URIHistoryLoad, History.LoadHistory)
}
