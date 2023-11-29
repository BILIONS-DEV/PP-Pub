package lang

type AdTag struct {
	Add  TYPETranslation `json:"add"`
	Edit TYPETranslation `json:"edit"`

	Status            TYPETranslation `json:"status"`
	Name              TYPETranslation `json:"name"`
	NameDesc          TYPETranslation `json:"name_desc"`
	NamePlaceholder   TYPETranslation `json:"name_placeholder"`
	AdType            TYPETranslation `json:"ad_type"`
	AdTypeDesc        TYPETranslation `json:"ad_type_desc"`
	AdTypePlaceholder TYPETranslation `json:"ad_type_placeholder"`

	GAMAdUnit                   TYPETranslation `json:"gam_ad_unit"`
	GAMAdUnitDesc               TYPETranslation `json:"gam_ad_unit_desc"`
	GAMAdUnitPlaceholder        TYPETranslation `json:"gam_ad_unit_placeholder"`
	PrimaryAdSize               TYPETranslation `json:"primary_ad_size"`
	SizeOnMobile                TYPETranslation `json:"size_on_mobile"`
	SizeOnMobilePlaceholder     TYPETranslation `json:"size_on_mobile_placeholder"`
	PrimaryAdSizeDesc           TYPETranslation `json:"primary_ad_size_desc"`
	PrimaryAdSizePlaceholder    TYPETranslation `json:"primary_ad_size_placeholder"`
	AdditionalAdSize            TYPETranslation `json:"additional_ad_size"`
	AdditionalAdSizeDesc        TYPETranslation `json:"additional_ad_size_desc"`
	AdditionalAdSizePlaceholder TYPETranslation `json:"additional_ad_size_placeholder"`
	BidOutstream                TYPETranslation `json:"bid_outstream"`
	BidOutstreamDesc            TYPETranslation `json:"bid_outstream_desc"`
	Passback                    TYPETranslation `json:"passback"`
	PassbackDesc                TYPETranslation `json:"passback_desc"`
	PassbackPlaceholder         TYPETranslation `json:"passback_placeholder"`
	DesktopConfig               TYPETranslation `json:"desktop_config"`
	MobileConfig                TYPETranslation `json:"mobile_config"`
}

type AdTagError struct {
	Add      TYPETranslation `json:"add"`
	Edit     TYPETranslation `json:"edit"`
	Delete   TYPETranslation `json:"delete"`
	List     TYPETranslation `json:"list"`
	NotFound TYPETranslation `json:"not_found"`
}
