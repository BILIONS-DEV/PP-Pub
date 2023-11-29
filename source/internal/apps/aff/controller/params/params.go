package params

import (
	"source/internal/entity/model"
)

type Handle interface {
	MakeParams(campaign model.CampaignModel, inputs ParamsCampaignDTO) (paramsString string)
}

type ParamsCampaignDTO struct {
	AdTitle       string `query:"ad_title" json:"ad_title"`
	PublisherName string `query:"publisher_name" json:"publisher_name"` // tt_pn
	PublisherId   string `query:"publisher_id" json:"publisher_id"`     // tt_pid
	SelectionId   string `query:"selection_id" json:"selection_id"`     // tt_ctv
	SectionName   string `query:"section_name" json:"section_name"`     // tt_dn
	CampaignId    string `query:"campaign_id" json:"campaign_id"`       // tt_cam
	AdId          string `query:"ad_id" json:"ad_id"`                   // tt_sc
	ClickPrice    string `query:"click_price" json:"click_price"`       // tt_cpc
	ClickId       string `query:"click_id" json:"click_id"`
	Referrer      string `query:"referrer" json:"referrer"`
}

func NewParam(paramType string) Handle {
	switch paramType {
	case "system1":
		return newSystem1()
	case "bodis":
		return NewBodis()
	case "adense":
		return NewAdsense()
	case "pixel":
		return NewPixel()
	default:
		return nil
	}
}
