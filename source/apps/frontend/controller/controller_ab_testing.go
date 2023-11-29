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
	"time"
)

type AbTesting struct{}

type AssignAbTesting struct {
	assign.Schema
	Params          payload.AbTestingIndex
	Countries       []model.CountryRecord
	Devices         []model.DeviceRecord
	Data            []model.InventoryRecord
	IsMoreData      bool
	IsMoreDevice    bool
	IsMoreGeography bool
	ListBoxCollapse []string

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

type AssignAbTestingAdd struct {
	assign.Schema
	ListBoxCollapse []string
	Bidders         []model.BidderRecord
	Floors          []model.FloorRecord
	UserIdModule    []model.ModuleUserIdRecord
	EndDate         string
}

type AssignAbTestingEdit struct {
	assign.Schema
	Row             model.AbTestingRecord
	ListBoxCollapse []string
	Bidders         []model.BidderRecord
	Floors          []model.FloorRecord
	UserIdModule    []model.ModuleUserIdRecord
	StartDate       string
	EndDate         string

	AdSizes              []model.AdSizeRecord
	AdFormats            []model.AdTypeRecord
	AdTags               []model.InventoryAdTagRecord
	Inventories          []model.InventoryRecord
	Countries            []model.CountryRecord
	Devices              []model.DeviceRecord
	InventoryOfAbTesting []int64
	AdTagOfAbTesting     []int64
	AdFormatOfAbTesting  []int64
	AdSizeOfAbTesting    []int64
	CountryOfAbTesting   []int64
	DeviceOfAbTesting    []int64
	AdSizesIncluded      []model.AdSizeRecord
	AdFormatsIncluded    []model.AdTypeRecord
	AdTagsIncluded       []model.InventoryAdTagRecord
	InventoriesIncluded  []model.InventoryRecord
	CountriesIncluded    []model.CountryRecord
	DevicesIncluded      []model.DeviceRecord

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

func (t *AbTesting) Index(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAbTesting)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	params := payload.AbTestingIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignAbTesting{Schema: assign.Get(ctx)}
	assigns.Params = params
	assigns.Domain = new(model.LineItem).GetFilter("domain", userLogin.Id, params.Domain)
	assigns.Format = new(model.LineItem).GetFilter("format", userLogin.Id, params.AdFormat)
	assigns.Size = new(model.LineItem).GetFilter("size", userLogin.Id, params.AdSize)
	assigns.AdTag = new(model.LineItem).GetFilter("adtag", userLogin.Id, params.AdTag)
	assigns.Country = new(model.LineItem).GetFilter("country", userLogin.Id, params.Country)
	assigns.Device = new(model.LineItem).GetFilter("device", userLogin.Id, params.Device)
	assigns.Title = config.TitleWithPrefix("A/B Test")
	return ctx.Render("ab_testing/index", assigns, view.LAYOUTMain)
}

func (t *AbTesting) Filter(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAbTesting)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get payload post
	var inputs payload.AbTestingFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Get data from model
	dataTable, err := new(model.AbTesting).GetByFilters(&inputs, userLogin.Id, GetLang(ctx))
	if err != nil {
		return err
	}
	// Print Data to JSON
	return ctx.JSON(dataTable)
}

func (t *AbTesting) Add(ctx *fiber.Ctx) error {
	assigns := AssignAbTestingAdd{Schema: assign.Get(ctx)}
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAbTestingAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	assigns.Bidders = new(model.Bidder).GetAllUser(userLogin.Id)
	assigns.UserIdModule, _ = new(model.ModuleUserId).GetAll()
	layoutDate := "01/02/2006"
	assigns.EndDate = time.Now().AddDate(0, 0, 7).Format(layoutDate)
	collapse := new(model.AbTesting).GetListBoxCollapse(userLogin.Id, 0, "abTesting", "add")
	assigns.ListBoxCollapse = collapse
	assigns.Floors = new(model.Floor).GetFloorDynamic(userLogin.Id)
	assigns.Title = config.TitleWithPrefix("Create A/B Test")
	return ctx.Render("ab_testing/add", assigns, view.LAYOUTMain)
}

func (t *AbTesting) AddPost(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAbTestingAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	//test
	//userLogin := model.UserRecord{TableUser: mysql.TableUser{
	//	Id: 1,
	//}}

	// Get Post Data
	inputs := payload.AbTestingCreate{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	rule, errs := new(model.AbTesting).Create(inputs, userLogin, userAdmin, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = rule
	}
	return ctx.JSON(response)
}

func (t *AbTesting) Edit(ctx *fiber.Ctx) error {
	assigns := AssignAbTestingEdit{Schema: assign.Get(ctx)}
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAbTestingEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	id := ctx.Query("id")
	idAbTesting, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	row, err := new(model.AbTesting).GetById(idAbTesting, userLogin.Id)
	assigns.Row = row
	if err != nil || row.Id == 0 {
		return fmt.Errorf(GetLang(ctx).Errors.AbTestingError.NotFound.ToString())
	}
	collapse := new(model.AbTesting).GetListBoxCollapse(userLogin.Id, idAbTesting, "abTesting", "edit")
	assigns.ListBoxCollapse = collapse
	assigns.Bidders = new(model.Bidder).GetAllUser(userLogin.Id)
	assigns.UserIdModule, _ = new(model.ModuleUserId).GetAll()
	layoutISO := "01/02/2006"
	if row.StartDate.Valid {
		assigns.StartDate = row.StartDate.Time.Format(layoutISO)
	}
	if row.EndDate.Valid {
		assigns.EndDate = row.EndDate.Time.Format(layoutISO)
	}

	targets := new(model.Target).GetTargetAbTesting(idAbTesting)
	mapInventory := make(map[int64]int)
	mapAdFormat := make(map[int64]int)
	mapAdSize := make(map[int64]int)
	mapAdTag := make(map[int64]int)
	mapGeo := make(map[int64]int)
	mapDevice := make(map[int64]int)
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
			mapAdTag[target.TagId] = 1
		}
		if target.GeoId != 0 {
			mapGeo[target.GeoId] = 1
		}
		if target.DeviceId != 0 {
			mapDevice[target.DeviceId] = 1
		}
	}

	// Lọc bỏ những id trùng nhau
	for inventoryId, _ := range mapInventory {
		assigns.InventoryOfAbTesting = append(assigns.InventoryOfAbTesting, inventoryId)
	}
	for adFormatId, _ := range mapAdFormat {
		assigns.AdFormatOfAbTesting = append(assigns.AdFormatOfAbTesting, adFormatId)
	}
	for adSizeId, _ := range mapAdSize {
		assigns.AdSizeOfAbTesting = append(assigns.AdSizeOfAbTesting, adSizeId)
	}
	for adTagId, _ := range mapAdTag {
		assigns.AdTagOfAbTesting = append(assigns.AdTagOfAbTesting, adTagId)
	}
	for geoId, _ := range mapGeo {
		assigns.CountryOfAbTesting = append(assigns.CountryOfAbTesting, geoId)
	}
	for deviceId, _ := range mapDevice {
		assigns.DeviceOfAbTesting = append(assigns.DeviceOfAbTesting, deviceId)
	}
	for _, v := range assigns.InventoryOfAbTesting {
		recInventory, _ := new(model.Inventory).GetById(v, userLogin.Id)
		if recInventory.Id == 0 {
			continue
		}
		assigns.InventoriesIncluded = append(assigns.InventoriesIncluded, recInventory)
	}
	for _, v := range assigns.AdTagOfAbTesting {
		recAdTags := new(model.InventoryAdTag).GetById(v)
		if recAdTags.Id == 0 {
			continue
		}
		assigns.AdTagsIncluded = append(assigns.AdTagsIncluded, recAdTags)
	}
	for _, v := range assigns.AdFormatOfAbTesting {
		recAdFormat := new(model.AdType).GetById(v)
		if recAdFormat.Id == 0 {
			continue
		}
		assigns.AdFormatsIncluded = append(assigns.AdFormatsIncluded, recAdFormat)
	}
	for _, v := range assigns.AdSizeOfAbTesting {
		recAdSize := new(model.AdSize).GetById(v)
		if recAdSize.Id == 0 {
			continue
		}
		assigns.AdSizesIncluded = append(assigns.AdSizesIncluded, recAdSize)
	}
	for _, v := range assigns.CountryOfAbTesting {
		recCountry := new(model.Country).GetById(v)
		if recCountry.Id == 0 {
			continue
		}
		assigns.CountriesIncluded = append(assigns.CountriesIncluded, recCountry)
	}
	for _, v := range assigns.DeviceOfAbTesting {
		recDevice := new(model.Device).GetById(v)
		if recDevice.Id == 0 {
			continue
		}
		assigns.DevicesIncluded = append(assigns.DevicesIncluded, recDevice)
	}

	assigns.Inventories, assigns.IsMoreInventory, assigns.InventoryLastPage = new(model.Inventory).LoadMoreDataPageEdit(userLogin.Id, assigns.InventoryOfAbTesting)
	assigns.AdSizes, assigns.IsMoreAdSize, assigns.SizeLastPage = new(model.AdSize).LoadMoreDataPageEdit(assigns.AdSizeOfAbTesting)
	assigns.AdFormats, assigns.IsMoreAdFormat, assigns.FormatLastPage = new(model.AdType).LoadMoreDataPageEdit(assigns.AdFormatOfAbTesting)
	assigns.AdTags, assigns.IsMoreTag, _, assigns.TagLastPage = new(model.InventoryAdTag).LoadMoreDataPageEdit(payload.FilterTarget{
		Inventory: assigns.InventoryOfAbTesting,
		Format:    assigns.AdFormatOfAbTesting,
		Size:      assigns.AdSizeOfAbTesting,
	}, userLogin.Id, assigns.AdTagOfAbTesting)
	assigns.Countries, assigns.IsMoreGeography, assigns.GeoLastPage = new(model.Country).LoadMoreDataPageEdit(assigns.CountryOfAbTesting)
	assigns.Devices, assigns.IsMoreDevice, assigns.DeviceLastPage = new(model.Device).LoadMoreDataPageEdit(assigns.DeviceOfAbTesting)

	assigns.Floors = new(model.Floor).GetFloorDynamic(userLogin.Id)
	assigns.Title = config.TitleWithPrefix("Edit A/B Test " + row.Name)
	return ctx.Render("ab_testing/edit", assigns, view.LAYOUTMain)
}

func (t *AbTesting) EditPost(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAbTestingEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get Post Data
	inputs := payload.AbTestingCreate{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	rule, errs := new(model.AbTesting).Edit(inputs, userLogin, userAdmin, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = rule
	}
	return ctx.JSON(response)
}

func (t *AbTesting) Delete(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAbTestingDel)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := payload.Delete{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}

	AbTestingModel := new(model.AbTesting)
	notify := AbTestingModel.Delete(inputs.Id, userLogin.Id, userAdmin, GetLang(ctx))
	return ctx.JSON(notify)
}

func (t *AbTesting) Collapse(ctx *fiber.Ctx) (err error) {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAbTestingCollapse)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := model.PageCollapseRecord{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	response := ajax.Responses{}
	inputs.UserId = userLogin.Id
	inputs.PageCollapse = "abTesting"
	errs := new(model.PageCollapse).HandleCollapse(inputs)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
	}
	return ctx.JSON(response)
}
