package model

func (AccountModel) TableName() string {
	return "account"
}

type AccountModel struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Object        string `gorm:"column:object" json:"object"`
	Name          string `gorm:"column:name" json:"name"`
	ShowName      string `gorm:"column:show_name" json:"show_name"`
	KeyCpcTaboola string `gorm:"column:key_cpc_taboola" json:"key_cpc_taboola"`
	RefreshToken  string `gorm:"column:refresh_token" json:"refresh_token"`
	Note          string `gorm:"column:note" json:"note"`
}

func (a *AccountModel) IsFound() bool {
	if a.ID > 0 {
		return true
	}
	return false
}
