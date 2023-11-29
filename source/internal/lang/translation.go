package lang

import (
	"strings"
	"unicode"
)

type Translation struct {
	Submit      TYPETranslation `json:"submit"`
	Create      TYPETranslation `json:"create"`
	Edit        TYPETranslation `json:"edit"`
	Delete      TYPETranslation `json:"delete"`
	Remove      TYPETranslation `json:"remove"`
	Run         TYPETranslation `json:"run"`
	Setup       TYPETranslation `json:"setup"`
	Domain      TYPETranslation `json:"domain"`
	Domains     TYPETranslation `json:"domains"`
	Inventory   TYPETranslation `json:"inventory"`
	Inventories TYPETranslation `json:"inventories"`
	AdTag       TYPETranslation `json:"ad_tag"`
	Config      TYPETranslation `json:"config"`
	Consent     TYPETranslation `json:"consent"`
	User        TYPETranslation `json:"user"`
	UserID      TYPETranslation `json:"user_id"`
	Main        TYPETranslation `json:"main"`

	Pages         Pages           `json:"pages"`
	Errors        Errors          `json:"errors"`
	ErrorRequired TYPETranslation `json:"error_required"`
}

type Pages struct {
	GAM       GAM       `json:"gam"`
	Inventory Inventory `json:"inventory"`
	AdTag     AdTag     `json:"ad_tag"`
	Floor     Floor     `json:"floor"`
	LineItem  LineItem  `json:"line_item"`
	AdsTxt    AdsTxt    `json:"ads_txt"`
	Blocking  Blocking  `json:"blocking"`
	Playlist  Playlist  `json:"playlist"`
	Content   Content   `json:"content"`
	Bidder    Bidder    `json:"bidder"`
	Config    Config    `json:"config"`
	Template  Template  `json:"template"`
	User      User      `json:"user"`
	Support   Support   `json:"support"`
	Channels  Channels  `json:"channels"`
	Payment   Payment   `json:"payment"`
	Ai        Ai        `json:"ai"`
}

type Errors struct {
	RecordExist    TYPETranslation `json:"record_exist"`
	RecordNotFound TYPETranslation `json:"record_not_found"`
	InventoryError InventoryError  `json:"inventory_error"`
	AdTagError     AdTagError      `json:"ad_tag_error"`
	FloorError     FloorError      `json:"floor_error"`
	LineItemError  LineItemError   `json:"line_item_error"`
	AdsTxtError    AdsTxtError     `json:"ads_txt_error"`
	BlockingError  BlockingError   `json:"blocking_error"`
	PlaylistError  PlaylistError   `json:"playlist_error"`
	ContentError   ContentError    `json:"content_error"`
	ChannelsError  ChannelsError   `json:"channels_error"`
	BidderError    BidderError     `json:"bidder_error"`
	ConfigError    ConfigError     `json:"config_error"`
	TemplateError  TemplateError   `json:"template_error"`
	GamError       GamError        `json:"gam_error"`
	UserError      UserError       `json:"user_error"`
	AiError        AiError         `json:"floor_error"`
}

type TYPETranslation string

// ToLower In thường toàn bộ
func (t TYPETranslation) ToLower() string {
	return strings.ToLower(string(t))
}

// ToUpper In hoa toàn bộ
func (t TYPETranslation) ToUpper() string {
	return strings.ToUpper(string(t))
}

// ToUpperFirstCharacter In hoa chữ cái đầu tiên
func (t TYPETranslation) ToUpperFirstCharacter() string {
	if len(t) == 0 {
		return ""
	}
	tmp := []rune(t)
	tmp[0] = unicode.ToUpper(tmp[0])
	return string(tmp)
}

// Title In hoa chữ cái đầu tiên
func (t TYPETranslation) Title() string {
	return strings.Title(string(t))
}

// ToTitle In hoa chữ cái đầu tiên
func (t TYPETranslation) ToTitle() string {
	return strings.ToTitle(string(t))
}

func (t TYPETranslation) ToString() string {
	return string(t)
}
