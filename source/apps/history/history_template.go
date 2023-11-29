package history

import (
	"encoding/json"
	"source/core/technology/mysql"
)

type Template struct {
	Detail    DetailTemplate
	CreatorId int64
	RecordOld mysql.TableTemplate
	RecordNew mysql.TableTemplate
}

func (t *Template) Page() string {
	return "Template"
}

type DetailTemplate int

const (
	DetailTemplateFE DetailTemplate = iota + 1
	DetailTemplateBE
)

func (t DetailTemplate) String() string {
	switch t {
	case DetailTemplateFE:
		return "template_fe"
	case DetailTemplateBE:
		return "template_be"
	}
	return ""
}

func (t DetailTemplate) App() string {
	switch t {
	case DetailTemplateFE:
		return "FE"
	case DetailTemplateBE:
		return "BE"
	}
	return ""
}

func (t *Template) Type() TYPEHistory {
	return TYPEHistoryTemplate
}

func (t *Template) Action() mysql.TYPEObjectType {
	if t.RecordOld.Id == 0 && t.RecordNew.Id != 0 {
		return mysql.TYPEObjectTypeAdd
	} else if t.RecordOld.Id != 0 && t.RecordNew.Id == 0 {
		return mysql.TYPEObjectTypeDel
	}
	return mysql.TYPEObjectTypeUpdate
}

func (t *Template) Data() mysql.TableHistory {
	switch t.Detail {
	case DetailTemplateFE:
		return t.getHistoryTemplateFE()
	case DetailTemplateBE:
		return t.getHistoryTemplateBE()
	}
	return mysql.TableHistory{}
}

func (t *Template) CompareData(history mysql.TableHistory) (res []ResponseCompare) {
	switch history.DetailType {
	case DetailTemplateFE.String():
		return t.compareDataTemplateFE(history)
	case DetailTemplateBE.String():
		return t.compareDataTemplateBE(history)
	}
	return []ResponseCompare{}
}

func (t *Template) getRootRecord() (record mysql.TableTemplate) {
	switch t.Action() {
	case mysql.TYPEObjectTypeAdd:
		return t.RecordNew
	case mysql.TYPEObjectTypeUpdate:
		return t.RecordNew
	case mysql.TYPEObjectTypeDel:
		return t.RecordOld
	}
	return
}

type templateFE struct {
	Name                      *string      `json:"name,omitempty"`
	PlayerType                *string      `json:"player_type,omitempty"`
	Appearance                appearance   `json:"appearance,omitempty"`
	MobileConfig              mobileConfig `json:"mobile_config,omitempty"`
	MainTitle                 *string      `json:"main_title,omitempty"`
	MainTitleText             *string      `json:"main_title_text,omitempty"`
	ContentTitle              *string      `json:"content_title,omitempty"`
	Description               *string      `json:"description,omitempty"`
	ControlsColor             *string      `json:"controls_color,omitempty"`
	ThemeColor                *string      `json:"theme_color,omitempty"`
	BackgroundColor           *string      `json:"background_color,omitempty"`
	TitleColor                *string      `json:"title_color,omitempty"`
	TitleBackgroundColor      *string      `json:"title_background_color,omitempty"`
	DescriptionColor          *string      `json:"description_color,omitempty"`
	DefaultSoundMode          *string      `json:"default_sound_mode,omitempty"`
	FullscreenButton          *string      `json:"fullscreen_button,omitempty"`
	NextPrevArrows            *string      `json:"next_prev_arrows,omitempty"`
	NextPrev10sec             *string      `json:"next_prev_10_sec,omitempty"`
	VideoConfig               *string      `json:"video_config,omitempty"`
	ShowViewsLikes            *string      `json:"show_views_likes,omitempty"`
	ShareButton               *string      `json:"share_button,omitempty"`
	EnableLogo                *string      `json:"enable_logo,omitempty"`
	CustomLogo                *string      `json:"custom_logo,omitempty"`
	Link                      *string      `json:"link,omitempty"`
	ClickThrough              *string      `json:"click_through,omitempty"`
	PoweredByPubPower         *string      `json:"powered_by_pub_power,omitempty"`
	WaitForAdBeforeContent    *string      `json:"wait_for_ad_before_content,omitempty"`
	VastRetry                 *int         `json:"vast_retry,omitempty"`
	AutoSkip                  *string      `json:"auto_skip,omitempty"`
	TimeToSkip                *int         `json:"time_to_skip,omitempty"`
	ShowAutoSkipButtons       *int         `json:"show_auto_skip_buttons,omitempty"`
	NumberOfPreRollAds        *int         `json:"number_of_pre_roll_ads,omitempty"`
	Delay                     *int         `json:"delay,omitempty,omitempty"`
	TitleEnable               *string      `json:"title_enable,omitempty"`
	ShowControls              *string      `json:"show_controls,omitempty"`
	SubTitle                  *string      `json:"sub_title,omitempty"`
	ActionButton              *string      `json:"action_button,omitempty"`
	SubTitleText              *string      `json:"sub_title_text,omitempty"`
	ActionButtonText          *string      `json:"action_button_text,omitempty"`
	MainTitleBackgroundColor  *string      `json:"main_title_background_color,omitempty"`
	MainTitleColor            *string      `json:"main_title_color,omitempty"`
	TitleHoverBackgroundColor *string      `json:"title_hover_background_color,omitempty"`
	ActionButtonColor         *string      `json:"action_button_color,omitempty"`
}

type appearance struct {
	PlayerLayout        *string `json:"player_layout,omitempty"`
	PlayerSize          *string `json:"player_size,omitempty"`
	Ratio               *string `json:"ratio,omitempty"`
	Width               *int    `json:"width,omitempty"`
	AutoStart           *string `json:"auto_start,omitempty"`
	PlayerMode          *string `json:"player_mode,omitempty"`
	FloatingOnDesktop   *string `json:"floating_on_desktop,omitempty"`
	CloseFloatingButton *string `json:"close_floating_button,omitempty"`
	FloatOnBottom       *string `json:"float_on_bottom,omitempty"`
	FloatingOnView      *string `json:"floating_on_view,omitempty"`
	FloatingWidth       *int    `json:"floating_width,omitempty"`
	Position            *string `json:"position,omitempty"`
	MarginBottom        *int    `json:"margin_bottom,omitempty"`
	MarginTop           *int    `json:"margin_top,omitempty"`
	MarginRight         *int    `json:"margin_right,omitempty"`
	MarginLeft          *int    `json:"margin_left,omitempty"`
	ColumnsNumber       *int    `json:"columns_number,omitempty"`
	ColumnsPosition     *string `json:"columns_position,omitempty"`
}

type mobileConfig struct {
	FloatingOnMobile    *string `json:"floating_on_mobile,omitempty"`
	CloseFloatingButton *string `json:"close_floating_button,omitempty"`
	FloatOnBottom       *string `json:"float_on_bottom,omitempty"`
	FloatingOnView      *string `json:"floating_on_view,omitempty"`
	FloatingWidth       *int    `json:"floating_width,omitempty"`
	Position            *string `json:"position,omitempty"`
	MarginBottom        *int    `json:"margin_bottom,omitempty"`
	MarginRight         *int    `json:"margin_right,omitempty"`
	MarginLeft          *int    `json:"margin_left,omitempty"`
}

func (t *Template) getHistoryTemplateFE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := templateFE{}
	newData := templateFE{}
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.Template,
		ObjectId:   t.getRootRecord().Id,
		ObjectName: t.getRootRecord().Name,
		ObjectType: t.Action(),
		DetailType: t.Detail.String(),
		App:        t.Detail.App(),
		UserId:     t.getRootRecord().UserId,
	}
	var bNewData, bOldData []byte
	if t.RecordNew.Id != 0 {
		newData.MakeData(t.RecordNew)
		bNewData, _ = json.Marshal(newData)
	}
	if t.RecordOld.Id != 0 {
		oldData.MakeData(t.RecordOld)
		bOldData, _ = json.Marshal(oldData)
	}
	if t.Action() == mysql.TYPEObjectTypeAdd {
		history.Title = "Add Template"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update Template"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete Template"
		history.OldData = string(bOldData)
	}
	return
}

func (rec *templateFE) MakeData(record mysql.TableTemplate) {
	rec.Name = &record.Name
	playerType := record.Type.String()
	rec.PlayerType = &playerType
	switch record.Type {
	case mysql.TYPEInStream:
		rec.makeTemplateInStream(record)
		break
	case mysql.TYPEOutStream:
		rec.makeTemplateOutStream(record)
		break
	case mysql.TYPETopArticles:
		rec.makeTemplateTopArticles(record)
		break
	}
}

func (rec *templateFE) makeTemplateInStream(record mysql.TableTemplate) {
	playerLayout := record.PlayerLayout.String()
	rec.Appearance.PlayerLayout = &playerLayout
	// Tạo main appearance
	rec.makeAppearance(record)
	// Tạo mobile config
	rec.makeMobileConfig(record)

	columnsPosition := ""
	if record.ColumnsPosition == 1 {
		columnsPosition = "Right"
	} else {
		columnsPosition = "Left"
	}
	if record.PlayerLayout == 6 {
		rec.Appearance.ColumnsNumber = &record.ColumnsNumber
		rec.Appearance.ColumnsPosition = &columnsPosition
	} else if record.PlayerLayout == 7 {
		rec.Appearance.ColumnsPosition = &columnsPosition
	}

	mainTitle := record.MainTitle.String()
	rec.MainTitle = &mainTitle

	titleEnable := record.TitleEnable.String()
	rec.TitleEnable = &titleEnable

	if record.PlayerLayout == 3 || record.PlayerLayout == 4 || record.PlayerLayout == 5 || record.PlayerLayout == 7 {
		description := record.DescriptionEnable.String()
		rec.Description = &description
	}

	showControls := record.ShowControls.String()
	rec.ShowControls = &showControls

	if record.MainTitle == mysql.TypeOn {
		rec.MainTitleText = &record.MainTitleText
	}

	if record.SubTitle == mysql.TypeOn {
		rec.SubTitleText = &record.SubTitleText
	}

	contentTitle := record.TitleEnable.String()
	rec.ContentTitle = &contentTitle

	// Color
	rec.ControlsColor = &record.ControlColor
	rec.BackgroundColor = &record.BackgroundColor
	rec.TitleColor = &record.TitleColor

	if record.PlayerLayout == 3 || record.PlayerLayout == 4 || record.PlayerLayout == 5 || record.PlayerLayout == 7 {
		rec.DescriptionColor = &record.DescriptionColor
	}

	// Controls
	fullScreenButton := record.FullscreenButton.String()
	rec.FullscreenButton = &fullScreenButton

	nextPrevArrowsButton := record.NextPrevArrowsButton.String()
	rec.NextPrevArrows = &nextPrevArrowsButton

	nextPrevTime := record.NextPrevTime.String()
	rec.NextPrev10sec = &nextPrevTime

	videoConfig := record.VideoConfig.String()
	rec.VideoConfig = &videoConfig

	showViewLikes := record.ShowStats.String()
	rec.ShowViewsLikes = &showViewLikes

	shareButton := record.ShareButton.String()
	rec.ShareButton = &shareButton

	// Logo
	enableLog := record.EnableLogo.String()
	rec.EnableLogo = &enableLog

	rec.Link = &record.Link

	rec.ClickThrough = &record.ClickThrough

	poweredByPubPower := record.PoweredBy.String()
	rec.PoweredByPubPower = &poweredByPubPower

	// Advertising
	rec.VastRetry = &record.VastRetry

	autoSkip := record.AutoSkip.String()
	rec.AutoSkip = &autoSkip

	if record.AutoSkip == mysql.TypeOn {
		rec.TimeToSkip = &record.TimeToSkip
		rec.ShowAutoSkipButtons = &record.ShowAutoSkipButton
		rec.NumberOfPreRollAds = &record.NumberOfPreRollAds
	} else {
		rec.Delay = &record.Delay
	}
}

func (rec *templateFE) makeTemplateOutStream(record mysql.TableTemplate) {
	playerLayout := record.PlayerLayout.String()
	rec.Appearance.PlayerLayout = &playerLayout
	// Tạo main appearance
	rec.makeAppearance(record)
	// Tạo mobile config
	rec.makeMobileConfig(record)

	rec.VastRetry = &record.VastRetry
	rec.Delay = &record.Delay
}

func (rec *templateFE) makeTemplateTopArticles(record mysql.TableTemplate) {
	playerLayout := record.PlayerLayout.String()
	rec.Appearance.PlayerLayout = &playerLayout
	// Tạo main appearance
	rec.makeAppearance(record)
	// Tạo mobile config
	rec.makeMobileConfig(record)

	// Text / Logo
	mainTitle := record.MainTitle.String()
	rec.MainTitle = &mainTitle

	subTitle := record.SubTitle.String()
	rec.SubTitle = &subTitle

	actionButton := record.ActionButton.String()
	rec.ActionButton = &actionButton

	if record.MainTitle == mysql.TypeOn {
		rec.MainTitleText = &record.MainTitleText
	}

	if record.SubTitle == mysql.TypeOn {
		rec.SubTitleText = &record.SubTitleText
	}

	if record.ActionButton == mysql.TypeOn {
		rec.ActionButtonText = &record.ActionButtonText
	}

	// Color
	rec.ThemeColor = &record.ThemeColor
	rec.TitleBackgroundColor = &record.TitleBackgroundColor
	rec.MainTitleBackgroundColor = &record.MainTitleBackgroundColor
	rec.MainTitleColor = &record.MainTitleColor
	rec.TitleHoverBackgroundColor = &record.TitleHoverBackgroundColor
	rec.ActionButtonColor = &record.ActionButtonColor

	enableLog := record.EnableLogo.String()
	rec.EnableLogo = &enableLog

	rec.Link = &record.Link
	rec.ClickThrough = &record.ClickThrough

	// Ad config
	rec.VastRetry = &record.VastRetry
	rec.Delay = &record.Delay
}

func (rec *templateFE) makeAppearance(record mysql.TableTemplate) {
	playerSize := record.Size.String()
	rec.Appearance.PlayerSize = &playerSize
	if record.Size == mysql.TypeSizeFixed {
		rec.Appearance.Width = &record.Width
	}

	autoStart := record.AutoStart.String()
	rec.Appearance.AutoStart = &autoStart

	playMode := record.PlayMode.String()
	rec.Appearance.PlayerMode = &playMode

	floatingOnDesktop := record.FloatingOnDesktop.String()
	rec.Appearance.FloatingOnDesktop = &floatingOnDesktop

	if record.FloatingOnDesktop == mysql.TypeOn {
		closeFloatingButton := record.CloseFloatingButtonDesktop.String()
		rec.Appearance.CloseFloatingButton = &closeFloatingButton

		floatOnBottom := record.FloatOnBottom.String()
		rec.Appearance.FloatOnBottom = &floatOnBottom

		if record.Type == mysql.TYPEInStream || record.Type == mysql.TYPETopArticles {
			floatingOnView := record.FloatingOnView.String()
			rec.Appearance.FloatingOnView = &floatingOnView
		}
		rec.Appearance.FloatingWidth = &record.FloatingWidth
		position := record.FloatingPositionDesktop.String()
		rec.Appearance.Position = &position
		switch record.FloatingPositionDesktop {
		case mysql.TypePositionDesktopBottomRight:
			rec.Appearance.MarginBottom = &record.MarginBottomDesktop
			rec.Appearance.MarginRight = &record.MarginRightDesktop
			break
		case mysql.TypePositionDesktopBottomLeft:
			rec.Appearance.MarginBottom = &record.MarginBottomDesktop
			rec.Appearance.MarginLeft = &record.MarginLeftDesktop
			break
		case mysql.TypePositionDesktopTopRight:
			rec.Appearance.MarginTop = &record.MarginTopDesktop
			rec.Appearance.MarginRight = &record.MarginRightDesktop
			break
		case mysql.TypePositionDesktopTopLeft:
			rec.Appearance.MarginTop = &record.MarginTopDesktop
			rec.Appearance.MarginLeft = &record.MarginLeftDesktop
			break

		}
	}
}

func (rec *templateFE) makeMobileConfig(record mysql.TableTemplate) {
	mobileConfig := mobileConfig{}
	floatingOnMobile := record.FloatingOnMobile.String()
	mobileConfig.FloatingOnMobile = &floatingOnMobile
	if record.FloatingOnMobile == mysql.TypeOn {
		closeFloatingButton := record.CloseFloatingButtonMobile.String()
		mobileConfig.CloseFloatingButton = &closeFloatingButton

		floatOnBottom := record.FloatOnBottomMobile.String()
		mobileConfig.FloatOnBottom = &floatOnBottom

		if record.Type == mysql.TYPEInStream || record.Type == mysql.TYPETopArticles {
			floatingOnView := record.FloatingOnViewMobile.String()
			mobileConfig.FloatingOnView = &floatingOnView
		}

		mobileConfig.FloatingWidth = &record.FloatingWidthMobile
		position := record.FloatingPositionMobile.String()
		mobileConfig.Position = &position
		switch record.FloatingPositionMobile {
		case mysql.TypePositionMobileBottomRight:
			mobileConfig.MarginBottom = &record.MarginBottomMobile
			mobileConfig.MarginRight = &record.MarginRightMobile
			break
		case mysql.TypePositionMobileBottomLeft:
			mobileConfig.MarginBottom = &record.MarginBottomMobile
			mobileConfig.MarginLeft = &record.MarginLeftMobile
			break
		}
	}
	rec.MobileConfig = mobileConfig
}

func (t *Template) compareDataTemplateFE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew templateFE
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Name
	res, err := makeResponseCompare("Name", recordOld.Name, recordNew.Name, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// PlayerType
	res, err = makeResponseCompare("Player Type", recordOld.PlayerType, recordNew.PlayerType, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Player Layout
	res, err = makeResponseCompare("Player Layout", recordOld.Appearance.PlayerLayout, recordNew.Appearance.PlayerLayout, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Columns Number
	res, err = makeResponseCompare("Columns Number", pointerIntToString(recordOld.Appearance.ColumnsNumber), pointerIntToString(recordNew.Appearance.ColumnsNumber), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Columns Position
	res, err = makeResponseCompare("Columns Position", recordOld.Appearance.ColumnsPosition, recordNew.Appearance.ColumnsPosition, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Player Size
	res, err = makeResponseCompare("Player Size", recordOld.Appearance.PlayerSize, recordNew.Appearance.PlayerSize, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Width
	res, err = makeResponseCompare("Width", pointerIntToString(recordOld.Appearance.Width), pointerIntToString(recordNew.Appearance.Width), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Auto Start
	res, err = makeResponseCompare("Auto Start", recordOld.Appearance.AutoStart, recordNew.Appearance.AutoStart, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Desktop
	// Floating On Desktop
	res, err = makeResponseCompare("Player Mode > Desktop > Floating On Desktop", recordOld.Appearance.FloatingOnDesktop, recordNew.Appearance.FloatingOnDesktop, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Close Floating Button
	res, err = makeResponseCompare("Player Mode > Desktop > Close Floating Button", recordOld.Appearance.CloseFloatingButton, recordNew.Appearance.CloseFloatingButton, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Float On Bottom
	res, err = makeResponseCompare("Player Mode > Desktop > Float On Bottom", recordOld.Appearance.FloatOnBottom, recordNew.Appearance.FloatOnBottom, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Floating On View
	res, err = makeResponseCompare("Player Mode > Desktop > Floating On View", recordOld.Appearance.FloatingOnView, recordNew.Appearance.FloatingOnView, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Floating Width
	res, err = makeResponseCompare("Player Mode > Desktop > Floating Width", pointerIntToString(recordOld.Appearance.FloatingWidth), pointerIntToString(recordNew.Appearance.FloatingWidth), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Position
	res, err = makeResponseCompare("Player Mode > Desktop > Position", recordOld.Appearance.Position, recordNew.Appearance.Position, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Margin Bottom
	res, err = makeResponseCompare("Player Mode > Desktop > Margin Bottom", pointerIntToString(recordOld.Appearance.MarginBottom), pointerIntToString(recordNew.Appearance.MarginBottom), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Margin Top
	res, err = makeResponseCompare("Player Mode > Desktop > Margin Top", pointerIntToString(recordOld.Appearance.MarginTop), pointerIntToString(recordNew.Appearance.MarginTop), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Margin Right
	res, err = makeResponseCompare("Player Mode > Desktop > Margin Right", pointerIntToString(recordOld.Appearance.MarginRight), pointerIntToString(recordNew.Appearance.MarginRight), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Margin Left
	res, err = makeResponseCompare("Player Mode > Desktop > Margin Left", pointerIntToString(recordOld.Appearance.MarginLeft), pointerIntToString(recordNew.Appearance.MarginLeft), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Mobile
	// Floating On Mobile
	res, err = makeResponseCompare("Player Mode > Mobile > Floating On Mobile", recordOld.MobileConfig.FloatingOnMobile, recordNew.MobileConfig.FloatingOnMobile, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Mobile Config Close Floating Button
	res, err = makeResponseCompare("Player Mode > Mobile > Close Floating Button", recordOld.MobileConfig.CloseFloatingButton, recordNew.MobileConfig.CloseFloatingButton, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Mobile Config Float On Bottom
	res, err = makeResponseCompare("Player Mode > Mobile > Float On Bottom", recordOld.MobileConfig.FloatOnBottom, recordNew.MobileConfig.FloatOnBottom, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Mobile Config Floating On View
	res, err = makeResponseCompare("Player Mode > Mobile > Floating On View", recordOld.MobileConfig.FloatingOnView, recordNew.MobileConfig.FloatingOnView, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Mobile Config Floating Width
	res, err = makeResponseCompare("Player Mode > Mobile > Floating Width", pointerIntToString(recordOld.MobileConfig.FloatingWidth), pointerIntToString(recordNew.MobileConfig.FloatingWidth), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Mobile Config Position
	res, err = makeResponseCompare("Player Mode > Mobile > Position", recordOld.MobileConfig.Position, recordNew.MobileConfig.Position, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Mobile Config Margin Bottom
	res, err = makeResponseCompare("Player Mode > Mobile > Margin Bottom", pointerIntToString(recordOld.MobileConfig.MarginBottom), pointerIntToString(recordNew.MobileConfig.MarginBottom), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Mobile Config Margin Right
	res, err = makeResponseCompare("Player Mode > Mobile > Margin Right", pointerIntToString(recordOld.MobileConfig.MarginRight), pointerIntToString(recordNew.MobileConfig.MarginRight), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Mobile Config Margin Left
	res, err = makeResponseCompare("Player Mode > Mobile > Margin Left", pointerIntToString(recordOld.MobileConfig.MarginLeft), pointerIntToString(recordNew.MobileConfig.MarginLeft), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}

	// ***Display Options***
	// Main Title
	res, err = makeResponseCompare("Show Main Title", recordOld.MainTitle, recordNew.MainTitle, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Sub Title
	res, err = makeResponseCompare("Show Sub Title", recordOld.SubTitle, recordNew.SubTitle, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Action Button
	res, err = makeResponseCompare("Show Action Button", recordOld.ActionButton, recordNew.ActionButton, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Content Title
	res, err = makeResponseCompare("Show Content Title", recordOld.TitleEnable, recordNew.TitleEnable, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Content Description
	res, err = makeResponseCompare("Show Content Description", recordOld.Description, recordNew.Description, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Controls
	res, err = makeResponseCompare("Show Controls", recordOld.ShowControls, recordNew.ShowControls, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Main Title Text
	res, err = makeResponseCompare("Main Title Text", recordOld.MainTitleText, recordNew.MainTitleText, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Sub Title Text
	res, err = makeResponseCompare("Sub Title Text", recordOld.SubTitleText, recordNew.SubTitleText, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Action Button Text
	res, err = makeResponseCompare("Action Button Text", recordOld.ActionButtonText, recordNew.ActionButtonText, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}

	// ***End Display Option***

	// *** Color ***
	// Theme Color
	res, err = makeResponseCompare("Theme Color", recordOld.ThemeColor, recordNew.ThemeColor, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Controls Color
	res, err = makeResponseCompare("Controls Color", recordOld.ControlsColor, recordNew.ControlsColor, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Background Color
	res, err = makeResponseCompare("Background Color", recordOld.BackgroundColor, recordNew.BackgroundColor, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Main Title Background Color
	res, err = makeResponseCompare("Main Title Background Color", recordOld.MainTitleBackgroundColor, recordNew.MainTitleBackgroundColor, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Main Title Color
	res, err = makeResponseCompare("Main Title Color", recordOld.MainTitleColor, recordNew.MainTitleColor, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Title Color
	res, err = makeResponseCompare("Title Color", recordOld.TitleColor, recordNew.TitleColor, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Title Background Color
	res, err = makeResponseCompare("Title Background Color", recordOld.TitleBackgroundColor, recordNew.TitleBackgroundColor, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Title Hover Background Color
	res, err = makeResponseCompare("Title Hover Background Color", recordOld.TitleHoverBackgroundColor, recordNew.TitleHoverBackgroundColor, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Action Button Color
	res, err = makeResponseCompare("Action Button Color", recordOld.ActionButtonColor, recordNew.ActionButtonColor, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Description Color
	res, err = makeResponseCompare("Description Color", recordOld.DescriptionColor, recordNew.DescriptionColor, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// *** End Color ***

	// *** Controls ***
	// Fullscreen Button
	res, err = makeResponseCompare("Fullscreen Button", recordOld.FullscreenButton, recordNew.FullscreenButton, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Next Prev Arrows
	res, err = makeResponseCompare("Next Prev Arrows", recordOld.NextPrevArrows, recordNew.NextPrevArrows, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Next Prev 10sec
	res, err = makeResponseCompare("Next Prev 10sec", recordOld.NextPrev10sec, recordNew.NextPrev10sec, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Video Config
	res, err = makeResponseCompare("Video Config", recordOld.VideoConfig, recordNew.VideoConfig, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Show Views Likes
	res, err = makeResponseCompare("Show Views Likes", recordOld.ShowViewsLikes, recordNew.ShowViewsLikes, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Share Button
	res, err = makeResponseCompare("Share Button", recordOld.ShareButton, recordNew.ShareButton, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// *** End Controls ***

	// *** Logo ***
	// Enable Logo
	res, err = makeResponseCompare("Show logo", recordOld.EnableLogo, recordNew.EnableLogo, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Powered By PubPower
	res, err = makeResponseCompare("Powered By", recordOld.PoweredByPubPower, recordNew.PoweredByPubPower, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Link
	res, err = makeResponseCompare("Link Logo", recordOld.Link, recordNew.Link, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Click Through
	res, err = makeResponseCompare("URL", recordOld.ClickThrough, recordNew.ClickThrough, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// *** End Logo ***

	// *** Advertising ***
	// Vast Retry
	res, err = makeResponseCompare("Vast Retry", pointerIntToString(recordOld.VastRetry), pointerIntToString(recordNew.VastRetry), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Auto Skip
	res, err = makeResponseCompare("Auto Skip", recordOld.AutoSkip, recordNew.AutoSkip, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Time To Skip
	res, err = makeResponseCompare("Time To Skip", pointerIntToString(recordOld.TimeToSkip), pointerIntToString(recordNew.TimeToSkip), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Show Auto Skip Buttons
	res, err = makeResponseCompare("Show Auto Skip Buttons", pointerIntToString(recordOld.ShowAutoSkipButtons), pointerIntToString(recordNew.ShowAutoSkipButtons), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Number Of PreRoll Ads
	res, err = makeResponseCompare("Number Of PreRoll Ads", pointerIntToString(recordOld.NumberOfPreRollAds), pointerIntToString(recordNew.NumberOfPreRollAds), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Delay
	res, err = makeResponseCompare("Delay", pointerIntToString(recordOld.Delay), pointerIntToString(recordNew.Delay), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// *** End Advertising ***
	return
}

func (t *Template) getHistoryTemplateBE() (history mysql.TableHistory) {
	history = t.getHistoryTemplateFE()
	return
}

func (t *Template) compareDataTemplateBE(history mysql.TableHistory) (responses []ResponseCompare) {
	responses = t.compareDataTemplateFE(history)
	return
}
