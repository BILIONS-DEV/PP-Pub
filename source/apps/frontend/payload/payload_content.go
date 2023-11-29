package payload

import (
	"source/core/technology/mysql"
	"source/pkg/datatable"
)

type ContentIndex struct {
	ContentFilterPostData
	QuerySearch string `query:"f_q" json:"f_q" form:"f_q"`
	// Category    []int  `query:"f_category" form:"f_category" json:"f_category"`
	Type []int `query:"f_type" form:"f_type" json:"f_type"`
}

type ContentAdd struct {
	ContentFilterPostData
	Channels int64 `query:"f_channels" json:"f_channels" form:"f_channels"`
}

type ContentFilterPayload struct {
	datatable.Request
	PostData *ContentFilterPostData `query:"postData"`
}

type ContentCreate struct {
	Id            int64                   `query:"id" json:"id" form:"id"`
	Title         string                  `query:"title" json:"title" form:"title"`
	ContentDesc   string                  `query:"content_description" json:"content_description" form:"content_description"`
	VideoUrl      string                  `query:"video_url" json:"video_url" form:"video_url"`
	Duration      int                     `query:"duration" json:"duration" form:"duration" `
	Thumb         string                  `query:"thumb" json:"thumb" form:"thumb"`
	VideoType     int                     `query:"video_type" json:"video_type" form:"video_type"`
	Category      int64                   `query:"category" json:"category" form:"category"`
	Channels      int64                   `query:"channels" json:"channels" form:"channels"`
	Keyword       []string                `query:"keyword" json:"keyword" form:"keyword"`
	Tag           []string                `query:"tag" json:"tag" form:"tag"`
	NameFile      string                  `query:"name_file" json:"name_file" form:"name_file"`
	AdStartTime   []string                `query:"ad_start_time" json:"ad_start_time" form:"ad_start_time"`
	ConfigAdBreak mysql.TYPEConfigAdBreak `query:"config_ad_break" json:"config_ad_break" form:"config_ad_break"`
	AdBreaks      []AdBreak               `query:"ad_breaks" json:"ad_breaks" form:"ad_breaks"`
}

type AdBreak struct {
	Type      string              `query:"type" json:"type" form:"type"`
	TimeBreak string              `query:"time_break" json:"time_break" form:"time_break"`
	BreakMode mysql.TYPEBreakMode `query:"break_mode" json:"break_mode" form:"break_mode"`
}

type ContentQuizCreate struct {
	Id        int64       `query:"id" json:"id" form:"id"`
	Title     string      `query:"title" json:"title" form:"title"`
	Category  int64       `query:"category" json:"category" form:"category"`
	Tag       []string    `query:"tag" json:"tag" form:"tag"`
	Questions []Questions `query:"questions" json:"questions" json:"questions"`
}

type Questions struct {
	Id             int64    `query:"id" json:"id" form:"id"`
	Title          string   `query:"title" json:"title" form:"title"`
	BackgroundType int64    `query:"background_type" json:"background_type" form:"background_type"`
	Background     string   `query:"background" json:"background" form:"background"`
	Type           int64    `query:"type" json:"type" form:"type"`
	Answer         []string `query:"answer" json:"answer" form:"answer"`
	PictureType    int64    `query:"picture_type" json:"picture_type" form:"picture_type"`
	Picture        string   `query:"picture" json:"picture" form:"picture"`
}

type ContentFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Category    interface{} `query:"f_category[]" json:"f_category" form:"f_category[]"`
	VideoType   interface{} `query:"f_video[]" json:"f_video" form:"f_video[]"`
	Type        interface{} `query:"f_type[]" json:"f_type" form:"f_type[]"`
	Page        int         `query:"page" json:"page" form:"page"`
	Limit       int         `query:"limit" json:"limit" form:"limit"`
	Start       int         `query:"start" json:"start" form:"start"`
	Length      int         `query:"length" json:"length" form:"length"`
}
