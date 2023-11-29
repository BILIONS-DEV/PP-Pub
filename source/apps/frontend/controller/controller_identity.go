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
	"strconv"
)

type Identity struct{}

type AssignIdentity struct {
	assign.Schema
	Params          payload.IdentityIndex
	Countries       []model.CountryRecord
	Devices         []model.DeviceRecord
	Data            []model.InventoryRecord
	IsMoreData      bool
	IsMoreDevice    bool
	IsMoreGeography bool
	ListBoxCollapse []string

	Domain       []model.SearchLineItem
	DomainSearch []model.SearchLineItem
}

type AssignIdentityAdd struct {
	assign.Schema
	ListBoxCollapse    []string
	ModuleUserId       []model.ModuleUserIdRecord
	ListIdOfInventory  []int64
	TargetAllInventory bool
	HaveAProfile       bool
}

type AssignIdentityEdit struct {
	assign.Schema
	Row                 model.IdentityRecord
	ListBoxCollapse     []string
	ModuleUserId        []model.ModuleUserIdRecord
	ModuleDefault       []payload.ModuleUserIdAssign
	Modules             []payload.ModuleUserIdAssign
	ListIdOfInventory   []int64
	ListModuleIdDefault []int64

	Inventories         []model.InventoryRecord
	InventoryOfIdentity []int64
	InventoriesIncluded []model.InventoryRecord

	IsMoreInventory   bool
	InventoryLastPage bool

	TargetAllInventory bool
	HaveAProfile       bool
}

func (t *Identity) Index(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIIdentity)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	params := payload.IdentityIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignIdentity{Schema: assign.Get(ctx)}
	assigns.Params = params
	assigns.Domain = new(model.LineItem).GetFilter("domain", userLogin.Id, params.Domain)
	assigns.Title = config.TitleWithPrefix("Profile")
	return ctx.Render("identity/index", assigns, view.LAYOUTMain)
}

func (t *Identity) Filter(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIIdentity)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get payload post
	var inputs payload.IdentityFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Get data from model
	dataTable, err := new(model.Identity).GetByFilters(&inputs, userLogin.Id, GetLang(ctx))
	if err != nil {
		return err
	}
	// Print Data to JSON
	return ctx.JSON(dataTable)
}

func (t *Identity) Add(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIIdentityAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignIdentityAdd{Schema: assign.Get(ctx)}

	IdentityModel := new(model.Identity)

	// Get Collapse
	collapse := IdentityModel.GetListBoxCollapse(userLogin.Id, 0, "Identity", "add")
	assigns.ListBoxCollapse = collapse

	// Get Module
	module, err := new(model.ModuleUserId).GetAll()
	if err != nil {
		return err
	}
	assigns.ModuleUserId = module

	// Get list ids profile (identity) target all inventory
	//profileTargetAllInventory := IdentityModel.GetIdsProfileTargetAllInventory(userLogin)
	//if len(profileTargetAllInventory) > 0 {
	//	assigns.TargetAllInventory = true
	//}
	//assigns.HaveAProfile = IdentityModel.IsHaveAProfile(userLogin, 0)

	// Set Title and Render view
	assigns.Title = config.TitleWithPrefix("New Profile")
	return ctx.Render("identity/add", assigns, view.LAYOUTMain)
}

func (t *Identity) AddPost(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIIdentityAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	// Get Post Data
	inputs := payload.IdentityCreate{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	rule, errs := new(model.Identity).Create(inputs, userLogin, userAdmin)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = rule
	}
	return ctx.JSON(response)
}

func (t *Identity) Edit(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIIdentityEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignIdentityEdit{Schema: assign.Get(ctx)}
	id := ctx.Query("id")
	idIdentity, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	row, err := new(model.Identity).GetById(idIdentity, userLogin.Id)
	assigns.Row = row
	if err != nil || row.Id == 0 {
		return fmt.Errorf(GetLang(ctx).Errors.IdentityError.NotFound.ToString())
	}
	collapse := new(model.Identity).GetListBoxCollapse(userLogin.Id, idIdentity, "Identity", "edit")
	assigns.ListBoxCollapse = collapse

	targets := new(model.Target).GetTargetIdentity(idIdentity)
	mapInventory := make(map[int64]int)
	for _, target := range targets {
		if target.InventoryId != 0 {
			mapInventory[target.InventoryId] = 1
		}
	}

	// Lọc bỏ những id trùng nhau
	for inventoryId, _ := range mapInventory {
		assigns.InventoryOfIdentity = append(assigns.InventoryOfIdentity, inventoryId)
	}

	for _, v := range assigns.InventoryOfIdentity {
		recInventory, _ := new(model.Inventory).GetById(v, userLogin.Id)
		if recInventory.Id == 0 {
			continue
		}
		assigns.InventoriesIncluded = append(assigns.InventoriesIncluded, recInventory)
	}
	assigns.Inventories, assigns.IsMoreInventory, assigns.InventoryLastPage = new(model.Inventory).LoadMoreDataPageEdit(userLogin.Id, assigns.InventoryOfIdentity)

	module, err := new(model.ModuleUserId).GetAll()
	if err != nil {
		return err
	}

	// Xử lý mudule user id
	moduleUsers, listIdOfInventory, _ := new(model.IdentityModuleInfo).GetModules(idIdentity)
	var moduleUserId payload.ModuleUserIdAssign
	var listModule []payload.ModuleUserIdAssign
	for _, val := range moduleUsers {
		moduleUserId.Id = val.Id
		moduleUserId.IdentityId = val.IdentityId
		moduleUserId.Name = val.Name
		moduleUserId.ModuleId = val.ModuleId
		moduleUserId.Params, _ = new(model.ModuleUserId).ConvertParamToList(val.Params)
		moduleUserId.Storage, _ = new(model.ModuleUserId).ConvertStorageToList(val.Storage)
		moduleUserId.AbTesting = val.AbTesting
		moduleUserId.Volume = val.Volume
		listModule = append(listModule, moduleUserId)
	}
	assigns.ModuleUserId = module
	assigns.Modules = listModule
	assigns.ListIdOfInventory = listIdOfInventory
	assigns.ListModuleIdDefault = []int64{1, 3, 4, 5}

	// Get list ids profile (identity) target all inventory
	//IdentityModel := new(model.Identity)
	//profileTargetAllInventory := IdentityModel.GetIdsProfileTargetAllInventory(userLogin)
	//if len(profileTargetAllInventory) > 0 && !utility.InArray(row.Id, profileTargetAllInventory, true) {
	//	assigns.TargetAllInventory = true
	//}
	//assigns.HaveAProfile = IdentityModel.IsHaveAProfile(userLogin, row.Id)

	assigns.Title = config.TitleWithPrefix("Edit Profile")
	return ctx.Render("identity/edit", assigns, view.LAYOUTMain)
}

func (t *Identity) EditPost(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIIdentityEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get Post Data
	inputs := payload.IdentityCreate{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	rule, errs := new(model.Identity).Edit(inputs, userLogin, userAdmin)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = rule
	}
	return ctx.JSON(response)
}

func (t *Identity) Delete(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIIdentityDel)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := payload.Delete{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}

	IdentityModel := new(model.Identity)
	notify := IdentityModel.Delete(inputs.Id, userLogin.Id, userAdmin, GetLang(ctx))
	return ctx.JSON(notify)
}

func (t *Identity) Collapse(ctx *fiber.Ctx) (err error) {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIIdentityCollapse)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := model.PageCollapseRecord{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	response := ajax.Responses{}
	inputs.UserId = userLogin.Id
	inputs.PageCollapse = "Identity"
	errs := new(model.PageCollapse).HandleCollapse(inputs)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
	}
	return ctx.JSON(response)
}

func (t *Identity) ChangeStatus(ctx *fiber.Ctx) (err error) {
	return ctx.JSON(fiber.Map{})
}
