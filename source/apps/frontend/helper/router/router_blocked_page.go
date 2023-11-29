package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func BlockedPage(app *fiber.App) {
	BlockedPage := new(controller.BlockedPage)
	app.Get(config.URIBlockedPageAdd, BlockedPage.Add)
	app.Post(config.URIBlockedPageAdd, BlockedPage.AddPost)
	app.Get(config.URIBlockedPageEdit, BlockedPage.Edit)
	app.Post(config.URIBlockedPageEdit, BlockedPage.EditPost)
	app.Post(config.URIBlockedPageImportCSV, BlockedPage.ImportCSV)
	app.Post(config.URIBlockedPageDel, BlockedPage.Delete)
}
