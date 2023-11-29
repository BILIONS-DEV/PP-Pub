package mysql

type TableInventoryConnectionDemand struct {
	Id          int64
	UserId      int64
	InventoryId int64
	BidderId    int64
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