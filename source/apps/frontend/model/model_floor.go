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

type Floor struct{}

type FloorRecord struct {
	mysql.TableFloor
}

func (FloorRecord) TableName() string {
	return mysql.Tables.Floor
}

type FloorRecordDatatable struct {
	Id          int64  `json:"id"`
	Description string `json:"description"`
	FloorValue  string `json:"floor_value"`
	Priority    int    `json:"priority"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	User        string `json:"user"`
	Inventory   string `json:"inventory"`
	AdTag       string `json:"ad_tag"`
	AdFormat    string `json:"ad_format"`
	AdSize      string `json:"ad_size"`
	Country     string `json:"country"`
	Device      string `json:"device"`
	Action      string `json:"action"`
}

type ListOfFloor struct {
	Id   int64
	List []string
}

func (t *Floor) GetAll() (records []FloorRecord) {
	mysql.Client.Find(&records)
	return
}

func (t *Floor) GetFloorDynamic(userId int64) (records []FloorRecord) {
	mysql.Client.Where("user_id = ? and floor_type = 1", userId).Find(&records)
	return
}

func (t *Floor) Create(inputs payload.FloorCreate, user UserRecord, userAdmin UserRecord) (record FloorRecord, errs []ajax.Error) {
	lang := lang.Translate
	// Xử lý lại inputs
	if inputs.FloorType == 1 { // Dynamic
		// Loại các target không dùng trong floor type dynamic
		inputs.ListAdFormat = []payload.ListTarget{}
		inputs.ListAdSize = []payload.ListTarget{}
		inputs.ListAdTag = []payload.ListTarget{}
		inputs.ListGeo = []payload.ListTarget{}
		inputs.ListDevice = []payload.ListTarget{}
	}

	// Validate inputs
	errs = t.ValidateCreate(inputs)
	if len(errs) > 0 {
		return
	}
	// Check exist
	// Insert to database
	record.makeRow(inputs)
	record.UserId = user.Id

	err := mysql.Client.Create(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.FloorError.Add.ToString(),
			})
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
		return
	}

	err = t.CreateTarget(record.Id, user.Id, inputs)
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.FloorError.Target.ToString(),
			})
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
		return
	}

	new(Inventory).UpdateRenderCacheWithFloor(record.Id, user.Id)
	// Push history
	recordNew, _ := new(Floor).GetById(record.Id, user.Id)
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = user.Id
	}
	_ = history.PushHistory(&history.Floor{
		Detail:    history.DetailFloorFE,
		CreatorId: creatorId,
		RecordOld: mysql.TableFloor{},
		RecordNew: recordNew.TableFloor,
	})
	return
}

func (t *Floor) Edit(inputs payload.FloorCreate, user UserRecord, userAdmin UserRecord) (record FloorRecord, errs []ajax.Error) {
	lang := lang.Translate
	// Xử lý lại inputs
	if inputs.FloorType == 1 { // Dynamic
		// Loại các target không dùng trong floor type dynamic
		inputs.ListAdFormat = []payload.ListTarget{}
		inputs.ListAdSize = []payload.ListTarget{}
		inputs.ListAdTag = []payload.ListTarget{}
		inputs.ListGeo = []payload.ListTarget{}
		inputs.ListDevice = []payload.ListTarget{}
	}

	// Validate inputs
	errs = t.ValidateEdit(inputs, user.Id)
	if len(errs) > 0 {
		return
	}
	// Check exist
	recordOld, err := t.GetById(inputs.Id, user.Id)
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}
	if recordOld.Id == 0 {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: "Floor does not exist!",
		})
	}
	record = recordOld
	// Insert to database
	record.makeRowEdit(inputs)
	// err = mysql.Client.Updates(&record).Update("status", record.Status).Where("id = ?", record.Id).Error
	err = mysql.Client.Where("id = ?", record.Id).Save(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.FloorError.Edit.ToString(),
			})
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}

	// Đoạn này để update lại cache vs những trường hợp domain khi chưa submit (ví dụ nếu pub xóa 1 domain trong target thì mình cũng phải clear floor này ra khỏi cache)
	new(Inventory).UpdateRenderCacheWithFloor(record.Id, user.Id)

	err = new(Target).DeleteTarget(TargetRecord{mysql.TableTarget{
		FloorId: record.Id,
	}})
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.FloorError.Edit.ToString(),
			})
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}
	err = t.CreateTarget(record.Id, user.Id, inputs)
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.FloorError.Target.ToString(),
			})
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}

	new(Inventory).UpdateRenderCacheWithFloor(record.Id, user.Id)
	// Push history
	recordNew, _ := new(Floor).GetById(record.Id, user.Id)
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = user.Id
	}
	_ = history.PushHistory(&history.Floor{
		Detail:    history.DetailFloorFE,
		CreatorId: creatorId,
		RecordOld: recordOld.TableFloor,
		RecordNew: recordNew.TableFloor,
	})
	return
}

func (this *Floor) Delete(id, userId int64, userAdmin UserRecord) fiber.Map {
	lang := lang.Translate
	record, _ := new(Floor).GetById(id, userId)
	err := mysql.Client.Model(&FloorRecord{}).Delete(&FloorRecord{}, "id = ? and user_id = ?", id, userId).Error
	if err != nil {
		if !utility.IsWindow() {
			return fiber.Map{
				"status":  "err",
				"message": lang.Errors.FloorError.Delete.ToString(),
				"id":      id,
			}
		}
		return fiber.Map{
			"status":  "err",
			"message": err.Error(),
			"id":      id,
		}
	}
	new(Inventory).UpdateRenderCacheWithFloor(id, userId)

	// History
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userId
	}
	_ = history.PushHistory(&history.Floor{
		Detail:    history.DetailFloorFE,
		CreatorId: creatorId,
		RecordOld: record.TableFloor,
		RecordNew: mysql.TableFloor{},
	})

	return fiber.Map{
		"status":  "success",
		"message": "done",
		"id":      id,
	}
}

func (rec *FloorRecord) makeRow(inputs payload.FloorCreate) {
	rec.Id = inputs.Id
	rec.Name = inputs.Name
	rec.Description = inputs.Description
	if rec.Description == "" {
		rec.Description = " "
	}
	rec.FloorValue = inputs.FloorValue
	rec.FloorType = inputs.FloorType
	rec.Status = inputs.Status
	rec.Priority = inputs.Priority
}

func (rec *FloorRecord) makeRowEdit(inputs payload.FloorCreate) {
	rec.Name = inputs.Name
	rec.Description = inputs.Description
	if rec.Description == "" {
		rec.Description = " "
	}
	rec.FloorValue = inputs.FloorValue
	rec.Status = inputs.Status
	rec.Priority = inputs.Priority
}

func (t *Floor) ValidateCreate(inputs payload.FloorCreate) (errs []ajax.Error) {
	if utility.ValidateString(inputs.Name) == "" {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Name is required",
		})
	}

	// if utility.ValidateString(inputs.Description) == "" {
	//	errs = append(errs, ajax.Error{
	//		Id:      "description",
	//		Message: "Description is required",
	//	})
	// }
	if inputs.FloorType == 2 {
		if inputs.FloorValue <= 0 {
			errs = append(errs, ajax.Error{
				Id:      "floor_value",
				Message: "Floor Value is invalid",
			})
		}
	}

	return
}

func (t *Floor) ValidateEdit(inputs payload.FloorCreate, userId int64) (errs []ajax.Error) {
	flag := t.VerificationRecord(inputs.Id, userId)
	if !flag {
		errs = append(errs, ajax.Error{
			Id:      "id",
			Message: "You don't own this floor",
		})
	}
	if utility.ValidateString(inputs.Name) == "" {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Name is required",
		})
	}

	// if utility.ValidateString(inputs.Description) == "" {
	//	errs = append(errs, ajax.Error{
	//		Id:      "description",
	//		Message: "Description is required",
	//	})
	// }
	if inputs.FloorType == 2 {
		if inputs.FloorValue <= 0 {
			errs = append(errs, ajax.Error{
				Id:      "floor_value",
				Message: "Floor Value is invalid",
			})
		}
	}

	return
}

func (t *Floor) GetByFilters(inputs *payload.FloorFilterPayload, userId int64, lang lang.Translation) (response datatable.Response, err error) {
	var inventories []FloorRecord
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
			err = fmt.Errorf(lang.Errors.FloorError.List.ToString())
		}
		return datatable.Response{}, err
	}
	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.MakeResponseDatatable(inventories)
	return
}

func (t *Floor) MakeResponseDatatable(Floors []FloorRecord) (records []FloorRecordDatatable) {
	for _, floor := range Floors {
		records = append(records, floor.Modify())
	}
	return
}

func (rec FloorRecord) Modify() (record FloorRecordDatatable) {
	record.Id = rec.Id
	record.Description = rec.Description
	record.Priority = rec.Priority
	user := new(User).GetById(rec.UserId)
	record.User = user.Email

	var listNameInventoryTarget []string
	var listNameAdFormatTarget []string
	var listNameAdSizeTarget []string
	var listNameAdTagTarget []string
	var listNameGeoTarget []string
	var listNameDeviceTarget []string

	targets := new(Target).GetTargetFloor(rec.Id)
	for _, target := range targets {
		if target.InventoryId != 0 {
			inventory, _ := new(Inventory).GetById(target.InventoryId, rec.UserId)
			listNameInventoryTarget = append(listNameInventoryTarget, inventory.Name)
		} else if target.AdFormatId != 0 {
			adFormat := new(AdType).GetById(target.AdFormatId)
			listNameAdFormatTarget = append(listNameAdFormatTarget, adFormat.Name)
		} else if target.AdSizeId != 0 {
			adSize := new(AdSize).GetById(target.AdSizeId)
			listNameAdSizeTarget = append(listNameAdSizeTarget, adSize.Name)
		} else if target.TagId != 0 {
			adTag := new(InventoryAdTag).GetById(target.TagId)
			listNameAdTagTarget = append(listNameAdTagTarget, adTag.Name)
		} else if target.GeoId != 0 {
			geo := new(Country).GetById(target.GeoId)
			listNameGeoTarget = append(listNameGeoTarget, geo.Name)
		} else if target.DeviceId != 0 {
			device := new(Device).GetById(target.DeviceId)
			listNameDeviceTarget = append(listNameDeviceTarget, device.Name)
		}
	}

	inventory := ListOfFloor{
		Id:   rec.Id,
		List: listNameInventoryTarget,
	}
	adTag := ListOfFloor{
		Id:   rec.Id,
		List: listNameAdTagTarget,
	}
	adFormat := ListOfFloor{
		Id:   rec.Id,
		List: listNameAdFormatTarget,
	}
	adSize := ListOfFloor{
		Id:   rec.Id,
		List: listNameAdSizeTarget,
	}
	country := ListOfFloor{
		Id:   rec.Id,
		List: listNameGeoTarget,
	}
	device := ListOfFloor{
		Id:   rec.Id,
		List: listNameDeviceTarget,
	}

	// p := message.NewPrinter(language.English)
	// a, _ := p.Printf("%v: %v%v\n", cur, currency.Symbol(cur), 123123123)
	// config := new(Config).GetByUserId(user.Id)
	// cur := new(Currency).GetByCode(config.Currency)
	record.Inventory = htmlblock.Render("floor/block_html/list_inventory.gohtml", inventory).String()
	record.AdTag = htmlblock.Render("floor/block_html/list_adtag.gohtml", adTag).String()
	record.AdFormat = htmlblock.Render("floor/block_html/list_adformat.gohtml", adFormat).String()
	record.AdSize = htmlblock.Render("floor/block_html/list_adsize.gohtml", adSize).String()
	record.Country = htmlblock.Render("floor/block_html/list_country.gohtml", country).String()
	record.Device = htmlblock.Render("floor/block_html/list_device.gohtml", device).String()
	record.Action = htmlblock.Render("floor/block_html/block.action.gohtml", rec).String()
	// record.FloorValue = fmt.Sprintf("%v%.2f", cur.Symbol, rec.FloorValue)
	if rec.FloorValue == 0 {
		record.FloorValue = "--"
	} else {
		record.FloorValue = fmt.Sprintf("%v%.2f", "$", rec.FloorValue)
	}
	record.Name = htmlblock.Render("floor/block_html/name_block.gohtml", rec).String()
	if rec.Status == 1 {
		record.Status = `<span class="">ON</span>`
	} else {
		record.Status = `<span class="">OFF</span>`
	}
	return
}

func (t *Floor) SetFilterStatus(inputs *payload.FloorFilterPayload) func(db *gorm.DB) *gorm.DB {
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

func (t *Floor) setFilterSearch(inputs *payload.FloorFilterPayload) func(db *gorm.DB) *gorm.DB {
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

func (t *Floor) setFilterTarget(inputs *payload.FloorFilterPayload, userId int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var flag bool

		floorIDs := []int64{}
		targets, _ := new(Target).GetTargetByFilterFloor(inputs, userId)
		if len(targets) > 0 {
			flag = true
		}
		if !flag {
			return db.Where("id IN ?", floorIDs)
		}
		for _, target := range targets {
			floorIDs = append(floorIDs, target.FloorId)
		}
		return db.Where("id IN ?", floorIDs)
	}
}

func (t *Floor) setOrder(inputs *payload.FloorFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(inputs.Order) > 0 {
			var orders []string
			for _, order := range inputs.Order {
				column := inputs.Columns[order.Column]
				orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
			}
			orderString := strings.Join(orders, ", ")
			return db.Order(orderString)
		} else {
			return db.Order("id DESC")
		}
	}
}

func (t *Floor) GetById(id, userId int64) (row FloorRecord, err error) {
	err = mysql.Client.Where("id = ? and user_id = ?", id, userId).Find(&row).Error
	// Get Rls
	row.TableFloor.GetRls()
	return
}

func (t *Floor) GetByUser(userId int64) (records []FloorRecord) {
	mysql.Client.Where("status = 1 AND user_id = ?", userId).Find(&records)
	return
}

func (t *Floor) CreateTarget(floorId, userId int64, inputs payload.FloorCreate) (err error) {
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
			FloorId:     floorId,
			InventoryId: inventory.Id,
		}}).Error
		if err != nil {
			return err
		}
	}
	for _, adFormat := range inputs.ListAdFormat {
		var recordTarget TargetRecord
		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
			UserId:     userId,
			FloorId:    floorId,
			AdFormatId: adFormat.Id,
		}}).Error
		if err != nil {
			return err
		}
	}
	for _, size := range inputs.ListAdSize {
		var recordTarget TargetRecord
		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
			UserId:   userId,
			FloorId:  floorId,
			AdSizeId: size.Id,
		}}).Error
		if err != nil {
			return err
		}
	}
	for _, geo := range inputs.ListGeo {
		var recordTarget TargetRecord
		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
			UserId:  userId,
			FloorId: floorId,
			GeoId:   geo.Id,
		}}).Error
		if err != nil {
			return err
		}
	}
	for _, device := range inputs.ListDevice {
		var recordTarget TargetRecord
		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
			UserId:   userId,
			FloorId:  floorId,
			DeviceId: device.Id,
		}}).Error
		if err != nil {
			return err
		}
	}

	for _, adTag := range inputs.ListAdTag {
		var recordTarget TargetRecord
		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
			UserId:  userId,
			FloorId: floorId,
			TagId:   adTag.Id,
		}}).Error
		if err != nil {
			return err
		}
	}
	new(Inventory).UpdateRenderCacheWithFloor(floorId, userId)
	return
}

func (t *Floor) VerificationRecord(id, userId int64) bool {
	var row FloorRecord
	err := mysql.Client.Model(&FloorRecord{}).Where("id = ? and user_id = ?", id, userId).Find(&row).Error
	if err != nil || row.Id == 0 {
		return false
	}
	return true
}

func (t *Floor) GetListBoxCollapse(userId, recordId int64, page, typ string) (list []string) {
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

func (t *Floor) AutoCreateDefaultFloor(user UserRecord) (errs []ajax.Error) {
	userAutoCreateFloor := new(User).GetById(234)
	floorDisplay := payload.FloorCreate{
		Name:        "Default Floor Display $0.01",
		Description: "",
		Status:      1,
		FloorType:   2,
		FloorValue:  0.01,
		Priority:    10,
		ListAdFormat: []payload.ListTarget{
			payload.ListTarget{
				Id:   1,
				Name: "Display",
			},
			payload.ListTarget{
				Id:   5,
				Name: "Sticky Banner",
			},
		},
		ListAdSize:      []payload.ListTarget{},
		ListAdTag:       []payload.ListTarget{},
		ListAdInventory: []payload.ListTarget{},
		ListGeo:         []payload.ListTarget{},
		ListDevice:      []payload.ListTarget{},
	}
	_, errs = t.Create(floorDisplay, user, userAutoCreateFloor)

	floorVideo := payload.FloorCreate{
		Name:        "Default Floor Video $0.10",
		Description: "",
		Status:      1,
		FloorType:   2,
		FloorValue:  0.10,
		Priority:    10,
		ListAdFormat: []payload.ListTarget{
			payload.ListTarget{
				Id:   2,
				Name: "Instream",
			},
			payload.ListTarget{
				Id:   3,
				Name: "Outstream",
			},
			payload.ListTarget{
				Id:   4,
				Name: "Top Articles",
			},
		},
		ListAdSize:      []payload.ListTarget{},
		ListAdTag:       []payload.ListTarget{},
		ListAdInventory: []payload.ListTarget{},
		ListGeo:         []payload.ListTarget{},
		ListDevice:      []payload.ListTarget{},
	}
	_, errs = t.Create(floorVideo, user, userAutoCreateFloor)

	return
}
