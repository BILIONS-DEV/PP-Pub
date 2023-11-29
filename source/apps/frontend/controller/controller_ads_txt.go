package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
	"source/pkg/ajax"
	"source/pkg/block"
	"source/pkg/cloudflare"
	"source/pkg/utility"
	"strconv"
)

type AdsTxt struct{}

type AssignAdsTxt struct {
	assign.Schema
	AdsTxtMissingLineSyncError string
	AdsTxtMissingLines         []model.AdsTxtMissingLine
	Domain                     model.InventoryRecord
}

type AssignAdsTxtIndex struct {
	assign.Schema
	Params payload.InventoryIndex
}

func (t *AdsTxt) Test(ctx *fiber.Ctx) error {
	assigns := AssignAdsTxtIndex{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("TEST JSX")
	return ctx.Render("ads_txt/test", assigns, view.LAYOUTMain)
}

func (t *AdsTxt) Index(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAdsTxt)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	params := payload.InventoryIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignAdsTxtIndex{Schema: assign.Get(ctx)}
	assigns.Params = params
	assigns.Title = config.TitleWithPrefix("Ads Txt for domain")
	return ctx.Render("ads_txt/index", assigns, view.LAYOUTMain)
}

func (t *AdsTxt) Filter(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAdsTxt)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get payload post
	var inputs payload.InventoryFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Get data from model
	dataTable, err := new(model.AdsTxt).GetByFilters(&inputs, userLogin.Id, GetLang(ctx))
	if err != nil {
		return err
	}
	// Print Data to JSON
	return ctx.JSON(dataTable)
}

func (t *AdsTxt) SaveAdsTxt(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAdsTxtDetail)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	domainIdString := ctx.Query("did")
	if domainIdString == "" {
		return ctx.JSON(ajax.Responses{
			Status:  "error",
			Message: "Not valid",
		})
	}

	domainId, err := strconv.ParseInt(domainIdString, 10, 64)
	if err != nil {
		return ctx.JSON(ajax.Responses{
			Status:  "error",
			Message: err.Error(),
		})
	}

	domain, err := new(model.Inventory).GetById(domainId, userLogin.Id)
	if err != nil {
		return ctx.JSON(ajax.Responses{
			Status:  "error",
			Message: err.Error(),
		})
	}

	adsTxtContent := ctx.FormValue("ads_txt")
	adsTxtCustom, err := new(model.AdsTxt).PushForInventory(domain, adsTxtContent, GetLang(ctx))
	if err != nil {
		return ctx.JSON(ajax.Responses{
			Status:  "error",
			Message: err.Error(),
		})
	}

	// Reset cache cho link ads.txt
	err = cloudflare.ResetLinkAds([]string{domain.Uuid})
	if err != nil {
		fmt.Println(err)
	}
	return ctx.JSON(ajax.Responses{
		Status:     "success",
		Message:    "ads.txt has been updated",
		DataObject: adsTxtCustom,
	})
}

func (t *AdsTxt) Detail(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIInventorySetup)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignAdsTxt{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Ads Txt for domain")

	domainIdString := ctx.Query("did")
	if domainIdString == "" {
		return ctx.JSON(ajax.Responses{
			Status:  "error",
			Message: "Not valid",
		})
	}

	domainId, err := strconv.ParseInt(domainIdString, 10, 64)
	if err != nil {
		if !utility.IsWindow() {
			return ctx.JSON(ajax.Responses{
				Status:  "error",
				Message: GetLang(ctx).Errors.InventoryError.NotFound.ToString(),
			})
		}
		return ctx.JSON(ajax.Responses{
			Status:  "error",
			Message: err.Error(),
		})
	}

	domain, err := new(model.Inventory).GetById(domainId, userLogin.Id)
	if err != nil {
		if !utility.IsWindow() {
			return ctx.JSON(ajax.Responses{
				Status:  "error",
				Message: GetLang(ctx).Errors.InventoryError.NotFound.ToString(),
			})
		}
		return ctx.JSON(ajax.Responses{
			Status:  "error",
			Message: err.Error(),
		})
	}
	if domain.Id == 0 {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Mỗi lần load scan lại adstxt
	_ = domain.ScanAdsTxt()
	adsMissingLine, syncError := domain.GetAllMissingAdsTxt()
	assigns.Domain = domain
	assigns.AdsTxtMissingLines = adsMissingLine
	assigns.AdsTxtMissingLineSyncError = syncError
	return ctx.Render("ads_txt/detail", assigns, view.LAYOUTMain)
}

func (t *AdsTxt) Load(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAdsTxtLoad)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignAdsTxt{Schema: assign.Get(ctx)}
	assigns.Title = config.TitleWithPrefix("Ads Txt missing for domain")

	payload := AdsTxtPayload{}
	err := ctx.BodyParser(&payload)
	if err != nil {
		return ctx.JSON(ajax.Responses{
			Status:  "error",
			Message: err.Error(),
		})
	}
	// Get Domain from payload
	domain, err := new(model.Inventory).GetById(payload.DomainId, userLogin.Id)
	if err != nil {
		if !utility.IsWindow() {
			return ctx.JSON(ajax.Responses{
				Status:  "error",
				Message: GetLang(ctx).Errors.InventoryError.NotFound.ToString(),
			})
		}
		return ctx.JSON(ajax.Responses{
			Status:  "error",
			Message: err.Error(),
		})
	}
	adsMissingLine, syncError := domain.GetAllMissingAdsTxt()
	assigns.Domain = domain
	assigns.AdsTxtMissingLines = adsMissingLine
	assigns.AdsTxtMissingLineSyncError = syncError
	html, err := block.Render("ads_txt/load.gohtml", assigns)
	if err != nil {
		return ctx.JSON(ajax.Responses{
			Status:  "error",
			Message: err.Error(),
		})
	}

	lastScanAdsTxt := domain.LastScanAdsTxt.Time.Format("2006-01-02 at 15:04:05")

	return ctx.JSON(ajax.Responses{
		Status: "success",
		DataObject: map[string]interface{}{
			"html":           html.String(),
			"lastScanAdsTxt": lastScanAdsTxt,
		},
	})

}

type AdsTxtPayload struct {
	DomainId int64 `form:"did"`
}

func (t *AdsTxt) Scan(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAdsTxtScan)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get Input
	payload := AdsTxtPayload{}
	err := ctx.BodyParser(&payload)
	if err != nil {
		return ctx.JSON(ajax.Responses{
			Status:  "error",
			Message: err.Error(),
		})
	}
	// Get Domain from payload
	domain, err := new(model.Inventory).GetById(payload.DomainId, userLogin.Id)
	if err != nil {
		if !utility.IsWindow() {
			return ctx.JSON(ajax.Responses{
				Status:  "error",
				Message: GetLang(ctx).Errors.InventoryError.NotFound.ToString(),
			})
		}
		return ctx.JSON(ajax.Responses{
			Status:  "error",
			Message: err.Error(),
		})
	}
	// Scan Ads.txt
	err = domain.ScanAdsTxt()
	if err != nil {
		if !utility.IsWindow() {
			return ctx.JSON(ajax.Responses{
				Status:  "error",
				Message: GetLang(ctx).Errors.InventoryError.ScanAdsTxt.ToString(),
			})
		}
		return ctx.JSON(ajax.Responses{
			Status:  "error",
			Message: err.Error(),
		})
	}
	// Return Success message
	return ctx.JSON(ajax.Responses{
		Status:  "success",
		Message: "All ads.txt lines are in sync!",
	})
}
