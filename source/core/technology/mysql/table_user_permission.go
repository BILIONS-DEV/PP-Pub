package mysql

type TableUserPermission struct {
	Id          int64          `gorm:"column:id" json:"id"`
	Name        string         `gorm:"column:name" json:"name"`
}