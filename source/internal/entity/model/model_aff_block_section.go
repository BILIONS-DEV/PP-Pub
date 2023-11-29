package model

func (AffBlockSectionModel) TableName() string {
	return "aff_block_section"
}

func (t *AffBlockSectionModel) IsFound() bool {
	if t.ID > 0 {
		return true
	}
	return false

}

type AffBlockSectionModel struct {
	ID            int64  `gorm:"column:id;primaryKey;autoIncrement"`
	TrafficSource string `gorm:"column:traffic_source"`
	CampaignID    string `gorm:"column:campaign_id"`
	SectionID     string `gorm:"column:section_id"`
}
