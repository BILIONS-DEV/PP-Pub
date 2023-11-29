package controller

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/view"
)

type Home struct{}

type AssignHome struct {
	assign.Schema
}

func (th *Home) Dashboard(ctx *fiber.Ctx) error {
	assigns := AssignHome{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Dashboards")
	return ctx.Render("home/dashboard", assigns, view.LAYOUTMain)
}

func (th *Home) Test(ctx *fiber.Ctx) error {
	assigns := AssignHome{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Test")
	return ctx.Render("home/test", assigns, view.LAYOUTEmpty)
}

func (t *Home) TestNpm(ctx *fiber.Ctx) error {
	assigns := AssignHome{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Test NPM")
	return ctx.Render("home/test", assigns, view.LAYOUTMain)
}
