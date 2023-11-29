package history

import (
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
	"source/core/technology/mysql"
	"strings"
)

type LineItem struct {
	Detail    DetailLineItem
	CreatorId int64
	RecordOld mysql.TableLineItem
	RecordNew mysql.TableLineItem
}

func (t *LineItem) Page() string {
	return "Demand"
}

type DetailLineItem int

const (
	DetailLineItemFE DetailLineItem = iota + 1
	DetailLineItemBE
)

func (t DetailLineItem) String() string {
	switch t {
	case DetailLineItemFE:
		return "line_item_fe"
	case DetailLineItemBE:
		return "line_item_be"
	}
	return ""
}

func (t DetailLineItem) App() string {
	switch t {
	case DetailLineItemFE:
		return "FE"
	case DetailLineItemBE:
		return "BE"
	}
	return ""
}

func (t *LineItem) Type() TYPEHistory {
	return TYPEHistoryLineItem
}

func (t *LineItem) Action() mysql.TYPEObjectType {
	if t.RecordOld.Id == 0 && t.RecordNew.Id != 0 {
		return mysql.TYPEObjectTypeAdd
	} else if t.RecordOld.Id != 0 && t.RecordNew.Id == 0 {
		return mysql.TYPEObjectTypeDel
	}
	return mysql.TYPEObjectTypeUpdate
}

func (t *LineItem) Data() mysql.TableHistory {
	switch t.Detail {
	case DetailLineItemFE:
		return t.getHistoryLineItemFE()
	case DetailLineItemBE:
		return t.getHistoryLineItemBE()
	}
	return mysql.TableHistory{}
}

func (t *LineItem) CompareData(history mysql.TableHistory) []ResponseCompare {
	switch history.DetailType {
	case DetailLineItemFE.String():
		return t.compareDataLineItemFE(history)
	case DetailLineItemBE.String():
		return t.compareDataLineItemBE(history)
	}
	return []ResponseCompare{}
}

func (t *LineItem) getRootRecord() (record mysql.TableLineItem) {
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

type lineItemFE struct {
	Name           *string              `json:"name,omitempty"`
	Description    *string              `json:"description,omitempty"`
	Type           *string              `json:"type,omitempty"`
	SelectAccount  *string              `json:"select_account,omitempty"`
	LinkedGAM      *string              `json:"linked_gam,omitempty"`
	ConnectionType *string              `json:"connection_type,omitempty"`
	LineItemType   *string              `json:"line_item_type,omitempty"`
	Priority       *int                 `json:"priority,omitempty"`
	Status         *string              `json:"status,omitempty"`
	Target         target               `json:"target,omitempty"`
	Bidder         *[]bidderForLineItem `json:"bidder,omitempty"`
	AdsenseAdSlots *[]adsenseAdSlot     `json:"adsense_ad_slots,omitempty"`
}

type target struct {
	Domain    *[]string `json:"domain,omitempty"`
	Format    *[]string `json:"format,omitempty"`
	Size      *[]string `json:"size,omitempty"`
	AdTag     *[]string `json:"ad_tag,omitempty"`
	Geography *[]string `json:"geography,omitempty"`
	Device    *[]string `json:"device,omitempty"`
}

type bidderForLineItem struct {
	Name   string  `json:"name"`
	Server string  `json:"server"`
	Params []param `json:"params"`
}

type adsenseAdSlot struct {
	Size     string `json:"size"`
	AdSlotId string `json:"ad_slot_id"`
}

type param struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (t *LineItem) getHistoryLineItemFE() (history mysql.TableHistory) {
	// Xử lý record old + new
	oldData := lineItemFE{}
	newData := lineItemFE{}
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.LineItem,
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
		history.Title = "Add Line Item"
		history.NewData = string(bNewData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update Line Item"
		history.NewData = string(bNewData)
		history.OldData = string(bOldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete Line Item"
		history.OldData = string(bOldData)
	}
	return
}

func (rec *lineItemFE) MakeData(record mysql.TableLineItem) {
	if record.ServerType == mysql.TYPEServerTypePrebid {
		rec.Name = &record.Name
		rec.Description = &record.Description
		serverType := record.ServerType.String()
		rec.Type = &serverType
		rec.Priority = &record.Priority
		status := record.Status.String()
		rec.Status = &status

		// Xử lý target
		rec.Target = makeTarget(record.Targets)

		// Xử lý bidder
		var bidders []bidderForLineItem
		for _, bidderInfor := range record.BidderInfo {
			var params []param
			for _, p := range bidderInfor.Params {
				params = append(params, param{
					Name:  p.Name,
					Type:  p.Type,
					Value: p.Value,
				})
			}
			bidders = append(bidders, bidderForLineItem{
				Name:   bidderInfor.Name,
				Server: bidderInfor.BidderType.String(),
				Params: params,
			})
		}
		rec.Bidder = &bidders
	} else if record.ServerType == mysql.TYPEServerTypeGoogle {
		rec.Name = &record.Name
		rec.Description = &record.Description
		serverType := record.ServerType.String()
		rec.Type = &serverType
		rec.Priority = &record.Priority
		rec.SelectAccount = &record.BidderGoogle.DisplayName
		var gamNetwork mysql.TableGamNetwork
		mysql.Client.Model(&mysql.TableGamNetwork{}).Where("id = ?", record.LinkedGam).Find(&gamNetwork)
		rec.LinkedGAM = &gamNetwork.NetworkName
		connectionType := record.ConnectionType.String()
		rec.ConnectionType = &connectionType
		// Connection Type là Line Items thì mới có Line Item Type
		if record.ConnectionType == mysql.TYPEConnectionTypeLineItems {
			lineItemType := record.GamLineItemType.String()
			rec.LineItemType = &lineItemType
		}

		// Xử lý target
		rec.Target = makeTarget(record.Targets)
		// Loại các target không dùng trong type google
		rec.Target.Size = nil
		rec.Target.Geography = nil
		rec.Target.Device = nil

		// Xử lý adsenseSlots
		var adsenseAdSlots []adsenseAdSlot
		if record.BidderGoogle.AccountType == mysql.TYPEAccountTypeAdsense && record.ConnectionType == mysql.TYPEConnectionTypeLineItems && record.GamLineItemType == mysql.TYPEGamLineItemTypeDisplay {
			for _, adsenseSlot := range record.AdsenseAdSlots {
				adsenseAdSlots = append(adsenseAdSlots, adsenseAdSlot{
					Size:     adsenseSlot.Size,
					AdSlotId: adsenseSlot.AdsenseAdSlotId,
				})
			}
		}
		if len(adsenseAdSlots) > 0 {
			rec.AdsenseAdSlots = &adsenseAdSlots
		}

	}
}

func makeTarget(targets []mysql.TableTarget) (target target) {
	var domain, format, size, adTag, geo, device []string
	for _, target := range targets {
		if target.InventoryId != 0 {
			if target.InventoryId == -1 {
				domain = append(domain, "all")
			} else {
				var inventoryName string
				mysql.Client.Model(&mysql.TableInventory{}).Select("name").Where("id = ?", target.InventoryId).Find(&inventoryName)
				if !govalidator.IsNull(inventoryName) {
					domain = append(domain, inventoryName)
				}
			}
		} else if target.AdFormatId != 0 {
			if target.AdFormatId == -1 {
				format = append(format, "all")
			} else {
				var typeName string
				mysql.Client.Model(&mysql.TableAdType{}).Select("name").Where("id = ?", target.AdFormatId).Find(&typeName)
				if !govalidator.IsNull(typeName) {
					format = append(format, typeName)
				}
			}
		} else if target.AdSizeId != 0 {
			if target.AdSizeId == -1 {
				size = append(size, "all")
			} else {
				var sizeName string
				mysql.Client.Model(&mysql.TableAdSize{}).Select("name").Where("id = ?", target.AdSizeId).Find(&sizeName)
				if !govalidator.IsNull(sizeName) {
					size = append(size, sizeName)
				}
			}
		} else if target.TagId != 0 {
			if target.TagId == -1 {
				adTag = append(adTag, "all")
			} else {
				var tagName string
				mysql.Client.Model(&mysql.TableInventoryAdTag{}).Select("name").Where("id = ?", target.TagId).Find(&tagName)
				if !govalidator.IsNull(tagName) {
					adTag = append(adTag, tagName)
				}
			}
		} else if target.GeoId != 0 {
			if target.GeoId == -1 {
				geo = append(geo, "all")
			} else {
				var geoName string
				mysql.Client.Model(&mysql.TableCountry{}).Select("name").Where("id = ?", target.GeoId).Find(&geoName)
				if !govalidator.IsNull(geoName) {
					geo = append(geo, geoName)
				}
			}
		} else if target.DeviceId != 0 {
			if target.DeviceId == -1 {
				device = append(device, "all")
			} else {
				var deviceName string
				mysql.Client.Model(&mysql.TableDevice{}).Select("name").Where("id = ?", target.DeviceId).Find(&deviceName)
				if !govalidator.IsNull(deviceName) {
					device = append(device, deviceName)
				}
			}
		}
	}

	if len(domain) > 0 {
		target.Domain = &domain
	}
	if len(format) > 0 {
		target.Format = &format
	}
	if len(size) > 0 {
		target.Size = &size
	}
	if len(adTag) > 0 {
		target.AdTag = &adTag
	}
	if len(geo) > 0 {
		target.Geography = &geo
	}
	if len(device) > 0 {
		target.Device = &device
	}

	return
}

func (t *LineItem) compareDataLineItemFE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew lineItemFE
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Name
	res, err := makeResponseCompare("Name", recordOld.Name, recordNew.Name, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Description
	res, err = makeResponseCompare("Description", recordOld.Description, recordNew.Description, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Type
	res, err = makeResponseCompare("Type", recordOld.Type, recordNew.Type, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Select Account
	res, err = makeResponseCompare("Select Account", recordOld.SelectAccount, recordNew.SelectAccount, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Linked GAM
	res, err = makeResponseCompare("Linked GAM", recordOld.LinkedGAM, recordNew.LinkedGAM, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Connection Type
	res, err = makeResponseCompare("Connection Type", recordOld.ConnectionType, recordNew.ConnectionType, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// LineItem Type
	res, err = makeResponseCompare("LineItem Type", recordOld.LineItemType, recordNew.LineItemType, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Priority
	res, err = makeResponseCompare("Priority", pointerIntToString(recordOld.Priority), pointerIntToString(recordNew.Priority), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Status
	res, err = makeResponseCompare("Status", recordOld.Status, recordNew.Status, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Target Domain
	res, err = makeResponseCompare("Target > Domain", pointerArrayStringToString(recordOld.Target.Domain), pointerArrayStringToString(recordNew.Target.Domain), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Target Format
	res, err = makeResponseCompare("Target > Format", pointerArrayStringToString(recordOld.Target.Format), pointerArrayStringToString(recordNew.Target.Format), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Target Size
	res, err = makeResponseCompare("Target > Size", pointerArrayStringToString(recordOld.Target.Size), pointerArrayStringToString(recordNew.Target.Size), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Target Ad Tag
	res, err = makeResponseCompare("Target > Ad Tag", pointerArrayStringToString(recordOld.Target.AdTag), pointerArrayStringToString(recordNew.Target.AdTag), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Target Geography
	res, err = makeResponseCompare("Target > Geography", pointerArrayStringToString(recordOld.Target.Geography), pointerArrayStringToString(recordNew.Target.Geography), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Target Device
	res, err = makeResponseCompare("Target > Device", pointerArrayStringToString(recordOld.Target.Device), pointerArrayStringToString(recordNew.Target.Device), history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// fmt.Printf("%+v \n", history)
	// fmt.Printf("%+v \n", recordOld.Bidder)
	// fmt.Printf("%+v \n", recordNew.Bidder)
	// Bidder
	respBidder, err := t.makeResponseFromBidderInfo(recordOld.Bidder, recordNew.Bidder)
	if err == nil {
		responses = append(responses, respBidder...)
	}
	// Adsense Ad Slot
	respAdsenseAdSlot, err := t.makeResponseFromAdsenseAdSlot(recordOld.AdsenseAdSlots, recordNew.AdsenseAdSlots)
	if err == nil {
		responses = append(responses, respAdsenseAdSlot...)
	}
	return
}

func (t *LineItem) makeResponseFromBidderInfo(bidderOld *[]bidderForLineItem, bidderNew *[]bidderForLineItem) (responses []ResponseCompare, err error) {
	if bidderOld == nil && bidderNew == nil {
		err = errors.New("no response")
		return
	}
	mapBidderOld := make(map[string]bidderForLineItem)
	mapBidderNew := make(map[string]bidderForLineItem)

	if bidderOld != nil {
		for _, bidder := range *bidderOld {
			name := strings.Title(bidder.Name)
			mapBidderOld[name] = bidder
		}
	}
	if bidderNew != nil {
		for _, bidder := range *bidderNew {
			name := strings.Title(bidder.Name)
			mapBidderNew[name] = bidder
		}
	}

	for name, bidderInfoOld := range mapBidderOld {
		name = strings.Title(name)
		if bidderInfoNew, ok := mapBidderNew[name]; ok {
			// Server Type
			res, err := makeResponseCompare("Bidder > "+name+" > Server Type", &bidderInfoOld.Server, &bidderInfoNew.Server, mysql.TYPEObjectTypeUpdate)
			if err == nil {
				responses = append(responses, res)
			}

			// Param
			mapParamBidderOld := make(map[string]param)
			mapParamBidderNew := make(map[string]param)

			for _, param := range bidderInfoOld.Params {
				nameParam := strings.Title(param.Name)
				mapParamBidderOld[nameParam] = param
			}
			for _, param := range bidderInfoNew.Params {
				nameParam := strings.Title(param.Name)
				mapParamBidderNew[nameParam] = param
			}

			for nameParam, paramOld := range mapParamBidderOld {
				nameParam = strings.Title(nameParam)
				// Kiểm tra nếu có tồn tại trong paramNew là update
				if paramNew, ok := mapParamBidderNew[nameParam]; ok {
					// Type
					res, err := makeResponseCompare("Bidder > "+name+" > Param > "+nameParam+" > Type", &paramOld.Type, &paramNew.Type, mysql.TYPEObjectTypeUpdate)
					if err == nil {
						responses = append(responses, res)
					}
					// Value
					res, err = makeResponseCompare("Bidder > "+name+" > Param > "+nameParam+" > Value", &paramOld.Value, &paramNew.Value, mysql.TYPEObjectTypeUpdate)
					if err == nil {
						responses = append(responses, res)
					}
					// Loại các param update để lại các param add mới trong mapParamBidderNew
					delete(mapParamBidderNew, nameParam)
				} else {
					// Nếu như param không tồn tại trong paramNew tức param này bị xóa
					// Type
					res, err := makeResponseCompare("Bidder > "+name+" > Param > "+nameParam+" > Type", &paramOld.Type, nil, mysql.TYPEObjectTypeDel)
					if err == nil {
						responses = append(responses, res)
					}
					// Value
					res, err = makeResponseCompare("Bidder > "+name+" > Param > "+nameParam+" > Value", &paramOld.Value, nil, mysql.TYPEObjectTypeDel)
					if err == nil {
						responses = append(responses, res)
					}
				}
			}
			// Sau khi đã loại các param đã tồn tại trong param old còn lại sẽ là các param mới
			for nameParam, paramNew := range mapParamBidderNew {
				nameParam = strings.Title(nameParam)
				// Type
				res, err := makeResponseCompare("Bidder > "+name+" > Param > "+nameParam+" > Type", nil, &paramNew.Value, mysql.TYPEObjectTypeAdd)
				if err == nil {
					responses = append(responses, res)
				}
				// Value
				res, err = makeResponseCompare("Bidder > "+name+" > Param > "+nameParam+" > Value", nil, &paramNew.Value, mysql.TYPEObjectTypeAdd)
				if err == nil {
					responses = append(responses, res)
				}
			}

			// Xóa các bidder update
			delete(mapBidderNew, name)
		} else {
			// Nếu chỉ tồn tại bidder old tức bidder này đã bị xóa toàn bộ các newData là nil
			// Server Type
			res, err := makeResponseCompare("Bidder > "+name+" > Server Type", &bidderInfoOld.Server, nil, mysql.TYPEObjectTypeDel)
			if err == nil {
				responses = append(responses, res)
			}
			// Param
			for _, param := range bidderInfoOld.Params {
				// Type
				res, err := makeResponseCompare("Bidder > "+name+" > Param > "+param.Name+" > Type", &param.Type, nil, mysql.TYPEObjectTypeDel)
				if err == nil {
					responses = append(responses, res)
				}
				// Value
				res, err = makeResponseCompare("Bidder > "+name+" > Param > "+param.Name+" > Value", &param.Type, nil, mysql.TYPEObjectTypeDel)
				if err == nil {
					responses = append(responses, res)
				}
			}
		}
	}
	// Các bidderNew còn lại đều là các bidder add mới
	for name, bidderInfoNew := range mapBidderNew {
		name = strings.Title(name)
		// Server Type
		res, err := makeResponseCompare("Bidder > "+name+" > Server Type", nil, &bidderInfoNew.Server, mysql.TYPEObjectTypeAdd)
		if err == nil {
			responses = append(responses, res)
		}
		// Param
		for _, param := range bidderInfoNew.Params {
			param.Name = strings.Title(param.Name)
			// Type
			res, err := makeResponseCompare("Bidder > "+name+" > Param > "+param.Name+" > Type", nil, &param.Type, mysql.TYPEObjectTypeAdd)
			if err == nil {
				responses = append(responses, res)
			}
			// Value
			res, err = makeResponseCompare("Bidder > "+name+" > Param > "+param.Name+" > Value", nil, &param.Type, mysql.TYPEObjectTypeAdd)
			if err == nil {
				responses = append(responses, res)
			}
		}
	}
	return
}

func (t *LineItem) makeResponseFromAdsenseAdSlot(adsenseAdSlotOld *[]adsenseAdSlot, adsenseAdSlotNew *[]adsenseAdSlot) (responses []ResponseCompare, err error) {
	if adsenseAdSlotOld == nil && adsenseAdSlotNew == nil {
		err = errors.New("no response")
		return
	}
	mapAdsenseAdSlotOld := make(map[string]adsenseAdSlot)
	mapAdsenseAdSlotNew := make(map[string]adsenseAdSlot)

	if adsenseAdSlotOld != nil {
		for _, adsenseAdSlot := range *adsenseAdSlotOld {
			mapAdsenseAdSlotOld[adsenseAdSlot.Size] = adsenseAdSlot
		}
	}
	if adsenseAdSlotNew != nil {
		for _, adsenseAdSlot := range *adsenseAdSlotNew {
			mapAdsenseAdSlotNew[adsenseAdSlot.Size] = adsenseAdSlot
		}
	}

	for size, adSlotOld := range mapAdsenseAdSlotOld {
		if adSlotNew, ok := mapAdsenseAdSlotNew[size]; ok {
			// Ad Slot Id
			res, err := makeResponseCompare("Adsense Ad Slot > Size "+size+" > Ad Slot Id", &adSlotOld.AdSlotId, &adSlotNew.AdSlotId, mysql.TYPEObjectTypeUpdate)
			if err == nil {
				responses = append(responses, res)
			}

			// Xóa các adslot đã update
			delete(mapAdsenseAdSlotNew, size)
		} else {
			// Nếu chỉ tồn tại adslot old tức adslot này đã bị xóa toàn bộ các newData là nil
			// Ad Slot Id
			res, err := makeResponseCompare("Adsense Ad Slot > Size "+size+" > Ad Slot Id", &adSlotOld.AdSlotId, nil, mysql.TYPEObjectTypeDel)
			if err == nil {
				responses = append(responses, res)
			}
		}
	}
	// Các adslotNew còn lại đều là các adslot add mới
	for size, adSlotNew := range mapAdsenseAdSlotNew {
		// Ad Slot Id
		res, err := makeResponseCompare("Adsense Ad Slot > Size "+size+" > Ad Slot Id", nil, &adSlotNew.AdSlotId, mysql.TYPEObjectTypeAdd)
		if err == nil {
			responses = append(responses, res)
		}
	}
	return
}

func (t *LineItem) getHistoryLineItemBE() (history mysql.TableHistory) {
	history = t.getHistoryLineItemFE()
	return
}

func (t *LineItem) compareDataLineItemBE(history mysql.TableHistory) (responses []ResponseCompare) {
	responses = t.compareDataLineItemFE(history)
	return
}
