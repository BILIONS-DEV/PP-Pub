package payload

import "source/pkg/datatable"

type BlockingIndex struct {
	InventoryAdTagFilterPostData
	QuerySearch string `query:"f_q" json:"f_q" form:"f_q"`
}

type BlockingFilterPayload struct {
	datatable.Request
	PostData *BlockingFilterPostData `query:"postData"`
}

type BlockingFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Status      interface{} `query:"f_status[]" json:"f_status" form:"f_status[]"`
	Type        interface{} `query:"f_type[]" json:"f_type" form:"f_type[]"`
	Page        int         `query:"page" json:"page" form:"page"`
	Limit       int         `query:"limit" json:"limit" form:"limit"`
	Start       int         `query:"start" json:"start" form:"start"`
	Length      int         `query:"length" json:"length" form:"length"`
	InventoryId int64       `query:"inventory_id" json:"inventory_id" form:"inventory_id"`
}

type BlockingAdd struct {
	Id               int64    `form:"column:id" json:"id"`
	RestrictionName  string   `form:"column:name" json:"restriction_name"`
	Inventories      []int64  `form:"column:inventories" json:"inventories"`
	AdvertiseDomains []string `form:"column:advertise_domains" json:"advertise_domains"`
	CreativeIds      []string  `form:"column:creative_ids" json:"creative_ids"`
}
