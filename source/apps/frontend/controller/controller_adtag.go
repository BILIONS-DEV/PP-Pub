package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"source/pkg/utility"
	"strconv"
)

type AdTag struct{}

type AssignAdTagIndex struct {
	assign.Schema
	Params    payload.InventoryAdTagIndex
	Countries []model.CountryRecord
	AdSizes   []model.AdSizeRecord
	AdFormats []model.AdTypeRecord
}

type AssignAdTagAdd struct {
	assign.Schema
	InventoryId        int64
	Inventory          model.InventoryRecord
	InventoryConfig    model.InventoryConfigRecord
	Templates          []model.TemplateRecord
	Playlists          []model.PlaylistRecord
	AdSizes            []model.AdSizeRecord
	AdTypes            []model.AdTypeRecord
	AdSizeSticky       []model.AdSizeRecord
	AdSizeStickyMobile []model.AdSizeRecord
	AdSizeAdditional   []model.AdSizeRecord
	DisabledGam        bool
	ListBoxCollapse    []string
	AdTagDisplay       []model.AdTagRecord
}

func (t *AdTag) Index(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAdTag)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	params := payload.InventoryAdTagIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignAdTagIndex{Schema: assign.Get(ctx)}
	assigns.Countries = new(model.Country).GetAll()
	assigns.AdSizes = new(model.AdSize).GetAll()
	assigns.AdFormats = new(model.AdType).GetAll(userLogin, userAdmin)
	assigns.Title = config.TitleWithPrefix("Add Ad Tag Video")
	return ctx.Render("adtag/index", assigns, view.LAYOUTMain)
}

func (t *AdTag) Filter(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAdTag)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	//pass login
	//userLogin := model.UserRecord{TableUser: mysql.TableUser{
	//	Id: 1,
	//}}
	// Get payload post
	var inputs payload.InventoryAdTagFilterPayload
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Get data from model
	dataTable, err := new(model.InventoryAdTag).GetByFilters(&inputs, userLogin, userAdmin, GetLang(ctx))
	if err != nil {
		return err
	}
	// Print Data to JSON
	return ctx.JSON(dataTable)
}

func (t *AdTag) Add(ctx *fiber.Ctx) error {
	assigns := AssignAdTagAdd{Schema: assign.Get(ctx)}
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAdTagAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	param := ctx.Query("inventoryId")
	if param == "" {
		return errors.New("inventory not found")
	}
	inventoryId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return err
	}
	inventory, err := new(model.Inventory).GetById(inventoryId, userLogin.Id)
	if err != nil || inventory.Id == 0 {
		return fmt.Errorf(GetLang(ctx).Errors.InventoryError.NotFound.ToString())
	}
	cfDomain := new(model.InventoryConfig).GetByInventoryId(inventoryId)
	if cfDomain.GamAutoCreate == 1 {
		assigns.DisabledGam = true
	} else {
		assigns.DisabledGam = false
	}
	template, _ := new(model.Template).GetAll(userLogin.Id)
	playlist, _ := new(model.Playlist).GetAll(userLogin.Id)
	collapse := new(model.AdTag).GetListBoxCollapse(userLogin.Id, 0, "adtag", "add")
	// assigns.AdSizeSticky = new(model.AdSize).GetSizeStickyDesktop()
	assigns.AdSizeSticky = new(model.AdSize).GetSizeDefaultStickyDesktop(1)
	// assigns.AdSizeStickyMobile = new(model.AdSize).GetSizeStickyMobile()
	assigns.AdSizeStickyMobile = new(model.AdSize).GetSizeDefaultStickyMobile()
	assigns.ListBoxCollapse = collapse
	assigns.InventoryId = inventoryId
	assigns.Inventory = inventory
	assigns.AdSizeAdditional = new(model.AdSize).GetAll()
	assigns.InventoryConfig = new(model.InventoryConfig).GetByInventoryId(inventoryId)
	assigns.AdTagDisplay = new(model.AdTag).GetAdTagDisplay(inventoryId)
	assigns.Templates = template
	assigns.Playlists = playlist
	assigns.AdSizes = new(model.AdSize).GetAll()
	assigns.AdTypes = new(model.AdType).GetAll(userLogin, userAdmin)
	assigns.Title = config.TitleWithPrefix("Add Ad Tag")
	return ctx.Render("adtag/add", assigns, view.LAYOUTEmpty)
}

func (t *AdTag) AddPost(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAdTagAdd)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	//userLogin := model.UserRecord{TableUser: mysql.TableUser{
	//	Id: 1,
	//}}
	// Get Post Data
	inputs := payload.AdTagAdd{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	row, errs := new(model.AdTag).Create(inputs, userLogin, userAdmin, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = row
	}
	return ctx.JSON(response)
}

type AssignAdTagEdit struct {
	assign.Schema
	Row                            model.AdTagRecord
	Inventory                      model.InventoryRecord
	InventoryConfig                model.InventoryConfigRecord
	Templates                      []model.TemplateRecord
	Playlists                      []model.PlaylistRecord
	AdSizes                        []model.AdSizeRecord
	AdTypes                        []model.AdTypeRecord
	AdSizeAdditional               []model.AdSizeRecord
	AdSizeAdditionalStickMobile    []model.AdSizeRecord
	AdSizeAdditionalSelected       []int64
	AdSizeAdditionalMobileSelected []int64
	AdSizeSticky                   []model.AdSizeRecord
	AdSizeStickyMobile             []model.AdSizeRecord
	EnableBidOutStream             bool
	DisabledGam                    bool
	ListBoxCollapse                []string
	EnableSizeOnMobile             bool
	AdTagDisplay                   []model.AdTagRecord
}

func (t *AdTag) Edit(ctx *fiber.Ctx) error {
	assigns := AssignAdTagEdit{Schema: assign.Get(ctx)}
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAdTagEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	id := ctx.Query("id")
	idSearch, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	param := ctx.Query("inventoryId")
	if param == "" {
		return errors.New("inventory not found")
	}
	inventoryId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return err
	}
	inventory, err := new(model.Inventory).GetById(inventoryId, userLogin.Id)
	if err != nil || inventory.Id == 0 {
		return fmt.Errorf(GetLang(ctx).Errors.InventoryError.NotFound.ToString())
	}
	assigns.Inventory = inventory
	assigns.InventoryConfig = new(model.InventoryConfig).GetByInventoryId(inventoryId)
	row, err := new(model.AdTag).CheckRecord(idSearch, inventory.Id, userLogin.Id)
	if err != nil || row.Id == 0 {
		return fmt.Errorf(GetLang(ctx).Errors.AdTagError.NotFound.ToString())
	}
	cfDomain := new(model.InventoryConfig).GetByInventoryId(inventoryId)
	if cfDomain.GamAutoCreate == 1 {
		assigns.DisabledGam = true
	} else {
		assigns.DisabledGam = false
	}
	assigns.Row = row
	collapse := new(model.AdTag).GetListBoxCollapse(userLogin.Id, idSearch, "adtag", "edit")
	assigns.ListBoxCollapse = collapse
	template, _ := new(model.Template).GetAll(userLogin.Id)
	playlist, _ := new(model.Playlist).GetAll(userLogin.Id)
	assigns.Templates = template
	assigns.Playlists = playlist
	assigns.AdSizes = new(model.AdSize).GetAll()
	assigns.AdTypes = new(model.AdType).GetAll(userLogin, userAdmin)
	assigns.AdTagDisplay = new(model.AdTag).GetAdTagDisplay(inventoryId)
	if row.Type == mysql.TYPEDisplay {
		recordPrimarySize := new(model.AdSize).GetById(row.PrimaryAdSize)
		if recordPrimarySize.Id != 0 {
			assigns.AdSizeAdditional = new(model.AdSize).GetSizeAdditional(recordPrimarySize, "")
		}
		recordsSizeAdditional := new(model.AdTagSizeAdditional).GetAllByTagId(row.TableInventoryAdTag.Id)
		for _, v := range recordsSizeAdditional {
			if v.Device == 2 {
				assigns.AdSizeAdditionalMobileSelected = append(assigns.AdSizeAdditionalMobileSelected, v.AdSizeId)
			} else {
				assigns.AdSizeAdditionalSelected = append(assigns.AdSizeAdditionalSelected, v.AdSizeId)
			}
		}
	}
	if row.Type == mysql.TYPEStickyBanner {
		recordPrimarySize := new(model.AdSize).GetById(row.PrimaryAdSize)
		if recordPrimarySize.Id != 0 {
			assigns.AdSizeAdditional = new(model.AdSize).GetSizeAdditional(recordPrimarySize, "")
		}
		recordPrimarySizeMobile := new(model.AdSize).GetById(row.PrimaryAdSizeMobile)
		if recordPrimarySizeMobile.Id != 0 {
			assigns.AdSizeAdditionalStickMobile = new(model.AdSize).GetSizeAdditional(recordPrimarySizeMobile, "sticky_mobile")
		}
		recordsSizeAdditional := new(model.AdTagSizeAdditional).GetAllByTagId(row.TableInventoryAdTag.Id)
		for _, v := range recordsSizeAdditional {
			if v.Device == 2 {
				assigns.AdSizeAdditionalMobileSelected = append(assigns.AdSizeAdditionalMobileSelected, v.AdSizeId)
			} else {
				assigns.AdSizeAdditionalSelected = append(assigns.AdSizeAdditionalSelected, v.AdSizeId)
			}
		}
	}
	// assigns.AdSizeSticky = new(model.AdSize).GetSizeStickyDesktop()
	assigns.AdSizeSticky = new(model.AdSize).GetSizeDefaultStickyDesktop(row.TableInventoryAdTag.PositionSticky.Base())
	// assigns.AdSizeStickyMobile = new(model.AdSize).GetSizeStickyMobile()
	assigns.AdSizeStickyMobile = new(model.AdSize).GetSizeDefaultStickyMobile()
	size := new(model.AdSize).GetById(assigns.Row.PrimaryAdSize)
	listSizeEnableBidOutStream := []string{"300x250", "336x280", "300x600", "970x250"}
	if utility.InArray(size.Name, listSizeEnableBidOutStream, false) {
		assigns.EnableBidOutStream = true
	} else {
		assigns.EnableBidOutStream = false
	}
	if size.Width > 320 {
		assigns.EnableSizeOnMobile = true
	} else {
		assigns.EnableSizeOnMobile = false
	}
	assigns.Title = config.TitleWithPrefix("Edit Ad Tag")
	return ctx.Render("adtag/edit", assigns, view.LAYOUTEmpty)
}

func (t *AdTag) EditPost(ctx *fiber.Ctx) error {
	// Get use login
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAdTagEdit)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get Post Data
	inputs := payload.AdTagAdd{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// Handle
	response := ajax.Responses{}
	newBidderAssign, errs := new(model.AdTag).Update(inputs, userLogin, userAdmin, GetLang(ctx))
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
		response.DataObject = newBidderAssign
	}
	return ctx.JSON(response)
}

func (t *AdTag) Delete(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAdTagDel)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := payload.Delete{}
	if err := json.Unmarshal(ctx.Body(), &inputs); err != nil {
		return err
	}
	model := new(model.InventoryAdTag)
	//notify := model.Delete(inputs.Id, userLogin.Id, GetLang(ctx))
	notify := model.Archived(inputs.Id, userLogin.Id, userAdmin)
	return ctx.JSON(notify)
}

func (t *AdTag) GetSizeAdditional(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAdTagGetSizeAdditional)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	id := ctx.Query("id")
	typ := ctx.Query("type")
	idPrimarySize, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	recordPrimarySize := new(model.AdSize).GetById(idPrimarySize)
	recordSizeAdditional := new(model.AdSize).GetSizeAdditional(recordPrimarySize, typ)
	return ctx.JSON(recordSizeAdditional)
}

func (t *AdTag) Collapse(ctx *fiber.Ctx) (err error) {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIAdTagCollapse)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	inputs := model.PageCollapseRecord{}
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	response := ajax.Responses{}
	inputs.UserId = userLogin.Id
	inputs.PageCollapse = "adtag"
	errs := new(model.PageCollapse).HandleCollapse(inputs)
	if len(errs) > 0 {
		response.Status = ajax.ERROR
		response.Errors = errs
	} else {
		response.Status = ajax.SUCCESS
	}
	return ctx.JSON(response)
}
