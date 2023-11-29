package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Bidder(app *fiber.App) {
	Bidder := new(controller.Bidder)
	app.Get(config.URISystemBidder, Bidder.Index)
	app.Post(config.URISystemBidder, Bidder.Filter)
	app.Get(config.URISystemBidderAdd, Bidder.Add)
	app.Post(config.URISystemBidderAdd, Bidder.AddPost)
	app.Get(config.URISystemBidderAddTemplate, Bidder.AddTemplate)
	app.Get(config.URISystemBidderEdit, Bidder.Edit)
	app.Post(config.URISystemBidderEdit, Bidder.EditPost)
	app.Post(config.URISystemBidderDel, Bidder.Delete)
	app.Get(config.URISystemBidderView, Bidder.Edit)
	app.Get(config.URISystemAddParam, Bidder.AddParamBidder)
	//app.Post(config.URISystemBidderUploadXlsx, Bidder.UploadHandleXlsx)
}
