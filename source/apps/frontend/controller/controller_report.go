package controller

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/view"
)

type Report struct{}

type AssignReport struct {
	assign.Schema
}

func (t *Report) Index(ctx *fiber.Ctx) error {
	assigns := AssignReport{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Report")
	return ctx.Render("report/index", assigns, view.LAYOUTMain)
}

func (t *Report) Filter(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}
