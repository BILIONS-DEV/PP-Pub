package dto

import (
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
)

type PayloadInventorySetup struct {
	datatable.Request
	Id          int64    `query:"id" json:"id" form:"id"`
	Tab         int64    `query:"tab" json:"tab" form:"tab"`
	QuerySearch string   `query:"f_q" json:"f_q" form:"f_q"`
	Status      []string `query:"f_status" form:"f_status" json:"f_status"`
	Type        []string `query:"f_type" form:"f_type" json:"f_type"`
}

type PayloadInventoryAdTagIndexPost struct {
	datatable.Request
	UserID   int64                               `json:"user_id"`
	PostData *PayloadInventoryAdTagIndexPostData `query:"postData"`
}

type PayloadInventoryAdTagIndexPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Status      interface{} `query:"f_status[]" json:"f_status" form:"f_status[]"`
	Type        interface{} `query:"f_type[]" json:"f_type" form:"f_type[]"`
}

type ResponseInventoryAdTagDatatable struct {
	*model.InventoryAdTagModel
	RowId  string `json:"DT_RowId"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Type   string `json:"type"`
	Size   string `json:"size"`
	Action string `json:"action"`
}

func (t *PayloadInventoryAdTagIndexPost) Validate() (errs []Error) {
	if t.UserID == 0 {
		errs = append(errs, Error{Message: "permission denied"})
	}
	return
}
