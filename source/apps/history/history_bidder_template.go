package history

import (
	"encoding/json"
	"errors"
	"fmt"
	"source/core/technology/mysql"
)

type BidderTemplate struct {
	Detail    DetailBidderTemplate
	CreatorId int64
	RecordOld mysql.TableBidderTemplate
	RecordNew mysql.TableBidderTemplate
}

func (t *BidderTemplate) Page() string {
	return "Bidder Template"
}

type DetailBidderTemplate int

const (
	DetailBidderTemplateBE DetailBidderTemplate = iota + 1
)

func (t DetailBidderTemplate) String() string {
	switch t {
	case DetailBidderTemplateBE:
		return "bidder_template_be"
	}
	return ""
}

func (t DetailBidderTemplate) App() string {
	switch t {
	case DetailBidderTemplateBE:
		return "BE"
	}
	return ""
}

func (t *BidderTemplate) Type() TYPEHistory {
	return TYPEHistoryBidderTemplate
}

func (t *BidderTemplate) Action() mysql.TYPEObjectType {
	if t.RecordOld.Id == 0 && t.RecordNew.Id != 0 {
		return mysql.TYPEObjectTypeAdd
	} else if t.RecordOld.Id != 0 && t.RecordNew.Id == 0 {
		return mysql.TYPEObjectTypeDel
	}
	return mysql.TYPEObjectTypeUpdate
}

func (t *BidderTemplate) Data() mysql.TableHistory {
	switch t.Detail {
	case DetailBidderTemplateBE:
		return t.getHistoryBidderTemplateBE()
	}
	return mysql.TableHistory{}
}

func (t *BidderTemplate) CompareData(history mysql.TableHistory) (res []ResponseCompare) {
	switch history.DetailType {
	case DetailBidderTemplateBE.String():
		return t.compareDataBidderTemplateBE(history)
	}
	return []ResponseCompare{}
}

func (t *BidderTemplate) getRootRecord() (record mysql.TableBidderTemplate) {
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

type bidderTemplateBE struct {
	PrebidModule  *string        `json:"prebid_module,omitempty"`
	DisplayName   *string        `json:"display_name,omitempty"`
	BidderCode    *string        `json:"bidder_code,omitempty"`
	BidderAlias   *string        `json:"bidder_alias,omitempty"`
	MediaTypes    *[]string      `json:"media_types,omitempty"`
	BidAdjustment *string        `json:"bid_adjustment,omitempty"`
	Logo          *string        `json:"logo,omitempty"`
	Params        *[]bidderParam `json:"params,omitempty"`
}

func (t *BidderTemplate) getHistoryBidderTemplateBE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := bidderTemplateBE{}
	newData := bidderTemplateBE{}
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.BidderTemplate,
		ObjectId:   t.getRootRecord().Id,
		ObjectName: t.getRootRecord().BidderCode,
		ObjectType: t.Action(),
		DetailType: t.Detail.String(),
		App:        t.Detail.App(),
		UserId:     t.getRootRecord().UserId,
		BidderId:   t.getRootRecord().Id,
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
		history.Title = "Add Bidder Template"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update Bidder Template"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete Bidder Template"
		history.OldData = string(bOldData)
	}
	return
}

func (rec *bidderTemplateBE) MakeData(record mysql.TableBidderTemplate) {
	rec.PrebidModule = &record.PrebidModule
	rec.DisplayName = &record.DisplayName
	rec.BidderCode = &record.BidderCode
	rec.BidderAlias = &record.BidderAlias
	bidAdjustment := fmt.Sprintf("%f", record.BidAdjustment)
	rec.BidAdjustment = &bidAdjustment
	rec.Logo = &record.Logo
	var mediaTypes []string
	for _, rlsBidderMediaType := range record.MediaTypes {
		var mediaType mysql.TableMediaType
		mysql.Client.Find(&mediaType, rlsBidderMediaType.MediaTypeId)
		mediaTypes = append(mediaTypes, mediaType.Name)
	}
	rec.MediaTypes = &mediaTypes
	var params []bidderParam
	for _, param := range record.Params {
		params = append(params, bidderParam{
			Key:  param.Name,
			Type: param.Type,
		})
	}
	rec.Params = &params
}

func (t *BidderTemplate) compareDataBidderTemplateBE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew bidderTemplateBE
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Prebid Module
	res, err := makeResponseCompare("Prebid Module", recordOld.PrebidModule, recordNew.PrebidModule, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Display Name
	res, err = makeResponseCompare("Display Name", recordOld.DisplayName, recordNew.DisplayName, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Bidder Code
	res, err = makeResponseCompare("Bidder Code", recordOld.BidderCode, recordNew.BidderCode, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Bidder Alias
	res, err = makeResponseCompare("Bidder Alias", recordOld.BidderAlias, recordNew.BidderAlias, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Media Types
	res, err = makeResponseCompare("Media Types", pointerArrayStringToString(recordOld.MediaTypes), pointerArrayStringToString(recordNew.MediaTypes), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// BidAdjustment
	res, err = makeResponseCompare("Bid Adjustment", recordOld.BidAdjustment, recordNew.BidAdjustment, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Logo
	res, err = makeResponseCompare("Logo", recordOld.Logo, recordNew.Logo, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Param
	responseParam, err := t.makeResponseCompareParam(recordOld.Params, recordNew.Params)
	if err == nil {
		responses = append(responses, responseParam...)
	}
	return
}

func (t *BidderTemplate) makeResponseCompareParam(oldData, newData *[]bidderParam) (responses []ResponseCompare, err error) {
	if oldData == nil && newData == nil {
		err = errors.New("no response")
		return
	}
	mapOldData := make(map[string]bidderParam)
	mapNewData := make(map[string]bidderParam)

	if oldData != nil {
		for _, param := range *oldData {
			mapOldData[param.Key] = param
		}
	}
	if newData != nil {
		for _, param := range *newData {
			mapNewData[param.Key] = param
		}
	}

	for name, paramOld := range mapOldData {
		if paramNew, ok := mapNewData[name]; ok {
			// Type
			res, err := makeResponseCompare("Params > "+name+" > Type", &paramOld.Type, &paramNew.Type, mysql.TYPEObjectTypeUpdate)
			if err == nil {
				responses = append(responses, res)
			}

			// Xóa các param update
			delete(mapNewData, name)
		} else {
			// Nếu chỉ tồn tại oldData tức param này đã bị xóa toàn bộ các newData là nil
			// Type
			res, err := makeResponseCompare("Params > "+name+" > Type", &paramOld.Type, nil, mysql.TYPEObjectTypeDel)
			if err == nil {
				responses = append(responses, res)
			}
		}
	}
	// Các dataNew còn lại đều là các param add mới
	for name, paramNew := range mapNewData {
		// Type
		res, err := makeResponseCompare("Params > "+name+" > Type", nil, &paramNew.Type, mysql.TYPEObjectTypeAdd)
		if err == nil {
			responses = append(responses, res)
		}
	}
	return
}
