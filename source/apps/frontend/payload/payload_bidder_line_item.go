package payload

import (
	"source/core/technology/mysql"
	"source/pkg/datatable"
)

type LineItemIndex struct {
	LineItemFilterPostData
	OrderColumn int      `query:"order_column" json:"order_column" form:"order_column"`
	OrderDir    string   `query:"order_dir" json:"order_dir" form:"order_dir"`
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

type LineItemFilterPayload struct {
	datatable.Request
	PostData *LineItemFilterPostData `query:"postData"`
}

type LineItemFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Status      interface{} `query:"f_status[]" json:"f_status" form:"f_status[]"`
	Type        interface{} `query:"f_type[]" json:"f_type" form:"f_type[]"`
	Domain      interface{} `query:"f_domain[]" json:"f_domain" form:"f_domain[]"`
	AdFormat    interface{} `query:"f_ad_format[]" json:"f_ad_format" form:"f_ad_format[]"`
	AdSize      interface{} `query:"f_ad_size[]" json:"f_ad_size" form:"f_ad_size[]"`
	AdTag       interface{} `query:"f_ad_tag[]" json:"f_ad_tag" form:"f_ad_tag[]"`
	Device      interface{} `query:"f_device[]" json:"f_device" form:"f_device[]"`
	Country     interface{} `query:"f_geo[]" json:"f_geo" form:"f_geo[]"`
	Page        int         `query:"page" json:"page" form:"page"`
	Limit       int         `query:"limit" json:"limit" form:"limit"`
	Start       int         `query:"start" json:"start" form:"start"`
	Length      int         `query:"length" json:"length" form:"length"`
}

type ListTarget struct {
	Id   int64  `gorm:"column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

type LineItemAdd struct {
	Id              int64                     `json:"id" form:"id"`
	Name            string                    `json:"name" form:"name"`
	Description     string                    `json:"description" form:"description"`
	ServerType      mysql.TYPEServerType      `json:"server_type" form:"server_type"`
	SelectAccount   []int64                   `json:"select_account" form:"select_account"`
	ConnectionType  mysql.TYPEConnectionType  `json:"connection_type" form:"connection_type"`
	GamLineItemType mysql.TYPEGamLineItemType `json:"line_item_type" form:"line_item_type"`
	LinkedGam       int64                     `json:"linked_gam" form:"linked_gam"`
	Status          string                    `json:"status" form:"status"`
	BidderParams    []BidderInfo              `json:"bidder_params" form:"bidder_params"`
	AdsenseAdSlots  []AdsenseAdSlot           `json:"adsense_ad_slots" form:"adsense_ad_slots"`
	BidderRate      int                       `json:"rate" form:"rate"`
	BidderVastUrl   string                    `json:"vast_url" form:"vast_url"`
	BidderAdTag     string                    `json:"ad_tag" form:"ad_tag"`
	Priority        int                       `json:"priority" form:"priority"`
	StartDate       string                    `json:"start_date" form:"start_date"`
	EndDate         string                    `json:"end_date" form:"end_date"`
	ListAdFormat    []ListTarget              `json:"listAdFormat" query:"listAdFormat" form:"listAdFormat"`
	ListAdSize      []ListTarget              `json:"listAdSize" query:"listAdSize" form:"listAdSize"`
	ListAdTag       []ListTarget              `json:"listAdtag" query:"listAdtag" form:"listAdtag"`
	ListAdInventory []ListTarget              `json:"listInventory" query:"listInventory" form:"listInventory"`
	ListGeo         []ListTarget              `json:"listGeo" query:"listGeo" form:"listGeo"`
	ListDevice      []ListTarget              `json:"listDevice" query:"listDevice" form:"listDevice"`
}

type LineItemEdit struct {
	Id              int64                `json:"id" form:"id"`
	Name            string               `json:"name" form:"name"`
	Description     string               `json:"description" form:"description"`
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

type FilterTarget struct {
	Inventory       []int64 `json:"inventory"`
	Format          []int64 `json:"format"`
	Size            []int64 `json:"size"`
	Language        []int64 `json:"language"`
	ExcludeLanguage []int64 `json:"exclude_language"`
	Channels        []int64 `json:"channels"`
	ExcludeChannels []int64 `json:"exclude_channels"`
	Keywords        []int64 `json:"keywords"`
	ExcludeKeywords []int64 `json:"exclude_keywords"`
}

type LineItem struct {
	BidderId  int64  `json:"bidder"`
	Placement string `json:"placement"`
	Publisher string `json:"publisher"`
}

type AdsenseAdSlot struct {
	Size     string `json:"size"`
	AdSlotId string `json:"ad_slot_id"`
}

type BidderInfo struct {
	BidderId       int64                  `json:"id"`
	BidderName     string                 `json:"name"`
	ConfigType     mysql.TYPEConfigType   `json:"config_type"`
	BidderType     mysql.TYPEBidderType   `json:"bidder_type"`
	BidderIndex    int                    `json:"bidder_index"`
	BidderTemplate int64                  `json:"bidder_template"`
	Params         map[string]ParamDetail `json:"params"`
	BidderParams   []BidderParam          `json:"-"`
	Link           string                 `json:"-"`
}

type ParamDetail struct {
	Value      string          `json:"value"`
	Type       string          `json:"type"`
	IsAddParam mysql.TypeOnOff `json:"addParam"`
}

type BidderParam struct {
	Param      mysql.TableLineItemBidderParams
	IsRequired bool
}
