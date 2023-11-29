package payload

import "source/core/technology/mysql"

type BidderCreate struct {
	BidAdjustment int                          `json:"bid_adjustment" form:"bid_adjustment"`
	RefreshAd     int                          `json:"refresh_ad" form:"refresh_ad"`
	AdFormat      []int64                      `json:"ad_format" form:"ad_format"`
	GeoIpOption   mysql.TYPEOption             `json:"geo_ip_option" form:"geo_ip_option"`
	GeoIpList     []string                     `json:"geo_ip_list" form:"geo_ip_list"`
	AdsTxt        string                       `json:"ads_txt" form:"ads_txt"`
	Status        mysql.TYPEBidderAssignStatus `json:"status" form:"status"`
}
