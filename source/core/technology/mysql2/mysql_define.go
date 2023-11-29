package mysql2

// TYPEStatus /**
type TYPEStatus int

const (
	TYPEStatusON TYPEStatus = iota + 1
	TYPEStatusOFF
)

func (t TYPEStatus) String() string {
	switch t {
	case TYPEStatusON:
		return "On"
	case TYPEStatusOFF:
		return "Off"
	}
	return ""
}

func (t TYPEStatus) StringRunPause() string {
	switch t {
	case TYPEStatusON:
		return "Running"
	case TYPEStatusOFF:
		return "Paused"
	}
	return ""
}
