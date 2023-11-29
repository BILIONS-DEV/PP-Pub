package controller

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
)

type Rule struct {
}

type AssignRuleIndex struct {
	assign.Schema
	Params payload.RuleIndex
}

func (t *Rule) Index(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIRule)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	params := payload.RuleIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignRuleIndex{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Rule")
	return ctx.Render("rule/index", assigns, view.LAYOUTMain)
}

func (t *Rule) Filter(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIRule)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get payload post
	var inputs payload.RuleFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	user := GetUserLogin(ctx)
	if !user.IsFound() {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get data from model
	dataTable, err := new(model.Rule).GetByFilters(&inputs, user)
	if err != nil {
		return err
	}
	// Print Data to JSON
	return ctx.JSON(dataTable)
}
