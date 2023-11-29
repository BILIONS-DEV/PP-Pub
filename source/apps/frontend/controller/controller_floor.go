package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"html/template"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
	"source/pkg/ajax"
	"strconv"
)

type Floor struct{}

type AssignFloor struct {
	assign.Schema
	Params          payload.FloorIndex
	Countries       []model.CountryRecord
	Devices         []model.DeviceRecord
	Data            []model.InventoryRecord
	IsMoreData      bool
	IsMoreDevice    bool
	IsMoreGeography bool
	ListBoxCollapse []string
	Currency        template.HTML
	Symbol          template.HTML

	Domain        []model.SearchLineItem
	Format        []model.SearchLineItem
	Size          []model.SearchLineItem
	AdTag         []model.SearchLineItem
	Country       []model.SearchLineItem
	Device        []model.SearchLineItem
	DomainSearch  []model.SearchLineItem
	FormatSearch  []model.SearchLineItem
	SizeSearch    []model.SearchLineItem
	AdTagSearch   []model.SearchLineItem
	CountrySearch []model.SearchLineItem
	DeviceSearch  []model.SearchLineItem
}

type AssignFloorEdit struct {
	assign.Schema
	Row             model.FloorRecord
	ListBoxCollapse []string
	Currency        template.HTML
	Symbol          template.HTML

	AdSizes             []model.AdSizeRecord
	AdFormats           []model.AdTypeRecord
	AdTags              []model.InventoryAdTagRecord
	Inventories         []model.InventoryRecord
	Countries           []model.CountryRecord
	Devices             []model.DeviceRecord
	InventoryOfFloor    []int64
	AdTagOfFloor        []int64
	AdFormatOfFloor     []int64
	AdSizeOfFloor       []int64
	CountryOfFloor      []int64
	DeviceOfFloor       []int64
	AdSizesIncluded     []model.AdSizeRecord
	AdFormatsIncluded   []model.AdTypeRecord
	AdTagsIncluded      []model.InventoryAdTagRecord
	InventoriesIncluded []model.InventoryRecord
	CountriesIncluded   []model.CountryRecord
	DevicesIncluded     []model.DeviceRecord

	IsMoreInventory   bool
	IsMoreTag         bool
	IsMoreAdSize      bool
	IsMoreAdFormat    bool
	IsMoreDevice      bool
	IsMoreGeography   bool
	InventoryLastPage bool
	TagLastPage       bool
	SizeLastPage      bool
	FormatLastPage    bool
	DeviceLastPage    bool
	GeoLastPage       bool
}

func (t *Floor) Index(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIFloor)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	params := payload.FloorIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignFloor{Schema: assign.Get(ctx)}
	assigns.Params = params
	assigns.Domain = new(model.LineItem).GetFilter("domain", userLogin.Id, params.Domain)
	assigns.Format = new(model.LineItem).GetFilter("format", userLogin.Id, params.AdFormat)
	assigns.Size = new(model.LineItem).GetFilter("size", userLogin.Id, params.AdSize)
	assigns.AdTag = new(model.LineItem).GetFilter("adtag", userLogin.Id, params.AdTag)
	assigns.Country = new(model.LineItem).GetFilter("country", userLogin.Id, params.Country)
	assigns.Device = new(model.LineItem).GetFilter("device", userLogin.Id, params.Device)
	assigns.Title = config.TitleWithPrefix("Floor Rule")
	return ctx.Render("floor/index", assigns, view.LAYOUTMain)
}

func (t *Floor) Filter(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIFloor)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get payload post
	var inputs payload.FloorFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Get data from model
	dataTable, err := new(model.Floor).GetByFilters(&inputs, userLogin.Id, GetLang(ctx))
	if err != nil {
		return err
	}
	// Print Data to JSON
	return ctx.JSON(dataTable)
}

func (t *Floor) Add(ctx *fiber.Ctx) error {
	assigns := AssignFloor{Schema: assign.Get(ctx)}
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIFloorAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns.Title = config.TitleWithPrefix("Add Floor")
	currency, symbol := new(model.Config).GetSymbolCurrencyByUserId(userLogin.Id)
	assigns.Currency = template.HTML(currency)
	assigns.Symbol = template.HTML(symbol)
	collapse := new(model.Floor).GetListBoxCollapse(userLogin.Id, 0, "floor", "add")
	assigns.ListBoxCollapse = collapse
	return ctx.Render("floor/add", assigns, view.LAYOUTMain)
}

func (t *Floor) AddPost(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIFloorAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	//test
	//userLogin := model.UserRecord{TableUser: mysql.TableUser{
	//	Id: 1,
	//}}

	// Get Post Data
	inputs := payload.FloorCreate{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	rule, errs := new(model.Floor).Create(inputs, userLogin, userAdmin)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = rule
	}
	return ctx.JSON(response)
}

func (t *Floor) Edit(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIFloorEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns := AssignFloorEdit{Schema: assign.Get(ctx)}
	id := ctx.Query("id")
	idFloor, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	row, err := new(model.Floor).GetById(idFloor, userLogin.Id)
	assigns.Row = row
	if err != nil || row.Id == 0 {
		return fmt.Errorf(GetLang(ctx).Errors.FloorError.NotFound.ToString())
	}
	collapse := new(model.Floor).GetListBoxCollapse(userLogin.Id, idFloor, "floor", "edit")
	assigns.ListBoxCollapse = collapse

	targets := new(model.Target).GetTargetFloor(idFloor)
	mapInventory := make(map[int64]int)
	mapAdFormat := make(map[int64]int)
	mapAdSize := make(map[int64]int)
	for _, target := range targets {
		if target.InventoryId != 0 {
			mapInventory[target.InventoryId] = 1
		}
		if target.AdFormatId != 0 {
			mapAdFormat[target.AdFormatId] = 1
		}
		if target.AdSizeId != 0 {
			mapAdSize[target.AdSizeId] = 1
		}
		if target.TagId != 0 {
			assigns.AdTagOfFloor = append(assigns.AdTagOfFloor, target.TagId)
		}
		if target.GeoId != 0 {
			assigns.CountryOfFloor = append(assigns.CountryOfFloor, target.GeoId)
		}
		if target.DeviceId != 0 {
			assigns.DeviceOfFloor = append(assigns.DeviceOfFloor, target.DeviceId)
		}
	}

	// Lọc bỏ những id trùng nhau
	for inventoryId, _ := range mapInventory {
		assigns.InventoryOfFloor = append(assigns.InventoryOfFloor, inventoryId)
	}
	for adFormatId, _ := range mapAdFormat {
		assigns.AdFormatOfFloor = append(assigns.AdFormatOfFloor, adFormatId)
	}
	for adSizeId, _ := range mapAdSize {
		assigns.AdSizeOfFloor = append(assigns.AdSizeOfFloor, adSizeId)
	}

	for _, v := range assigns.InventoryOfFloor {
		recInventory, _ := new(model.Inventory).GetById(v, userLogin.Id)
		if recInventory.Id == 0 {
			continue
		}
		assigns.InventoriesIncluded = append(assigns.InventoriesIncluded, recInventory)
	}
	for _, v := range assigns.AdTagOfFloor {
		recAdTags := new(model.InventoryAdTag).GetById(v)
		if recAdTags.Id == 0 {
			continue
		}
		assigns.AdTagsIncluded = append(assigns.AdTagsIncluded, recAdTags)
	}
	for _, v := range assigns.AdFormatOfFloor {
		recAdFormat := new(model.AdType).GetById(v)
		if recAdFormat.Id == 0 {
			continue
		}
		assigns.AdFormatsIncluded = append(assigns.AdFormatsIncluded, recAdFormat)
	}
	for _, v := range assigns.AdSizeOfFloor {
		recAdSize := new(model.AdSize).GetById(v)
		if recAdSize.Id == 0 {
			continue
		}
		assigns.AdSizesIncluded = append(assigns.AdSizesIncluded, recAdSize)
	}
	for _, v := range assigns.CountryOfFloor {
		recCountry := new(model.Country).GetById(v)
		if recCountry.Id == 0 {
			continue
		}
		assigns.CountriesIncluded = append(assigns.CountriesIncluded, recCountry)
	}
	for _, v := range assigns.DeviceOfFloor {
		recDevice := new(model.Device).GetById(v)
		if recDevice.Id == 0 {
			continue
		}
		assigns.DevicesIncluded = append(assigns.DevicesIncluded, recDevice)
	}

	assigns.Inventories, assigns.IsMoreInventory, assigns.InventoryLastPage = new(model.Inventory).LoadMoreDataPageEdit(userLogin.Id, assigns.InventoryOfFloor)
	assigns.AdSizes, assigns.IsMoreAdSize, assigns.SizeLastPage = new(model.AdSize).LoadMoreDataPageEdit(assigns.AdSizeOfFloor)
	assigns.AdFormats, assigns.IsMoreAdFormat, assigns.FormatLastPage = new(model.AdType).LoadMoreDataPageEdit(assigns.AdFormatOfFloor)
	assigns.AdTags, assigns.IsMoreTag, _, assigns.TagLastPage = new(model.InventoryAdTag).LoadMoreDataPageEdit(payload.FilterTarget{
		Inventory: assigns.InventoryOfFloor,
		Format:    assigns.AdFormatOfFloor,
		Size:      assigns.AdSizeOfFloor,
	}, userLogin.Id, assigns.AdTagOfFloor)
	assigns.Countries, assigns.IsMoreGeography, assigns.GeoLastPage = new(model.Country).LoadMoreDataPageEdit(assigns.CountryOfFloor)
	assigns.Devices, assigns.IsMoreDevice, assigns.DeviceLastPage = new(model.Device).LoadMoreDataPageEdit(assigns.DeviceOfFloor)
	curency, symbol := new(model.Config).GetSymbolCurrencyByUserId(userLogin.Id)
	assigns.Currency = template.HTML(curency)
	assigns.Symbol = template.HTML(symbol)
	assigns.Title = config.TitleWithPrefix("Edit Floor")
	return ctx.Render("floor/edit", assigns, view.LAYOUTMain)
}

func (t *Floor) EditPost(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIFloorEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get Post Data
	inputs := payload.FloorCreate{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	rule, errs := new(model.Floor).Edit(inputs, userLogin, userAdmin)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = rule
	}
	return ctx.JSON(response)
}

func (t *Floor) Delete(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIFloorDel)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := payload.Delete{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}

	floorModel := new(model.Floor)
	notify := floorModel.Delete(inputs.Id, userLogin.Id, userAdmin)
	return ctx.JSON(notify)
}

func (t *Floor) Collapse(ctx *fiber.Ctx) (err error) {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIFloorCollapse)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := model.PageCollapseRecord{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	response := ajax.Responses{}
	inputs.UserId = userLogin.Id
	inputs.PageCollapse = "floor"
	errs := new(model.PageCollapse).HandleCollapse(inputs)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
	}
	return ctx.JSON(response)
}
