package params

import (
	"github.com/google/go-querystring/query"
	"source/internal/entity/model"
	"strconv"
	"strings"
	"time"
)

type bodis struct {
}

func NewBodis() *bodis {
	return &bodis{}
}

type bodisQuery struct {
	RefAdnetwork string `url:"ref_adnetwork,omitempty"`
	RefKeyword   string `url:"ref_keyword,omitempty"`
	Terms        string `url:"terms,omitempty"`
	Subid1       string `url:"subid1,omitempty"`
	Subid2       string `url:"subid2,omitempty"`
	Subid3       string `url:"subid3,omitempty"`
	Subid4       string `url:"subid4,omitempty"`
	Subid5       string `url:"subid5,omitempty"`
	RefPubsite   string `url:"ref_pubsite,omitempty"`
	TtPn         string `url:"tt_pn,omitempty"`
	TtPid        string `url:"tt_pid,omitempty"`
	TtCtv        string `url:"tt_ctv,omitempty"`
	TtDn         string `url:"tt_dn,omitempty"`
	TtCam        string `url:"tt_cam,omitempty"`
	TtSc         string `url:"tt_sc,omitempty"`
	ClickId      string `url:"click_id,omitempty"`
	Cache        string `url:"cache,omitempty"`
}

func (t *bodis) MakeParams(campaign model.CampaignModel, inputs ParamsCampaignDTO) (paramString string) {
	bodisQueryBuid := bodisQuery{
		RefAdnetwork: campaign.TrafficSource,
		RefKeyword:   inputs.AdTitle,
		Terms:        "",
		Subid1:       campaign.TrafficSource,
		Subid2:       "campaignID-" + inputs.CampaignId,
		Subid3:       "adID-" + inputs.AdId,
		Subid4:       "publisherID-" + inputs.PublisherId,
		Subid5:       "selectionID-" + inputs.SelectionId,
		RefPubsite:   "",
		TtPn:         inputs.PublisherName,
		TtPid:        inputs.PublisherId,
		TtCtv:        inputs.SelectionId,
		TtDn:         inputs.SectionName,
		TtCam:        inputs.CampaignId,
		TtSc:         inputs.AdId,
		ClickId:      inputs.ClickId,
		Cache:        strconv.FormatInt(time.Now().UnixNano(), 10),
	}

	// => &terms = list keyword
	if len(campaign.Keywords) > 0 {
		for _, keywordCampaign := range campaign.Keywords {
			if keywordCampaign.Keyword == "" {
				continue
			}
			if strings.Index(bodisQueryBuid.Terms, "terms=") == -1 {
				bodisQueryBuid.Terms = keywordCampaign.Keyword
			} else {
				bodisQueryBuid.Terms = bodisQueryBuid.Terms + "," + keywordCampaign.Keyword
			}

		}
	}

	// => &ref_pubsite kiểm tra theo TrafficSource sẽ lấy Inputs khác nhau
	switch campaign.TrafficSource {
	case "Outbrain":
		bodisQueryBuid.RefPubsite = inputs.SectionName
		break
	case "PocPoc", "Mgid":
		bodisQueryBuid.RefPubsite = inputs.Referrer
		break
	case "Taboola":
		bodisQueryBuid.RefPubsite = inputs.PublisherName
		break
	}

	queryString, _ := query.Values(bodisQueryBuid)
	return queryString.Encode()
}
