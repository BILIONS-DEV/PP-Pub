package model

func (CampaignModel) TableName() string {
	return "campaign"
}

type CampaignModel struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement;type:int(11)" json:"id"`
	Name string `gorm:"column:name" json:"name"`
	//TrafficSource string                 `gorm:"column:traffic_source" json:"traffic_source"`
	//DemandSource  string                 `gorm:"column:demand_source" json:"demand_source"`
	//PixelId       string                 `gorm:"column:pixel_id" json:"pixel_id"`
	//LandingPages  string                 `gorm:"column:landing_pages" json:"landing_pages"`
	//MainKeyword   string                 `gorm:"column:main_keyword" json:"main_keyword"`
	//Channel       string                 `gorm:"column:channel" json:"channel"`
	//GD            string                 `gorm:"column:gd" json:"gd"`
	//Params        string                 `gorm:"column:params" json:"params"`
	//UrlTrackImp   string                 `gorm:"column:url_track_imp" json:"url_track_imp"`
	//URLTrackClick string                 `gorm:"column:url_track_click" json:"url_track_click"`
	Keywords []CampaignKeywordModel `gorm:"foreignKey:CampaignID;references:ID;constraint:OnUpdate:CASCADE;"`
	//Creative      CreativeModel          `gorm:"foreignKey:CampaignID;references:ID"`
	//CreatedAt     time.Time              `gorm:"column:created_at" json:"created_at"`
	//DeletedAt     gorm.DeletedAt         `gorm:"column:deleted_at" json:"deleted_at"`
}

func (t CampaignModel) Validate() (err error) {
	// if {
	//
	// }
	return
}

func (CampaignKeywordModel) TableName() string {
	return "campaign_keyword"
}

type CampaignKeywordModel struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement;type:int(11)" json:"id"`
	CampaignID int64  `gorm:"column:campaign_id" json:"campaign_id"`
	Keyword    string `gorm:"column:keyword" json:"keyword"`
}

type CreativeModel struct {
	ID         int64                `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CampaignID int64                `gorm:"column:campaign_id" json:"campaign_id"`
	SizeName   string               `gorm:"column:size_name" json:"size_name"`
	Titles     []CreativeTitleModel `gorm:"foreignKey:CreativeID;references:ID"`
	Images     []CreativeImageModel `gorm:"foreignKey:CreativeID;references:ID"`
}

func (CreativeModel) TableName() string {
	return "campaign_creative"
}

type CreativeTitleModel struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CreativeID int64  `gorm:"column:creative_id" json:"creative_id"`
	Title      string `gorm:"column:title" json:"title"`
}

func (CreativeTitleModel) TableName() string {
	return "creative_title"
}

type CreativeImageModel struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CreativeID int64  `gorm:"column:creative_id" json:"creative_id"`
	Image      string `gorm:"column:image" json:"image"`
}

func (CreativeImageModel) TableName() string {
	return "creative_image"
}

// type ResponseData struct {
// 	Subid                      interface{} `gorm:"column:subid" json:"subid"`
// 	Visits                     int         `gorm:"column:visits" json:"visits"`
// 	LandingPageVisits          int         `gorm:"column:landing_page_visits" json:"landing_page_visits"`
// 	Zeroclick_visits           int         `gorm:"column:zeroclick_visits" json:"zeroclick_visits"`
// 	Clicks                     int         `gorm:"column:clicks" json:"clicks"`
// 	CreditedRevenue            float64     `gorm:"column:credited_revenue" json:"credited_revenue"`
// 	LandingPageCreditedRevenue float64     `gorm:"column:landing_page_credited_revenue" json:"landing_page_credited_revenue"`
// 	ZeroclickCreditedRevenue   float64     `gorm:"column:zeroclick_credited_revenue" json:"zeroclick_credited_revenue"`
// 	Ctr                        float64     `gorm:"column:ctr" json:"ctr"`
// 	Epc                        float64     `gorm:"column:epc" json:"epc"`
// 	Rpm                        float64     `gorm:"column:rpm" json:"rpm"`
// 	LandingPageRpm             float64     `gorm:"column:landing_page_rpm" json:"landing_page_rpm"`
// 	ZeroclickRpm               float64     `gorm:"column:zeroclick_rpm" json:"zeroclick_rpm"`
// 	ClicksSpamRatio            float64     `gorm:"column:clicks_spam_ratio" json:"clicks_spam_ratio"`
// 	IsFinalized                float64     `gorm:"column:is_finalized" json:"is_finalized"`
// }
