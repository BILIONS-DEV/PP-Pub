package model

func (LandingPagesDemand) TableName() string {
	return "landing_pages_demand"
}

type LandingPagesDemand struct {
	ID           int64  `gorm:"column:id" json:"id"`
	DemandSource string `gorm:"column:demand_source" json:"demand_source"`
	UserID       int64  `gorm:"column:user_id" json:"user_id"`
	Name         string `gorm:"column:name" json:"name"`
	LandingPage  string `gorm:"column:landing_page" json:"landing_page"`
	MainKeyword  string `gorm:"column:main_keyword" json:"main_keyword"`
}
