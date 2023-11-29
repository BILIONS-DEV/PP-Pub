package mysql

import "time"

type TableQizUsers struct {
	ID                int64     `gorm:"column:id"`
	UserLogin         string    `gorm:"column:user_login"`
	UserPass          string    `gorm:"column:user_pass"`
	UserNicename      string    `gorm:"column:user_nicename"`
	UserEmail         string    `gorm:"column:user_email"`
	UserUrl           string    `gorm:"column:user_url"`
	UserRegistered    time.Time `gorm:"column:user_registered"`
	UserActivationKey string    `gorm:"column:user_activation_key"`
	UserStatus        int64     `gorm:"column:user_status"`
	DisplayName       string    `gorm:"column:display_name"`
}

func (TableQizUsers) TableName() string {
	return Tables.QizUsers
}
