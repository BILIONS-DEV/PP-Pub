package payload

import "source/pkg/datatable"

type InventoryConnectionDemandIndex struct {
	InventoryConnectionDemandFilterPostData
	QuerySearch string   `query:"f_q" json:"f_q" form:"f_q"`
	Status      []string `query:"f_status" form:"f_status" json:"f_status"`
}

type InventoryConnectionDemandFilterPayload struct {
	datatable.Request
	PostData *InventoryConnectionDemandFilterPostData `query:"postData"`
}

type InventoryConnectionDemandFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Status      interface{} `query:"f_status[]" json:"f_status" form:"f_status[]"`
	Type        interface{} `query:"f_type[]" json:"f_type" form:"f_type[]"`
	Page        int         `query:"page" json:"page" form:"page"`
	Limit       int         `query:"limit" json:"limit" form:"limit"`
	Start       int         `query:"start" json:"start" form:"start"`
	Length      int         `query:"length" json:"length" form:"length"`
	InventoryId int64       `query:"inventory_id" json:"inventory_id" form:"inventory_id"`
}
