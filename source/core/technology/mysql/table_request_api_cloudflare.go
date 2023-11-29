package mysql

type TableRequestApiCloudFlare struct {
	Id     int64                   `gorm:"column:id" json:"id"`
	Uuid   string                  `gorm:"column:uuid" json:"uuid"`
	Status TYPEStatusApiCloudflare `gorm:"column:status" json:"status"`
}

func (TableRequestApiCloudFlare) TableName() string {
	return Tables.RequestApiCloudflare
}

type TYPEStatusApiCloudflare int

const (
	TYPEStatusApiCloudflareWaiting TYPEStatusApiCloudflare = iota + 1
	TYPEStatusApiCloudflarePending
)
