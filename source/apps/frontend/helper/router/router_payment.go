package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Payment(app *fiber.App) {
	Payment := new(controller.Payment)
	app.Get(config.URIPayment, Payment.Index)
	app.Post(config.URIPayment, Payment.IndexFilter)
	// app.Get(config.URIPaymentPreview, Payment.Preivew)
	// app.Get(config.URIPaymentExport, Payment.ExportPDF)
}
