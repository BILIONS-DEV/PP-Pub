package mysql

import (
	"gorm.io/gorm"
	"time"
)

type TableTemplate struct {
	Id                         int64                     `gorm:"column:id" json:"id"`
	UserId                     int64                     `gorm:"column:user_id" json:"user_id"`
	Name                       string                    `gorm:"column:name" json:"name"`
	Type                       TYPEAdType                `gorm:"column:type" json:"type"`
	IsDefault                  TypeOnOff                 `gorm:"column:is_default" json:"is_default"`
	PlayerLayout               TypePlayerLayout          `gorm:"column:player_layout" json:"player_layout"`
	Size                       TypeSize                  `gorm:"column:size" json:"size"`
	MaxWidth                   int                       `gorm:"column:max_width" json:"max_width"`
	Width                      int                       `gorm:"column:width" json:"width"`
	PlayMode                   TypePlayMode              `gorm:"column:play_mode" json:"play_mode"`
	CloseFloatingButtonDesktop TypeOnOff                 `gorm:"column:close_floating_button_desktop" json:"close_floating_button_desktop"`
	FloatOnBottom              TypeOnOff                 `gorm:"column:float_on_bottom" json:"float_on_bottom"`
	FloatingOnView             TypeOnOff                 `gorm:"column:floating_on_view" json:"floating_on_view"`
	FloatingOnImpression       TypeOnOff                 `gorm:"column:floating_on_impression" json:"floating_on_impression"`
	FloatingOnAdFetched        TypeOnOff                 `gorm:"column:floating_on_ad_fetched" json:"floating_on_ad_fetched"`
	FloatingWidth              int                       `gorm:"column:floating_width" json:"floating_width"`
	FloatingPositionDesktop    TypePositionDesktop       `gorm:"column:floating_position_desktop" json:"floating_position_desktop"`
	MarginTopDesktop           int                       `gorm:"column:margin_top_desktop" json:"margin_top_desktop"`
	MarginBottomDesktop        int                       `gorm:"column:margin_bottom_desktop" json:"margin_bottom_desktop"`
	MarginRightDesktop         int                       `gorm:"column:margin_right_desktop" json:"margin_right_desktop"`
	MarginLeftDesktop          int                       `gorm:"column:margin_left_desktop" json:"margin_left_desktop"`
	ColumnsNumber              int                       `gorm:"column:columns_number" json:"columns_number"`
	ColumnsPosition            int                       `gorm:"column:columns_position" json:"columns_position"`
	CloseFloatingButtonMobile  TypeOnOff                 `gorm:"column:close_floating_button_mobile" json:"close_floating_button_mobile"`
	FloatOnBottomMobile        TypeOnOff                 `gorm:"column:float_on_bottom_mobile" json:"float_on_bottom_mobile"`
	FloatingOnViewMobile       TypeOnOff                 `gorm:"column:floating_on_view_mobile" json:"floating_on_view_mobile"`
	FloatingOnAdFetchedMobile  TypeOnOff                 `gorm:"column:floating_on_ad_fetched_mobile" json:"floating_on_ad_fetched_mobile"`
	FloatingWidthMobile        int                       `gorm:"column:floating_width_mobile" json:"floating_width_mobile"`
	FloatingPositionMobile     TypePositionMobile        `gorm:"column:floating_position_mobile" json:"floating_position_mobile"`
	MarginBottomMobile         int                       `gorm:"column:margin_bottom_mobile" json:"margin_bottom_mobile"`
	MarginRightMobile          int                       `gorm:"column:margin_right_mobile" json:"margin_right_mobile"`
	MarginLeftMobile           int                       `gorm:"column:margin_left_mobile" json:"margin_left_mobile"`
	MainTitle                  TypeOnOff                 `gorm:"column:main_title" json:"main_title"`
	MainTitleText              string                    `gorm:"column:main_title_text" json:"main_title_text"`
	SubTitle                   TypeOnOff                 `gorm:"column:sub_title" json:"sub_title"`
	SubTitleText               string                    `gorm:"column:sub_title_text" json:"sub_title_text"`
	ActionButton               TypeOnOff                 `gorm:"column:action_button" json:"action_button"`
	ActionButtonText           string                    `gorm:"column:action_button_text" json:"action_button_text"`
	TitleEnable                TypeOnOff                 `gorm:"column:title_enable" json:"title_enable"`
	DescriptionEnable          TypeOnOff                 `gorm:"column:description_enable" json:"description_enable"`
	ControlColor               string                    `gorm:"column:control_color" json:"control_color"`
	ThemeColor                 string                    `gorm:"column:theme_color" json:"theme_color"`
	BackgroundColor            string                    `gorm:"column:background_color" json:"background_color"`
	MainTitleBackgroundColor   string                    `gorm:"column:main_title_background_color" json:"main_title_background_color"`
	TitleColor                 string                    `gorm:"column:title_color" json:"title_color"`
	MainTitleColor             string                    `gorm:"column:main_title_color" json:"main_title_color"`
	TitleBackgroundColor       string                    `gorm:"column:title_background_color" json:"title_background_color"`
	TitleHoverBackgroundColor  string                    `gorm:"column:title_hover_background_color" json:"title_hover_background_color"`
	ActionButtonColor          string                    `gorm:"column:action_button_color" json:"action_button_color"`
	DescriptionColor           string                    `gorm:"column:description_color" json:"description_color"`
	AdvertiserColor            string                    `gorm:"column:advertiser_color" json:"advertiser_color"`
	DefaultSoundMode           TypeOnOff                 `gorm:"column:default_sound_mode" json:"default_sound_mode"`
	FullscreenButton           TypeOnOff                 `gorm:"column:fullscreen_button" json:"fullscreen_button"`
	NextPrevArrowsButton       TypeOnOff                 `gorm:"column:next_prev_arrows_button" json:"next_prev_arrows_button"`
	NextPrevTime               TypeOnOff                 `gorm:"column:next_prev_time" json:"next_prev_time"`
	VideoConfig                TypeOnOff                 `gorm:"column:video_config" json:"video_config"`
	ShowStats                  TypeOnOff                 `gorm:"column:show_stats" json:"show_stats"`
	ShareButton                TypeOnOff                 `gorm:"column:share_button" json:"share_button"`
	CustomLogo                 TypeOnOff                 `gorm:"column:custom_logo" json:"custom_logo"`
	Link                       string                    `gorm:"column:link" json:"link"`
	ClickThrough               string                    `gorm:"column:click_through" json:"click_through"`
	WaitForAd                  TypeOnOff                 `gorm:"column:wait_for_ad" json:"wait_for_ad"`
	AdvertisementScenario      TypeAdvertisementScenario `gorm:"column:advertisement_scenario" json:"advertisement_scenario"`
	VastRetry                  int                       `gorm:"column:vast_retry" json:"vast_retry"`
	Delay                      int                       `gorm:"column:delay" json:"delay"`
	AutoSkip                   TypeOnOff                 `gorm:"column:auto_skip" json:"auto_skip"`
	TimeToSkip                 int                       `gorm:"column:time_to_skip" json:"time_to_skip"`
	ShowAutoSkipButton         int                       `gorm:"column:show_auto_skip_button" json:"show_auto_skip_button"`
	NumberOfPreRollAds         int                       `gorm:"column:number_of_pre_roll_ads" json:"number_of_pre_roll_ads"`
	FloatingOnDesktop          TypeOnOff                 `gorm:"column:floating_on_desktop" json:"floating_on_desktop"`
	FloatingOnMobile           TypeOnOff                 `gorm:"column:floating_on_mobile" json:"floating_on_mobile"`
	PoweredBy                  TypeOnOff                 `gorm:"column:powered_by" json:"powered_by"`
	EnableLogo                 TypeOnOff                 `gorm:"column:enable_logo" json:"enable_logo"`
	PreRoll                    TypeOnOff                 `gorm:"column:pre_roll" json:"pre_roll"`
	MidRoll                    TypeOnOff                 `gorm:"column:mid_roll" json:"mid_roll"`
	PostRoll                   TypeOnOff                 `gorm:"column:post_roll" json:"post_roll"`
	AutoStart                  TYPEAutoStart             `gorm:"column:auto_start" json:"auto_start"`
	ShowControls               TypeOnOff                 `gorm:"column:show_controls" json:"show_controls"`
	PubPowerLogo               TypeOnOff                 `gorm:"column:pubpower_logo" json:"pubpower_logo"`
	AdStyle                    TYPETemplateAdStyle       `gorm:"column:ad_style" json:"ad_style"`
	AdSize                     int64                     `gorm:"column:ad_size" json:"ad_size"`
	Mode                       string                    `gorm:"column:mode" json:"mode"`
	Template                   string                    `gorm:"column:template" json:"template"`
	Rows                       int                       `gorm:"column:rows" json:"rows"`
	Columns                    int                       `gorm:"column:columns" json:"columns"`
	SponsoredBrand             TYPESponsoredBrand        `gorm:"column:sponsored_brand" json:"sponsored_brand"`
	CreatedAt                  time.Time                 `gorm:"column:created_at" json:"created_at"`
	UpdatedAt                  time.Time                 `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt                  gorm.DeletedAt            `gorm:"column:deleted_at" json:"deleted_at"`
	SizeModel                  TableAdSize               `gorm:"-" json:"size_model"`
}

func (TableTemplate) TableName() string {
	return Tables.Template
}

type TypeAdvertisementScenario int

const (
	TypeAdvertisementScenarioDefault TypeAdvertisementScenario = iota + 1
	TypeAdvertisementScenarioAuto
	TypeAdvertisementScenarioManual
)

func (t TypeAdvertisementScenario) String() string {
	switch t {
	case TypeAdvertisementScenarioDefault:
		return "Default"
	case TypeAdvertisementScenarioAuto:
		return "Auto"
	case TypeAdvertisementScenarioManual:
		return "Manual"
	}
	return ""
}

type TypeSize int

const (
	TypeSizeResponsive TypeSize = iota + 1
	TypeSizeFixed
)

func (t TypeSize) String() string {
	switch t {
	case TypeSizeResponsive:
		return "Responsive"
	case TypeSizeFixed:
		return "Fixed"
	}
	return ""
}

type TypeOnOff int

const (
	TypeOn TypeOnOff = iota + 1
	TypeOff
)

func (t TypeOnOff) String() string {
	switch t {
	case TypeOn:
		return "On"
	case TypeOff:
		return "Off"
	}
	return ""
}

func (t TypeOnOff) Boolean() bool {
	switch t {
	case TypeOn:
		return true
	case TypeOff:
		return false
	}
	return false
}

type TypePositionDesktop int

const (
	TypePositionDesktopBottomRight TypePositionDesktop = iota + 1
	TypePositionDesktopBottomLeft
	TypePositionDesktopTopRight
	TypePositionDesktopTopLeft
)

func (t TypePositionDesktop) String() string {
	switch t {
	case TypePositionDesktopBottomRight:
		return "Bottom Right"
	case TypePositionDesktopBottomLeft:
		return "Bottom Left"
	case TypePositionDesktopTopRight:
		return "Top Right"
	case TypePositionDesktopTopLeft:
		return "Top Left"
	}
	return ""
}

func (t TypePositionDesktop) Int() int {
	switch t {
	case TypePositionDesktopBottomRight:
		return 1
	case TypePositionDesktopBottomLeft:
		return 2
	case TypePositionDesktopTopRight:
		return 3
	case TypePositionDesktopTopLeft:
		return 4
	}
	return 0
}

type TypePositionMobile int

const (
	TypePositionMobileBottomRight TypePositionMobile = iota + 1
	TypePositionMobileBottomLeft
)

func (t TypePositionMobile) String() string {
	switch t {
	case TypePositionMobileBottomRight:
		return "Bottom Right"
	case TypePositionMobileBottomLeft:
		return "Bottom Left"
	}
	return ""
}

func (t TypePositionMobile) Int() int {
	switch t {
	case TypePositionMobileBottomRight:
		return 1
	case TypePositionMobileBottomLeft:
		return 2
	}
	return 0
}

type TypePlayerLayout int

const (
	TypePlayerLayoutBasic TypePlayerLayout = iota + 1
	TypePlayerLayoutClassic
	TypePlayerLayoutSmall
	TypePlayerLayoutInContent
	TypePlayerLayoutSide
	TypePlayerLayoutInContentThumb
	TypePlayerLayoutInContentText
	TypePlayerLayoutTopArticle
	TypePlayerLayoutNone
)

func (t TypePlayerLayout) String() string {
	switch t {
	case TypePlayerLayoutBasic:
		return "Basic"
	case TypePlayerLayoutClassic:
		return "Classic"
	case TypePlayerLayoutSmall:
		return "Small"
	case TypePlayerLayoutInContent:
		return "In Content"
	case TypePlayerLayoutSide:
		return "Side"
	case TypePlayerLayoutInContentThumb:
		return "In Content Thumb"
	case TypePlayerLayoutInContentText:
		return "In Content Text"
	case TypePlayerLayoutTopArticle:
		return "Top Article"
	case TypePlayerLayoutNone:
		return "None"
	}
	return ""
}

func (t TypePlayerLayout) Int() int {
	switch t {
	case TypePlayerLayoutBasic:
		return 1
	case TypePlayerLayoutClassic:
		return 2
	case TypePlayerLayoutSmall:
		return 3
	case TypePlayerLayoutInContent:
		return 4
	case TypePlayerLayoutSide:
		return 5
	case TypePlayerLayoutInContentThumb:
		return 6
	case TypePlayerLayoutInContentText:
		return 7
	case TypePlayerLayoutTopArticle:
		return 8
	case TypePlayerLayoutNone:
		return 9
	}
	return 0
}

type TypePlayMode int

const (
	TypePlayModeInline TypePlayMode = iota + 1
	TypePlayModeFloating
)

func (t TypePlayMode) String() string {
	switch t {
	case TypePlayModeInline:
		return "Inline"
	case TypePlayModeFloating:
		return "Floating"
	}
	return ""
}

func (t TypePlayMode) Int() int {
	switch t {
	case TypePlayModeInline:
		return 1
	case TypePlayModeFloating:
		return 2
	}
	return 0
}

type TYPEAutoStart int

const (
	TYPEAutoStartOn TYPEAutoStart = iota + 1
	TYPEAutoStartOff
	TYPEAutoStartWhenIsView
	TYPEAutoStartWhenIsAfterFirst
)

func (t TYPEAutoStart) String() string {
	switch t {
	case TYPEAutoStartOn:
		return "On"
	case TYPEAutoStartOff:
		return "Off"
	case TYPEAutoStartWhenIsView:
		return "When player is in view"
	case TYPEAutoStartWhenIsAfterFirst:
		return "After first ad is finished"
	}
	return ""
}

func (t TYPEAutoStart) Code() string {
	switch t {
	case TYPEAutoStartOn:
		return "on"
	case TYPEAutoStartOff:
		return "off"
	case TYPEAutoStartWhenIsView:
		return ""
	case TYPEAutoStartWhenIsAfterFirst:
		return "afterAd"
	}
	return ""
}

func (t TYPEAutoStart) Int() int {
	switch t {
	case TYPEAutoStartOn:
		return 1
	case TYPEAutoStartOff:
		return 2
	case TYPEAutoStartWhenIsView:
		return 3
	case TYPEAutoStartWhenIsAfterFirst:
		return 4
	}
	return 0
}

type TYPETemplateAdStyle string

const (
	TYPETemplateAdStyleMultiple = "multiple"
	TYPETemplateAdStyleSingle   = "single"
)

type TYPESponsoredBrand string

const (
	TYPESponsoredBrandTrue  = "true"
	TYPESponsoredBrandFalse = "false"
)

func (t TYPESponsoredBrand) Boolean() bool {
	switch t {
	case TYPESponsoredBrandTrue:
		return true
	case TYPESponsoredBrandFalse:
		return false
	}
	return false
}
