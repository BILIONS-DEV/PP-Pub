package payload

import (
	"source/core/technology/mysql"
	"source/pkg/datatable"
)

type SystemIndex struct {
	SystemFilterPostData
	QuerySearch string   `query:"f_q" json:"f_q" form:"f_q"`
	Status      []string `query:"f_status" form:"f_status" json:"f_status"`
	MediaType   []string `query:"f_media" form:"f_media" json:"f_media"`
}

type SystemFilterPayload struct {
	datatable.Request
	PostData *SystemFilterPostData `query:"postData"`
}

type SystemCreate struct {
	Id                int64                 `query:"id" json:"id" form:"id"`
	DisplayName       string                `query:"display_name" json:"display_name" form:"display_name"`
	BidderAlias       string                `query:"bidder_alias" json:"bidder_alias" form:"bidder_alias"`
	PubId             string                `query:"pub_id" json:"pub_id" form:"pub_id"`
	LinkedGam         int64                 `query:"linked_gam" json:"linked_gam" form:"linked_gam"`
	AccountType       mysql.TYPEAccountType `query:"account_type" json:"account_type" form:"account_type"`
	BidderId          int64                 `query:"bidder_id" json:"bidder_id" form:"bidder_id"`
	MediaType         []string              `query:"media_type" json:"media_type" form:"media_type"`
	BidAdjustment     float64               `query:"bid_adjustment" json:"bid_adjustment" form:"bid_adjustment"`
	RPM               float64               `query:"rpm" json:"rpm" form:"rpm"`
	AdsTxt            string                `query:"ads_text" json:"ads_txt"`
	XlsxPath          string                `query:"xlsx_file" json:"xlsx_file"`
	IsLock            mysql.TYPEIsLock      `query:"-" json:"-"`
	IsDefault         mysql.TypeOnOff       `query:"-" json:"-"`
	Params            []Param               `query:"params" json:"params" form:"params"`
	SupplyChainDomain string                `query:"supply_chain_domain" json:"supply_chain_domain" json:"supply_chain_domain"`
}

type SystemFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Status      interface{} `query:"f_status[]" json:"f_status" form:"f_status[]"`
	MediaType   interface{} `query:"f_media[]" json:"f_media" form:"f_media[]"`
	Page        int         `query:"page" json:"page" form:"page"`
	Limit       int         `query:"limit" json:"limit" form:"limit"`
	Start       int         `query:"start" json:"start" form:"start"`
	Length      int         `query:"length" json:"length" form:"length"`
}

type SystemConfig struct {
	Id            int64 `json:"id"`
	PrebidTimeOut int   `json:"prebid_time_out"`
	AdRefreshTime int   `json:"ad_refresh_time"`
	Currency      int64 `json:"currency"`
}

type ResponseCurrency struct {
	DataAsOf    string      `json:"dataAsOf"`
	GeneratedAt string      `json:"generatedAt"`
	Conversions Conversions `json:"conversions"`
}

type Conversions struct {
	USD map[string]float64 `json:"USD"`
}

type Param struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Index int    `json:"index"`
}
