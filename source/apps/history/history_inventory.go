package history

import (
	"encoding/json"
	"fmt"
	"source/core/technology/mysql"
)

type Inventory struct {
	Detail    DetailInventory
	CreatorId int64
	RecordOld mysql.TableInventory
	RecordNew mysql.TableInventory
}

func (t *Inventory) Page() string {
	return "Supply"
}

type DetailInventory int

const (
	DetailInventorySubmitFE DetailInventory = iota + 1
	DetailInventoryConfigFE
	DetailInventoryConsentFE
	DetailInventoryAdsTxtFE
	DetailInventoryConnectionFE
	DetailInventoryConnectionBE
	DetailInventoryChangeStatusDomainBE
	DetailInventoryChangeStatusBidderBE
	DetailInventoryConfigBE
	DetailInventoryConsentBE
	DetailInventoryDeleteFE
	DetailInventoryDeleteBE
)

func (t DetailInventory) String() string {
	switch t {
	case DetailInventorySubmitFE:
		return "inventory_submit_fe"
	case DetailInventoryConfigFE:
		return "inventory_config_fe"
	case DetailInventoryAdsTxtFE:
		return "inventory_adstxt_fe"
	case DetailInventoryConnectionFE:
		return "inventory_connection_fe"
	case DetailInventoryConsentFE:
		return "inventory_consent_fe"
	case DetailInventoryChangeStatusDomainBE:
		return "inventory_change_status_domain_be"
	case DetailInventoryChangeStatusBidderBE:
		return "inventory_change_status_bidder_be"
	case DetailInventoryConfigBE:
		return "inventory_config_be"
	case DetailInventoryConsentBE:
		return "inventory_consent_be"
	case DetailInventoryConnectionBE:
		return "inventory_connection_be"
	case DetailInventoryDeleteFE:
		return "inventory_delete_fe"
	case DetailInventoryDeleteBE:
		return "inventory_delete_be"
	}
	return ""
}

func (t DetailInventory) App() string {
	switch t {
	case DetailInventorySubmitFE, DetailInventoryConfigFE, DetailInventoryConsentFE, DetailInventoryAdsTxtFE, DetailInventoryConnectionFE, DetailInventoryDeleteFE:
		return "FE"
	case DetailInventoryChangeStatusDomainBE, DetailInventoryChangeStatusBidderBE, DetailInventoryConfigBE, DetailInventoryConsentBE, DetailInventoryConnectionBE, DetailInventoryDeleteBE:
		return "BE"
	}
	return ""
}

func (t *Inventory) Type() TYPEHistory {
	return TYPEHistoryInventory
}

func (t *Inventory) Action() mysql.TYPEObjectType {
	if t.RecordOld.Id == 0 && t.RecordNew.Id != 0 {
		return mysql.TYPEObjectTypeAdd
	} else if t.RecordOld.Id != 0 && t.RecordNew.Id == 0 {
		return mysql.TYPEObjectTypeDel
	}
	return mysql.TYPEObjectTypeUpdate
}

func (t *Inventory) getRootRecord() (record mysql.TableInventory) {
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

func (t *Inventory) Data() mysql.TableHistory {
	switch t.Detail {
	case DetailInventorySubmitFE:
		return t.getHistoryInventoryFESubmit()
	case DetailInventoryConfigFE:
		return t.getHistoryInventoryFEConfig()
	case DetailInventoryConsentFE:
		return t.getHistoryInventoryFEConsent()
	case DetailInventoryAdsTxtFE:
		return t.getHistoryInventoryAdsTxtFE()
	case DetailInventoryConnectionFE:
		return t.getHistoryInventoryConnectionFE()
	case DetailInventoryChangeStatusDomainBE:
		return t.getHistoryInventoryChangeStatusDomainBE()
	case DetailInventoryChangeStatusBidderBE:
		return t.getHistoryInventoryChangeStatusBidderBE()
	case DetailInventoryConfigBE:
		return t.getHistoryInventoryConfigBE()
	case DetailInventoryConsentBE:
		return t.getHistoryInventoryConsentBE()
	case DetailInventoryConnectionBE:
		return t.getHistoryInventoryConnectionBE()
	case DetailInventoryDeleteFE:
		return t.getHistoryInventoryDeleteFE()
	case DetailInventoryDeleteBE:
		return t.getHistoryInventoryDeleteBE()
	}
	return mysql.TableHistory{}
}

func (t *Inventory) CompareData(history mysql.TableHistory) (res []ResponseCompare) {
	switch history.DetailType {
	case DetailInventorySubmitFE.String():
		return t.compareInventoryFESubmit(history)
	case DetailInventoryConfigFE.String():
		return t.compareInventoryFEConfig(history)
	case DetailInventoryConsentFE.String():
		return t.compareInventoryFEConsent(history)
	case DetailInventoryAdsTxtFE.String():
		return t.compareInventoryAdsTxtFE(history)
	case DetailInventoryConnectionFE.String():
		return t.compareInventoryConnectionFE(history)
	case DetailInventoryChangeStatusDomainBE.String():
		return t.compareInventoryChangeStatusDomainBE(history)
	case DetailInventoryChangeStatusBidderBE.String():
		return t.compareInventoryChangeStatusBidderBE(history)
	case DetailInventoryConfigBE.String():
		return t.compareHistoryInventoryConfigBE(history)
	case DetailInventoryConsentBE.String():
		return t.compareHistoryInventoryConsentBE(history)
	case DetailInventoryConnectionBE.String():
		return t.compareHistoryInventoryConnectionBE(history)
	case DetailInventoryDeleteFE.String():
		return t.compareInventoryDeleteFE(history)
	case DetailInventoryDeleteBE.String():
		return t.compareInventoryDeleteBE(history)
	}
	return []ResponseCompare{}
}

func (t *Inventory) getHistoryInventoryFESubmit() (history mysql.TableHistory) {
	history = mysql.TableHistory{
		CreatorId:   t.CreatorId,
		Title:       "Submit Inventory",
		Object:      mysql.Tables.Inventory,
		ObjectName:  t.getRootRecord().Name,
		ObjectId:    t.getRootRecord().Id,
		ObjectType:  t.Action(),
		DetailType:  t.Detail.String(),
		App:         t.Detail.App(),
		UserId:      t.getRootRecord().UserId,
		InventoryId: t.getRootRecord().Id,
	}

	newData := make(map[string]interface{})
	newData["domain"] = t.RecordNew.Name
	history.CreatedAt = t.RecordNew.CreatedAt

	jsonNewData, _ := json.Marshal(newData)
	history.NewData = string(jsonNewData)
	return
}

func (t *Inventory) getHistoryInventoryDeleteFE() (history mysql.TableHistory) {
	history = mysql.TableHistory{
		CreatorId:   t.CreatorId,
		Title:       "Delete Inventory",
		Object:      mysql.Tables.Inventory,
		ObjectId:    t.getRootRecord().Id,
		ObjectType:  t.Action(),
		DetailType:  t.Detail.String(),
		App:         t.Detail.App(),
		UserId:      t.getRootRecord().UserId,
		InventoryId: t.getRootRecord().Id,
	}

	oldData := make(map[string]interface{})
	oldData["domain"] = t.RecordOld.Name
	history.CreatedAt = t.RecordOld.CreatedAt

	jsonOldData, _ := json.Marshal(oldData)
	history.OldData = string(jsonOldData)
	return
}

type inventoryFEConfig struct {
	GamAutoCreate       *string `json:"gam_auto_create,omitempty"`
	SafeFrame           *string `json:"safe_frame,omitempty"`
	PassbackRenderMode  *string `json:"passback_render_mode,omitempty"`
	PrebidTimeOut       *int    `json:"prebid_time_out,omitempty"`
	LoadAdType          *string `json:"load_ad_type,omitempty"`
	FetchMarginPercent  *int    `json:"fetch_margin_percent,omitempty"`
	RenderMarginPercent *int    `json:"render_margin_percent,omitempty"`
	MobileScaling       *int    `json:"mobile_scaling,omitempty"`
	AdRefresh           *string `json:"ad_refresh,omitempty"`
	AdRefreshTime       *int    `json:"ad_refresh_time,omitempty"`
	AdRefreshType       *string `json:"ad_refresh_type,omitempty"`
	DirectSales         *string `json:"direct_sales,omitempty"`
}

func (t *Inventory) getHistoryInventoryFEConfig() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := inventoryFEConfig{}
	newData := inventoryFEConfig{}
	history = mysql.TableHistory{
		CreatorId:   t.CreatorId,
		Object:      mysql.Tables.Inventory,
		ObjectId:    t.getRootRecord().Id,
		ObjectName:  t.getRootRecord().Name,
		ObjectType:  t.Action(),
		DetailType:  t.Detail.String(),
		App:         t.Detail.App(),
		UserId:      t.getRootRecord().UserId,
		InventoryId: t.getRootRecord().Id,
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
		history.Title = "Add Config"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update Config"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete Config"
		history.OldData = string(bOldData)
	}
	return
}

func (rec *inventoryFEConfig) MakeData(record mysql.TableInventory) {
	gamAuto := record.Config.GamAutoCreate.String()
	rec.GamAutoCreate = &gamAuto

	safeFrame := record.Config.SafeFrame.String()
	rec.SafeFrame = &safeFrame

	pbRenderMode := record.Config.PbRenderMode.String()
	rec.PassbackRenderMode = &pbRenderMode

	rec.PrebidTimeOut = &record.Config.PrebidTimeOut

	rec.LoadAdType = &record.Config.LoadAdType
	if record.Config.LoadAdType == "lazyload" {
		rec.FetchMarginPercent = &record.Config.FetchMarginPercent
		rec.RenderMarginPercent = &record.Config.RenderMarginPercent
		rec.MobileScaling = &record.Config.MobileScaling
	}

	adRefresh := record.Config.AdRefresh.String()
	rec.AdRefresh = &adRefresh
	if record.Config.AdRefresh == mysql.TypeOn {
		rec.AdRefreshTime = &record.Config.AdRefreshTime
		rec.AdRefreshType = &record.Config.LoadAdRefresh
	}

	directSales := record.Config.DirectSales.String()
	rec.DirectSales = &directSales
}

type inventoryFEConsent struct {
	GDPR        *string `json:"gdpr,omitempty"`
	TimeoutGDPR *int    `json:"timeout_gdpr,omitempty"`
	CCPA        *string `json:"ccpa,omitempty"`
	TimeoutCCPA *int    `json:"timeout_ccpa,omitempty"`
	CustomBrand *string `json:"custom_brand,omitempty"`
	Logo        *string `json:"logo,omitempty"`
	Title       *string `json:"title,omitempty"`
	Content     *string `json:"mobile_scaling,omitempty"`
}

func (t *Inventory) getHistoryInventoryFEConsent() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := inventoryFEConsent{}
	newData := inventoryFEConsent{}
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.Inventory,
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
		history.Title = "Add Consent"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update Consent"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete Consent"
		history.OldData = string(bOldData)
	}
	return
}

func (rec *inventoryFEConsent) MakeData(record mysql.TableInventory) {
	gdpr := mysql.TypeOnOff(record.Config.Gdpr).String()
	rec.GDPR = &gdpr
	if mysql.TypeOnOff(record.Config.Gdpr) == mysql.TypeOn {
		rec.TimeoutGDPR = &record.Config.GdprTimeout
	}

	ccpa := mysql.TypeOnOff(record.Config.Ccpa).String()
	rec.CCPA = &ccpa
	if mysql.TypeOnOff(record.Config.Ccpa) == mysql.TypeOn {
		rec.TimeoutCCPA = &record.Config.CcpaTimeout
	}

	customBrand := mysql.TypeOnOff(record.Config.CustomBrand).String()
	rec.CustomBrand = &customBrand
	if mysql.TypeOnOff(record.Config.CustomBrand) == mysql.TypeOn {
		rec.Logo = &record.Config.Logo
		rec.Title = &record.Config.Title
		rec.Content = &record.Config.Content
	}
}

func (t *Inventory) compareInventoryFESubmit(history mysql.TableHistory) (responses []ResponseCompare) {
	recordNew := make(map[string]interface{})
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Submit
	domainName := fmt.Sprintf("%v", recordNew["name"])
	res, err := makeResponseCompare("Submit Domain ", nil, &domainName, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}

func (t *Inventory) compareInventoryDeleteFE(history mysql.TableHistory) (responses []ResponseCompare) {
	recordOld := make(map[string]interface{})
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)

	// Xử lý compare từng row

	// Delete
	domainName := fmt.Sprintf("%v", recordOld["name"])
	res, err := makeResponseCompare("Submit Domain ", &domainName, nil, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}

func (t *Inventory) compareInventoryDeleteBE(history mysql.TableHistory) (responses []ResponseCompare) {
	responses = t.compareInventoryDeleteFE(history)
	return
}

func (t *Inventory) compareInventoryFEConfig(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew inventoryFEConfig
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Direct Sales
	res, err := makeResponseCompare("Direct Sales", recordOld.DirectSales, recordNew.DirectSales, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Gam Auto Create
	res, err = makeResponseCompare("Gam Auto Create", recordOld.GamAutoCreate, recordNew.GamAutoCreate, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Safe Frame
	res, err = makeResponseCompare("Safe Frame", recordOld.SafeFrame, recordNew.SafeFrame, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Passback Render Mode
	res, err = makeResponseCompare("Passback Render Mode", recordOld.PassbackRenderMode, recordNew.PassbackRenderMode, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Prebid Time Out
	res, err = makeResponseCompare("Prebid Timeout", pointerIntToString(recordOld.PrebidTimeOut), pointerIntToString(recordNew.PrebidTimeOut), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Load Ad Type
	res, err = makeResponseCompare("Load Ad Type", recordOld.LoadAdType, recordNew.LoadAdType, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Fetch Margin Percent
	res, err = makeResponseCompare("Fetch Margin Percent", pointerIntToString(recordOld.FetchMarginPercent), pointerIntToString(recordNew.FetchMarginPercent), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Render Margin Percent
	res, err = makeResponseCompare("Render Margin Percent", pointerIntToString(recordOld.RenderMarginPercent), pointerIntToString(recordNew.RenderMarginPercent), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Mobile Scaling
	res, err = makeResponseCompare("Mobile Scaling", pointerIntToString(recordOld.MobileScaling), pointerIntToString(recordNew.MobileScaling), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Ad Refresh
	res, err = makeResponseCompare("Ad Refresh", recordOld.AdRefresh, recordNew.AdRefresh, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Ad Refresh Time
	res, err = makeResponseCompare("Ad Refresh Time", pointerIntToString(recordOld.AdRefreshTime), pointerIntToString(recordNew.AdRefreshTime), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Ad Refresh Type
	res, err = makeResponseCompare("Ad Refresh Type", recordOld.AdRefreshType, recordNew.AdRefreshType, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}

func (t *Inventory) compareInventoryFEConsent(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew inventoryFEConsent
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// GDPR
	res, err := makeResponseCompare("Gam Auto Create", recordOld.GDPR, recordNew.GDPR, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Timeout GDPR
	res, err = makeResponseCompare("Timeout GDPR", pointerIntToString(recordOld.TimeoutGDPR), pointerIntToString(recordNew.TimeoutGDPR), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// CCPA
	res, err = makeResponseCompare("CCPA", recordOld.CCPA, recordNew.CCPA, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Timeout CCPA
	res, err = makeResponseCompare("Timeout CCPA", pointerIntToString(recordOld.TimeoutCCPA), pointerIntToString(recordNew.TimeoutCCPA), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Custom Brand
	res, err = makeResponseCompare("Custom Brand", recordOld.CustomBrand, recordNew.CustomBrand, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Logo
	res, err = makeResponseCompare("Logo", recordOld.Logo, recordNew.Logo, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Title
	res, err = makeResponseCompare("Title", recordOld.Title, recordNew.Title, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Content
	res, err = makeResponseCompare("Content", recordOld.Content, recordNew.Content, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}

func (t *Inventory) getHistoryInventoryAdsTxtFE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := make(map[string]string)
	newData := make(map[string]string)
	history = mysql.TableHistory{
		CreatorId:   t.CreatorId,
		Object:      mysql.Tables.Inventory,
		ObjectId:    t.getRootRecord().Id,
		ObjectName:  t.getRootRecord().Name,
		ObjectType:  t.Action(),
		DetailType:  t.Detail.String(),
		App:         t.Detail.App(),
		UserId:      t.getRootRecord().UserId,
		InventoryId: t.getRootRecord().Id,
	}
	var bNewData, bOldData []byte
	if t.RecordNew.Id != 0 {
		newData["inventory"] = t.RecordNew.Name
		newData["adstxt"] = string(t.RecordNew.AdsTxtCustom)
		bNewData, _ = json.Marshal(newData)
	}
	if t.RecordOld.Id != 0 {
		oldData["inventory"] = t.RecordOld.Name
		oldData["adstxt"] = string(t.RecordOld.AdsTxtCustom)
		bOldData, _ = json.Marshal(oldData)
	}
	history.Title = "Update ads.txt"
	history.NewData = string(bNewData)
	history.OldData = string(bOldData)
	return
}

func (t *Inventory) compareInventoryAdsTxtFE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew map[string]string
	recordOld = make(map[string]string)
	recordNew = make(map[string]string)
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Status
	adstxtOld := recordOld["adstxt"]
	adstxtNew := recordNew["adstxt"]
	res, err := makeResponseCompare(recordNew["inventory"]+" > Ads.txt", &adstxtOld, &adstxtNew, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}

func (t *Inventory) getHistoryInventoryConnectionFE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := make(map[string]string)
	newData := make(map[string]string)
	history = mysql.TableHistory{
		CreatorId:   t.CreatorId,
		Object:      mysql.Tables.InventoryConnectionDemand,
		ObjectId:    t.getRootRecord().ConnectionDemand.InventoryId,
		ObjectType:  t.Action(),
		DetailType:  t.Detail.String(),
		App:         t.Detail.App(),
		UserId:      t.getRootRecord().UserId,
		InventoryId: t.getRootRecord().Id,
		BidderId:    t.getRootRecord().ConnectionDemand.BidderId,
	}
	var bNewData, bOldData []byte
	if t.RecordOld.Id != 0 {
		var bidder mysql.TableBidder
		mysql.Client.Find(&bidder, t.RecordOld.ConnectionDemand.BidderId)
		oldData["inventory"] = t.RecordOld.Name
		oldData["bidder"] = bidder.BidderCode
		oldData["status"] = t.RecordOld.ConnectionDemand.Status.String()
		bOldData, _ = json.Marshal(oldData)
		if bidder.BidderTemplateId == 1 {
			history.ObjectName = bidder.DisplayName
		} else {
			history.ObjectName = bidder.BidderCode
		}
	}
	if t.RecordNew.Id != 0 {
		history.ObjectName = t.RecordNew.Name
		var bidder mysql.TableBidder
		mysql.Client.Find(&bidder, t.RecordNew.ConnectionDemand.BidderId)
		newData["inventory"] = t.RecordNew.Name
		newData["bidder"] = bidder.BidderCode
		newData["status"] = t.RecordNew.ConnectionDemand.Status.String()
		bNewData, _ = json.Marshal(newData)
		if bidder.BidderTemplateId == 1 {
			history.ObjectName = bidder.DisplayName
		} else {
			history.ObjectName = bidder.BidderCode
		}
	}
	history.Title = "Change Status Connection Demand"
	history.NewData = string(bNewData)
	history.OldData = string(bOldData)
	return
}

func (t *Inventory) compareInventoryConnectionFE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew map[string]string
	recordOld = make(map[string]string)
	recordNew = make(map[string]string)
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Status
	statusOld := recordOld["status"]
	statusNew := recordNew["status"]
	res, err := makeResponseCompare(recordNew["inventory"]+" > "+recordNew["bidder"]+" > Status", &statusOld, &statusNew, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}

func (t *Inventory) getHistoryInventoryChangeStatusDomainBE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := make(map[string]string)
	newData := make(map[string]string)
	history = mysql.TableHistory{
		CreatorId:   t.CreatorId,
		Object:      mysql.Tables.Inventory,
		ObjectId:    t.getRootRecord().Id,
		ObjectName:  t.getRootRecord().Name,
		ObjectType:  t.Action(),
		DetailType:  t.Detail.String(),
		App:         t.Detail.App(),
		UserId:      t.getRootRecord().UserId,
		InventoryId: t.getRootRecord().Id,
	}
	var bNewData, bOldData []byte
	if t.RecordNew.Id != 0 {
		newData["status"] = t.RecordNew.Status.String()
		bNewData, _ = json.Marshal(newData)
	}
	if t.RecordOld.Id != 0 {
		oldData["status"] = t.RecordOld.Status.String()
		bOldData, _ = json.Marshal(oldData)
	}
	history.Title = "Change Status"
	history.NewData = string(bNewData)
	history.OldData = string(bOldData)
	return
}

func (t *Inventory) compareInventoryChangeStatusDomainBE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew map[string]string
	recordOld = make(map[string]string)
	recordNew = make(map[string]string)
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Status
	statusOld := recordOld["status"]
	statusNew := recordNew["status"]
	res, err := makeResponseCompare("Status", &statusOld, &statusNew, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}

func (t *Inventory) getHistoryInventoryChangeStatusBidderBE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := make(map[string]string)
	newData := make(map[string]string)
	history = mysql.TableHistory{
		CreatorId:   t.CreatorId,
		Object:      mysql.Tables.Inventory,
		ObjectId:    t.getRootRecord().RlsBidderSystem.Id,
		ObjectType:  t.Action(),
		DetailType:  t.Detail.String(),
		App:         t.Detail.App(),
		UserId:      t.getRootRecord().UserId,
		InventoryId: t.getRootRecord().Id,
	}
	var bNewData, bOldData []byte
	if t.RecordOld.Id != 0 {
		var bidder mysql.TableBidder
		mysql.Client.Find(&bidder, t.RecordOld.RlsBidderSystem.BidderId)
		if bidder.BidderTemplateId == 1 {
			oldData["bidder"] = bidder.DisplayName
		} else {
			oldData["bidder"] = bidder.BidderCode
		}
		oldData["status"] = t.RecordOld.RlsBidderSystem.Status.String()
		bOldData, _ = json.Marshal(oldData)
		history.ObjectName = oldData["bidder"]
		history.BidderId = bidder.Id
	}
	if t.RecordNew.Id != 0 {
		var bidder mysql.TableBidder
		mysql.Client.Find(&bidder, t.RecordNew.RlsBidderSystem.BidderId)
		if bidder.BidderTemplateId == 1 {
			newData["bidder"] = bidder.DisplayName
		} else {
			newData["bidder"] = bidder.BidderCode
		}
		newData["status"] = t.RecordNew.RlsBidderSystem.Status.String()
		bNewData, _ = json.Marshal(newData)
		history.ObjectName = newData["bidder"]
		history.BidderId = bidder.Id
	}
	history.Title = "Change Status Bidder"
	history.NewData = string(bNewData)
	history.OldData = string(bOldData)
	return
}

func (t *Inventory) compareInventoryChangeStatusBidderBE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew map[string]string
	recordOld = make(map[string]string)
	recordNew = make(map[string]string)
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Status
	statusOld := recordOld["status"]
	statusNew := recordNew["status"]
	res, err := makeResponseCompare(recordNew["bidder"]+" > Status", &statusOld, &statusNew, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}

type inventoryConfigBE struct {
	Name                *string `json:"name,omitempty"`
	AdsTxtUrl           *string `json:"ads_txt_url,omitempty"`
	AdsTxtSync          *string `json:"ads_txt_sync,omitempty"`
	RevenueSharing      *string `json:"revenue_sharing,omitempty"`
	Status              *string `json:"status,omitempty"`
	JsMode              *string `json:"js_mode,omitempty"`
	VpaidMode           *string `json:"vpaid_mode,omitempty"`
	PassbackRenderMode  *string `json:"passback_render_mode,omitempty"`
	GAMAdUnitAutoCreate *string `json:"gam_ad_unit_auto_create,omitempty"`
	SafeFrame           *string `json:"safe_frame,omitempty"`
	PrebidTimeOut       *int    `json:"prebid_time_out,omitempty"`
	LoadAdType          *string `json:"load_ad_type,omitempty"`
	FetchMarginPercent  *int    `json:"fetch_margin_percent,omitempty"`
	RenderMarginPercent *int    `json:"render_margin_percent,omitempty"`
	MobileScaling       *int    `json:"mobile_scaling,omitempty"`
	AdRefresh           *string `json:"ad_refresh,omitempty"`
	AdRefreshTime       *int    `json:"ad_refresh_time,omitempty"`
	AdRefreshType       *string `json:"ad_refresh_type,omitempty"`
	CustomAdsTxt        *string `json:"custom_ads_txt,omitempty"`
	DirectSales         *string `json:"direct_sales,omitempty"`
}

func (t *Inventory) getHistoryInventoryConfigBE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := inventoryConfigBE{}
	newData := inventoryConfigBE{}
	history = mysql.TableHistory{
		CreatorId:   t.CreatorId,
		Object:      mysql.Tables.Inventory,
		ObjectId:    t.getRootRecord().Id,
		ObjectName:  t.getRootRecord().Name,
		ObjectType:  t.Action(),
		DetailType:  t.Detail.String(),
		App:         t.Detail.App(),
		UserId:      t.getRootRecord().UserId,
		InventoryId: t.getRootRecord().Id,
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
		history.Title = "Add Config"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update Config"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete Config"
		history.OldData = string(bOldData)
	}
	return
}

func (rec *inventoryConfigBE) MakeData(record mysql.TableInventory) {
	rec.Name = &record.Name

	rec.AdsTxtUrl = &record.AdsTxtUrl

	adsTxtSync := record.SyncAdsTxt.String()
	rec.AdsTxtSync = &adsTxtSync

	customAdsTxt := string(record.AdsTxtCustom)
	rec.CustomAdsTxt = &customAdsTxt

	status := record.Status.String()
	rec.Status = &status

	rec.JsMode = &record.JsMode
	rec.VpaidMode = &record.JsMode

	gamAuto := record.Config.GamAutoCreate.String()
	rec.GAMAdUnitAutoCreate = &gamAuto

	safeFrame := record.Config.SafeFrame.String()
	rec.SafeFrame = &safeFrame

	pbRenderMode := record.Config.PbRenderMode.String()
	rec.PassbackRenderMode = &pbRenderMode

	rec.PrebidTimeOut = &record.Config.PrebidTimeOut

	rec.LoadAdType = &record.Config.LoadAdType
	if record.Config.LoadAdType == "lazyload" {
		rec.FetchMarginPercent = &record.Config.FetchMarginPercent
		rec.RenderMarginPercent = &record.Config.RenderMarginPercent
		rec.MobileScaling = &record.Config.MobileScaling
	}

	adRefresh := record.Config.AdRefresh.String()
	rec.AdRefresh = &adRefresh
	if record.Config.AdRefresh == mysql.TypeOn {
		rec.AdRefreshTime = &record.Config.AdRefreshTime
		rec.AdRefreshType = &record.Config.LoadAdRefresh
	}

	directSales := record.Config.DirectSales.String()
	rec.DirectSales = &directSales
}

func (t *Inventory) compareHistoryInventoryConfigBE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew inventoryConfigBE
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Name
	res, err := makeResponseCompare("Name", recordOld.Name, recordNew.Name, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// AdsTxt Url
	res, err = makeResponseCompare("AdsTxt Url", recordOld.AdsTxtUrl, recordNew.AdsTxtUrl, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// AdsTxt Sync
	res, err = makeResponseCompare("AdsTxt Sync", recordOld.AdsTxtSync, recordNew.AdsTxtSync, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Status
	res, err = makeResponseCompare("Status", recordOld.Status, recordNew.Status, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Js Mode
	res, err = makeResponseCompare("Js Mode", recordOld.JsMode, recordNew.JsMode, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Vpaid Mode
	res, err = makeResponseCompare("Vpaid Mode", recordOld.VpaidMode, recordNew.VpaidMode, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Direct Sales
	res, err = makeResponseCompare("Direct Sales", recordOld.DirectSales, recordNew.DirectSales, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Gam Auto Create
	res, err = makeResponseCompare("Gam Auto Create", recordOld.GAMAdUnitAutoCreate, recordNew.GAMAdUnitAutoCreate, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Safe Frame
	res, err = makeResponseCompare("Safe Frame", recordOld.SafeFrame, recordNew.SafeFrame, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Passback Render Mode
	res, err = makeResponseCompare("Passback Render Mode", recordOld.PassbackRenderMode, recordNew.PassbackRenderMode, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Prebid Time Out
	res, err = makeResponseCompare("Prebid Timeout", pointerIntToString(recordOld.PrebidTimeOut), pointerIntToString(recordNew.PrebidTimeOut), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Load Ad Type
	res, err = makeResponseCompare("Load Ad Type", recordOld.LoadAdType, recordNew.LoadAdType, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Fetch Margin Percent
	res, err = makeResponseCompare("Fetch Margin Percent", pointerIntToString(recordOld.FetchMarginPercent), pointerIntToString(recordNew.FetchMarginPercent), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Render Margin Percent
	res, err = makeResponseCompare("Render Margin Percent", pointerIntToString(recordOld.RenderMarginPercent), pointerIntToString(recordNew.RenderMarginPercent), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Mobile Scaling
	res, err = makeResponseCompare("Mobile Scaling", pointerIntToString(recordOld.MobileScaling), pointerIntToString(recordNew.MobileScaling), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Ad Refresh
	res, err = makeResponseCompare("Ad Refresh", recordOld.AdRefresh, recordNew.AdRefresh, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Ad Refresh Time
	res, err = makeResponseCompare("Ad Refresh Time", pointerIntToString(recordOld.AdRefreshTime), pointerIntToString(recordNew.AdRefreshTime), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Ad Refresh Type
	res, err = makeResponseCompare("Ad Refresh Type", recordOld.AdRefreshType, recordNew.AdRefreshType, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Custom Ads Txt
	res, err = makeResponseCompare("Custom AdsTxt", recordOld.CustomAdsTxt, recordNew.CustomAdsTxt, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}

func (t *Inventory) getHistoryInventoryConsentBE() (history mysql.TableHistory) {
	history = t.getHistoryInventoryFEConsent()
	return
}

func (t *Inventory) getHistoryInventoryDeleteBE() (history mysql.TableHistory) {
	history = t.getHistoryInventoryDeleteFE()
	return
}

func (t *Inventory) compareHistoryInventoryConsentBE(history mysql.TableHistory) (responses []ResponseCompare) {
	responses = t.compareInventoryFEConsent(history)
	return
}

func (t *Inventory) getHistoryInventoryConnectionBE() (history mysql.TableHistory) {
	history = t.getHistoryInventoryConnectionFE()
	return
}

func (t *Inventory) compareHistoryInventoryConnectionBE(history mysql.TableHistory) (responses []ResponseCompare) {
	responses = t.compareInventoryConnectionFE(history)
	return
}
