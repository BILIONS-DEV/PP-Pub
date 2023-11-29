package model

/**
TYPEOnOff
*/

type TYPEState int

const (
	On TYPEState = iota + 1
	Off
)

func (t TYPEState) String() string {
	switch t {
	case On:
		return "on"
	case Off:
		return "off"
	}
	return ""
}
func (t TYPEState) StringIsLock() string {
	switch t {
	case On:
		return "lock"
	case Off:
		return "unlock"
	}
	return ""
}
func (t TYPEState) StringIsDisabled() string {
	switch t {
	case On:
		return "enabled"
	case Off:
		return "disabled"
	}
	return ""
}
func (t TYPEState) Boolean() bool {
	switch t {
	case On:
		return true
	case Off:
		return false
	}
	return false
}

/**
TYPEStatus
*/

type TYPEStatus int

const (
	StatusApproved TYPEStatus = iota + 1
	StatusPending
	StatusReject
)

func (s TYPEStatus) String() string {
	switch s {
	case StatusApproved:
		return "approved"
	case StatusPending:
		return "pending"
	case StatusReject:
		return "rejected"
	default:
		return ""
	}
}
func (s TYPEStatus) Int() int {
	switch s {
	case StatusApproved:
		return 1
	case StatusPending:
		return 2
	case StatusReject:
		return 3
	default:
		return 0
	}
}
