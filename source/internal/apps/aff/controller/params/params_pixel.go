package params

import (
	"encoding/hex"
	"github.com/google/go-querystring/query"
	"source/internal/entity/model"
	"strconv"
	"strings"
	"time"
)

type pixel struct{}

func NewPixel() *pixel {
	return &pixel{}
}

type pixelQuery struct {
	TtPn  string `url:"tt_pn,omitempty"`
	TtPid string `url:"tt_pid,omitempty"`
	TtCtv string `url:"tt_ctv,omitempty"`
	TtDn  string `url:"tt_dn,omitempty"`
	TtCam string `url:"tt_cam,omitempty"`
	TtSc  string `url:"tt_sc,omitempty"`
	TtCpc string `url:"tt_cpc,omitempty"`
	TtCid string `url:"tt_cid,omitempty"`
	TtPa  string `url:"tt_pa,omitempty"`
	TtTfs string `url:"tt_tfs,omitempty"`
	TtRid string `url:"tt_rid,omitempty"`
	Uuid  string `url:"uuid,omitempty"`
}

func (t *pixel) MakeParams(campaign model.CampaignModel, inputs ParamsCampaignDTO) (paramsString string) {
	pixelQueryBuild := pixelQuery{
		TtPn:  inputs.PublisherName,
		TtPid: inputs.PublisherId,
		TtCtv: inputs.SelectionId,
		TtDn:  inputs.SectionName,
		TtCam: inputs.CampaignId,
		TtSc:  inputs.AdId,
		TtCpc: inputs.ClickPrice,
		TtCid: inputs.ClickId,
		TtPa:  strings.ToLower(campaign.DemandSource),
		TtTfs: strings.ToLower(campaign.TrafficSource),
		TtRid: strconv.FormatInt(campaign.ID, 10),
		Uuid:  "",
	}
	uuid := []byte("aff_pp" + strconv.FormatInt(campaign.ID, 10) + strconv.Itoa(time.Now().Nanosecond()))
	pixelQueryBuild.Uuid = hex.EncodeToString(uuid[:])

	queryString, _ := query.Values(pixelQueryBuild)
	return queryString.Encode()
}
