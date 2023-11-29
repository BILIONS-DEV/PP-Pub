package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Channels(app *fiber.App) {
	Channels := new(controller.Channels)
	app.Get(config.URIChannels, Channels.Index)
	app.Post(config.URIChannels, Channels.Filter)
	app.Get(config.URIChannelsAdd, Channels.Add)
	app.Post(config.URIChannelsAdd, Channels.AddPost)
	app.Get(config.URIChannelsEdit, Channels.Edit)
	app.Post(config.URIChannelsEdit, Channels.EditPost)
	app.Post(config.URIChannelsDel, Channels.Delete)
}
