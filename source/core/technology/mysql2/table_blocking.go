package mysql2

import "gorm.io/gorm"

type TableBlocking struct {
	gorm.Model
	ID           int64                      `gorm:"column:id"`
	UserID       int64                      `gorm:"column:user_id"`
	Name         string                     `gorm:"column:restriction_name"`
	Restrictions []TableBlockingRestriction `gorm:"foreignKey:blocking_id"`
}

func (TableBlocking) TableName() string {
	return Tables.Blocking
}
