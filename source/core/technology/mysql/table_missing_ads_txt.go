package mysql

func (TableMissingAdsTxt) TableName() string {
	return Tables.MissingAdsTxt
}

type TableMissingAdsTxt struct {
	Line         string         `gorm:"column:line" json:"line"`
	UserId       int64          `gorm:"column:user_id" json:"user_id"`
	InventoryId  int64          `gorm:"column:inventory_id" json:"inventory_id"`
	Domain       string         `gorm:"column:domain" json:"domain"`
	AdsTxtUrl    string         `gorm:"column:ads_txt_url" json:"ads_txt_url"`
	ErrorMessage string         `gorm:"column:error_message" json:"error_message"`
}
