package model

import (
	"source/apps/frontend/ggapi"
	"source/apps/frontend/lang"
	"source/apps/frontend/payload"
	"source/apps/history"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"source/pkg/pagination"
	"source/pkg/utility"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	// "fmt"
)

type AdTag struct{}

type AdTagRecord struct {
	mysql.TableInventoryAdTag
}

func (AdTagRecord) TableName() string {
	return mysql.Tables.InventoryAdTag
}

func (t *AdTag) Create(inputs payload.AdTagAdd, userLogin UserRecord, userAdmin UserRecord, lang lang.Translation) (record AdTagRecord, errs []ajax.Error) {
	// Validate inputs
	errs = t.ValidateCreate(inputs, userLogin)
	if len(errs) > 0 {
		return
	}
	// Insert to database
	record, _ = t.makeInfoCreate(inputs, userLogin)
	err := mysql.Client.Create(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.AdTagError.Add.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}
	//Mỗi khi tạo mới adTag tiến hành tìm và push lại các line target đến tag
	record.PushToQueueWorkerLineItemDfp()

	//Create Size Additional
	inputs.Id = record.Id
	t.createSizeAdditional(record, inputs)
	new(Inventory).ResetCacheWorker(inputs.InventoryId)

	// Kiểm tra xử lý cho adsense adslot
	_ = new(ApiAdsense).HandlerAdTagPostAPIAdSlot(record)

	// Push history
	record.GetFullData()
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userLogin.Id
	}
	_ = history.PushHistory(&history.AdTag{
		Detail:    history.DetailAdTagFE,
		CreatorId: creatorId,
		RecordOld: mysql.TableInventoryAdTag{},
		RecordNew: record.TableInventoryAdTag,
	})
	return
}

func (t *AdTag) Update(inputs payload.AdTagAdd, userLogin UserRecord, userAdmin UserRecord, lang lang.Translation) (record AdTagRecord, errs []ajax.Error) {
	recordOld, _ := t.GetById(inputs.Id, inputs.InventoryId)
	// Đặt lại các giá field không thay đổi từ recordOld
	inputs.Type = recordOld.Type
	if inputs.Type == mysql.TYPEDisplay {
		inputs.PrimaryAdSize = recordOld.PrimaryAdSize
	}
	if inputs.Type == mysql.TYPEStickyBanner {
		inputs.SizeSticky = recordOld.PrimaryAdSize
	}
	if inputs.Type == mysql.TYPEStickyBanner && recordOld.PrimaryAdSizeMobile != 0 {
		inputs.SizeStickyMobile = recordOld.PrimaryAdSizeMobile
	}

	// Validate inputs
	errs = t.ValidateEdit(inputs, recordOld, userLogin)
	if len(errs) > 0 {
		return
	}
	record = recordOld
	// Insert to database
	_ = record.makeInfoUpdate(inputs, userLogin, recordOld)
	err := mysql.Client.Save(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.AdTagError.Edit.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}

	//Delete Additional Old
	new(AdTagSizeAdditional).DeleteAllUnscopedByTagId(record.Id)

	//Create Size Additional
	inputs.Id = record.Id
	t.createSizeAdditional(record, inputs)
	new(Inventory).ResetCacheWorker(inputs.InventoryId)
	// Push history
	record.GetFullData()
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userLogin.Id
	}
	_ = history.PushHistory(&history.AdTag{
		Detail:    history.DetailAdTagFE,
		CreatorId: creatorId,
		RecordOld: recordOld.TableInventoryAdTag,
		RecordNew: record.TableInventoryAdTag,
	})
	return
}

func (t *AdTag) GetAll() (records []AdTagRecord) {
	mysql.Client.Find(&records)
	return
}

func (t *AdTag) CountData(value string) (count int64) {
	mysql.Client.Model(&AdTagRecord{}).Where("name like ?", "%"+value+"%").Count(&count)
	return
}

func (t *AdTag) LoadMoreData(key, value string) (rows []AdTagRecord, isMoreData bool) {
	page, offset := pagination.Pagination(key, 10)
	mysql.Client.Where("name like ?", "%"+value+"%").Limit(10).Offset(offset).Find(&rows)
	total := t.CountData(value)
	totalPages := int(total) / 10
	if (int(total) % 10) != 0 {
		totalPages++
	}
	if page < totalPages {
		isMoreData = true
	}
	return
}

func (t *AdTag) makeInfoCreate(inputs payload.AdTagAdd, userLogin UserRecord) (record AdTagRecord, mapRecord map[string]interface{}) {
	// mapRecord để xử lý các trường hợp đặc biêt
	mapRecord = make(map[string]interface{})
	cfDomain := new(InventoryConfig).GetByInventoryId(inputs.InventoryId)
	// record.Id = inputs.Id
	record.Name = inputs.Name
	record.Type = inputs.Type
	record.UserId = userLogin.Id
	record.Status = inputs.Status
	record.InventoryId = inputs.InventoryId
	record.PrimaryAdSizeMobile = 0
	switch inputs.Type {
	case mysql.TYPEDisplay:
		if cfDomain.GamAutoCreate == 2 {
			record.Gam = inputs.GamDisplay
		}
		record.AdSize = inputs.AdSize
		if inputs.AdSize == mysql.TYPEAdSizeFixed {
			record.PrimaryAdSize = inputs.PrimaryAdSize
		} else if inputs.AdSize == mysql.TYPEAdSizeResponsive {
			record.ResponsiveType = inputs.ResponsiveType
		}
		record.PassBack = strings.TrimSpace(inputs.PassBack)
		record.PrimaryAdSizeMobile = inputs.SizeOnMobile
		record.PassBackMobile = strings.TrimSpace(inputs.PassBackMobile)
		record.AdRefresh = inputs.AdRefresh
		record.AdRefreshTime = inputs.AdRefreshTime
		record.Sticky4k = inputs.Sticky4k
		break
	case mysql.TYPEInStream:
		if cfDomain.GamAutoCreate == 2 {
			record.Gam = inputs.GamInStream
		}
		record.TemplateId = inputs.TemplateId
		record.ContentSource = inputs.ContentSource
		switch record.ContentSource {
		case mysql.TypeContentSourcePlaylist:
			record.PlaylistId = inputs.PlaylistId
			break
		case mysql.TypeContentSourceFeed:
			record.FeedUrl = inputs.FeedUrl
			break
		}
		record.Renderer = inputs.RendererInstream
		record.BannerAD = inputs.BannerAD
		record.VideoAD = inputs.VideoAD
		break
	case mysql.TYPEOutStream:
		record.TemplateId = inputs.TemplateOutStream
		record.PassBackType = inputs.PassBackTypeOutStream
		if record.PassBackType == 1 {
			record.InlineTag = inputs.InlineTagOutStream
		} else {
			record.PassBack = strings.TrimSpace(inputs.PassBackOutStream)
		}
		record.Renderer = inputs.RendererOutStream
		break
	case mysql.TYPETopArticles:
		record.TemplateId = inputs.TemplateArticles
		record.ContentSource = inputs.ContentSourceArticles
		switch record.ContentSource {
		case mysql.TypeContentSourceFeed:
			record.FeedUrl = inputs.FeedArticles
			break
		case mysql.TypeContentSourceAuto:
			record.RelatedContent = mysql.TYPERelatedContentMostViewed // Mặc định là mostviewed
			break
		}
		break
	case mysql.TYPEStickyBanner:
		if cfDomain.GamAutoCreate == 2 {
			record.GamSticky = inputs.GamSticky
		}
		record.ShiftContent = inputs.ShiftContent
		record.PrimaryAdSize = inputs.SizeSticky
		record.PrimaryAdSizeMobile = inputs.SizeStickyMobile
		record.PositionSticky = inputs.PositionSticky
		record.PositionStickyMobile = inputs.PositionStickyMobile
		record.CloseButtonSticky = inputs.CloseButtonSticky
		record.CloseButtonStickyMobile = inputs.CloseButtonStickyMobile
		record.EnableStickyDesktop = inputs.EnableStickyDesktop
		record.EnableStickyMobile = inputs.EnableStickyMobile
		record.PassBack = strings.TrimSpace(inputs.PassBackSticky)
		record.PassBackMobile = strings.TrimSpace(inputs.PassBackStickyMobile)
		break
	case mysql.TYPEInterstitial:
		record.FrequencyCaps = inputs.FrequencyCaps
		break
	case mysql.TYPEPlayZone:
		record.TemplateId = inputs.TemplatePlayZone
		record.ContentType = inputs.ContentType
		record.ContentSource = inputs.ContentSourcePlayZone
		record.MainTitle = inputs.MainTitle
		record.BackgroundColor = inputs.BackgroundColor
		record.TitleColor = inputs.TitleColor
		record.PassBack = strings.TrimSpace(inputs.PassBackPlayZone)
		break
	case mysql.TYPEPNative:
		record.TemplateId = inputs.TemplateNative
		break
	}
	record.Uuid = uuid.New().String()
	record.BidOutStream = inputs.BidOutStream
	return
}

func (t *AdTag) createSizeAdditional(record AdTagRecord, inputs payload.AdTagAdd) {
	modelAdSizeAdditional := new(AdTagSizeAdditional)
	if record.Type == mysql.TYPEDisplay {
		for _, v := range inputs.AdditionalAdSize {
			if v > 0 && record.PrimaryAdSize > 0 {
				modelAdSizeAdditional.Create(inputs.Id, 1, v)
			}
		}
		for _, v := range inputs.AdditionalAdSizeMobile {
			if v > 0 && record.PrimaryAdSizeMobile > 0 {
				modelAdSizeAdditional.Create(inputs.Id, 2, v)
			}
		}
	} else if record.Type == mysql.TYPEStickyBanner {
		for _, v := range inputs.AdditionalAdSizeDesktopStick {
			if v > 0 {
				modelAdSizeAdditional.Create(inputs.Id, 1, v)
			}
		}
		for _, v := range inputs.AdditionalAdSizeMobileStick {
			if v > 0 {
				modelAdSizeAdditional.Create(inputs.Id, 2, v)
			}
		}
	}
}

func (adTag *AdTagRecord) makeInfoUpdate(inputs payload.AdTagAdd, userLogin UserRecord, oldRecord AdTagRecord) (mapRecord map[string]interface{}) {
	// mapRecord để xử lý các trường hợp đặc biêt
	mapRecord = make(map[string]interface{})
	cfDomain := new(InventoryConfig).GetByInventoryId(inputs.InventoryId)
	adTag.Id = inputs.Id
	adTag.Name = inputs.Name
	adTag.Type = oldRecord.Type
	adTag.UserId = userLogin.Id
	adTag.Status = inputs.Status
	adTag.InventoryId = inputs.InventoryId
	adTag.BidOutStream = inputs.BidOutStream
	adTag.SizeOnMobile = 0
	adTag.PrimaryAdSizeMobile = 0
	adTag.PrimaryAdSize = 94
	switch oldRecord.Type {
	case mysql.TYPEDisplay:
		if cfDomain.GamAutoCreate == 2 {
			adTag.Gam = inputs.GamDisplay
		} else {
			adTag.Gam = oldRecord.Gam
		}
		adTag.AdSize = inputs.AdSize
		adTag.PrimaryAdSize = inputs.PrimaryAdSize
		adTag.ResponsiveType = inputs.ResponsiveType
		adTag.PassBack = strings.TrimSpace(inputs.PassBack)
		adTag.PrimaryAdSizeMobile = inputs.SizeOnMobile
		adTag.PassBackMobile = strings.TrimSpace(inputs.PassBackMobile)
		adTag.AdRefresh = inputs.AdRefresh
		adTag.AdRefreshTime = inputs.AdRefreshTime
		adTag.Sticky4k = inputs.Sticky4k
		break
	case mysql.TYPEInStream:
		if cfDomain.GamAutoCreate == 2 {
			adTag.Gam = inputs.GamInStream
		} else {
			adTag.Gam = oldRecord.Gam
		}
		adTag.Renderer = oldRecord.Renderer
		adTag.TemplateId = inputs.TemplateId
		adTag.ContentSource = inputs.ContentSource
		switch adTag.ContentSource {
		case mysql.TypeContentSourcePlaylist:
			adTag.PlaylistId = inputs.PlaylistId
			break
		case mysql.TypeContentSourceFeed:
			adTag.FeedUrl = inputs.FeedUrl
			break
		}
		adTag.BannerAD = inputs.BannerAD
		adTag.VideoAD = inputs.VideoAD
		break
	case mysql.TYPEOutStream:
		adTag.TemplateId = inputs.TemplateOutStream
		adTag.PassBackType = inputs.PassBackTypeOutStream
		if adTag.PassBackType == 1 {
			adTag.InlineTag = inputs.InlineTagOutStream
		} else {
			adTag.PassBack = strings.TrimSpace(inputs.PassBackOutStream)
		}
		adTag.Renderer = oldRecord.Renderer
		break
	case mysql.TYPETopArticles:
		adTag.TemplateId = inputs.TemplateArticles
		adTag.ContentSource = inputs.ContentSourceArticles
		switch adTag.ContentSource {
		case mysql.TypeContentSourceFeed:
			adTag.FeedUrl = inputs.FeedArticles
			break
		case mysql.TypeContentSourceAuto:
			adTag.RelatedContent = mysql.TYPERelatedContentMostViewed
			break
		}
		break
	case mysql.TYPEStickyBanner:
		if cfDomain.GamAutoCreate == 2 {
			adTag.GamSticky = inputs.GamSticky
		} else {
			adTag.GamSticky = oldRecord.GamSticky
		}
		adTag.ShiftContent = inputs.ShiftContent
		if inputs.SizeSticky != 0 {
			adTag.PrimaryAdSize = inputs.SizeSticky
		}
		adTag.PrimaryAdSizeMobile = inputs.SizeStickyMobile
		adTag.PositionSticky = inputs.PositionSticky
		adTag.PositionStickyMobile = inputs.PositionStickyMobile
		adTag.CloseButtonSticky = inputs.CloseButtonSticky
		adTag.CloseButtonStickyMobile = inputs.CloseButtonStickyMobile
		adTag.EnableStickyDesktop = inputs.EnableStickyDesktop
		adTag.EnableStickyMobile = inputs.EnableStickyMobile
		adTag.PassBack = strings.TrimSpace(inputs.PassBackSticky)
		adTag.PassBackMobile = strings.TrimSpace(inputs.PassBackStickyMobile)
		break
	case mysql.TYPEInterstitial:
		adTag.FrequencyCaps = inputs.FrequencyCaps
		break
	case mysql.TYPEPlayZone:
		adTag.TemplateId = inputs.TemplatePlayZone
		adTag.ContentType = inputs.ContentType
		adTag.ContentSource = inputs.ContentSourcePlayZone
		adTag.MainTitle = inputs.MainTitle
		adTag.BackgroundColor = inputs.BackgroundColor
		adTag.TitleColor = inputs.TitleColor
		adTag.PassBack = strings.TrimSpace(inputs.PassBackPlayZone)
		break
	case mysql.TYPEPNative:
		adTag.TemplateId = inputs.TemplateNative
		break
	}
	adTag.CreatedAt = oldRecord.CreatedAt
	return

}

func (t *AdTag) ValidateCreate(inputs payload.AdTagAdd, user UserRecord) (errs []ajax.Error) {
	if utility.ValidateString(inputs.Name) == "" {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Name is required",
		})
	}
	flagName := t.VerificationRecordName(inputs.Name, inputs.Id, user.Id, inputs.InventoryId)
	if !flagName {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Ad Tag Name already exist",
		})
	}
	if inputs.Type == 0 {
		errs = append(errs, ajax.Error{
			Id:      "ad_tag_type",
			Message: "Ad Type is required",
		})
	}
	cfDomain := new(InventoryConfig).GetByInventoryId(inputs.InventoryId)
	switch inputs.Type {
	case mysql.TYPEDisplay:
		if utility.ValidateString(inputs.Sticky4k) == "" {
			errs = append(errs, ajax.Error{
				Id:      "sticky4k",
				Message: "Sticky 4k is required",
			})
		}
		if inputs.AdSize == mysql.TYPEAdSizeFixed {
			if inputs.PrimaryAdSize == 0 {
				errs = append(errs, ajax.Error{
					Id:      "primary_ad_size",
					Message: "Size is required",
				})
			} else {
				size := new(AdSize).GetById(inputs.PrimaryAdSize)
				listSizeEnableBidOutStream := []string{"300x250", "336x280", "300x600", "970x250"}
				if !utility.InArray(size.Name, listSizeEnableBidOutStream, false) {
					if inputs.BidOutStream == mysql.TypeOn {
						errs = append(errs, ajax.Error{
							Id:      "bid_out_stream",
							Message: "Bid Out Stream of size " + size.Name + " can't on",
						})
					}
				}
			}
		}
		if cfDomain.GamAutoCreate == 2 {
			if utility.ValidateString(inputs.GamDisplay) != "" {
				errs = append(errs, t.ValidateGAM(inputs.GamDisplay, "gam_display", user)...)
			}
		}
		break
	case mysql.TYPEInStream:
		if inputs.RendererInstream == 1 {
			if inputs.ContentSource == 0 {
				errs = append(errs, ajax.Error{
					Id:      "content_source",
					Message: "Content Source is required",
				})
			}
			if inputs.ContentSource == 2 {
				if utility.ValidateString(inputs.FeedUrl) == "" {
					errs = append(errs, ajax.Error{
						Id:      "feed",
						Message: "Feed Url is required",
					})
				}
				if isUrl := utility.IsUrl(inputs.FeedUrl); !isUrl {
					errs = append(errs, ajax.Error{
						Id:      "feed",
						Message: "malformed Feed Url",
					})
				}
			}
			if cfDomain.GamAutoCreate == 2 {
				if utility.ValidateString(inputs.GamInStream) != "" {
					errs = append(errs, t.ValidateGAM(inputs.GamInStream, "gam_instream", user)...)
				}
			}
			if inputs.ContentSource == 1 {
				if inputs.PlaylistId == 0 {
					errs = append(errs, ajax.Error{
						Id:      "playlist",
						Message: "Playlist is required",
					})
				}
			}
			if inputs.TemplateId == 0 {
				errs = append(errs, ajax.Error{
					Id:      "template",
					Message: "Template is required",
				})
			}
		}
		break
	case mysql.TYPEOutStream:
		if inputs.RendererOutStream == 1 {
			if inputs.TemplateOutStream == 0 {
				errs = append(errs, ajax.Error{
					Id:      "template_outstream",
					Message: "Template is required",
				})
			}
			// if utility.ValidateString(inputs.PassBackOutStream) == "" {
			//	errs = append(errs, ajax.Error{
			//		Id:      "pass_back_outstream",
			//		Message: "Pass Back is required",
			//	})
			// }
		}
		break
	case mysql.TYPETopArticles:
		if inputs.ContentSourceArticles == 2 {
			if utility.ValidateString(inputs.FeedArticles) == "" {
				errs = append(errs, ajax.Error{
					Id:      "feed_articles",
					Message: "Feed Url is required",
				})
			}
			if isUrl := utility.IsUrl(inputs.FeedArticles); !isUrl {
				errs = append(errs, ajax.Error{
					Id:      "feed_articles",
					Message: "malformed Feed Url",
				})
			}
		} else if inputs.ContentSourceArticles == 4 {
			if inputs.TemplateArticles == 0 {
				errs = append(errs, ajax.Error{
					Id:      "template_articles",
					Message: "Template is required",
				})
			}
		}
		break
	case mysql.TYPEStickyBanner:
		if cfDomain.GamAutoCreate == 2 {
			if utility.ValidateString(inputs.GamSticky) != "" {
				errs = append(errs, t.ValidateGAM(inputs.GamSticky, "gam_sticky", user)...)
			}
		}

		//if inputs.EnableStickyDesktop == mysql.On && inputs.SizeSticky == 0 {
		//	errs = append(errs, ajax.Error{
		//		Id:      "size_sticky",
		//		Message: "Size is required",
		//	})
		//}
		//
		//if inputs.EnableStickyMobile == mysql.On && inputs.SizeStickyMobile == 0 {
		//	errs = append(errs, ajax.Error{
		//		Id:      "size_sticky_mobile",
		//		Message: "Size is required",
		//	})
		//}
		break
	case mysql.TYPEInterstitial:
		if t.ExistsInterstitial(inputs.InventoryId) {
			errs = append(errs, ajax.Error{
				Id:      "name",
				Message: "Interstitial already exists!",
			})
		}
		break
	case mysql.TYPEPlayZone:
		if inputs.ContentType == 1 {
			if utility.ValidateString(inputs.MainTitle) == "" {
				errs = append(errs, ajax.Error{
					Id:      "main_title",
					Message: lang.Translate.ErrorRequired.ToString(),
				})
			}
			if utility.ValidateString(inputs.BackgroundColor) == "" {
				errs = append(errs, ajax.Error{
					Id:      "background_color",
					Message: lang.Translate.ErrorRequired.ToString(),
				})
			}
			if utility.ValidateString(inputs.TitleColor) == "" {
				errs = append(errs, ajax.Error{
					Id:      "title_color",
					Message: lang.Translate.ErrorRequired.ToString(),
				})
			}
		}
		break
	case mysql.TYPEPNative:
		if inputs.TemplateNative == 0 {
			errs = append(errs, ajax.Error{
				Id:      "template_native",
				Message: "Template is required",
			})
		}
		break
	}

	return
}

func (t *AdTag) ValidateEdit(inputs payload.AdTagAdd, row AdTagRecord, user UserRecord) (errs []ajax.Error) {
	flag := t.VerificationRecord(inputs.Id, user.Id, inputs.InventoryId)
	if !flag {
		errs = append(errs, ajax.Error{
			Id:      "id",
			Message: "You don't own this Ad Tag",
		})
	}

	if utility.ValidateString(inputs.Name) == "" {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Name is required",
		})
	}
	flagName := t.VerificationRecordName(inputs.Name, inputs.Id, user.Id, inputs.InventoryId)
	if !flagName {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Ad Tag Name already exist",
		})
	}

	if inputs.Type != row.Type {
		errs = append(errs, ajax.Error{
			Id:      "ad_tag_type",
			Message: "This type isn't of the ad tag",
		})
	}
	cfDomain := new(InventoryConfig).GetByInventoryId(inputs.InventoryId)
	switch inputs.Type {
	case mysql.TYPEDisplay:
		if utility.ValidateString(inputs.Sticky4k) == "" {
			errs = append(errs, ajax.Error{
				Id:      "sticky4k",
				Message: "Sticky 4k is required",
			})
		}
		if inputs.AdSize == mysql.TYPEAdSizeFixed {
			if inputs.PrimaryAdSize == 0 {
				errs = append(errs, ajax.Error{
					Id:      "primary_ad_size",
					Message: "Size is required",
				})
			} else {
				size := new(AdSize).GetById(inputs.PrimaryAdSize)
				listSizeEnableBidOutStream := []string{"300x250", "336x280", "300x600", "970x250"}
				if !utility.InArray(size.Name, listSizeEnableBidOutStream, false) {
					if inputs.BidOutStream == mysql.TypeOn {
						errs = append(errs, ajax.Error{
							Id:      "bid_out_stream",
							Message: "Bid Out Stream of size " + size.Name + " can't on",
						})
					}
				}
			}
		}
		if cfDomain.GamAutoCreate == 2 {
			if utility.ValidateString(inputs.GamDisplay) != "" {
				errs = append(errs, t.ValidateGAM(inputs.GamDisplay, "gam_display", user)...)
			}
		}
		break
	case mysql.TYPEInStream:
		if inputs.RendererInstream == 1 {
			if inputs.ContentSource == 0 {
				errs = append(errs, ajax.Error{
					Id:      "content_source",
					Message: "Content Source is required",
				})
			}
			if inputs.ContentSource == 2 {
				if utility.ValidateString(inputs.FeedUrl) == "" {
					errs = append(errs, ajax.Error{
						Id:      "feed",
						Message: "Feed Url is required",
					})
				}
				if isUrl := utility.IsUrl(inputs.FeedUrl); !isUrl {
					errs = append(errs, ajax.Error{
						Id:      "feed",
						Message: "malformed Feed Url",
					})
				}
			}
			if cfDomain.GamAutoCreate == 2 {
				if utility.ValidateString(inputs.GamInStream) != "" {
					errs = append(errs, t.ValidateGAM(inputs.GamInStream, "gam_instream", user)...)
				}
			}
			if inputs.ContentSource == 1 {
				if inputs.PlaylistId == 0 {
					errs = append(errs, ajax.Error{
						Id:      "playlist",
						Message: "Playlist is required",
					})
				}
			}
			if inputs.TemplateId == 0 {
				errs = append(errs, ajax.Error{
					Id:      "template",
					Message: "Template is required",
				})
			}
		}
		break
	case mysql.TYPEOutStream:
		if inputs.RendererOutStream == 1 {
			if inputs.TemplateOutStream == 0 {
				errs = append(errs, ajax.Error{
					Id:      "template_outstream",
					Message: "Template is required",
				})
			}
			// if utility.ValidateString(inputs.PassBackOutStream) == "" {
			//	errs = append(errs, ajax.Error{
			//		Id:      "pass_back_outstream",
			//		Message: "Pass Back is required",
			//	})
			// }
		}
		break
	case mysql.TYPETopArticles:
		if inputs.ContentSourceArticles == 2 {
			if utility.ValidateString(inputs.FeedArticles) == "" {
				errs = append(errs, ajax.Error{
					Id:      "feed_articles",
					Message: "Feed Url is required",
				})
			}
			if isUrl := utility.IsUrl(inputs.FeedArticles); !isUrl {
				errs = append(errs, ajax.Error{
					Id:      "feed_articles",
					Message: "malformed Feed Url",
				})
			}
		} else if inputs.ContentSourceArticles == 4 {
			if inputs.TemplateArticles == 0 {
				errs = append(errs, ajax.Error{
					Id:      "template_articles",
					Message: "Template is required",
				})
			}
		}
		break
	case mysql.TYPEStickyBanner:
		if cfDomain.GamAutoCreate == 2 {
			if utility.ValidateString(inputs.GamSticky) != "" {
				errs = append(errs, t.ValidateGAM(inputs.GamSticky, "gam_sticky", user)...)
			}
		}

		if inputs.EnableStickyDesktop == mysql.On && inputs.SizeSticky == 0 {
			errs = append(errs, ajax.Error{
				Id:      "size_sticky",
				Message: "Size is required",
			})
		}

		if inputs.EnableStickyMobile == mysql.On && inputs.SizeStickyMobile == 0 {
			errs = append(errs, ajax.Error{
				Id:      "size_sticky_mobile",
				Message: "Size is required",
			})
		}
		break

	case mysql.TYPEInterstitial:

		break
	case mysql.TYPEPlayZone:
		if inputs.ContentType == 1 {
			if utility.ValidateString(inputs.MainTitle) == "" {
				errs = append(errs, ajax.Error{
					Id:      "main_title",
					Message: lang.Translate.ErrorRequired.ToString(),
				})
			}
			if utility.ValidateString(inputs.BackgroundColor) == "" {
				errs = append(errs, ajax.Error{
					Id:      "background_color",
					Message: lang.Translate.ErrorRequired.ToString(),
				})
			}
			if utility.ValidateString(inputs.TitleColor) == "" {
				errs = append(errs, ajax.Error{
					Id:      "title_color",
					Message: lang.Translate.ErrorRequired.ToString(),
				})
			}
		}
		break
	case mysql.TYPEPNative:
		if inputs.TemplateNative == 0 {
			errs = append(errs, ajax.Error{
				Id:      "template_native",
				Message: "Template is required",
			})
		}
		break
	}
	return
}

func (t *AdTag) GetById(id, inventoryId int64) (row AdTagRecord, err error) {
	err = mysql.Client.Where("id = ? and inventory_id = ?", id, inventoryId).Find(&row).Error
	// Get Rls
	row.GetFullData()
	return
}

func (t *AdTag) GetDetail(id int64) (row AdTagRecord, err error) {
	err = mysql.Client.Where("id = ?", id).Find(&row).Error
	return
}

func (t *AdTag) CheckRecord(id, inventoryId, userId int64) (row AdTagRecord, err error) {
	err = mysql.Client.Model(&AdTagRecord{}).Where("id = ? and inventory_id = ? and user_id = ?", id, inventoryId, userId).Find(&row).Error
	return
}

func (t *AdTag) GetListBoxCollapse(userId, recordId int64, page, typ string) (list []string) {
	switch typ {
	case "add":
		mysql.Client.Select("box_collapse").Model(PageCollapseRecord{}).Where("user_id = ? and page_collapse = ? and is_collapse = ? and page_type = ?", userId, page, 1, typ).Find(&list)
		return
	case "edit":
		mysql.Client.Select("box_collapse").Model(PageCollapseRecord{}).Where("user_id = ? and page_collapse = ? and is_collapse = ? and page_type = ? and page_id = ?", userId, page, 1, typ, recordId).Find(&list)
		return
	}
	return
}

func (t *AdTag) VerificationRecord(id, userId, inventoryId int64) bool {
	var row AdTagRecord
	err := mysql.Client.Model(&AdTagRecord{}).Where("id = ? and user_id = ? and inventory_id = ?", id, userId, inventoryId).Find(&row).Error
	if err != nil || row.Id == 0 {
		return false
	}
	return true
}

func (t *AdTag) VerificationRecordName(name string, id, userId, inventoryId int64) bool {
	var row AdTagRecord
	err := mysql.Client.Model(&AdTagRecord{}).Where("name = ? and id != ? and user_id = ? and inventory_id = ?", name, id, userId, inventoryId).Find(&row).Error
	if err != nil || row.Id != 0 {
		return false
	}
	return true
}

func (adTag *AdTagRecord) PushToQueueWorkerLineItemDfp() {
	var listLineItemFilterDomain []int64
	mysql.Client.
		//Debug().
		Table(mysql.Tables.Target).
		Select("line_item_id").
		Where("line_item_id != 0 and (inventory_id = -1 or inventory_id = ? ) and (user_id = ? or user_id = 0)", adTag.InventoryId, adTag.UserId).
		Group("line_item_id").
		Find(&listLineItemFilterDomain)
	// fmt.Println(listLineItemFilterDomain)

	var listLineItemFilter []int64
	mysql.Client.
		//Debug().
		Table(mysql.Tables.LineItem).
		Select("id").
		Where("id in (" + utility.ArrayToString(listLineItemFilterDomain, ",") + ") and server_type = 2").
		Group("id").
		Find(&listLineItemFilter)
	//fmt.Println(listLineItemFilter)

	if len(listLineItemFilter) < 1 {
		return
	}

	// Từ các line item tiến hành đặt push_line_item_dfp = 1 để xếp vào hàng đợi cho worker
	for _, lineItemId := range listLineItemFilter {
		new(LineItem).PushToQueueWorkerLineItemDfp(lineItemId)
	}
}

func (t *AdTag) GetAdTagDisplay(inventoryId int64) (rows []AdTagRecord) {
	mysql.Client.Where("type = 1 and inventory_id = ?", inventoryId).Find(&rows)
	return
}

func (t *AdTag) GetAdTagForTool(inventoryId int64) (rows []AdTagRecord) {
	mysql.Client.Where("((type = ? and ad_size = 1) or type = ?) and inventory_id = ? and status != 3", mysql.TYPEDisplay, mysql.TYPEStickyBanner, inventoryId).Find(&rows)
	return
}

func (t *AdTag) ValidateGAM(GAM string, idGam string, user UserRecord) (errs []ajax.Error) {
	// GAM đúng /{Network_Id}/{AdUnit_Name}
	splitGam := strings.Split(GAM, "/")
	// Khi split đúng theo GAM đúng sẽ đc 3 phần trong đó index đầu sẽ rỗng
	if len(splitGam) != 3 {
		errs = append(errs, ajax.Error{
			Id:      idGam,
			Message: "GAM isn't valid, ex GAM: /123456/adUnitName",
		})
		return
	}
	if !govalidator.IsNull(splitGam[0]) {
		errs = append(errs, ajax.Error{
			Id:      idGam,
			Message: "GAM isn't valid, ex GAM: /123456/adUnitName",
		})
		return
	}
	//Vị trí index 1 sẽ là networkId và index 2 sẽ là adUnitName
	networkId, err := strconv.ParseInt(splitGam[1], 10, 64)
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      idGam,
			Message: "GAM isn't valid, ex GAM: /123456/adUnitName",
		})
		return
	}
	adUnitName := splitGam[2]
	if networkId < 1 { // networkId không được rỗng
		errs = append(errs, ajax.Error{
			Id:      idGam,
			Message: "GAM isn't valid, ex GAM: /123456/adUnitName",
		})
		return
	}
	if govalidator.IsNull(adUnitName) { // adUnitName không được rỗng
		errs = append(errs, ajax.Error{
			Id:      idGam,
			Message: "GAM isn't valid, ex GAM: /123456/adUnitName",
		})
		return
	}
	// Kiểm tra xem networkId có tồn tại trong GAM của user không
	gamNetwork := new(GamNetwork).GetByNetworkId(networkId, user.Id)
	if gamNetwork.Id == 0 {
		errs = append(errs, ajax.Error{
			Id:      idGam,
			Message: "NetworkId " + splitGam[1] + " not found in the account",
		})
		return
	}
	// Từ gamNetwork lấy gam Id để dùng api
	gam := new(Gam).GetById(gamNetwork.GamId)

	// Kiểm tra xem adUnitName có tồn tại trong DFP không
	checkExistAdUnit, err := ggapi.CheckAdUnitByName(gam.RefreshToken, gamNetwork.NetworkId, gamNetwork.NetworkName, adUnitName)
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      idGam,
			Message: "AdUnitName " + adUnitName + " not found in the DFP",
		})
		return
	}
	if !checkExistAdUnit {

		errs = append(errs, ajax.Error{
			Id:      idGam,
			Message: "AdUnitName " + adUnitName + " not found in the DFP",
		})
		return
	}
	return
}

func (t *AdTag) ExistsInterstitial(inventoryId int64) bool {
	var record InventoryAdTagRecord
	mysql.Client.Where("inventory_id = ? and type = ?", inventoryId, mysql.TYPEInterstitial).Find(&record)
	if record.Id != 0 {
		return true
	} else {
		return false
	}
}
