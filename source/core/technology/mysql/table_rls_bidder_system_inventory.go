package mysql

import "time"

type TableRlsBidderSystemInventory struct {
	Id            int64                              `gorm:"column:id" json:"id"`
	BidderId      int64                              `gorm:"column:bidder_id" json:"bidder_id"`
	InventoryName string                             `gorm:"column:inventory_name" json:"inventory_name"`
	Status        TYPERlsBidderSystemInventoryStatus `gorm:"column:status" json:"status"`
	Reason        string                             `gorm:"column:reason" json:"reason"`
	UpdatedAt     time.Time                          `gorm:"column:updated_at" json:"updated_at"`
}

func (TableRlsBidderSystemInventory) TableName() string {
	return Tables.RlsBidderSystemInventory
}

type TYPERlsBidderSystemInventoryStatus int

const (
	RlsBidderSystemInventoryPending TYPERlsBidderSystemInventoryStatus = iota + 1
	RlsBidderSystemInventorySubmitted
	RlsBidderSystemInventoryApproved
	RlsBidderSystemInventoryRejected
	RlsBidderSystemInventoryQueue
	RlsBidderSystemInventoryNotfound
	RlsBidderSystemInventoryApprovedS2S
	RlsBidderSystemInventoryApprovedClient
	RlsBidderSystemInventoryNotUse
)

func (p TYPERlsBidderSystemInventoryStatus) ColorClass() string {
	switch p {
	case RlsBidderSystemInventoryPending:
		return "custom-bd-text-dark"

	case RlsBidderSystemInventorySubmitted:
		return "custom-bd-text-warning"

	case RlsBidderSystemInventoryApproved:
		return "custom-bd-text-success"

	case RlsBidderSystemInventoryRejected:
		return "custom-bd-text-danger"

	case RlsBidderSystemInventoryQueue:
		return "custom-bd-text-info"

	case RlsBidderSystemInventoryNotfound:
		return "custom-bd-text-muted"

	case RlsBidderSystemInventoryApprovedS2S:
		return "custom-bd-text-success"

	case RlsBidderSystemInventoryApprovedClient:
		return "custom-bd-text-success"

	case RlsBidderSystemInventoryNotUse:
		return "custom-bd-text-notuse"

	default:
		return ""
	}
}

func (p TYPERlsBidderSystemInventoryStatus) String() string {
	switch p {
	case RlsBidderSystemInventoryPending:
		return "pending"

	case RlsBidderSystemInventorySubmitted:
		return "submited"

	case RlsBidderSystemInventoryApproved:
		return "approved"

	case RlsBidderSystemInventoryRejected:
		return "rejected"

	case RlsBidderSystemInventoryQueue:
		return "queue"

	case RlsBidderSystemInventoryNotfound:
		return "notfound"

	case RlsBidderSystemInventoryApprovedS2S:
		return "approved s2s"

	case RlsBidderSystemInventoryApprovedClient:
		return "approved client"

	case RlsBidderSystemInventoryNotUse:
		return "not use"

	default:
		return ""
	}
}

func GetStatusInventoryBidderByString(status string) TYPERlsBidderSystemInventoryStatus {
	switch status {
	case "pending":
		return RlsBidderSystemInventoryPending

	case "submited":
		return RlsBidderSystemInventorySubmitted

	case "approved":
		return RlsBidderSystemInventoryApproved

	case "rejected":
		return RlsBidderSystemInventoryRejected

	case "queue":
		return RlsBidderSystemInventoryQueue

	case "notfound":
		return RlsBidderSystemInventoryNotfound

	case "approved s2s":
		return RlsBidderSystemInventoryApprovedS2S

	case "approved client":
		return RlsBidderSystemInventoryApprovedClient

	case "not use":
		return RlsBidderSystemInventoryNotUse

	default:
		return RlsBidderSystemInventoryPending
	}
}

func (p TYPERlsBidderSystemInventoryStatus) TextStatus() string {
	switch p {
	case RlsBidderSystemInventoryApproved:
		return "All"

	case RlsBidderSystemInventoryApprovedS2S:
		return "S2S"

	case RlsBidderSystemInventoryApprovedClient:
		return "Client"

	default:
		return ""
	}
}

func (p TYPERlsBidderSystemInventoryStatus) IsApproved() bool {
	switch p {
	case RlsBidderSystemInventoryApproved:
		return true

	case RlsBidderSystemInventoryApprovedS2S:
		return true

	case RlsBidderSystemInventoryApprovedClient:
		return true

	default:
		return false
	}
}
