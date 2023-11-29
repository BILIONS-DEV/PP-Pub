package mysql

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type TableAbTesting struct {
	Id                int64             `gorm:"column:id" json:"id"`
	UserId            int64             `gorm:"column:user_id" json:"user_id"`
	Name              string            `gorm:"column:name" json:"name"`
	Description       string            `gorm:"column:description" json:"description"`
	TestType          TYPETestType      `gorm:"column:test_type" json:"test_type"`
	BidderId          int64             `gorm:"column:bidder_id" json:"bidder_id"`
	UserIdModuleId    int64             `gorm:"column:module_userid_id" json:"module_userid_id"`
	TestGroupSize     TYPETestGroupSize `gorm:"column:test_group_size" json:"test_group_size"`
	DynamicFloorPrice int64             `gorm:"column:dynamic_floor_price" json:"dynamic_floor_price"`
	HardPriceFloor    float64           `gorm:"column:hard_price_floor" json:"hard_price_floor"`
	Status            TypeOnOff         `gorm:"column:status" json:"status"`
	StartDate         sql.NullTime      `gorm:"column:start_date" json:"start_date"`
	EndDate           sql.NullTime      `gorm:"column:end_date" json:"end_date"`
	Priority          int               `gorm:"column:priority" json:"priority"`
	AbTestFloor       int               `gorm:"column:abtest_floor" json:"abtest_floor"`
	CreatedAt         time.Time         `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time         `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt         gorm.DeletedAt    `gorm:"column:deleted_at" json:"deleted_at"`
}

func (TableAbTesting) TableName() string {
	return Tables.AbTesting
}

type TYPETestType int

const (
	TYPETestTypeAuctionTimeOut TYPETestType = iota + 1
	TYPETestTypeClientVsServer
	TYPETestTypeUserIdModule
	TYPETestTypeDynamicHardPriceFloor
)

func (t TYPETestType) String() string {
	switch t {
	case TYPETestTypeAuctionTimeOut:
		return "Auction Timeout"
	case TYPETestTypeClientVsServer:
		return "Client vs Server"
	case TYPETestTypeUserIdModule:
		return "User Id Module"
	case TYPETestTypeDynamicHardPriceFloor:
		return "Dynamic vs Hard Price Floor"
	default:
		return ""
	}
}

func (t TYPETestType) CheckValid() bool {
	switch t {
	case
		TYPETestTypeAuctionTimeOut,
		TYPETestTypeClientVsServer,
		TYPETestTypeUserIdModule,
		TYPETestTypeDynamicHardPriceFloor:
		return true
	default:
		return false
	}
}

type TYPETestGroupSize int

const (
	TYPETestGroupSize1 TYPETestGroupSize = iota + 1
	TYPETestGroupSize5
	TYPETestGroupSize10
	TYPETestGroupSize25
	TYPETestGroupSize50
)

func (t TYPETestGroupSize) CheckValid() bool {
	switch t {
	case
		TYPETestGroupSize1,
		TYPETestGroupSize5,
		TYPETestGroupSize10,
		TYPETestGroupSize25,
		TYPETestGroupSize50:
		return true
	default:
		return false
	}
}

func (t TYPETestGroupSize) Int() int {
	switch t {
	case TYPETestGroupSize1:
		return 1
	case TYPETestGroupSize5:
		return 2
	case TYPETestGroupSize10:
		return 3
	case TYPETestGroupSize25:
		return 4
	case TYPETestGroupSize50:
		return 5
	default:
		return 1
	}
}

func (t TYPETestGroupSize) Value() int {
	switch t {
	case TYPETestGroupSize1:
		return 1
	case TYPETestGroupSize5:
		return 5
	case TYPETestGroupSize10:
		return 10
	case TYPETestGroupSize25:
		return 25
	case TYPETestGroupSize50:
		return 50
	default:
		return 1
	}
}
