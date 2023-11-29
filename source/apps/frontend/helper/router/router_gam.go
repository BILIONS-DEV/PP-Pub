package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Gam(app *fiber.App) {
	Gam := new(controller.Gam)
	app.Get(config.URIGam, Gam.Index)
	app.Post(config.URIGam, Gam.Filter)
	app.Get(config.URIGamAdd, Gam.Add)
	app.Get(config.URIGamEdit, Gam.Edit)
	app.Post(config.URIGamPushLine, Gam.PushLine)
	app.Get(config.URIGamConnect, Gam.Connect)
	app.Get(config.URIGamCallback, Gam.Callback)
	app.Post(config.URIGamSelectNetwork, Gam.SelectNetwork)
	app.Post(config.URIGamGetNetworks, Gam.GetNetworks)
	app.Get(config.URIGam+"/gettoken", Gam.GetToken)
	app.Post(config.URIGamCheckApiAccess, Gam.CheckApiAccess)
}
