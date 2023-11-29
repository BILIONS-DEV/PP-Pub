package mysql

type TableRlsConnectionMCM struct {
	Id           int64             `gorm:"column:id" json:"id"`
	BidderId     int64             `gorm:"column:bidder_id" json:"bidder_id"`
	NetworkId    int64             `gorm:"column:network_id" json:"network_id"`
	Status       TYPEConnectionMCM `gorm:"column:status" json:"status"`
}

type TYPEConnectionMCM int

const (
	TYPEConnectionMCMTypeAccept TYPEConnectionMCM = iota + 1
	TYPEConnectionMCMTypePending
	TYPEConnectionMCMTypeReject
)

func (t TYPEConnectionMCM) String() string {
	switch t {
	case TYPEConnectionMCMTypeAccept:
		return "accept"
	case TYPEConnectionMCMTypePending:
		return "pending"
	case TYPEConnectionMCMTypeReject:
		return "reject"
	default:
		return "pending"
	}
}

func (p TYPEConnectionMCM) ColorClass() string {
	switch p {
	case TYPEConnectionMCMTypeAccept:
		return "custom-bd-text-success"

	case TYPEConnectionMCMTypePending:
		return "custom-bd-text-dark"

	case TYPEConnectionMCMTypeReject:
		return "custom-bd-text-danger"

	default:
		return "custom-bd-text-dark"
	}
}
