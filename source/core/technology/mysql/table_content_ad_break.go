package mysql

type TableContentAdBreak struct {
	Id          int64         `json:"id"`
	ContentId   int64         `json:"content_id"`
	Type        string        `json:"type"`
	TimeAdBreak string        `json:"time_ad_break"`
	BreakMode   TYPEBreakMode `json:"break_mode"`
}

func (TableContentAdBreak) TableName() string {
	return Tables.ContentAdBreak
}

type TYPEBreakMode int

const (
	TYPEBreakModeSecondsIntoVideo TYPEBreakMode = iota + 1
	TYPEBreakModeTimeCode
	TYPEBreakModePercentOfVideo
)

func (s TYPEBreakMode) String() string {
	switch s {
	case TYPEBreakModeSecondsIntoVideo:
		return "Seconds Into Video"
	case TYPEBreakModeTimeCode:
		return "Time Code"
	case TYPEBreakModePercentOfVideo:
		return "Percent Of Video"
	default:
		return ""
	}
}

func (s TYPEBreakMode) Int() int {
	switch s {
	case TYPEBreakModeSecondsIntoVideo:
		return 1
	case TYPEBreakModeTimeCode:
		return 2
	case TYPEBreakModePercentOfVideo:
		return 3
	default:
		return 0
	}
}
