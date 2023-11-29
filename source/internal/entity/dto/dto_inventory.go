package dto

import (
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
)

type PayloadInventoryIndex struct {
	datatable.Request
	OrderColumn int      `query:"order_column" json:"order_column" form:"order_column"`
	OrderDir    string   `query:"order_dir" json:"order_dir" form:"order_dir"`
	QuerySearch string   `query:"f_q" json:"f_q" form:"f_q"`
	Status      []string `query:"f_status" form:"f_status" json:"f_status"`
	Type        []string `query:"f_type" form:"f_type" json:"f_type"`
	AdsTxtSync  []string `query:"f_ads_sync" form:"f_ads_sync" json:"f_ads_sync"`
}

type PayloadInventoryIndexPost struct {
	datatable.Request
	UserID   int64                          `json:"user_id"`
	PostData *PayloadInventoryIndexPostData `query:"postData"`
}

type PayloadInventoryIndexPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Status      interface{} `query:"f_status[]" json:"f_status" form:"f_status[]"`
	Type        interface{} `query:"f_type[]" json:"f_type" form:"f_type[]"`
	AdsSync     interface{} `query:"f_ads_sync[]" json:"f_ads_sync" form:"f_ads_sync[]"`
}

type ResponseInventoryDatatable struct {
	*model.InventoryModel
	RowId        string `json:"DT_RowId"`
	Name         string `json:"name"`
	NameForAdmin string `json:"name_for_admin"`
	Status       string `json:"status"`
	Live         string `json:"live"`
	Type         string `json:"type"`
	SyncAdsTxt   string `json:"sync_ads_txt"`
	Action       string `json:"action"`
}

func (t *PayloadInventoryIndexPost) Validate() (errs []Error) {
	if t.UserID == 0 {
		errs = append(errs, Error{Message: "permission denied"})
	}
	return
}
