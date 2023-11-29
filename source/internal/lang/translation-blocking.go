package lang

type Blocking struct {
	Title            TYPETranslation `json:"title"`
	Add              TYPETranslation `json:"add"`
	Edit             TYPETranslation `json:"edit"`
	SearchBlocking   TYPETranslation `json:"search_blocking"`
	Run              TYPETranslation `json:"run"`
	Top              TYPETranslation `json:"top"`
	AddRule          TYPETranslation `json:"add_rule"`
	Name             TYPETranslation `json:"name"`
	NamePlaceholder  TYPETranslation `json:"name_placeholder"`
	Inventory        TYPETranslation `json:"inventory"`
	AllInventory     TYPETranslation `json:"all_inventory"`
	Restrictions     TYPETranslation `json:"restrictions"`
	AdvertiserDomain TYPETranslation `json:"advertiser_domain"`
	Clear            TYPETranslation `json:"clear"`
	Placeholder      TYPETranslation `json:"placeholder"`
	AddDomain        TYPETranslation `json:"add_domain"`
	Button           TYPETranslation `json:"button"`
}

type BlockingError struct {
	Add      TYPETranslation `json:"add"`
	Edit     TYPETranslation `json:"edit"`
	Delete   TYPETranslation `json:"delete"`
	List     TYPETranslation `json:"list"`
	NotFound TYPETranslation `json:"not_found"`
}
