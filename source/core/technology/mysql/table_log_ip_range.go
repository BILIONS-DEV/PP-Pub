package mysql

type TableLogIpRange struct {
	Id      int64  `gorm:"id" json:"id"`
	Type    string `gorm:"type" json:"type"`
	IpRange string `gorm:"ip_range" json:"ip_range"`
}

func (TableLogIpRange) TableName() string {
	return Tables.LogIpRange
}
