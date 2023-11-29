package lang

type LineItem struct {
	Title                          TYPETranslation `json:"title"`
	Add                            TYPETranslation `json:"add"`
	Edit                           TYPETranslation `json:"edit"`
	Top                            TYPETranslation `json:"top"`
	SearchDemand                   TYPETranslation `json:"search_demand"`
	Run                            TYPETranslation `json:"run"`
	Main                           TYPETranslation `json:"main"`
	LineItemName                   TYPETranslation `json:"line_item_name"`
	Description                    TYPETranslation `json:"description"`
	Type                           TYPETranslation `json:"type"`
	Priority                       TYPETranslation `json:"priority"`
	Status                         TYPETranslation `json:"status"`
	DemandPartners                 TYPETranslation `json:"demand_partners"`
	AdsenseAdSlot                  TYPETranslation `json:"adsense_ad_slot"`
	SelectBidder                   TYPETranslation `json:"select_bidder"`
	SelectBidderPlaceholder        TYPETranslation `json:"select_bidder_placeholder"`
	SelectAccount                  TYPETranslation `json:"select_account"`
	SelectConnectionType           TYPETranslation `json:"select_connection_type"`
	ConnectionTypeDes              TYPETranslation `json:"connection_type_des"`
	SelectLineItemType             TYPETranslation `json:"select_line_item_type"`
	SelectAdsenseAdSlot            TYPETranslation `json:"select_adsense_ad_slot"`
	SelectAdsenseAdSlotPlaceholder TYPETranslation `json:"select_adsense_ad_slot_placeholder"`
	ItemAdsenseAdSlotPlaceholder   TYPETranslation `json:"item_adsense_ad_slot_placeholder"`
	Targeting                      TYPETranslation `json:"targeting"`
	StartDate                      TYPETranslation `json:"start_date"`
	EndDate                        TYPETranslation `json:"end_date"`
	Target                         TYPETranslation `json:"target"`
	Domains                        TYPETranslation `json:"domains"`
	AllDomains                     TYPETranslation `json:"all_domains"`
	Clear                          TYPETranslation `json:"clear"`
	Format                         TYPETranslation `json:"format"`
	AllFormat                      TYPETranslation `json:"all_format"`
	Size                           TYPETranslation `json:"size"`
	AllSizes                       TYPETranslation `json:"all_sizes"`
	AdTag                          TYPETranslation `json:"ad_tag"`
	AllAdTag                       TYPETranslation `json:"all_ad_tag"`
	Geography                      TYPETranslation `json:"geography"`
	AllGeographies                 TYPETranslation `json:"all_geographies"`
	Devices                        TYPETranslation `json:"devices"`
	AllDevices                     TYPETranslation `json:"all_devices"`
	Button                         TYPETranslation `json:"button"`
	SearchInventory                TYPETranslation `json:"search_inventory"`
	SearchAdFormat                 TYPETranslation `json:"search_ad_format"`
	SearchAdSize                   TYPETranslation `json:"search_ad_size"`
	SearchAdTag                    TYPETranslation `json:"search_ad_tag"`
	SearchGeography                TYPETranslation `json:"search_geography"`
	SearchDevice                   TYPETranslation `json:"search_device"`
	LinkedGam                      TYPETranslation `json:"linked_gam"`
}

type LineItemError struct {
	Add             TYPETranslation `json:"add"`
	Edit            TYPETranslation `json:"edit"`
	Delete          TYPETranslation `json:"delete"`
	List            TYPETranslation `json:"list"`
	NotFound        TYPETranslation `json:"not_found"`
	LineItemAccount TYPETranslation `json:"line_item_account"`
	BidderInfo      TYPETranslation `json:"bidder_info"`
	BidderParams    TYPETranslation `json:"bidder_params"`
	Target          TYPETranslation `json:"target"`
}
