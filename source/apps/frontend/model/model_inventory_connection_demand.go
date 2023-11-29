package model

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"source/apps/frontend/lang"
	"source/apps/frontend/payload"
	"source/apps/history"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"source/pkg/datatable"
	"source/pkg/htmlblock"
	"source/pkg/pagination"
	"source/pkg/utility"
	"strconv"
	"strings"
)

type InventoryConnectionDemand struct{}

type InventoryConnectionDemandRecord struct {
	mysql.TableInventoryConnectionDemand
}

func (InventoryConnectionDemandRecord) TableName() string {
	return mysql.Tables.InventoryConnectionDemand
}

func (t *InventoryConnectionDemand) GetStatus(inventoryId, bidderId int64) (status mysql.TYPEStatusInventoryConnectionDemand) {
	var inventoryConnectionDemandRecord InventoryConnectionDemandRecord
	mysql.Client.Where("bidder_id = ? and inventory_id = ?", bidderId, inventoryId).Find(&inventoryConnectionDemandRecord)
	if inventoryConnectionDemandRecord.Id > 0 {
		status = inventoryConnectionDemandRecord.Status
	} else {
		status = mysql.TYPEStatusConnectionDemandWaiting
	}
	return
}

func (t *InventoryConnectionDemand) GetConnection(inventoryId, bidderId int64) (record InventoryConnectionDemandRecord) {
	mysql.Client.Where("bidder_id = ? and inventory_id = ?", bidderId, inventoryId).Find(&record)
	return
}

func (t *InventoryConnectionDemand) GetByFilters(inputs *payload.InventoryConnectionDemandFilterPayload, user UserRecord) (response datatable.Response, err error) {
	// Lấy theo list bidder domain được approve từ BE
	inventory, _ := new(Inventory).GetById(inputs.PostData.InventoryId, user.Id)
	var bidders []BidderRecord
	var total int64
	err = mysql.Client.Select("bidder.*, rls_bidder_system_inventory.status").Where("rls_bidder_system_inventory.inventory_name = ? and rls_bidder_system_inventory.status = 3", inventory.Name).
		Joins("inner join rls_bidder_system_inventory on rls_bidder_system_inventory.bidder_id = bidder.id").
		Scopes(
			t.SetFilterStatus(inputs),
			t.setFilterSearch(inputs),
		).
		Model(&bidders).Count(&total).
		Scopes(
			t.setOrder(inputs),
			pagination.Paginate(pagination.Params{
				Limit:  inputs.Length,
				Offset: inputs.Start,
			}),
		).
		Find(&bidders).Error
	if err != nil {
		if !utility.IsWindow() {
			err = fmt.Errorf(lang.Translate.Errors.AdTagError.List.ToString())
		}
		return datatable.Response{}, err
	}
	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.MakeResponseDatatable(bidders, inventory)
	return
}

type InventoryConnectionDemandRecordDatatable struct {
	BidderRecord
	RowId      string `json:"DT_RowId"`
	Demand     string `json:"demand"`
	DemandName string `json:"demand_name"`
	Status     string `json:"status"`
	Action     string `json:"action"`
}

func (t *InventoryConnectionDemand) MakeResponseDatatable(bidders []BidderRecord, inventory InventoryRecord) (records []InventoryConnectionDemandRecordDatatable) {
	for _, bidder := range bidders {
		if bidder.BidderTemplateId == 1 { // Bỏ qua bidder google
			continue
		}
		bidderTemplate := new(BidderTemplate).GetById(bidder.BidderTemplateId)
		status := new(InventoryConnectionDemand).GetStatus(inventory.Id, bidder.Id)
		rec := InventoryConnectionDemandRecordDatatable{
			BidderRecord: bidder,
			RowId:        strconv.FormatInt(bidder.Id, 10),
			Demand:       htmlblock.Render("supply/connection/block.demand.gohtml", fiber.Map{"LinkLogo": bidderTemplate.Logo}).String(),
			DemandName:   htmlblock.Render("supply/connection/block.demand_name.gohtml", bidder).String(),
			Status:       htmlblock.Render("supply/connection/block.status.gohtml", fiber.Map{"Status": status, "BidderId": bidder.Id}).String(),
			Action:       htmlblock.Render("supply/connection/block.action.gohtml", fiber.Map{"Status": status, "BidderId": bidder.Id, "InventoryId": inventory.Id}).String(),
		}
		records = append(records, rec)
	}
	return
}

func (t *InventoryConnectionDemand) SetFilterStatus(inputs *payload.InventoryConnectionDemandFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.Status != nil {
			switch inputs.PostData.Status.(type) {
			case string, int:
				if inputs.PostData.Status != "" {
					return db.Where("status = ?", inputs.PostData.Status)
				}
			case []string, []interface{}:
				return db.Where("status IN ?", inputs.PostData.Status)
			}
		}
		return db
	}
}

func (t *InventoryConnectionDemand) setFilterSearch(inputs *payload.InventoryConnectionDemandFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var flag bool
		// Search from form of datatable <- not use
		if inputs.Search != nil && inputs.Search.Value != "" {
			flag = true
		}
		// Search from form filter
		if inputs.PostData.QuerySearch != "" {
			flag = true
		}
		if !flag {
			return db
		}
		return db.Where("bidder_code LIKE ?", "%"+inputs.PostData.QuerySearch+"%")
	}
}

func (t *InventoryConnectionDemand) setOrder(inputs *payload.InventoryConnectionDemandFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(inputs.Order) > 0 {
			var orders []string
			for _, order := range inputs.Order {
				column := inputs.Columns[order.Column]
				if column.Data == "demand" { // Trường hợp order column 0 thì order theo thời gian update và id tạo trước hay sau
					column.Data = "rls_bidder_system_inventory.updated_at"
					orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
					column.Data = "rls_bidder_system_inventory.id"
					orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
					continue
				}
				if column.Data == "demand_name" {
					column.Data = "bidder_code"
				}
				orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
			}
			orderString := strings.Join(orders, ", ")
			return db.Order(orderString)
		}
		return db
	}
}

func (t *InventoryConnectionDemand) ChangeStatus(inputs payload.PayloadInventoryChangeStatusConnection, user UserRecord, userAdmin UserRecord) (resp ajax.Responses) {
	// Validate
	if inputs.InventoryId < 1 {
		resp.Status = "error"
		resp.Errors = append(resp.Errors, ajax.Error{
			Id:      "",
			Message: "Inventory is required",
		})
		return
	}
	if inputs.BidderId < 1 {
		resp.Status = "error"
		resp.Errors = append(resp.Errors, ajax.Error{
			Id:      "",
			Message: "Bidder is required",
		})
		return
	}
	var err error
	// Tiến hành đổi status kiểm tra xem record cho inventory connection đã được tạo chưa nếu đã tạo thì update status nếu chưa thì tạo mới
	record := t.GetConnection(inputs.InventoryId, inputs.BidderId)
	connectionDemandOld := record
	if record.Id > 0 {
		err = mysql.Client.Model(&InventoryConnectionDemandRecord{}).Where("inventory_id = ? and bidder_id = ?", inputs.InventoryId, inputs.BidderId).Update("status", inputs.Status).Error
	} else {
		err = mysql.Client.Create(&InventoryConnectionDemandRecord{mysql.TableInventoryConnectionDemand{
			InventoryId: inputs.InventoryId,
			BidderId:    inputs.BidderId,
			Status:      inputs.Status,
			UserId:      user.Id,
		}}).Error
		// Nếu như chưa có record tạo một rls waiting mặc định để push history
		connectionDemandOld = InventoryConnectionDemandRecord{mysql.TableInventoryConnectionDemand{
			UserId:      user.Id,
			InventoryId: inputs.InventoryId,
			BidderId:    inputs.BidderId,
			Status:      mysql.TYPEStatusConnectionDemandWaiting,
		}}
	}
	if err != nil {
		resp.Status = "error"
		if utility.IsWindow() {
			resp.Errors = append(resp.Errors, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
			return
		} else {
			resp.Errors = append(resp.Errors, ajax.Error{
				Id:      "",
				Message: "Error!",
			})
			return
		}
	}
	// Get bidder system
	bidderSystem := new(Bidder).GetById(inputs.BidderId, 0)
	if bidderSystem.Id == 0 {
		resp.Status = "error"
		resp.Errors = append(resp.Errors, ajax.Error{
			Id:      "",
			Message: "Error!",
		})
		return
	}
	// Get inventory
	inventory, _ := new(Inventory).GetById(inputs.InventoryId, user.Id)

	// Nếu bidder là apacdex tiến hành xử lý line system cho apacdex cho từng inventory của user
	if bidderSystem.BidderCode == "apacdex" {
		// Get line system apacdex của inventory
		lineSystemApd, err := new(LineItem).GetLineSystemApacdex(inventory.Id)
		if err != nil {
			resp.Status = "error"
			resp.Errors = append(resp.Errors, ajax.Error{
				Id:      "",
				Message: "Error!",
			})
			return
		}
		if lineSystemApd.Id != 0 {
			// Trường hợp có r tiến hành update lại line
			new(LineItem).UpdateLineSystemApacdex(lineSystemApd, inputs.BidderId)
		} else {
			// Nếu chưa có thì tạo line auto mới
			new(LineItem).AutoCreateLineSystemApacdex(inventory, bidderSystem)
		}
	}

	// Kiểm tra lại số lượng bidder system đang waitting
	countConnectWaiting := 0
	for _, bidderId := range new(RlsBidderSystemInventory).GetListIdBidderApprove(inventory.Name) {
		bidder := new(Bidder).GetByIdNoCheckUser(bidderId)
		if bidder.Id == 0 || bidder.BidderTemplateId == 1 { // Bỏ qua nếu k tìm thấy bidder hoặc là bidder google
			continue
		}
		status := new(InventoryConnectionDemand).GetStatus(inventory.Id, bidder.Id)
		if status == 2 {
			countConnectWaiting++
		}
	}
	resp.Status = "success"
	resp.DataObject = fiber.Map{
		"id":           "box_status_" + strconv.FormatInt(inputs.BidderId, 10),
		"status":       htmlblock.Render("supply/connection/block.status.gohtml", fiber.Map{"Status": inputs.Status, "BidderId": inputs.BidderId}).String(),
		"countWaiting": countConnectWaiting,
	}
	new(Inventory).ResetCacheWorker(inventory.Id)

	// Push history
	connectionDemandNew := t.GetConnection(inputs.InventoryId, inputs.BidderId)
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = user.Id
	}
	historyInventory := history.Inventory{
		Detail:    history.DetailInventoryConnectionFE,
		CreatorId: creatorId,
	}

	inventory.ConnectionDemand = connectionDemandOld.TableInventoryConnectionDemand
	historyInventory.RecordOld = inventory.TableInventory
	inventory.ConnectionDemand = connectionDemandNew.TableInventoryConnectionDemand
	historyInventory.RecordNew = inventory.TableInventory

	_ = history.PushHistory(&historyInventory)

	return
}

type InventoryConnectionDemands struct {
	BidderRecord
	Status string `json:"status"`
	Logo   string `json:"logo"`
}

func (t *InventoryConnectionDemand) GetByInventory(inventory InventoryRecord) (records []InventoryConnectionDemands) {
	// Lấy theo list bidder domain được approve từ BE
	var bidders []BidderRecord
	err := mysql.Client.Select("bidder.*, rls_bidder_system_inventory.status").
		Where("rls_bidder_system_inventory.inventory_name = ? and bidder.bidder_template_id != 1 and (rls_bidder_system_inventory.status = 3 or rls_bidder_system_inventory.status = 7 or rls_bidder_system_inventory.status = 8)", inventory.Name).
		Joins("inner join rls_bidder_system_inventory on rls_bidder_system_inventory.bidder_id = bidder.id").
		Model(&bidders).Find(&bidders).Error
	if err != nil || len(bidders) == 0 {
		return
	}

	for _, bidder := range bidders {
		bidderTemplate := new(BidderTemplate).GetById(bidder.BidderTemplateId)
		status := new(InventoryConnectionDemand).GetStatus(inventory.Id, bidder.Id)
		rec := InventoryConnectionDemands{
			BidderRecord: bidder,
			Logo:         bidderTemplate.Logo,
		}
		if status == 1 {
			rec.Status = "Live"
		} else {
			rec.Status = "Waiting"
		}
		records = append(records, rec)
	}
	return
}
