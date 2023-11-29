package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func LineItem(app *fiber.App) {
	LineItem := new(controller.LineItem)
	app.Get(config.URILineItem, LineItem.Index)
	app.Post(config.URILineItem, LineItem.Filter)
	app.Get(config.URILineItemChoose, LineItem.Choose)
	//app.Get(config.URILineItemCreate, LineItem.Create)
	//app.Post(config.URILineItemCreate, LineItem.CreatePost)
	app.Get(config.URILineItemAdd, LineItem.Add)
	app.Post(config.URILineItemAdd, LineItem.AddPost)
	app.Get(config.URILineItemEdit, LineItem.Edit)
	app.Post(config.URILineItemEdit, LineItem.EditPost)
	app.Post(config.URILineItemDelete, LineItem.Delete)
	app.Get(config.URILineItemView, LineItem.Edit)
	app.Get(config.URILineItemLoadParam, LineItem.BidderParam)
	app.Get(config.URISearchDomain, LineItem.SearchDomain)
	app.Get(config.URISearchAdFormat, LineItem.SearchAdFormat)
	app.Get(config.URISearchAdTag, LineItem.SearchAdTag)
	app.Get(config.URISearchAdSize, LineItem.SearchAdSize)
	app.Get(config.URISearchDevice, LineItem.SearchDevice)
	app.Get(config.URISearchCountry, LineItem.SearchCountry)
	app.Post(config.URILineItemCollapse, LineItem.Collapse)
	app.Get(config.URILineItemCheckParam, LineItem.CheckParam)
	app.Get(config.URILineItemListLinkedGam, LineItem.ListLinkedGam)
	app.Get(config.URIAddParamBidder, LineItem.AddParamBidder)
}
