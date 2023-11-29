package router

import (
	"source/apps/frontend/config"
	"source/apps/frontend/controller"

	"github.com/gofiber/fiber/v2"
)

func AdTagV2(app *fiber.App) {
	AdTag := new(controller.AdTag)
	// app.Get(config.URIAdTag, AdTag.Index)
	app.Get(config.URIAdTagAddV2, AdTag.AddV2)
	app.Post(config.URIAdTagAddV2, AdTag.AddPostV2)
	app.Get(config.URIAdTagEditV2, AdTag.EditV2)
	// app.Post(config.URIAdTagEdit, AdTag.EditPost)
	// app.Post(config.URIAdTagDel, AdTag.Delete)
	// app.Get(config.URIAdTagGetSizeAdditional, AdTag.GetSizeAdditional)
	// app.Post(config.URIAdTagCollapse, AdTag.Collapse)
}
