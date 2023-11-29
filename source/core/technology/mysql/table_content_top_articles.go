package mysql

import (
	"gorm.io/gorm"
)

func (TableContentTopArticles) TableName() string {
	return Tables.ContentTopArticles
}

type TableContentTopArticles struct {
	Id              int64          `gorm:"column:id" json:"id"`
	Domain          string         `gorm:"column:domain" json:"domain"`
	Type            string         `gorm:"column:type" json:"type"`
	FeedUrl         string         `gorm:"column:feed_url" json:"feed_url"`
	VideoUrlMp4     string         `gorm:"column:video_url_mp4" json:"video_url_mp4"`
	VideoUrlM3u8    string         `gorm:"column:video_url_m3u8" json:"video_url_m3u8"`
	VideoUrlYoutube string         `gorm:"column:video_url_youtube" json:"video_url_youtube"`
	YoutubeId       string         `gorm:"column:youtube_id" json:"youtube_id"`
	Info            string         `gorm:"column:info" json:"info"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}
