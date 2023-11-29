package controllerv2

import (
	"github.com/gofiber/fiber/v2"
	"github.com/syyongx/php2go"
	"source/apps/frontend/view"
	"source/internal/entity/dto"
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
	adsUC "source/internal/usecase/ads"
	"source/pkg/htmlblock"
	"strings"
)

type campaign struct {
	*handler
}

func (h *handler) InitRoutesCampaign(app fiber.Router) {
	this := campaign{h}
	// app.Post("/login", this.Login)
	app.Get("/campaigns", this.Index)
	app.Post("/campaigns", this.IndexPost)
	app.Post("/campaigns/change-status-ad", this.ChangeStatusAd)

}

type AssignCampaign struct {
	Assign
	Params      dto.PayloadAdsIndex
	Inventories []model.InventoryModel
}

func (t *campaign) Index(ctx *fiber.Ctx) (err error) {
	// get user login
	userLogin := getUserLogin(ctx)

	params := dto.PayloadAdsIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignCampaign{Assign: newAssign(ctx, "Campaign")}
	if params.OrderColumn == 0 && params.OrderDir == "" {
		params.OrderColumn = 2
		params.OrderDir = "desc"
	}
	// assigns
	assigns.Params = params
	assigns.Inventories, _ = t.useCases.Inventory.GetByPublisher(userLogin.ID)
	return ctx.Render("campaign/index", assigns, view.LAYOUTMain)
}
func (t *campaign) IndexPost(ctx *fiber.Ctx) (err error) {
	// get user login
	userLogin := getUserLogin(ctx)

	var payload dto.PayloadAdsIndexPost
	if err = ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.MakeResponseError(err),
		)
	}

	if errs := payload.Validate(); len(errs) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.MakeResponseErrorWithID(errs...),
		)
	}
	// fmt.Printf("%+v \n", payload)
	var inputUC = adsUC.InputFilterUC{
		Request:     payload.Request,
		User:        userLogin,
		QuerySearch: payload.PostData.QuerySearch,
		Status:      payload.PostData.Status,
		Inventory:   payload.PostData.Inventory,
		OrderColumn: payload.OrderColumn,
		OrderDir:    payload.OrderDir,
	}
	var records []model.AdsModel
	var totalRecords int64
	if records, totalRecords, err = t.useCases.Ads.Filter(&inputUC); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.Fail(err),
		)
	}
	return ctx.JSON(
		datatable.Response{
			Draw:            payload.Draw,
			RecordsTotal:    totalRecords,
			RecordsFiltered: totalRecords,
			Data:            t.makeResponseDatatable(records),
		},
	)
}

func (t *campaign) makeResponseDatatable(ads []model.AdsModel) (records []dto.ResponseAdsDatatable) {
	for _, ad := range ads {
		// fmt.Printf("%#v\n", ad.Audience)
		// var targets = `<div style="height: 30px; text-align: center; margin: 9px auto;">All target...</div>`
		var targets = ""
		var allTarget = false
		// var targets = ``
		if len(ad.Audience.Category.Include) == 0 && len(ad.Audience.Category.Exclude) == 0 && len(ad.Audience.Device.Include) == 0 && len(ad.Audience.Device.Exclude) == 0 && len(ad.Audience.Locations.Include) == 0 && len(ad.Audience.Locations.Exclude) == 0 && len(ad.Audience.Languages.Include) == 0 && len(ad.Audience.Languages.Exclude) == 0 {
			allTarget = true
		}
		targets = htmlblock.Render("campaign/index/target.gohtml", fiber.Map{
			"Audience":  ad.Audience,
			"AllTarget": allTarget,
		}).String()

		var StatusAd = "on"
		if !ad.StatusAd {
			StatusAd = "off"
		}
		var width string
		var height string
		if ad.Size != "" {
			size := strings.Split(ad.Size, "x")
			width = size[0]
			height = size[1]
		}
		rec := dto.ResponseAdsDatatable{
			AdPreview: t.block.RenderToString("campaign/index/block.media.gohtml", fiber.Map{
				"MediaUrl":    ad.MediaURL,
				"MediaType":   ad.MediaType,
				"Thumbnail":   ad.Thumbnail,
				"CTA":         ad.CTA,
				"HeadLine":    ad.Headline,
				"Description": ad.Description,
				"ClickURL":    ad.ClickURL,
				"SiteName":    ad.SiteName,
				"AdType":      ad.AdType,
				"Height":      height,
				"Width":       width,
			}),
			AdInfo: t.block.RenderToString("campaign/index/block.ad_info.gohtml", fiber.Map{
				"Ad":     ad,
				"AdType": strings.Title(ad.AdType),
			}),
			Headline: ad.Headline,
			SiteName: ad.SiteName,
			ClickUrl: ad.ClickURL,
			Target: htmlblock.Render("campaign/index/block.target.gohtml", fiber.Map{
				"Target": (targets),
			}).String(),
			Impressions: php2go.NumberFormat(float64(ad.Impressions), 0, ".", ","),
			Action: t.block.RenderToString("campaign/index/block.action.gohtml", fiber.Map{
				"Action": StatusAd,
				"ID":     ad.ID,
			}),
		}
		records = append(records, rec)
	}
	return
}

func (t *campaign) ChangeStatusAd(ctx *fiber.Ctx) (err error) {
	// get user login
	userLogin := getUserLogin(ctx)
	var payload dto.PayloadChangeActionAd
	if err = ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.Fail(err),
		)
	}
	errs := payload.Validate()
	if len(errs) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.Fail(errs...),
		)
	}

	if err = t.useCases.Ads.ChangeActionAd(&adsUC.ChangeActionAdUC{
		UserLogin: userLogin,
		AdID:      payload.ID,
		Action:    payload.Action,
		Inventory: payload.Inventory,
		Placement: payload.Placement,
	}); err != nil {
		return ctx.JSON(
			dto.Fail(err),
		)
	}
	// fmt.Printf("%+v\n", payload)
	return ctx.JSON(dto.OK(nil, "Success"))
}
