package model

func (Section) TableName() string {
	return "section"
}

type Section struct {
	TrafficSource string `gorm:"column:traffic_source;primaryKey" json:"traffic_source"`
	SectionID     string `gorm:"column:section_id;primaryKey" json:"section_id"`
	SectionName   string `gorm:"column:section_name" json:"section_name"`
	Referrer      string `gorm:"column:referrer"  json:"referrer"`
}
