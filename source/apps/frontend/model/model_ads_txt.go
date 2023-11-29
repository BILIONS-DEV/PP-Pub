package model

import (
	"fmt"
	"gorm.io/gorm"
	"source/apps/frontend/lang"
	"source/apps/frontend/payload"
	"source/apps/history"
	"source/core/technology/mysql"
	"source/pkg/adstxt"
	"source/pkg/datatable"
	"source/pkg/htmlblock"
	"source/pkg/pagination"
	"source/pkg/utility"
	"strings"
)

type AdsTxt struct{}

type AdsTxtRecord struct {
	mysql.TableInventory
}

func (AdsTxtRecord) TableName() string {
	return mysql.Tables.Inventory
}

type AdsTxtDatatable struct {
	Url          string `json:"name"`
	AdsTxtStatus string `json:"sync_ads_txt"`
	Action       string `json:"action"`
	LastScanned  string `json:"last_scan_ads_txt"`
}

func (t *AdsTxt) GetByFilters(inputs *payload.InventoryFilterPayload, userId int64, lang lang.Translation) (response datatable.Response, err error) {
	var inventories []AdsTxtRecord
	var total int64
	err = mysql.Client.Where("user_id = ?", userId).
		Scopes(
			t.setFilterStatus(inputs),
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

func (t *AdsTxt) setFilterStatus(inputs *payload.InventoryFilterPayload) func(db *gorm.DB) *gorm.DB {
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

func (t *AdsTxt) setFilterSearch(inputs *payload.InventoryFilterPayload) func(db *gorm.DB) *gorm.DB {
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

func (t *AdsTxt) setOrder(inputs *payload.InventoryFilterPayload) func(db *gorm.DB) *gorm.DB {
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
		return db.Order("id desc")
	}
}

func (t *AdsTxt) MakeResponseDatatable(inventories []AdsTxtRecord) (records []AdsTxtDatatable) {
	for _, inventory := range inventories {
		rec := AdsTxtDatatable{
			Url:          inventory.Name,
			AdsTxtStatus: htmlblock.Render("ads_txt/block_html/block.sync_ads_txt.gohtml", inventory).String(),
			Action:       htmlblock.Render("ads_txt/block_html/action.gohtml", inventory).String(),
			LastScanned:  `<span class="fw-12">` + inventory.GetDateScan() + `</span>`,
		}
		records = append(records, rec)
	}
	return
}

func (row *AdsTxtRecord) GetDateScan() (date string) {
	layout := "01/02/2006 at 15:04 PM" //m/d/y
	if !row.LastScanAdsTxt.Valid {
		return `<span class="text-muted">--</span>`
	} else {
		date = row.LastScanAdsTxt.Time.Format(layout)
	}
	return
}

func (t *AdsTxt) PushForInventory(inventory InventoryRecord, adsTxtContent string, lang lang.Translation) (adsTxtCustom string, err error) {
	// Luu record old cho history
	recordOld := inventory
	// Lấy ra các line hệ thống đang dùng
	//lineSystem, _ := inventory.GetAdsTxt()
	// Check duplicate bỏ các dòng đã có trong hệ thống
	lineCustom := adstxt.StandardizedWithText(adsTxtContent)
	//var lineSave []string
	//for _, line := range lineCustom {
	//	if !utility.InArray(line, lineSystem, false) {
	//		lineSave = append(lineSave, line)
	//	}
	//}
	// Save adstxt custom của pub
	adsTxtCustom = adstxt.Standardized(lineCustom).ToString()
	err = mysql.Client.Table(mysql.Tables.Inventory).Model(&inventory).Update("ads_txt_custom", adsTxtCustom).Error
	if err != nil && !utility.IsWindow() {
		err = fmt.Errorf(lang.Errors.InventoryError.SaveAdsTxt.ToString())
		return
	}
	// Push history
	inventory.AdsTxtCustom = mysql.TYPEAdsTxtCustom(adsTxtCustom)
	recordNew := inventory
	_ = history.PushHistory(&history.Inventory{
		Detail:    history.DetailInventoryAdsTxtFE,
		CreatorId: inventory.UserId,
		RecordOld: recordOld.TableInventory,
		RecordNew: recordNew.TableInventory,
	})
	return
}
