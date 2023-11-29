package controller

import (
	// "encoding/json"
	// "fmt"
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/view"
	// "source/pkg/ajax"
	// "source/pkg/utility"
	// "strconv"
)

type Support struct{}

type AssignSupport struct {
	assign.Schema
	// Categories []model.CategoryRecord
	// Params     payload.ContentIndex
}

func (t *Support) Index(ctx *fiber.Ctx) error {
	// params := payload.ContentIndex{}
	// if err := ctx.QueryParser(&params); err != nil {
	// 	return err
	// }
	assigns := AssignSupport{Schema: assign.Get(ctx)}
	// assigns.Categories = new(model.Category).GetAll()
	// assigns.Params = params
	assigns.Title = config.TitleWithPrefix("Support")
	assigns.LANG.Pages.Support.Title = "Support"
	return ctx.Render("support/index", assigns, view.LAYOUTMain)
}

func (t *Support) ProductDescription(ctx *fiber.Ctx) error {
	// params := payload.ContentIndex{}
	// if err := ctx.QueryParser(&params); err != nil {
	// 	return err
	// }
	assigns := AssignSupport{Schema: assign.Get(ctx)}
	// assigns.Categories = new(model.Category).GetAll()
	// assigns.Params = params
	assigns.Title = config.TitleWithPrefix("Support")
	assigns.LANG.Pages.Support.Title = "Support"
	return ctx.Render("support/product_description", assigns, view.LAYOUTMain)
}
func (t *Support) TichketsNew(ctx *fiber.Ctx) error {
	// params := payload.ContentIndex{}
	// if err := ctx.QueryParser(&params); err != nil {
	// 	return err
	// }
	assigns := AssignSupport{Schema: assign.Get(ctx)}
	// assigns.Categories = new(model.Category).GetAll()
	// assigns.Params = params
	assigns.Title = config.TitleWithPrefix("Support")
	assigns.LANG.Pages.Support.Title = "Support"
	return ctx.Render("support/tickets_new", assigns, view.LAYOUTMain)
}
