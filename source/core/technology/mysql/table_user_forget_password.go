package mysql

import "time"

type TableUserForgetPassword struct {
	Id          int64     `gorm:"column:id" json:"id"`
	UserId      int64     `gorm:"column:user_id" json:"user_id"`
	Email       string    `gorm:"column:email" json:"email"`
	Uuid        string    `gorm:"column:uuid" json:"uuid"`
	IsUsed      int       `gorm:"column:is_used" json:"is_used"`
	CreatedDate time.Time `gorm:"column:created_date" json:"created_date"`
	ExpiredTime time.Time `gorm:"column:expired_time" json:"expired_time"`
}
