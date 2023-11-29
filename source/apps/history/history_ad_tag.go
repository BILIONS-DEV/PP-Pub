package history

import (
	"encoding/json"
	"source/core/technology/mysql"
	"strconv"
)

type AdTag struct {
	Detail    DetailAdTag
	CreatorId int64
	RecordOld mysql.TableInventoryAdTag
	RecordNew mysql.TableInventoryAdTag
}

func (t *AdTag) Page() string {
	return "Supply"
}

type DetailAdTag int

const (
	DetailAdTagFE DetailAdTag = iota + 1
	DetailAdTagBE
)

func (t DetailAdTag) String() string {
	switch t {
	case DetailAdTagFE:
		return "ad_tag_fe"
	case DetailAdTagBE:
		return "ad_tag_be"
	}
	return ""
}

func (t DetailAdTag) App() string {
	switch t {
	case DetailAdTagFE:
		return "FE"
	case DetailAdTagBE:
		return "BE"
	}
	return ""
}

func (t *AdTag) Type() TYPEHistory {
	return TYPEHistoryAdTag
}

func (t *AdTag) Action() mysql.TYPEObjectType {
	if t.RecordOld.Id == 0 && t.RecordNew.Id != 0 {
		return mysql.TYPEObjectTypeAdd
	} else if t.RecordOld.Id != 0 && t.RecordNew.Id == 0 {
		return mysql.TYPEObjectTypeDel
	}
	return mysql.TYPEObjectTypeUpdate
}

func (t *AdTag) Data() mysql.TableHistory {
	switch t.Detail {
	case DetailAdTagFE:
		return t.getHistoryAdTagFE()
	case DetailAdTagBE:
		return t.getHistoryAdTagBE()
	}
	return mysql.TableHistory{}
}

func (t *AdTag) CompareData(history mysql.TableHistory) (res []ResponseCompare) {
	switch history.DetailType {
	case DetailAdTagFE.String():
		return t.compareDataAdTagFE(history)
	case DetailAdTagBE.String():
		return t.compareDataAdTagBE(history)
	}
	return []ResponseCompare{}
}

func (t *AdTag) getRootRecord() (record mysql.TableInventoryAdTag) {
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

type adTagFE struct {
	Name                    *string   `json:"name,omitempty"`
	AdType                  *string   `json:"ad_type,omitempty"`
	Status                  *string   `json:"status,omitempty"`
	Renderer                *string   `json:"renderer,omitempty"`
	GAMAdUnit               *string   `json:"gam_ad_unit,omitempty"`
	PrimaryAdSize           *string   `json:"primary_ad_size,omitempty"`
	AdditionalAdSizes       *[]string `json:"additional_ad_sizes,omitempty"`
	BidOutstream            *string   `json:"bid_outstream,omitempty"`
	Passback                *string   `json:"passback,omitempty"`
	Template                *string   `json:"template,omitempty"`
	ContentSource           *string   `json:"content_source,omitempty"`
	Playlist                *string   `json:"playlist,omitempty"`
	PassbackType            *string   `json:"passback_type,omitempty"`
	InlineTag               *string   `json:"inline_tag,omitempty"`
	Position                *string   `json:"position,omitempty"`
	CloseButton             *string   `json:"close_button,omitempty"`
	FrequencyCaps           *int      `json:"frequency_caps,omitempty"`
	ContentType             *string   `json:"content_type,omitempty"`
	Feed                    *string   `json:"feed,omitempty"`
	MainTitle               *string   `json:"main_title,omitempty"`
	BackgroundColor         *string   `json:"background_color,omitempty"`
	TitleColor              *string   `json:"title_color,omitempty"`
	AdSize                  *string   `json:"ad_size,omitempty"`
	ResponsiveType          *string   `json:"responsive_type,omitempty"`
	ShiftContent            *string   `json:"shift_content,omitempty"`
	EnableDesktop           *string   `json:"enable_desktop"`
	EnableMobile            *string   `json:"enable_mobile"`
	PrimaryAdSizeMobile     *string   `json:"primary_ad_size_mobile,omitempty"`
	AdditionalAdSizesMobile *[]string `json:"additional_ad_sizes_mobile,omitempty"`
	PositionMobile          *string   `json:"position_mobile,omitempty"`
	CloseButtonMobile       *string   `json:"close_button_mobile,omitempty"`
}

func (t *AdTag) getHistoryAdTagBE() (history mysql.TableHistory) {
	return t.getHistoryAdTagFE()
}

func (t *AdTag) getHistoryAdTagFE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := adTagFE{}
	newData := adTagFE{}
	history = mysql.TableHistory{
		CreatorId:   t.CreatorId,
		Object:      mysql.Tables.InventoryAdTag,
		ObjectId:    t.getRootRecord().Id,
		ObjectName:  t.getRootRecord().Name,
		ObjectType:  t.Action(),
		DetailType:  t.Detail.String(),
		App:         t.Detail.App(),
		UserId:      t.getRootRecord().UserId,
		InventoryId: t.RecordNew.InventoryId,
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
		history.Title = "Add Ad Tag"
		history.NewData = string(bNewData)
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update Ad Tag"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete Ad Tag"
		history.OldData = string(bOldData)
	}
	return
}

func (rec *adTagFE) MakeData(record mysql.TableInventoryAdTag) {
	rec.Name = &record.Name
	adType := record.Type.String()
	rec.AdType = &adType
	status := record.Status.String()
	rec.Status = &status
	switch record.Type {
	case mysql.TYPEDisplay:
		rec.makeTagDisplay(record)
		break
	case mysql.TYPEInStream:
		rec.makeTagInStream(record)
		break
	case mysql.TYPEOutStream:
		rec.makeTagOutStream(record)
		break
	case mysql.TYPETopArticles:
		rec.makeTagPinZone(record)
		break
	case mysql.TYPEStickyBanner:
		rec.makeTagStickyBanner(record)
		break
	case mysql.TYPEInterstitial:
		rec.makeTagInterstitial(record)
	case mysql.TYPEPlayZone:
		rec.makeTagRelatedZone(record)
		break
	}
}

func (rec *adTagFE) makeTagDisplay(record mysql.TableInventoryAdTag) {
	var inventoryConfig mysql.TableInventoryConfig
	mysql.Client.Where("inventory_id = ?", record.InventoryId).Find(&inventoryConfig)
	if inventoryConfig.GamAutoCreate == mysql.TypeOff {
		rec.GAMAdUnit = &record.Gam
	}
	adSize := record.AdSize.String()
	rec.AdSize = &adSize
	if record.AdSize == mysql.TYPEAdSizeFixed {
		var primaryAdsize mysql.TableAdSize
		mysql.Client.Find(&primaryAdsize, record.PrimaryAdSize)
		rec.PrimaryAdSize = &primaryAdsize.Name
		var additionalSize []string
		for _, size := range record.AdditionalAdSize {
			additionalSize = append(additionalSize, size.Name)
		}
		if len(additionalSize) > 0 {
			rec.AdditionalAdSizes = &additionalSize
		}
	} else if record.AdSize == mysql.TYPEAdSizeResponsive {
		responsiveType := record.ResponsiveType.String()
		rec.ResponsiveType = &responsiveType
	}
	rec.Passback = &record.PassBack
}

func (rec *adTagFE) makeTagInStream(record mysql.TableInventoryAdTag) {
	var inventoryConfig mysql.TableInventoryConfig
	mysql.Client.Where("inventory_id = ?", record.InventoryId).Find(&inventoryConfig)
	if inventoryConfig.GamAutoCreate == mysql.TypeOff {
		rec.GAMAdUnit = &record.Gam
	}

	var renderer string
	if record.Renderer == 1 {
		renderer = "Valueimpression Player"
		rec.Template = &record.Template.Name
		contentSource := record.ContentSource.String()
		rec.ContentSource = &contentSource
		rec.Playlist = &record.Playlist.Name
	} else {
		renderer = record.Renderer.String()
	}
	rec.Renderer = &renderer
}

func (rec *adTagFE) makeTagOutStream(record mysql.TableInventoryAdTag) {
	var inventoryConfig mysql.TableInventoryConfig
	mysql.Client.Where("inventory_id = ?", record.InventoryId).Find(&inventoryConfig)
	if inventoryConfig.GamAutoCreate == mysql.TypeOff {
		rec.GAMAdUnit = &record.Gam
	}

	var renderer string
	if record.Renderer == 1 {
		renderer = "Valueimpression Player"
		rec.Template = &record.Template.Name
		passbackType := record.PassBackType.String()
		rec.PassbackType = &passbackType
		if record.PassBackType == mysql.TYPEPassBackTypeInline {
			var inlineTag mysql.TableInventoryAdTag
			mysql.Client.Find(&inlineTag, record.InlineTag)
			rec.InlineTag = &inlineTag.Name
		} else if record.PassBackType == mysql.TYPEPassBackTypeCustom {
			rec.Passback = &record.PassBack
		}
	} else {
		renderer = record.Renderer.String()
	}
	rec.Renderer = &renderer
}

func (rec *adTagFE) makeTagPinZone(record mysql.TableInventoryAdTag) {
	var inventoryConfig mysql.TableInventoryConfig
	mysql.Client.Where("inventory_id = ?", record.InventoryId).Find(&inventoryConfig)
	if inventoryConfig.GamAutoCreate == mysql.TypeOff {
		rec.GAMAdUnit = &record.Gam
	}

	rec.Template = &record.Template.Name
	contentSource := record.ContentSource.String()
	rec.ContentSource = &contentSource
	if record.ContentSource == mysql.TypeContentSourceAuto {
		contentType := record.RelatedContent.String()
		rec.ContentType = &contentType
	} else if record.ContentSource == mysql.TypeContentSourceFeed {
		rec.Feed = &record.FeedUrl
	}

}

func (rec *adTagFE) makeTagStickyBanner(record mysql.TableInventoryAdTag) {
	var inventoryConfig mysql.TableInventoryConfig
	mysql.Client.Where("inventory_id = ?", record.InventoryId).Find(&inventoryConfig)
	if inventoryConfig.GamAutoCreate == mysql.TypeOff {
		rec.GAMAdUnit = &record.Gam
	}
	shiftContent := record.ShiftContent.String()
	rec.ShiftContent = &shiftContent

	// Desktop
	enableDesktop := record.EnableStickyDesktop.String()
	rec.EnableDesktop = &enableDesktop
	if record.EnableStickyDesktop == mysql.On {
		var primaryAdsize mysql.TableAdSize
		mysql.Client.Find(&primaryAdsize, record.PrimaryAdSize)
		rec.PrimaryAdSize = &primaryAdsize.Name
		var additionalSize []string
		for _, size := range record.AdditionalAdSize {
			additionalSize = append(additionalSize, size.Name)
		}
		if len(additionalSize) > 0 {
			rec.AdditionalAdSizes = &additionalSize
		}
		position := record.PositionSticky.String()
		rec.Position = &position
		var closeButton string
		if record.CloseButtonSticky == 1 {
			closeButton = "On"
		} else {
			closeButton = "Off"
		}
		rec.CloseButton = &closeButton
	}

	// Mobile
	enableMobile := record.EnableStickyMobile.String()
	rec.EnableMobile = &enableMobile
	if record.EnableStickyMobile == mysql.On {
		var primaryAdsizeMobile mysql.TableAdSize
		mysql.Client.Find(&primaryAdsizeMobile, record.PrimaryAdSizeMobile)
		rec.PrimaryAdSizeMobile = &primaryAdsizeMobile.Name
		var additionalSizeMobile []string
		for _, size := range record.AdditionalAdSize {
			additionalSizeMobile = append(additionalSizeMobile, size.Name)
		}
		if len(additionalSizeMobile) > 0 {
			rec.AdditionalAdSizesMobile = &additionalSizeMobile
		}
		positionMobile := record.PositionStickyMobile.String()
		rec.PositionMobile = &positionMobile
		var closeButtonMobile string
		if record.CloseButtonStickyMobile == 1 {
			closeButtonMobile = "On"
		} else {
			closeButtonMobile = "Off"
		}
		rec.CloseButtonMobile = &closeButtonMobile
	}
}

func (rec *adTagFE) makeTagInterstitial(record mysql.TableInventoryAdTag) {
	var inventoryConfig mysql.TableInventoryConfig
	mysql.Client.Where("inventory_id = ?", record.InventoryId).Find(&inventoryConfig)
	if inventoryConfig.GamAutoCreate == mysql.TypeOff {
		rec.GAMAdUnit = &record.Gam
	}
	rec.FrequencyCaps = &record.FrequencyCaps
}

func (rec *adTagFE) makeTagRelatedZone(record mysql.TableInventoryAdTag) {
	var inventoryConfig mysql.TableInventoryConfig
	mysql.Client.Where("inventory_id = ?", record.InventoryId).Find(&inventoryConfig)
	if inventoryConfig.GamAutoCreate == mysql.TypeOff {
		rec.GAMAdUnit = &record.Gam
	}

	template := "Template " + strconv.FormatInt(record.TemplateId, 10)
	rec.Template = &template

	contentType := record.RelatedContent.String()
	rec.ContentType = &contentType

	contentSource := record.ContentSource.String()
	rec.ContentSource = &contentSource

	rec.MainTitle = &record.MainTitle
	rec.Passback = &record.PassBack
	rec.BackgroundColor = &record.BackgroundColor
	rec.TitleColor = &record.TitleColor
}

func (t *AdTag) compareDataAdTagFE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew adTagFE
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	var adType string
	if recordOld.AdType != nil {
		adType = *recordOld.AdType
	}
	if recordNew.AdType != nil {
		adType = *recordNew.AdType
	}
	// Xử lý compare từng row

	// Name
	res, err := makeResponseCompare("Name", recordOld.Name, recordNew.Name, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Ad Type
	res, err = makeResponseCompare("Ad Type", recordOld.AdType, recordNew.AdType, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Status
	res, err = makeResponseCompare("Status", recordOld.Status, recordNew.Status, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Renderer
	res, err = makeResponseCompare("Renderer", recordOld.Renderer, recordNew.Renderer, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// GAM AdUnit
	res, err = makeResponseCompare("GAM AdUnit", recordOld.GAMAdUnit, recordNew.GAMAdUnit, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Ad Size
	res, err = makeResponseCompare("Ad Size", recordOld.AdSize, recordNew.AdSize, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// ShiftContent
	res, err = makeResponseCompare("Shift Content", recordOld.ShiftContent, recordNew.ShiftContent, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Enable Desktop
	res, err = makeResponseCompare("Enable Desktop", recordOld.EnableDesktop, recordNew.EnableDesktop, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Primary Ad Size
	if adType == mysql.TYPEStickyBanner.String() {
		res, err = makeResponseCompare("Desktop > Primary Ad Size", recordOld.PrimaryAdSize, recordNew.PrimaryAdSize, history.ObjectType)
		if err == nil {
			responses = append(responses, res)
		}
	} else {
		res, err = makeResponseCompare("Primary Ad Size", recordOld.PrimaryAdSize, recordNew.PrimaryAdSize, history.ObjectType)
		if err == nil {
			responses = append(responses, res)
		}
	}
	// Additional Ad Size
	if adType == mysql.TYPEStickyBanner.String() {
		res, err = makeResponseCompare("Desktop > Additional Ad Size", pointerArrayStringToString(recordOld.AdditionalAdSizes), pointerArrayStringToString(recordNew.AdditionalAdSizes), history.ObjectType)
		if err == nil {
			responses = append(responses, res)
		}
	} else {
		res, err = makeResponseCompare("Additional Ad Size", pointerArrayStringToString(recordOld.AdditionalAdSizes), pointerArrayStringToString(recordNew.AdditionalAdSizes), history.ObjectType)
		if err == nil {
			responses = append(responses, res)
		}
	}
	// Responsive Type
	res, err = makeResponseCompare("Responsive Type", recordOld.ResponsiveType, recordNew.ResponsiveType, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Bid Outstream
	//res, err = makeResponseCompare("Bid Outstream", recordOld.BidOutstream, recordNew.BidOutstream)
	//if err == nil {
	//	responses = append(responses, res)
	//}
	// Passback
	res, err = makeResponseCompare("Passback", recordOld.Passback, recordNew.Passback, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Template
	res, err = makeResponseCompare("Template", recordOld.Template, recordNew.Template, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Content Source
	res, err = makeResponseCompare("Content Source", recordOld.ContentSource, recordNew.ContentSource, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Playlist
	res, err = makeResponseCompare("Playlist", recordOld.Playlist, recordNew.Playlist, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Passback Type
	res, err = makeResponseCompare("Passback Type", recordOld.PassbackType, recordNew.PassbackType, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Inline Tag
	res, err = makeResponseCompare("InlineTag", recordOld.InlineTag, recordNew.InlineTag, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Position
	if adType == mysql.TYPEStickyBanner.String() {
		res, err = makeResponseCompare("Desktop > Position", recordOld.Position, recordNew.Position, history.ObjectType)
		if err == nil {
			responses = append(responses, res)
		}
	} else {
		res, err = makeResponseCompare("Position", recordOld.Position, recordNew.Position, history.ObjectType)
		if err == nil {
			responses = append(responses, res)
		}
	}
	// Close Button
	if adType == mysql.TYPEStickyBanner.String() {
		res, err = makeResponseCompare("Desktop > Close Button", recordOld.CloseButton, recordNew.CloseButton, history.ObjectType)
		if err == nil {
			responses = append(responses, res)
		}
	} else {
		res, err = makeResponseCompare("Close Button", recordOld.CloseButton, recordNew.CloseButton, history.ObjectType)
		if err == nil {
			responses = append(responses, res)
		}
	}
	// Frequency Caps
	res, err = makeResponseCompare("Frequency Caps", pointerIntToString(recordOld.FrequencyCaps), pointerIntToString(recordNew.FrequencyCaps), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Content Type
	res, err = makeResponseCompare("Content Type", recordOld.ContentType, recordNew.ContentType, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Feed
	res, err = makeResponseCompare("Feed", recordOld.Feed, recordNew.Feed, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Main Title
	res, err = makeResponseCompare("Main Title", recordOld.MainTitle, recordNew.MainTitle, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Background Color
	res, err = makeResponseCompare("Background Color", recordOld.BackgroundColor, recordNew.BackgroundColor, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Title Color
	res, err = makeResponseCompare("Title Color", recordOld.TitleColor, recordNew.TitleColor, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Enable Mobile
	res, err = makeResponseCompare("Enable Mobile", recordOld.EnableMobile, recordNew.EnableMobile, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Primary Ad Size Mobile
	if adType == mysql.TYPEStickyBanner.String() {
		res, err = makeResponseCompare("Mobile > Primary Ad Size", recordOld.PrimaryAdSizeMobile, recordNew.PrimaryAdSizeMobile, history.ObjectType)
		if err == nil {
			responses = append(responses, res)
		}
	}
	// Additional Ad Size Mobile
	if adType == mysql.TYPEStickyBanner.String() {
		res, err = makeResponseCompare("Mobile > Additional Ad Size", pointerArrayStringToString(recordOld.AdditionalAdSizesMobile), pointerArrayStringToString(recordNew.AdditionalAdSizesMobile), history.ObjectType)
		if err == nil {
			responses = append(responses, res)
		}
	}
	// Position Mobile
	if adType == mysql.TYPEStickyBanner.String() {
		res, err = makeResponseCompare("Mobile > Position", recordOld.PositionMobile, recordNew.PositionMobile, history.ObjectType)
		if err == nil {
			responses = append(responses, res)
		}
	}
	// Close Button Mobile
	if adType == mysql.TYPEStickyBanner.String() {
		res, err = makeResponseCompare("Mobile > Close Button", recordOld.CloseButtonMobile, recordNew.CloseButtonMobile, history.ObjectType)
		if err == nil {
			responses = append(responses, res)
		}
	}
	return
}

func (t *AdTag) compareDataAdTagBE(history mysql.TableHistory) (responses []ResponseCompare) {
	return t.compareDataAdTagFE(history)
}
