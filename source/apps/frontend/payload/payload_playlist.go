package payload

import "source/pkg/datatable"

type PlaylistIndex struct {
	PlaylistFilterPostData
	//Status []string `query:"f_status[]" form:"f_status json:"f_status""`
	Status     []string `query:"f_status" form:"f_status" json:"f_status"`
	Permission []string `query:"f_permission" form:"f_permission" json:"f_permission"`
	SearchBy   []string `query:"f_search_by" form:"f_search_by" json:"f_search_by"`
}

type PlaylistFilterPayload struct {
	datatable.Request
	PostData *PlaylistFilterPostData `query:"postData"`
}

type PlaylistCreate struct {
	Id             int64  `query:"id" json:"id" form:"id"`
	Name           string `query:"name" json:"name" form:"name"`
	Description    string `query:"description" json:"description" form:"description"`
	OrderingMethod string `query:"ordering_method" json:"ordering_method" form:"ordering_method"`
	VideosLimit    int64 `query:"videos_limit" json:"videos_limit" form:"videos_limit"`
	// Channels    int64   `query:"channels" json:"channels" form:"channels"`
	// Category    int64   `query:"category" json:"category" form:"category"`
	// Category    int64   `query:"category" json:"category" form:"category"`
	// Content     []string `query:"content" json:"content" form:"content"`
	TypeCategory int64                `query:"type_category" json:"type_category" form:"type_category"`
	TypeChannels int64                `query:"type_channels" json:"type_channels" form:"type_channels"`
	TypeKeywords int64                `query:"type_keywords" json:"type_keywords" form:"type_keywords"`
	TypeLanguage int64                `query:"type_language" json:"type_language" form:"type_language"`
	TypeVideos   int64                `query:"type_videos" json:"type_videos" form:"type_videos"`
	ListChannels []ListPlaylistConfig `json:"listChannels" query:"listChannels" form:"listChannels"`
	ListLanguage []ListPlaylistConfig `json:"listLanguage" query:"listLanguage" form:"listLanguage"`
	ListCategory []ListPlaylistConfig `json:"listCategory" query:"listCategory" form:"listCategory"`
	ListKeywords []ListPlaylistConfig `json:"listKeywords" query:"listKeywords" form:"listKeywords"`
	ListVideos   []ListPlaylistConfig `json:"listVideos" query:"listVideos" form:"listVideos"`
}

type ListPlaylistConfig struct {
	Id   int64  `gorm:"column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
	Type string `gorm:"column:type" json:"type"`
}

type PlaylistFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Status      interface{} `query:"f_status[]" json:"f_status" form:"f_status[]"`
	Page        int         `query:"page" json:"page" form:"page"`
	Limit       int         `query:"limit" json:"limit" form:"limit"`
	Start       int         `query:"start" json:"start" form:"start"`
	Length      int         `query:"length" json:"length" form:"length"`
}
