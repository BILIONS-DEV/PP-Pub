package mysql

import (
	"gorm.io/gorm"
	"time"
)

type TableTestAdTag struct {
	Id               int64          `gorm:"column:id" json:"id"`
	Name             string         `gorm:"column:name" json:"name"`
	UserId           int64          `gorm:"column:user_id" json:"user_id"`
	InventoryId      int64          `gorm:"column:inventory_id" json:"inventory_id"`
	Gam              string         `gorm:"column:gam" json:"gam"`
	Type             int            `gorm:"column:type" json:"type"`
	PassBack         string         `gorm:"column:pass_back" json:"pass_back"`
	Status           int            `gorm:"column:status" json:"status"`
	PrimaryAdSize    int            `gorm:"column:primary_ad_size" json:"primary_ad_size"`
	AdditionalAdSize string         `gorm:"column:additional_ad_size" json:"additional_ad_size"`
	TemplateId       int64          `gorm:"column:template_id" json:"template_id"`
	ContentSource    int            `gorm:"column:content_source" json:"content_source"`
	PlaylistId       int64          `gorm:"column:playlist_id" json:"playlist_id"`
	FeedUrl          string         `gorm:"column:feed_url" json:"feed_url"`
	Uuid             string         `gorm:"column:uuid" json:"uuid"`
	CreatedAt        time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
}
