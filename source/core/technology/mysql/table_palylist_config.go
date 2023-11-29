package mysql

type TablePlaylistConfig struct {
	Id               int64
	UserId           int64 `gorm:"column:user_id" json:"user_id"`
	PlaylistId       int64 `gorm:"column:playlist_id" json:"playlist_id"`
	Type             int64 `gorm:"column:type" json:"type"`
	LanguageId       int64 `gorm:"column:language_id" json:"language_id"`
	ChannelsId       int64 `gorm:"column:channels_id" json:"channels_id"`
	ContentId        int64 `gorm:"column:content_id" json:"content_id"`
	ContentKeywordId int64 `gorm:"column:content_keyword_id" json:"content_keyword_id"`
	CategoryId       int64 `gorm:"column:category_id" json:"category_id"`
}

func (TablePlaylistConfig) TableName() string {
	return Tables.PlaylistConfig
}
