package mysql2

import "gorm.io/gorm"

type TableLineItem struct {
	Targets []TableTarget `gorm:"foreignKey:line_item_id"`
	gorm.Model
	ID             int64                      `gorm:"column:id"`
	UserID         int64                      `gorm:"column:user_id"`
	Name           string                     `gorm:"column:name"`
	Description    string                     `gorm:"column:description"`
	ServerType     TYPELineItemServerType     `gorm:"column:server_type"`
	Status         TYPEStatus                 `gorm:"column:status"`
	ConnectionType TYPELineItemConnectionType `gorm:"column:connection_type"`
}

func (TableLineItem) TableName() string {
	return Tables.LineItem
}

// TYPELineItemServerType : Type of bidder
type TYPELineItemServerType int

const (
	TYPELineItemServerTypePrebid TYPELineItemServerType = iota + 1
	TYPELineItemServerTypeGoogle
)

func (t TYPELineItemServerType) String() string {
	switch t {
	case TYPELineItemServerTypePrebid:
		return "Prebid"
	case TYPELineItemServerTypeGoogle:
		return "Google"
	}
	return ""
}

// TYPELineItemConnectionType : Kiểu tạo dữ liệu trên GAM của pub
type TYPELineItemConnectionType int

const (
	TYPELineItemConnectionTypeAdUnits = iota + 1
	TYPELineItemConnectionTypeLineItem
	TYPELineItemConnectionTypeMcm
)

func (t TYPELineItemConnectionType) String() string {
	switch t {
	case TYPELineItemConnectionTypeAdUnits:
		return "AdUnits"
	case TYPELineItemConnectionTypeLineItem:
		return "LineItems"
	case TYPELineItemConnectionTypeMcm:
		return "MCM"
	}
	return ""
}
