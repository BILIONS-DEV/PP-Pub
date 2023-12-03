package controllerv2

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/view"
	"source/internal/entity/dto"
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
	inventoryUC "source/internal/usecase/inventory"
	"strconv"
)

type inventory struct {
	*handler
}

func (h *handler) InitRoutesInventory(app fiber.Router) {
	this := inventory{h}
	app.Get(config.URIInventory, this.Index)
	app.Post(config.URIInventory, this.IndexPost)
	app.Get(config.URIInventorySetup, this.Setup)
	app.Post(config.URIInventorySetup, this.SetupPost)
}

type AssignInventory struct {
	Assign
	Params dto.PayloadInventoryIndex
}

func (t *inventory) Index(ctx *fiber.Ctx) (err error) {
	params := dto.PayloadInventoryIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignInventory{Assign: newAssign(ctx, "Supply")}
	// assigns
	assigns.Params = params

	return ctx.Render("supply/index", assigns, view.LAYOUTMain)
}

func (t *inventory) IndexPost(ctx *fiber.Ctx) (err error) {
	// get user login
	userLogin := getUserLogin(ctx)
	var payload dto.PayloadInventoryIndexPost
	if err = ctx.BodyParser(&payload); err != nil {
		fmt.Printf("%+v\n", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.MakeResponseError(err),
		)
	}
	payload.UserID = userLogin.ID
	fmt.Println("data: ", userLogin.ID)
	if errs := payload.Validate(); len(errs) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.MakeResponseErrorWithID(errs...),
		)
	}
	// fmt.Printf("%+v \n", payload)
	var inputUC = inventoryUC.InputFilterUC{
		Request:     payload.Request,
		UserID:      userLogin.ID,
		QuerySearch: payload.PostData.QuerySearch,
		Status:      payload.PostData.Status,
		Type:        payload.PostData.Type,
		AdsSync:     payload.PostData.AdsSync,
	}
	var records []*model.InventoryModel
	var totalRecords int64
	if totalRecords, records, err = t.useCases.Inventory.Filter(&inputUC); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			dto.Fail(err),
		)
	}

	return ctx.JSON(
		datatable.Response{
			Draw:            payload.Draw,
			RecordsTotal:    totalRecords,
			RecordsFiltered: totalRecords,
			Data:            t.makeResponseDatatable(records),
		},
	)
}

func (t *inventory) makeResponseDatatable(inventories []*model.InventoryModel) (records []dto.ResponseInventoryDatatable) {
	for _, inventory := range inventories {
		rec := dto.ResponseInventoryDatatable{
			InventoryModel: inventory,
			RowId:          strconv.FormatInt(inventory.ID, 10),
			Name:           t.block.RenderToString("supply/index/block.name.gohtml", inventory),
			Type:           inventory.Type.String(),
			Status:         t.block.RenderToString("supply/index/block.status.gohtml", inventory),
			Live:           t.block.RenderToString("supply/index/block.live.gohtml", inventory),
			SyncAdsTxt:     t.block.RenderToString("supply/index/block.sync_ads_txt.gohtml", inventory),
			Action:         t.block.RenderToString("supply/index/block.action.gohtml", inventory),
		}
		records = append(records, rec)
	}
	return
}

type AssignInventorySetup struct {
	Assign
	AdTagParams                dto.PayloadInventorySetup
	Row                        *model.InventoryModel
	GamNetworks                []model.GamNetworkModel
	AdTypes                    []model.AdTypeModel
	CountConnectWaiting        int
	AdsTxtMissingLineSyncError string
	AdsTxtMissingLines         []model.TableMissingAdsTxtModel
}

func (t *inventory) Setup(ctx *fiber.Ctx) (err error) {
	userLogin := getUserLogin(ctx)

	params := dto.PayloadInventorySetup{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}

	// bỏ tab 2,3,5 với user permission managed service
	if userLogin.Permission == model.UserPermissionManagedService && (params.Tab == 2 || params.Tab == 3 || params.Tab == 5) {
		return ctx.Status(fiber.StatusNotFound).SendString("permission denied")
	}

	assigns := AssignInventorySetup{Assign: newAssign(ctx, "Setup Inventory")}
	// assigns
	assigns.AdTagParams = params
	// Get inventory
	inventory, _ := t.useCases.Inventory.GetById(params.Id)
	assigns.Row = inventory
	assigns.GamNetworks, _ = t.useCases.GamNetwork.GetByUser(userLogin.ID)
	assigns.AdTypes, _ = t.useCases.AdType.GetByUser(userLogin.ID)

	assigns.CountConnectWaiting = t.useCases.Inventory.CountConnectWaiting(inventory)

	return ctx.Render("supply/setup", assigns, view.LAYOUTMain)
}

func (t *inventory) SetupPost(ctx *fiber.Ctx) (err error) {
	return nil
}
