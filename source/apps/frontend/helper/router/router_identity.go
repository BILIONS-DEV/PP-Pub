package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Identity(app *fiber.App) {
	Identity := new(controller.Identity)
	app.Get(config.URIIdentity, Identity.Index)
	app.Post(config.URIIdentity, Identity.Filter)
	app.Get(config.URIIdentityAdd, Identity.Add)
	app.Post(config.URIIdentityAdd, Identity.AddPost)
	app.Get(config.URIIdentityEdit, Identity.Edit)
	app.Post(config.URIIdentityEdit, Identity.EditPost)
	app.Post(config.URIIdentityDel, Identity.Delete)
	app.Post(config.URIIdentityCollapse, Identity.Collapse)
	app.Post(config.URIIdentityChangeStatus, Identity.ChangeStatus)
}
