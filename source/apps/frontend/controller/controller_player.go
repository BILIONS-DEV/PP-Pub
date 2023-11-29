package controller

import (
	"encoding/json"
	"fmt"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Player struct{}

type AssignPlayer struct {
	assign.Schema
	Params payload.PlayerIndex
}

type AssignPlayerAdd struct {
	assign.Schema
	AdTypes         []model.AdTypeRecord
	Sizes           []model.AdSizeRecord
	Config          payload.PlayerConfig
	ListBoxCollapse []string
}

type AssignPlayerEdit struct {
	assign.Schema
	Row             model.PlayerRecord
	Sizes           []model.AdSizeRecord
	AdTypes         []model.AdTypeRecord
	Config          payload.PlayerConfig
	ListBoxCollapse []string
}

func (t *Player) Index(ctx *fiber.Ctx) error {
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
	return ctx.Render("player/index", assigns, view.LAYOUTMain)
}

func (t *Player) Filter(ctx *fiber.Ctx) error {
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
	dataTable, err := new(model.Template).GetByFilters(&inputs, userLogin, GetLang(ctx))
	if err != nil {
		return err
	}
	// Print Data to JSON
	return ctx.JSON(dataTable)
}

func (t *Player) AddTemplate(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlayerAddTemplate)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignPlayerAdd{Schema: assign.Get(ctx)}
	assigns.AdTypes = new(model.AdType).GetAllTypeVideo()
	assigns.Sizes = new(model.AdSize).GetAllSizeForNative()
	assigns.TemplatePath = "/player/template/add"
	collapse := new(model.Player).GetListBoxCollapse(userLogin.Id, 0, "template", "add")
	assigns.ListBoxCollapse = collapse
	assigns.Title = config.TitleWithPrefix("Add Template")
	return ctx.Render("player/add", assigns, view.LAYOUTTemplate)
}

func (t *Player) AddPost(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlayerAddTemplate)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	//userLogin := model.UserRecord{TableUser: mysql.TableUser{
	//	Id: 1,
	//}}
	// Get Post Data
	inputs := payload.TemplateCreate{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	data, errs := new(model.Player).Create(inputs, userLogin, userAdmin)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = data
	}
	return ctx.JSON(response)
}

func (t *Player) Edit(ctx *fiber.Ctx) error {
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
	return ctx.Render("player/edit", assigns, view.LAYOUTTemplate)
}

func (t *Player) EditPost(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlayerEditTemplate)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	//userLogin := model.UserRecord{TableUser: mysql.TableUser{
	//	Id: 1,
	//}}
	// Get Post Data
	inputs := payload.TemplateCreate{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	data, errs := new(model.Player).Edit(inputs, userLogin, userAdmin)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		configTemplate := payload.PlayerConfig{}
		configTemplate.Contents = new(model.Player).MakeContents(data.Type)
		configTemplate.Template = new(model.Player).MakeTemplate(data)
		configTemplate.TopArticle = new(model.Player).MakeTopArticle()
		configTemplate.Info = new(model.Player).MakeInfo()
		response.Status = ajax.SUCCESS
		response.DataObject = configTemplate
	}
	return ctx.JSON(response)
}

func (t *Player) View(ctx *fiber.Ctx) error {
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
	assigns.AdTypes = new(model.AdType).GetAllTypeVideo()
	//var contentTemplate model.ContentRecord
	//mysql.Client.Where("id = 0").Find(&contentTemplate)
	configTemplate := payload.PlayerConfig{}
	configTemplate.Contents = new(model.Player).MakeContents(row.Type)
	configTemplate.Template = new(model.Player).MakeTemplate(row)
	configTemplate.TopArticle = new(model.Player).MakeTopArticle()
	configTemplate.Info = new(model.Player).MakeInfo()
	assigns.Config = configTemplate
	assigns.Title = config.TitleWithPrefix("View template")
	return ctx.Render("player/view", assigns, view.LAYOUTTemplate)
}

func (t *Player) Delete(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlayerDelTemplate)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := payload.Delete{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}
	m := new(model.Player)
	notify := m.Delete(inputs.Id, userLogin.Id, userAdmin, GetLang(ctx))
	return ctx.JSON(notify)
}

func (t *Player) Duplicate(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlayerDuplicateTemplate)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := payload.Delete{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return fmt.Errorf(GetLang(ctx).Errors.TemplateError.Duplicate.ToString())
	}
	notify := new(model.Player).Duplicate(inputs.Id, userLogin)
	return ctx.JSON(notify)
}

func (t *Player) Collapse(ctx *fiber.Ctx) (err error) {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPlayerCollapse)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := model.PageCollapseRecord{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	response := ajax.Responses{}
	inputs.UserId = userLogin.Id
	inputs.PageCollapse = "template"
	errs := new(model.PageCollapse).HandleCollapse(inputs)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
	}
	return ctx.JSON(response)
}

func (t *Player) Preview(ctx *fiber.Ctx) (err error) {
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
	configTemplate.TopArticle = new(model.Player).MakeTopArticle()
	configTemplate.Info = new(model.Player).MakeInfo()
	return ctx.JSON(configTemplate)
}
