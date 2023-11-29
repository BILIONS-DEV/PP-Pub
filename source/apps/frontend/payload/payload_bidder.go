package payload

import "source/pkg/datatable"

type BidderIndex struct {
	InventoryFilterPostData
	//Status []string `query:"f_status[]" form:"f_status json:"f_status""`
	Status     []string `query:"f_status" form:"f_status" json:"f_status"`
	Permission []string `query:"f_permission" form:"f_permission" json:"f_permission"`
	SearchBy   []string `query:"f_search_by" form:"f_search_by" json:"f_search_by"`
}

type BidderFilterPayload struct {
	datatable.Request
	PostData *InventoryFilterPostData `query:"postData"`
}

type BidderFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Status      interface{} `query:"f_status[]" json:"f_status" form:"f_status[]"`
	Page        int         `query:"page" json:"page" form:"page"`
	Limit       int         `query:"limit" json:"limit" form:"limit"`
	Start       int         `query:"start" json:"start" form:"start"`
	Length      int         `query:"length" json:"length" form:"length"`
}
