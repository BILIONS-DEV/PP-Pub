package model

func (CampaignTrafficSourceID) TableName() string {
	return "campaign_traffic_source_id"
}

type CampaignTrafficSourceID struct {
	CampaignID      int64  `gorm:"column:campaign_id;primaryKey" json:"campaign_id"`
	TrafficSourceID string `gorm:"column:traffic_source_id;primaryKey" json:"traffic_source_id"`
}
