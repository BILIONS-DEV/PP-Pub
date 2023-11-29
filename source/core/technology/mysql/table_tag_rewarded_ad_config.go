package mysql

type TableTagRewardedAdConfig struct {
	TagID          int64  `gorm:"column:tag_id;primaryKey" json:"tag_id"`
	Title          string `gorm:"column:title" json:"title"`
	BtnApproved    string `gorm:"column:btn_approved" json:"btn_approved"`
	BtnCancel      string `gorm:"column:btn_cancel" json:"btn_cancel"`
	Amount         string `gorm:"column:amount" json:"amount"`
	Type           string `gorm:"column:type" json:"type"`
	MessageSuccess string `gorm:"column:message_success" json:"message_success"`
	BtnClose       string `gorm:"column:btn_close" json:"btn_close"`
}

func (TableTagRewardedAdConfig) TableName() string {
	return Tables.TagRewardedAdConfig
}
