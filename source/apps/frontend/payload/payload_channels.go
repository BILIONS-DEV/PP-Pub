package payload

import (
	"source/pkg/datatable"
)

type ChannelsIndex struct {
	ChannelsFilterPostData
	QuerySearch string `query:"f_q" json:"f_q" form:"f_q"`
	Category    []int  `query:"f_category" form:"f_category" json:"f_category"`
}

type ChannelsFilterPayload struct {
	datatable.Request
	PostData *ChannelsFilterPostData `query:"postData"`
}

type ChannelsCreate struct {
	Id          int64    `query:"id" json:"id" form:"id"`
	Name        string   `query:"name" json:"name" form:"name"`
	Description string   `query:"description" json:"description" form:"description"`
	Category    int64    `query:"category" json:"category" form:"category"`
	Keyword     []string `query:"keyword" json:"keyword" form:"keyword"`
	Language    int64   `query:"language" json:"language" form:"language"`
}

type ChannelsFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Category    interface{} `query:"f_category[]" json:"f_category" form:"f_category[]"`
	Page        int         `query:"page" json:"page" form:"page"`
	Limit       int         `query:"limit" json:"limit" form:"limit"`
	Start       int         `query:"start" json:"start" form:"start"`
	Length      int         `query:"length" json:"length" form:"length"`
}
