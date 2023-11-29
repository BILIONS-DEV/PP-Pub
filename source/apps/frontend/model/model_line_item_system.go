package model
//
//import (
//	"fmt"
//	"github.com/asaskevich/govalidator"
//	"github.com/gofiber/fiber/v2"
//	"gorm.io/gorm"
//	"source/apps/frontend/lang"
//	"source/apps/frontend/payload"
//	"source/core/technology/mysql"
//	"source/pkg/ajax"
//	"source/pkg/datatable"
//	"source/pkg/htmlblock"
//	"source/pkg/pagination"
//	"source/pkg/utility"
//	"strconv"
//	"strings"
//	"time"
//)
//
//type LineItemSystem struct{}
//
//func (t *LineItemSystem) makeRecord(payload payload.LineItemSystemAdd) (rec LineItemRecord, err error) {
//	rec = LineItemRecord{mysql.TableLineItem{
//		Id:           payload.Id,
//		UserId:       0,
//		Name:         payload.Name,
//		Description:  payload.Description,
//		ServerType:   payload.ServerType,
//		Type:         payload.BidderType,
//		LinkedGam:    payload.LinkedGam,
//		Rate:         payload.BidderRate,
//		VastUrl:      payload.BidderVastUrl,
//		AdTag:        payload.BidderAdTag,
//		Priority:     payload.Priority,
//		LineItemType: 2,
//	}}
//
//	layoutISO := "01/02/2006"
//	if payload.StartDate != "" {
//		startDate, err := time.Parse(layoutISO, payload.StartDate)
//		if err != nil {
//			return rec, err
//		}
//		err = rec.StartDate.Scan(startDate)
//		if err != nil {
//			return rec, err
//		}
//	}
//	if payload.EndDate != "" {
//		endDate, err := time.Parse(layoutISO, payload.EndDate)
//		if err != nil {
//			return rec, err
//		}
//		err = rec.EndDate.Scan(endDate)
//		if err != nil {
//			return rec, err
//		}
//	}
//
//	//rec.StartDate.Time = startDate
//	//rec.EndDate.Time = endDate
//	if payload.Status == "on" {
//		rec.Status = 1
//	} else {
//		rec.Status = 2
//	}
//	return
//}
//
//func (t *LineItemSystem) AddLineItem(inputs payload.LineItemSystemAdd, user UserRecord) (record LineItemRecord, errs []ajax.Error) {
//	// Validate inputs
//	errs = t.Validate(inputs)
//	if len(errs) > 0 {
//		return
//	}
//	var err error
//	// Make record use insert to database
//	record, err = t.makeRecord(inputs)
//	if err != nil {
//		if !utility.IsWindow() {
//			errs = append(errs, ajax.Error{
//				Id:      "",
//				Message: lang.Translate.Errors.LineItemError.Add.ToString(),
//			})
//		} else {
//			errs = append(errs, ajax.Error{
//				Id:      "",
//				Message: err.Error(),
//			})
//		}
//		return
//	}
//
//	// Insert to database
//	err = mysql.Client.Create(&record).Error
//	if err != nil {
//		if !utility.IsWindow() {
//			errs = append(errs, ajax.Error{
//				Id:      "",
//				Message: lang.Translate.Errors.LineItemError.Add.ToString(),
//			})
//		} else {
//			errs = append(errs, ajax.Error{
//				Id:      "",
//				Message: err.Error(),
//			})
//		}
//		return
//	}
//
//	switch inputs.ServerType {
//	case mysql.TYPEServerTypeGoogle: //=> CASE: LineItemType là GOOGLE
//		//Tạo line_item account nếu server là google
//		for _, v := range inputs.SelectAccount {
//			if v != "" {
//				bidderId, _ := strconv.ParseInt(v, 10, 64)
//				err = new(LineItemAccount).Create(record.Id, bidderId)
//				if err != nil {
//					if !utility.IsWindow() {
//						errs = append(errs, ajax.Error{
//							Id:      "",
//							Message: lang.Translate.Errors.LineItemError.LineItemAccount.ToString(),
//						})
//					} else {
//						errs = append(errs, ajax.Error{
//							Id:      "",
//							Message: err.Error(),
//						})
//					}
//					return
//				}
//			}
//		}
//		// Tạo Adsense Ad Slot
//		errs = append(errs, t.CreateAdsenseAdSlot(record.Id, inputs)...)
//
//	case mysql.TYPEServerTypePrebid: //=>> CASE: LineItemType là PREBID
//		//Tạo bidder info + params cho line item trong bảng line_item_bidder_info nếu server là prebid
//		err = t.CreateBidderInfo(record.Id, inputs)
//		if err != nil {
//			if !utility.IsWindow() {
//				errs = append(errs, ajax.Error{
//					Id:      "",
//					Message: lang.Translate.Errors.LineItemError.BidderInfo.ToString(),
//				})
//			} else {
//				errs = append(errs, ajax.Error{
//					Id:      "",
//					Message: err.Error(),
//				})
//			}
//			return
//		}
//	}
//
//	//Tạo target cho line item vào bảng rl
//	err = t.CreateTarget(record.Id, user.Id, inputs)
//	if err != nil {
//		if !utility.IsWindow() {
//			errs = append(errs, ajax.Error{
//				Id:      "",
//				Message: lang.Translate.Errors.LineItemError.Target.ToString(),
//			})
//		} else {
//			errs = append(errs, ajax.Error{
//				Id:      "",
//				Message: err.Error(),
//			})
//		}
//		return
//	}
//	new(Inventory).UpdateRenderCacheWithLineItem(record.Id, user.Id)
//
//	//Nếu server type là google thì set line vào hàng đợi worker push_line_item_dfp
//	if record.ServerType == mysql.TYPEServerTypeGoogle {
//		new(LineItem).PushToQueueWorkerLineItemDfp(record.Id)
//	}
//
//	return
//}
//
//func (t *LineItemSystem) UpdateLineItem(inputs payload.LineItemSystemAdd, user UserRecord, lang lang.Translation) (record LineItemRecord, errs []ajax.Error) {
//	lineItem := t.GetById(inputs.Id)
//	if lineItem.Id == 0 {
//		errs = append(errs, ajax.Error{
//			Id:      "id",
//			Message: "record not found",
//		})
//		return
//	}
//	// Validate inputs
//	errs = t.Validate(inputs)
//	if len(errs) > 0 {
//		return
//	}
//	// Insert to database
//	var err error
//	inputs.ServerType = lineItem.ServerType
//	record, err = t.makeRecord(inputs)
//	if err != nil {
//		if !utility.IsWindow() {
//			errs = append(errs, ajax.Error{
//				Id:      "",
//				Message: lang.Errors.LineItemError.Edit.ToString(),
//			})
//		} else {
//			errs = append(errs, ajax.Error{
//				Id:      "",
//				Message: err.Error(),
//			})
//		}
//		return
//	}
//	//fmt.Printf("%+v \n", record)
//	err = mysql.Client.Updates(&record).Error
//	if err != nil {
//		if !utility.IsWindow() {
//			errs = append(errs, ajax.Error{
//				Id:      "",
//				Message: lang.Errors.LineItemError.Edit.ToString(),
//			})
//		} else {
//			errs = append(errs, ajax.Error{
//				Id:      "",
//				Message: err.Error(),
//			})
//		}
//		return
//	}
//
//	//Xóa toàn bộ Adsense Ad Slot cũ
//	new(LineItemAdsenseAdSlot).DeleteByLineItem(record.Id)
//
//	// Tạo Adsense Ad Slot
//	errs = append(errs, t.CreateAdsenseAdSlot(record.Id, inputs)...)
//
//	////Xóa toàn bộ account theo line item cũ
//	//new(LineItemAccount).DeleteAccount(record.Id)
//	//
//	////Tạo line_item account nếu server là google
//	//if inputs.ServerType == mysql.TYPEServerTypeGoogle {
//	//	for _, v := range inputs.SelectAccount {
//	//		if v != "" {
//	//			bidderId, _ := strconv.ParseInt(v, 10, 64)
//	//			err = new(LineItemAccount).Create(record.Id, bidderId)
//	//			if err != nil {
//	//				if !utility.IsWindow() {
//	//					errs = append(errs, ajax.Error{
//	//						Id:      "",
//	//						Message: lang.Errors.LineItemError.LineItemAccount.ToString(),
//	//					})
//	//				} else {
//	//					errs = append(errs, ajax.Error{
//	//						Id:      "",
//	//						Message: err.Error(),
//	//					})
//	//				}
//	//				return
//	//			}
//	//		}
//	//	}
//	//}
//
//	//Xóa toàn bộ bidder info + param cũ của line item
//	new(LineItemBidderInfo).DeleteBidderInfo(record.Id)
//	new(LineItemBidderParams).DeleteByLineItemId(record.Id)
//
//	//Tạo bidder info cho line item trong bảng line_item_bidder_info nếu server là prebid
//	if inputs.ServerType == mysql.TYPEServerTypePrebid {
//		err = t.CreateBidderInfo(record.Id, inputs)
//		if err != nil {
//			if !utility.IsWindow() {
//				errs = append(errs, ajax.Error{
//					Id:      "",
//					Message: lang.Errors.LineItemError.BidderInfo.ToString(),
//				})
//			} else {
//				errs = append(errs, ajax.Error{
//					Id:      "",
//					Message: err.Error(),
//				})
//			}
//			return
//		}
//	}
//
//	//Reset cache lại các domain với target cũ
//	new(Inventory).UpdateRenderCacheWithLineItem(record.Id, user.Id)
//
//	//Xóa toàn bộ target cũ để tạo mới list target nhận đc
//	err = new(Target).DeleteTarget(TargetRecord{mysql.TableTarget{
//		LineItemId: record.Id,
//	}})
//	if err != nil {
//		if !utility.IsWindow() {
//			errs = append(errs, ajax.Error{
//				Id:      "",
//				Message: lang.Errors.LineItemError.Target.ToString(),
//			})
//		} else {
//			errs = append(errs, ajax.Error{
//				Id:      "",
//				Message: err.Error(),
//			})
//		}
//		return
//	}
//	//Tạo target cho line item vào bảng rl
//	err = t.CreateTarget(record.Id, user.Id, inputs)
//	if err != nil {
//		if !utility.IsWindow() {
//			errs = append(errs, ajax.Error{
//				Id:      "",
//				Message: lang.Errors.LineItemError.Target.ToString(),
//			})
//		} else {
//			errs = append(errs, ajax.Error{
//				Id:      "",
//				Message: err.Error(),
//			})
//		}
//		return
//	}
//	//Reset cache lại các domain với target mới
//	new(Inventory).UpdateRenderCacheWithLineItem(record.Id, user.Id)
//
//	//Nếu server type là google thì set line vào hàng đợi worker push_line_item_dfp
//	if record.ServerType == 2 {
//		new(LineItem).PushToQueueWorkerLineItemDfp(record.Id)
//	}
//	return
//}
//
//func (t *LineItemSystem) Validate(inputs payload.LineItemSystemAdd) (errs []ajax.Error) {
//	if utility.ValidateString(inputs.Name) == "" {
//		errs = append(errs, ajax.Error{
//			Id:      "name",
//			Message: lang.Translate.ErrorRequired.ToString(),
//		})
//	}
//	if inputs.ServerType == mysql.TYPEServerTypePrebid {
//		if len(inputs.BidderParams) == 0 {
//			errs = append(errs, ajax.Error{
//				Id:      "select-bidder",
//				Message: "Choose at least one bidder",
//			})
//		} else {
//			err := t.ValidateBidder(inputs)
//			errs = append(errs, err...)
//		}
//	} else if inputs.ServerType == mysql.TYPEServerTypeGoogle {
//		if len(inputs.SelectAccount) == 0 {
//			errs = append(errs, ajax.Error{
//				Id:      "select_account",
//				Message: lang.Translate.ErrorRequired.ToString(),
//			})
//		} else {
//			if len(inputs.AdsenseAdSlots) == 0 {
//				errs = append(errs, ajax.Error{
//					Id:      "select-adsense-ad-slot-size",
//					Message: lang.Translate.ErrorRequired.ToString(),
//				})
//			} else {
//				errs = append(errs, t.ValidateAdsenseAdSlot(inputs)...)
//			}
//		}
//		if inputs.LinkedGam == 0 {
//			errs = append(errs, ajax.Error{
//				Id:      "linked_gam",
//				Message: lang.Translate.ErrorRequired.ToString(),
//			})
//		}
//	}
//
//	if inputs.BidderType == 3 {
//		if inputs.BidderRate <= 0 {
//			errs = append(errs, ajax.Error{
//				Id:      "rate",
//				Message: lang.Translate.ErrorRequired.ToString(),
//			})
//		}
//		if utility.ValidateString(inputs.BidderVastUrl) == "" {
//			errs = append(errs, ajax.Error{
//				Id:      "vast_url",
//				Message: lang.Translate.ErrorRequired.ToString(),
//			})
//		}
//	}
//	if inputs.BidderType == 4 {
//		if inputs.BidderRate <= 0 {
//			errs = append(errs, ajax.Error{
//				Id:      "rate",
//				Message: "Rate is required",
//			})
//		}
//		if utility.ValidateString(inputs.BidderAdTag) == "" {
//			errs = append(errs, ajax.Error{
//				Id:      "ad_tag",
//				Message: "Ad Tag is required",
//			})
//		}
//	}
//
//	if len(inputs.ListAdInventory) == 0 {
//		errs = append(errs, ajax.Error{
//			Id:      "text_for_domain",
//			Message: "Target domains is required",
//		})
//	}
//	return
//}
//
//func (t *LineItemSystem) ValidateAdsenseAdSlot(inputs payload.LineItemSystemAdd) (errs []ajax.Error) {
//	for _, adSlot := range inputs.AdsenseAdSlots {
//		if adSlot.AdSlotId == "" {
//			errs = append(errs, ajax.Error{
//				Id:      "adsense-ad-slot-" + adSlot.Size,
//				Message: "Ad Slot ID is required",
//			})
//		}
//	}
//	return
//}
//
//func (t *LineItemSystem) ValidateBidder(inputs payload.LineItemSystemAdd) (errs []ajax.Error) {
//	var listIdTypeClient []int64
//	for _, bidderInfo := range inputs.BidderParams {
//		if bidderInfo.BidderType != mysql.TYPEBidderTypePrebidClient && bidderInfo.BidderType != mysql.TYPEBidderTypePrebidServer {
//			errs = append(errs, ajax.Error{
//				Id:      "",
//				Message: " Bidder " + bidderInfo.BidderName + " type invalid",
//			})
//			return
//		}
//		if bidderInfo.BidderType == 1 {
//			listIdTypeClient = append(listIdTypeClient, bidderInfo.BidderId)
//		}
//		flag := t.Count(bidderInfo.BidderId, listIdTypeClient)
//		if !flag {
//			errs = append(errs, ajax.Error{
//				Id:      "",
//				Message: " Bidder " + bidderInfo.BidderName + " doesn't allow selection of more than 2 bidders with type is Prebid Client",
//			})
//			return
//		}
//		bidderParams := new(BidderParams).GetByBidderId(bidderInfo.BidderId)
//		for _, bidderParam := range bidderParams {
//			idParam := strconv.FormatInt(bidderInfo.BidderId, 10) + "-" + bidderParam.Name + "-" + strconv.Itoa(bidderInfo.BidderIndex)
//			value := bidderInfo.Params[bidderParam.Name]
//			switch bidderParam.Type {
//			case "int":
//				if value != "" {
//					if !govalidator.IsInt(bidderInfo.Params[bidderParam.Name]) {
//						errs = append(errs, ajax.Error{
//							Id:      idParam,
//							Message: "param " + bidderParam.Name + " value is int",
//						})
//					}
//				}
//				break
//			case "float":
//				if value != "" {
//					if !govalidator.IsFloat(bidderInfo.Params[bidderParam.Name]) {
//						errs = append(errs, ajax.Error{
//							Id:      idParam,
//							Message: "param " + bidderParam.Name + " value is float",
//						})
//					}
//				}
//				break
//			case "json":
//				if value != "" {
//					if !govalidator.IsJSON(bidderInfo.Params[bidderParam.Name]) {
//						errs = append(errs, ajax.Error{
//							Id:      idParam,
//							Message: "param " + bidderParam.Name + " value is json",
//						})
//					}
//				}
//			case "boolean":
//				if value == "" {
//					value = "false"
//				} else if value != "true" && value != "false" {
//					errs = append(errs, ajax.Error{
//						Id:      idParam,
//						Message: "param " + bidderParam.Name + " value is true and false",
//					})
//				}
//				break
//			}
//		}
//	}
//	return
//}
//
//func (this *LineItemSystem) Delete(id, userId int64, lang lang.Translation) fiber.Map {
//	err := mysql.Client.Model(&LineItemRecord{}).Delete(&LineItemRecord{}, "id = ? and user_id = ?", id, userId).Error
//	if err != nil {
//		if !utility.IsWindow() {
//			return fiber.Map{
//				"status":  "err",
//				"message": lang.Errors.LineItemError.Delete.ToString(),
//				"id":      id,
//			}
//		}
//		return fiber.Map{
//			"status":  "err",
//			"message": err.Error(),
//			"id":      id,
//		}
//	} else {
//		new(Inventory).UpdateRenderCacheWithLineItem(id, userId)
//		return fiber.Map{
//			"status":  "success",
//			"message": "done",
//			"id":      id,
//		}
//	}
//}
//
//func (t *LineItemSystem) GetByFilters(inputs *payload.LineItemFilterPayload, lang lang.Translation) (response datatable.Response, err error) {
//	var bidders []LineItemRecord
//	var total int64
//	if inputs.PostData.Domain != nil || inputs.PostData.AdFormat != nil || inputs.PostData.AdSize != nil || inputs.PostData.AdTag != nil || inputs.PostData.Device != nil || inputs.PostData.Country != nil {
//		err = mysql.Client.Where("line_item.line_item_type = 2").
//			Scopes(
//				t.SetFilterStatus(inputs),
//				t.setFilterSearch(inputs),
//				t.setFilterType(inputs),
//				t.setFilterDomain(inputs),
//				t.setFilterAdFormat(inputs),
//				t.setFilterAdSize(inputs),
//				t.setFilterAdTag(inputs),
//				t.setFilterDevice(inputs),
//				t.setFilterCountry(inputs),
//			).
//			Joins("inner join target on target.line_item_id = line_item.id").
//			Model(&bidders).Count(&total).
//			Scopes(
//				t.setOrder(inputs),
//				pagination.Paginate(pagination.Params{
//					Limit:  inputs.Length,
//					Offset: inputs.Start,
//				}),
//			).
//			Select("line_item.*").
//			Group("line_item.id").
//			Find(&bidders).Error
//		if err != nil {
//			if !utility.IsWindow() {
//				err = fmt.Errorf(lang.Errors.LineItemError.List.ToString())
//			}
//			return datatable.Response{}, err
//		}
//	} else {
//		err = mysql.Client.Where("line_item_type = 2").
//			Scopes(
//				t.SetFilterStatus(inputs),
//				t.setFilterSearch(inputs),
//				t.setFilterType(inputs),
//				t.setFilterDomain(inputs),
//				t.setFilterAdFormat(inputs),
//				t.setFilterAdSize(inputs),
//				t.setFilterAdTag(inputs),
//				t.setFilterDevice(inputs),
//				t.setFilterCountry(inputs),
//			).
//			Model(&bidders).Count(&total).
//			Scopes(
//				t.setOrder(inputs),
//				pagination.Paginate(pagination.Params{
//					Limit:  inputs.Length,
//					Offset: inputs.Start,
//				}),
//			).
//			Select("line_item.*").
//			Group("line_item.id").
//			Find(&bidders).Error
//		if err != nil {
//			if !utility.IsWindow() {
//				err = fmt.Errorf(lang.Errors.LineItemError.List.ToString())
//			}
//			return datatable.Response{}, err
//		}
//	}
//	response.Draw = inputs.Draw
//	response.RecordsFiltered = total
//	response.RecordsTotal = total
//	response.Data = t.MakeResponseDatatable(bidders)
//	return
//}
//
//func (t *LineItemSystem) MakeResponseDatatable(bidders []LineItemRecord) (records []BidderRecordDatatable) {
//	for _, bidder := range bidders {
//		target := t.getTarget(bidder)
//		stringHtml := htmlblock.Render("line-item/index/block.target.gohtml", target).String()
//		target.HtmlContent = stringHtml
//
//		//fmt.Println(target.StringInventory)
//		rec := BidderRecordDatatable{
//			LineItemRecord: bidder,
//			RowId:          "bidder_" + strconv.FormatInt(bidder.Id, 10),
//			Name:           htmlblock.Render("_backend/line_item_system/index/block.name.gohtml", target).String(),
//			Target:         htmlblock.Render("_backend/line_item_system/index/block.target-button.gohtml", target).String(),
//			Status:         htmlblock.Render("_backend/line_item_system/index/block.status.gohtml", bidder).String(),
//			Type:           htmlblock.Render("_backend/line_item_system/index/block.type.gohtml", bidder).String(),
//			Rate:           htmlblock.Render("_backend/line_item_system/index/block.rate.gohtml", bidder).String(),
//			Action:         htmlblock.Render("_backend/line_item_system/index/block.action.gohtml", bidder).String(),
//		}
//		records = append(records, rec)
//	}
//	return
//}
//
//func (t *LineItemSystem) SetFilterStatus(inputs *payload.LineItemFilterPayload) func(db *gorm.DB) *gorm.DB {
//	return func(db *gorm.DB) *gorm.DB {
//		if inputs.PostData.Status != nil {
//			switch inputs.PostData.Status.(type) {
//			case string, int:
//				if inputs.PostData.Status != "" {
//					return db.Where("status = ?", inputs.PostData.Status)
//				}
//			case []string, []interface{}:
//				return db.Where("status IN ?", inputs.PostData.Status)
//			}
//		}
//		return db
//	}
//}
//
//func (t *LineItemSystem) setFilterType(inputs *payload.LineItemFilterPayload) func(db *gorm.DB) *gorm.DB {
//	return func(db *gorm.DB) *gorm.DB {
//		if inputs.PostData.Type != nil {
//			switch inputs.PostData.Type.(type) {
//			case string, int:
//				if inputs.PostData.Type != "" {
//					return db.Where("type = ?", inputs.PostData.Type)
//				}
//			case []string, []interface{}:
//				return db.Where("type IN ?", inputs.PostData.Type)
//			}
//		}
//		return db
//	}
//}
//
//func (t *LineItemSystem) setFilterSearch(inputs *payload.LineItemFilterPayload) func(db *gorm.DB) *gorm.DB {
//	return func(db *gorm.DB) *gorm.DB {
//		var flag bool
//		// Search from form of datatable <- not use
//		if inputs.Search != nil && inputs.Search.Value != "" {
//			flag = true
//		}
//		// Search from form filter
//		if inputs.PostData.QuerySearch != "" {
//			flag = true
//		}
//		if !flag {
//			return db
//		}
//		return db.Where("name LIKE ?", "%"+inputs.PostData.QuerySearch+"%")
//	}
//}
//
//func (t *LineItemSystem) setOrder(inputs *payload.LineItemFilterPayload) func(db *gorm.DB) *gorm.DB {
//	return func(db *gorm.DB) *gorm.DB {
//		if len(inputs.Order) > 0 {
//			var orders []string
//			for _, order := range inputs.Order {
//				column := inputs.Columns[order.Column]
//				orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
//			}
//			var orderString string
//			if inputs.PostData.Domain != nil {
//				orderString = strings.Join(orders, ", ")
//				orderString = "line_item." + orderString
//			} else {
//				orderString = strings.Join(orders, ", ")
//			}
//			return db.Order(orderString)
//		}
//		return db
//	}
//}
//
//func (t *LineItemSystem) setFilterDomain(inputs *payload.LineItemFilterPayload) func(db *gorm.DB) *gorm.DB {
//	return func(db *gorm.DB) *gorm.DB {
//		if inputs.PostData.Domain != nil {
//			switch inputs.PostData.Domain.(type) {
//			case string, int:
//				if inputs.PostData.Domain != "" {
//					return db.Where("target.inventory_id = ? or target.inventory_id = -1", inputs.PostData.Domain).Group("line_item.id")
//				}
//			case []string, []interface{}:
//				return db.Where("target.inventory_id IN ? or target.inventory_id = -1", inputs.PostData.Domain).Group("line_item.id")
//			}
//		}
//		return db
//	}
//}
//
//func (t *LineItemSystem) setFilterAdFormat(inputs *payload.LineItemFilterPayload) func(db *gorm.DB) *gorm.DB {
//	return func(db *gorm.DB) *gorm.DB {
//		if inputs.PostData.AdFormat != nil {
//			switch inputs.PostData.AdFormat.(type) {
//			case string, int:
//				if inputs.PostData.AdFormat != "" {
//					return db.Where("target.ad_format_id = ? or target.ad_format_id = -1", inputs.PostData.AdFormat).Group("line_item.id")
//				}
//			case []string, []interface{}:
//				return db.Where("target.ad_format_id IN ? or target.ad_format_id = -1", inputs.PostData.AdFormat).Group("line_item.id")
//			}
//		}
//		return db
//	}
//}
//
//func (t *LineItemSystem) setFilterAdSize(inputs *payload.LineItemFilterPayload) func(db *gorm.DB) *gorm.DB {
//	return func(db *gorm.DB) *gorm.DB {
//		if inputs.PostData.AdSize != nil {
//			switch inputs.PostData.AdSize.(type) {
//			case string, int:
//				if inputs.PostData.AdSize != "" {
//					return db.Where("target.ad_size_id = ? or target.ad_size_id = -1", inputs.PostData.AdSize).Group("line_item.id")
//				}
//			case []string, []interface{}:
//				return db.Where("target.ad_size_id IN ? or target.ad_size_id = -1", inputs.PostData.AdSize).Group("line_item.id")
//			}
//		}
//		return db
//	}
//}
//
//func (t *LineItemSystem) setFilterAdTag(inputs *payload.LineItemFilterPayload) func(db *gorm.DB) *gorm.DB {
//	return func(db *gorm.DB) *gorm.DB {
//		if inputs.PostData.AdTag != nil {
//			switch inputs.PostData.AdTag.(type) {
//			case string, int:
//				if inputs.PostData.AdTag != "" {
//					return db.Where("target.tag_id = ? or target.tag_id = -1", inputs.PostData.AdTag).Group("line_item.id")
//				}
//			case []string, []interface{}:
//				return db.Where("target.tag_id IN ? or target.tag_id = -1", inputs.PostData.AdTag).Group("line_item.id")
//			}
//		}
//		return db
//	}
//}
//
//func (t *LineItemSystem) setFilterDevice(inputs *payload.LineItemFilterPayload) func(db *gorm.DB) *gorm.DB {
//	return func(db *gorm.DB) *gorm.DB {
//		if inputs.PostData.Device != nil {
//			switch inputs.PostData.Device.(type) {
//			case string, int:
//				if inputs.PostData.Device != "" {
//					return db.Where("target.device_id = ? or target.device_id = -1", inputs.PostData.Device).Group("line_item.id")
//				}
//			case []string, []interface{}:
//				return db.Where("target.device_id IN ? or target.device_id = -1", inputs.PostData.Device).Group("line_item.id")
//			}
//		}
//		return db
//	}
//}
//
//func (t *LineItemSystem) setFilterCountry(inputs *payload.LineItemFilterPayload) func(db *gorm.DB) *gorm.DB {
//	return func(db *gorm.DB) *gorm.DB {
//		if inputs.PostData.Country != nil {
//			switch inputs.PostData.Country.(type) {
//			case string, int:
//				if inputs.PostData.Country != "" {
//					return db.Where("target.geo_id = ? or target.geo_id = -1", inputs.PostData.Country).Group("line_item.id")
//				}
//			case []string, []interface{}:
//				return db.Where("target.geo_id IN ? or target.geo_id = -1", inputs.PostData.Country).Group("line_item.id")
//			}
//		}
//		return db
//	}
//}
//
//func (t *LineItemSystem) getTarget(lineItem LineItemRecord) Targets {
//	assignTarget := Targets{
//		LineItem: lineItem,
//	}
//	var listNameInventoryTarget []string
//	var listNameAdFormatTarget []string
//	var listNameAdSizeTarget []string
//	var listNameAdTagTarget []string
//	var listNameGeoTarget []string
//	var listNameDeviceTarget []string
//
//	targets := new(Target).GetTargetLineItem(lineItem.Id)
//
//	mapInventory := make(map[int64]int)
//	mapAdFormat := make(map[int64]int)
//	mapAdSize := make(map[int64]int)
//	mapAdTag := make(map[int64]int)
//	mapGeo := make(map[int64]int)
//	mapDevice := make(map[int64]int)
//	for _, target := range targets {
//		if target.InventoryId != 0 {
//			mapInventory[target.InventoryId] = 1
//		}
//		if target.AdFormatId != 0 {
//			mapAdFormat[target.AdFormatId] = 1
//		}
//		if target.AdSizeId != 0 {
//			mapAdSize[target.AdSizeId] = 1
//		}
//		if target.TagId != 0 {
//			mapAdTag[target.TagId] = 1
//		}
//		if target.GeoId != 0 {
//			mapGeo[target.GeoId] = 1
//		}
//		if target.DeviceId != 0 {
//			mapDevice[target.DeviceId] = 1
//		}
//	}
//
//	// Lọc bỏ những id trùng nhau
//	for inventoryId, _ := range mapInventory {
//		inventory, _ := new(Inventory).GetByIdSystem(inventoryId)
//		listNameInventoryTarget = append(listNameInventoryTarget, inventory.Name)
//	}
//	for adFormatId, _ := range mapAdFormat {
//		adFormat := new(AdType).GetById(adFormatId)
//		listNameAdFormatTarget = append(listNameAdFormatTarget, adFormat.Name)
//	}
//	for adSizeId, _ := range mapAdSize {
//		adSize := new(AdSize).GetById(adSizeId)
//		listNameAdSizeTarget = append(listNameAdSizeTarget, adSize.Name)
//	}
//	for adTagId, _ := range mapAdTag {
//		adTag := new(InventoryAdTag).GetById(adTagId)
//		listNameAdTagTarget = append(listNameAdTagTarget, adTag.Name)
//	}
//	for geoId, _ := range mapGeo {
//		geo := new(Country).GetById(geoId)
//		listNameGeoTarget = append(listNameGeoTarget, geo.Name)
//	}
//	for deviceId, _ := range mapDevice {
//		device := new(Device).GetById(deviceId)
//		listNameDeviceTarget = append(listNameDeviceTarget, device.Name)
//	}
//
//	assignTarget.TextListInventory = strings.Join(listNameInventoryTarget, ", ")
//	assignTarget.TextListAdTag = strings.Join(listNameAdTagTarget, ", ")
//	assignTarget.TextListAdSize = strings.Join(listNameAdSizeTarget, ", ")
//	assignTarget.TextListAdFormat = strings.Join(listNameAdFormatTarget, ", ")
//	assignTarget.TextListGeo = strings.Join(listNameGeoTarget, ", ")
//	assignTarget.TextListDevice = strings.Join(listNameDeviceTarget, ", ")
//	return assignTarget
//}
//
//func (t *LineItemSystem) CreateAdsenseAdSlot(lineItemId int64, inputs payload.LineItemSystemAdd) (errs []ajax.Error) {
//	for _, adsenseAdSlot := range inputs.AdsenseAdSlots {
//		if err := new(LineItemAdsenseAdSlot).Push(lineItemId, adsenseAdSlot.Size, adsenseAdSlot.AdSlotId); err != nil {
//			errs = append(errs, ajax.Error{
//				Id:      `adsense-ad-slot-` + adsenseAdSlot.Size,
//				Message: err.Error(),
//			})
//		}
//	}
//	return
//}
//
//func (t *LineItemSystem) CreateBidderInfo(lineItemId int64, inputs payload.LineItemSystemAdd) (err error) {
//	for _, bidderInfo := range inputs.BidderParams {
//		// Bỏ qua các bidder không có param
//		if len(bidderInfo.Params) < 1 {
//			continue
//		}
//		// Tạo bảng bidder info trước rồi lấy id của bảng này để tạo rls
//		recBidderInfo := new(LineItemBidderInfo).CreateBidderInfo(LineItemBidderInfoRecord{mysql.TableLineItemBidderInfo{
//			LineItemId: lineItemId,
//			BidderId:   bidderInfo.BidderId,
//			Name:       bidderInfo.BidderName,
//			BidderType: bidderInfo.BidderType,
//		}})
//		// Lấy toàn bộ các param theo id của bidder
//		bidderParams := new(BidderParams).GetByBidderId(bidderInfo.BidderId)
//		for _, bidderParam := range bidderParams {
//			var value string
//			// Kiểm tra xem param có nằm trong param add lên không nếu không thì bỏ qua nếu có thì lấy value để tạo bảng row line_item_bidder_params
//			if v, ok := bidderInfo.Params[bidderParam.Name]; ok {
//				value = v
//			} else {
//				continue
//			}
//			err = mysql.Client.Create(&LineItemBidderParamsRecord{mysql.TableLineItemBidderParams{
//				LineItemId:       lineItemId,
//				BidderId:         bidderParam.BidderId,
//				LineItemBidderId: recBidderInfo.Id,
//				Name:             bidderParam.Name,
//				Type:             bidderParam.Type,
//				Value:            value,
//			}}).Error
//		}
//	}
//	return
//}
//
//func (t *LineItemSystem) CreateTarget(lineItemId, userId int64, inputs payload.LineItemSystemAdd) (err error) {
//	all := int64(-1)
//	// Kiểm tra nếu đầu vào input list target = 0 thì thêm một target = 0 thể hiện select all
//	if len(inputs.ListAdInventory) == 0 {
//		inputs.ListAdInventory = []payload.ListTarget{
//			{
//				Id: all,
//			},
//		}
//	}
//	if len(inputs.ListAdFormat) == 0 {
//		inputs.ListAdFormat = []payload.ListTarget{
//			{
//				Id: all,
//			},
//		}
//	}
//	if len(inputs.ListGeo) == 0 {
//		inputs.ListGeo = []payload.ListTarget{
//			{
//				Id: all,
//			},
//		}
//	}
//	if len(inputs.ListAdSize) == 0 {
//		inputs.ListAdSize = []payload.ListTarget{
//			{
//				Id: all,
//			},
//		}
//	}
//	if len(inputs.ListAdTag) == 0 {
//		inputs.ListAdTag = []payload.ListTarget{
//			{
//				Id: all,
//			},
//		}
//	}
//	if len(inputs.ListDevice) == 0 {
//		inputs.ListDevice = []payload.ListTarget{
//			{
//				Id: all,
//			},
//		}
//	}
//	//If type là gg thì mặc định size và device là all
//	if inputs.ServerType == mysql.TYPEServerTypeGoogle {
//		inputs.ListAdSize = []payload.ListTarget{
//			{
//				Id: all,
//			},
//		}
//		inputs.ListDevice = []payload.ListTarget{
//			{
//				Id: all,
//			},
//		}
//	}
//	for _, inventory := range inputs.ListAdInventory {
//		var recordTarget TargetRecord
//		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
//			UserId:      userId,
//			LineItemId:  lineItemId,
//			InventoryId: inventory.Id,
//		}}).Error
//		if err != nil {
//			return err
//		}
//	}
//	for _, adFormat := range inputs.ListAdFormat {
//		var recordTarget TargetRecord
//		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
//			UserId:     userId,
//			LineItemId: lineItemId,
//			AdFormatId: adFormat.Id,
//		}}).Error
//		if err != nil {
//			return err
//		}
//	}
//	for _, size := range inputs.ListAdSize {
//		var recordTarget TargetRecord
//		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
//			UserId:     userId,
//			LineItemId: lineItemId,
//			AdSizeId:   size.Id,
//		}}).Error
//		if err != nil {
//			return err
//		}
//	}
//
//	for _, geo := range inputs.ListGeo {
//		var recordTarget TargetRecord
//		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
//			UserId:     userId,
//			LineItemId: lineItemId,
//			GeoId:      geo.Id,
//		}}).Error
//		if err != nil {
//			return err
//		}
//	}
//	for _, device := range inputs.ListDevice {
//		var recordTarget TargetRecord
//		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
//			UserId:     userId,
//			LineItemId: lineItemId,
//			DeviceId:   device.Id,
//		}}).Error
//		if err != nil {
//			return err
//		}
//	}
//
//	for _, adTag := range inputs.ListAdTag {
//		var recordTarget TargetRecord
//		err := mysql.Client.Table("target").FirstOrCreate(&recordTarget, TargetRecord{mysql.TableTarget{
//			UserId:     userId,
//			LineItemId: lineItemId,
//			TagId:      adTag.Id,
//		}}).Error
//		if err != nil {
//			return err
//		}
//	}
//
//	return
//}
//
//func (t *LineItemSystem) GetInventoryByName(search string) (rows []SearchLineItem) {
//	q := "%" + search + "%"
//	mysql.Client.Raw("SELECT id,name FROM `inventory` WHERE name like ? limit 20", q).Find(&rows)
//	return
//}
//
//func (t *LineItemSystem) GetAdFormatByName(search string) (rows []SearchLineItem) {
//	q := "%" + search + "%"
//	mysql.Client.Raw("SELECT id,name FROM `ad_type` WHERE name like ? limit 20", q).Find(&rows)
//	return
//}
//
//func (t *LineItemSystem) GetAdSizeByName(search string) (rows []SearchLineItem) {
//	q := "%" + search + "%"
//	mysql.Client.Raw("SELECT id,name FROM `ad_size` WHERE name like ? limit 10", q).Find(&rows)
//	return
//}
//
//func (t *LineItemSystem) GetAdTagByName(search string) (rows []SearchLineItem) {
//	q := "%" + search + "%"
//	mysql.Client.Raw("SELECT id,name FROM `inventory_ad_tag` WHERE name like ? limit 10", q).Find(&rows)
//	return
//}
//
//func (t *LineItemSystem) GetDeviceByName(search string) (rows []SearchLineItem) {
//	q := "%" + search + "%"
//	mysql.Client.Raw("SELECT id,name FROM `devices` WHERE name like ? limit 10", q).Find(&rows)
//	return
//}
//
//func (t *LineItemSystem) GetCountryByName(search string) (rows []SearchLineItem) {
//	q := "%" + search + "%"
//	mysql.Client.Raw("SELECT id,name FROM `country` WHERE name like ? limit 10", q).Find(&rows)
//	return
//}
//
//func (t *LineItemSystem) GetSearched(typ string, search []string) (rows []SearchLineItem) {
//	switch typ {
//	case "domain":
//		mysql.Client.Raw("SELECT id,name FROM `inventory` WHERE id in ?", search).Find(&rows)
//		break
//	case "format":
//		mysql.Client.Raw("SELECT id,name FROM `ad_type` WHERE id in ?", search).Find(&rows)
//		break
//	case "size":
//		mysql.Client.Raw("SELECT id,name FROM `ad_size` WHERE id in ?", search).Find(&rows)
//		break
//	case "adtag":
//		mysql.Client.Raw("SELECT id,name FROM `inventory_ad_tag` WHERE id in ?", search).Find(&rows)
//		break
//	case "country":
//		mysql.Client.Raw("SELECT id,name FROM `country` WHERE id in ?", search).Find(&rows)
//		break
//	case "device":
//		mysql.Client.Raw("SELECT id,name FROM `devices` WHERE id in ?", search).Find(&rows)
//		break
//	}
//	return
//}
//
//func (t *LineItemSystem) GetOfUserById(id, userId int64) (row LineItemRecord) {
//	mysql.Client.Model(&LineItemRecord{}).Where("id = ? and user_id = ?", id, userId).Find(&row)
//	return
//}
//
//func (t *LineItemSystem) InArray(array []SearchLineItem, id int64) (index int, flag bool) {
//	for i, v := range array {
//		if v.Id == id {
//			return i, true
//		}
//	}
//	return 0, false
//}
//
//func (t *LineItemSystem) GetFilter(option string, searched []string) (rows []SearchLineItem) {
//	switch option {
//	case "domain":
//		rows = t.GetInventoryByName("")
//		//fmt.Printf("%+v\n", rows)
//		search := t.GetSearched(option, searched)
//		for i, v := range search {
//			index, flag := t.InArray(rows, v.Id)
//			if flag {
//				rows[index].Selected = true
//			} else {
//				search[i].Selected = true
//				rows = append(rows, search[i])
//			}
//		}
//		break
//	case "format":
//		rows = t.GetAdFormatByName("")
//		search := t.GetSearched(option, searched)
//		for i, v := range search {
//			index, flag := t.InArray(rows, v.Id)
//			if flag {
//				rows[index].Selected = true
//			} else {
//				search[i].Selected = true
//				rows = append(rows, search[i])
//			}
//		}
//		break
//	case "size":
//		rows = t.GetAdSizeByName("")
//		search := t.GetSearched(option, searched)
//		for i, v := range search {
//			index, flag := t.InArray(rows, v.Id)
//			if flag {
//				rows[index].Selected = true
//			} else {
//				search[i].Selected = true
//				rows = append(rows, search[i])
//			}
//		}
//		break
//	case "adtag":
//		rows = t.GetAdTagByName("")
//		search := t.GetSearched(option, searched)
//		for i, v := range search {
//			index, flag := t.InArray(rows, v.Id)
//			if flag {
//				rows[index].Selected = true
//			} else {
//				search[i].Selected = true
//				rows = append(rows, search[i])
//			}
//		}
//		break
//	case "country":
//		rows = t.GetCountryByName("")
//		search := t.GetSearched(option, searched)
//		for i, v := range search {
//			index, flag := t.InArray(rows, v.Id)
//			if flag {
//				rows[index].Selected = true
//			} else {
//				search[i].Selected = true
//				rows = append(rows, search[i])
//			}
//		}
//		break
//	case "device":
//		rows = t.GetDeviceByName("")
//		search := t.GetSearched(option, searched)
//		for i, v := range search {
//			index, flag := t.InArray(rows, v.Id)
//			if flag {
//				rows[index].Selected = true
//			} else {
//				search[i].Selected = true
//				rows = append(rows, search[i])
//			}
//		}
//		break
//	}
//	return
//}
//
//func (t *LineItemSystem) GetListBoxCollapse(userId, lineItemId int64, page, typ string) (list []string) {
//	switch typ {
//	case "add":
//		mysql.Client.Select("box_collapse").Model(PageCollapseRecord{}).Where("user_id = ? and page_collapse = ? and is_collapse = ? and page_type = ?", userId, page, 1, typ).Find(&list)
//		return
//	case "edit":
//		mysql.Client.Select("box_collapse").Model(PageCollapseRecord{}).Where("user_id = ? and page_collapse = ? and is_collapse = ? and page_type = ? and page_id = ?", userId, page, 1, typ, lineItemId).Find(&list)
//		return
//	}
//	return
//}
//
//func (t *LineItemSystem) Count(id int64, listId []int64) bool {
//	var count int
//	for _, v := range listId {
//		if v == id {
//			count++
//			if count > 1 {
//				return false
//			}
//		}
//	}
//	return true
//}
//
//func (t *LineItemSystem) PushToQueueWorkerLineItemDfp(lineItemId int64) {
//	//Push to queue worker với lineId và type là google
//	mysql.Client.Model(LineItemRecord{}).Where("id = ? and server_type = 2 and push_line_item_dfp != 4", lineItemId).Update("push_line_item_dfp", 1)
//	return
//}
//
//func (t *LineItemSystem) GetById(id int64) (record LineItemRecord) {
//	mysql.Client.Where("id = ? and line_item_type = 2", id).Find(&record)
//	return
//}