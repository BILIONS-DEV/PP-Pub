package model

import (
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"source/apps/frontend/lang"
	"source/apps/frontend/payload"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"source/pkg/datatable"
	"source/pkg/htmlblock"
	"source/pkg/pagination"
	"source/pkg/utility"
	"strconv"
	"strings"
	"time"
)

type LineItemV2 struct{}

type LineItemRecordV2 struct {
	mysql.TableLineItemV2
}

func (LineItemRecordV2) TableName() string {
	return mysql.Tables.LineItemV2
}

type ListOfLineItemV2 struct {
	Id   int64
	List []string
}

type SearchLineItemV2 struct {
	Id       int64  `json:"id"`
	Selected bool   `json:"selected"`
	Name     string `json:"name"`
}

func (LineItemV2) GetById(id, userId int64) (row LineItemRecordV2, err error) {
	err = mysql.Client.
		//Debug().
		Where("id = ? and user_id = ?", id, userId).Find(&row).Error
	if row.Id == 0 {
		err = errors.New("Record not found")
		return
	}
	// Get các rls của line item
	row.TableLineItemV2.GetRls()
	return
}

func (LineItemV2) GetByUser(userId int64) (records []LineItemRecordV2) {
	mysql.Client.Where("status = 1 and user_id = ?", userId).Find(&records)
	return
}

func (rec *LineItemRecordV2) makeRecord(payload payload.LineItemAddV2, userId int64) (err error) {
	rec.Id = payload.Id
	rec.UserId = userId
	rec.Name = payload.Name
	rec.Description = payload.Description
	rec.ServerType = payload.ServerType
	rec.LinkedGam = payload.LinkedGam
	rec.ConnectionType = payload.ConnectionType
	rec.GamLineItemType = payload.GamLineItemType
	rec.Rate = payload.BidderRate
	rec.VastUrl = payload.BidderVastUrl
	rec.AdTag = payload.BidderAdTag
	rec.Priority = payload.Priority
	rec.LineItemType = 1
	if payload.Id == 0 {
		rec.IsLock = mysql.TYPEIsLockTypeUnlock
		rec.AutoCreate = mysql.Off
	}

	layoutISO := "01/02/2006"
	if payload.StartDate != "" {
		startDate, err := time.Parse(layoutISO, payload.StartDate)
		if err != nil {
			return err
		}
		err = rec.StartDate.Scan(startDate)
		if err != nil {
			return err
		}
	}
	if payload.EndDate != "" {
		endDate, err := time.Parse(layoutISO, payload.EndDate)
		if err != nil {
			return err
		}
		err = rec.EndDate.Scan(endDate)
		if err != nil {
			return err
		}
	}

	//rec.StartDate.Time = startDate
	//rec.EndDate.Time = endDate
	if payload.Status == "on" {
		rec.Status = mysql.On
	} else {
		rec.Status = mysql.Off
	}
	return
}

func (t *LineItemV2) GetLineSystemApacdex(inventoryId int64) (record LineItemRecordV2, err error) {
	err = mysql.Client.Where("apd_inventory = ? && user_id = 0", inventoryId).Find(&record).Error
	return
}

func (t *LineItemV2) AutoCreateLineSystemApacdex(inventory InventoryRecord, bidder BidderRecord) (errs []ajax.Error) {
	// Check và lấy các inventory đã được system approve và pub active
	status := new(InventoryConnectionDemand).GetStatus(inventory.Id, bidder.Id)
	if status != mysql.TYPEStatusConnectionDemandLive {
		// Nếu status chưa đổi live thì bỏ qua
		return
	}

	layoutTime := "3:04:05 PM, January 2, 2006"
	timeConfigTime := time.Now().Format(layoutTime)
	// Tạo line item system cho apacdex
	lineItem := LineItemRecordV2{mysql.TableLineItemV2{
		UserId:       0,
		ApdInventory: inventory.Id,
		Name:         "[Default] - Target demand Apacdex - " + inventory.Name + "[" + strconv.FormatInt(inventory.Id, 10) + "]",
		Description:  "Auto create at " + timeConfigTime,
		ServerType:   1,
		Status:       1,
		IsLock:       mysql.TYPEIsLockTypeLock,
		AutoCreate:   mysql.On,
		Priority:     1,
		LineItemType: 2,
	}}
	// Insert to database
	err := mysql.Client.Create(&lineItem).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Translate.Errors.LineItemError.Add.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}
	// Tạo bidder info client của line item
	bidderInfo := new(LineItemBidderInfoV2).CreateBidderInfo(LineItemBidderInfoRecordV2{mysql.TableLineItemBidderInfoV2{
		LineItemId: lineItem.Id,
		BidderId:   bidder.Id,
		Name:       bidder.BidderCode,
		BidderType: 1,
	}})
	// Tạo rls bidder param cho bidder info
	err = mysql.Client.Create(&LineItemBidderParamsRecordV2{mysql.TableLineItemBidderParamsV2{
		LineItemId:       lineItem.Id,
		BidderId:         bidder.Id,
		LineItemBidderId: bidderInfo.Id,
		Name:             "siteId",
		Type:             "string",
		Value:            strconv.FormatInt(inventory.ApacSiteId, 10),
	}}).Error

	// Tạo bidder info s2s của line item
	bidderInfo = new(LineItemBidderInfoV2).CreateBidderInfo(LineItemBidderInfoRecordV2{mysql.TableLineItemBidderInfoV2{
		LineItemId: lineItem.Id,
		BidderId:   bidder.Id,
		Name:       bidder.BidderCode,
		BidderType: 2,
	}})
	// Tạo rls bidder param cho bidder info
	err = mysql.Client.Create(&LineItemBidderParamsRecordV2{mysql.TableLineItemBidderParamsV2{
		LineItemId:       lineItem.Id,
		BidderId:         bidder.Id,
		LineItemBidderId: bidderInfo.Id,
		Name:             "siteId",
		Type:             "string",
		Value:            strconv.FormatInt(inventory.ApacSiteId, 10),
	}}).Error

	// Tạo target domain cho nhưng domain có bidder apacdex live
	_ = new(TargetV2).DeleteTargetInventory(TargetRecordV2{mysql.TableTargetV2{
		LineItemId: lineItem.Id,
	}})
	// Tạo target cho inventory
	mysql.Client.Create(&TargetRecordV2{mysql.TableTargetV2{
		UserId:      0,
		LineItemId:  lineItem.Id,
		InventoryId: inventory.Id,
	}})
	// Tạo target all cho tất cả target
	mysql.Client.Create(&TargetRecordV2{mysql.TableTargetV2{
		UserId:     0,
		LineItemId: lineItem.Id,
		TagId:      -1,
	}})
	mysql.Client.Create(&TargetRecordV2{mysql.TableTargetV2{
		UserId:     0,
		LineItemId: lineItem.Id,
		AdSizeId:   -1,
	}})
	mysql.Client.Create(&TargetRecordV2{mysql.TableTargetV2{
		UserId:     0,
		LineItemId: lineItem.Id,
		AdFormatId: -1,
	}})
	mysql.Client.Create(&TargetRecordV2{mysql.TableTargetV2{
		UserId:     0,
		LineItemId: lineItem.Id,
		GeoId:      -1,
	}})
	mysql.Client.Create(&TargetRecordV2{mysql.TableTargetV2{
		UserId:     0,
		LineItemId: lineItem.Id,
		DeviceId:   -1,
	}})
	return
}

func (t *LineItemV2) UpdateLineSystemApacdex(lineItem LineItemRecordV2, bidderId int64) (errs []ajax.Error) {

	// Get status connection inventory
	status := new(InventoryConnectionDemand).GetStatus(lineItem.ApdInventory, bidderId)
	// Trường hợp live
	if status == mysql.TYPEStatusConnectionDemandLive {
		// Nếu line item đang off thì tiến hành on lại
		if lineItem.Status == mysql.Off {
			mysql.Client.Model(LineItemRecordV2{}).Where("id = ?", lineItem.Id).Update("status", mysql.On)
		}

		// Xóa toàn bộ target inventory cũ
		_ = new(TargetV2).DeleteTargetInventory(TargetRecordV2{mysql.TableTargetV2{
			LineItemId: lineItem.Id,
		}})
		// Tạo target cho inventory
		mysql.Client.Create(&TargetRecordV2{mysql.TableTargetV2{
			UserId:      0,
			LineItemId:  lineItem.Id,
			InventoryId: lineItem.ApdInventory,
		}})
	} else {
		// Trường hợp connection off update line item về off
		mysql.Client.Model(LineItemRecordV2{}).Where("id = ?", lineItem.Id).Update("status", mysql.Off)
	}
	return
}

func (t *LineItemV2) AddLineItem(inputs payload.LineItemAddV2, user UserRecord, userAdmin UserRecord) (record LineItemRecordV2, errs []ajax.Error) {
	// Validate inputs
	errs = t.Validate(inputs, user)
	if len(errs) > 0 {
		return
	}
	var err error
	// Make record use insert to database
	err = record.makeRecord(inputs, user.Id)
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Translate.Errors.LineItemError.Add.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}

	// Insert to database
	err = mysql.Client.Create(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Translate.Errors.LineItemError.Add.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}

	switch inputs.ServerType {
	case mysql.TYPEServerTypeGoogle: //=> CASE: LineItemType là GOOGLE
		//Tạo line_item account nếu server là google
		for _, v := range inputs.SelectAccount {
			if v != 0 {
				//bidderId, _ := strconv.ParseInt(v, 10, 64)
				bidderId := v
				err = new(LineItemAccount).Create(record.Id, bidderId)
				if err != nil {
					if !utility.IsWindow() {
						errs = append(errs, ajax.Error{
							Id:      "",
							Message: lang.Translate.Errors.LineItemError.LineItemAccount.ToString(),
						})
					} else {
						errs = append(errs, ajax.Error{
							Id:      "",
							Message: err.Error(),
						})
					}
					return
				}
			}
		}
		// Tạo Adsense Ad Slot
		errs = append(errs, t.CreateAdsenseAdSlot(record.Id, inputs)...)

	case mysql.TYPEServerTypePrebid: //=>> CASE: LineItemType là PREBID
		//Tạo bidder info + params cho line item trong bảng line_item_bidder_info nếu server là prebid
		err = t.CreateBidderInfo(record.Id, inputs)
		if err != nil {
			if !utility.IsWindow() {
				errs = append(errs, ajax.Error{
					Id:      "",
					Message: lang.Translate.Errors.LineItemError.BidderInfo.ToString(),
				})
			} else {
				errs = append(errs, ajax.Error{
					Id:      "",
					Message: err.Error(),
				})
			}
			return
		}
	}

	//Tạo target cho line item vào bảng rl
	err = t.CreateTarget(record.Id, user.Id, inputs)
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Translate.Errors.LineItemError.Target.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}
	new(Inventory).UpdateRenderCacheWithLineItem(record.Id, user.Id)

	//Nếu server type là google thì set line vào hàng đợi worker push_line_item_dfp
	if record.ServerType == mysql.TYPEServerTypeGoogle {
		new(LineItemV2).PushToQueueWorkerLineItemDfp(record.Id)
	}

	//recordNew, _ := new(LineItemV2).GetById(record.Id, record.UserId)

	// Push History
	//var creatorId int64
	//if userAdmin.Id != 0 {
	//	creatorId = userAdmin.Id
	//} else {
	//	creatorId = user.Id
	//}
	//_ = history.PushHistory(&history.LineItem{
	//	Detail:    history.DetailLineItemFE,
	//	CreatorId: creatorId,
	//	RecordOld: mysql.TableLineItem{},
	//	RecordNew: recordNew.TableLineItem,
	//})
	return
}

func (t *LineItemV2) UpdateLineItem(inputs payload.LineItemAddV2, user UserRecord, userAdmin UserRecord) (recordNew LineItemRecordV2, errs []ajax.Error) {
	lang := lang.Translate
	recordOld := t.GetOfUserById(inputs.Id, user.Id)
	if recordOld.Id == 0 {
		errs = append(errs, ajax.Error{
			Id:      "id",
			Message: "You don't own this line item",
		})
		return
	} else if recordOld.AutoCreate == mysql.On { // Đây là line auto create
		if recordOld.IsLock == mysql.TYPEIsLockTypeLock {
			errs = append(errs, ajax.Error{
				Id:      "id",
				Message: "You can't edit this line item",
			})
			return
		}
		var status mysql.TYPEOnOff
		if inputs.Status == "on" {
			status = mysql.On
		} else {
			status = mysql.Off
		}
		mysql.Client.Model(&LineItemRecordV2{}).Where("id = ?", inputs.Id).Update("status", status)
		new(Inventory).ResetCacheAll(user.Id)
		return
	}
	// Gán lại các giá trị không thay đổi trong edit
	recLineItemAccount, _ := new(LineItemAccount).GetByLineItem(recordOld.Id)
	inputs.ServerType = recordOld.ServerType
	inputs.SelectAccount[0] = new(Bidder).GetById(recLineItemAccount.BidderId, user.Id).Id
	if recordOld.GamLineItemType != 0 {
		inputs.GamLineItemType = recordOld.GamLineItemType
	}
	inputs.LinkedGam = recordOld.LinkedGam

	// Validate inputs
	errs = t.Validate(inputs, user)
	if len(errs) > 0 {
		return
	}
	// Insert to database
	var err error
	recordNew = recordOld
	err = recordNew.makeRecord(inputs, user.Id)
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.LineItemError.Edit.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}
	//fmt.Printf("%+v \n", record)
	err = mysql.Client.Save(&recordNew).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.LineItemError.Edit.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}

	//Xóa toàn bộ Adsense Ad Slot cũ
	new(LineItemAdsenseAdSlot).DeleteByLineItem(recordNew.Id)

	// Tạo Adsense Ad Slot
	errs = append(errs, t.CreateAdsenseAdSlot(recordNew.Id, inputs)...)

	////Xóa toàn bộ account theo line item cũ
	//new(LineItemAccount).DeleteAccount(record.Id)
	//
	////Tạo line_item account nếu server là google
	//if inputs.ServerType == mysql.TYPEServerTypeGoogle {
	//	for _, v := range inputs.SelectAccount {
	//		if v != "" {
	//			bidderId, _ := strconv.ParseInt(v, 10, 64)
	//			err = new(LineItemAccount).Create(record.Id, bidderId)
	//			if err != nil {
	//				if !utility.IsWindow() {
	//					errs = append(errs, ajax.Error{
	//						Id:      "",
	//						Message: lang.Errors.LineItemError.LineItemAccount.ToString(),
	//					})
	//				} else {
	//					errs = append(errs, ajax.Error{
	//						Id:      "",
	//						Message: err.Error(),
	//					})
	//				}
	//				return
	//			}
	//		}
	//	}
	//}

	//Xóa toàn bộ bidder info + param cũ của line item
	new(LineItemBidderInfoV2).DeleteBidderInfo(recordNew.Id)
	new(LineItemBidderParamsV2).DeleteByLineItemId(recordNew.Id)

	//Tạo bidder info cho line item trong bảng line_item_bidder_info nếu server là prebid
	if inputs.ServerType == mysql.TYPEServerTypePrebid {
		err = t.CreateBidderInfo(recordNew.Id, inputs)
		if err != nil {
			if !utility.IsWindow() {
				errs = append(errs, ajax.Error{
					Id:      "",
					Message: lang.Errors.LineItemError.BidderInfo.ToString(),
				})
			} else {
				errs = append(errs, ajax.Error{
					Id:      "",
					Message: err.Error(),
				})
			}
			return
		}
	}

	//Reset cache lại các domain với target cũ
	new(Inventory).UpdateRenderCacheWithLineItem(recordNew.Id, user.Id)

	//Xóa toàn bộ target cũ để tạo mới list target nhận đc
	err = new(TargetV2).DeleteTarget(TargetRecordV2{mysql.TableTargetV2{
		LineItemId: recordNew.Id,
	}})
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.LineItemError.Target.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}
	//Tạo target cho line item vào bảng rl
	err = t.CreateTarget(recordNew.Id, user.Id, inputs)
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.LineItemError.Target.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}
	//Reset cache lại các domain với target mới
	new(Inventory).UpdateRenderCacheWithLineItem(recordNew.Id, user.Id)

	//Nếu server type là google thì set line vào hàng đợi worker push_line_item_dfp
	if recordNew.ServerType == 2 {
		new(LineItemV2).PushToQueueWorkerLineItemDfp(recordNew.Id)
	}
	// Get all lại data recordNew
	//recordNew, _ = new(LineItemV2).GetById(recordNew.Id, recordNew.UserId)

	// Push History
	//var creatorId int64
	//if userAdmin.Id != 0 {
	//	creatorId = userAdmin.Id
	//} else {
	//	creatorId = user.Id
	//}
	//_ = history.PushHistory(&history.LineItem{
	//	Detail:    history.DetailLineItemFE,
	//	CreatorId: creatorId,
	//	RecordOld: recordOld.TableLineItem,
	//	RecordNew: recordNew.TableLineItem,
	//})
	return
}

func (t *LineItemV2) Validate(inputs payload.LineItemAddV2, user UserRecord) (errs []ajax.Error) {
	if utility.ValidateString(inputs.Name) == "" {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: lang.Translate.ErrorRequired.ToString(),
		})
	}
	if inputs.ServerType == mysql.TYPEServerTypePrebid {
		if len(inputs.BidderParams) == 0 {
			errs = append(errs, ajax.Error{
				Id:      "select-bidder",
				Message: "Choose at least one bidder",
			})
		} else {
			err := t.ValidateBidder(inputs, user)
			errs = append(errs, err...)
		}
	} else if inputs.ServerType == mysql.TYPEServerTypeGoogle {
		if len(inputs.SelectAccount) == 0 {
			errs = append(errs, ajax.Error{
				Id:      "select_account",
				Message: lang.Translate.ErrorRequired.ToString(),
			})
		}
		if inputs.ConnectionType == 0 {
			errs = append(errs, ajax.Error{
				Id:      "select_connection_type",
				Message: lang.Translate.ErrorRequired.ToString(),
			})
		} else {
			if inputs.ConnectionType == mysql.TYPEConnectionTypeLineItems {
				if inputs.GamLineItemType == 0 {
					errs = append(errs, ajax.Error{
						Id:      "select_line_item_type",
						Message: lang.Translate.ErrorRequired.ToString(),
					})
				}
			}
		}
		if len(inputs.SelectAccount) > 0 && inputs.ConnectionType != 0 && inputs.GamLineItemType != 0 {
			bidderGG := new(Bidder).GetById(inputs.SelectAccount[0], user.Id)
			if bidderGG.AccountType == mysql.TYPEAccountTypeAdsense {
				if inputs.ConnectionType == mysql.TYPEConnectionTypeLineItems && inputs.GamLineItemType == mysql.TYPEGamLineItemTypeDisplay {
					if len(inputs.AdsenseAdSlots) == 0 {
						errs = append(errs, ajax.Error{
							Id:      "select-adsense-ad-slot-size",
							Message: lang.Translate.ErrorRequired.ToString(),
						})
					} else {
						errs = append(errs, t.ValidateAdsenseAdSlot(inputs)...)
					}
				}
			}
		}

		// Validate nếu bidder là pp_adx
		if len(inputs.SelectAccount) > 0 {
			bidderGG := new(Bidder).GetByIdNoCheckUser(inputs.SelectAccount[0])
			if bidderGG.Id > 0 && bidderGG.UserId == 0 && bidderGG.AccountType.IsAdx() {
				if inputs.ConnectionType != mysql.TYPEConnectionTypeMCM {
					errs = append(errs, ajax.Error{
						Id:      "select_connection_type",
						Message: "connection type not accept",
					})
				} else {
					if inputs.LinkedGam > 0 {
						rlsConnectionMCM := new(RlsConnectionMCM).GetStatus(bidderGG.Id, inputs.LinkedGam, user.Id)
						if rlsConnectionMCM.Status != mysql.TYPEConnectionMCMTypeAccept {
							errs = append(errs, ajax.Error{
								Id:      "linked_gam",
								Message: "GAM not accept",
							})
						}
					}
				}
			}
		}

		if inputs.LinkedGam == 0 {
			errs = append(errs, ajax.Error{
				Id:      "linked_gam",
				Message: lang.Translate.ErrorRequired.ToString(),
			})
		}
	}

	if len(inputs.ListAdInventory) == 0 {
		errs = append(errs, ajax.Error{
			Id:      "text_for_domain",
			Message: "Target domains is required",
		})
	}
	return
}

func (t *LineItemV2) ValidateAdsenseAdSlot(inputs payload.LineItemAddV2) (errs []ajax.Error) {
	for _, adSlot := range inputs.AdsenseAdSlots {
		if adSlot.AdSlotId == "" {
			errs = append(errs, ajax.Error{
				Id:      "adsense-ad-slot-" + adSlot.Size,
				Message: "Ad Slot ID is required",
			})
		}
	}
	return
}

func (t *LineItemV2) ValidateBidder(inputs payload.LineItemAddV2, user UserRecord) (errs []ajax.Error) {
	var listIdTypeClient []string
	for _, bidderInfo := range inputs.BidderParams {
		if bidderInfo.BidderType != mysql.TYPEBidderTypePrebidClient && bidderInfo.BidderType != mysql.TYPEBidderTypePrebidServer {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: " Bidder " + bidderInfo.BidderName + " type invalid",
			})
			return
		}
		keyCheckSelection := strconv.FormatInt(bidderInfo.BidderId, 10) + "_" + string(bidderInfo.ConfigType)
		if bidderInfo.BidderType == 1 {
			listIdTypeClient = append(listIdTypeClient, keyCheckSelection)
		}
		flag := t.Count(keyCheckSelection, listIdTypeClient)
		if !flag {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: " Bidder " + bidderInfo.BidderName + " doesn't allow selection of more than 2 bidders with type is Prebid Client",
			})
			return
		}
		// Xử lý các param đã tồn tại trong bidder
		bidderParams := new(BidderParams).GetByBidderId(bidderInfo.BidderId)
		for _, bidderParam := range bidderParams {
			idParam := strconv.FormatInt(bidderInfo.BidderId, 10) + "-" + bidderParam.Name + "-" + strconv.Itoa(bidderInfo.BidderIndex)
			param, exists := bidderInfo.Params[bidderParam.Name]
			if !exists {
				continue
			}
			//required := t.CheckParamRequired(bidderInfo.BidderName, bidderParam.Name)
			required := new(PbBidder).IsParamRequiredByBidder(bidderInfo.BidderName, bidderParam.Name)
			if required && govalidator.IsNull(param.Value) {
				errs = append(errs, ajax.Error{
					Id:      idParam,
					Message: "(*) Required",
				})
			}
			errs = append(errs, t.validateValueOfType(bidderParam.Name, idParam, bidderParam.Type, param.Value)...)
		}
		// Xử lý các param được add mới trong line item
		for paramName, param := range bidderInfo.Params {
			if param.IsAddParam != mysql.TypeOn {
				continue
			}
			idParam := strconv.FormatInt(bidderInfo.BidderId, 10) + "-" + paramName + "-" + strconv.Itoa(bidderInfo.BidderIndex)
			required := new(PbBidder).IsParamRequiredByBidder(bidderInfo.BidderName, paramName)
			if required && govalidator.IsNull(param.Value) {
				errs = append(errs, ajax.Error{
					Id:      idParam,
					Message: "(*) Required",
				})
			}
			errs = append(errs, t.validateValueOfType(paramName, idParam, param.Type, param.Value)...)
		}
	}
	return
}

func (t *LineItemV2) validateValueOfType(paramName string, idParam string, typ string, value string) (errs []ajax.Error) {
	switch typ {
	case "int":
		if value != "" {
			if !govalidator.IsInt(value) {
				errs = append(errs, ajax.Error{
					Id:      idParam,
					Message: "param " + paramName + " value is int",
				})
			}
		}
		break
	case "float":
		if value != "" {
			if !govalidator.IsFloat(value) {
				errs = append(errs, ajax.Error{
					Id:      idParam,
					Message: "param " + paramName + " value is float",
				})
			}
		}
		break
	case "json":
		if value != "" {
			if !govalidator.IsJSON(value) {
				errs = append(errs, ajax.Error{
					Id:      idParam,
					Message: "param " + paramName + " value is json",
				})
			}
		}
	case "boolean":
		if value == "" {
			value = "false"
		} else if value != "true" && value != "false" {
			errs = append(errs, ajax.Error{
				Id:      idParam,
				Message: "param " + paramName + " value is true and false",
			})
		}
		break
	}
	return
}

func (t *LineItemV2) Delete(id, userId int64, userAdmin UserRecord, lang lang.Translation) fiber.Map {
	record, _ := new(LineItemV2).GetById(id, userId)
	if record.AutoCreate == mysql.On {
		return fiber.Map{
			"status":  "err",
			"message": "you cannot delete this line item",
			"id":      id,
		}
	}
	err := mysql.Client.Model(&LineItemRecordV2{}).Delete(&LineItemRecordV2{}, "id = ? and user_id = ?", id, userId).Error
	if err != nil {
		if !utility.IsWindow() {
			return fiber.Map{
				"status":  "err",
				"message": lang.Errors.LineItemError.Delete.ToString(),
				"id":      id,
			}
		}
		return fiber.Map{
			"status":  "err",
			"message": err.Error(),
			"id":      id,
		}
	} else {
		//// History
		//var creatorId int64
		//if userAdmin.Id != 0 {
		//	creatorId = userAdmin.Id
		//} else {
		//	creatorId = userId
		//}
		//_ = history.PushHistory(&history.LineItem{
		//	Detail:    history.DetailLineItemFE,
		//	CreatorId: creatorId,
		//	RecordOld: record.TableLineItem,
		//	RecordNew: mysql.TableLineItem{},
		//})

		new(Inventory).UpdateRenderCacheWithLineItem(id, userId)
		return fiber.Map{
			"status":  "success",
			"message": "done",
			"id":      id,
		}
	}
}

func (t *LineItemV2) GetByFilters(inputs *payload.LineItemFilterPayloadV2, user UserRecord, lang lang.Translation) (response datatable.Response, err error) {
	var bidders []LineItemRecordV2
	var total int64
	err = mysql.Client.Where("line_item_type = 1 and user_id = ?", user.Id).
		Scopes(
			t.SetFilterStatus(inputs),
			t.setFilterSearch(inputs),
			t.setFilterType(inputs),
			t.setFilterTarget(inputs, user),
		).
		Model(&bidders).Count(&total).
		Scopes(
			t.setOrder(inputs),
			pagination.Paginate(pagination.Params{
				Limit:  inputs.Length,
				Offset: inputs.Start,
			}),
		).
		Select("line_item_v2.*").
		Group("line_item_v2.id").
		Find(&bidders).Error
	if err != nil {
		if !utility.IsWindow() {
			err = fmt.Errorf(lang.Errors.LineItemError.List.ToString())
		}
		return datatable.Response{}, err
	}
	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.MakeResponseDatatable(bidders, user.Id)
	return
}

type BidderRecordDatatableV2 struct {
	LineItemRecordV2
	RowId  string `json:"DT_RowId"`
	Name   string `json:"name"`
	Target string `json:"target"`
	Status string `json:"status"`
	Type   string `json:"server_type"`
	Rate   string `json:"rate"`
	Action string `json:"action"`
}

type TargetsV2 struct {
	LineItem          LineItemRecordV2
	HtmlContent       string
	TextListInventory string
	TextListAdTag     string
	TextListAdSize    string
	TextListAdFormat  string
	TextListGeo       string
	TextListDevice    string
}

func (t *LineItemV2) MakeResponseDatatable(bidders []LineItemRecordV2, userId int64) (records []BidderRecordDatatableV2) {
	for _, bidder := range bidders {
		target := t.getTarget(bidder, userId)
		stringHtml := htmlblock.Render("line-item-v2/index/block.target.gohtml", target).String()
		target.HtmlContent = stringHtml

		//fmt.Println(target.StringInventory)
		rec := BidderRecordDatatableV2{
			LineItemRecordV2: bidder,
			RowId:            "bidder_" + strconv.FormatInt(bidder.Id, 10),
			Name:             htmlblock.Render("line-item-v2/index/block.name.gohtml", target).String(),
			Target:           htmlblock.Render("line-item-v2/index/block.target-button.gohtml", target).String(),
			Status:           htmlblock.Render("line-item-v2/index/block.status.gohtml", bidder).String(),
			Type:             htmlblock.Render("line-item-v2/index/block.type.gohtml", bidder).String(),
			Rate:             htmlblock.Render("line-item-v2/index/block.rate.gohtml", bidder).String(),
			Action:           htmlblock.Render("line-item-v2/index/block.action.gohtml", bidder).String(),
		}
		records = append(records, rec)
	}
	return
}

func (t *LineItemV2) SetFilterStatus(inputs *payload.LineItemFilterPayloadV2) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.Status != nil {
			switch inputs.PostData.Status.(type) {
			case string, int:
				if inputs.PostData.Status != "" {
					return db.Where("status = ?", inputs.PostData.Status)
				}
			case []string, []interface{}:
				return db.Where("status IN ?", inputs.PostData.Status)
			}
		}
		return db
	}
}

func (t *LineItemV2) setFilterType(inputs *payload.LineItemFilterPayloadV2) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.Type != nil {
			switch inputs.PostData.Type.(type) {
			case string, int:
				if inputs.PostData.Type != "" {
					return db.Where("server_type = ?", inputs.PostData.Type)
				}
			case []string, []interface{}:
				return db.Where("server_type IN ?", inputs.PostData.Type)
			}
		}
		return db
	}
}

func (t *LineItemV2) setFilterSearch(inputs *payload.LineItemFilterPayloadV2) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var flag bool
		// Search from form of datatable <- not use
		if inputs.Search != nil && inputs.Search.Value != "" {
			flag = true
		}
		// Search from form filter
		if inputs.PostData.QuerySearch != "" {
			flag = true
		}
		if !flag {
			return db
		}
		return db.Where("name LIKE ?", "%"+inputs.PostData.QuerySearch+"%")
	}
}

func (t *LineItemV2) setOrder(inputs *payload.LineItemFilterPayloadV2) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var orders []string
		orders = append(orders, "is_lock DESC")
		if len(inputs.Order) <= 0 {
			orders = append(orders, "id DESC")
		} else {
			for _, order := range inputs.Order {
				column := inputs.Columns[order.Column]
				orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
			}
		}
		var orderString string
		if inputs.PostData.Domain != nil {
			orderString = strings.Join(orders, ", ")
			orderString = "line_item_v2." + orderString
		} else {
			orderString = strings.Join(orders, ", ")
		}
		return db.Order(orderString)
	}
}

func (t *LineItemV2) setFilterTarget(inputs *payload.LineItemFilterPayloadV2, user UserRecord) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var listLineItemFilterAdTag []int64
		if inputs.PostData.AdTag != nil {
			var listDomain, listFormat, listSize []string
			switch inputs.PostData.AdTag.(type) {
			case string:
				if inputs.PostData.AdTag != "" {
					adTagId, _ := strconv.ParseInt(inputs.PostData.AdTag.(string), 10, 64)
					adTag := new(InventoryAdTag).GetById(adTagId)
					adTag.GetFullData()
					mysql.Client.
						Table(mysql.Tables.TargetV2).
						Select("line_item_id").
						Where("line_item_id != 0 and (tag_id = -1 or tag_id = ?) and user_id = ?", inputs.PostData.AdTag, user.Id).
						Group("line_item_id").
						Find(&listLineItemFilterAdTag)

					if !utility.InArray(strconv.FormatInt(adTag.InventoryId, 10), listDomain, true) {
						listDomain = append(listDomain, strconv.FormatInt(adTag.InventoryId, 10))
					}
					if !utility.InArray(strconv.FormatInt(adTag.Type.Int(), 10), listFormat, true) {
						listFormat = append(listFormat, strconv.FormatInt(adTag.Type.Int(), 10))
					}
					if adTag.Type.IsBanner() {
						if !utility.InArray(strconv.FormatInt(adTag.PrimaryAdSize, 10), listSize, true) {
							listSize = append(listSize, strconv.FormatInt(adTag.PrimaryAdSize, 10))
						}
						for _, size := range adTag.AdditionalAdSize {
							if !utility.InArray(strconv.FormatInt(size.Id, 10), listSize, true) {
								listSize = append(listSize, strconv.FormatInt(size.Id, 10))
							}
						}
					}
				}
			case []interface{}:
				for _, v := range inputs.PostData.AdTag.([]interface{}) {
					adTagId, _ := strconv.ParseInt(v.(string), 10, 64)
					adTag := new(InventoryAdTag).GetById(adTagId)
					adTag.GetFullData()
					if !utility.InArray(strconv.FormatInt(adTag.InventoryId, 10), listDomain, true) {
						listDomain = append(listDomain, strconv.FormatInt(adTag.InventoryId, 10))
					}
					if !utility.InArray(strconv.FormatInt(adTag.Type.Int(), 10), listFormat, true) {
						listFormat = append(listFormat, strconv.FormatInt(adTag.Type.Int(), 10))
					}
					if adTag.Type.IsBanner() {
						if !utility.InArray(strconv.FormatInt(adTag.PrimaryAdSize, 10), listSize, true) {
							listSize = append(listSize, strconv.FormatInt(adTag.PrimaryAdSize, 10))
						}
						for _, size := range adTag.AdditionalAdSize {
							if !utility.InArray(strconv.FormatInt(size.Id, 10), listSize, true) {
								listSize = append(listSize, strconv.FormatInt(size.Id, 10))
							}
						}
					}
				}
				mysql.Client.
					Table(mysql.Tables.TargetV2).
					Select("line_item_id").
					Where("line_item_id != 0 and (tag_id = -1 or tag_id in ?) and user_id = ?", inputs.PostData.AdTag, user.Id).
					Group("line_item_id").
					Find(&listLineItemFilterAdTag)
			}
			inputs.PostData.Domain = listDomain
			inputs.PostData.AdFormat = listFormat
			inputs.PostData.AdSize = listSize
		} else {
			mysql.Client.
				Table(mysql.Tables.TargetV2).
				Select("line_item_id").
				Where("line_item_id != 0 and tag_id != 0 and user_id = ?", user.Id).
				Group("line_item_id").
				Find(&listLineItemFilterAdTag)
		}
		//fmt.Println(listLineItemFilterAdTag)

		var listLineItemFilterDomain []int64
		if inputs.PostData.Domain != nil {
			switch inputs.PostData.Domain.(type) {
			case string, int:
				if inputs.PostData.Domain != "" {
					mysql.Client.
						Table(mysql.Tables.TargetV2).
						Select("line_item_id").
						Where("line_item_id != 0 and (inventory_id = -1 or inventory_id = ?) and user_id = ? and line_item_id IN ?", inputs.PostData.Domain, user.Id, listLineItemFilterAdTag).
						Group("line_item_id").
						Find(&listLineItemFilterDomain)
				}
			case []string, []interface{}:
				mysql.Client.
					Table(mysql.Tables.TargetV2).
					Select("line_item_id").
					Where("line_item_id != 0 and (inventory_id = -1 or inventory_id in ?) and user_id = ? and line_item_id IN ?", inputs.PostData.Domain, user.Id, listLineItemFilterAdTag).
					Group("line_item_id").
					Find(&listLineItemFilterDomain)
			}
		} else {
			mysql.Client.
				Table(mysql.Tables.TargetV2).
				Select("line_item_id").
				Where("line_item_id != 0 and user_id = ? and line_item_id IN ?", user.Id, listLineItemFilterAdTag).
				Group("line_item_id").
				Find(&listLineItemFilterDomain)
		}
		//fmt.Println(listLineItemFilterDomain)

		var listLineItemFilterFormat []int64
		if inputs.PostData.AdFormat != nil {
			switch inputs.PostData.AdFormat.(type) {
			case string, int:
				if inputs.PostData.AdFormat != "" {
					mysql.Client.
						Table(mysql.Tables.TargetV2).
						Select("line_item_id").
						Where("line_item_id != 0 and (ad_format_id = -1 or ad_format_id = ?) and user_id = ? and line_item_id IN ?", inputs.PostData.AdFormat, user.Id, listLineItemFilterDomain).
						Group("line_item_id").
						Find(&listLineItemFilterFormat)
				}
			case []string, []interface{}:
				mysql.Client.
					Table(mysql.Tables.TargetV2).
					Select("line_item_id").
					Where("line_item_id != 0 and (ad_format_id = -1 or ad_format_id in ?) and user_id = ? and line_item_id IN ?", inputs.PostData.AdFormat, user.Id, listLineItemFilterDomain).
					Group("line_item_id").
					Find(&listLineItemFilterFormat)
			}
		} else {
			mysql.Client.
				Table(mysql.Tables.TargetV2).
				Select("line_item_id").
				Where("line_item_id != 0 and user_id = ? and line_item_id IN ?", user.Id, listLineItemFilterDomain).
				Group("line_item_id").
				Find(&listLineItemFilterFormat)
			//fmt.Println(listLineItemFilterAdTag)
		}

		var listLineItemFilterSize []int64
		if inputs.PostData.AdSize != nil {
			switch inputs.PostData.AdSize.(type) {
			case string, int:
				if inputs.PostData.AdSize != "" {
					mysql.Client.
						Table(mysql.Tables.TargetV2).
						Select("line_item_id").
						Where("line_item_id != 0 and (ad_size_id = -1 or ad_size_id = ?) and user_id = ? and line_item_id IN ?", inputs.PostData.AdSize, user.Id, listLineItemFilterFormat).
						Group("line_item_id").
						Find(&listLineItemFilterSize)
				}
			case []string, []interface{}:
				mysql.Client.
					Table(mysql.Tables.TargetV2).
					Select("line_item_id").
					Where("line_item_id != 0 and (ad_size_id = -1 or ad_size_id in ?) and user_id = ? and line_item_id IN ?", inputs.PostData.AdSize, user.Id, listLineItemFilterFormat).
					Group("line_item_id").
					Find(&listLineItemFilterSize)
			}
		} else {
			mysql.Client.
				Table(mysql.Tables.TargetV2).
				Select("line_item_id").
				Where("line_item_id != 0 and user_id = ? and line_item_id IN ?", user.Id, listLineItemFilterFormat).
				Group("line_item_id").
				Find(&listLineItemFilterSize)
			//fmt.Println(listLineItemFilterAdTag)
		}

		var listLineItemFilterGeo []int64
		if inputs.PostData.Country != nil {
			switch inputs.PostData.Country.(type) {
			case string, int:
				if inputs.PostData.Country != "" {
					mysql.Client.
						Table(mysql.Tables.TargetV2).
						Select("line_item_id").
						Where("line_item_id != 0 and (geo_id = -1 or geo_id = ?) and user_id = ? and line_item_id IN ?", inputs.PostData.Country, user.Id, listLineItemFilterSize).
						Group("line_item_id").
						Find(&listLineItemFilterGeo)
				}
			case []string, []interface{}:
				mysql.Client.
					Table(mysql.Tables.TargetV2).
					Select("line_item_id").
					Where("line_item_id != 0 and (geo_id = -1 or geo_id in ?) and user_id = ? and line_item_id IN ?", inputs.PostData.Country, user.Id, listLineItemFilterSize).
					Group("line_item_id").
					Find(&listLineItemFilterGeo)
			}
		} else {
			mysql.Client.
				Table(mysql.Tables.TargetV2).
				Select("line_item_id").
				Where("line_item_id != 0 and user_id = ? and line_item_id IN ?", user.Id, listLineItemFilterSize).
				Group("line_item_id").
				Find(&listLineItemFilterGeo)
			//fmt.Println(listLineItemFilterAdTag)
		}

		var listLineItemFilter []int64
		if inputs.PostData.Device != nil {
			switch inputs.PostData.Device.(type) {
			case string, int:
				if inputs.PostData.Device != "" {
					mysql.Client.
						Table(mysql.Tables.TargetV2).
						Select("line_item_id").
						Where("line_item_id != 0 and (device_id = -1 or device_id = ?) and user_id = ? and line_item_id IN ?", inputs.PostData.Country, user.Id, listLineItemFilterGeo).
						Group("line_item_id").
						Find(&listLineItemFilter)
				}
			case []string, []interface{}:
				mysql.Client.
					Table(mysql.Tables.TargetV2).
					Select("line_item_id").
					Where("line_item_id != 0 and (device_id = -1 or device_id in ?) and user_id = ? and line_item_id IN ?", inputs.PostData.Country, user.Id, listLineItemFilterGeo).
					Group("line_item_id").
					Find(&listLineItemFilter)
			}
		} else {
			mysql.Client.
				Table(mysql.Tables.TargetV2).
				Select("line_item_id").
				Where("line_item_id != 0 and user_id = ? and line_item_id IN ?", user.Id, listLineItemFilterGeo).
				Group("line_item_id").
				Find(&listLineItemFilter)
			//fmt.Println(listLineItemFilterAdTag)
		}

		return db.Where("line_item_v2.id in ?", listLineItemFilter)
	}
}

func (t *LineItemV2) getTarget(lineItem LineItemRecordV2, userId int64) TargetsV2 {
	assignTarget := TargetsV2{
		LineItem: lineItem,
	}
	var listNameInventoryTarget []string
	var listNameAdFormatTarget []string
	var listNameAdSizeTarget []string
	var listNameAdTagTarget []string
	var listNameGeoTarget []string
	var listNameDeviceTarget []string

	targets := new(TargetV2).GetTargetLineItem(lineItem.Id)

	mapInventory := make(map[int64]int)
	mapAdFormat := make(map[int64]int)
	mapAdSize := make(map[int64]int)
	mapAdTag := make(map[int64]int)
	mapGeo := make(map[int64]int)
	mapDevice := make(map[int64]int)
	for _, target := range targets {
		if target.InventoryId != 0 {
			mapInventory[target.InventoryId] = 1
		}
		if target.AdFormatId != 0 {
			mapAdFormat[target.AdFormatId] = 1
		}
		if target.AdSizeId != 0 {
			mapAdSize[target.AdSizeId] = 1
		}
		if target.TagId != 0 {
			mapAdTag[target.TagId] = 1
		}
		if target.GeoId != 0 {
			mapGeo[target.GeoId] = 1
		}
		if target.DeviceId != 0 {
			mapDevice[target.DeviceId] = 1
		}
	}

	// Lọc bỏ những id trùng nhau
	for inventoryId, _ := range mapInventory {
		inventory, _ := new(Inventory).GetById(inventoryId, userId)
		listNameInventoryTarget = append(listNameInventoryTarget, inventory.Name)
	}
	for adFormatId, _ := range mapAdFormat {
		adFormat := new(AdType).GetById(adFormatId)
		listNameAdFormatTarget = append(listNameAdFormatTarget, adFormat.Name)
	}
	for adSizeId, _ := range mapAdSize {
		adSize := new(AdSize).GetById(adSizeId)
		listNameAdSizeTarget = append(listNameAdSizeTarget, adSize.Name)
	}
	for adTagId, _ := range mapAdTag {
		adTag := new(InventoryAdTag).GetById(adTagId)
		listNameAdTagTarget = append(listNameAdTagTarget, adTag.Name)
	}
	for geoId, _ := range mapGeo {
		geo := new(Country).GetById(geoId)
		listNameGeoTarget = append(listNameGeoTarget, geo.Name)
	}
	for deviceId, _ := range mapDevice {
		device := new(Device).GetById(deviceId)
		listNameDeviceTarget = append(listNameDeviceTarget, device.Name)
	}

	assignTarget.TextListInventory = strings.Join(listNameInventoryTarget, ", ")
	assignTarget.TextListAdTag = strings.Join(listNameAdTagTarget, ", ")
	assignTarget.TextListAdSize = strings.Join(listNameAdSizeTarget, ", ")
	assignTarget.TextListAdFormat = strings.Join(listNameAdFormatTarget, ", ")
	assignTarget.TextListGeo = strings.Join(listNameGeoTarget, ", ")
	assignTarget.TextListDevice = strings.Join(listNameDeviceTarget, ", ")
	return assignTarget
}

func (t *LineItemV2) CreateAdsenseAdSlot(lineItemId int64, inputs payload.LineItemAddV2) (errs []ajax.Error) {
	for _, adsenseAdSlot := range inputs.AdsenseAdSlots {
		if err := new(LineItemAdsenseAdSlot).Push(lineItemId, adsenseAdSlot.Size, adsenseAdSlot.AdSlotId); err != nil {
			errs = append(errs, ajax.Error{
				Id:      `adsense-ad-slot-` + adsenseAdSlot.Size,
				Message: err.Error(),
			})
		}
	}
	return
}

func (t *LineItemV2) CreateBidderInfo(lineItemId int64, inputs payload.LineItemAddV2) (err error) {
	for _, bidderInfo := range inputs.BidderParams {
		// Bỏ qua các bidder không có param
		if len(bidderInfo.Params) < 1 {
			continue
		}
		// Tạo bảng bidder info trước rồi lấy id của bảng này để tạo rls
		recBidderInfo := new(LineItemBidderInfoV2).CreateBidderInfo(LineItemBidderInfoRecordV2{mysql.TableLineItemBidderInfoV2{
			LineItemId: lineItemId,
			BidderId:   bidderInfo.BidderId,
			ConfigType: bidderInfo.ConfigType,
			BidderType: bidderInfo.BidderType,
			Name:       bidderInfo.BidderName,
		}})
		// Lấy toàn bộ các param theo id của bidder
		bidderParams := new(BidderParams).GetByBidderId(bidderInfo.BidderId)
		for _, bidderParam := range bidderParams {
			// Kiểm tra xem param có nằm trong param add lên không nếu không thì bỏ qua nếu có thì lấy value để tạo bảng row line_item_bidder_params
			param, exists := bidderInfo.Params[bidderParam.Name]
			//fmt.Println(bidderInfo.Params)
			//fmt.Println(param)
			//fmt.Println(exists)
			if !exists {
				continue
			}
			//fmt.Println(param)
			err = mysql.Client.Create(&LineItemBidderParamsRecordV2{mysql.TableLineItemBidderParamsV2{
				LineItemId:       lineItemId,
				BidderId:         bidderParam.BidderId,
				LineItemBidderId: recBidderInfo.Id,
				Name:             bidderParam.Name,
				Type:             bidderParam.Type,
				Value:            param.Value,
				//IsAddParam:       mysql.TypeOff,
			}}).Error
			delete(bidderInfo.Params, bidderParam.Name)
		}
		// Xử lý các param được add mới trong line item
		for paramName, param := range bidderInfo.Params {
			err = mysql.Client.Create(&LineItemBidderParamsRecordV2{mysql.TableLineItemBidderParamsV2{
				LineItemId:       lineItemId,
				BidderId:         bidderInfo.BidderId,
				LineItemBidderId: recBidderInfo.Id,
				Name:             paramName,
				Type:             param.Type,
				Value:            param.Value,
				//IsAddParam:       mysql.TypeOff,
			}}).Error
		}
	}
	return
}

func (t *LineItemV2) CreateTarget(lineItemId, userId int64, inputs payload.LineItemAddV2) (err error) {
	all := int64(-1)
	// Kiểm tra nếu đầu vào input list target = 0 thì thêm một target = 0 thể hiện select all
	if len(inputs.ListAdInventory) == 0 {
		inputs.ListAdInventory = []payload.ListTargetV2{
			{
				Id: all,
			},
		}
	}
	if len(inputs.ListAdFormat) == 0 {
		inputs.ListAdFormat = []payload.ListTargetV2{
			{
				Id: all,
			},
		}
	}
	if len(inputs.ListGeo) == 0 {
		inputs.ListGeo = []payload.ListTargetV2{
			{
				Id: all,
			},
		}
	}
	if len(inputs.ListAdSize) == 0 {
		inputs.ListAdSize = []payload.ListTargetV2{
			{
				Id: all,
			},
		}
	}
	if len(inputs.ListAdTag) == 0 {
		inputs.ListAdTag = []payload.ListTargetV2{
			{
				Id: all,
			},
		}
	}
	if len(inputs.ListDevice) == 0 {
		inputs.ListDevice = []payload.ListTargetV2{
			{
				Id: all,
			},
		}
	}
	//If type là gg thì mặc định size và device là all
	if inputs.ServerType == mysql.TYPEServerTypeGoogle {
		inputs.ListAdSize = []payload.ListTargetV2{
			{
				Id: all,
			},
		}
		inputs.ListDevice = []payload.ListTargetV2{
			{
				Id: all,
			},
		}
	}
	for _, inventory := range inputs.ListAdInventory {
		var recordTarget TargetRecordV2
		err := mysql.Client.Table(mysql.Tables.TargetV2).FirstOrCreate(&recordTarget, TargetRecordV2{mysql.TableTargetV2{
			UserId:      userId,
			LineItemId:  lineItemId,
			InventoryId: inventory.Id,
		}}).Error
		if err != nil {
			return err
		}
	}
	for _, adFormat := range inputs.ListAdFormat {
		var recordTarget TargetRecordV2
		err := mysql.Client.Table(mysql.Tables.TargetV2).FirstOrCreate(&recordTarget, TargetRecordV2{mysql.TableTargetV2{
			UserId:     userId,
			LineItemId: lineItemId,
			AdFormatId: adFormat.Id,
		}}).Error
		if err != nil {
			return err
		}
	}
	for _, size := range inputs.ListAdSize {
		var recordTarget TargetRecordV2
		err := mysql.Client.Table(mysql.Tables.TargetV2).FirstOrCreate(&recordTarget, TargetRecordV2{mysql.TableTargetV2{
			UserId:     userId,
			LineItemId: lineItemId,
			AdSizeId:   size.Id,
		}}).Error
		if err != nil {
			return err
		}
	}

	for _, geo := range inputs.ListGeo {
		var recordTarget TargetRecordV2
		err := mysql.Client.Table(mysql.Tables.TargetV2).FirstOrCreate(&recordTarget, TargetRecordV2{mysql.TableTargetV2{
			UserId:     userId,
			LineItemId: lineItemId,
			GeoId:      geo.Id,
		}}).Error
		if err != nil {
			return err
		}
	}
	for _, device := range inputs.ListDevice {
		var recordTarget TargetRecordV2
		err := mysql.Client.Table(mysql.Tables.TargetV2).FirstOrCreate(&recordTarget, TargetRecordV2{mysql.TableTargetV2{
			UserId:     userId,
			LineItemId: lineItemId,
			DeviceId:   device.Id,
		}}).Error
		if err != nil {
			return err
		}
	}

	for _, adTag := range inputs.ListAdTag {
		var recordTarget TargetRecordV2
		err := mysql.Client.Table(mysql.Tables.TargetV2).FirstOrCreate(&recordTarget, TargetRecordV2{mysql.TableTargetV2{
			UserId:     userId,
			LineItemId: lineItemId,
			TagId:      adTag.Id,
		}}).Error
		if err != nil {
			return err
		}
	}

	return
}

func (t *LineItemV2) GetInventoryByName(search string, userId int64) (rows []SearchLineItemV2) {
	//q := "%" + search + "%"
	//mysql.Client.Raw("SELECT id,name FROM `inventory` WHERE name like ? and user_id = ? limit 20", q, userId).Find(&rows)
	query := mysql.Client.Table("inventory").Select("id,name").Where("user_id = ? AND deleted_at IS NULL", userId)
	if search != "" {
		query = query.Where("name LIKE %?%", search)
	}
	query.Find(&rows)
	return
}

func (t *LineItemV2) GetAdFormatByName(search string) (rows []SearchLineItemV2) {
	q := "%" + search + "%"
	mysql.Client.Raw("SELECT id,name FROM `ad_type` WHERE name like ? limit 20", q).Find(&rows)
	return
}

func (t *LineItemV2) GetAdSizeByName(search string) (rows []SearchLineItemV2) {
	q := "%" + search + "%"
	mysql.Client.Raw("SELECT id,name FROM `ad_size` WHERE name like ? limit 10", q).Find(&rows)
	return
}

func (t *LineItemV2) GetAdTagByName(search string, userId int64) (rows []SearchLineItemV2) {
	q := "%" + search + "%"
	mysql.Client.Raw("SELECT id,name FROM `inventory_ad_tag` WHERE name like ? and user_id = ? limit 10", q, userId).Find(&rows)
	return
}

func (t *LineItemV2) GetDeviceByName(search string) (rows []SearchLineItemV2) {
	q := "%" + search + "%"
	mysql.Client.Raw("SELECT id,name FROM `devices` WHERE name like ? limit 10", q).Find(&rows)
	return
}

func (t *LineItemV2) GetCountryByName(search string) (rows []SearchLineItemV2) {
	q := "%" + search + "%"
	mysql.Client.Raw("SELECT id,name FROM `country` WHERE name like ? limit 10", q).Find(&rows)
	return
}

func (t *LineItemV2) GetSearched(typ string, search []string, userId int64) (rows []SearchLineItemV2) {
	switch typ {
	case "domain":
		mysql.Client.Raw("SELECT id,name FROM `inventory` WHERE id in ? and user_id = ?", search, userId).Find(&rows)
		//query := mysql.Client.Table("inventory").
		//	Select("id, name").Where("user_id = ? AND deleted_at IS NULL", userId)
		////if len(search) > 0 {
		//query = query.Where("id IN ?", search)
		////}
		//query.Find(&rows)
		break
	case "format":
		mysql.Client.Raw("SELECT id,name FROM `ad_type` WHERE id in ?", search).Find(&rows)
		break
	case "size":
		mysql.Client.Raw("SELECT id,name FROM `ad_size` WHERE id in ?", search).Find(&rows)
		break
	case "adtag":
		mysql.Client.Raw("SELECT id,name FROM `inventory_ad_tag` WHERE id in ? and user_id = ?", search, userId).Find(&rows)
		break
	case "country":
		mysql.Client.Raw("SELECT id,name FROM `country` WHERE id in ?", search).Find(&rows)
		break
	case "device":
		mysql.Client.Raw("SELECT id,name FROM `devices` WHERE id in ?", search).Find(&rows)
		break
	}
	return
}

func (t *LineItemV2) GetOfUserById(id, userId int64) (row LineItemRecordV2) {
	mysql.Client.Model(&LineItemRecordV2{}).Where("id = ? and user_id = ?", id, userId).Find(&row)
	//Get Rls
	row.GetRls()
	return
}

func (t *LineItemV2) InArray(array []SearchLineItemV2, id int64) (index int, flag bool) {
	for i, v := range array {
		if v.Id == id {
			return i, true
		}
	}
	return 0, false
}

func (t *LineItemV2) GetFilter(option string, userId int64, searched []string) (rows []SearchLineItemV2) {
	switch option {
	case "domain":
		rows = t.GetInventoryByName("", userId)
		search := t.GetSearched(option, searched, userId)
		for i, v := range search {
			index, flag := t.InArray(rows, v.Id)
			if flag {
				rows[index].Selected = true
			} else {
				search[i].Selected = true
				rows = append(rows, search[i])
			}
		}
		break
	case "format":
		rows = t.GetAdFormatByName("")
		search := t.GetSearched(option, searched, userId)
		for i, v := range search {
			index, flag := t.InArray(rows, v.Id)
			if flag {
				rows[index].Selected = true
			} else {
				search[i].Selected = true
				rows = append(rows, search[i])
			}
		}
		break
	case "size":
		rows = t.GetAdSizeByName("")
		search := t.GetSearched(option, searched, userId)
		for i, v := range search {
			index, flag := t.InArray(rows, v.Id)
			if flag {
				rows[index].Selected = true
			} else {
				search[i].Selected = true
				rows = append(rows, search[i])
			}
		}
		break
	case "adtag":
		rows = t.GetAdTagByName("", userId)
		search := t.GetSearched(option, searched, userId)
		for i, v := range search {
			index, flag := t.InArray(rows, v.Id)
			if flag {
				rows[index].Selected = true
			} else {
				search[i].Selected = true
				rows = append(rows, search[i])
			}
		}
		break
	case "country":
		rows = t.GetCountryByName("")
		search := t.GetSearched(option, searched, userId)
		for i, v := range search {
			index, flag := t.InArray(rows, v.Id)
			if flag {
				rows[index].Selected = true
			} else {
				search[i].Selected = true
				rows = append(rows, search[i])
			}
		}
		break
	case "device":
		rows = t.GetDeviceByName("")
		search := t.GetSearched(option, searched, userId)
		for i, v := range search {
			index, flag := t.InArray(rows, v.Id)
			if flag {
				rows[index].Selected = true
			} else {
				search[i].Selected = true
				rows = append(rows, search[i])
			}
		}
		break
	}
	return
}

func (t *LineItemV2) GetListBoxCollapse(userId, lineItemId int64, page, typ string) (list []string) {
	switch typ {
	case "add":
		mysql.Client.Select("box_collapse").Model(PageCollapseRecord{}).Where("user_id = ? and page_collapse = ? and is_collapse = ? and page_type = ?", userId, page, 1, typ).Find(&list)
		return
	case "edit":
		mysql.Client.Select("box_collapse").Model(PageCollapseRecord{}).Where("user_id = ? and page_collapse = ? and is_collapse = ? and page_type = ? and page_id = ?", userId, page, 1, typ, lineItemId).Find(&list)
		return
	}
	return
}

func (t *LineItemV2) Count(id string, listId []string) bool {
	var count int
	for _, v := range listId {
		if v == id {
			count++
			if count > 1 {
				return false
			}
		}
	}
	return true
}

func (t *LineItemV2) CheckParamRequired(name, param string) (check bool) {
	recPbBidder := new(PbBidder).GetPbBidderByBidderCode(name)
	recPbBidderParam := new(PbBidder).GetPbBidderParamCheckRequired(recPbBidder.Id, param)
	if recPbBidderParam.Scope == "required" {
		check = true
	}
	return
}

func (t *LineItemV2) PushToQueueWorkerLineItemDfp(lineItemId int64) {
	//Push to queue worker với lineId và type là google
	mysql.Client.Model(LineItemRecordV2{}).Where("id = ? and server_type = 2", lineItemId).Update("push_line_item_dfp", 1)
	return
}

func (t *LineItemV2) BuildLineNameAdsenseDisplay(inventoryName string, account string) string {
	return "[Default] - " + inventoryName + " - Adsense(" + account + ") - Display"
}

func (t *LineItemV2) BuildLineNameAdsenseVideo(inventoryName string, account string) string {
	return "[Default] - " + inventoryName + " - Adsense(" + account + ") - Video"
}

func (t *LineItemV2) GetByName(name string) (record LineItemRecordV2) {
	mysql.Client.Where("name = ?", name).Find(&record)
	return
}
