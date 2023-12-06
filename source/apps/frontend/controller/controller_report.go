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

func (t *Report) Dimension(ctx *fiber.Ctx) error {
	assigns := AssignReport{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Reports By Dimension")
	return ctx.Render("report/dimension", assigns, view.LAYOUTMain)
}

func (t *Report) Saved(ctx *fiber.Ctx) error {
	assigns := AssignReport{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Reports Saved Queries")
	return ctx.Render("report/saved", assigns, view.LAYOUTMain)
}
