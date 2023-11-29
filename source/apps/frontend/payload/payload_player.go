package payload

import (
	"encoding/json"
	"source/core/technology/mysql"
	"source/pkg/datatable"

	"golang.org/x/net/html"
)

type TemplateCreate struct {
	Id                         int64                           `json:"id"`
	Name                       string                          `json:"name"`
	Type                       mysql.TYPEAdType                `json:"type"`
	IsDefault                  string                          `json:"is_default"`
	PlayerLayout               mysql.TypePlayerLayout          `json:"player_layout"`
	Size                       mysql.TypeSize                  `json:"size"`
	MaxWidth                   int                             `json:"max_width"`
	Width                      int                             `json:"width"`
	PlayMode                   mysql.TypePlayMode              `json:"play_mode"`
	CloseFloatingButtonDesktop mysql.TypeOnOff                 `json:"close_floating_button_desktop"`
	FloatOnBottom              mysql.TypeOnOff                 `json:"float_on_bottom"`
	FloatingOnView             mysql.TypeOnOff                 `json:"floating_on_view"`
	FloatingOnImpression       mysql.TypeOnOff                 `json:"floating_on_impression"`
	FloatingOnAdFetched        mysql.TypeOnOff                 `json:"floating_on_ad_fetched"`
	FloatingWidth              int                             `json:"floating_width"`
	FloatingPositionDesktop    mysql.TypePositionDesktop       `json:"floating_position_desktop"`
	MarginTopDesktop           int                             `json:"margin_top_desktop"`
	MarginBottomDesktop        int                             `json:"margin_bottom_desktop"`
	MarginRightDesktop         int                             `json:"margin_right_desktop"`
	MarginLeftDesktop          int                             `json:"margin_left_desktop"`
	CloseFloatingButtonMobile  mysql.TypeOnOff                 `json:"close_floating_button_mobile"`
	FloatOnBottomMobile        mysql.TypeOnOff                 `json:"float_on_bottom_mobile"`
	FloatingOnViewMobile       mysql.TypeOnOff                 `json:"floating_on_view_mobile"`
	FloatingOnAdFetchedMobile  mysql.TypeOnOff                 `json:"floating_on_ad_fetched_mobile"`
	FloatingWidthMobile        int                             `json:"floating_width_mobile"`
	FloatingPositionMobile     mysql.TypePositionMobile        `json:"floating_position_mobile"`
	MarginBottomMobile         int                             `json:"margin_bottom_mobile"`
	MarginRightMobile          int                             `json:"margin_right_mobile"`
	MarginLeftMobile           int                             `json:"margin_left_mobile"`
	ColumnsNumber              int                             `json:"columns_number"`
	ColumnsPosition            int                             `json:"columns_position"`
	MainTitle                  mysql.TypeOnOff                 `json:"main_title"`
	MainTitleText              string                          `json:"main_title_text"`
	SubTitle                   mysql.TypeOnOff                 `json:"sub_title"`
	SubTitleText               string                          `json:"sub_title_text"`
	ActionButton               mysql.TypeOnOff                 `json:"action_button"`
	ActionButtonText           string                          `json:"action_button_text"`
	TitleEnable                mysql.TypeOnOff                 `json:"title_enable"`
	DescriptionEnable          mysql.TypeOnOff                 `json:"description_enable"`
	ControlColor               string                          `json:"control_color"`
	ThemeColor                 string                          `json:"theme_color"`
	BackgroundColor            string                          `json:"background_color"`
	MainTitleBackgroundColor   string                          `json:"main_title_background_color"`
	TitleColor                 string                          `json:"title_color"`
	MainTitleColor             string                          `json:"main_title_color"`
	TitleBackgroundColor       string                          `json:"title_background_color"`
	TitleHoverBackgroundColor  string                          `json:"title_hover_background_color"`
	ActionButtonColor          string                          `json:"action_button_color"`
	DescriptionColor           string                          `json:"description_color"`
	AdvertiserColor            string                          `json:"advertiser_color"`
	DefaultSoundMode           mysql.TypeOnOff                 `json:"default_sound_mode"`
	FullscreenButton           mysql.TypeOnOff                 `json:"fullscreen_button"`
	NextPrevArrowsButton       mysql.TypeOnOff                 `json:"next_prev_arrows_button"`
	NextPrevTime               mysql.TypeOnOff                 `json:"next_prev_time"`
	VideoConfig                mysql.TypeOnOff                 `json:"video_config"`
	ShowStats                  mysql.TypeOnOff                 `json:"show_stats"`
	ShareButton                mysql.TypeOnOff                 `json:"share_button"`
	CustomLogo                 mysql.TypeOnOff                 `json:"custom_logo"`
	Link                       string                          `json:"link"`
	ClickThrough               string                          `json:"click_through"`
	WaitForAd                  mysql.TypeOnOff                 `json:"wait_for_ad"`
	VastRetry                  int                             `json:"vast_retry"`
	AutoSkip                   mysql.TypeOnOff                 `json:"auto_skip"`
	TimeToSkip                 int                             `json:"time_to_skip"`
	ShowAutoSkipButton         int                             `json:"show_auto_skip_button"`
	FloatingOnDesktop          mysql.TypeOnOff                 `json:"floating_on_desktop"`
	FloatingOnMobile           mysql.TypeOnOff                 `json:"floating_on_mobile"`
	PoweredBy                  mysql.TypeOnOff                 `json:"powered_by"`
	EnableLogo                 mysql.TypeOnOff                 `json:"enable_logo"`
	MainTitleTopArticle        string                          `json:"main_title_top_article"`
	EnableLogoTopArticle       mysql.TypeOnOff                 `json:"enable_logo_top_article"`
	CustomLogoTopArticle       mysql.TypeOnOff                 `json:"custom_logo_top_article"`
	LinkTopArticle             string                          `json:"link_top_article"`
	ClickThroughTopArticle     string                          `json:"click_through_top_article"`
	PoweredByTopArticle        mysql.TypeOnOff                 `json:"powered_by_top_article"`
	Delay                      int                             `json:"delay"`
	AdvertisementScenario      mysql.TypeAdvertisementScenario `json:"advertisement_scenario"`
	PreRoll                    mysql.TypeOnOff                 `json:"pre_roll"`
	MidRoll                    mysql.TypeOnOff                 `json:"mid_roll"`
	PostRoll                   mysql.TypeOnOff                 `json:"post_roll"`
	AutoStart                  mysql.TYPEAutoStart             `json:"auto_start"`
	ShowControls               mysql.TypeOnOff                 `json:"show_controls"`
	NumberOfPreRollAds         int                             `json:"number_of_pre_roll_ads"`
	PubPowerLogo               mysql.TypeOnOff                 `json:"pubpower_logo"`
	AdStyle                    string                          `json:"ad_style"`
	AdSize                     int64                           `json:"ad_size"`
	Mode                       string                          `json:"mode"`
	Rows                       int                             `json:"rows"`
	Columns                    int                             `json:"columns"`
}

type PlayerIndex struct {
	PlayerFilterPostData
	QuerySearch string   `query:"f_q" json:"f_q" form:"f_q"`
	Type        []string `query:"f_type" form:"f_type" json:"f_type"`
	Size        []string `query:"f_size" form:"f_size" json:"f_size"`
}

type PlayerFilterPayload struct {
	datatable.Request
	PostData *PlayerFilterPostData `query:"postData"`
}

type PlayerFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Type        interface{} `query:"f_type[]" json:"f_type" form:"f_type[]"`
	Size        interface{} `query:"f_size[]" json:"f_size" form:"f_size[]"`
	Page        int         `query:"page" json:"page" form:"page"`
	Limit       int         `query:"limit" json:"limit" form:"limit"`
	Start       int         `query:"start" json:"start" form:"start"`
	Length      int         `query:"length" json:"length" form:"length"`
}

type PlayerConfig struct {
	Contents   []Content    `json:"contents"`
	Template   Template     `json:"template"`
	TopArticle []TopArticle `json:"topArticle,optimize,omitempty"`
	Info       []Info       `json:"info,omitempty"`
}

type Info struct {
	Text string `json:"text"`
	Link string `json:"link"`
}

type Template struct {
	VideoTempName string        `json:"videoTempName"`
	AdType        string        `json:"adType"`
	Appearance    *Appearance   `json:"appearance"`
	MobileConfig  *MobileConfig `json:"mobileConfig"`
	Text          *Text         `json:"text"`
	Color         *Color        `json:"color"`
	Controls      *Controls     `json:"controls"`
	LogoBand      *LogoBand     `json:"logoBand"`
	AdConfig      *AdConfig     `json:"adConfig"`
	OrderMethod   string        `json:"orderMethod"`
	LogoLeft      string        `json:"logoLeft,omitempty"`
	LogoRight     string        `json:"logoRight,omitempty"`
	ReadMore      string        `json:"readMore,omitempty"`
	AutoStart     string        `json:"autoStart,omitempty"`
}

type MobileConfig struct {
	CloseFloatingBtn bool `json:"closeFloatingBtn"`
	FloatOnBottom    bool `json:"floatOnBottom"`
	FloatingOnView   bool `json:"floatingOnView"`
	Width            int  `json:"width"`
	Position         int  `json:"position"`
	MarginBot        int  `json:"margin-bot"`
	MarginLeftRight  int  `json:"margin-left-right"`
}

type Appearance struct {
	PlayerLayout    *PlayerLayout    `json:"playerLayout"`
	PlayerSize      string           `json:"playerSize"`
	MaxWidth        int              `json:"maxWidth"`
	FloatingSetting *FloatingSetting `json:"floatingSetting"`
	ColumnSetting   ColumnSetting    `json:"columnSetting"`
}

type ColumnSetting struct {
	ColumnPosition string `json:"columnPosition"`
	ColumnNumber   int    `json:"columnNumber"`
}

type PlayerLayout struct {
	Type      mysql.TypePlayerLayout `json:"type"`
	ListVideo []ListVideo            `json:"listVideo"`
}

type ListVideo struct {
	Title       string `json:"title"`
	ViewInfo    int    `json:"viewInfo"`
	Like        int    `json:"like"`
	Description string `json:"description"`
}

type FloatingSetting struct {
	CloseFloatingBtn bool `json:"closeFloatingBtn"`
	FloatOnBottom    bool `json:"floatOnBottom"`
	FloatingOnView   bool `json:"floatingOnView"`
	Width            int  `json:"width"`
	Position         int  `json:"position"`
	MarginTopBot     int  `json:"margin-top-bot"`
	MarginLeftRight  int  `json:"margin-left-right"`
}

type Text struct {
	MainTitle     string `json:"mainTitle"`
	TitleOn       bool   `json:"titleOn"`
	DescriptionOn bool   `json:"descriptionOn"`
	ActionButton  bool   `json:"actionButton"`
	ReadMore      string `json:"readMore"`
}

type Color struct {
	Controls                 string `json:"controls"`
	Background               string `json:"background"`
	Title                    string `json:"title"`
	Description              string `json:"description"`
	Theme                    string `json:"theme"`
	TitleBackground          string `json:"title-background"`
	MainTitle                string `json:"mainTitle"`
	MainTitleBackground      string `json:"mainTitleBackground"`
	TitleTextBackgroundHover string `json:"titleTextBackgroundHover"`
	ActionButtonColor        string `json:"actionButtonColor"`
}

type Controls struct {
	DefaultSoundMode bool `json:"defaultSoundMode"`
	Fullscreen       bool `json:"fullscreen"`
	NextPrevArrow    bool `json:"nextPrevArrow"`
	NextPrevSkip     bool `json:"nextPrevSkip"`
	VideoConfig      bool `json:"videoConfig"`
	ViewsLikes       bool `json:"viewsLikes"`
	Share            bool `json:"share"`
}

type LogoBand struct {
	EndableLogo      bool        `json:"endableLogo"`
	PoweredByApacdex bool        `json:"poweredByApacdex"`
	PoweredText      string      `json:"poweredText"`
	LogoByPubPower   bool        `json:"logoByPubPower"`
	CustomLogo       *CustomLogo `json:"customLogo"`
}

type CustomLogo struct {
	Link         string `json:"link"`
	ClickThrough string `json:"clickThrough"`
}

type AdConfig struct {
	VastRetry int       `json:"vastRetry"`
	AutoSkip  *AutoSkip `json:"autoSkip"`
	Delay     int       `json:"delay"`
}

type AutoSkip struct {
	TimeToSkip  int `json:"timeToSkip"`
	AutoSkipBtn int `json:"autoSkipBtn"`
	AdsNums     int `json:"adsNums"`
}

type Content struct {
	CreateTime string   `json:"create_time"`
	DeletedAt  string   `json:"deleted_at"`
	Des        string   `json:"des"`
	ID         string   `json:"id"`
	IsDefault  string   `json:"is_default"`
	Link       string   `json:"link"`
	Thumb      string   `json:"thumb"`
	Title      string   `json:"title"`
	UserID     string   `json:"user_id"`
	VideoURL   VideoURL `json:"video_url"`
}

type VideoURL struct {
	M3U8 string `json:"m3u8"`
	Mp4  string `json:"mp4"`
	Ogg  string `json:"ogg"`
}

type TopArticle struct {
	Title string `json:"title"`
	Image string `json:"image"`
	Link  string `json:"link"`
}

func (t PlayerConfig) String() string {
	jsonConfig, _ := json.Marshal(t)
	return html.UnescapeString(string(jsonConfig))
}
