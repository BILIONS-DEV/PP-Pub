package dto

import (
	"errors"
	// "fmt"
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
	"source/pkg/utility"
)

type PayloadCampaignFilter struct {
	datatable.Request
	OrderColumn   int                            `query:"order_column" json:"order_column" form:"order_column"`
	OrderDir      string                         `query:"order_dir" json:"order_dir" form:"order_dir"`
	Page          int                            `query:"page" json:"page" form:"page"`
	QuerySearch   string                         `query:"f_q" json:"f_q" form:"f_q"`
	Limit         int                            `query:"limit" json:"limit" form:"limit"`
	Start         int                            `query:"start" json:"start" form:"start"`
	Length        int                            `query:"length" json:"length" form:"length"`
	TrafficSource string                         `query:"traffic_source" json:"traffic_source" form:"traffic_source"`
	DemandSource  string                         `query:"demand_source" json:"demand_source" form:"demand_source"`
	PostData      *PayloadCampaignFilterPostData `query:"postData"`
}

type PayloadCampaignFilterPostData struct {
	QuerySearch   string `query:"f_q" json:"f_q" form:"f_q"`
	TrafficSource string `query:"traffic_source" json:"traffic_source" form:"traffic_source"`
	DemandSource  string `query:"demand_source" json:"demand_source" form:"demand_source"`
}

func (t *PayloadCampaignFilter) Validate() (err error) {
	return
}

type PayloadCampaignAdd struct {
	Id            int                `json:"id"`
	Name          string             `json:"name"`
	TrafficSource string             `json:"traffic_source"`
	DemandSource  string             `json:"demand_source"`
	UserID        int64              `json:"user_id"`
	CreativeID    int                `json:"creative_id"`
	Vertical      string             `json:"vertical"`
	LandingPages  string             `json:"landing_pages"`
	MainKeyword   string             `json:"main_keyword"`
	Channel       string             `json:"channel"`
	GD            string             `json:"gd"`
	UserFlow      string             `json:"user_flow"`
	Keywords      []CampaignKeywords `json:"keywords"`
	// SiteName string `json:"site_name"`
	// Titles          []CreativeTitle  `json:"titles"`
	// Images          []CreativeImage  `json:"images"`
	Account         string           `json:"account"`
	AccountAdsense  string           `json:"account_adsense"`
	LandingPagesGeo []LandingPageGeo `json:"landing_pages_geo"`
	Devices         []string         `json:"device"`
	Bidding         string           `json:"bidding"`
	Cpc             float64          `json:"cpc"`
	ExpectedCpa     float64          `json:"expected_cpa"`
	Budget          float64          `json:"budget"`
	Locations       []string         `json:"location"`
	LocationType    string           `json:"location_type"`
	Countries       []string         `json:"country"`
	Terms           string           `json:"terms"`
	Adtitle         string           `json:"ad_title"`
}

type CampaignKeywords struct {
	Id         int    `json:"id"`
	CampaignId int    `json:"campaign_id"`
	Keyword    string `json:"keyword"`
}

type LandingPageGeo struct {
	ID          int64  `json:"id"`
	Geo         string `json:"geo"`
	LandingPage string `json:"landing_page"`
}

type PayloadDelete struct {
	Id int64 `json:"id"`
}

func (p *PayloadCampaignAdd) ToModel() model.CampaignModel {
	record := model.CampaignModel{}
	ID := int64(p.Id)
	record.ID = ID
	record.Name = p.Name
	record.TrafficSource = p.TrafficSource
	record.DemandSource = p.DemandSource
	record.UserID = p.UserID
	record.CreativeID = int64(p.CreativeID)
	if len(p.LandingPages) > 0 {
		if p.LandingPages[len(p.LandingPages)-1:] == "/" {
			last := len(p.LandingPages) - 1
			p.LandingPages = p.LandingPages[:last]
		}
	}
	record.LandingPages = p.LandingPages
	record.MainKeyword = p.MainKeyword
	record.Account = p.Account
	record.Params = ""
	record.UrlTrackImp = ""
	record.URLTrackClick = ""
	Keywords := []model.CampaignKeywordModel{}
	for _, value := range p.Keywords {
		if value.Keyword == "" {
			continue
		}
		var keyword model.CampaignKeywordModel
		// keyword.ID = int64(value.Id)
		keyword.CampaignID = int64(value.CampaignId)
		keyword.Keyword = value.Keyword
		Keywords = append(Keywords, keyword)
	}
	record.Keywords = Keywords

	var Devices []model.CampaignDeviceModel
	for _, device := range p.Devices {
		if device == "" {
			continue
		}
		var Device model.CampaignDeviceModel
		Device.Deivce = device
		Devices = append(Devices, Device)
	}
	record.Device = Devices
	record.Bidding = p.Bidding
	record.Cpc = p.Cpc
	record.ExpectedCpa = p.ExpectedCpa
	record.Budget = p.Budget

	switch record.DemandSource {
	case "adsense":
		record.Terms = p.Terms
		record.Channel = p.Channel
		record.AccountAdsense = p.AccountAdsense
		break
	case "codefuel":
		record.GD = p.GD
		record.UserFlow = p.UserFlow
		record.Vertical = p.Vertical
		break
	case "tonic":
		var landingPages []model.CampaignLandingPageModel
		for _, value := range p.LandingPagesGeo {
			if value.LandingPage == "" || value.Geo == "" {
				continue
			}
			var landingPage model.CampaignLandingPageModel
			landingPage.LandingPage = value.LandingPage
			landingPage.Country = value.Geo
			landingPage.CampaignID = record.ID
			landingPages = append(landingPages, landingPage)
		}
		record.LandingPageGeo = landingPages
		break
	default:
		break
	}

	switch record.TrafficSource {
	case "Outbrain":
		var Locations []model.CampaignLocationModel
		for _, location := range p.Locations {
			if location == "" {
				continue
			}
			var Location model.CampaignLocationModel
			Location.Location = location
			Location.LocationType = p.LocationType
			Locations = append(Locations, Location)
		}
		record.Location = Locations
		break
	case "Google":
		record.AdTitle = p.Adtitle
		break

	default:
		var Countries []model.CampaignLocationModel
		for _, country := range p.Countries {
			if country == "" {
				continue
			}
			var Country model.CampaignLocationModel
			Country.Location = country
			Country.LocationType = p.LocationType
			Countries = append(Countries, Country)
		}
		record.Location = Countries
		break
	}

	if record.ID == 0 {
		record.Flag = "new"
	}
	return record
}

func (p *PayloadCampaignAdd) ToCampaignGroupModel() model.CampaignGroupModel {
	record := model.CampaignGroupModel{}
	record.ID = int64(p.Id)
	record.Name = p.Name
	record.TrafficSource = p.TrafficSource
	record.DemandSource = p.DemandSource
	record.UserID = p.UserID
	record.CreativeID = int64(p.CreativeID)

	return record
}

// func (t *PayloadCampaignFilter) ToCondition() map[string]interface{} {
// 	var condition = make(map[string]interface{})
// 	if t.SubID != "" {
// 		condition["subid"] = t.SubID
// 	}
//
// 	return condition
// }

// type PayloadGetHistory struct {
// 	Id     string  `json:"id" form:"id"`
// 	Object string `json:"object" form:"object"`
// }

type PayloadLandingPages struct {
	DemandSource string   `json:"demand_source"`
	Name         []string `json:"name"`
	LandingPages []string `json:"landing_pages"`
	MainKeyword  []string `json:"main_keyword"`
}

func (t *PayloadLandingPages) Validate() (err error) {
	if utility.ValidateString(t.DemandSource) == "" {
		err = errors.New("Demand Source is required")
		return
	}
	if len(t.LandingPages) == 0 {
		err = errors.New("Landing Pages is required")
	}
	if len(t.Name) == 0 {
		err = errors.New("Name is required")
	}
	return
}
