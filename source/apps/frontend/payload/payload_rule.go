package payload

import (
	"source/core/technology/mysql"
	"source/pkg/datatable"
	"time"
)

type Rule struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type RuleIndex struct {
	RuleFilterPostData
	OrderColumn int      `query:"order_column" json:"order_column" form:"order_column"`
	OrderDir    string   `query:"order_dir" json:"order_dir" form:"order_dir"`
	Status      []string `query:"f_status" form:"f_status" json:"f_status"`
	Type        []string `query:"f_type" form:"f_type" json:"f_type"`
}

type RuleFilterPayload struct {
	datatable.Request
	PostData *RuleFilterPostData `query:"postData"`
}

type RuleFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Status      interface{} `query:"f_status[]" json:"f_status" form:"f_status[]"`
	Type        interface{} `query:"f_type[]" json:"f_type" form:"f_type[]"`
	Page        int         `query:"page" json:"page" form:"page"`
	Limit       int         `query:"limit" json:"limit" form:"limit"`
	Start       int         `query:"start" json:"start" form:"start"`
	Length      int         `query:"length" json:"length" form:"length"`
}

type LoadMoreData struct {
	Data        string  `json:"data"`
	ListFilter  []int64 `json:"list_filter"`
	IsMoreData  bool    `json:"is_more_data"`
	IsSearch    bool    `json:"is_search"`
	LastPage    bool    `json:"last_page"`
	CurrentPage int     `json:"current_page"`
	Total       int     `json:"total"`
	TotalAll    int64   `json:"total_all"`
}

type SelectAll struct {
	Data         string       `json:"data"`
	ListFilter   []int64      `json:"list_filter"`
	ListSelected []ListTarget `json:"list_selected"`
}

type ListTargetCheck struct {
	Id       int64  `gorm:"column:id" json:"id"`
	Name     string `gorm:"column:name" json:"name"`
	Selected bool   `gorm:"column:selected" json:"selected"`
}

type RuleSubmit struct {
	// Main
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Status      mysql.TYPEOnOff    `json:"status"`
	Type        mysql.TYPERuleType `json:"type"`

	// Blocked Page
	Domain int64    `json:"domain"`
	Pages  []string `json:"pages"`
}
