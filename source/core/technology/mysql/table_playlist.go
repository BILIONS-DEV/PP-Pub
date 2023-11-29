package mysql

import (
	"gorm.io/gorm"
	"time"
)

type TablePlaylist struct {
	Id             int64     `gorm:"column:id" json:"id"`
	Name           string    `gorm:"column:name" json:"name"`
	UserId         int64     `gorm:"column:user_id" json:"user_id"`
	Description    string    `gorm:"column:description" json:"description"`
	OrderingMethod string    `gorm:"column:ordering_method" json:"ordering_method"`
	VideosLimit    int64     `gorm:"column:videos_limit" json:"videos_limit"`
	IsDefault      TypeOnOff `gorm:"column:is_default" json:"is_default"`
	// Channels    int64 `gorm:"column:channels" json:"channels"`
	// Category    int64 `gorm:"column:category" json:"category"`
	// Content     string         `gorm:"column:content" json:"content"`
	CreatedAt         time.Time             `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time             `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt         gorm.DeletedAt        `gorm:"column:deleted_at" json:"deleted_at"`
	ChannelsAndVideos []TablePlaylistConfig `gorm:"-"`
}

func (TablePlaylist) TableName() string {
	return Tables.Playlist
}

func (rec *TablePlaylist) GetRls() {
	var channelsAndVideos []TablePlaylistConfig
	Client.Where(TablePlaylistConfig{PlaylistId: rec.Id}).Find(&channelsAndVideos)
	rec.ChannelsAndVideos = channelsAndVideos
}
