package model

import (
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"time"
)

func (ReportOutBrainModel) TableName() string {
	return "report_outbrain"
}

type ReportOutBrainModel struct {
	ID           int64   `gorm:"column:id;primaryKey;autoIncrement"`
	MarketerID   string  `gorm:"column:marketer_id"`
	CampaignID   string  `gorm:"column:campaign_id"`
	CampaignName string  `gorm:"column:campaign_name"`
	SectionID    string  `gorm:"column:section_id"`
	Time         string  `gorm:"column:time"`
	Spend        float64 `gorm:"column:spend"`
	Clicks       int64   `gorm:"column:clicks"`
	Conversions  int64   `gorm:"column:conversions"`
}

func (t *ReportOutBrainModel) IsFound() bool {
	if t.ID > 0 {
		return true
	}
	return false
}

func (t *ReportOutBrainModel) Validate() (err error) {
	if govalidator.IsNull(t.MarketerID) {
		return errors.New("MarketerID required")
	}
	if govalidator.IsNull(t.CampaignID) {
		return errors.New("CampaignID required")
	}
	if govalidator.IsNull(t.SectionID) {
		return errors.New("SectionID required")
	}
	if govalidator.IsNull(t.Time) {
		return errors.New("Time required")
	}
	if t.Spend == 0 {
		return errors.New("Spend > 0")
	}
	return
}

func (t *ReportOutBrainModel) ToMessageKafka() (request TrackingAdsMessage) {
	loc, _ := time.LoadLocation("America/New_York")
	layout := "2006-01-02 15:04:05"
	date, err := time.ParseInLocation(layout, t.Time, loc)
	if err != nil {
		fmt.Println(err)
		return
	}
	layout = "2006-01-02T15:04:05"
	request = TrackingAdsMessage{
		TrafficSource: "Outbrain",
		SelectionId:   t.SectionID,
		Campaign:      t.CampaignID,
		CampaignName:  t.CampaignName,
		Time:          date.UTC().Format(layout),
		Amount:        t.Spend,
	}
	return
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
		Id              string `json:"id"`
		Name            string `json:"name"`
		Shared          bool   `json:"shared"`
		Amount          int    `json:"amount"`
		AmountRemaining int    `json:"amountRemaining"`
		CreationTime    string `json:"creationTime"`
		LastModified    string `json:"lastModified"`
		StartDate       string `json:"startDate"`
		RunForever      bool   `json:"runForever"`
		Type            string `json:"type"`
		Currency        string `json:"currency"`
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
	ContentType     string `json:"contentType"`
	CpcPerAdEnabled bool   `json:"cpcPerAdEnabled"`
	BlockedSites    struct {
		BlockedPublishers []struct {
			Id           string `json:"id"`
			Name         string `json:"name"`
			CreationTime string `json:"creationTime"`
		} `json:"blockedPublishers"`
		BlockedSections []struct {
			Id           string `json:"id"`
			Name         string `json:"name"`
			CreationTime string `json:"creationTime"`
			Publisher    struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			} `json:"publisher"`
		} `json:"blockedSections"`
		MarketerBlockedSites struct {
		} `json:"marketerBlockedSites"`
	} `json:"blockedSites"`
	Bids struct {
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

type OutBrainCampaignUpdate struct {
	Name      string  `json:"name"`
	Enabled   bool    `json:"enabled"`
	Cpc       float64 `json:"cpc"`
	Targeting struct {
		Platform          []string `json:"platform"`
		Locations         []string `json:"locations"`
		ExcludedLocations []string `json:"excludedLocations"`
		CustomAudience    struct {
			IncludedSegments []string `json:"includedSegments"`
			ExcludedSegments []string `json:"excludedSegments"`
		} `json:"customAudience"`
		OperatingSystems []string `json:"operatingSystems"`
		Browsers         []string `json:"browsers"`
	} `json:"targeting"`
	SuffixTrackingCode string `json:"suffixTrackingCode"`
	PrefixTrackingCode struct {
		Prefix string `json:"prefix"`
		Encode bool   `json:"encode"`
	} `json:"prefixTrackingCode"`
	CpcPerAdEnabled bool `json:"cpcPerAdEnabled"`
	BlockedSites    struct {
		BlockedPublishers []struct {
			Id string `json:"id"`
		} `json:"blockedPublishers"`
		BlockedSections []struct {
			Id string `json:"id"`
		} `json:"blockedSections"`
	} `json:"blockedSites"`
	Archived bool `json:"archived"`
	Bids     struct {
		BySection []struct {
			SectionId     string  `json:"sectionId"`
			CpcAdjustment float64 `json:"cpcAdjustment"`
		} `json:"bySection"`
	} `json:"bids"`
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
	Objective string `json:"objective"`
	Pixels    struct {
		TrackingPixels   []string `json:"trackingPixels"`
		ImpressionPixels []string `json:"impressionPixels"`
	} `json:"pixels"`
}
