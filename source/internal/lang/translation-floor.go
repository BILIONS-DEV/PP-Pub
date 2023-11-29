package lang

type Floor struct {
	Title           TYPETranslation `json:"title"`
	AddFloor        TYPETranslation `json:"add_floor"`
	EditFloor       TYPETranslation `json:"edit_floor"`
	Main            TYPETranslation `json:"main"`
	PriceFloor      TYPETranslation `json:"price_floor"`
	Targeting       TYPETranslation `json:"targeting"`
	Run             TYPETranslation `json:"run"`
	Status          TYPETranslation `json:"status"`
	Name            TYPETranslation `json:"name"`
	NamePlaceHolder TYPETranslation `json:"name_placeholder"`
	Description     TYPETranslation `json:"description"`
	FloorValue      TYPETranslation `json:"floor_value"`
	Priority        TYPETranslation `json:"priority"`
	Target          TYPETranslation `json:"target"`
	Domains         TYPETranslation `json:"domains"`
	AllDomains      TYPETranslation `json:"all_domains"`
	Clear           TYPETranslation `json:"clear"`
	Format          TYPETranslation `json:"format"`
	AllFormat       TYPETranslation `json:"all_format"`
	Size            TYPETranslation `json:"size"`
	AllSizes        TYPETranslation `json:"all_sizes"`
	AdTag           TYPETranslation `json:"ad_tag"`
	AllAdTag        TYPETranslation `json:"all_ad_tag"`
	Geography       TYPETranslation `json:"geography"`
	AllGeographies  TYPETranslation `json:"all_geographies"`
	Devices         TYPETranslation `json:"devices"`
	AllDevices      TYPETranslation `json:"all_devices"`
	Button          TYPETranslation `json:"button"`
	SearchInventory TYPETranslation `json:"search_inventory"`
	SearchAdFormat  TYPETranslation `json:"search_ad_format"`
	SearchAdSize    TYPETranslation `json:"search_ad_size"`
	SearchAdTag     TYPETranslation `json:"search_ad_tag"`
	SearchGeography TYPETranslation `json:"search_geography"`
	SearchDevice    TYPETranslation `json:"search_device"`
	SearchFloor     TYPETranslation `json:"search_floor"`
}

type FloorError struct {
	Add      TYPETranslation `json:"add"`
	Edit     TYPETranslation `json:"edit"`
	Delete   TYPETranslation `json:"delete"`
	List     TYPETranslation `json:"list"`
	NotFound TYPETranslation `json:"not_found"`
	Target   TYPETranslation `json:"target"`
}
