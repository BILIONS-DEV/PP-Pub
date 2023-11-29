package router

import (
	"source/apps/frontend/config"
	"source/apps/frontend/controller"

	"github.com/gofiber/fiber/v2"
)

func Player(app *fiber.App) {
	Player := new(controller.Player)
	app.Get(config.URIPlayerTemplate, Player.Index)
	app.Post(config.URIPlayerTemplate, Player.Filter)
	app.Get(config.URIPlayerAddTemplate, Player.AddTemplate)
	app.Post(config.URIPlayerAddTemplate, Player.AddPost)
	app.Get(config.URIPlayerEditTemplate, Player.Edit)
	app.Post(config.URIPlayerEditTemplate, Player.EditPost)
	app.Get(config.URIPlayerViewTemplate, Player.View)
	app.Post(config.URIPlayerDelTemplate, Player.Delete)
	app.Post(config.URIPlayerDuplicateTemplate, Player.Duplicate)
	app.Post(config.URIPlayerCollapse, Player.Collapse)
	app.Post(config.URIPlayerPreview, Player.Preview)
}
