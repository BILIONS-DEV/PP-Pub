package controller

import (
	"fmt"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
	"source/core/technology/mysql"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (t *Inventory) SetupV2(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIInventorySetup)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	id := ctx.Query("id")
	// start := ctx.Query("start")
	// length := ctx.Query("length")
	idSearch, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	tab, err := strconv.ParseInt(ctx.Query("tab"), 10, 64)
	if err != nil {
		tab = 1
	}
	// bỏ tab 2,3,5 với user permission managed service
	if userLogin.Permission == mysql.UserPermissionManagedService && !userAdmin.IsFound() {
		if tab == 2 || tab == 3 || tab == 5 {
			return ctx.SendStatus(fiber.StatusNotFound)
		}
	}

	assigns := AssignInventory{Schema: assign.Get(ctx)}
	if err := ctx.QueryParser(&assigns.Params); err != nil {
		return err
	}
	params := payload.InventoryAdTagIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	row, err := new(model.Inventory).GetById(idSearch, userLogin.Id)
	if err != nil || row.Id == 0 || row.Status != 1 {
		return fmt.Errorf(GetLang(ctx).Errors.InventoryError.NotFound.ToString())
	}
	cf := new(model.InventoryConfig).GetByInventoryId(idSearch)
	_ = row.ScanAdsTxt()

	module, err := new(model.ModuleUserId).GetAll()
	if err != nil {
		return err
	}
	collapse := new(model.Inventory).GetListBoxCollapse(userLogin.Id, idSearch, "inventory")
	assigns.Tab = tab
	assigns.Row = row
	assigns.Config = cf
	assigns.ModuleUserId = module
	assigns.ParamsAdTag = params
	assigns.ListBoxCollapse = collapse
	assigns.GamNetworks = new(model.GamNetwork).GetByUser(userLogin.Id)
	assigns.AdTypes = new(model.AdType).GetAll(userLogin, userAdmin)
	assigns.Bidders = new(model.InventoryConnectionDemand).GetByInventory(row)
	assigns.CountConnectWaiting = 0
	for _, bidderId := range new(model.RlsBidderSystemInventory).GetListIdBidderApprove(row.Name) {
		bidder := new(model.Bidder).GetByIdNoCheckUser(bidderId)
		if bidder.Id == 0 || bidder.BidderTemplateId == 1 { // Bỏ qua nếu k tìm thấy bidder hoặc là bidder google
			continue
		}
		status := new(model.InventoryConnectionDemand).GetStatus(row.Id, bidder.Id)
		if status == 2 {
			assigns.CountConnectWaiting++
		}
	}

	adsMissingLine, syncError := row.GetAllMissingAdsTxt()
	assigns.AdsTxtMissingLines = adsMissingLine
	assigns.AdsTxtMissingLineSyncError = syncError

	assigns.Title = config.TitleWithPrefix("Setup Inventory")
	assigns.LANG.Title = "Setup Inventory"
	return ctx.Render("supply_v2/setup", assigns, view.LAYOUTMain)
}
