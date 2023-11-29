package mysql

import "time"

type TableDfpLineItem struct {
	Id           int64           `gorm:"column:id" json:"id"`
	Type         TypeDfpCreative `gorm:"column:type" json:"type"`
	PPLineItemId int64           `gorm:"column:pp_line_item_id" json:"pp_line_item_id"`
	UserId       int64           `gorm:"column:user_id" json:"user_id"`
	NetworkId    int64           `gorm:"column:network_id" json:"network_id"`
	NetworkName  string          `gorm:"column:network_name" json:"network_name"`
	OrderId      string          `gorm:"column:order_id" json:"order_id"`
	LineItemName string          `gorm:"column:line_item_name" json:"line_item_name"`
	LineItemId   string          `gorm:"column:line_item_id" json:"line_item_id"`
	Active       TypeOnOff       `gorm:"column:active" json:"active"`
	CreatedAt    time.Time       `gorm:"column:created_at" json:"created_at"`
}
