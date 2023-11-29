package payload

import "source/pkg/datatable"

type LoadHistory struct {
	Id     int64  `json:"id" form:"id"`
	Object string `json:"object" form:"object"`
}

type HistoryIndex struct {
	HistoryFilterPostData
	QuerySearch string `query:"search" json:"search" form:"search"`
	ObjectPage  string `query:"object_page" form:"object_page" json:"object_page"`
	ObjectId    string `query:"object_id" form:"object_id" json:"object_id"`
}

type HistoryFilterPostData struct {
	QuerySearch string `query:"search" json:"search" form:"search"`
	ObjectPage  string `query:"object_page[]" json:"object_page" form:"object_page[]"`
	ObjectId    string `query:"object_id[]" json:"object_id" form:"object_id[]"`
	Page        int    `query:"page" json:"page" form:"page"`
	Limit       int    `query:"limit" json:"limit" form:"limit"`
	Start       int    `query:"start" json:"start" form:"start"`
	Length      int    `query:"length" json:"length" form:"length"`
}

type HistoryFilterPayload struct {
	datatable.Request
	PostData *HistoryFilterPostData `query:"postData"`
}

type LoadObjectByPage struct {
	ObjectPage string `json:"object_page" form:"object_page"`
}

type ObjectPage struct {
	ID   int64
	Name string
}