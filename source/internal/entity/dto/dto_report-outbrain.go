package dto

// => Tạo struct nhận response từ report
type Metadata struct {
	ID                           string  `json:"id"`
	Name                         string  `json:"name"`
	CampaignOnAir                bool    `json:"campaignOnAir"`
	OnAirReason                  string  `json:"onAirReason"`
	Enabled                      bool    `json:"enabled"`
	IsVideo                      bool    `json:"isVideo"`
	Cpc                          float64 `json:"cpc"`
	TargetedSegments             int     `json:"targetedSegments"`
	OptimizationExperimentStatus string  `json:"optimizationExperimentStatus"`
	CreativeFormat               string  `json:"creativeFormat"`
}
type Result struct {
	Metadata Metadata `json:"metadata"`
}
type ResponseCampaignOutBrain struct {
	Results []Result `json:"results"`
}

type OutBrainCampaign struct {
	Id                            string  `json:"id"`
	Name                          string  `json:"name"`
	MarketerId                    string  `json:"marketerId"`
	LastModified                  string  `json:"lastModified"`
	CreationTime                  string  `json:"creationTime"`
	CampaignOnAir                 bool    `json:"campaignOnAir"`
	OnAirReason                   string  `json:"onAirReason"`
	Enabled                       bool    `json:"enabled"`
	AutoArchived                  bool    `json:"autoArchived"`
	Currency                      string  `json:"currency"`
	Cpc                           float64 `json:"cpc"`
	MinimumCpc                    float64 `json:"minimumCpc"`
	AutoExpirationOfPromotedLinks int     `json:"autoExpirationOfPromotedLinks"`
	AmountSpent                   int     `json:"amountSpent"`
	Targeting                     struct {
		Platform  []string `json:"platform"`
		Locations []struct {
			Name          string `json:"name"`
			CanonicalName string `json:"canonicalName"`
			Id            string `json:"id"`
			GeoType       string `json:"geoType"`
			Parent        struct {
				GeoType       string `json:"geoType"`
				CanonicalName string `json:"canonicalName"`
				Id            string `json:"id"`
				Name          string `json:"name"`
			} `json:"parent"`
		} `json:"locations"`
		OperatingSystems []string `json:"operatingSystems"`
		Browsers         []string `json:"browsers"`
	} `json:"targeting"`
	Feeds  []string `json:"feeds"`
	Budget struct {
		Id              string  `json:"id"`
		Name            string  `json:"name"`
		Shared          bool    `json:"shared"`
		Amount          float64 `json:"amount"`
		AmountRemaining float64 `json:"amountRemaining"`
		CreationTime    string  `json:"creationTime"`
		LastModified    string  `json:"lastModified"`
		StartDate       string  `json:"startDate"`
		RunForever      bool    `json:"runForever"`
		Type            string  `json:"type"`
		Currency        string  `json:"currency"`
	} `json:"budget"`
	LiveStatus struct {
		CampaignOnAir bool    `json:"campaignOnAir"`
		OnAirReason   string  `json:"onAirReason"`
		AmountSpent   float64 `json:"amountSpent"`
	} `json:"liveStatus"`
	SuffixTrackingCode string `json:"suffixTrackingCode"`
	PrefixTrackingCode struct {
		Prefix string `json:"prefix"`
		Encode bool   `json:"encode"`
	} `json:"prefixTrackingCode"`
	ContentType     string              `json:"contentType"`
	CpcPerAdEnabled bool                `json:"cpcPerAdEnabled"`
	BlockedSites    OutBrainBlockedSite `json:"blockedSites"`
	Bids            struct {
		BySection []struct {
			SectionId     string  `json:"sectionId"`
			SectionName   string  `json:"sectionName"`
			PublisherId   string  `json:"publisherId"`
			CpcAdjustment float64 `json:"cpcAdjustment"`
			CampaignId    string  `json:"campaignId"`
			CreationTime  string  `json:"creationTime"`
			LastModified  string  `json:"lastModified"`
		} `json:"bySection"`
	} `json:"bids"`
	CampaignOptimization struct {
		OptimizationType string `json:"optimizationType"`
	} `json:"campaignOptimization"`
	OnAirType  string `json:"onAirType"`
	Scheduling struct {
		Cpc []struct {
			StartDay      string  `json:"startDay"`
			EndDay        string  `json:"endDay"`
			StartHour     int     `json:"startHour"`
			EndHour       int     `json:"endHour"`
			CpcAdjustment float64 `json:"cpcAdjustment"`
		} `json:"cpc"`
		OnAir []struct {
			StartDay  string `json:"startDay"`
			EndDay    string `json:"endDay"`
			StartHour int    `json:"startHour"`
			EndHour   int    `json:"endHour"`
		} `json:"onAir"`
	} `json:"scheduling"`
	Objective      string `json:"objective"`
	CreativeFormat string `json:"creativeFormat"`
	Pixels         struct {
		TrackingPixels   []string `json:"trackingPixels"`
		ImpressionPixels []string `json:"impressionPixels"`
	} `json:"pixels"`
}

type OutBrainBlockedSite struct {
	BlockedPublishers    []OutBrainBlockedPublisher  `json:"blockedPublishers"`
	BlockedSections      []OutBrainBlockedSection    `json:"blockedSections"`
	MarketerBlockedSites OutBrainMarketerBlockedSite `json:"marketerBlockedSites"`
}

type OutBrainBlockedSection struct {
	Id           string                          `json:"id"`
	Name         string                          `json:"name"`
	Publisher    OutBrainBlockedSectionPublisher `json:"publisher"`
	CreationTime string                          `json:"creationTime"`
	ModifiedBy   string                          `json:"modifiedBy"`
}

type OutBrainBlockedSectionPublisher struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type OutBrainBlockedPublisher struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	CreationTime string `json:"creationTime"`
	ModifiedBy   string `json:"modifiedBy"`
}

type OutBrainMarketerBlockedSite struct {
	BlockedPublishers []OutBrainBlockedPublisher `json:"blockedPublishers"`
	BlockedSections   []OutBrainBlockedSection   `json:"blockedSections"`
}
