package mysql

type TableChannelsKeyword struct {
	Id         int64  `gorm:"column:id" json:"id"`
	ChannelsId int64  `gorm:"column:channels_id" json:"channels_id"`
	Keyword    string `gorm:"column:keyword" json:"keyword"`
}

func (TableChannelsKeyword) TableName() string {
	return Tables.ChannelsKeyword
}