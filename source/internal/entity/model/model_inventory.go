package model

import (
	"database/sql"
	"encoding/json"
	"gorm.io/gorm"
	"strings"
	"time"
)

func (InventoryModel) TableName() string {
	return "inventory"
}

// Inventory : cấu hình các trường cho table `inventory` trên mysql
type InventoryModel struct {
	ID                  int64                            `gorm:"column:id" json:"id"`
	UserID              int64                            `gorm:"column:user_id" json:"user_id"`
	Presenter           int64                            `gorm:"column:presenter" json:"presenter"`
	Name                string                           `gorm:"column:name" json:"name"`
	Domain              string                           `gorm:"column:domain" json:"domain"`
	Type                TYPEInventoryType                `gorm:"column:type" json:"type"`
	Status              TYPEStatus                       `gorm:"column:status" json:"status"`
	JsMode              string                           `gorm:"column:js_mode" json:"js_mode"`
	VpaidMode           string                           `gorm:"column:vpaid_mode" json:"vpaid_mode"`
	PrebidJs            string                           `gorm:"column:prebid_js" json:"prebid_js"`
	SyncAdsTxt          TYPEInventorySyncAdsTxt          `gorm:"column:sync_ads_txt" json:"sync_ads_txt"`
	LastScanAdsTxt      sql.NullTime                     `gorm:"column:last_scan_ads_txt" json:"last_scan_ads_txt"`
	AdsTxtUrl           string                           `gorm:"column:ads_txt_url" json:"ads_txt_url"`
	AdsTxtCustom        TYPEAdsTxtCustom                 `gorm:"column:ads_txt_custom" json:"ads_txt_custom"`
	AdsTxtCustomByAdmin TYPEAdsTxtCustom                 `gorm:"column:ads_txt_custom_by_admin" json:"ads_txt_custom_by_admin"`
	CreatedAt           time.Time                        `gorm:"column:created_at" json:"created_at"`
	DeletedAt           gorm.DeletedAt                   `gorm:"column:deleted_at" json:"deleted_at"`
	IabCategories       string                           `gorm:"column:iab_categories" json:"iab_categories"`
	Uuid                string                           `gorm:"column:uuid" json:"uuid"`
	Requests            int64                            `gorm:"column:requests" json:"requests"`
	Impressions         int64                            `gorm:"column:impressions" json:"impressions"`
	Revenue             float64                          `gorm:"column:revenue" json:"revenue"`
	ApacSiteId          int64                            `gorm:"column:apac_siteid" json:"apac_siteid"`
	CachedAt            time.Time                        `gorm:"column:cached_at"  json:"cached_at"`
	RenderCache         int                              `gorm:"column:render_cache" json:"render_cache"`
	AdTag               []InventoryAdTagModel            `gorm:"foreignKey:InventoryID;references:ID"`
	Config              InventoryConfigModel             `gorm:"foreignKey:InventoryID;references:ID"`
	ConnectionDemand    []InventoryConnectionDemandModel `gorm:"foreignKey:InventoryID;references:ID"`
	RlsBidderSystem     []RlsBidderSystemInventoryModel  `gorm:"foreignKey:InventoryName;references:Name"` // Rls để check chang status trong BE
}

func (i *InventoryModel) IsFound() bool {
	if i.ID > 0 {
		return true
	}
	return false
}

func (i *InventoryModel) ToJSON() string {
	jsonEncode, _ := json.Marshal(i)
	return string(jsonEncode)
}

/**
Config TYPE
*/

type TYPEAdsTxtCustom string

func (t TYPEAdsTxtCustom) ToArray() (array []string) {
	array = strings.Split(string(t), "\n")
	return
}

type TYPEInventoryType int

const (
	InventoryTypeWeb = iota + 1
	InventoryTypeApp
)

func (t TYPEInventoryType) String() string {
	switch t {
	case InventoryTypeWeb:
		return "web"
	case InventoryTypeApp:
		return "app"
	default:
		return ""
	}
}

type TYPEInventorySyncAdsTxt int

const (
	InventorySyncAdsTxt = iota + 1
	InventorySyncAdsTxtNotIn
	InventorySyncAdsTxtError
)

func (t TYPEInventorySyncAdsTxt) String() string {
	switch t {
	case InventorySyncAdsTxt:
		return "In Sync"
	case InventorySyncAdsTxtNotIn:
		return "Not In Sync"
	case InventorySyncAdsTxtError:
		return "Does Not Exist"
	default:
		return ""
	}
}
