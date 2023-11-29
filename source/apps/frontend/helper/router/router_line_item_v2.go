package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func LineItemV2(app *fiber.App) {
	LineItem := new(controller.LineItemV2)
	app.Get(config.URILineItemV2, LineItem.Index)
	app.Post(config.URILineItemV2, LineItem.Filter)
	app.Get(config.URILineItemV2Choose, LineItem.Choose)
	//app.Get(config.URILineItemV2Create, LineItem.Create)
	//app.Post(config.URILineItemV2Create, LineItem.CreatePost)
	app.Get(config.URILineItemV2Add, LineItem.Add)
	app.Post(config.URILineItemV2Add, LineItem.AddPost)
	app.Get(config.URILineItemV2Edit, LineItem.Edit)
	app.Post(config.URILineItemV2Edit, LineItem.EditPost)
	app.Post(config.URILineItemV2Delete, LineItem.Delete)
	app.Get(config.URILineItemV2View, LineItem.Edit)
	app.Get(config.URILineItemV2LoadParam, LineItem.BidderParam)
	app.Get(config.URISearchDomainV2, LineItem.SearchDomain)
	app.Get(config.URISearchAdFormatV2, LineItem.SearchAdFormat)
	app.Get(config.URISearchAdTagV2, LineItem.SearchAdTag)
	app.Get(config.URISearchAdSizeV2, LineItem.SearchAdSize)
	app.Get(config.URISearchDeviceV2, LineItem.SearchDevice)
	app.Get(config.URISearchCountryV2, LineItem.SearchCountry)
	app.Post(config.URILineItemV2Collapse, LineItem.Collapse)
	app.Get(config.URILineItemV2CheckParam, LineItem.CheckParam)
	app.Get(config.URILineItemV2ListLinkedGam, LineItem.ListLinkedGam)
	app.Get(config.URIAddParamBidderV2, LineItem.AddParamBidder)
}
