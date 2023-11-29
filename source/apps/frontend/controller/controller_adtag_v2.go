package controller

import (
	"errors"
	"fmt"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"source/pkg/helpers"
	"source/pkg/utility"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (t *AdTag) AddV2(ctx *fiber.Ctx) error {
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

	if helpers.IsHaiMode() {
		// assigns.InventoryConfig.TableInventoryConfig.GamAutoCreate = 1
	}
	return ctx.Render("adtag_v2/add", assigns, view.LAYOUTEmpty)
}

func (t *AdTag) AddPostV2(ctx *fiber.Ctx) error {
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

func (t *AdTag) EditV2(ctx *fiber.Ctx) error {
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
	return ctx.Render("adtag_v2/edit", assigns, view.LAYOUTEmpty)
}
