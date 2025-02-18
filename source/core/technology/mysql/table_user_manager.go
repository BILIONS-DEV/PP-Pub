package mysql

import (
	"time"
)

func (TableUserManager) TableName() string {
	return Tables.UserManager
}

type TableUserManager struct {
	Id          int64     `gorm:"column:id" json:"id"`
	PresenterId int64     `gorm:"column:presenter_id" json:"presenter_id"`
	FirstName   string    `gorm:"column:first_name" json:"first_name"`
	LastName    string    `gorm:"column:last_name" json:"last_name"`
	Email       string    `gorm:"column:email" json:"email"`
	Phone       string    `gorm:"column:phone" json:"phone"`
	Telegram    string    `gorm:"column:telegram" json:"telegram"`
	Skype       string    `gorm:"column:skype" json:"skype"`
	Linkedin    string    `gorm:"column:linkedin" json:"linkedin"`
	Avatar      string    `gorm:"column:avatar" json:"avatar"`
	Whatsapp    string    `gorm:"column:whatsapp" json:"whatsapp"`
	Wechat      string    `gorm:"column:wechat" json:"wechat"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
}

func (TableUserManagerPub) TableName() string {
	return Tables.UserManagerPub
}

type TableUserManagerPub struct {
	Id            int64     `gorm:"column:id" json:"id"`
	UserManagerId string    `gorm:"column:user_manager_id" json:"user_manager_id"`
	PresenterId   int64     `gorm:"column:presenter_id" json:"presenter_id"`
	Pubid        string    `gorm:"column:pub_id" json:"pub_id"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
}
