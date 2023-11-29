package payload

import (
	"source/core/technology/mysql"
	"source/pkg/datatable"
)

type FloorIndex struct {
	FloorFilterPostData
	QuerySearch string   `query:"f_q" json:"f_q" form:"f_q"`
	Status      []string `query:"f_status" form:"f_status" json:"f_status"`
	Type        []string `query:"f_type" form:"f_type" json:"f_type"`
	Domain      []string `query:"f_domain" form:"f_domain" json:"f_domain"`
	AdFormat    []string `query:"f_ad_format" form:"f_ad_format" json:"f_ad_format"`
	AdSize      []string `query:"f_ad_size" form:"f_ad_size" json:"f_ad_size"`
	AdTag       []string `query:"f_ad_tag" form:"f_ad_tag" json:"f_ad_tag"`
	Device      []string `query:"f_device" form:"f_device" json:"f_device"`
	Country     []string `query:"f_geo" form:"f_geo" json:"f_geo"`
}

type FloorFilterPayload struct {
	datatable.Request
	PostData *FloorFilterPostData `query:"postData"`
}

type FloorFilterPostData struct {
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

type FloorCreate struct {
	Id              int64           `json:"id" query:"id" form:"id"`
	Name            string          `json:"name" query:"name" form:"name"`
	Description     string          `json:"description" query:"description" form:"description"`
	Status          int             `json:"status" query:"status" form:"status"`
	FloorType       mysql.TYPEFloor `json:"floor_type" query:"floor_type" form:"floor_type"`
	FloorValue      float64         `json:"floor_value" query:"floor_value" form:"floor_value"`
	Priority        int             `json:"priority" query:"priority" form:"priority"`
	ListAdFormat    []ListTarget    `json:"listAdFormat" query:"listAdFormat" form:"listAdFormat"`
	ListAdSize      []ListTarget    `json:"listAdSize" query:"listAdSize" form:"listAdSize"`
	ListAdTag       []ListTarget    `json:"listAdtag" query:"listAdtag" form:"listAdtag"`
	ListAdInventory []ListTarget    `json:"listData" query:"listData" form:"listData"`
	ListGeo         []ListTarget    `json:"listGeo" query:"listGeo" form:"listGeo"`
	ListDevice      []ListTarget    `json:"listDevice" query:"listDevice" form:"listDevice"`
}

type FloorEdit struct {
	Id              int64        `json:"id" query:"id" form:"id"`
	Name            string       `json:"name" query:"name" form:"name"`
	Description     string       `json:"description" query:"description" form:"description"`
	Status          int          `json:"status" query:"status" form:"status"`
	FloorValue      float64      `json:"floor_value" query:"floor_value" form:"floor_value"`
	Priority        int          `json:"priority" query:"priority" form:"priority"`
	ListAdFormat    []ListTarget `json:"listAdFormat" query:"listAdFormat" form:"listAdFormat"`
	ListAdSize      []ListTarget `json:"listAdSize" query:"listAdSize" form:"listAdSize"`
	ListAdTag       []ListTarget `json:"listAdtag" query:"listAdtag" form:"listAdtag"`
	ListAdInventory []ListTarget `json:"listData" query:"listData" form:"listData"`
	ListGeo         []ListTarget `json:"listGeo" query:"listGeo" form:"listGeo"`
	ListDevice      []ListTarget `json:"listDevice" query:"listDevice" form:"listDevice"`
}
