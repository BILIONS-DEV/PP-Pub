package router

import (
	"source/apps/frontend/config"
	"source/apps/frontend/controller"

	"github.com/gofiber/fiber/v2"
)

func PlayerV2(app *fiber.App) {
	PlayerV2 := new(controller.PlayerV2)
	app.Get(config.URIPlayerTemplateV2, PlayerV2.Index)
	app.Post(config.URIPlayerTemplateV2, PlayerV2.Filter)
	app.Get(config.URIPlayerAddTemplateV2, PlayerV2.AddTemplate)
	// app.Post(config.URIPlayerAddTemplate, Player.AddPost)
	app.Get(config.URIPlayerEditTemplateV2, PlayerV2.Edit)
	// app.Post(config.URIPlayerEditTemplate, Player.EditPost)
	app.Get(config.URIPlayerViewTemplateV2, PlayerV2.View)
	// app.Post(config.URIPlayerDelTemplate, Player.Delete)
	// app.Post(config.URIPlayerDuplicateTemplate, Player.Duplicate)
	// app.Post(config.URIPlayerCollapse, Player.Collapse)
	app.Post(config.URIPlayerPreviewV2, PlayerV2.Preview)
}
