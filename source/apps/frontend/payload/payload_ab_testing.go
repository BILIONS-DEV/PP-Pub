package payload

import (
	"source/core/technology/mysql"
	"source/pkg/datatable"
)

type AbTestingIndex struct {
	AbTestingFilterPostData
	OrderColumn int      `query:"order_column" json:"order_column" form:"order_column"`
	OrderDir    string   `query:"order_dir" json:"order_dir" form:"order_dir"`
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

type AbTestingFilterPayload struct {
	datatable.Request
	PostData *AbTestingFilterPostData `query:"postData"`
}

type AbTestingFilterPostData struct {
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

type AbTestingCreate struct {
	Id                int64                   `json:"id" query:"id" form:"id"`
	Name              string                  `json:"name" query:"name" form:"name"`
	Description       string                  `json:"description" query:"description" form:"description"`
	TestType          mysql.TYPETestType      `json:"test_type" query:"test_type" form:"test_type"`
	Bidder            int64                   `json:"select_bidder" query:"select_bidder" form:"select_bidder"`
	UserIdModule      int64                   `json:"select_user_id_module" query:"select_user_id_module" form:"select_user_id_module"`
	TestGroupSize     mysql.TYPETestGroupSize `json:"test_group_size" query:"test_group_size" form:"test_group_size"`
	DynamicFloorPrice int64                   `json:"dynamic_floor_price" query:"dynamic_floor_price" form:"dynamic_floor_price"`
	HardPriceFloor    float64                 `json:"hard_price_floor" query:"hard_price_floor" form:"hard_price_floor"`
	Status            mysql.TypeOnOff         `json:"status" query:"status" form:"status"`
	StartDate         string                  `json:"start_date" query:"start_date" form:"start_date"`
	EndDate           string                  `json:"end_date" query:"end_date" form:"end_date"`
	ListAdFormat      []ListTarget            `json:"listAdFormat" query:"listAdFormat" form:"listAdFormat"`
	ListAdSize        []ListTarget            `json:"listAdSize" query:"listAdSize" form:"listAdSize"`
	ListAdTag         []ListTarget            `json:"listAdtag" query:"listAdtag" form:"listAdtag"`
	ListAdInventory   []ListTarget            `json:"listData" query:"listData" form:"listData"`
	ListGeo           []ListTarget            `json:"listGeo" query:"listGeo" form:"listGeo"`
	ListDevice        []ListTarget            `json:"listDevice" query:"listDevice" form:"listDevice"`
}
