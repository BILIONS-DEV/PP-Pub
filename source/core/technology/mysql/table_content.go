package mysql

import (
	"gorm.io/gorm"
	"source/pkg/utility"
	"time"
)

type TableContent struct {
	Id            int64                 `gorm:"column:id" json:"id"`
	Uuid          string                `gorm:"column:uuid" json:"uuid"`
	Title         string                `gorm:"column:title" json:"title"`
	Status        TYPEStatus            `gorm:"column:status" json:"status"`
	Type          int64                 `gorm:"column:type" json:"type"`
	ContentDesc   string                `gorm:"column:content_desc" json:"content_desc"`
	Thumb         string                `gorm:"column:thumb" json:"thumb"`
	VideoUrl      string                `gorm:"column:video_url" json:"video_url"`
	VideoType     int                   `gorm:"column:video_type" json:"video_type"`
	Category      int64                 `gorm:"column:category" json:"category"`
	Channels      int64                 `gorm:"column:channels" json:"channels"`
	Tag           int64                 `gorm:"column:tag" json:"tag"`
	UserId        int64                 `gorm:"column:user_id" json:"user_id"`
	NameFile      string                `gorm:"column:name_file" json:"name_file"`
	Duration      int                   `gorm:"column:duration" json:"duration"`
	ConfigAdBreak TYPEConfigAdBreak     `gorm:"column:config_ad_break" json:"config_ad_break"`
	CreatedAt     time.Time             `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time             `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     gorm.DeletedAt        `gorm:"column:deleted_at" json:"deleted_at"`
	Keywords      []TableContentKeyword `gorm:"-" as:"-"`
	AdBreaks      []TableContentAdBreak `gorm:"-" as:"-"`
}

func (TableContent) TableName() string {
	return Tables.Content
}

func (rec *TableContent) GetRls() {
	var keywords []TableContentKeyword
	Client.Where("content_id = ?", rec.Id).Find(&keywords)
	rec.Keywords = keywords

	var adBreaks []TableContentAdBreak
	Client.Where("content_id = ?", rec.Id).Find(&adBreaks)
	rec.AdBreaks = adBreaks
}

func (rec *TableContent) RenderUuid() {
	if rec.Uuid == "" {
		uid := utility.RandSeq(9)
		recordCheck := TableContent{}
		Client.Where("uuid = ?", uid).Find(recordCheck)
		if recordCheck.Id != 0 {
			rec.RenderUuid()
			return
		} else {
			Client.Model(TableContent{}).Where("id = ?", rec.Id).Update("uuid", uid)
			rec.Uuid = uid
		}
	}
}

type TYPEConfigAdBreak int

const (
	TYPEConfigAdBreakManual TYPEConfigAdBreak = iota + 1
	TYPEConfigAdBreakAuto
)

func (s TYPEConfigAdBreak) String() string {
	switch s {
	case TYPEConfigAdBreakManual:
		return "Add Breaks Manual"
	case TYPEConfigAdBreakAuto:
		return "Auto"
	default:
		return ""
	}
}

func (s TYPEConfigAdBreak) Int() int {
	switch s {
	case TYPEConfigAdBreakManual:
		return 1
	case TYPEConfigAdBreakAuto:
		return 2
	default:
		return 0
	}
}
