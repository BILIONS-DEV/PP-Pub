package controller

import "github.com/gofiber/fiber/v2"

type userController struct {
	*handler
}

func (h *handler) InitRoutesUser(app fiber.Router) {
	ctl := userController{h}
	app.Get("/user", ctl.Index)
}

func (t userController) Index(ctx *fiber.Ctx) (err error) {
	return ctx.JSON(fiber.Map{
		"hello": "world",
	})
}
