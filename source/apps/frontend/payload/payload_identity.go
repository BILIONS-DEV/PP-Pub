package payload

import (
	"source/core/technology/mysql"
	"source/pkg/datatable"
)

type IdentityIndex struct {
	IdentityFilterPostData
	OrderColumn int      `query:"order_column" json:"order_column" form:"order_column"`
	OrderDir    string   `query:"order_dir" json:"order_dir" form:"order_dir"`
	QuerySearch string   `query:"f_q" json:"f_q" form:"f_q"`
	Status      []string `query:"f_status" form:"f_status" json:"f_status"`
	Type        []string `query:"f_type" form:"f_type" json:"f_type"`
	Domain      []string `query:"f_domain" form:"f_domain" json:"f_domain"`
}

type IdentityFilterPayload struct {
	datatable.Request
	PostData *IdentityFilterPostData `query:"postData"`
}

type IdentityFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Status      interface{} `query:"f_status[]" json:"f_status" form:"f_status[]"`
	Type        interface{} `query:"f_type[]" json:"f_type" form:"f_type[]"`
	Domain      interface{} `query:"f_domain[]" json:"f_domain" form:"f_domain[]"`
	AdFormat    interface{} `query:"f_ad_format[]" json:"f_ad_format" form:"f_ad_format[]"`
	AdSize      interface{} `query:"f_ad_size[]" json:"f_ad_size" form:"f_ad_size[]"`
	AdTag       interface{} `query:"f_ad_tag[]" json:"f_ad_tag" form:"f_ad_tag[]"`
	Device      interface{} `query:"f_device[]" json:"f_device" form:"f_device[]"`
	Country     interface{} `query:"f_geo[]" json:"f_geo" form:"f_geo[]"`
	Page        int         `query:"page" json:"page" form:"page"`
	Limit       int         `query:"limit" json:"limit" form:"limit"`
	Start       int         `query:"start" json:"start" form:"start"`
	Length      int         `query:"length" json:"length" form:"length"`
}

type IdentityCreate struct {
	Id              int64           `json:"id" query:"id" form:"id"`
	Name            string          `json:"name" query:"name" form:"name"`
	Description     string          `json:"description" query:"description" form:"description"`
	AuctionDelay    int             `form:"column:auction_delay" json:"auction_delay"`
	SyncDelay       int             `form:"column:sync_delay" json:"sync_delay"`
	Status          mysql.TypeOnOff `json:"status" query:"status" form:"status"`
	Priority        int             `json:"priority" query:"priority" form:"priority"`
	ModuleParams    []ModuleInfo    `json:"module_params" form:"module_params"`
	ListAdInventory []ListTarget    `json:"listInventory" query:"listInventory" form:"listInventory"`
	IsDefault       mysql.TypeOnOff `json:"-"`
}

type IdentityEdit struct {
	Id              int64           `json:"id" query:"id" form:"id"`
	Name            string          `json:"name" query:"name" form:"name"`
	Description     string          `json:"description" query:"description" form:"description"`
	Status          mysql.TypeOnOff `json:"status" query:"status" form:"status"`
	Priority        int             `json:"priority" query:"priority" form:"priority"`
	ModuleParams    []ModuleInfo    `json:"module_params" form:"module_params"`
	ListAdInventory []ListTarget    `json:"listInventory" query:"listInventory" form:"listInventory"`
}
