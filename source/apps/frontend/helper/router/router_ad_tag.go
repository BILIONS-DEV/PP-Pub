package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func AdTag(app *fiber.App) {
	AdTag := new(controller.AdTag)
	app.Get(config.URIAdTag, AdTag.Index)
	app.Get(config.URIAdTagAdd, AdTag.Add)
	app.Post(config.URIAdTagAdd, AdTag.AddPost)
	app.Get(config.URIAdTagEdit, AdTag.Edit)
	app.Post(config.URIAdTagEdit, AdTag.EditPost)
	app.Post(config.URIAdTagDel, AdTag.Delete)
	app.Get(config.URIAdTagGetSizeAdditional, AdTag.GetSizeAdditional)
	app.Post(config.URIAdTagCollapse, AdTag.Collapse)
}
