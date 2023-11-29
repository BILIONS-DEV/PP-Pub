package payload

import (
	"source/core/technology/mysql"
	"source/pkg/datatable"
)

type LineItemSystemIndex struct {
	InventoryAdTagFilterPostData
	QuerySearch string   `query:"f_q" json:"f_q" form:"f_q"`
	Status      []string `query:"f_status" form:"f_status" json:"f_status"`
	Type        []string `query:"f_type" form:"f_type" json:"f_type"`
	Domain      []string `query:"f_domain" form:"f_domain" json:"f_domain"`
	AdFormat    []string `query:"f_ad_format" form:"f_ad_format" json:"f_ad_format"`
	AdSize      []string `query:"f_ad_size" form:"f_ad_size" json:"f_ad_size"`
	AdTag       []string `query:"f_ad_tag" form:"f_ad_tag" json:"f_ad_tag"`
	Device      []string `query:"f_device" form:"f_device" json:"f_device"`
	Country     []string `query:"f_geo" form:"f_geo" json:"f_geo"`
}

type LineItemSystemFilterPayload struct {
	datatable.Request
	PostData *LineItemSystemFilterPostData `query:"postData"`
}

type LineItemSystemFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Status      interface{} `query:"f_status[]" json:"f_status" form:"f_status[]"`
	Type        interface{} `query:"f_type[]" json:"f_type" form:"f_type[]"`
	Page        int         `query:"page" json:"page" form:"page"`
	Limit       int         `query:"limit" json:"limit" form:"limit"`
	Start       int         `query:"start" json:"start" form:"start"`
	Length      int         `query:"length" json:"length" form:"length"`
}

type LineItemSystemAdd struct {
	Id              int64                `json:"id" form:"id"`
	Name            string               `json:"name" form:"name"`
	Description     string               `json:"description" form:"description"`
	ServerType      mysql.TYPEServerType `json:"server_type" form:"server_type"`
	SelectAccount   []string             `json:"select_account" form:"select_account"`
	LinkedGam       int64                `json:"linked_gam" form:"linked_gam"`
	Status          string               `json:"status" form:"status"`
	BidderType      mysql.TYPEBidderType `json:"type" form:"type"`
	BidderParams    []BidderInfo         `json:"bidder_params" form:"bidder_params"`
	AdsenseAdSlots  []AdsenseAdSlot      `json:"adsense_ad_slots" form:"adsense_ad_slots"`
	BidderRate      int                  `json:"rate" form:"rate"`
	BidderVastUrl   string               `json:"vast_url" form:"vast_url"`
	BidderAdTag     string               `json:"ad_tag" form:"ad_tag"`
	Priority        int                  `json:"priority" form:"priority"`
	StartDate       string               `json:"start_date" form:"start_date"`
	EndDate         string               `json:"end_date" form:"end_date"`
	ListAdFormat    []ListTarget         `json:"listAdFormat" query:"listAdFormat" form:"listAdFormat"`
	ListAdSize      []ListTarget         `json:"listAdSize" query:"listAdSize" form:"listAdSize"`
	ListAdTag       []ListTarget         `json:"listAdtag" query:"listAdtag" form:"listAdtag"`
	ListAdInventory []ListTarget         `json:"listInventory" query:"listInventory" form:"listInventory"`
	ListGeo         []ListTarget         `json:"listGeo" query:"listGeo" form:"listGeo"`
	ListDevice      []ListTarget         `json:"listDevice" query:"listDevice" form:"listDevice"`
}
