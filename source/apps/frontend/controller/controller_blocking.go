package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
	"source/pkg/ajax"
	"source/pkg/utility"
	"strconv"
)

type Blocking struct{}

type AssignBlockingIndex struct {
	assign.Schema
	Params          payload.BlockingIndex
	ListBoxCollapse []string
}

func (t *Blocking) Index(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIBlocking)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	params := payload.BlockingIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignBlockingIndex{Schema: assign.Get(ctx)}
	assigns.Params = params
	assigns.Title = config.TitleWithPrefix("Blocking")
	assigns.LANG.Title = "Blocking"
	return ctx.Render("blocking/index", assigns, view.LAYOUTMain)
}

func (t *Blocking) Filter(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIBlocking)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get payload post
	var inputs payload.BlockingFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Get data from model
	dataTable, err := new(model.Blocking).GetByFilters(&inputs, userLogin.Id, GetLang(ctx))
	if err != nil {
		return err
	}
	// Print Data to JSON
	return ctx.JSON(dataTable)
}

type AssignBlockingAdd struct {
	assign.Schema
	Inventories     []model.InventoryRecord
	ListBoxCollapse []string
}

func (t *Blocking) Add(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIBlockingAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignBlockingAdd{Schema: assign.Get(ctx)}
	assigns.Inventories = new(model.Inventory).GetByUser(userLogin.Id)
	assigns.Title = config.TitleWithPrefix("Add Blocking")
	return ctx.Render("blocking/add", assigns, view.LAYOUTMain)
}

func (t *Blocking) AddPost(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIBlockingAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	//userLogin := model.UserRecord{TableUser: mysql.TableUser{
	//	Id: 1,
	//}}
	// Get Post Data
	inputs := payload.BlockingAdd{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	row, errs := new(model.Blocking).Create(inputs, userLogin, userAdmin)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = row
	}
	return ctx.JSON(response)
}

type AssignBlockingEdit struct {
	assign.Schema
	Row                     model.BlockingRecord
	ListBoxCollapse         []string
	Inventories             []model.InventoryRecord
	ListInventoryIdSelected []int64
	ListAdvertiseDomains    []model.BlockingRestrictionsRecord
	ListCreativeId          []model.BlockingRestrictionsRecord
}

func (t *Blocking) Edit(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIBlockingEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	id, err := strconv.ParseInt(ctx.Query("id"), 10, 64)
	if err != nil {
		return err
	}
	assigns := AssignBlockingEdit{Schema: assign.Get(ctx)}
	assigns.Row, err = new(model.Blocking).GetById(id, userLogin.Id)
	if err != nil || assigns.Row.Id == 0 {
		return fmt.Errorf(GetLang(ctx).Errors.BlockingError.NotFound.ToString())
	}
	assigns.Inventories = new(model.Inventory).GetByUser(userLogin.Id)
	recordRlBlockingInventory, err := new(model.RlBlockingInventory).GetByBlockingId(assigns.Row.Id)
	if err != nil {
		return err
	}
	for _, v := range recordRlBlockingInventory {
		assigns.ListInventoryIdSelected = append(assigns.ListInventoryIdSelected, v.InventoryId)
	}
	assigns.ListAdvertiseDomains, err = new(model.BlockingRestrictions).GetAdvertiseDomainByBlocking(assigns.Row.Id)
	if err != nil {
		return err
	}
	assigns.ListCreativeId, err = new(model.BlockingRestrictions).GetCreativeIdByBlocking(assigns.Row.Id)
	if err != nil {
		return err
	}
	assigns.Title = config.TitleWithPrefix("Edit Blocking")
	return ctx.Render("blocking/edit", assigns, view.LAYOUTMain)
}

func (t *Blocking) EditPost(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIBlockingEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	//userLogin := model.UserRecord{TableUser: mysql.TableUser{
	//	Id: 1,
	//}}
	// Get Post Data
	inputs := payload.BlockingAdd{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	row, errs := new(model.Blocking).Update(inputs, userLogin, userAdmin, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = row
	}
	return ctx.JSON(response)
}

func (t *Blocking) LoadInventory(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIBlockingLoadInventory)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	var recInventory []model.InventoryRecord
	keyword := ctx.Query("keyword")
	if len(keyword) == 0 {
		recInventory, _ = new(model.Inventory).GetByUserIdLimit(userLogin.Id, 10)
	} else {
		recInventory, _ = new(model.Inventory).GetByUserIdSearch(userLogin.Id, keyword)
	}
	type Response struct {
		Id   int64  `json:"id"`
		Name string `json:"name"`
	}
	var response []Response
	for _, v := range recInventory {
		response = append(response, Response{
			Id:   v.Id,
			Name: v.Name,
		})
	}
	return ctx.JSON(response)
}

func (t *Blocking) Delete(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIBlockingDelete)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := payload.Delete{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}

	blocking := new(model.Blocking)
	notify := blocking.Delete(inputs.Id, userLogin.Id, userAdmin, GetLang(ctx))
	return ctx.JSON(notify)
}

func (t *Blocking) ValidateDomain(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIBlockingValidateDomain)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	type ValidateDomain struct {
		ListDomain []string `json:"listDomain"`
	}
	inputs := ValidateDomain{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}
	type Response struct {
		ListError []string `json:"listError"`
		ListValid []string `json:"listValid"`
	}
	res := Response{}
	for _, domain := range inputs.ListDomain {
		if domain == "" {
			continue
		}
		if !utility.ValidateDomainName(domain) {
			res.ListError = append(res.ListError, domain)
		} else {
			res.ListValid = append(res.ListValid, domain)
		}
	}
	return ctx.JSON(res)
}
