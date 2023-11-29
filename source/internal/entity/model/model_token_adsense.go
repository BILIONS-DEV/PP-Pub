package model

func (TokenAdsense) TableName() string {
	return "token_adsense"
}

type TokenAdsense struct {
	App    string                 `gorm:"column:app;primaryKey" json:"app"`
	Email  string                 `gorm:"column:email;primaryKey" json:"email"`
	Object TYPETokenAdsenseObject `gorm:"column:object" json:"object"`
	Token  string                 `gorm:"column:token" json:"token"`
	Status string                 `gorm:"column:status" json:"status"`
}

type TYPETokenAdsenseObject string

const (
	TYPETokenAdsenseObjectAdsense2  = "adsense_2"
	TYPETokenAdsenseObjectAdsense3  = "adsense_3"
	TYPETokenAdsenseObjectAdsense24 = "adsense_24"
	TYPETokenAdsenseObjectAdsense30 = "adsense_30"
	TYPETokenAdsenseObjectAdsense31 = "adsense_31"
)

func (t TYPETokenAdsenseObject) IsAdsense2() bool {
	if t == TYPETokenAdsenseObjectAdsense2 {
		return true
	}
	return false
}

func (t TYPETokenAdsenseObject) IsAdsense3() bool {
	if t == TYPETokenAdsenseObjectAdsense3 {
		return true
	}
	return false
}

func (t TYPETokenAdsenseObject) IsAdsense24() bool {
	if t == TYPETokenAdsenseObjectAdsense24 {
		return true
	}
	return false
}

func (t TYPETokenAdsenseObject) IsAdsense30() bool {
	if t == TYPETokenAdsenseObjectAdsense30 {
		return true
	}
	return false
}

func (t TYPETokenAdsenseObject) IsAdsense31() bool {
	if t == TYPETokenAdsenseObjectAdsense31 {
		return true
	}
	return false
}
