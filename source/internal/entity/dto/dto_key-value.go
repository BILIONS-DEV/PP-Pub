package dto

import (
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
	"source/internal/errors"
)

type PayloadKeyValueSubmit struct {
	ID     int64          `json:"id"`
	Key    string         `json:"key"`
	Values []PayloadValue `json:"values"`
}

type PayloadValue struct {
	ID    int64  `json:"id"`
	Value string `json:"value"`
}

func (p *PayloadKeyValueSubmit) Validate() (errs []error) {
	if p.Key == "" {
		errs = append(errs, errors.NewWithID(`"Key" is required`, "key"))
	}

	if len(p.Values) == 0 {
		errs = append(errs, errors.NewWithID(`"Value" is required`, "value"))
	}
	return
}

func (p *PayloadKeyValueSubmit) ToModel() model.KeyModel {
	record := model.KeyModel{}
	record.ID = p.ID
	record.KeyName = p.Key
	for _, value := range p.Values {
		if value.Value == "" {
			continue
		}
		record.Value = append(record.Value, model.ValueModel{
			ID:    value.ID,
			KeyID: record.ID,
			Value: value.Value,
		})
	}
	return record
}

type PayloadKeyValueIndex struct {
	datatable.Request
	QuerySearch string `query:"f_q" json:"f_q" form:"f_q"`
}

type PayloadKeyValueIndexPost struct {
	datatable.Request
	UserID   int64                         `json:"user_id"`
	PostData *PayloadKeyValueIndexPostData `query:"postData"`
}

type PayloadKeyValueIndexPostData struct {
	QuerySearch string `query:"f_q" json:"f_q" form:"f_q"`
}

func (t *PayloadKeyValueIndexPost) Validate() (errs []Error) {
	if t.UserID == 0 {
		errs = append(errs, Error{Message: "permission denied"})
	}
	return
}

type ResponseKeyValueDatatable struct {
	ID      int64  `json:"id"`
	KeyName string `json:"key_name"`
	RowId   string `json:"DT_RowId"`
	Action  string `json:"action"`
}

type PayloadKeyValueEdit struct {
	ID int64 `json:"id"`
}

func (t *PayloadKeyValueEdit) Validate() (errs []error) {
	if t.ID == 0 {
		errs = append(errs, errors.NewWithID(`"id" is missing`, "id"))
	}
	return
}
