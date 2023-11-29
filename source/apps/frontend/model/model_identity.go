package model

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
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

type Identity struct{}

type IdentityRecord struct {
	mysql.TableIdentity
}

func (IdentityRecord) TableName() string {
	return mysql.Tables.Identity
}

type IdentityRecordDatatable struct {
	Id            int64  `json:"id"`
	Description   string `json:"description"`
	IdentityValue string `json:"identity_value"`
	Priority      string `json:"priority"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	User          string `json:"user"`
	Inventory     string `json:"inventory"`
	AdTag         string `json:"ad_tag"`
	AdFormat      string `json:"ad_format"`
	AdSize        string `json:"ad_size"`
	Country       string `json:"country"`
	Device        string `json:"device"`
	Action        string `json:"action"`
}

type ListOfIdentity struct {
	Id   int64
	List []string
}

func (t *Identity) GetAll() (records []IdentityRecord) {
	mysql.Client.Find(&records)
	return
}

func (t *Identity) IsHaveAProfile(user UserRecord, currentIdentityId int64) (have bool) {
	var record IdentityRecord
	mysql.Client.
		Select("id").
		Where("id != ? AND user_id = ? AND status = ?", currentIdentityId, user.Id, mysql.TypeOn).
		Last(&record)
	if record.Id > 0 {
		have = true
	}
	return
}

func (t *Identity) GetIdsProfileTargetAllInventory(user UserRecord) (Ids []int64) {
	var identityTargetAllInventory []IdentityRecord
	//inventoryIds := new(Inventory).GetAllIdsOfUser(user.Id)
	tableTarget := mysql.Tables.Target
	tableIdentity := mysql.Tables.Identity
	mysql.Client.
		Select(tableIdentity+".id").
		Joins("LEFT JOIN "+tableTarget+" ON "+tableTarget+".identity_id = "+tableIdentity+".id").
		Where(tableTarget+".inventory_id = ?", -1).
		Where(tableIdentity+".status = ?", mysql.TypeOn).
		Find(&identityTargetAllInventory)
	for _, identity := range identityTargetAllInventory {
		Ids = append(Ids, identity.Id)
	}
	return
}

func (t *Identity) Create(inputs payload.IdentityCreate, user UserRecord, userAdmin UserRecord) (record IdentityRecord, errs []ajax.Error) {
	lang := lang.Translate
	var err error

	// Validate inputs
	errs = t.ValidateCreate(inputs, user, lang)
	if len(errs) > 0 {
		return
	}

	//if len(inputs.ListAdInventory) == 0 {
	//	// trường hợp target all domain thì sẽ disable toàn bộ các profile
	//	err = t.DisableAllProfile(user)
	//} else {
	//	// trường hợp target 1 domain thì sẽ disable toàn bộ các profile target all domain
	//	err = t.DisableAllProfileTargetAll(user)
	//}

	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.IdentityError.Add.ToString(),
			})
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
		return
	}

	// Insert to database
	record.makeRow(inputs, user)
	err = mysql.Client.Create(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.IdentityError.Add.ToString(),
			})
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
		return
	}

	// Gán id input với id record vừa tạo được
	inputs.Id = record.Id

	// Tạo module userid info
	err = t.CreateModuleInfo(inputs)
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.IdentityError.ModuleUserId.ToString(),
			})
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
		return
	}

	// Tạo target
	err = t.CreateTarget(record.Id, user.Id, inputs)
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.IdentityError.Target.ToString(),
			})
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
		return
	}

	// Render Cache
	t.RenderCacheWithIdentity(record.Id, user.Id)
	// Get full lại record new
	record, _ = new(Identity).GetById(record.Id, user.Id)

	_ = history.PushHistory(&history.Identity{
		Detail:    history.DetailIdentityFE,
		CreatorId: user.Id,
		RecordOld: mysql.TableIdentity{},
		RecordNew: record.TableIdentity,
	})
	return
}

func (t *Identity) DisableAllProfile(user UserRecord) (err error) {
	err = mysql.Client.Model(&IdentityRecord{}).
		Where("identity.user_id = ? AND identity.status = ? AND identity.deleted_at is NULL", user.Id, mysql.TypeOn).
		Update("identity.status", mysql.TypeOff).Error
	return
}

func (t *Identity) DisableAllProfileTargetAll(user UserRecord) (err error) {
	profileTargetAllInventory := t.GetIdsProfileTargetAllInventory(user)
	if len(profileTargetAllInventory) > 0 {
		err = mysql.Client.Model(&IdentityRecord{}).
			Where("id IN ?", profileTargetAllInventory).
			Update("status", mysql.TypeOff).Error
	}
	return
}

func (t *Identity) Edit(inputs payload.IdentityCreate, user UserRecord, userAdmin UserRecord) (record IdentityRecord, errs []ajax.Error) {
	lang := lang.Translate
	var err error

	recordOld, _ := t.VerificationRecord(inputs.Id, user.Id)
	if recordOld.Id == 0 {
		errs = append(errs, ajax.Error{
			Id:      "id",
			Message: "You don't own this Identity",
		})
	}
	// Gán các giá trị cũ lại để khi save all những giá trị không thay đổi không bị mất
	record = recordOld
	// Validate inputs
	errs = t.ValidateCreate(inputs, user, lang)
	if len(errs) > 0 {
		return
	}

	// Nếu pub chuyển status từ OFF -> ON
	//if recordOld.Status == mysql.TypeOff && inputs.Status == mysql.TypeOn {
	//	if len(inputs.ListAdInventory) == 0 {
	//		// trường hợp target all domain thì sẽ disable toàn bộ các profile
	//		err = t.DisableAllProfile(user)
	//	} else {
	//		// trường hợp target 1 domain thì sẽ disable toàn bộ các profile target all domain
	//		err = t.DisableAllProfileTargetAll(user)
	//	}
	//	if err != nil {
	//		if !utility.IsWindow() {
	//			errs = append(errs, ajax.Error{
	//				Id:      "",
	//				Message: lang.Errors.IdentityError.Add.ToString(),
	//			})
	//			return
	//		}
	//		errs = append(errs, ajax.Error{
	//			Id:      "",
	//			Message: err.Error(),
	//		})
	//		return
	//	}
	//}

	// Insert to database
	record.makeRow(inputs, user)
	err = mysql.Client.Save(&record).Where("id = ?", record.Id).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.IdentityError.Edit.ToString(),
			})
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}

	// Xóa các row lưu module cũ
	_ = t.DeleteModuleInfo(record.Id)
	// Tạo module userid info
	err = t.CreateModuleInfo(inputs)
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.IdentityError.ModuleUserId.ToString(),
			})
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
		return
	}

	// Đoạn này để update lại cache vs những trường hợp domain khi chưa submit (ví dụ nếu pub xóa 1 domain trong target thì mình cũng phải clear Identity này ra khỏi cache)
	t.RenderCacheWithIdentity(record.Id, user.Id)

	err = new(Target).DeleteTarget(TargetRecord{mysql.TableTarget{
		IdentityId: record.Id,
	}})
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.IdentityError.Edit.ToString(),
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
				Message: lang.Errors.IdentityError.Target.ToString(),
			})
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}

	t.RenderCacheWithIdentity(record.Id, user.Id)
	// Get full lại record new
	record, _ = new(Identity).GetById(record.Id, user.Id)

	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = user.Id
	}
	_ = history.PushHistory(&history.Identity{
		Detail:    history.DetailIdentityFE,
		CreatorId: creatorId,
		RecordOld: recordOld.TableIdentity,
		RecordNew: record.TableIdentity,
	})
	return
}

func (t *Identity) ValidateModule(inputs payload.IdentityCreate) (errs []ajax.Error) {
	// var listIdTypeClient []int64
	for _, ModuleInfo := range inputs.ModuleParams {
		module := new(ModuleUserId).GetById(ModuleInfo.ModuleId)
		var Param []payload.ParamModuleUserId
		var Storage []payload.StorageModuleUserId
		json.Unmarshal([]byte(module.Params), &Param)
		json.Unmarshal([]byte(module.Storage), &Storage)

		for _, ModuleParam := range ModuleInfo.Params {
			idParam := strconv.FormatInt(ModuleInfo.ModuleId, 10) + "-" + ModuleParam.Name
			value := ModuleParam.Template

			for _, _Param := range Param {
				switch _Param.Type {
				case "int":
					if value != "" && ModuleParam.Name == _Param.Name {
						if !govalidator.IsInt(ModuleParam.Template) {
							errs = append(errs, ajax.Error{
								Id:      idParam,
								Message: "param " + ModuleParam.Name + " value is int",
							})
						}
					}
					break
				case "float":
					if value != "" && ModuleParam.Name == _Param.Name {
						if !govalidator.IsFloat(ModuleParam.Template) {
							errs = append(errs, ajax.Error{
								Id:      idParam,
								Message: "param " + ModuleParam.Name + " value is float",
							})
						}
					}
					break
				case "json":
					if value != "" && ModuleParam.Name == _Param.Name {
						if !govalidator.IsJSON(ModuleParam.Template) {
							errs = append(errs, ajax.Error{
								Id:      idParam,
								Message: "param " + ModuleParam.Name + " value is json",
							})
						}
					}
				case "boolean":
					if value != "" && ModuleParam.Name == _Param.Name {
						errs = append(errs, ajax.Error{
							Id:      idParam,
							Message: "param " + ModuleParam.Name + " value is true and false",
						})
					}
					break
				}
			}
		}

		for _, ModuleStorage := range ModuleInfo.Storage {
			idStorage := strconv.FormatInt(ModuleInfo.ModuleId, 10) + "-" + ModuleStorage.Name
			value := ModuleStorage.Template

			for _, _Storage := range Storage {
				switch _Storage.Type {
				case "int":
					if value != "" && ModuleStorage.Name == _Storage.Name {
						if !govalidator.IsInt(ModuleStorage.Template) {
							errs = append(errs, ajax.Error{
								Id:      idStorage,
								Message: "storage " + ModuleStorage.Name + " value is int",
							})
						}
					}
					break
				case "float":
					if value != "" && ModuleStorage.Name == _Storage.Name {
						if !govalidator.IsFloat(ModuleStorage.Template) {
							errs = append(errs, ajax.Error{
								Id:      idStorage,
								Message: "storage " + ModuleStorage.Name + " value is float",
							})
						}
					}
					break
				case "json":
					if value != "" && ModuleStorage.Name == _Storage.Name {
						if !govalidator.IsJSON(ModuleStorage.Template) {
							errs = append(errs, ajax.Error{
								Id:      idStorage,
								Message: "storage " + ModuleStorage.Name + " value is json",
							})
						}
					}
				case "boolean":
					if value != "" && ModuleStorage.Name == _Storage.Name {
						errs = append(errs, ajax.Error{
							Id:      idStorage,
							Message: "storage " + ModuleStorage.Name + " value is true and false",
						})
					}
					break
				}
			}
		}
	}
	return
}

func (t *Identity) Delete(id, userId int64, userAdmin UserRecord, lang lang.Translation) fiber.Map {
	row, er := t.GetById(id, userId)
	if er != nil {
		return fiber.Map{
			"status":  "err",
			"message": er.Error(),
			"id":      id,
		}
	}
	if row.IsDefault == 1 {
		return fiber.Map{
			"status":  "err",
			"message": "Identity default does not exist!",
			"id":      id,
		}
	}

	err := mysql.Client.Model(&IdentityRecord{}).Delete(&IdentityRecord{}, "id = ? and user_id = ?", id, userId).Error
	if err != nil {
		if !utility.IsWindow() {
			return fiber.Map{
				"status":  "err",
				"message": lang.Errors.IdentityError.Delete.ToString(),
				"id":      id,
			}
		}
		return fiber.Map{
			"status":  "err",
			"message": err.Error(),
			"id":      id,
		}
	}
	// Xóa các rls của identity
	_ = t.DeleteModuleInfo(id)

	// Render lại cache
	t.RenderCacheWithIdentity(id, userId)

	// History
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userId
	}
	_ = history.PushHistory(&history.Identity{
		Detail:    history.DetailIdentityFE,
		CreatorId: creatorId,
		RecordOld: row.TableIdentity,
		RecordNew: mysql.TableIdentity{},
	})

	return fiber.Map{
		"status":  "success",
		"message": "done",
		"id":      id,
	}
}

func (rec *IdentityRecord) makeRow(inputs payload.IdentityCreate, user UserRecord) {
	if inputs.Id == 0 { // Là trường hợp add
		if inputs.IsDefault == 0 {
			rec.IsDefault = mysql.TypeOff
		} else {
			rec.IsDefault = inputs.IsDefault
		}
	}
	rec.Id = inputs.Id
	rec.UserId = user.Id
	rec.Name = inputs.Name
	rec.Description = inputs.Description
	if rec.Description == "" {
		rec.Description = " "
	}
	rec.AuctionDelay = inputs.AuctionDelay
	rec.SyncDelay = inputs.SyncDelay
	rec.Status = inputs.Status
	rec.Priority = inputs.Priority
}

func (rec *IdentityRecord) makeRowEdit(inputs payload.IdentityCreate, user UserRecord) {
	rec.Id = inputs.Id
	rec.UserId = user.Id
	rec.Name = inputs.Name
	rec.Description = inputs.Description
	if rec.Description == "" {
		rec.Description = " "
	}
	rec.Status = inputs.Status
	rec.Priority = inputs.Priority
}

func (t *Identity) ValidateCreate(inputs payload.IdentityCreate, user UserRecord, lang lang.Translation) (errs []ajax.Error) {
	if utility.ValidateString(inputs.Name) == "" {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: lang.ErrorRequired.ToString(),
		})
	}

	if len(inputs.ModuleParams) > 0 {
		errs = append(errs, t.ValidateModule(inputs)...)
	} else {
		errs = append(errs, ajax.Error{
			Id:      "select-module",
			Message: lang.Errors.IdentityError.ChooseDomain.ToString(),
		})
	}

	if len(inputs.ListAdInventory) > 0 {
		var listInventoryTargeted []int64
		targets := new(Target).GetAllTargetIdentityValidate(user.Id, inputs.Id)
		for _, target := range targets {
			listInventoryTargeted = append(listInventoryTargeted, target.InventoryId)
		}
		textDomainError := ""
		for _, inventory := range inputs.ListAdInventory {
			if utility.InArray(inventory.Id, listInventoryTargeted, false) {
				if textDomainError == "" {
					textDomainError += inventory.Name
				} else if len(textDomainError) > 20 {
					textDomainError += " ,..."
					break
				} else {
					textDomainError += " ," + inventory.Name
				}
			}
		}
		if textDomainError != "" {
			errs = append(errs, ajax.Error{
				Id:      "text_for_domain",
				Message: textDomainError + lang.Errors.IdentityError.DomainTargeted.ToString() + " in another profile",
			})
		}
	} else {
		if inputs.IsDefault != mysql.TypeOn {
			errs = append(errs, ajax.Error{
				Id:      "text_for_domain",
				Message: lang.Errors.IdentityError.TargetDomain.ToString(),
			})
		}
	}
	return
}

func (t *Identity) ValidateEdit(inputs payload.IdentityCreate, userId int64) (errs []ajax.Error) {
	if utility.ValidateString(inputs.Name) == "" {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Name is required",
		})
	}

	if len(inputs.ModuleParams) > 0 {
		errs = append(errs, t.ValidateModule(inputs)...)
	} else {
		errs = append(errs, ajax.Error{
			Id:      "select-module",
			Message: "Choose at least one module",
		})
	}
	// if utility.ValidateString(inputs.Description) == "" {
	//	errs = append(errs, ajax.Error{
	//		Id:      "description",
	//		Message: "Description is required",
	//	})
	// }
	return
}

func (t *Identity) GetByFilters(inputs *payload.IdentityFilterPayload, userId int64, lang lang.Translation) (response datatable.Response, err error) {
	var inventories []IdentityRecord
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
			err = fmt.Errorf(lang.Errors.IdentityError.List.ToString())
		}
		return datatable.Response{}, err
	}
	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.MakeResponseDatatable(inventories)
	return
}

func (t *Identity) MakeResponseDatatable(Identitys []IdentityRecord) (records []IdentityRecordDatatable) {
	for _, identity := range Identitys {
		records = append(records, identity.Modify())
	}
	return
}

func (rec IdentityRecord) Modify() (record IdentityRecordDatatable) {
	record.Id = rec.Id
	record.Description = rec.Description
	//record.Priority = rec.Priority
	user := new(User).GetById(rec.UserId)
	record.User = user.Email

	var listNameInventoryTarget []string

	targets := new(Target).GetTargetIdentity(rec.Id)
	for _, target := range targets {
		if target.InventoryId != 0 {
			inventory, _ := new(Inventory).GetById(target.InventoryId, rec.UserId)
			listNameInventoryTarget = append(listNameInventoryTarget, inventory.Name)
		}
	}

	inventory := ListOfIdentity{
		Id:   rec.Id,
		List: listNameInventoryTarget,
	}

	record.Inventory = htmlblock.Render("identity/index/list_inventory.gohtml", inventory).String()
	record.Action = htmlblock.Render("identity/index/block.action.gohtml", rec).String()
	record.Name = htmlblock.Render("identity/index/name_block.gohtml", rec).String()
	if rec.IsDefault == 1 {
		record.Priority = `2 <span class="text-muted">(lowest)</span>`
	} else {
		record.Priority = `1 <span class="text-muted">(highest)</span>`
	}
	if rec.Status == 1 {
		record.Status = `<span class="">ON</span>`
	} else {
		record.Status = `<span class="">OFF</span>`
	}
	return
}

func (t *Identity) SetFilterStatus(inputs *payload.IdentityFilterPayload) func(db *gorm.DB) *gorm.DB {
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

func (t *Identity) setFilterSearch(inputs *payload.IdentityFilterPayload) func(db *gorm.DB) *gorm.DB {
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

func (t *Identity) setFilterTarget(inputs *payload.IdentityFilterPayload, userId int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var flag bool

		identityIDs := []int64{}
		targets, _ := new(Target).GetTargetByFilterIdentity(inputs, userId)
		if len(targets) > 0 {
			flag = true
		}
		if !flag {
			return db.Where("id IN ?", identityIDs)
		}
		for _, target := range targets {
			identityIDs = append(identityIDs, target.IdentityId)
		}
		return db.Where("id IN ?", identityIDs)
	}
}

func (t *Identity) setOrder(inputs *payload.IdentityFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var orders []string
		orders = append(orders, "is_default ASC")
		if len(inputs.Order) <= 0 {
			orders = append(orders, "id DESC")
		} else {
			for _, order := range inputs.Order {
				column := inputs.Columns[order.Column]
				orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
			}
		}
		orderString := strings.Join(orders, ", ")
		return db.Order(orderString)
	}
}

func (t *Identity) GetById(id, userId int64) (row IdentityRecord, err error) {
	err = mysql.Client.Where("id = ? and user_id = ?", id, userId).Find(&row).Error

	// Get rls
	row.TableIdentity.GetRls()
	return
}

func (t *Identity) GetByUser(userId int64) (records []IdentityRecord) {
	mysql.Client.Where("status = 1 and user_id = ?", userId).Find(&records)
	return
}

func (t *Identity) CreateTarget(identityId, userId int64, inputs payload.IdentityCreate) (err error) {
	all := int64(-1)
	// Kiểm tra nếu đầu vào input list target = 0 thì thêm một target = 0 thể hiện select all
	if len(inputs.ListAdInventory) == 0 {
		inputs.ListAdInventory = []payload.ListTarget{
			{
				Id: all,
			},
		}
	}

	for _, inventory := range inputs.ListAdInventory {
		var recordTarget TargetRecord
		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
			UserId:      userId,
			IdentityId:  identityId,
			InventoryId: inventory.Id,
		}}).Error
		if err != nil {
			return err
		}
	}
	//new(Inventory).UpdateRenderCacheWithIdentity(identityId, userId)
	return
}

func (t *Identity) VerificationRecord(id, userId int64) (row IdentityRecord, err error) {
	err = mysql.Client.Model(&IdentityRecord{}).Where("id = ? and user_id = ?", id, userId).Find(&row).Error
	// Get rls
	row.TableIdentity.GetRls()
	return
}

func (t *Identity) GetListBoxCollapse(userId, recordId int64, page, typ string) (list []string) {
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

func (t *Identity) CreateModuleInfo(input payload.IdentityCreate) (err error) {
	rows, err := new(IdentityModuleInfo).MakeRecord(input.ModuleParams, input.Id)
	if err != nil {
		return err
	}
	for _, row := range rows {
		err = new(IdentityModuleInfo).CreateModuleInfo(row)
		if err != nil {
			return err
		}
	}
	return
}

func (t *Identity) DeleteModuleInfo(identityId int64) (err error) {
	err = mysql.Client.Where(IdentityModuleInfoRecord{mysql.TableIdentityModuleInfo{
		IdentityId: identityId,
	}}).Delete(IdentityModuleInfoRecord{}).Error
	return
}

func (t *Identity) AutoCreateDefaultProfile(user UserRecord) (errs []ajax.Error) {
	layoutTime := "3:04:05 PM, January 2, 2006"
	timeConfigTime := time.Now().Format(layoutTime)
	identityCreate := payload.IdentityCreate{
		Name:         "Default Identity Profile - Target All Domains",
		Description:  "Auto create at " + timeConfigTime,
		AuctionDelay: 0,
		SyncDelay:    3000,
		Status:       mysql.TypeOn,
		IsDefault:    mysql.TypeOn,
		Priority:     2,
		ModuleParams: []payload.ModuleInfo{
			payload.ModuleInfo{
				ModuleId:   4,
				ModuleName: "criteo",
				Params:     []payload.ParamModuleUserId{},
				Storage:    []payload.StorageModuleUserId{},
				AbTesting:  0,
				Volume:     0,
			},
			payload.ModuleInfo{
				ModuleId:   1,
				ModuleName: "flocId",
				Params: []payload.ParamModuleUserId{payload.ParamModuleUserId{
					Name:     "token",
					Type:     "string",
					Template: "A3dHTSoNUMjjERBLlrvJSelNnwWUCwVQhZ5tNQ+sll7y+LkPPVZXtB77u2y7CweRIxiYaGw GXNlW1/dFp8VMEgIAAAB+eyJvcmlnaW4iOiJodHRwczovL3NoYXJlZGlkLm9yZzo0NDMiLC JmZWF0dXJlIjoiSW50ZXJlc3RDb2hvcnRBUEkiLCJleHBpcnkiOjE2MjYyMjA3OTksImlzU 3ViZG9tYWluIjp0cnVlLCJpc1RoaXJkUGFydHkiOnRydWV9",
				}},
				Storage:   []payload.StorageModuleUserId{},
				AbTesting: 0,
				Volume:    0,
			},
			payload.ModuleInfo{
				ModuleId:   3,
				ModuleName: "pubCommonId",
				Params:     []payload.ParamModuleUserId{},
				Storage: []payload.StorageModuleUserId{
					payload.StorageModuleUserId{
						Name:     "type",
						Type:     "string",
						Template: "cookie",
					},
					payload.StorageModuleUserId{
						Name:     "name",
						Type:     "string",
						Template: "_pubcid",
					},
					payload.StorageModuleUserId{
						Name:     "expires",
						Type:     "int",
						Template: "365",
					},
				},
				AbTesting: 0,
				Volume:    0,
			},
			payload.ModuleInfo{
				ModuleId:    5,
				ModuleName:  "id5Id",
				ModuleIndex: 0,
				Params: []payload.ParamModuleUserId{
					payload.ParamModuleUserId{
						Name:     "partner",
						Type:     "int",
						Template: "696",
					},
					payload.ParamModuleUserId{
						Name:     "pd",
						Type:     "string",
						Template: "",
					},
				},
				Storage: []payload.StorageModuleUserId{
					payload.StorageModuleUserId{
						Name:     "type",
						Type:     "string",
						Template: "html5",
					},
					payload.StorageModuleUserId{
						Name:     "name",
						Type:     "string",
						Template: "id5id",
					},
					payload.StorageModuleUserId{
						Name:     "expires",
						Type:     "int",
						Template: "90",
					},
					payload.StorageModuleUserId{
						Name:     "refreshInSeconds",
						Type:     "int",
						Template: "28800",
					},
				},
				AbTesting: 0,
				Volume:    0,
			},
			payload.ModuleInfo{
				ModuleId:   7,
				ModuleName: "sharedId",
				Params: []payload.ParamModuleUserId{
					payload.ParamModuleUserId{
						Name:     "syncTime",
						Type:     "int",
						Template: "60",
					},
				},
				Storage: []payload.StorageModuleUserId{
					payload.StorageModuleUserId{
						Name:     "type",
						Type:     "string",
						Template: "cookie",
					},
					payload.StorageModuleUserId{
						Name:     "name",
						Type:     "string",
						Template: "sharedid",
					},
					payload.StorageModuleUserId{
						Name:     "expires",
						Type:     "int",
						Template: "28",
					},
				},
				AbTesting: 0,
				Volume:    0,
			},
			payload.ModuleInfo{
				ModuleId:   8,
				ModuleName: "amxId",
				Params:     []payload.ParamModuleUserId{},
				Storage: []payload.StorageModuleUserId{
					payload.StorageModuleUserId{
						Name:     "name",
						Type:     "string",
						Template: "amxId",
					},
					payload.StorageModuleUserId{
						Name:     "type",
						Type:     "string",
						Template: "html5",
					},
					payload.StorageModuleUserId{
						Name:     "expires",
						Type:     "int",
						Template: "15",
					},
				},
				AbTesting: 0,
				Volume:    0,
			},
		},
		ListAdInventory: []payload.ListTarget{},
	}
	_, errs = t.Create(identityCreate, user, UserRecord{})
	//fmt.Println(record)
	return
}

func (t *Identity) RenderCacheWithIdentity(identityId int64, userId int64) {
	var listInventoryId []int64
	targets := new(Target).GetTargetIdentity(identityId)
	for _, target := range targets {
		if target.InventoryId != 0 {
			listInventoryId = append(listInventoryId, target.InventoryId)
		}
	}
	if len(listInventoryId) < 1 {
		return
	}
	if len(listInventoryId) == 1 {
		inventoryId := listInventoryId[0]
		// Nếu trường hợp inventory == -1 là all thì reset cache all
		if inventoryId == -1 {
			new(Inventory).ResetCacheAll(userId)
		} else {
			// ngược lại là trường hợp target cụ thể 1 domain thì reset cache domain ngay
			new(Inventory).ResetCacheWorker(inventoryId)
		}
	} else {
		for _, inventoryId := range listInventoryId {
			new(Inventory).ResetCacheWorker(inventoryId)
		}
	}
}
