package controller

import (
	"fmt"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
	"source/core/technology/mysql"
	"source/pkg/helpers"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const viewDir = "player-v2"

type PlayerV2 struct{}

func (t *PlayerV2) Index(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlayerTemplate)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	params := payload.PlayerIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignPlayer{Schema: assign.Get(ctx)}
	assigns.Params = params
	assigns.Title = config.TitleWithPrefix("Template")
	return ctx.Render(viewDir+"/index", assigns, view.LAYOUTMain)
}

func (t *PlayerV2) Filter(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlayerTemplate)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get payload post
	var inputs payload.PlayerFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Get data from model
	dataTable, err := new(model.Template).GetByFiltersV2(&inputs, userLogin, GetLang(ctx))
	if err != nil {
		return err
	}
	// Print Data to JSON
	return ctx.JSON(dataTable)
}

func (t *PlayerV2) AddTemplate(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlayerAddTemplate)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignPlayerAdd{Schema: assign.Get(ctx)}
	assigns.AdTypes = new(model.AdType).GetAllTypeVideo()
	assigns.Sizes = new(model.AdSize).GetAllSizeForNative()
	assigns.TemplatePath = "/player-v2/template/add"
	collapse := new(model.Player).GetListBoxCollapse(userLogin.Id, 0, "template", "add")
	assigns.ListBoxCollapse = collapse
	assigns.Title = config.TitleWithPrefix("Add Template")

	if helpers.IsHaiMode() {
		// assigns.AdTypes = append(assigns.AdTypes,
		// 	model.AdTypeRecord{TableAdType: mysql.TableAdType{Id: 9, Name: "Video"}})
	}
	return ctx.Render(viewDir+"/add", assigns, view.LAYOUTTemplateV2)
}

func (t *PlayerV2) Edit(ctx *fiber.Ctx) error {
	assigns := AssignPlayerEdit{Schema: assign.Get(ctx)}
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlayerEditTemplate)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	id := ctx.Query("id")
	idSearch, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	row, err := new(model.Player).GetById(idSearch, userLogin.Id)
	if err != nil || row.Id == 0 {
		return fmt.Errorf(GetLang(ctx).Errors.TemplateError.NotFound.ToString())
	}
	assigns.Row = row
	assigns.AdTypes = new(model.AdType).GetAllTypeVideo()
	assigns.Sizes = new(model.AdSize).GetAllSizeForNative()
	//var contentTemplate model.ContentRecord
	//mysql.Client.Where("id = 0").Find(&contentTemplate)
	configTemplate := payload.PlayerConfig{}
	configTemplate.Contents = new(model.Player).MakeContents(row.Type)
	configTemplate.Template = new(model.Player).MakeTemplate(row)
	configTemplate.TopArticle = new(model.Player).MakeTopArticle()
	configTemplate.Info = new(model.Player).MakeInfo()
	assigns.Config = configTemplate
	collapse := new(model.Player).GetListBoxCollapse(userLogin.Id, idSearch, "template", "edit")
	assigns.ListBoxCollapse = collapse
	assigns.Title = config.TitleWithPrefix("Edit Template")

	return ctx.Render(viewDir+"/edit", assigns, view.LAYOUTTemplateV2)
}

func (t *PlayerV2) Preview(ctx *fiber.Ctx) (err error) {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlayerPreview)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := payload.TemplateCreate{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	templateRecord := model.PlayerRecord{}

	templateRecord.MakeRow(inputs, userLogin)
	if inputs.Id != 0 {
		tempOld, _ := new(model.Player).GetById(inputs.Id, userLogin.Id)
		templateRecord.Type = tempOld.Type
	}
	configTemplate := payload.PlayerConfig{}
	configTemplate.Contents = new(model.Player).MakeContents(templateRecord.Type)
	configTemplate.Template = new(model.Player).MakeTemplate(templateRecord)
	// configTemplate.TopArticle = new(model.Player).MakeTopArticle()
	configTemplate.Info = new(model.Player).MakeInfo()
	return ctx.JSON(configTemplate)
}

func (t *PlayerV2) View(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlayerViewTemplate)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignPlayerEdit{Schema: assign.Get(ctx)}
	id := ctx.Query("id")
	idSearch, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	row, err := new(model.Template).GetTemplateDefaultById(idSearch)
	if err != nil || row.Id == 0 || row.IsDefault != mysql.TypeOn {
		return fmt.Errorf(GetLang(ctx).Errors.TemplateError.NotFound.ToString())
	}
	assigns.Row = row
	// assigns.AdTypes = new(model.AdType).GetAllTypeVideo()
	//var contentTemplate model.ContentRecord
	//mysql.Client.Where("id = 0").Find(&contentTemplate)
	configTemplate := payload.PlayerConfig{}
	configTemplate.Contents = new(model.Player).MakeContents(row.Type)
	configTemplate.Template = new(model.Player).MakeTemplate(row)
	configTemplate.TopArticle = new(model.Player).MakeTopArticle()
	configTemplate.Info = new(model.Player).MakeInfo()
	assigns.Config = configTemplate
	assigns.Title = config.TitleWithPrefix("View template")
	return ctx.Render(viewDir+"/view", assigns, view.LAYOUTTemplateV2)
}
