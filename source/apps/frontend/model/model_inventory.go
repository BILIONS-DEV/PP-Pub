package model

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"source/apps/frontend/lang"
	"source/apps/frontend/payload"
	"source/apps/history"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"source/pkg/block"
	"source/pkg/datatable"
	"source/pkg/logger"
	"source/pkg/pagination"
	"source/pkg/scanAds"
	"source/pkg/telegram"
	"source/pkg/utility"
	"strconv"
	"strings"
	"time"
)

type Inventory struct{}

type InventoryRecord struct {
	mysql.TableInventory
}

func (InventoryRecord) TableName() string {
	return mysql.Tables.Inventory
}

func (t *Inventory) GetByFiltersForAdmin(inputs *payload.InventoryFilterPayload) (response datatable.Response, err error) {
	var inventories []InventoryRecord
	var total int64
	err = mysql.Client.
		Scopes(
			t.SetFilterStatus(inputs),
			t.setFilterSearch(inputs),
			t.SetFilterWebsiteLive(inputs),
			t.SetFilterAdsSync(inputs),
			t.SetFilterUser(inputs),
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
			err = fmt.Errorf(lang.Translate.Errors.InventoryError.List.ToString())
		}
		return datatable.Response{}, err
	}
	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.MakeResponseDatatable(inventories)
	return
}

func (t *Inventory) GetByFilters(inputs *payload.InventoryFilterPayload, userId int64, lang lang.Translation) (response datatable.Response, err error) {
	var inventories []InventoryRecord
	var total int64
	err = mysql.Client.Where("user_id = ?", userId).
		Scopes(
			t.SetFilterStatus(inputs),
			t.setFilterSearch(inputs),
			t.SetFilterWebsiteLive(inputs),
			t.SetFilterAdsSync(inputs),
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
			err = fmt.Errorf(lang.Errors.InventoryError.List.ToString())
		}
		return datatable.Response{}, err
	}
	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.MakeResponseDatatable(inventories)
	return
}

type InventoryRecordDatatable struct {
	InventoryRecord
	RowId        string `json:"DT_RowId"`
	Name         string `json:"name"`
	NameForAdmin string `json:"name_for_admin"`
	Status       string `json:"status"`
	Live         string `json:"live"`
	SyncAdsTxt   string `json:"sync_ads_txt"`
	Action       string `json:"action"`
}

func (t *Inventory) MakeResponseDatatable(inventories []InventoryRecord) (records []InventoryRecordDatatable) {
	for _, inventory := range inventories {
		rec := InventoryRecordDatatable{
			InventoryRecord: inventory,
			RowId:           strconv.FormatInt(inventory.Id, 10),
			Name:            block.RenderToString("supply/index/block.name.gohtml", inventory),
			NameForAdmin: block.RenderToString("supply/index/block.name-for-admin.gohtml", fiber.Map{
				"inventory":       inventory,
				"userOfInventory": new(User).GetEmailById(inventory.UserId),
			}),
			Status:     block.RenderToString("supply/index/block.status.gohtml", inventory),
			Live:       block.RenderToString("supply/index/block.live.gohtml", inventory),
			SyncAdsTxt: block.RenderToString("supply/index/block.sync_ads_txt.gohtml", inventory),
			Action:     block.RenderToString("supply/index/block.action.gohtml", inventory),
		}
		records = append(records, rec)
	}
	return
}

func (t *Inventory) SetFilterStatus(inputs *payload.InventoryFilterPayload) func(db *gorm.DB) *gorm.DB {
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

func (t *Inventory) setFilterSearch(inputs *payload.InventoryFilterPayload) func(db *gorm.DB) *gorm.DB {
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
		return db.Where("name  LIKE ?", "%"+inputs.PostData.QuerySearch+"%")
	}
}

func (t *Inventory) SetFilterWebsiteLive(inputs *payload.InventoryFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.WebLive != nil {
			switch inputs.PostData.WebLive.(type) {
			case string, int:
				if inputs.PostData.WebLive != "" {
					if inputs.PostData.WebLive == "1" {
						return db.Where("requests >= ?", 50)
					} else {
						return db.Where("requests < ?", 50)
					}
				}
			case []string, []interface{}:
				return db
			}
		}
		return db
	}
}

func (t *Inventory) SetFilterAdsSync(inputs *payload.InventoryFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.AdsTxtSync != nil {
			switch inputs.PostData.AdsTxtSync.(type) {
			case string, int:
				if inputs.PostData.AdsTxtSync != "" {
					return db.Where("sync_ads_txt = ?", inputs.PostData.AdsTxtSync)
				}
			case []string, []interface{}:
				return db.Where("sync_ads_txt IN ?", inputs.PostData.AdsTxtSync)
			}
		}
		return db
	}
}

func (t *Inventory) SetFilterUser(inputs *payload.InventoryFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.User != nil {
			switch inputs.PostData.User.(type) {
			case string, int:
				if inputs.PostData.User != "" {
					return db.Where("user_id = ?", inputs.PostData.User)
				}
			case []string, []interface{}:
				return db.Where("user_id IN ?", inputs.PostData.User)
			}
		}
		return db
	}
}

func (t *Inventory) setOrder(inputs *payload.InventoryFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(inputs.Order) > 0 {
			var orders []string
			for _, order := range inputs.Order {
				column := inputs.Columns[order.Column]
				if column.Data == "name_for_admin" {
					column.Data = "name"
				}
				if column.Data == "live" {
					column.Data = "requests"
				}
				orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
			}
			orderString := strings.Join(orders, ", ")
			return db.Order(orderString)
		}
		return db.Order("id desc")
	}
}

func (t *Inventory) Submit(inputs *payload.InventorySubmit, userId int64, lang lang.Translation) (errs []ajax.Error) {
	inventories := utility.SplitLines(inputs.Inventories)
	totalInventory := len(inventories)
	if totalInventory > 0 {
		limitGoroutine := make(chan bool, 10)
		respChannel := make(chan ajax.Error, totalInventory)
		for _, inventoryName := range inventories {
			limitGoroutine <- true
			go func(inventoryName string, respChannel chan ajax.Error) {
				inventoryName = strings.TrimSpace(inventoryName)
				// inventoryName empty
				if inventoryName == "" {
					respChannel <- ajax.Error{
						Id:      "",
						Message: "",
					}
					<-limitGoroutine
					return
				}
				// inventory validate
				inventoryName = t.ClearUpInventoryName(inventoryName)
				rootDomain := strings.TrimSpace(inventoryName)
				rootDomain, err := utility.GetRootDomain(inventoryName)
				if err != nil {
					if !utility.IsWindow() {
						respChannel <- ajax.Error{
							Id:      inventoryName,
							Message: lang.Errors.InventoryError.RootDomain.ToString(),
						}
					} else {
						respChannel <- ajax.Error{
							Id:      inventoryName,
							Message: err.Error(),
						}
					}
					<-limitGoroutine
					return
				}

				// validate
				var record InventoryRecord
				mysql.Client.Select("id").
					Where(InventoryRecord{mysql.TableInventory{Name: rootDomain, UserId: userId}}).
					Last(&record)
				if record.Id > 0 {
					respChannel <- ajax.Error{
						Id:      inventoryName,
						Message: lang.Pages.Inventory.SubmitDomain.Errors.AlreadyExist.ToUpperFirstCharacter(),
					}
					<-limitGoroutine
					return
				}
				// check unscoped record trong trường domain đã xóa nhưng pub muốn update lại trong user
				mysql.Client.Unscoped().
					Where(InventoryRecord{mysql.TableInventory{Name: rootDomain, UserId: userId}}).
					Last(&record)
				// Nếu inventory đã tồn tại khi bỏ deleted_at thì tiến hành update
				var errI error
				if record.Id != 0 {
					record.Status = mysql.StatusPending
					record.DeletedAt.Valid = false
					errI = mysql.Client.Unscoped().Updates(&record).Error
				} else {
					record = InventoryRecord{mysql.TableInventory{
						UserId:         userId,
						Type:           mysql.InventoryTypeWeb,
						Status:         mysql.StatusPending,
						Name:           rootDomain,
						Domain:         rootDomain,
						Uuid:           uuid.New().String(),
						IabCategories:  "",
						LastScanAdsTxt: sql.NullTime{},
						CachedAt:       time.Now(),
						AdsTxtCustom:   "",
					}}
					errI = mysql.Client.Create(&record).Error
				}
				if errI != nil {
					if !utility.IsWindow() {
						respChannel <- ajax.Error{
							Id:      inventoryName,
							Message: lang.Errors.InventoryError.Submit.ToString(),
						}
					} else {
						respChannel <- ajax.Error{
							Id:      inventoryName,
							Message: errI.Error(),
						}
					}
					<-limitGoroutine
					return
				}
				new(InventoryConfig).MakeRowDefault(record.Id)
				_ = history.PushHistory(&history.Inventory{
					Detail:    history.DetailInventorySubmitFE,
					CreatorId: userId,
					RecordNew: record.TableInventory,
				})
				if !utility.IsWindow() {
					userModel := new(User).GetById(userId)
					managerModel := new(User).GetById(userModel.Presenter)
					// Gửi thông báo domain mới tạo vào nhóm sale telegram
					_ = telegram.SendMessageGroupPubPowerNotify(record.Name, userModel.Email, managerModel.Email)
				}
				// success
				respChannel <- ajax.Error{
					Id:      inventoryName,
					Message: "",
				}
				<-limitGoroutine
				return
			}(inventoryName, respChannel)
		}
		// consumer get data from channel
		for i := 1; i <= totalInventory; i++ {
			resp := <-respChannel
			if resp.Id == "" && resp.Message == "" {
				continue
			}
			errs = append(errs, resp)
		}
	}
	return
}

func (t *Inventory) ValidateInventory(inventoryName string) (output string, err error) {
	output, err = utility.GetRootDomain(inventoryName)
	return
}

func (t *Inventory) ClearUpInventoryName(inventoryName string) (output string) {
	output = strings.ReplaceAll(inventoryName, "https://", "")
	output = strings.ReplaceAll(output, "https://", "")
	return
}

func (t *Inventory) GetAll() (records []InventoryRecord) {
	mysql.Client.Find(&records)
	return
}

func (t *Inventory) GetByUser(userId int64) (records []InventoryRecord) {
	mysql.Client.Where("user_id = ?", userId).Find(&records)
	return
}

func (t *Inventory) GetByUserId(userId int64) (row []InventoryRecord, err error) {
	err = mysql.Client.Where("user_id = ?", userId).Unscoped().Find(&row).Error
	return
}

func (t *Inventory) CountData(value string, userId int64) (count int64) {
	mysql.Client.Model(&InventoryRecord{}).Where("name like ? and user_id = ?", "%"+value+"%", userId).Count(&count)
	return
}

func (t *Inventory) CountDataSystem(value string) (count int64) {
	mysql.Client.Model(&InventoryRecord{}).Where("name like ?", "%"+value+"%").Count(&count)
	return
}

func (t *Inventory) CountDataPageEdit(userId int64, listId []int64) (count int64) {
	if len(listId) > 0 {
		mysql.Client.Model(&InventoryRecord{}).Where("user_id = ? and id not in ?", userId, listId).Count(&count)
	} else {
		mysql.Client.Model(&InventoryRecord{}).Where("user_id = ?", userId).Count(&count)
	}

	return
}

func (t *Inventory) CountDataPageEditSystem(listId []int64) (count int64) {
	if len(listId) > 0 {
		mysql.Client.Model(&InventoryRecord{}).Where("id not in ?", listId).Count(&count)
	} else {
		mysql.Client.Model(&InventoryRecord{}).Count(&count)
	}

	return
}

func (t *Inventory) LoadMoreData(key, value string, userId int64, listSelected []int64) (rows []InventoryRecord, isMoreData, lastPage bool) {
	limit := 10
	page, offset := pagination.Pagination(key, limit)
	if len(listSelected) > 0 {
		mysql.Client.Where("name like ? AND user_id = ? AND id NOT IN ? AND status = ?", "%"+value+"%", userId, listSelected, mysql.StatusApproved).
			Limit(limit).Offset(offset).Find(&rows)
	} else {
		mysql.Client.Where("name like ? AND user_id = ? AND status = ?", "%"+value+"%", userId, mysql.StatusApproved).Limit(limit).Offset(offset).Find(&rows)
	}
	total := t.CountData(value, userId)
	totalPages := int(total) / limit
	if (int(total) % limit) != 0 {
		totalPages++
	}
	if page < totalPages {
		isMoreData = true
	}
	if page >= totalPages || len(rows) == 0 {
		isMoreData = false
		lastPage = true
	}
	return
}

func (t *Inventory) LoadMoreDataSystem(key, value string, listSelected []int64) (rows []InventoryRecord, isMoreData, lastPage bool) {
	limit := 10
	page, offset := pagination.Pagination(key, limit)
	if len(listSelected) > 0 {
		mysql.Client.Where("name LIKE ? AND id NOT IN ? AND status = ?", "%"+value+"%", listSelected, mysql.StatusApproved).Limit(limit).Offset(offset).Find(&rows)
	} else {
		mysql.Client.Where("name LIKE ? AND status = ?", "%"+value+"%", mysql.StatusApproved).Limit(limit).Offset(offset).Find(&rows)
	}
	total := t.CountDataSystem(value)
	totalPages := int(total) / limit
	if (int(total) % limit) != 0 {
		totalPages++
	}
	if page < totalPages {
		isMoreData = true
	}
	if page >= totalPages || len(rows) == 0 {
		isMoreData = false
		lastPage = true
	}
	return
}

func (t *Inventory) LoadMoreDataPageEdit(userId int64, listSelected []int64) (rows []InventoryRecord, isMoreData, lastPage bool) {
	limit := 10
	page, offset := pagination.Pagination("1", limit)
	if len(listSelected) > 0 {
		mysql.Client.Where("user_id = ? and id not in ? and status = ?", userId, listSelected, mysql.StatusApproved).Find(&rows)
	} else {
		mysql.Client.Where("user_id = ? and status = ?", userId, mysql.StatusApproved).Limit(limit).Offset(offset).Find(&rows)
	}
	if len(rows) > 10 {
		rows = rows[0:9]
	}
	total := t.CountDataPageEdit(userId, listSelected)
	totalPages := int(total) / limit
	if (int(total) % limit) != 0 {
		totalPages++
	}
	if page < totalPages {
		isMoreData = true
	}
	if page >= totalPages || len(rows) == 0 {
		isMoreData = false
		lastPage = true
	}
	return
}

func (t *Inventory) LoadMoreDataPageEditSystem(listSelected []int64) (rows []InventoryRecord, isMoreData, lastPage bool) {
	limit := 10
	page, offset := pagination.Pagination("1", limit)
	if len(listSelected) > 0 {
		mysql.Client.Where("id not in ?", listSelected).Find(&rows)
	} else {
		mysql.Client.Limit(limit).Offset(offset).Find(&rows)
	}
	if len(rows) > 10 {
		rows = rows[0:9]
	}
	total := t.CountDataPageEditSystem(listSelected)
	totalPages := int(total) / limit
	if (int(total) % limit) != 0 {
		totalPages++
	}
	if page < totalPages {
		isMoreData = true
	}
	if page >= totalPages || len(rows) == 0 {
		isMoreData = false
		lastPage = true
	}
	return
}

func (t *Inventory) GetById(id, userId int64) (row InventoryRecord, err error) {
	err = mysql.Client.Where("id = ? and user_id = ?", id, userId).Find(&row).Error
	row.GetFullData()
	return
}

func (t *Inventory) GetByIdSystem(id int64) (row InventoryRecord, err error) {
	err = mysql.Client.Where("id = ?", id).Find(&row).Error
	return
}

func (t *Inventory) GetByUserIdLimit(userId int64, limit int) (row []InventoryRecord, err error) {
	err = mysql.Client.Where("user_id = ?", userId).Limit(limit).Find(&row).Error
	return
}

func (t *Inventory) GetByUserIdSearch(userId int64, search string) (row []InventoryRecord, err error) {
	err = mysql.Client.Where("user_id = ? and name LIKE ?", userId, "%"+search+"%").Limit(20).Find(&row).Error
	return
}

func (t *Inventory) GetAllIdsOfUser(userId int64) (ids []int64) {
	var inventories []InventoryRecord
	mysql.Client.Select("id").Where("user_id = ? AND status = ?", userId, mysql.StatusApproved).Find(&inventories)
	for _, inventory := range inventories {
		ids = append(ids, inventory.Id)
	}
	return
}

type AdsTxtMissingLine struct {
	Status bool
	Text   string
}

func (rec InventoryRecord) GetAllMissingAdsTxt() (listMissing []AdsTxtMissingLine, syncError string) {
	// samples := Samples()
	samples, _ := rec.GetAdsTxt()
	missingRecords := new(MissingAdsTxt).GetByInventory(rec)
	// Nếu không tìm thấy row nào trong database => có 2 trường hợp.
	// 		1. Nếu last_scan_ads = nil -> chưa được scan lần nào -> tất cả đều đc đánh dấu là chưa thành công.
	// 		2. Nếu last_scan_ads != nil -> đã được scan và k có row missing nào -> tất cả được đánh dấu là đã thành công
	if len(missingRecords) == 0 {
		var status = true
		if !rec.LastScanAdsTxt.Valid { // Nếu chưa được scan thì tất cả các dòng mặc định là false
			status = false
		}
		for _, lineTxt := range samples {
			if strings.TrimSpace(lineTxt) == "" {
				continue
			}
			index := strings.Index(lineTxt, "#")
			if index == 0 || index == 1 {
				continue
			}
			listMissing = append(listMissing, AdsTxtMissingLine{
				Status: status,
				Text:   lineTxt,
			})
		}
		return
	}
	// Nếu chỉ lấy được 1 record trong db & record đó lại là scan_ads_url lỗi thì đánh dấu tất cả là chưa scan
	// Đây là trường hợp link ads.txt của pub bị lỗi, k quét được.
	if len(missingRecords) == 1 && missingRecords[0].ErrorMessage != "" {
		for _, lineTxt := range samples {
			if strings.TrimSpace(lineTxt) == "" {
				continue
			}
			index := strings.Index(lineTxt, "#")
			if index == 0 || index == 1 {
				continue
			}
			listMissing = append(listMissing, AdsTxtMissingLine{
				Status: false,
				Text:   lineTxt,
			})
		}
		syncError = missingRecords[0].ErrorMessage
		return
	}
	// Trường hợp nếu lấy được nhiều hơn 1 row missing thì sẽ kiểm tra từng dòng xem dòng nào bị thiếu thì dánh dấu status=false
	for _, lineTxt := range samples {
		if strings.TrimSpace(lineTxt) == "" {
			continue
		}
		index := strings.Index(lineTxt, "#")
		if index == 0 || index == 1 {
			continue
		}
		status := true
		for _, missingRecord := range missingRecords {
			if lineTxt == missingRecord.Line {
				status = false
			}
		}
		listMissing = append(listMissing, AdsTxtMissingLine{
			Status: status,
			Text:   lineTxt,
		})
	}
	return

}

func (rec InventoryRecord) GetAdsTxtUrls() (listURL []string) {
	if rec.AdsTxtUrl != "" {
		return []string{rec.AdsTxtUrl}
	}
	listURL = append(listURL, fmt.Sprintf("https://%s/ads.txt", rec.Domain))
	listURL = append(listURL, fmt.Sprintf("https://www.%s/ads.txt", rec.Domain))
	listURL = append(listURL, fmt.Sprintf("http://%s/ads.txt", rec.Domain))
	listURL = append(listURL, fmt.Sprintf("http://www.%s/ads.txt", rec.Domain))
	return
}

func (rec InventoryRecord) ScanAdsTxt() (err error) {
	// Get Ads.txt url of domain
	adsTxtUrls := rec.GetAdsTxtUrls()
	if len(adsTxtUrls) == 0 {
		return errors.New("ads.txt url not found")
	}
	// Scan Ads.txt
	scanAdsTxtResponse := scanAds.Response{}
	samples, _ := rec.GetAdsTxt()
	for _, adsTxtUrl := range adsTxtUrls {
		// scanAdsTxtResponse = scanAds.ScanUrl(adsTxtUrl, Samples())
		scanAdsTxtResponse = scanAds.ScanUrl(adsTxtUrl, samples)
		if !scanAdsTxtResponse.ErrorStatus {
			rec.AdsTxtUrl = adsTxtUrl
			break
		}
	}
	// If error
	if scanAdsTxtResponse.ErrorStatus {
		err = new(MissingAdsTxt).ScanErrorSaved(rec, scanAdsTxtResponse)
		if err != nil {
			return
		}
		err = rec.AdsTxtScanned(mysql.InventorySyncAdsTxtError)
		if err != nil {
			return
		}
		return scanAdsTxtResponse.Error
	}
	// Update table missing_ads_txt
	err, syncStatus := new(MissingAdsTxt).ScanSuccessSaved(rec, scanAdsTxtResponse)
	if err != nil {
		return
	}
	// Save ads_txt_url & last_scan_ads_txt
	err = rec.AdsTxtScanned(syncStatus)
	if err != nil {
		return
	}
	return
}

func (rec *InventoryRecord) AdsTxtScanned(SyncAdsTxt mysql.TYPEInventorySyncAdsTxt) (err error) {
	data := InventoryRecord{mysql.TableInventory{
		LastScanAdsTxt: sql.NullTime{Time: time.Now(), Valid: true},
		AdsTxtUrl:      rec.AdsTxtUrl,
		SyncAdsTxt:     SyncAdsTxt,
	}}
	err = mysql.Client.Model(&rec).Updates(data).Error
	return
}

func (t *Inventory) GetByIdForFilter(id, userId int64) (row InventoryRecord, err error) {
	err = mysql.Client.Unscoped().Where("id = ? and user_id = ? and status = ?", id, userId, mysql.StatusApproved).Find(&row).Error
	if row.Id == 0 {
		err = errors.New("Record not found")
	}
	return
}

func (t *Inventory) GetByIdForFilterSystem(id int64) (row InventoryRecord, err error) {
	err = mysql.Client.Unscoped().Where("id = ?", id).Find(&row).Error
	if row.Id == 0 {
		err = errors.New("Record not found")
	}
	return
}

func (t *Inventory) ValidateSetup(inputs payload.GeneralInventory) (errs []ajax.Error) {

	if inputs.PrebidTimeout == 0 || inputs.PrebidTimeout < 0 {
		errs = append(errs, ajax.Error{
			Id:      "prebid_timeout",
			Message: "PrebidTimeout is required",
		})
	}

	if inputs.AdRefresh == 1 {
		if inputs.AdRefreshTime == 0 || inputs.AdRefreshTime < 0 {
			errs = append(errs, ajax.Error{
				Id:      "ad_refresh_time",
				Message: "Ad Refresh Time is required",
			})
		}
	}

	if inputs.Gdpr == 1 {
		if inputs.GdprTimeout == 0 || inputs.GdprTimeout < 0 {
			errs = append(errs, ajax.Error{
				Id:      "gdpr_timeout",
				Message: "Gdpr Timeout is required",
			})
		}
	}

	if inputs.Ccpa == 1 {
		if inputs.CcpaTimeout == 0 || inputs.CcpaTimeout < 0 {
			errs = append(errs, ajax.Error{
				Id:      "ccpa_timeout",
				Message: "Ccpa Timeout is required",
			})
		}
	}

	if inputs.CustomBrand == 1 {
		if utility.ValidateString(inputs.Logo) == "" {
			errs = append(errs, ajax.Error{
				Id:      "custom_brand_logo",
				Message: "Logo is required",
			})
		}
		if utility.ValidateString(inputs.Title) == "" {
			errs = append(errs, ajax.Error{
				Id:      "custom_brand_title",
				Message: "Title is required",
			})
		}
		if utility.ValidateString(inputs.Content) == "" {
			errs = append(errs, ajax.Error{
				Id:      "custom_brand_content",
				Message: "Content is required",
			})
		}
	}
	return
}

func (t *Inventory) Delete(id, userId int64, userAdmin UserRecord, lang lang.Translation) fiber.Map {
	inventory, _ := new(Inventory).GetById(id, userId)
	// Xóa và cập nhập trạng thái reject cho domain
	err := mysql.Client.Model(&InventoryRecord{}).
		Where("id = ? and user_id = ?", id, userId).
		Updates(&InventoryRecord{mysql.TableInventory{DeletedAt: gorm.DeletedAt{
			Time:  time.Now(),
			Valid: true},
			Status: mysql.StatusReject}}).
		Error
	if err != nil {
		logger.Error(err.Error())
		return fiber.Map{
			"status":  "err",
			"message": "error",
			"id":      id,
		}
	}
	// Đặt tất cả ad tag trong domain về trạng thái not live
	err = mysql.Client.Model(&InventoryAdTagRecord{}).
		Where("inventory_id = ? and user_id = ?", id, userId).
		Updates(&InventoryAdTagRecord{mysql.TableInventoryAdTag{
			Status: mysql.TypeStatusAdTagNotLive}}).
		Error
	t.ResetCacheWorker(id)
	new(InventoryAdTag).DeleteTagAfterDeleteDomain(id, userId)
	new(InventoryConfig).DeleteConfigAfterDeleteDomain(id)
	if err != nil {
		if !utility.IsWindow() {
			return fiber.Map{
				"status":  "err",
				"message": lang.Errors.InventoryError.Delete.ToString(),
				"id":      id,
			}
		}
		return fiber.Map{
			"status":  "err",
			"message": err.Error(),
			"id":      id,
		}
	}
	// History
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userId
	}
	_ = history.PushHistory(&history.Inventory{
		Detail:    history.DetailInventoryDeleteFE,
		CreatorId: creatorId,
		RecordOld: inventory.TableInventory,
		RecordNew: mysql.TableInventory{},
	})
	return fiber.Map{
		"status":  "success",
		"message": "done",
		"id":      id,
	}
}

func (t *Inventory) GetListBoxCollapse(userId, inventoryId int64, page string) (list []string) {
	mysql.Client.Select("box_collapse").Model(&PageCollapseRecord{}).Where("user_id = ? and page_collapse = ? and is_collapse = ? and page_id = ?", userId, page, 1, inventoryId).Find(&list)
	return
}

func (t *Inventory) UpdateRenderCacheWithLineItem(lineItemId int64, userId int64) {
	var listInventory []int64
	var listAdFormat []int64
	var listSize []int64
	var listAdTag []int64
	mysql.Client.Model(&TargetRecord{}).Select("inventory_id").Where("line_item_id = ? and inventory_id != 0", lineItemId).Find(&listInventory)
	mysql.Client.Model(&TargetRecord{}).Select("ad_format_id").Where("line_item_id = ? and ad_format_id != 0", lineItemId).Find(&listAdFormat)
	mysql.Client.Model(&TargetRecord{}).Select("ad_size_id").Where("line_item_id = ? and ad_size_id != 0", lineItemId).Find(&listSize)
	mysql.Client.Model(&TargetRecord{}).Select("tag_id").Where("line_item_id = ? and tag_id != 0", lineItemId).Find(&listAdTag)
	// Khởi tạo giá trị rỗng nếu không tìm record trong trường hợp lỗi
	if len(listInventory) == 0 {
		listInventory = append(listInventory, 0)
	}
	if len(listAdFormat) == 0 {
		listAdFormat = append(listAdFormat, 0)
	}
	if len(listSize) == 0 {
		listSize = append(listSize, 0)
	}
	if len(listAdTag) == 0 {
		listAdTag = append(listAdTag, 0)
	}
	// Tìm kiếm các inventory phù hợp với target để render cache
	if listInventory[0] == -1 {
		if listAdTag[0] == -1 {
			if listAdFormat[0] == -1 {
				if listSize[0] == -1 {
					t.ResetCacheAll(userId)
				} else {
					for _, sizeId := range listSize {
						var adTags []InventoryAdTagRecord
						mysql.Client.Where("primary_ad_size = ? and user_id = ?", sizeId, userId).Find(&adTags)
						for _, adTag := range adTags {
							t.ResetCacheWorker(adTag.InventoryId)
						}
					}
				}
			} else {
				if listSize[0] == -1 {
					for _, formatId := range listAdFormat {
						var adTags []InventoryAdTagRecord
						mysql.Client.Where("type = ? and user_id = ?", formatId, userId).Find(&adTags)
						for _, adTag := range adTags {
							t.ResetCacheWorker(adTag.InventoryId)
						}
					}
				} else {
					for _, formatId := range listAdFormat {
						for _, sizeId := range listSize {
							var adTags []InventoryAdTagRecord
							mysql.Client.Where("primary_ad_size = ? and type = ? and user_id = ?", sizeId, formatId, userId).Find(&adTags)

							for _, adTag := range adTags {
								t.ResetCacheWorker(adTag.InventoryId)
							}
						}
					}
				}
			}
		} else {
			for _, adTagId := range listAdTag {
				adTagRecord := new(InventoryAdTag).GetById(adTagId)
				t.ResetCacheWorker(adTagRecord.InventoryId)
			}
		}
	} else {
		for _, inventoryId := range listInventory {
			t.ResetCacheWorker(inventoryId)
		}
	}
	return
}

func (t *Inventory) UpdateRenderCacheWithFloor(floorId int64, userId int64) {
	var listInventory []int64
	var listAdFormat []int64
	var listSize []int64
	var listAdTag []int64
	mysql.Client.Model(&TargetRecord{}).Select("inventory_id").Where("floor_id = ? and inventory_id != 0", floorId).Find(&listInventory)
	mysql.Client.Model(&TargetRecord{}).Select("ad_format_id").Where("floor_id = ? and ad_format_id != 0", floorId).Find(&listAdFormat)
	mysql.Client.Model(&TargetRecord{}).Select("ad_size_id").Where("floor_id = ? and ad_size_id != 0", floorId).Find(&listSize)
	mysql.Client.Model(&TargetRecord{}).Select("tag_id").Where("floor_id = ? and tag_id != 0", floorId).Find(&listAdTag)
	// Khởi tạo giá trị rỗng nếu không tìm record trong trường hợp lỗi
	if len(listInventory) == 0 {
		listInventory = append(listInventory, 0)
	}
	if len(listAdFormat) == 0 {
		listAdFormat = append(listAdFormat, 0)
	}
	if len(listSize) == 0 {
		listSize = append(listSize, 0)
	}
	if len(listAdTag) == 0 {
		listAdTag = append(listAdTag, 0)
	}
	// Tìm kiếm các inventory phù hợp với target để render cache
	if listInventory[0] == -1 {
		if listAdTag[0] == -1 {
			if listAdFormat[0] == -1 {
				if listSize[0] == -1 {
					t.ResetCacheAll(userId)
				} else {
					for _, sizeId := range listSize {
						var adTags []InventoryAdTagRecord
						mysql.Client.Where("primary_ad_size = ? and user_id = ?", sizeId, userId).Find(&adTags)
						for _, adTag := range adTags {
							t.ResetCacheWorker(adTag.InventoryId)
						}
					}
				}
			} else {
				if listSize[0] == -1 {
					for _, formatId := range listAdFormat {
						var adTags []InventoryAdTagRecord
						mysql.Client.Where("type = ? and user_id = ?", formatId, userId).Find(&adTags)
						for _, adTag := range adTags {
							t.ResetCacheWorker(adTag.InventoryId)
						}
					}
				} else {
					for _, formatId := range listAdFormat {
						for _, sizeId := range listSize {
							var adTags []InventoryAdTagRecord
							mysql.Client.Where("primary_ad_size = ? and type = ? and user_id = ?", sizeId, formatId, userId).Find(&adTags)

							for _, adTag := range adTags {
								t.ResetCacheWorker(adTag.InventoryId)
							}
						}
					}
				}
			}
		} else {
			for _, adTagId := range listAdTag {
				adTagRecord := new(InventoryAdTag).GetById(adTagId)
				t.ResetCacheWorker(adTagRecord.InventoryId)
			}
		}
	} else {
		for _, inventoryId := range listInventory {
			t.ResetCacheWorker(inventoryId)
		}
	}

	return
}

func (t *Inventory) ResetCacheWorker(inventoryId int64) {
	mysql.Client.Model(&InventoryRecord{}).Where("id = ?", inventoryId).Update("render_cache", 1)
}

func (t *Inventory) ResetCacheAll(userId int64) {
	mysql.Client.Model(&InventoryRecord{}).Where("user_id = ?", userId).Update("render_cache", 1)
}

func (t *Inventory) GetByName(name string) (rows []InventoryRecord, err error) {
	err = mysql.Client.Where("name like ?", name).Find(&rows).Error
	return
}
