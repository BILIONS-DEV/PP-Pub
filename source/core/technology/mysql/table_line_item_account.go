package mysql

type TableLineItemAccount struct {
	Id         int64 `gorm:"column:id" json:"id"`
	LineItemId int64 `gorm:"column:line_item_id" json:"line_item_id"`
	BidderId   int64 `gorm:"column:bidder_id" json:"bidder_id"`
}

func (TableLineItemAccount) TableName() string {
	return Tables.LineItemAccount
}