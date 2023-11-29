package lang

type Bidder struct {
	Title              TYPETranslation `json:"title"`
	Add                TYPETranslation `json:"add"`
	Edit               TYPETranslation `json:"edit"`
	SearchBidder       TYPETranslation `json:"search_bidder"`
	Run                TYPETranslation `json:"run"`
	Top                TYPETranslation `json:"top"`
	Main               TYPETranslation `json:"main"`
	SelectBidder       TYPETranslation `json:"select_bidder"`
	SelectABidder      TYPETranslation `json:"select_a_bidder"`
	MediaTypes         TYPETranslation `json:"media_types"`
	SelectMedia        TYPETranslation `json:"select_media"`
	BidAdjustment      TYPETranslation `json:"bid_adjustment"`
	BidContent         TYPETranslation `json:"bid_content"`
	AdsTxt             TYPETranslation `json:"ads_txt"`
	Params             TYPETranslation `json:"params"`
	AddParam           TYPETranslation `json:"add_param"`
	DisplayName        TYPETranslation `json:"display_name"`
	DisplayNameContent TYPETranslation `json:"display_name_content"`
	PubId              TYPETranslation `json:"pub_id"`
	PubIdDesGoogle     TYPETranslation `json:"pub_id_des_google"`
	PubIdDesAmz        TYPETranslation `json:"pub_id_des_amz"`
	LinkGam            TYPETranslation `json:"link_gam"`
	SelectAccountType  TYPETranslation `json:"select_account_type"`
	SelectGam          TYPETranslation `json:"select_gam"`
	Button             TYPETranslation `json:"button"`
	BidderAlias        TYPETranslation `json:"bidder_alias"`
	IsDefault          TYPETranslation `json:"is_default"`
}

type BidderError struct {
	Add         TYPETranslation `json:"add"`
	Edit        TYPETranslation `json:"edit"`
	Delete      TYPETranslation `json:"delete"`
	List        TYPETranslation `json:"list"`
	NotFound    TYPETranslation `json:"not_found"`
	MediaBidder TYPETranslation `json:"media_bidder"`
}
