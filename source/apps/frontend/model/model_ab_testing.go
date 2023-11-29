package model

import (
	"database/sql"
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
	"time"
)

type AbTesting struct{}

type AbTestingRecord struct {
	mysql.TableAbTesting
}

func (AbTestingRecord) TableName() string {
	return mysql.Tables.AbTesting
}

type AbTestingRecordDatatable struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	TestType      string `json:"test_type"`
	Bidder        string `json:"bidder"`
	TestGroupSize string `json:"test_group_size"`
	Action        string `json:"action"`
	Status        string `json:"status"`
}

type ListOfAbTesting struct {
	Id   int64
	List []string
}

func (t *AbTesting) GetAll() (records []AbTestingRecord) {
	mysql.Client.Find(&records)
	return
}

func (t *AbTesting) Create(inputs payload.AbTestingCreate, userLogin UserRecord, userAdmin UserRecord, lang lang.Translation) (record AbTestingRecord, errs []ajax.Error) {
	// Validate inputs
	errs = t.ValidateCreate(inputs)
	if len(errs) > 0 {
		return
	}
	// Bỏ qua input target không dùng nếu testType là auction
	if inputs.TestType != mysql.TYPETestTypeAuctionTimeOut {
		inputs.ListAdSize = []payload.ListTarget{}
		inputs.ListAdFormat = []payload.ListTarget{}
		inputs.ListAdTag = []payload.ListTarget{}
	} else {
		inputs.ListAdSize = []payload.ListTarget{}
		inputs.ListAdFormat = []payload.ListTarget{}
		inputs.ListAdTag = []payload.ListTarget{}
		inputs.ListGeo = []payload.ListTarget{}
		inputs.ListDevice = []payload.ListTarget{}
	}

	// Insert to database
	record.makeRow(inputs)
	record.UserId = userLogin.Id
	err := mysql.Client.Create(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.AbTestingError.Add.ToString(),
			})
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
		return
	}

	err = t.CreateTarget(record.Id, userLogin.Id, inputs)
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.AbTestingError.Target.ToString(),
			})
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
		return
	}

	//new(Inventory).UpdateRenderCacheWithAbTesting(record.Id, userLogin.Id)
	if inputs.TestType == mysql.TYPETestTypeDynamicHardPriceFloor {
		new(Inventory).UpdateRenderCacheWithFloor(inputs.DynamicFloorPrice, userLogin.Id)
	}
	// Get recordNew
	recordNew, _ := new(AbTesting).GetById(record.Id, userLogin.Id)
	// Push history
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userLogin.Id
	}
	_ = history.PushHistory(&history.ABTest{
		Detail:    history.DetailABTestFE,
		CreatorId: creatorId,
		RecordOld: mysql.TableAbTesting{},
		RecordNew: recordNew.TableAbTesting,
	})
	return
}

func (t *AbTesting) Edit(inputs payload.AbTestingCreate, userLogin UserRecord, userAdmin UserRecord, lang lang.Translation) (record AbTestingRecord, errs []ajax.Error) {
	recordOld, _ := t.VerificationRecord(inputs.Id, userLogin.Id)
	if recordOld.Id < 1 {
		errs = append(errs, ajax.Error{
			Id:      "id",
			Message: "You don't own this A/B Testing",
		})
	}
	// Bỏ qua input target không dùng nếu testType là auction
	if recordOld.TestType == mysql.TYPETestTypeAuctionTimeOut {
		inputs.ListAdSize = []payload.ListTarget{}
		inputs.ListAdFormat = []payload.ListTarget{}
		inputs.ListAdTag = []payload.ListTarget{}
	}
	// Gán lại inputs các giá trị không thay đổi
	inputs.TestType = recordOld.TestType

	// Validate inputs
	errs = t.ValidateCreate(inputs)
	if len(errs) > 0 {
		return
	}
	// Đặt record mới bằng recordOld để save all k mất giá dữ liệu các field k thay đổi
	record = recordOld

	// Insert to database
	record.makeRowEdit(inputs)
	record.UserId = userLogin.Id
	err := mysql.Client.Save(&record).Where("id = ?", record.Id).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.AbTestingError.Edit.ToString(),
			})
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}

	// Đoạn này để update lại cache vs những trường hợp domain khi chưa submit (ví dụ nếu pub xóa 1 domain trong target thì mình cũng phải clear AbTesting này ra khỏi cache)
	//new(Inventory).UpdateRenderCacheWithAbTesting(record.Id, userLogin.Id)

	err = new(Target).DeleteTarget(TargetRecord{mysql.TableTarget{
		AbTestingId: record.Id,
	}})
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.AbTestingError.Edit.ToString(),
			})
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}
	err = t.CreateTarget(record.Id, userLogin.Id, inputs)
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.AbTestingError.Target.ToString(),
			})
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}

	if inputs.TestType == mysql.TYPETestTypeDynamicHardPriceFloor {
		new(Inventory).UpdateRenderCacheWithFloor(inputs.DynamicFloorPrice, userLogin.Id)
	}

	// Get recordNew
	recordNew, _ := new(AbTesting).GetById(record.Id, userLogin.Id)
	// Push history
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userLogin.Id
	}
	_ = history.PushHistory(&history.ABTest{
		Detail:    history.DetailABTestFE,
		CreatorId: creatorId,
		RecordOld: recordOld.TableAbTesting,
		RecordNew: recordNew.TableAbTesting,
	})
	return
}

func (this *AbTesting) Delete(id, userId int64, userAdmin UserRecord, lang lang.Translation) fiber.Map {
	record, _ := new(AbTesting).GetById(id, userId)
	err := mysql.Client.Model(&AbTestingRecord{}).Delete(&AbTestingRecord{}, "id = ? and user_id = ?", id, userId).Error
	if err != nil {
		if !utility.IsWindow() {
			return fiber.Map{
				"status":  "err",
				"message": lang.Errors.AbTestingError.Delete.ToString(),
				"id":      id,
			}
		}
		return fiber.Map{
			"status":  "err",
			"message": err.Error(),
			"id":      id,
		}
	}
	//new(Inventory).UpdateRenderCacheWithAbTesting(id, userId)
	// History
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userId
	}
	_ = history.PushHistory(&history.ABTest{
		Detail:    history.DetailABTestFE,
		CreatorId: creatorId,
		RecordOld: record.TableAbTesting,
		RecordNew: mysql.TableAbTesting{},
	})
	new(Inventory).ResetCacheAll(userId)
	return fiber.Map{
		"status":  "success",
		"message": "done",
		"id":      id,
	}
}

func (rec *AbTestingRecord) makeRow(inputs payload.AbTestingCreate) {
	rec.Name = inputs.Name
	rec.Description = inputs.Description
	rec.TestType = inputs.TestType
	rec.Status = inputs.Status
	layoutISO := "01/02/2006"
	if inputs.StartDate != "" {
		startDate, err := time.Parse(layoutISO, inputs.StartDate)
		if err != nil {
			return
		}
		err = rec.StartDate.Scan(startDate)
		if err != nil {
			return
		}
	} else {
		rec.StartDate = sql.NullTime{}
	}
	if inputs.EndDate != "" {
		endDate, err := time.Parse(layoutISO, inputs.EndDate)
		if err != nil {
			return
		}
		err = rec.EndDate.Scan(endDate)
		if err != nil {
			return
		}
	} else {
		rec.EndDate = sql.NullTime{}
	}
	switch rec.TestType {
	case mysql.TYPETestTypeAuctionTimeOut:
		rec.TestGroupSize = inputs.TestGroupSize
		break
	case mysql.TYPETestTypeClientVsServer:
		rec.BidderId = inputs.Bidder
		rec.TestGroupSize = inputs.TestGroupSize
		break
	case mysql.TYPETestTypeUserIdModule:
		rec.UserIdModuleId = inputs.UserIdModule
		rec.TestGroupSize = inputs.TestGroupSize
		break
	case mysql.TYPETestTypeDynamicHardPriceFloor:
		rec.DynamicFloorPrice = inputs.DynamicFloorPrice
		rec.HardPriceFloor = inputs.HardPriceFloor
		break
	}
}

func (rec *AbTestingRecord) makeRowEdit(inputs payload.AbTestingCreate) {
	rec.Id = inputs.Id
	rec.Name = inputs.Name
	rec.Description = inputs.Description
	//rec.TestType = inputs.TestType
	rec.BidderId = inputs.Bidder
	rec.TestGroupSize = inputs.TestGroupSize
	rec.Status = inputs.Status
	layoutISO := "01/02/2006"
	if inputs.StartDate != "" {
		startDate, err := time.Parse(layoutISO, inputs.StartDate)
		if err != nil {
			return
		}
		err = rec.StartDate.Scan(startDate)
		if err != nil {
			return
		}
	} else {
		rec.StartDate = sql.NullTime{}
	}
	if inputs.EndDate != "" {
		endDate, err := time.Parse(layoutISO, inputs.EndDate)
		if err != nil {
			return
		}
		err = rec.EndDate.Scan(endDate)
		if err != nil {
			return
		}
	} else {
		rec.EndDate = sql.NullTime{}
	}
	switch rec.TestType {
	case mysql.TYPETestTypeAuctionTimeOut:
		rec.TestGroupSize = inputs.TestGroupSize
		break
	case mysql.TYPETestTypeClientVsServer:
		rec.BidderId = inputs.Bidder
		rec.TestGroupSize = inputs.TestGroupSize
		break
	case mysql.TYPETestTypeUserIdModule:
		rec.UserIdModuleId = inputs.UserIdModule
		rec.TestGroupSize = inputs.TestGroupSize
		break
	case mysql.TYPETestTypeDynamicHardPriceFloor:
		rec.DynamicFloorPrice = inputs.DynamicFloorPrice
		rec.HardPriceFloor = inputs.HardPriceFloor
		break
	}
}

func (t *AbTesting) ValidateCreate(inputs payload.AbTestingCreate) (errs []ajax.Error) {
	lang := lang.Translate
	if utility.ValidateString(inputs.Name) == "" {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: lang.ErrorRequired.ToString(),
		})
	}
	if !inputs.TestType.CheckValid() {
		errs = append(errs, ajax.Error{
			Id:      "test_type",
			Message: lang.ErrorRequired.ToString(),
		})
	} else if inputs.TestType == mysql.TYPETestTypeUserIdModule {
		if inputs.UserIdModule == 0 {
			errs = append(errs, ajax.Error{
				Id:      "select_user_id_module",
				Message: lang.ErrorRequired.ToString(),
			})
		}
	} else if inputs.TestType == mysql.TYPETestTypeClientVsServer {
		if inputs.Bidder == 0 {
			errs = append(errs, ajax.Error{
				Id:      "select-bidder",
				Message: lang.ErrorRequired.ToString(),
			})
		}
	} else if inputs.TestType == mysql.TYPETestTypeDynamicHardPriceFloor {
		if inputs.DynamicFloorPrice == 0 {
			errs = append(errs, ajax.Error{
				Id:      "dynamic_floor_price",
				Message: lang.ErrorRequired.ToString(),
			})
		}
		if inputs.HardPriceFloor == 0 {
			errs = append(errs, ajax.Error{
				Id:      "hard_price_floor",
				Message: lang.ErrorRequired.ToString(),
			})
		}
	}
	if !inputs.TestGroupSize.CheckValid() && inputs.TestType != mysql.TYPETestTypeDynamicHardPriceFloor {
		errs = append(errs, ajax.Error{
			Id:      "test_group_size",
			Message: lang.ErrorRequired.ToString(),
		})
	}

	return
}

func (t *AbTesting) ValidateEdit(inputs payload.AbTestingCreate, userId int64) (errs []ajax.Error) {
	if utility.ValidateString(inputs.Name) == "" {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Name is required",
		})
	}
	if inputs.TestType == mysql.TYPETestTypeDynamicHardPriceFloor {
		if inputs.DynamicFloorPrice == 0 {
			errs = append(errs, ajax.Error{
				Id:      "dynamic_floor_price",
				Message: lang.Translate.ErrorRequired.ToString(),
			})
		}
		if inputs.HardPriceFloor == 0 {
			errs = append(errs, ajax.Error{
				Id:      "hard_price_floor",
				Message: lang.Translate.ErrorRequired.ToString(),
			})
		}
	}
	// if utility.ValidateString(inputs.Description) == "" {
	//	errs = append(errs, ajax.Error{
	//		Id:      "description",
	//		Message: "Description is required",
	//	})
	// }
	return
}

func (t *AbTesting) GetByFilters(inputs *payload.AbTestingFilterPayload, userId int64, lang lang.Translation) (response datatable.Response, err error) {
	var inventories []AbTestingRecord
	var total int64

	err = mysql.Client.Where("user_id = ?", userId).
		Scopes(
			t.setFilterTarget(inputs, userId),
			t.SetFilterStatus(inputs),
			t.setFilterSearch(inputs),
		).
		Model(&inventories).Count(&total).
		Scopes(
			t.setOrder(inputs),
			pagination.Paginate(pagination.Params{
				Limit:  inputs.Length,
				Offset: inputs.Start,
			}),
		).
		Find(&inventories).Error
	if err != nil {
		if !utility.IsWindow() {
			err = fmt.Errorf(lang.Errors.AbTestingError.List.ToString())
		}
		return datatable.Response{}, err
	}
	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.MakeResponseDatatable(inventories)
	return
}

func (t *AbTesting) MakeResponseDatatable(AbTestings []AbTestingRecord) (records []AbTestingRecordDatatable) {
	for _, abTesting := range AbTestings {
		records = append(records, abTesting.Modify())
	}
	return
}

func (rec AbTestingRecord) Modify() (record AbTestingRecordDatatable) {
	record.Id = rec.Id
	//record.Description = rec.Description

	record.Action = htmlblock.Render("ab_testing/block_html/block.action.gohtml", rec).String()
	record.Name = htmlblock.Render("ab_testing/block_html/name_block.gohtml", rec).String()
	record.TestType = rec.TestType.String()
	bidder := new(Bidder).GetById(rec.BidderId, rec.UserId)
	record.Bidder = bidder.BidderCode
	if rec.TestType == mysql.TYPETestTypeDynamicHardPriceFloor {
		record.TestGroupSize = "100%"
	} else {
		record.TestGroupSize = strconv.Itoa(rec.TestGroupSize.Value()) + "%"
	}
	if rec.Status == 1 {
		record.Status = `<span class="">ON</span>`
	} else {
		record.Status = `<span class="">OFF</span>`
	}
	return
}

func (t *AbTesting) SetFilterStatus(inputs *payload.AbTestingFilterPayload) func(db *gorm.DB) *gorm.DB {
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

func (t *AbTesting) setFilterSearch(inputs *payload.AbTestingFilterPayload) func(db *gorm.DB) *gorm.DB {
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
		return db.Where("name LIKE ?", "%"+inputs.PostData.QuerySearch+"%")
	}
}

func (t *AbTesting) setFilterTarget(inputs *payload.AbTestingFilterPayload, userId int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var flag bool

		abTestingIDs := []int64{}
		targets, _ := new(Target).GetTargetByFilterAbTesting(inputs, userId)
		if len(targets) > 0 {
			flag = true
		}
		if !flag {
			return db.Where("id IN ?", abTestingIDs)
		}
		for _, target := range targets {
			abTestingIDs = append(abTestingIDs, target.AbTestingId)
		}
		return db.Where("id IN ?", abTestingIDs)
	}
}

func (t *AbTesting) setOrder(inputs *payload.AbTestingFilterPayload) func(db *gorm.DB) *gorm.DB {
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
		return db
	}
}

func (t *AbTesting) GetById(id, userId int64) (row AbTestingRecord, err error) {
	err = mysql.Client.Where("id = ? and user_id = ?", id, userId).Find(&row).Error
	// if row.Id == 0 || err != nil {
	//	return
	// }
	return
}

func (t *AbTesting) CreateTarget(abTestingId, userId int64, inputs payload.AbTestingCreate) (err error) {
	all := int64(-1)
	// Kiểm tra nếu đầu vào input list target = 0 thì thêm một target = 0 thể hiện select all
	if len(inputs.ListAdInventory) == 0 {
		inputs.ListAdInventory = []payload.ListTarget{
			{
				Id: all,
			},
		}
	}
	if len(inputs.ListAdFormat) == 0 {
		inputs.ListAdFormat = []payload.ListTarget{
			{
				Id: all,
			},
		}
	}
	if len(inputs.ListAdSize) == 0 {
		inputs.ListAdSize = []payload.ListTarget{
			{
				Id: all,
			},
		}
	}
	if len(inputs.ListAdTag) == 0 {
		inputs.ListAdTag = []payload.ListTarget{
			{
				Id: all,
			},
		}
	}
	if len(inputs.ListGeo) == 0 {
		inputs.ListGeo = []payload.ListTarget{
			{
				Id: all,
			},
		}
	}
	if len(inputs.ListDevice) == 0 {
		inputs.ListDevice = []payload.ListTarget{
			{
				Id: all,
			},
		}
	}
	for _, inventory := range inputs.ListAdInventory {
		var recordTarget TargetRecord
		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
			UserId:      userId,
			AbTestingId: abTestingId,
			InventoryId: inventory.Id,
		}}).Error
		if err != nil {
			return err
		}
	}
	for _, adFormat := range inputs.ListAdFormat {
		var recordTarget TargetRecord
		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
			UserId:      userId,
			AbTestingId: abTestingId,
			AdFormatId:  adFormat.Id,
		}}).Error
		if err != nil {
			return err
		}
	}
	for _, size := range inputs.ListAdSize {
		var recordTarget TargetRecord
		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
			UserId:      userId,
			AbTestingId: abTestingId,
			AdSizeId:    size.Id,
		}}).Error
		if err != nil {
			return err
		}
	}
	for _, geo := range inputs.ListGeo {
		var recordTarget TargetRecord
		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
			UserId:      userId,
			AbTestingId: abTestingId,
			GeoId:       geo.Id,
		}}).Error
		if err != nil {
			return err
		}
	}
	for _, device := range inputs.ListDevice {
		var recordTarget TargetRecord
		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
			UserId:      userId,
			AbTestingId: abTestingId,
			DeviceId:    device.Id,
		}}).Error
		if err != nil {
			return err
		}
	}

	for _, adTag := range inputs.ListAdTag {
		var recordTarget TargetRecord
		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
			UserId:      userId,
			AbTestingId: abTestingId,
			TagId:       adTag.Id,
		}}).Error
		if err != nil {
			return err
		}
	}
	//new(Inventory).UpdateRenderCacheWithAbTesting(abTestingId, userId)
	return
}

func (t *AbTesting) VerificationRecord(id, userId int64) (record AbTestingRecord, err error) {
	err = mysql.Client.Model(&AbTestingRecord{}).Where("id = ? and user_id = ?", id, userId).Find(&record).Error
	return
}

func (t *AbTesting) GetListBoxCollapse(userId, recordId int64, page, typ string) (list []string) {
	switch typ {
	case "add":
		mysql.Client.Select("box_collapse").Model(&PageCollapseRecord{}).Where("user_id = ? and page_collapse = ? and is_collapse = ? and page_type = ?", userId, page, 1, typ).Find(&list)
		return
	case "edit":
		mysql.Client.Select("box_collapse").Model(PageCollapseRecord{}).Where("user_id = ? and page_collapse = ? and is_collapse = ? and page_type = ? and page_id = ?", userId, page, 1, typ, recordId).Find(&list)
		return
	}
	return
}
