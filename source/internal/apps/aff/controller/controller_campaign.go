package controller

import (
	"github.com/gofiber/fiber/v2"
	"source/internal/apps/aff/controller/params"
	"source/internal/entity/dto"
	"source/internal/entity/model"
)

type campaignController struct {
	*handler
}

func (h *handler) InitRoutesCampaign(app fiber.Router) {
	ctl := campaignController{h}
	routerGroup := app.Group("/campaign")
	routerGroup.Get("/", ctl.Index)
	routerGroup.Get("/add", ctl.Add)
	routerGroup.Get("/edit", ctl.Edit)
	routerGroup.Get("/:campID", ctl.BuildCampaign)
}

func (t *campaignController) Index(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"page": "/campaign/index",
	})
}

type PayloadAddCampaignDTO struct {
	Name          string `json:"name"`
	TrafficSource string `json:"traffic_source"`
	DemandSource  string `json:"demand_source"`
}

func (t *campaignController) Add(ctx *fiber.Ctx) (err error) {
	// => get payload
	var payload PayloadAddCampaignDTO
	if err = ctx.BodyParser(&payload); err != nil {
		return ctx.JSON(dto.Fail(err))
	}

	// => convert payload to Model
	var rec = model.CampaignModel{
		ID:   0,
		Name: payload.Name,
	}

	// => logic handle
	if err = t.UseCases.Campaign.AddCampaign2(&rec); err != nil {
		return ctx.JSON(dto.Fail(err))
	}

	return ctx.JSON(dto.OK(nil, "done"))
}

type PayloadEditCampaignDTO struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	TrafficSource string `json:"traffic_source"`
	DemandSource  string `json:"demand_source"`
}

func (t *campaignController) Edit(ctx *fiber.Ctx) (err error) {
	// var ID int
	// if ID, err = ctx.ParamsInt("uid"); err != nil {
	// 	return ctx.JSON(dto.Fail(err))
	// }
	var payload PayloadEditCampaignDTO
	if err = ctx.BodyParser(&payload); err != nil {
		return ctx.JSON(dto.Fail(err))
	}

	var rec = model.CampaignModel{
		ID:            payload.ID,
		Name:          payload.Name,
		TrafficSource: payload.TrafficSource,
		DemandSource:  payload.DemandSource,
	}
	if err = t.UseCases.Campaign.EditCampaign2(&rec); err != nil {
		return ctx.JSON(dto.Fail(err))
	}
	return ctx.JSON(dto.OK(nil, "done"))
}

func (t *campaignController) BuildCampaign(ctx *fiber.Ctx) (err error) {
	var ID int
	if ID, err = ctx.ParamsInt("campID"); err != nil {
		return ctx.JSON(dto.Fail(err))
	}

	campaign := t.UseCases.Campaign.GetCampaignById(int64(ID))
	if campaign.ID == 0 {
		return ctx.JSON(dto.FailWithString("Campaign does not exist!"))
	}

	var payload params.ParamsCampaignDTO
	if err = ctx.QueryParser(&payload); err != nil {
		return ctx.JSON(dto.Fail(err))
	}

	return ctx.JSON(dto.OK(nil, "done"))
}
