package params

import (
	"github.com/google/go-querystring/query"
	"source/internal/entity/model"
)

type adsense struct{}

func NewAdsense() *adsense {
	return &adsense{}
}

type adsenseQuery struct {
	Channel string `url:"channel,omitempty"`
}

func (t *adsense) MakeParams(campaign model.CampaignModel, inputs ParamsCampaignDTO) (paramsString string) {
	adsenseQueryBuild := adsenseQuery{
		Channel: campaign.Channel,
	}
	queryString, _ := query.Values(adsenseQueryBuild)
	return queryString.Encode()
}
