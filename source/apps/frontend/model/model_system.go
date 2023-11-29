package model

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"math"
	"regexp"
	"sort"
	"source/apps/frontend/lang"
	"source/apps/frontend/payload"
	"source/apps/history"
	"source/core/technology/mysql"
	"source/pkg/adstxt"
	"source/pkg/ajax"
	"source/pkg/cloudflare"
	"source/pkg/datatable"
	"source/pkg/htmlblock"
	"source/pkg/pagination"
	"source/pkg/utility"
	"strconv"
	"strings"
	"time"
)

type System struct{}

//func (t *System) AutoCreateApacdex(user UserRecord, adsTxt string) (record BidderRecord, errs []ajax.Error) {
//	adsTxt = "quantumdex.io, PP" + strconv.FormatInt(user.Id, 10) + ", DIRECT\ninterdogmedia.com, PP" + strconv.FormatInt(user.Id, 10) + ", DIRECT\napacdex.com, PP" + strconv.FormatInt(user.Id, 10) + ", DIRECT\nappnexus.com, 10273, RESELLER, f5ab79cb980f11d1\nappnexus.com, 11395, RESELLER, f5ab79cb980f11d1\nadvertising.com, 28643, RESELLER #VerizonVideo\nyahoo.com, 58754, RESELLER, e1a5b5b6e3255540\naol.com, 58754, RESELLER, e1a5b5b6e3255540\npubmatic.com, 157940, RESELLER, 5d62403b186f2ace\nonetag.com, 2bb78272a859ca6, RESELLER\nsharethrough.com, cc26d15a, RESELLER, d53b998a7bd4ecd2"
//	lang := lang.Translation{}
//	inputs := payload.SystemCreate{
//		AccountType:   1,
//		BidderId:      95,
//		MediaType:     []string{"1", "2"},
//		BidAdjustment: 1,
//		AdsTxt:        adsTxt,
//		IsLock:        mysql.TYPEIsLockTypeLock,
//		IsDefault:     mysql.TypeOn,
//	}
//	record, errs = t.Create(inputs, user, lang)
//	return
//}

func (t *System) Create(inputs payload.SystemCreate, userLogin UserRecord, userAdmin UserRecord, lang lang.Translation) (record BidderRecord, errs []ajax.Error) {
	errs = t.ValidateCreate(inputs, userLogin, lang)
	if len(errs) > 0 {
		return
	}
	record.makeRow(inputs, userLogin)
	err := mysql.Client.Create(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.BidderError.Add.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}
	//Bidder template id = 1 và 2 mặc định là Google và amazon được tự thêm mediatype theo user post lên
	if record.BidderTemplateId == 1 || record.BidderTemplateId == 2 {
		for _, v := range inputs.MediaType {
			mediaTypeId, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				err = new(RlBidderMediaType).Create(RlsBidderMediaTypeRecord{mysql.TableRlsBidderMediaType{
					BidderId:    record.Id,
					MediaTypeId: mediaTypeId,
				}})
				if err != nil {
					if !utility.IsWindow() {
						errs = append(errs, ajax.Error{
							Id:      "",
							Message: lang.Errors.BidderError.MediaBidder.ToString(),
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
		// Nếu bidder = 2 là amazon thì add param mẫu
		if inputs.BidderId == 2 {
			// Tạo param cho bidder
			bidderTemplate := new(BidderTemplate).GetById(inputs.BidderId)
			record.BidderCode = bidderTemplate.BidderCode
			// Get param của bidder template
			bidderTemplateParams := new(BidderTemplateParams).GetByBidderTemplateId(inputs.BidderId)
			//Tạo các param với bidder template tạo từ backend
			mapTemplateParam := make(map[string]string)
			for _, bidderTemplateParam := range bidderTemplateParams {
				//Lưu lại nam và type của bidder template param vào map để check với param inputs
				mapTemplateParam[bidderTemplateParam.Name] = bidderTemplateParam.Type
				mysql.Client.Create(&BidderParamsRecord{mysql.TableBidderParams{
					UserId:           userLogin.Id,
					BidderId:         record.Id,
					BidderTemplateId: bidderTemplate.Id,
					Name:             bidderTemplateParam.Name,
					Type:             bidderTemplateParam.Type,
					Template:         bidderTemplateParam.Template,
				}})
			}
		}
	} else {
		//Lấy list media type của bidder created by admin add vào bidder của pub tạo
		mediaTypeTemplates := new(RlBidderMediaType).GetByBidderTemplateId(record.BidderTemplateId)
		for _, v := range mediaTypeTemplates {
			rlsMediaType := RlsBidderMediaTypeRecord{mysql.TableRlsBidderMediaType{
				BidderId:    record.Id,
				MediaTypeId: v.MediaTypeId,
			}}
			err = new(RlBidderMediaType).Create(rlsMediaType)
			if err != nil {
				if !utility.IsWindow() {
					errs = append(errs, ajax.Error{
						Id:      "",
						Message: lang.Errors.BidderError.MediaBidder.ToString(),
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

		// Tạo param cho bidder
		bidderTemplate := new(BidderTemplate).GetById(inputs.BidderId)
		record.BidderCode = bidderTemplate.BidderCode
		// Get param của bidder template
		bidderTemplateParams := new(BidderTemplateParams).GetByBidderTemplateId(inputs.BidderId)
		//Tạo các param với bidder template tạo từ backend
		mapTemplateParam := make(map[string]string)
		for _, bidderTemplateParam := range bidderTemplateParams {
			//Lưu lại nam và type của bidder template param vào map để check với param inputs
			mapTemplateParam[bidderTemplateParam.Name] = bidderTemplateParam.Type
			mysql.Client.Create(&BidderParamsRecord{mysql.TableBidderParams{
				UserId:           userLogin.Id,
				BidderId:         record.Id,
				BidderTemplateId: bidderTemplate.Id,
				Name:             bidderTemplateParam.Name,
				Type:             bidderTemplateParam.Type,
				Template:         bidderTemplateParam.Template,
			}})
		}
		//Tạo các param còn lại ngoài các param mẫu do người dùng add thêm
		if len(inputs.Params) > 0 {
			for _, param := range inputs.Params {
				// Kiểm tra nếu param mới truyền vào thuộc param template thì bỏ qua
				if _, ok := mapTemplateParam[param.Name]; ok {
					continue
				}
				//fmt.Println(param)
				mysql.Client.Create(&BidderParamsRecord{mysql.TableBidderParams{
					UserId:   userLogin.Id,
					BidderId: record.Id,
					Name:     param.Name,
					Type:     param.Type,
				}})
			}
		}
	}
	// Trong trường hợp bidder amazon add và xlsx path khác rỗng thì parse file và log lại cpm
	if record.BidderTemplateId == 2 && record.XlsxPath != "" {
		_ = t.LogXlsxForAmz(record.Id, record.XlsxPath)
	}
	// Xử lý render cache lại toàn bộ domain
	new(Inventory).ResetCacheAll(userLogin.Id)

	// Reset link ads
	var listUuid []string
	inventories, _ := new(Inventory).GetByUserId(record.UserId)
	for _, inventory := range inventories {
		listUuid = append(listUuid, inventory.Uuid)
	}
	err = cloudflare.ResetLinkAds(listUuid)
	if err != nil {
		fmt.Println(err)
	}

	record.GetById(record.Id)
	// Log history
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userLogin.Id
	}
	_ = history.PushHistory(&history.Bidder{
		Detail:    history.DetailBidderFE,
		CreatorId: creatorId,
		RecordNew: record.TableBidder,
	})
	return
}

func (t *System) UpdateBidder(inputs payload.SystemCreate, userLogin UserRecord, userAdmin UserRecord, lang lang.Translation) (recordNew BidderRecord, errs []ajax.Error) {
	// Kiểm tra record xem có tồn tại không
	recordOld, _ := t.VerificationRecord(inputs.Id, userLogin.Id)
	if recordOld.Id < 1 {
		errs = append(errs, ajax.Error{
			Id:      "id",
			Message: "You don't own this bidder",
		})
		return
	}
	if recordOld.IsLock == mysql.TYPEIsLockTypeLock {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: "You cannot edit this bidder",
		})
		return
	}
	//Get all data
	recordOld.GetById(recordOld.Id)
	recordNew = recordOld
	// Validate inputs
	errs = t.ValidateUpdate(inputs, userLogin, lang)
	if len(errs) > 0 {
		return
	}
	// Insert to database
	var err error
	recordNew.makeRow(inputs, userLogin)
	err = mysql.Client.Save(&recordNew).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.BidderError.Edit.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}

	//Bidder template id = 1 là Google được tự thêm mediatye theo user post lên
	if recordNew.BidderTemplateId == 1 {
		// Chỉ xử lý media type với bidder google
		mysql.Client.Where("bidder_id = ?", recordNew.Id).Delete(RlsBidderMediaTypeRecord{})
		for _, v := range inputs.MediaType {
			mediaTypeId, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				err = new(RlBidderMediaType).Create(RlsBidderMediaTypeRecord{mysql.TableRlsBidderMediaType{
					BidderId:    recordNew.Id,
					MediaTypeId: mediaTypeId,
				}})
				if err != nil {
					if !utility.IsWindow() {
						errs = append(errs, ajax.Error{
							Id:      "",
							Message: lang.Errors.BidderError.MediaBidder.ToString(),
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
	} else if inputs.BidderId == 2 {
		// Nếu bidder = 2 là amazon thì giữ nguyên param mẫu và không đổi lại mediatype đã add

	} else { // Ngược lại nếu không phải bidder google và amz thì mặc định lấy lại media type đã đặt sẵn của bidder template

		//Xóa toàn bộ rls media type + param cũ của bidder
		mysql.Client.Where("bidder_id = ?", recordNew.Id).Delete(BidderParamsRecord{})

		// Tạo param cho bidder
		bidderTemplate := new(BidderTemplate).GetById(inputs.BidderId)
		recordNew.BidderCode = bidderTemplate.BidderCode
		// Get param của bidder template
		bidderTemplateParams := new(BidderTemplateParams).GetByBidderTemplateId(inputs.BidderId)
		//Tạo các param với bidder template tạo từ backend
		mapTemplateParam := make(map[string]string)
		for _, bidderTemplateParam := range bidderTemplateParams {
			//Lưu lại nam và type của bidder template param vào map để check với param inputs
			mapTemplateParam[bidderTemplateParam.Name] = bidderTemplateParam.Type
			mysql.Client.Create(&BidderParamsRecord{mysql.TableBidderParams{
				UserId:           userLogin.Id,
				BidderId:         recordNew.Id,
				BidderTemplateId: bidderTemplate.Id,
				Name:             bidderTemplateParam.Name,
				Type:             bidderTemplateParam.Type,
				Template:         bidderTemplateParam.Template,
			}})
		}
		//Tạo các param còn lại ngoài các param mẫu do người dùng add thêm
		if len(inputs.Params) > 0 {
			for _, param := range inputs.Params {
				// Kiểm tra nếu param mới truyền vào thuộc param template thì bỏ qua
				if _, ok := mapTemplateParam[param.Name]; ok {
					continue
				}
				mysql.Client.Create(&BidderParamsRecord{mysql.TableBidderParams{
					UserId:   userLogin.Id,
					BidderId: recordNew.Id,
					Name:     param.Name,
					Type:     param.Type,
				}})
			}
		}
	}

	// Trong trường hợp bidder amazon add và xlsx path khác rỗng thì parse file và log lại cpm
	if recordNew.BidderTemplateId == 2 && recordNew.XlsxPath != recordOld.XlsxPath {
		// Xóa toàn bộ các log cũ của bidder
		new(LogCpmAmz).DeleteByBidderId(recordNew.Id)
		// Log lại cpm từ file xlsx mới
		_ = t.LogXlsxForAmz(recordNew.Id, recordNew.XlsxPath)
	}

	//Reset cache
	new(Inventory).ResetCacheAll(userLogin.Id)

	// Reset link ads
	if recordOld.AdsTxt != recordNew.AdsTxt {
		var listUuid []string
		inventories, _ := new(Inventory).GetByUserId(recordNew.UserId)
		for _, inventory := range inventories {
			listUuid = append(listUuid, inventory.Uuid)
		}
		err = cloudflare.ResetLinkAds(listUuid)
		if err != nil {
			fmt.Println(err)
		}
	}
	//Push History
	recordNew.GetById(recordNew.Id)
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userLogin.Id
	}
	_ = history.PushHistory(&history.Bidder{
		Detail:    history.DetailBidderFE,
		CreatorId: creatorId,
		RecordOld: recordOld.TableBidder,
		RecordNew: recordNew.TableBidder,
	})
	return
}

func (t *System) SaveConfig(inputs payload.SystemConfig, user UserRecord, lang lang.Translation) (record ConfigRecord, errs []ajax.Error) {
	recordOld, _ := t.VerificationRecordConfig(user.Id)
	if recordOld.Id == 0 {
		errs = append(errs, ajax.Error{
			Id:      "id",
			Message: lang.Errors.ConfigError.Save.ToString(),
		})
		return
	}
	// Validate inputs
	errs = t.ValidateConfig(inputs, recordOld)
	if len(errs) > 0 {
		return
	}

	// Update to database
	var err error
	record = t.makeRowConfig(inputs)
	if !govalidator.IsNull(recordOld.Currency) {
		record.Currency = recordOld.Currency
		record.GranularityMultiplier = recordOld.GranularityMultiplier
	} else {
		currency := new(Currency).GetById(inputs.Currency)
		record.Currency = currency.Code
		record.GranularityMultiplier = currency.GranularityMultiplier
	}
	//fmt.Printf("%+v \n", record)
	err = mysql.Client.Where(ConfigRecord{mysql.TableConfig{UserId: user.Id}}).Updates(&record).Error
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}

	new(Inventory).ResetCacheAll(user.Id)
	return
}

func (t *System) Delete(id, userId int64, userAdmin UserRecord, lang lang.Translation) fiber.Map {
	record := new(Bidder).GetById(id, userId)
	if record.IsLock == mysql.TYPEIsLockTypeLock {
		return fiber.Map{
			"status":  "err",
			"message": "you cannot delete this bidder",
			"id":      id,
		}
	}
	err := mysql.Client.Model(&BidderRecord{}).Delete(&BidderRecord{}, "id = ? and user_id = ?", id, userId).Error
	if err != nil {
		if !utility.IsWindow() {
			return fiber.Map{
				"status":  "err",
				"message": lang.Errors.BidderError.Delete.ToString(),
				"id":      id,
			}
		}
		return fiber.Map{
			"status":  "err",
			"message": err.Error(),
			"id":      id,
		}
	}

	//Delete các relation ship
	//new(LineItemBidderInfo).DeleteBidderByBidderId(id)
	//new(LineItemBidderParams).DeleteByBidderId(id)
	//new(RlBidderMediaType).DeleteByBidderId(id)
	//new(LineItemAccount).DeleteAccountByBidder(id)

	// Xử lý render cache lại toàn bộ domain
	new(Inventory).ResetCacheAll(userId)

	// History
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userId
	}
	_ = history.PushHistory(&history.Bidder{
		Detail:    history.DetailBidderFE,
		CreatorId: creatorId,
		RecordOld: record.TableBidder,
		RecordNew: mysql.TableBidder{},
	})

	return fiber.Map{
		"status":  "success",
		"message": "done",
		"id":      id,
	}

}

func (record *BidderRecord) makeRow(inputs payload.SystemCreate, user UserRecord) {
	//record.Id = inputs.Id
	// Nếu inputs.Id = 0 đây là trường hợp add
	if inputs.Id == 0 {
		// Tạo name và bidder code theo bidder template
		bidderTemplate := new(BidderTemplate).GetById(inputs.BidderId)
		record.BidderCode = bidderTemplate.BidderCode
		if inputs.BidderId == 1 {
			record.DisplayName = inputs.DisplayName
			record.PubId = inputs.PubId
		}
		record.UserId = user.Id
		record.BidderTemplateId = inputs.BidderId
		record.AccountType = inputs.AccountType
		record.IsLock = mysql.TYPEIsLockTypeUnlock
		record.IsDefault = mysql.TypeOff
	}
	record.BidderType = 1 // Bidder pub tạo = 1
	record.BidAdjustment = &inputs.BidAdjustment
	record.RPM = inputs.RPM
	record.AdsTxt = adstxt.StandardizedWithText(inputs.AdsTxt).ToString()
	record.Status = 1 // Mặc định là 1 = on
	// Nếu bidder = 1 là google or bidder = 2 là amazon thì thêm pubId
	if inputs.BidderId == 1 || inputs.BidderId == 2 {
		record.PubId = inputs.PubId
	}
	if inputs.IsLock != 0 {
		record.IsLock = inputs.IsLock
		record.IsDefault = inputs.IsDefault
	}
	// Nếu bidder là Amz
	if inputs.BidderId == 2 {
		// Nếu như đổi linked GAM thì đặt lại scan
		if inputs.LinkedGam != record.LinkedGam {
			record.LinkedGam = inputs.LinkedGam
			record.ScanAmz = 0
			record.LastScanAmz = time.Now().Add(-24 * time.Hour)
		}
	}
	// Nếu user là network thì lưu thêm supply chain domain
	if user.Permission == mysql.UserPermissionNetwork {
		record.SellerDomain = inputs.SupplyChainDomain
	}
	return
}

func (t *System) makeRowConfig(inputs payload.SystemConfig) (record ConfigRecord) {
	record = ConfigRecord{mysql.TableConfig{
		PrebidTimeOut: inputs.PrebidTimeOut,
		AdRefreshTime: inputs.AdRefreshTime,
	}}
	return
}

func (t *System) ValidateCreate(inputs payload.SystemCreate, user UserRecord, lang lang.Translation) (errs []ajax.Error) {
	// Xử lý riêng cho từng bidder
	switch inputs.BidderId {
	case 0: //=> case đặc biệt.
		errs = append(errs, ajax.Error{
			Id:      "select-bidder",
			Message: lang.ErrorRequired.ToString(),
		})
		errs = append(errs, t.ValidateParam(inputs)...)
	case 1: //=> bidder GOOGLE
		errs = append(errs, t.ValidateBidderGoogle(inputs, user, lang)...)
	case 2: //=> bidder AMAZON
		errs = append(errs, t.ValidateBidderAmazon(inputs, user)...)
	default: //=> Các bidder khác
		if new(Bidder).IsExists(inputs.BidderId, user.Id, inputs.Id) {
			errs = append(errs, ajax.Error{
				Id:      "select-bidder",
				Message: "Bidder already exists",
			})
		}
		errs = append(errs, t.ValidateParam(inputs)...)
	}

	// Xử lý cho toàn bộ các bidder
	if inputs.BidAdjustment < 0 || inputs.BidAdjustment > 1 {
		errs = append(errs, ajax.Error{
			Id:      "bid_adjustment",
			Message: "Bid Adjustment only accept values from 0 to 1",
		})
	}

	// Kiểm tra supply chain nếu permission là network
	if user.Permission == 6 {
		if govalidator.IsNull(inputs.SupplyChainDomain) {
			errs = append(errs, ajax.Error{
				Id:      "supply_chain_domain",
				Message: lang.ErrorRequired.ToString(),
			})
		} else if !govalidator.IsDNSName(inputs.SupplyChainDomain) {
			errs = append(errs, ajax.Error{
				Id:      "supply_chain_domain",
				Message: "Domain require",
			})
		}
	}
	return
}

func (t *System) ValidateUpdate(inputs payload.SystemCreate, user UserRecord, lang lang.Translation) (errs []ajax.Error) {
	// Xử lý riêng cho từng bidder
	switch inputs.BidderId {
	case 0: // Case đặc biệt
		errs = append(errs, ajax.Error{
			Id:      "select-bidder",
			Message: lang.ErrorRequired.ToString(),
		})

	case 1: // Bidder GOOGLE
		if len(inputs.MediaType) < 1 {
			errs = append(errs, ajax.Error{
				Id:      "media_type",
				Message: lang.ErrorRequired.ToString(),
			})
		}
		if !govalidator.IsInt(strings.TrimSpace(inputs.PubId)) {
			errs = append(errs, ajax.Error{
				Id:      "pub_id",
				Message: "Pub ID must be a number with minimum 16 length",
			})
		} else if utility.ValidateString(inputs.PubId) == "" {
			errs = append(errs, ajax.Error{
				Id:      "pub_id",
				Message: lang.ErrorRequired.ToString(),
			})
		} else if i, _ := strconv.Atoi(inputs.PubId); i < 1 {
			errs = append(errs, ajax.Error{
				Id:      "pub_id",
				Message: "Pub ID must be a number with minimum 16 length",
			})
		} else if len(strings.TrimSpace(inputs.PubId)) != 16 {
			errs = append(errs, ajax.Error{
				Id:      "pub_id",
				Message: "Pub ID must be 16 numbers",
			})
		}
	case 2: //=> Bidder AMAZON
		errs = append(errs, t.ValidateBidderAmazon(inputs, user)...)
	default: //=> Các bidder khác
		//Get data bidder đang update
		record := new(Bidder).GetById(inputs.Id, user.Id)
		if record.BidderTemplateId != inputs.BidderId {
			//Sửa chọn bidder kiểm tra xem bidder đã tồn tại chưa
			if new(Bidder).IsExists(inputs.BidderId, user.Id, inputs.Id) {
				errs = append(errs, ajax.Error{
					Id:      "select-bidder",
					Message: "Bidder is existed",
				})
			}
		}
		errs = append(errs, t.ValidateParam(inputs)...)

	}

	// Xử lý cho tất cả các bidder
	if inputs.BidAdjustment < 0 || inputs.BidAdjustment > 1 {
		errs = append(errs, ajax.Error{
			Id:      "bid_adjustment",
			Message: "Bid Adjustment only accept values from 0 to 1",
		})
	}

	// Kiểm tra supply chain nếu permission là network
	if user.Permission == 6 {
		if govalidator.IsNull(inputs.SupplyChainDomain) {
			errs = append(errs, ajax.Error{
				Id:      "supply_chain_domain",
				Message: lang.ErrorRequired.ToString(),
			})
		} else if !govalidator.IsDNSName(inputs.SupplyChainDomain) {
			errs = append(errs, ajax.Error{
				Id:      "supply_chain_domain",
				Message: "Domain require",
			})
		}
	}
	return
}

func (t *System) ValidateParam(inputs payload.SystemCreate) (errs []ajax.Error) {
	mapCheckParam := new(PbBidder).GetMapCheckParamFloor()
	for _, param := range inputs.Params {
		if govalidator.IsNull(param.Name) {
			errs = append(errs, ajax.Error{
				Id:      "param_name_" + strconv.Itoa(param.Index),
				Message: "Param name is not null!",
			})
			continue
		}
		if _, ok := mapCheckParam[param.Name]; ok {
			errs = append(errs, ajax.Error{
				Id:      "param_name_" + strconv.Itoa(param.Index),
				Message: "Param " + param.Name + " isn't valid!",
			})
		}
	}
	return
}

func (t *System) ValidateBidderAmazon(inputs payload.SystemCreate, user UserRecord) (errs []ajax.Error) {
	if new(Bidder).IsExists(inputs.BidderId, user.Id, inputs.Id) {
		errs = append(errs, ajax.Error{
			Id:      "select-bidder",
			Message: "Bidder already exists",
		})
	}
	if govalidator.IsNull(inputs.PubId) {
		errs = append(errs, ajax.Error{
			Id:      "pub_id",
			Message: lang.Translate.ErrorRequired.ToString(),
		})
	} else {
		r := regexp.MustCompile(`^[a-zA-Z0-9\-]+$`)
		if !r.MatchString(inputs.PubId) {
			errs = append(errs, ajax.Error{
				Id:      "pub_id",
				Message: `PubID contains only uppercase, normal case, number and "-".`,
			})
		}
	}
	if inputs.LinkedGam == 0 {
		errs = append(errs, ajax.Error{
			Id:      "linked_gam",
			Message: lang.Translate.ErrorRequired.ToString(),
		})
	}
	return
}

func (t *System) ValidateBidderGoogle(inputs payload.SystemCreate, user UserRecord, lang lang.Translation) (errs []ajax.Error) {
	// bidder id = 1 là google chỉ cần lưu display name
	if utility.ValidateString(inputs.DisplayName) == "" {
		errs = append(errs, ajax.Error{
			Id:      "display_name",
			Message: lang.ErrorRequired.ToString(),
		})
	} else if len(inputs.DisplayName) > 15 {
		errs = append(errs, ajax.Error{
			Id:      "display_name",
			Message: "Maximum Display Name is 15 characters",
		})
	} else if new(Bidder).IsDisplayNameGoogleUnique(inputs.DisplayName, user.Id) {
		errs = append(errs, ajax.Error{
			Id:      "display_name",
			Message: "Display Name already exists",
		})
	} else if new(BidderTemplate).CheckUniqueGoogleDisplayNameWithPrebidDemand(inputs.DisplayName) {
		errs = append(errs, ajax.Error{
			Id:      "display_name",
			Message: "The display name cannot be the same as the list of bidders already in place",
		})
	} else {
		r := regexp.MustCompile(`^[a-z0-9_]+$`)
		if !r.MatchString(inputs.DisplayName) {
			errs = append(errs, ajax.Error{
				Id:      "display_name",
				Message: `Display Name contains only normal case, number and "_". Max 15 characters`,
			})
		}
	}
	if !govalidator.IsInt(strings.TrimSpace(inputs.PubId)) {
		errs = append(errs, ajax.Error{
			Id:      "pub_id",
			Message: "Pub ID must be a number with minimum 16 length",
		})
	} else if utility.ValidateString(inputs.PubId) == "" {
		errs = append(errs, ajax.Error{
			Id:      "pub_id",
			Message: lang.ErrorRequired.ToString(),
		})
	} else if i, _ := strconv.Atoi(inputs.PubId); i < 1 {
		errs = append(errs, ajax.Error{
			Id:      "pub_id",
			Message: "Pub ID must be a number with minimum 16 length",
		})
	} else if len(strings.TrimSpace(inputs.PubId)) != 16 {
		errs = append(errs, ajax.Error{

			Id:      "pub_id",
			Message: "Pub ID must be 16 numbers",
		})
	}
	if len(inputs.MediaType) < 1 {
		errs = append(errs, ajax.Error{
			Id:      "media_type",
			Message: lang.ErrorRequired.ToString(),
		})
	}
	return
}

func (t *System) ValidateConfig(inputs payload.SystemConfig, record ConfigRecord) (errs []ajax.Error) {
	// if inputs.PrebidTimeOut == 0 {
	// 	errs = append(errs, ajax.Error{
	// 		Id:      "prebid_time_out",
	// 		Message: "(*) required",
	// 	})
	// }

	// if inputs.AdRefreshTime == 0 {
	// 	errs = append(errs, ajax.Error{
	// 		Id:      "ad_refresh_time",
	// 		Message: "(*) required",
	// 	})
	// }
	if govalidator.IsNull(record.Currency) {
		if inputs.Currency == 0 {
			errs = append(errs, ajax.Error{
				Id:      "currency",
				Message: "(*) required",
			})
		}
	}
	return
}

// GetAll get all bidder default từ system
func (t *System) GetListBidderIdDefault() (listBidderId []int64) {
	mysql.Client.Model(&BidderRecord{}).Select("id").Where("is_default = 1 and user_id = 0").Find(&listBidderId)
	return
}

func (t *System) GetByFilters(inputs *payload.SystemFilterPayload, user UserRecord, lang lang.Translation) (response datatable.Response, err error) {
	//var listBidderDefaultApproved []int64
	//listBidderIdDefault := t.GetListBidderIdDefault()
	//for _, bidderId := range listBidderIdDefault {
	//	if new(RlsBidderSystemInventory).CheckApproveBidderAdxByUser(bidderId, user.Id) {
	//		listBidderDefaultApproved = append(listBidderDefaultApproved, bidderId)
	//	}
	//}
	var bidders []BidderRecord
	var total int64
	err = mysql.Client.Where("user_id = ? ", user.Id).
		Scopes(
			t.SetFilterStatus(inputs),
			t.setFilterSearch(inputs),
			t.SetFilterMediaType(inputs),
		).
		Model(&bidders).Count(&total).
		Scopes(
			t.setOrder(inputs),
			pagination.Paginate(pagination.Params{
				Limit:  inputs.Length,
				Offset: inputs.Start,
			}),
		).
		Find(&bidders).Error
	if err != nil {
		if !utility.IsWindow() {
			err = fmt.Errorf(lang.Errors.BidderError.List.ToString())
		}
		return datatable.Response{}, err
	}
	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.MakeResponseDatatable(bidders, user)
	return
}

type BidderDatatable struct {
	BidderRecord
	RowId      string `json:"DT_RowId"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	BidderCode string `json:"bidder_code"`
	MediaType  string `json:"media_type"`
	Status     string `json:"status"`
	Action     string `json:"action"`
}

func (t *System) MakeResponseDatatable(bidders []BidderRecord, user UserRecord) (records []BidderDatatable) {
	for _, bidder := range bidders {
		var recordMediaTypes []MediaTypeRecord
		recordRlMediaType := new(RlBidderMediaType).GetByBidderId(bidder.Id)
		for _, v := range recordRlMediaType {
			recordMediaType := new(MediaType).GetById(v.MediaTypeId)
			recordMediaTypes = append(recordMediaTypes, recordMediaType)
		}

		bidderCode := new(BidderTemplate).GetBidderCodeById(bidder.BidderTemplateId)
		if bidder.DisplayName != "" {
			// bidderCode += fmt.Sprintf("<div><small class='text-muted' title='display name'>%s</small></div>", bidder.DisplayName)
			bidderCode = bidder.DisplayName
		}

		layoutTime := "3:04:05 PM, January 2, 2006"
		timeConfigTime := bidder.CreatedAt.Format(layoutTime)
		var description string
		if bidder.IsLock == mysql.TYPEIsLockTypeLock {
			if bidder.IsDefault != mysql.TypeOn {
				description = "Auto created at " + timeConfigTime
			}
		} else {
			description = "Created at " + timeConfigTime
		}
		bidAdjustment := 1.0
		if bidder.AccountType.IsAdx() && bidder.IsDefault == mysql.TypeOn {
			bidder.BidAdjustment = &bidAdjustment
		}
		rec := BidderDatatable{
			BidderRecord: bidder,
			RowId:        "bidder_" + strconv.FormatInt(bidder.Id, 10),
			Name:         htmlblock.Render("bidder/index/name.block.gohtml", fiber.Map{"bidder": bidder, "description": description}).String(),
			Type:         htmlblock.Render("bidder/index/type.block.gohtml", bidder).String(),
			BidderCode:   bidderCode,
			MediaType:    htmlblock.Render("bidder/index/media_type.block.gohtml", recordMediaTypes).String(),
			Status:       htmlblock.Render("bidder/index/status.block.gohtml", bidder).String(),
			Action:       htmlblock.Render("bidder/index/action.block.gohtml", bidder).String(),
		}
		records = append(records, rec)
	}
	return
}

func (t *System) SetFilterStatus(inputs *payload.SystemFilterPayload) func(db *gorm.DB) *gorm.DB {
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

func (t *System) SetFilterMediaType(inputs *payload.SystemFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.MediaType != nil {
			switch inputs.PostData.MediaType.(type) {
			case string, int:
				var listBidderId []int64
				mysql.Client.Model(RlsBidderMediaTypeRecord{}).
					Select("bidder_id").
					Where("bidder_id != 0 and media_type_id = ?", inputs.PostData.MediaType).
					Find(&listBidderId)
				if len(listBidderId) != 0 {
					return db.Where("id IN ?", listBidderId)
				}
			case []string, []interface{}:
				var listBidderId []int64
				mysql.Client.Model(RlsBidderMediaTypeRecord{}).
					Select("bidder_id").
					Where("bidder_id != 0 and media_type_id IN ?", inputs.PostData.MediaType).
					Find(&listBidderId)
				if len(listBidderId) != 0 {
					return db.Where("id IN ?", listBidderId)
				}
			}
		}
		return db
	}
}

func (t *System) setFilterSearch(inputs *payload.SystemFilterPayload) func(db *gorm.DB) *gorm.DB {
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
		return db.Where("display_name LIKE ? OR bidder_code LIKE ? OR show_on_pub LIKE ?", "%"+inputs.PostData.QuerySearch+"%", "%"+inputs.PostData.QuerySearch+"%", "%"+inputs.PostData.QuerySearch+"%")
	}
}

func (t *System) setOrder(inputs *payload.SystemFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var orders []string
		orders = append(orders, "is_lock DESC")
		if len(inputs.Order) > 0 {
			for _, order := range inputs.Order {
				column := inputs.Columns[order.Column]
				if column.Data == "name" {
					column.Data = "bidder_code"
				}
				orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
			}
		} else {
			orders = append(orders, "id DESC")
		}
		orderString := strings.Join(orders, ", ")
		return db.Order(orderString)
	}
}

func (t *System) VerificationRecord(id, userId int64) (record BidderRecord, err error) {
	err = mysql.Client.Model(&BidderRecord{}).Where("id = ? and user_id = ?", id, userId).Find(&record).Error
	return
}

func (t *System) VerificationRecordConfig(userId int64) (record ConfigRecord, err error) {
	err = mysql.Client.Model(&ConfigRecord{}).Where("user_id = ?", userId).Find(&record).Error
	return
}

func (t *System) LogXlsxForAmz(bidderId int64, xlsxPath string) (err error) {
	f, err := excelize.OpenFile(AssetsPath + "/../" + xlsxPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Get all the rows in the Sheet1.
	cols, err := f.GetCols("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	var rowCustomTarget, rowRate []string
	for _, col := range cols {
		checkColCustomTarget := false
		checkColRate := false
		for _, rowCel := range col {
			if checkColCustomTarget {
				rowCustomTarget = append(rowCustomTarget, rowCel)
			} else if checkColRate {
				rowRate = append(rowRate, rowCel)
			}
			if rowCel == "Custom targeting" {
				checkColCustomTarget = true
			} else if rowCel == "Rate" {
				checkColRate = true
			}
		}
	}
	// Tạo mảng rates để sort rate float64
	var rates []float64
	// Tạo map chứa key và cmp và đưa dữ liệu vào từ row
	mapCPM := make(map[float64]string)
	if len(rowRate) > 0 && len(rowCustomTarget) == len(rowRate) {
		for i, _ := range rowCustomTarget {
			rate, _ := strconv.ParseFloat(rowRate[i], 64)
			rates = append(rates, rate)
			// Customtarget có dạng amznbid = "yd4mio" phân tích để lấy được "yd4mio"
			splitCustomTarget := strings.Split(rowCustomTarget[i], "\"")
			key := splitCustomTarget[1]
			mapCPM[rate] = key
		}
	}
	sort.Float64s(rates)
	for _, rate := range rates {
		//Log vào db
		mysql.Client.Create(&LogCpmAmzRecord{mysql.TableLogCpmAmz{
			BidderId: bidderId,
			XlsxPath: xlsxPath,
			Target:   mapCPM[rate],
			Value:    math.Round(rate*100) / 100,
		}})
	}
	return
}

func (t *System) BuildCpmFromXlsxAmz(xlsxPath string) (objectCpm map[string]float64, err error) {
	objectCpm = make(map[string]float64)
	f, err := excelize.OpenFile(AssetsPath + "/../" + xlsxPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Get all the rows in the Sheet1.
	cols, err := f.GetCols("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	var rowCustomTarget, rowRate []string
	for _, col := range cols {
		checkColCustomTarget := false
		checkColRate := false
		for _, rowCel := range col {
			if checkColCustomTarget {
				rowCustomTarget = append(rowCustomTarget, rowCel)
			} else if checkColRate {
				rowRate = append(rowRate, rowCel)
			}
			if rowCel == "Custom targeting" {
				checkColCustomTarget = true
			} else if rowCel == "Rate" {
				checkColRate = true
			}
		}
	}
	// Tạo map chứa key và cmp và đưa dữ liệu vào từ row
	if len(rowRate) > 0 && len(rowCustomTarget) == len(rowRate) {
		for i, _ := range rowCustomTarget {
			rate, _ := strconv.ParseFloat(rowRate[i], 64)
			// Customtarget có dạng amznbid = "yd4mio" phân tích để lấy được "yd4mio"
			splitCustomTarget := strings.Split(rowCustomTarget[i], "\"")
			key := splitCustomTarget[1]
			objectCpm[key] = rate
		}
	}

	return
}
