package payload

import (
	"source/pkg/datatable"
)

type LanguageIndex struct {
	LanguageFilterPostData
	QuerySearch string `query:"f_q" json:"f_q" form:"f_q"`
	Category    []int  `query:"f_category" form:"f_category" json:"f_category"`
}

type LanguageFilterPayload struct {
	datatable.Request
	PostData *LanguageFilterPostData `query:"postData"`
}

type LanguageFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Page        int         `query:"page" json:"page" form:"page"`
	Limit       int         `query:"limit" json:"limit" form:"limit"`
	Start       int         `query:"start" json:"start" form:"start"`
	Length      int         `query:"length" json:"length" form:"length"`
}
