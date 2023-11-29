package history2

import (
	"encoding/json"
	"source/internal/entity/model"
	"source/internal/errors"
	"strconv"
)

type historyInventory struct {
	*historyUsecase
}

func newHistoryInventory(historyUsecase *historyUsecase) *historyInventory {
	return &historyInventory{historyUsecase: historyUsecase}
}

func (t *historyInventory) push(input *inputPush) (err error) {
	oldRecord := input.oldData.(model.InventoryModel)
	newRecord := input.newData.(model.InventoryModel)
	var rootRecord model.InventoryModel
	if newRecord.ID != 0 {
		rootRecord = newRecord
	} else {
		rootRecord = oldRecord
	}
	record := model.HistoryModel{
		CreatorID:  input.creatorID,
		UserID:     rootRecord.UserID,
		Title:      input.title,
		Object:     rootRecord.TableName(),
		ObjectID:   rootRecord.ID,
		ObjectName: rootRecord.Name,
		ObjectType: input.objectType,
		DetailType: input.detailType,
		OldData:    oldRecord.ToJSON(),
		NewData:    newRecord.ToJSON(),
	}
	if input.objectType == model.HistoryObjectTYPEAdd {
		record.CreatedAt = rootRecord.CreatedAt
	}
	return t.Repos.History.Save(&record)
}

func (t *historyInventory) validate(push *inputPush) (err error) {
	if push.oldData == nil && push.newData == nil {
		return errors.New("empty data")
	}
	return
}

func (t *historyInventory) compare(record *model.HistoryModel) (resp []model.CompareData, err error) {
	switch model.HistoryDetailTypeTYPE(record.DetailType) {
	case model.DetailInventorySubmitFE:
		resp, err = t.compareInventorySubmitFE(record)
		break
	case model.DetailInventoryConfigFE:
		resp, err = t.compareInventoryConfigFE(record)
		break
	case model.DetailInventoryConsentFE:
		resp, err = t.compareInventoryConsentFE(record)
		break
	case model.DetailInventoryAdsTxtFE:
		resp, err = t.compareInventoryAdsTxtFE(record)
		break
	}
	return
}

func (t *historyInventory) filter() {
	// ở đây em xác định được là bảng inventory rồi => em chỉ cần where("object_id = ?",inventory.ID) là được mà
	return
}

func (t *historyInventory) compareInventorySubmitFE(record *model.HistoryModel) (resp []model.CompareData, err error) {
	var recordNew model.InventoryModel
	_ = json.Unmarshal([]byte(record.NewData), &recordNew)

	// Xử lý compare từng row

	// Submit
	resp = append(resp, model.CompareData{
		Action:  record.ObjectType.String(),
		Text:    "Submit Domain",
		OldData: "",
		NewData: recordNew.Name,
	})
	return
}

func (t *historyInventory) compareInventoryConfigFE(record *model.HistoryModel) (resp []model.CompareData, err error) {
	var recordOld, recordNew model.InventoryModel
	_ = json.Unmarshal([]byte(record.OldData), &recordOld)
	_ = json.Unmarshal([]byte(record.NewData), &recordNew)

	// Xử lý compare từng row

	//=> Gam Auto Create
	gamAutoCreateOld := recordOld.Config.GamAutoCreate.String()
	gamAutoCreateNew := recordNew.Config.GamAutoCreate.String()
	res, err := makeResponseCompare("Gam Auto Create", &gamAutoCreateOld, &gamAutoCreateNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	//=> Safe Frame
	safeFrameOld := recordOld.Config.SafeFrame.String()
	safeFrameNew := recordNew.Config.SafeFrame.String()
	res, err = makeResponseCompare("Safe Frame", &safeFrameOld, &safeFrameNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	//=> Passback Render Mode
	passBackRenderModeOld := recordOld.Config.PbRenderMode.String()
	passBackRenderModeNew := recordNew.Config.PbRenderMode.String()
	res, err = makeResponseCompare("Passback Render Mode", &passBackRenderModeOld, &passBackRenderModeNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	//=> Prebid Timeout
	prebidTimeOutOld := strconv.Itoa(recordOld.Config.PrebidTimeOut)
	prebidTimeOutNew := strconv.Itoa(recordNew.Config.PrebidTimeOut)
	res, err = makeResponseCompare("Prebid Timeout", &prebidTimeOutOld, &prebidTimeOutNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	//=> Ad Refresh
	adRefreshOld := recordOld.Config.AdRefresh.String()
	adRefreshNew := recordNew.Config.AdRefresh.String()
	res, err = makeResponseCompare("Prebid Timeout", &adRefreshOld, &adRefreshNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	//=> Ad Refresh Time
	var adRefreshTimeOld, adRefreshTimeNew *string
	if recordOld.Config.AdRefresh != recordNew.Config.AdRefresh {
		if recordNew.Config.AdRefresh == model.On {
			adRefreshTimeOld = nil

			adRefreshTimeNewConvert := strconv.Itoa(recordNew.Config.AdRefreshTime)
			adRefreshTimeNew = &adRefreshTimeNewConvert
		} else {
			adRefreshTimeNew = nil

			adRefreshTimeOldConvert := strconv.Itoa(recordOld.Config.AdRefreshTime)
			adRefreshTimeOld = &adRefreshTimeOldConvert
		}
	} else {
		adRefreshTimeOldConvert := strconv.Itoa(recordOld.Config.AdRefreshTime)
		adRefreshTimeOld = &adRefreshTimeOldConvert
		adRefreshTimeNewConvert := strconv.Itoa(recordNew.Config.AdRefreshTime)
		adRefreshTimeNew = &adRefreshTimeNewConvert
	}
	res, err = makeResponseCompare("Ad Refresh Time", adRefreshTimeOld, adRefreshTimeNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	//=> Ad Refresh Type
	var adRefreshTypeOld, adRefreshTypeNew *string
	if recordOld.Config.AdRefresh != recordNew.Config.AdRefresh {
		if recordNew.Config.AdRefresh == model.On {
			adRefreshTypeOld = nil

			adRefreshTypeNewConvert := recordNew.Config.LoadAdRefresh
			adRefreshTypeNew = &adRefreshTypeNewConvert
		} else {
			adRefreshTypeNew = nil

			adRefreshTypeOldConvert := recordOld.Config.LoadAdRefresh
			adRefreshTypeOld = &adRefreshTypeOldConvert
		}
	} else {
		adRefreshTypeOldConvert := recordOld.Config.LoadAdRefresh
		adRefreshTypeOld = &adRefreshTypeOldConvert
		adRefreshTypeNewConvert := recordNew.Config.LoadAdRefresh
		adRefreshTypeNew = &adRefreshTypeNewConvert
	}
	res, err = makeResponseCompare("Ad Refresh Type", adRefreshTypeOld, adRefreshTypeNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	//=> Load Ad Type
	loadAdTypeOld := recordOld.Config.LoadAdType
	loadAdTypeNew := recordNew.Config.LoadAdType
	res, err = makeResponseCompare("Load Ad Type", &loadAdTypeOld, &loadAdTypeNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	//=> Fetch Margin Percent
	var fetchMarginPercentOld, fetchMarginPercentNew *string
	if recordOld.Config.LoadAdType != recordNew.Config.LoadAdType {
		if recordNew.Config.LoadAdType == "lazyload" {
			fetchMarginPercentOld = nil

			fetchMarginPercentNewConvert := strconv.Itoa(recordNew.Config.FetchMarginPercent)
			fetchMarginPercentNew = &fetchMarginPercentNewConvert
		} else {
			adRefreshTypeNew = nil

			fetchMarginPercentOldConvert := strconv.Itoa(recordOld.Config.FetchMarginPercent)
			fetchMarginPercentOld = &fetchMarginPercentOldConvert
		}
	} else {
		fetchMarginPercentOldConvert := strconv.Itoa(recordOld.Config.FetchMarginPercent)
		fetchMarginPercentOld = &fetchMarginPercentOldConvert
		fetchMarginPercentNewConvert := strconv.Itoa(recordNew.Config.FetchMarginPercent)
		fetchMarginPercentNew = &fetchMarginPercentNewConvert
	}
	res, err = makeResponseCompare("Fetch Margin Percent", fetchMarginPercentOld, fetchMarginPercentNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	//=> Render Margin Percent
	var renderMarginPercentOld, renderMarginPercentNew *string
	if recordOld.Config.LoadAdType != recordNew.Config.LoadAdType {
		if recordNew.Config.LoadAdType == "lazyload" {
			renderMarginPercentOld = nil

			renderMarginPercentNewConvert := strconv.Itoa(recordNew.Config.RenderMarginPercent)
			renderMarginPercentNew = &renderMarginPercentNewConvert
		} else {
			adRefreshTypeNew = nil

			renderMarginPercentOldConvert := strconv.Itoa(recordOld.Config.RenderMarginPercent)
			renderMarginPercentOld = &renderMarginPercentOldConvert
		}
	} else {
		renderMarginPercentOldConvert := strconv.Itoa(recordOld.Config.RenderMarginPercent)
		renderMarginPercentOld = &renderMarginPercentOldConvert
		renderMarginPercentNewConvert := strconv.Itoa(recordNew.Config.RenderMarginPercent)
		renderMarginPercentNew = &renderMarginPercentNewConvert
	}
	res, err = makeResponseCompare("Render Margin Percent", renderMarginPercentOld, renderMarginPercentNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	//=> Mobile Scaling
	var mobileScalingOld, mobileScalingNew *string
	if recordOld.Config.LoadAdType != recordNew.Config.LoadAdType {
		if recordNew.Config.LoadAdType == "lazyload" {
			mobileScalingOld = nil

			mobileScalingNewConvert := strconv.Itoa(recordNew.Config.MobileScaling)
			mobileScalingNew = &mobileScalingNewConvert
		} else {
			adRefreshTypeNew = nil

			mobileScalingOldConvert := strconv.Itoa(recordOld.Config.MobileScaling)
			mobileScalingOld = &mobileScalingOldConvert
		}
	} else {
		mobileScalingOldConvert := strconv.Itoa(recordOld.Config.MobileScaling)
		mobileScalingOld = &mobileScalingOldConvert
		mobileScalingNewConvert := strconv.Itoa(recordNew.Config.MobileScaling)
		mobileScalingNew = &mobileScalingNewConvert
	}
	res, err = makeResponseCompare("Mobile Scaling", mobileScalingOld, mobileScalingNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}
	return
}

func (t *historyInventory) compareInventoryConsentFE(record *model.HistoryModel) (resp []model.CompareData, err error) {
	var recordOld, recordNew model.InventoryModel
	_ = json.Unmarshal([]byte(record.OldData), &recordOld)
	_ = json.Unmarshal([]byte(record.NewData), &recordNew)

	// Xử lý compare từng row

	//=> GDPR
	gdprOld := recordOld.Config.Gdpr.String()
	gdprNew := recordNew.Config.Gdpr.String()
	res, err := makeResponseCompare("GDPR", &gdprOld, &gdprNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	//=> GDPR Timeout
	var gdprTimeoutOld, gdprTimeoutNew *string
	if recordOld.Config.Gdpr != recordNew.Config.Gdpr {
		if recordNew.Config.AdRefresh == model.On {
			gdprTimeoutOld = nil

			gdprTimeoutNewConvert := strconv.Itoa(recordNew.Config.GdprTimeout)
			gdprTimeoutNew = &gdprTimeoutNewConvert
		} else {
			gdprTimeoutNew = nil

			gdprTimeoutOldConvert := strconv.Itoa(recordOld.Config.GdprTimeout)
			gdprTimeoutOld = &gdprTimeoutOldConvert
		}
	} else {
		gdprTimeoutOldConvert := strconv.Itoa(recordOld.Config.GdprTimeout)
		gdprTimeoutOld = &gdprTimeoutOldConvert
		gdprTimeoutNewConvert := strconv.Itoa(recordNew.Config.GdprTimeout)
		gdprTimeoutNew = &gdprTimeoutNewConvert
	}
	res, err = makeResponseCompare("GDPR Timeout", gdprTimeoutOld, gdprTimeoutNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	//=> CCPA
	ccpaOld := recordOld.Config.Ccpa.String()
	ccpaNew := recordNew.Config.Ccpa.String()
	res, err = makeResponseCompare("CCPA", &ccpaOld, &ccpaNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	//=> CCPA Timeout
	var ccpaTimeoutOld, ccpaTimeoutNew *string
	if recordOld.Config.Ccpa != recordNew.Config.Ccpa {
		if recordNew.Config.AdRefresh == model.On {
			ccpaTimeoutOld = nil

			ccpaTimeoutNewConvert := strconv.Itoa(recordNew.Config.CcpaTimeout)
			ccpaTimeoutNew = &ccpaTimeoutNewConvert
		} else {
			ccpaTimeoutNew = nil

			ccpaTimeoutOldConvert := strconv.Itoa(recordOld.Config.CcpaTimeout)
			ccpaTimeoutOld = &ccpaTimeoutOldConvert
		}
	} else {
		ccpaTimeoutOldConvert := strconv.Itoa(recordOld.Config.CcpaTimeout)
		ccpaTimeoutOld = &ccpaTimeoutOldConvert
		ccpaTimeoutNewConvert := strconv.Itoa(recordNew.Config.CcpaTimeout)
		ccpaTimeoutNew = &ccpaTimeoutNewConvert
	}
	res, err = makeResponseCompare("CCPA Timeout", ccpaTimeoutOld, ccpaTimeoutNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	//=> Custom Brand
	customBrandOld := recordOld.Config.CustomBrand.String()
	customBrandNew := recordNew.Config.CustomBrand.String()
	res, err = makeResponseCompare("Custom Brand", &customBrandOld, &customBrandNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	//=> Logo
	var logoOld, logoNew *string
	if recordOld.Config.CustomBrand != recordNew.Config.CustomBrand {
		if recordNew.Config.AdRefresh == model.On {
			logoOld = nil

			logoNewConvert := recordNew.Config.Logo
			logoNew = &logoNewConvert
		} else {
			logoNew = nil

			logoOldConvert := recordOld.Config.Logo
			logoOld = &logoOldConvert
		}
	} else {
		logoOldConvert := recordOld.Config.Logo
		logoOld = &logoOldConvert
		logoNewConvert := recordNew.Config.Logo
		logoNew = &logoNewConvert
	}
	res, err = makeResponseCompare("Logo", logoOld, logoNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	//=> Title
	var titleOld, titleNew *string
	if recordOld.Config.CustomBrand != recordNew.Config.CustomBrand {
		if recordNew.Config.AdRefresh == model.On {
			titleOld = nil

			titleNewConvert := recordNew.Config.Title
			titleNew = &titleNewConvert
		} else {
			titleNew = nil

			titleOldConvert := recordOld.Config.Title
			titleOld = &titleOldConvert
		}
	} else {
		titleOldConvert := recordOld.Config.Title
		titleOld = &titleOldConvert
		titleNewConvert := recordNew.Config.Title
		titleNew = &titleNewConvert
	}
	res, err = makeResponseCompare("Title", titleOld, titleNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	//=> Content
	var contentOld, contentNew *string
	if recordOld.Config.CustomBrand != recordNew.Config.CustomBrand {
		if recordNew.Config.AdRefresh == model.On {
			contentOld = nil

			contentNewConvert := recordNew.Config.Content
			contentNew = &contentNewConvert
		} else {
			contentNew = nil

			contentOldConvert := recordOld.Config.Content
			contentOld = &contentOldConvert
		}
	} else {
		contentOldConvert := recordOld.Config.Content
		contentOld = &contentOldConvert
		contentNewConvert := recordNew.Config.Content
		contentNew = &contentNewConvert
	}
	res, err = makeResponseCompare("Content", contentOld, contentNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}

	return
}

func (t *historyInventory) compareInventoryAdsTxtFE(record *model.HistoryModel) (resp []model.CompareData, err error) {
	var recordOld, recordNew model.InventoryModel
	_ = json.Unmarshal([]byte(record.OldData), &recordOld)
	_ = json.Unmarshal([]byte(record.NewData), &recordNew)

	// Xử lý compare từng row

	//=> Ads.Txt
	adsTxtOld := string(recordOld.AdsTxtCustom)
	adsTxtNew := string(recordNew.AdsTxtCustom)
	res, err := makeResponseCompare("Ads.Txt", &adsTxtOld, &adsTxtNew, record.ObjectType)
	if err == nil {
		resp = append(resp, res)
	}
	return
}
