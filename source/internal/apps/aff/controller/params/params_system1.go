package params

import (
	"github.com/google/go-querystring/query"
	"source/internal/entity/model"
	"strings"
)

type system1 struct{}

func newSystem1() *system1 {
	return &system1{}
}

type system1Query struct {
	Compkey string `url:"compkey,omitempty"`
	Ref     string `url:"ref,omitempty"`
	Rskey   string `url:"rskey,omitempty"`
	SubId   string `url:"sub_id,omitempty"`
	SubId2  string `url:"sub_id2,omitempty"`
	// Obid      string `url:"obid,omitempty"`
	ForcekeyA string `url:"forcekeyA,omitempty"`
	ForcekeyB string `url:"forcekeyB,omitempty"`
	ForcekeyC string `url:"forcekeyC,omitempty"`
	ForcekeyD string `url:"forcekeyD,omitempty"`
	// Obclick   string `url:"obclick,omitempty"`
	// Obclickid string `url:"obclickid,omitempty"`
	ClickId string `url:"click_id,omitempty"`
}

func (t *system1) MakeParams(campaign model.CampaignModel, inputs ParamsCampaignDTO) (paramsString string) {
	system1QueryBuild := system1Query{
		Compkey:   campaign.MainKeyword,
		Ref:       "",
		Rskey:     inputs.AdTitle,
		SubId:     strings.ToLower(campaign.TrafficSource) + "_" + inputs.CampaignId,
		SubId2:    inputs.PublisherId + "_" + inputs.SelectionId,
		ForcekeyA: "",
		ForcekeyB: "",
		ForcekeyC: "",
		ForcekeyD: "",
		ClickId:   inputs.ClickId,
	}
	switch campaign.TrafficSource {
	case "Outbrain":
		system1QueryBuild.Ref = inputs.PublisherName
		// system1QueryBuild.Obid = campaign.PixelId
		break
	case "Mgid":
		system1QueryBuild.Ref = inputs.Referrer
		break
	case "PocPoc":
		system1QueryBuild.Ref = inputs.Referrer
		break
	case "Taboola":
		system1QueryBuild.Ref = inputs.PublisherName
		break
	}

	if len(campaign.Keywords) > 0 {
		for _, keywordCampaign := range campaign.Keywords {
			if keywordCampaign.Keyword == "" {
				continue
			}
			if system1QueryBuild.ForcekeyA == "" {
				system1QueryBuild.ForcekeyA = keywordCampaign.Keyword
			} else if system1QueryBuild.ForcekeyB == "" {
				system1QueryBuild.ForcekeyB = keywordCampaign.Keyword
			} else if system1QueryBuild.ForcekeyC == "" {
				system1QueryBuild.ForcekeyC = keywordCampaign.Keyword
			} else if system1QueryBuild.ForcekeyD == "" {
				system1QueryBuild.ForcekeyD = keywordCampaign.Keyword
			}
		}
	}
	// if inputs.ClickId != "" {
	// 	system1QueryBuild.Obclick = "Click"
	// }
	queryString, _ := query.Values(system1QueryBuild)
	return queryString.Encode()
}
