package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Inventory(app *fiber.App) {
	Inventory := new(controller.Inventory)
	//app.Get(config.URIInventory, Inventory.Index)
	//app.Post(config.URIInventory, Inventory.Filter)
	app.Post(config.URIInventorySubmit, Inventory.Submit)
	app.Get(config.URIInventorySetup, Inventory.Setup)
	app.Post(config.URIInventorySetup, Inventory.SetupConfig)
	app.Post(config.URIInventoryConsent, Inventory.SetupConsent)
	app.Post(config.URIInventoryUserId, Inventory.SetupUserId)
	app.Post(config.URIInventoryDelete, Inventory.Delete)
	app.Post(config.URIInventoryAdTag, Inventory.FilterAdTag)
	app.Get(config.URIInventoryLoadParam, Inventory.LoadModuleParam)
	app.Post(config.URIInventoryCollapse, Inventory.Collapse)
	app.Get(config.URIInventoryCopyAdTag, Inventory.CopyAdTag)
	app.Post(config.URIInventoryBuildScript, Inventory.BuildScript)
	app.Post(config.URIInventoryConnection, Inventory.FilterConnection)
	app.Post(config.URIInventoryChangeStatusConnection, Inventory.ChangeStatusConnection)
}
