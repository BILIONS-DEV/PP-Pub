package payload

import (
	"source/core/technology/mysql"
	"source/pkg/datatable"
)

type InventoryAdTagSubmit struct {
	InventoryId int64                 `query:"inventory_id" form:"inventory_id" json:"inventory_id"`
	Name        string                `query:"name" form:"name" json:"name"`
	FloorPrice  float64               `query:"floor_price" form:"floor_price" json:"floor_price"`
	AdTagType   mysql.TYPEAdType      `query:"ad_tag_type" form:"ad_tag_type" json:"ad_tag_type"`
	AdTagSize   int                   `query:"ad_tag_size" form:"ad_tag_size" json:"ad_tag_size"`
	Gam         string                `query:"gam" form:"gam" json:"gam"`
	PassBack    string                `query:"pass_back" form:"pass_back" json:"pass_back"`
	Status      mysql.TypeStatusAdTag `query:"status" form:"status" json:"status"`
}

type InventoryAdTagIndex struct {
	InventoryAdTagFilterPostData
	OrderColumn int      `query:"order_column" json:"order_column" form:"order_column"`
	OrderDir    string   `query:"order_dir" json:"order_dir" form:"order_dir"`
	QuerySearch string   `query:"f_q" json:"f_q" form:"f_q"`
	Status      []string `query:"f_status" form:"f_status" json:"f_status"`
	Type        []int64  `query:"f_type" form:"f_type" json:"f_type"`
}

type InventoryAdTagFilterPayload struct {
	datatable.Request
	PostData *InventoryAdTagFilterPostData `query:"postData"`
}

type InventoryAdTagFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Status      interface{} `query:"f_status[]" json:"f_status" form:"f_status[]"`
	Type        interface{} `query:"f_type[]" json:"f_type" form:"f_type[]"`
	Page        int         `query:"page" json:"page" form:"page"`
	Limit       int         `query:"limit" json:"limit" form:"limit"`
	Start       int         `query:"start" json:"start" form:"start"`
	Length      int         `query:"length" json:"length" form:"length"`
	InventoryId int64       `query:"inventory_id" json:"inventory_id" form:"inventory_id"`
}
