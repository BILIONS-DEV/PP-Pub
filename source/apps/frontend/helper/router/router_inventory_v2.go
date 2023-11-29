package router

import (
	"source/apps/frontend/config"
	"source/apps/frontend/controller"

	"github.com/gofiber/fiber/v2"
)

func InventoryV2(app *fiber.App) {
	Inventory := new(controller.Inventory)
	//app.Get(config.URIInventory, Inventory.Index)
	//app.Post(config.URIInventory, Inventory.Filter)
	// app.Post(config.URIInventorySubmit, Inventory.Submit)
	app.Get(config.URIInventorySetupV2, Inventory.SetupV2)
	// app.Post(config.URIInventorySetup, Inventory.SetupConfig)
	// app.Post(config.URIInventoryConsent, Inventory.SetupConsent)
	// app.Post(config.URIInventoryUserId, Inventory.SetupUserId)
	// app.Post(config.URIInventoryDelete, Inventory.Delete)
	// app.Post(config.URIInventoryAdTag, Inventory.FilterAdTag)
	// app.Get(config.URIInventoryLoadParam, Inventory.LoadModuleParam)
	// app.Post(config.URIInventoryCollapse, Inventory.Collapse)
	// app.Get(config.URIInventoryCopyAdTag, Inventory.CopyAdTag)
	// app.Post(config.URIInventoryBuildScript, Inventory.BuildScript)
	// app.Post(config.URIInventoryConnection, Inventory.FilterConnection)
	// app.Post(config.URIInventoryChangeStatusConnection, Inventory.ChangeStatusConnection)
}
