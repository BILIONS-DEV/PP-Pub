package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/controller"
)

func Home(app *fiber.App) {
	Home := new(controller.Home)
	app.Get("/", Home.Dashboard)
	app.Get("/dashboards", Home.Dashboard)
	//app.Get("/test", Home.Test)
	app.Get("/test", Home.TestNpm)
}

func HomeBackend(app fiber.Router) {

}
