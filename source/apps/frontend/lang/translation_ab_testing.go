package lang

type AbTesting struct {
	Title               TYPETranslation `json:"title"`
	AddAbTesting        TYPETranslation `json:"add_ab_testing"`
	EditAbTesting       TYPETranslation `json:"edit_ab_testing"`
	Main                TYPETranslation `json:"main"`
	Targeting           TYPETranslation `json:"targeting"`
	Run                 TYPETranslation `json:"run"`
	TestType            TYPETranslation `json:"test_type"`
	Bidder              TYPETranslation `json:"bidder"`
	SelectABidder       TYPETranslation `json:"select_a_bidder"`
	UserIdModule        TYPETranslation `json:"user_id_module"`
	SelectAUserIdModule TYPETranslation `json:"select_a_user_id_module"`
	TestGroupSize       TYPETranslation `json:"test_group_size"`
	Status              TYPETranslation `json:"status"`
	Name                TYPETranslation `json:"name"`
	NamePlaceHolder     TYPETranslation `json:"name_placeholder"`
	Description         TYPETranslation `json:"description"`
	Priority            TYPETranslation `json:"priority"`
	StartDate           TYPETranslation `json:"start_date"`
	EndDate             TYPETranslation `json:"end_date"`
	Target              TYPETranslation `json:"target"`
	Domains             TYPETranslation `json:"domains"`
	AllDomains          TYPETranslation `json:"all_domains"`
	Clear               TYPETranslation `json:"clear"`
	Format              TYPETranslation `json:"format"`
	AllFormat           TYPETranslation `json:"all_format"`
	Size                TYPETranslation `json:"size"`
	AllSizes            TYPETranslation `json:"all_sizes"`
	AdTag               TYPETranslation `json:"ad_tag"`
	AllAdTag            TYPETranslation `json:"all_ad_tag"`
	Geography           TYPETranslation `json:"geography"`
	AllGeographies      TYPETranslation `json:"all_geographies"`
	Devices             TYPETranslation `json:"devices"`
	AllDevices          TYPETranslation `json:"all_devices"`
	Button              TYPETranslation `json:"button"`
	SearchInventory     TYPETranslation `json:"search_inventory"`
	SearchAdFormat      TYPETranslation `json:"search_ad_format"`
	SearchAdSize        TYPETranslation `json:"search_ad_size"`
	SearchAdTag         TYPETranslation `json:"search_ad_tag"`
	SearchGeography     TYPETranslation `json:"search_geography"`
	SearchDevice        TYPETranslation `json:"search_device"`
	SearchAbTesting     TYPETranslation `json:"search_ab_testing"`
	DynamicPriceFloor   TYPETranslation `json:"dynamic_price_floor"`
	HardPriceFloor      TYPETranslation `json:"hard_price_floor"`
}

type AbTestingError struct {
	Add      TYPETranslation `json:"add"`
	Edit     TYPETranslation `json:"edit"`
	Delete   TYPETranslation `json:"delete"`
	List     TYPETranslation `json:"list"`
	NotFound TYPETranslation `json:"not_found"`
	Target   TYPETranslation `json:"target"`
}
