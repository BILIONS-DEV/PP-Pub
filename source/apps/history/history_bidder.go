package history

import (
	"encoding/json"
	"errors"
	"fmt"
	"source/core/technology/mysql"
	"strconv"
)

type Bidder struct {
	Detail    DetailBidder
	CreatorId int64
	RecordOld mysql.TableBidder
	RecordNew mysql.TableBidder
}

func (t *Bidder) Page() string {
	return "Bidder"
}

type DetailBidder int

const (
	DetailBidderFE DetailBidder = iota + 1
	DetailBidderAdxConnectionMcmBE
	DetailBidderBE
)

func (t DetailBidder) String() string {
	switch t {
	case DetailBidderFE:
		return "bidder_fe"
	case DetailBidderAdxConnectionMcmBE:
		return "bidder_adx_connection_mcm_be"
	case DetailBidderBE:
		return "bidder_be"
	}
	return ""
}

func (t DetailBidder) App() string {
	switch t {
	case DetailBidderFE:
		return "FE"
	case DetailBidderAdxConnectionMcmBE, DetailBidderBE:
		return "BE"
	}
	return ""
}

func (t *Bidder) Type() TYPEHistory {
	return TYPEHistoryBidder
}

func (t *Bidder) Action() mysql.TYPEObjectType {
	if t.RecordOld.Id == 0 && t.RecordNew.Id != 0 {
		return mysql.TYPEObjectTypeAdd
	} else if t.RecordOld.Id != 0 && t.RecordNew.Id == 0 {
		return mysql.TYPEObjectTypeDel
	}
	return mysql.TYPEObjectTypeUpdate
}

func (t *Bidder) Data() mysql.TableHistory {
	switch t.Detail {
	case DetailBidderFE:
		return t.getHistoryBidderFE()
	case DetailBidderAdxConnectionMcmBE:
		return t.getHistoryBidderAdxConnectionMcmBE()
	case DetailBidderBE:
		return t.getHistoryBidderBE()
	}
	return mysql.TableHistory{}
}

func (t *Bidder) CompareData(history mysql.TableHistory) (res []ResponseCompare) {
	switch history.DetailType {
	case DetailBidderFE.String():
		return t.compareDataBidderFE(history)
	case DetailBidderAdxConnectionMcmBE.String():
		return t.compareDataBidderAdxConnectionMcmBE(history)
	case DetailBidderBE.String():
		return t.compareDataBidderBE(history)
	}
	return []ResponseCompare{}
}

func (t *Bidder) getRootRecord() (record mysql.TableBidder) {
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

type bidderFE struct {
	Bidder            *string        `json:"bidder,omitempty"`
	AccountType       *string        `json:"account_type,omitempty"`
	DisplayName       *string        `json:"display_name,omitempty"`
	MediaTypes        *[]string      `json:"media_types,omitempty"`
	BidAdjustment     *float64       `json:"bid_adjustment,omitempty"`
	RPM               *string        `json:"rpm,omitempty"`
	AdsTxt            *string        `json:"ads_txt,omitempty"`
	PubId             *string        `json:"pub_id,omitempty"`
	Params            *[]bidderParam `json:"params,omitempty"`
	SupplyChainDomain *string        `json:"supply_chain_domain"`
}

type bidderParam struct {
	Key  string `json:"key"`
	Type string `json:"type"`
}

func (t *Bidder) getHistoryBidderFE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := bidderFE{}
	newData := bidderFE{}
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.Bidder,
		ObjectId:   t.getRootRecord().Id,
		ObjectType: t.Action(),
		DetailType: t.Detail.String(),
		App:        t.Detail.App(),
		UserId:     t.getRootRecord().UserId,
		BidderId:   t.getRootRecord().Id,
	}
	if t.getRootRecord().BidderTemplateId == 1 {
		history.ObjectName = t.getRootRecord().DisplayName
	} else {
		history.ObjectName = t.getRootRecord().BidderCode
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
		history.Title = "Add Bidder"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update Bidder"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete Bidder"
		history.OldData = string(bOldData)
	}
	return
}

func (rec *bidderFE) MakeData(record mysql.TableBidder) {
	var user mysql.TableUser
	mysql.Client.Find(&user, record.UserId)

	rec.Bidder = &record.BidderCode
	rpm := fmt.Sprintf("%f", record.RPM)
	rec.RPM = &rpm
	if user.Permission == mysql.UserPermissionNetwork {
		rec.SupplyChainDomain = &record.SellerDomain
	}
	if record.BidderTemplateId == 1 { // Xử lý khi bidder là google
		accountType := record.AccountType.String()
		rec.AccountType = &accountType
		rec.DisplayName = &record.DisplayName
		var mediaTypes []string
		for _, rlsBidderMediaType := range record.MediaTypes {
			var mediaType mysql.TableMediaType
			mysql.Client.Find(&mediaType, rlsBidderMediaType.MediaTypeId)
			mediaTypes = append(mediaTypes, mediaType.Name)
		}
		rec.MediaTypes = &mediaTypes
		rec.BidAdjustment = record.BidAdjustment
		rec.AdsTxt = &record.AdsTxt
		rec.PubId = &record.PubId
	} else if record.BidderTemplateId == 2 { // Bidder là amz
		var mediaTypes []string
		for _, rlsBidderMediaType := range record.MediaTypes {
			var mediaType mysql.TableMediaType
			mysql.Client.Find(&mediaType, rlsBidderMediaType.MediaTypeId)
			mediaTypes = append(mediaTypes, mediaType.Name)
		}
		rec.MediaTypes = &mediaTypes
		rec.BidAdjustment = record.BidAdjustment
		rec.AdsTxt = &record.AdsTxt
		rec.PubId = &record.PubId
		var params []bidderParam
		for _, param := range record.Params {
			params = append(params, bidderParam{
				Key:  param.Name,
				Type: param.Type,
			})
		}
		rec.Params = &params
	} else {
		var mediaTypes []string
		for _, rlsBidderMediaType := range record.MediaTypes {
			var mediaType mysql.TableMediaType
			mysql.Client.Find(&mediaType, rlsBidderMediaType.MediaTypeId)
			mediaTypes = append(mediaTypes, mediaType.Name)
		}
		rec.MediaTypes = &mediaTypes
		rec.BidAdjustment = record.BidAdjustment
		rec.AdsTxt = &record.AdsTxt
		var params []bidderParam
		for _, param := range record.Params {
			params = append(params, bidderParam{
				Key:  param.Name,
				Type: param.Type,
			})
		}
		rec.Params = &params
	}
}

func (t *Bidder) compareDataBidderFE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew bidderFE
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Bidder
	res, err := makeResponseCompare("Bidder", recordOld.Bidder, recordNew.Bidder, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Account Type
	res, err = makeResponseCompare("Account Type", recordOld.AccountType, recordNew.AccountType, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Display Name
	res, err = makeResponseCompare("Display Name", recordOld.DisplayName, recordNew.DisplayName, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Media Types
	res, err = makeResponseCompare("Media Types", pointerArrayStringToString(recordOld.MediaTypes), pointerArrayStringToString(recordNew.MediaTypes), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// BidAdjustment
	res, err = makeResponseCompare("Bid Adjustment", pointerFloatToString(recordOld.BidAdjustment), pointerFloatToString(recordNew.BidAdjustment), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// RPM
	res, err = makeResponseCompare("RPM", recordOld.RPM, recordNew.RPM, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// AdsTxt
	res, err = makeResponseCompare("AdsTxt", recordOld.AdsTxt, recordNew.AdsTxt, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Pub Id
	res, err = makeResponseCompare("Pub Id", recordOld.PubId, recordNew.PubId, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Supply Chain Domain
	res, err = makeResponseCompare("Supply Chain Domain", recordOld.SupplyChainDomain, recordNew.SupplyChainDomain, history.ObjectType)
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

func (t *Bidder) makeResponseCompareParam(oldData, newData *[]bidderParam) (responses []ResponseCompare, err error) {
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

func (t *Bidder) getHistoryBidderAdxConnectionMcmBE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := make(map[string]string)
	newData := make(map[string]string)
	history = mysql.TableHistory{
		CreatorId:   t.CreatorId,
		Object:      mysql.Tables.RlsConnectionMCM,
		ObjectId:    t.getRootRecord().Id,
		ObjectType:  t.Action(),
		DetailType:  t.Detail.String(),
		App:         t.Detail.App(),
		UserId:      t.getRootRecord().UserId,
		InventoryId: t.getRootRecord().Id,
	}
	if t.getRootRecord().BidderTemplateId == 1 {
		history.ObjectName = t.getRootRecord().DisplayName
	} else {
		history.ObjectName = t.getRootRecord().BidderCode
	}
	var bNewData, bOldData []byte
	if t.RecordNew.Id != 0 {
		var connectionMCM mysql.TableRlsConnectionMCM
		mysql.Client.Find(&connectionMCM, t.RecordNew.Id)
		newData["bidderName"] = t.RecordNew.DisplayName
		newData["networkId"] = strconv.FormatInt(t.RecordNew.RlsConnectionMCM.NetworkId, 10)
		newData["status"] = t.RecordNew.RlsConnectionMCM.Status.String()
		bNewData, _ = json.Marshal(newData)
	}
	if t.RecordOld.Id != 0 {
		var connectionMCM mysql.TableRlsConnectionMCM
		mysql.Client.Find(&connectionMCM, t.RecordOld.Id)
		oldData["bidderName"] = t.RecordOld.DisplayName
		oldData["networkId"] = strconv.FormatInt(t.RecordOld.RlsConnectionMCM.NetworkId, 10)
		oldData["status"] = t.RecordOld.RlsConnectionMCM.Status.String()
		bOldData, _ = json.Marshal(oldData)
	}
	history.Title = "Change Status Connection MCM"
	history.NewData = string(bNewData)
	history.OldData = string(bOldData)
	return
}

func (t *Bidder) compareDataBidderAdxConnectionMcmBE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew map[string]string
	recordOld = make(map[string]string)
	recordNew = make(map[string]string)
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Status
	statusOld := recordOld["status"]
	statusNew := recordNew["status"]
	res, err := makeResponseCompare(recordNew["bidderName"]+" > "+recordNew["networkId"]+" > Status", &statusOld, &statusNew, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}

type bidderBE struct {
	Bidder           *string        `json:"bidder,omitempty"`
	AccountType      *string        `json:"account_type,omitempty"`
	DisplayName      *string        `json:"display_name,omitempty"`
	AliasCode        *string        `json:"alias_code,omitempty"`
	AliasName        *string        `json:"alias_name,omitempty"`
	ShowOnPub        *string        `json:"show_on_pub,omitempty"`
	MediaTypes       *[]string      `json:"media_types,omitempty"`
	BidAdjustment    *string        `json:"bid_adjustment,omitempty"`
	RPM              *string        `json:"rpm,omitempty"`
	AdsTxt           *string        `json:"ads_txt,omitempty"`
	PubId            *string        `json:"pub_id,omitempty"`
	LinkedAccount    *string        `json:"linked_account,omitempty"`
	LinkedGAM        *string        `json:"linked_gam,omitempty"`
	IsDefault        *string        `json:"is_default,omitempty"`
	SellerJsonDomain *string        `json:"seller_json_domain,omitempty"`
	Params           *[]bidderParam `json:"params,omitempty"`
}

func (t *Bidder) getHistoryBidderBE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := bidderBE{}
	newData := bidderBE{}
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.Bidder,
		ObjectId:   t.getRootRecord().Id,
		ObjectType: t.Action(),
		DetailType: t.Detail.String(),
		App:        t.Detail.App(),
		UserId:     t.getRootRecord().UserId,
		BidderId:   t.getRootRecord().Id,
	}
	if t.getRootRecord().BidderTemplateId == 1 {
		history.ObjectName = t.getRootRecord().DisplayName
	} else {
		history.ObjectName = t.getRootRecord().BidderCode
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
		history.Title = "Add Bidder System"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update Bidder System"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete Bidder System"
		history.OldData = string(bOldData)
	}
	return
}

func (rec *bidderBE) MakeData(record mysql.TableBidder) {
	rec.Bidder = &record.BidderCode
	rec.AliasCode = &record.BidderAlias
	rec.AliasName = &record.AliasName
	rec.ShowOnPub = &record.ShowOnPub
	var bidAdjustment string
	if record.BidAdjustment != nil {
		bidAdjustment = fmt.Sprintf("%f", *record.BidAdjustment)
	} else {
		bidAdjustment = "0"
	}
	rec.BidAdjustment = &bidAdjustment
	rpm := fmt.Sprintf("%f", record.RPM)
	rec.RPM = &rpm
	isDefault := record.IsDefault.String()
	rec.IsDefault = &isDefault
	rec.SellerJsonDomain = &record.SellerDomain
	// Xử lý các phần riêng cho từng loại bidder
	if record.BidderTemplateId == 1 { // Xử lý khi bidder là google
		accountType := record.AccountType.String()
		rec.AccountType = &accountType
		rec.DisplayName = &record.DisplayName
		var mediaTypes []string
		for _, rlsBidderMediaType := range record.MediaTypes {
			var mediaType mysql.TableMediaType
			mysql.Client.Find(&mediaType, rlsBidderMediaType.MediaTypeId)
			mediaTypes = append(mediaTypes, mediaType.Name)
		}
		rec.MediaTypes = &mediaTypes
		rec.AdsTxt = &record.AdsTxt
		rec.PubId = &record.PubId
		linkedAccount := record.LinkedAccount.String()
		rec.LinkedAccount = &linkedAccount
		rec.LinkedGAM = &record.GAM.NetworkName
	} else if record.BidderTemplateId == 2 { // Bidder là amz
		var mediaTypes []string
		for _, rlsBidderMediaType := range record.MediaTypes {
			var mediaType mysql.TableMediaType
			mysql.Client.Find(&mediaType, rlsBidderMediaType.MediaTypeId)
			mediaTypes = append(mediaTypes, mediaType.Name)
		}
		rec.MediaTypes = &mediaTypes
		rec.AdsTxt = &record.AdsTxt
		rec.PubId = &record.PubId
		var params []bidderParam
		for _, param := range record.Params {
			params = append(params, bidderParam{
				Key:  param.Name,
				Type: param.Type,
			})
		}
		rec.Params = &params
	} else {
		var mediaTypes []string
		for _, rlsBidderMediaType := range record.MediaTypes {
			var mediaType mysql.TableMediaType
			mysql.Client.Find(&mediaType, rlsBidderMediaType.MediaTypeId)
			mediaTypes = append(mediaTypes, mediaType.Name)
		}
		rec.MediaTypes = &mediaTypes
		rec.AdsTxt = &record.AdsTxt
		var params []bidderParam
		for _, param := range record.Params {
			params = append(params, bidderParam{
				Key:  param.Name,
				Type: param.Type,
			})
		}
		rec.Params = &params
	}
}

func (t *Bidder) compareDataBidderBE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew bidderBE
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Bidder
	res, err := makeResponseCompare("Bidder", recordOld.Bidder, recordNew.Bidder, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Alias Code
	res, err = makeResponseCompare("Alias Code", recordOld.AliasCode, recordNew.AliasCode, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Alias Name
	res, err = makeResponseCompare("Alias Name", recordOld.AliasName, recordNew.AliasName, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Show On Pub
	res, err = makeResponseCompare("Show On Pub", recordOld.ShowOnPub, recordNew.ShowOnPub, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Account Type
	res, err = makeResponseCompare("Account Type", recordOld.AccountType, recordNew.AccountType, history.ObjectType)
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
	// RPM
	res, err = makeResponseCompare("RPM", recordOld.RPM, recordNew.RPM, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// AdsTxt
	res, err = makeResponseCompare("AdsTxt", recordOld.AdsTxt, recordNew.AdsTxt, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Pub Id
	res, err = makeResponseCompare("Pub Id", recordOld.PubId, recordNew.PubId, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Linked Account
	res, err = makeResponseCompare("Linked Account", recordOld.LinkedAccount, recordNew.LinkedAccount, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Linked GAM
	res, err = makeResponseCompare("Linked GAM", recordOld.LinkedGAM, recordNew.LinkedGAM, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Is Default
	res, err = makeResponseCompare("Is Default", recordOld.IsDefault, recordNew.IsDefault, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Seller.Json Domain
	res, err = makeResponseCompare("Seller.Json Domain", recordOld.SellerJsonDomain, recordNew.SellerJsonDomain, history.ObjectType)
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
