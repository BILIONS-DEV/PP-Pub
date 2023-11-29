package mysql

type TableRlsUserPermission struct {
	Id           int64 `gorm:"column:id" json:"id"`
	UserId       int64 `gorm:"column:user_id" json:"user_id"`
	PermissionId int64 `gorm:"column:permission_id" json:"permission_id"`
}
