package model

type InventoryConnectionDemandModel struct {
	ID          int64
	UserID      int64
	InventoryID int64
	BidderID    int64
	Status      TYPEStatusInventoryConnectionDemand
}

type TYPEStatusInventoryConnectionDemand int

const (
	TYPEStatusConnectionDemandLive TYPEStatusInventoryConnectionDemand = iota + 1
	TYPEStatusConnectionDemandWaiting
)

func (t TYPEStatusInventoryConnectionDemand) String() string {
	switch t {
	case TYPEStatusConnectionDemandLive:
		return "live"
	case TYPEStatusConnectionDemandWaiting:
		return "waiting"
	default:
		return "waiting"
	}
}
