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
	"strings"
)

type Blocking struct{}

type BlockingRecord struct {
	mysql.TableBlocking
}

func (BlockingRecord) TableName() string {
	return mysql.Tables.Blocking
}

type BlockingRecordDatatable struct {
	BlockingRecord
	RowId           string `json:"DT_RowId"`
	RestrictionName string `json:"restriction_name"`
	Action          string `json:"action"`
}

func (t *Blocking) GetByFilters(inputs *payload.BlockingFilterPayload, userId int64, lang lang.Translation) (response datatable.Response, err error) {
	var blockings []BlockingRecord
	var total int64
	err = mysql.Client.Where("user_id = ?", userId).
		Scopes(
			t.setFilterSearch(inputs),
		).
		Model(&blockings).Count(&total).
		Scopes(
			t.setOrder(inputs),
			pagination.Paginate(pagination.Params{
				Limit:  inputs.Length,
				Offset: inputs.Start,
			}),
		).
		Find(&blockings).Error
	if err != nil {
		if !utility.IsWindow() {
			err = fmt.Errorf(lang.Errors.BlockingError.List.ToString())
		}
		return datatable.Response{}, err
	}
	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.MakeResponseDatatable(blockings)
	return
}

func (t *Blocking) setFilterSearch(inputs *payload.BlockingFilterPayload) func(db *gorm.DB) *gorm.DB {
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
		return db.Where("restriction_name  LIKE ?", "%"+inputs.PostData.QuerySearch+"%")
	}
}

func (t *Blocking) setOrder(inputs *payload.BlockingFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(inputs.Order) > 0 {
			var orders []string
			for _, order := range inputs.Order {
				column := inputs.Columns[order.Column]
				orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
			}
			orderString := strings.Join(orders, ", ")
			return db.Order(orderString)
		}
		return db.Order("id DESC")
	}
}

func (t *Blocking) MakeResponseDatatable(blockings []BlockingRecord) (records []BlockingRecordDatatable) {
	for _, blocking := range blockings {
		rec := BlockingRecordDatatable{
			BlockingRecord:  blocking,
			RestrictionName: htmlblock.Render("blocking/index/block.name.gohtml", blocking).String(),
			Action:          htmlblock.Render("blocking/index/block.action.gohtml", blocking).String(),
		}
		records = append(records, rec)
	}
	return
}

func (t *Blocking) Create(inputs payload.BlockingAdd, user UserRecord, userAdmin UserRecord) (record BlockingRecord, errs []ajax.Error) {
	lang := lang.Translate
	// Validate inputs
	errs = t.ValidateCreate(inputs)
	if len(errs) > 0 {
		return
	}
	// Insert to database
	record = t.makeRecord(inputs, user)
	err := mysql.Client.Create(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.BlockingError.Add.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}
	// Create rl blocking vs inventory
	for _, inventoryId := range inputs.Inventories {
		_ = new(RlBlockingInventory).CreateRl(record.Id, inventoryId)
	}

	// Create blocking restrictions
	// Domain
	for _, advertiseDomain := range inputs.AdvertiseDomains {
		_ = new(BlockingRestrictions).CreateAdvertiseDomain(record.Id, advertiseDomain)
	}
	// CreativeId
	for _, creativeId := range inputs.CreativeIds {
		_ = new(BlockingRestrictions).CreateCreativeId(record.Id, creativeId)
	}

	// Reset cache
	for _, inventoryId := range inputs.Inventories {
		new(Inventory).ResetCacheWorker(inventoryId)
	}
	// Push History
	recordNew, _ := new(Blocking).GetById(record.Id, user.Id)
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = user.Id
	}
	_ = history.PushHistory(&history.Blocking{
		Detail:    history.DetailBlockingFE,
		CreatorId: creatorId,
		RecordOld: mysql.TableBlocking{},
		RecordNew: recordNew.TableBlocking,
	})
	return
}

func (t *Blocking) Update(inputs payload.BlockingAdd, user UserRecord, userAdmin UserRecord, lang lang.Translation) (record BlockingRecord, errs []ajax.Error) {
	recordOld, _ := t.GetById(inputs.Id, user.Id)
	if recordOld.Id == 0 {
		errs = append(errs, ajax.Error{
			Id:      "id",
			Message: "You don't own this blocking",
		})
		return
	}

	// Validate inputs
	errs = t.ValidateEdit(inputs, user.Id)
	if len(errs) > 0 {
		return
	}
	// Insert to database
	record = t.makeRecord(inputs, user)
	err := mysql.Client.Updates(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.BlockingError.Edit.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}
	var listInventoryResetCache []int64
	// Lấy ra listInventory cũ
	rlsBlockingInventories, _ := new(RlBlockingInventory).GetByBlockingId(record.Id)
	for _, rls := range rlsBlockingInventories {
		listInventoryResetCache = append(listInventoryResetCache, rls.InventoryId)
	}

	// Delete all rl blocking vs inventory cũ
	_ = new(RlBlockingInventory).DeleteByBlockingId(record.Id)

	// Create rls blocking vs inventory
	for _, inventoryId := range inputs.Inventories {
		_ = new(RlBlockingInventory).CreateRl(record.Id, inventoryId)
		// Check và thêm inventoryId để resetCache nếu chưa có trong list
		if !utility.InArray(inventoryId, listInventoryResetCache, false) {
			listInventoryResetCache = append(listInventoryResetCache, inventoryId)
		}
	}

	// Delete all blocking restrictions cũ
	_ = new(BlockingRestrictions).DeleteByBlockingId(record.Id)

	// Create blocking restrictions
	// Domain
	for _, advertiseDomain := range inputs.AdvertiseDomains {
		_ = new(BlockingRestrictions).CreateAdvertiseDomain(record.Id, advertiseDomain)
	}
	// CreativeId
	for _, creativeId := range inputs.CreativeIds {
		_ = new(BlockingRestrictions).CreateCreativeId(record.Id, creativeId)
	}
	// Reset cache
	for _, inventoryId := range listInventoryResetCache {
		new(Inventory).ResetCacheWorker(inventoryId)
	}
	// Push History
	recordNew, _ := new(Blocking).GetById(record.Id, user.Id)
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = user.Id
	}
	_ = history.PushHistory(&history.Blocking{
		Detail:    history.DetailBlockingFE,
		CreatorId: creatorId,
		RecordOld: recordOld.TableBlocking,
		RecordNew: recordNew.TableBlocking,
	})
	return
}

func (t *Blocking) ValidateCreate(inputs payload.BlockingAdd) (errs []ajax.Error) {
	if utility.ValidateString(inputs.RestrictionName) == "" {
		errs = append(errs, ajax.Error{
			Id:      "restriction_name",
			Message: "Restriction Name is required",
		})
	}
	if len(inputs.AdvertiseDomains) == 0 && len(inputs.CreativeIds) == 0 {
		errs = append(errs, ajax.Error{
			Id:      "textarea-add-custom",
			Message: "Advertise domain is required",
		})
	}
	var domainString []string
	for _, advertiseDomain := range inputs.AdvertiseDomains {
		flag := utility.ValidateDomainName(advertiseDomain)
		if !flag {
			domainString = append(domainString, advertiseDomain)
		}
	}
	if len(domainString) > 0 {
		if len(domainString) == 1 {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: "This domain " + strings.Join(domainString, ",") + " isn't valid,please check again!",
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: "These domains " + strings.Join(domainString, ",") + " aren't valid,please check again!",
			})
		}
	}
	return
}

func (t *Blocking) ValidateEdit(inputs payload.BlockingAdd, userId int64) (errs []ajax.Error) {
	flag := t.VerificationRecord(inputs.Id, userId)
	if !flag {
		errs = append(errs, ajax.Error{
			Id:      "id",
			Message: "You don't own this blocking",
		})
	}
	if utility.ValidateString(inputs.RestrictionName) == "" {
		errs = append(errs, ajax.Error{
			Id:      "restriction_name",
			Message: "Restriction Name is required",
		})
	}
	if len(inputs.AdvertiseDomains) == 0 && len(inputs.CreativeIds) == 0 {
		errs = append(errs, ajax.Error{
			Id:      "textarea-add-custom",
			Message: "Advertise domain is required",
		})
	}
	var domainString []string
	for _, advertiseDomain := range inputs.AdvertiseDomains {
		f := utility.ValidateDomainName(advertiseDomain)
		if !f {
			domainString = append(domainString, advertiseDomain)
		}
	}
	if len(domainString) > 0 {
		if len(domainString) == 1 {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: "This domain " + strings.Join(domainString, ",") + " isn't valid,please check again!",
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: "These domains " + strings.Join(domainString, ",") + " aren't valid,please check again!",
			})
		}
	}
	return
}

func (t *Blocking) makeRecord(inputs payload.BlockingAdd, user UserRecord) (record BlockingRecord) {
	record = BlockingRecord{mysql.TableBlocking{
		Id:              inputs.Id,
		UserId:          user.Id,
		RestrictionName: inputs.RestrictionName,
	}}
	return
}

func (t *Blocking) GetById(id int64, userId int64) (row BlockingRecord, err error) {
	err = mysql.Client.Where("id = ? and user_id = ?", id, userId).Find(&row).Error
	// Get rls
	row.TableBlocking.GetRls()
	return
}

func (t *Blocking) GetByUser(userId int64) (records []BlockingRecord) {
	mysql.Client.Where("user_id = ?", userId).Find(&records)
	return
}

func (t *Blocking) Delete(id, userId int64, userAdmin UserRecord, lang lang.Translation) fiber.Map {
	record, _ := new(Blocking).GetById(id, userId)
	err := mysql.Client.Model(&BlockingRecord{}).Delete(&BlockingRecord{}, "id = ? and user_id = ?", id, userId).Error
	if err != nil {
		if !utility.IsWindow() {
			return fiber.Map{
				"status":  "err",
				"message": lang.Errors.BlockingError.Delete.ToString(),
				"id":      id,
			}
		}
		return fiber.Map{
			"status":  "err",
			"message": err.Error(),
			"id":      id,
		}
	}
	new(BlockingRestrictions).DeleteByBlockingId(id)

	// History
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userId
	}
	_ = history.PushHistory(&history.Blocking{
		Detail:    history.DetailBlockingFE,
		CreatorId: creatorId,
		RecordOld: record.TableBlocking,
		RecordNew: mysql.TableBlocking{},
	})

	return fiber.Map{
		"status":  "success",
		"message": "done",
		"id":      id,
	}

}

func (t *Blocking) VerificationRecord(id, userId int64) bool {
	var row BlockingRecord
	err := mysql.Client.Model(&BlockingRecord{}).Where("id = ? and user_id = ?", id, userId).Find(&row).Error
	if err != nil || row.Id == 0 {
		return false
	}
	return true
}
