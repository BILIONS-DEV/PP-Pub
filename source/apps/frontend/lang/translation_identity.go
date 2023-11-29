package lang

type Identity struct {
	Title            TYPETranslation `json:"title"`
	AddIdentity      TYPETranslation `json:"add_identity"`
	EditIdentity     TYPETranslation `json:"edit_identity"`
	Main             TYPETranslation `json:"main"`
	Targeting        TYPETranslation `json:"targeting"`
	Run              TYPETranslation `json:"run"`
	Status           TYPETranslation `json:"status"`
	Name             TYPETranslation `json:"name"`
	NamePlaceHolder  TYPETranslation `json:"name_placeholder"`
	Description      TYPETranslation `json:"description"`
	Priority         TYPETranslation `json:"priority"`
	PriorityDesc     TYPETranslation `json:"priority_desc"`
	Target           TYPETranslation `json:"target"`
	Domains          TYPETranslation `json:"domains"`
	AllDomains       TYPETranslation `json:"all_domains"`
	Clear            TYPETranslation `json:"clear"`
	Button           TYPETranslation `json:"button"`
	SearchInventory  TYPETranslation `json:"search_inventory"`
	SearchIdentity   TYPETranslation `json:"search_identity"`
	SyncDelay        TYPETranslation `json:"sync_delay"`
	SyncDelayDesc    TYPETranslation `json:"sync_delay_desc"`
	AuctionDelay     TYPETranslation `json:"auction_delay"`
	AuctionDelayDesc TYPETranslation `json:"auction_delay_desc"`
}

type IdentityError struct {
	Add            TYPETranslation `json:"add"`
	Edit           TYPETranslation `json:"edit"`
	Delete         TYPETranslation `json:"delete"`
	List           TYPETranslation `json:"list"`
	NotFound       TYPETranslation `json:"not_found"`
	Target         TYPETranslation `json:"target"`
	ModuleUserId   TYPETranslation `json:"module_user_id"`
	ChooseDomain   TYPETranslation `json:"choose_domain"`
	TargetDomain   TYPETranslation `json:"target_domain"`
	DomainTargeted TYPETranslation `json:"domain_targeted"`
}
