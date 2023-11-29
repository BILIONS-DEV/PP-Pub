package model

import (
	"fmt"
	"gorm.io/gorm"
	"source/apps/frontend/payload"
	"source/apps/history"
	"time"
	"source/apps/frontend/lang"
	// "reflect"
	"source/core/technology/mysql"
	"source/pkg/pagination"
	"source/pkg/utility"
)

type History struct{}

type HistoryRecord struct {
	mysql.TableHistory
}

func (HistoryRecord) TableName() string {
	return mysql.Tables.History
}

type ResponseHistory struct {
	Row        mysql.TableHistory
	Compare    []history.ResponseCompare
	CreateTime string
}

func (t *History) GetByFilters(inputs *payload.HistoryFilterPayload, user UserRecord) (histories []HistoryRecord, err error) {
	var total int64
	err = mysql.Client.Where("user_id = ?", user.Id).
		Scopes(
			t.SetFilterObject(inputs),
			t.SetFilterObjectId(inputs),
		).
		Model(&histories).Count(&total).
		Order("id desc").
		Scopes(
			pagination.Paginate(pagination.Params{
				Limit:  inputs.Length,
				Offset: inputs.Start,
			}),
		).
		Find(&histories).Error
	if err != nil {
		if !utility.IsWindow() {
			var lang lang.Translation
			err = fmt.Errorf(lang.Errors.HistoryError.List.ToString())
		}
		return
	}

	return
}

func (t *History) SetFilterObject(inputs *payload.HistoryFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// if inputs.PostData.DetailType != nil {
		// 	switch inputs.PostData.DetailType.(type) {
		// 	case string, int:
		// 		if inputs.PostData.DetailType != "" {
		// 			return db.Where("detail_type = ?", inputs.PostData.DetailType)
		// 		}
		// 	case []string, []interface{}:
		// 		return db.Where("detail_type IN ?", inputs.PostData.DetailType)
		// 	}
		// }
		if inputs.PostData.ObjectPage != "" {
			return db.Where("page = ?", inputs.PostData.ObjectPage)
		}
		return db
	}
}

func (t *History) SetFilterObjectId(inputs *payload.HistoryFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// if inputs.PostData.ObjectId != nil {
		// 	switch inputs.PostData.ObjectId.(type) {
		// 	case string, int:
		// 		if inputs.PostData.ObjectId != "" {
		// 			return db.Where("object_id = ?", inputs.PostData.ObjectId)
		// 		}
		// 	case []string, []interface{}:
		// 		return db.Where("object_id IN ?", inputs.PostData.ObjectId)
		// 	}
		// }
		if inputs.PostData.ObjectId != "" {
			return db.Where("object_id = ?", inputs.PostData.ObjectId)
		}
		return db
	}
}

func (t *History) LoadHistories(inputs payload.LoadHistory, userId int64) (responseHistory []ResponseHistory, err error) {
	histories, err := t.GetHistories(inputs.Object, inputs.Id, userId)
	if err != nil {
		return
	}
	if len(histories) == 0 {
		switch inputs.Object {
		case history.DetailInventoryAdsTxtFE.String():
			break
		case history.DetailInventoryConnectionFE.String():
			break
		default:
			err = t.AddHistory(inputs.Object, inputs.Id, userId)
			histories, err = t.GetHistories(inputs.Object, inputs.Id, userId)
			if err != nil || len(histories) == 0 {
				return
			}
		}
	}
	flagCreate := false
	for _, value := range histories {
		if value.ObjectType == 1 {
			flagCreate = true
		}
	}
	if flagCreate == false {
		err = t.AddHistory(inputs.Object, inputs.Id, userId)
		histories, err = t.GetHistories(inputs.Object, inputs.Id, userId)
	}
	for _, value := range histories {
		var History ResponseHistory
		History.Row = value
		History.Compare, err = history.CompareDataHistory(value)
		if len(History.Compare) == 0 {
			continue
		}
		// timeNow := time.Now()
		// timeZone, _ := timeNow.Zone()
		loc, err := time.LoadLocation("America/New_York")
		if err != nil {
			panic(err)
		}
		// fmt.Printf("%+v\n", value.CreatedAt.In(loc))
		History.CreateTime = value.CreatedAt.In(loc).Format("01/02/2006 15:04:05 AM")
		responseHistory = append(responseHistory, History)
	}
	return
}

func (t *History) AddHistory(object string, objectId, userId int64) (err error) {
	switch object {
	case "line_item_fe":
		recordNew, _ := new(LineItem).GetById(objectId, userId)
		_ = history.PushHistory(&history.LineItem{
			Detail:    history.DetailLineItemFE,
			CreatorId: userId,
			RecordOld: mysql.TableLineItem{},
			RecordNew: recordNew.TableLineItem,
		})
		break
	case "inventory_config_fe":
		recordNew, _ := new(Inventory).GetById(objectId, userId)
		err = history.PushHistory(&history.Inventory{
			Detail:    history.DetailInventoryConfigFE,
			CreatorId: userId,
			RecordNew: recordNew.TableInventory,
		})
		break
	case "inventory_consent_fe":
		recordNew, _ := new(Inventory).GetById(objectId, userId)
		_ = history.PushHistory(&history.Inventory{
			Detail:    history.DetailInventoryConsentFE,
			CreatorId: userId,
			RecordNew: recordNew.TableInventory,
		})
		break
	case "inventory_connection_fe":
		// 	recordNew, _ := new(Inventory).GetById(objectId, userId)
		// 	_ = history.PushHistory(&history.Inventory{
		// 		Detail:    history.DetailInventoryConnectionFE,
		// 		CreatorId: userId,
		// 		RecordNew: recordNew.TableInventory,
		// 	})
		break
	case "inventory_adstxt_fe":
		// recordNew, _ := new(Inventory).GetById(objectId, userId)
		// _ = history.PushHistory(&history.Inventory{
		// 	Detail:    history.DetailInventoryAdsTxtFE,
		// 	CreatorId: userId,
		// 	RecordNew: recordNew.TableInventory,
		// })
		break
	case "ad_tag_fe":
		recordNew := new(InventoryAdTag).GetById(objectId)
		_ = history.PushHistory(&history.AdTag{
			Detail:    history.DetailAdTagFE,
			CreatorId: userId,
			RecordNew: recordNew.TableInventoryAdTag,
		})
		break
	case "blocking_fe":
		recordNew, _ := new(Blocking).GetById(objectId, userId)
		_ = history.PushHistory(&history.Blocking{
			Detail:    history.DetailBlockingFE,
			CreatorId: userId,
			RecordNew: recordNew.TableBlocking,
		})
		break
	case "template_fe":
		recordNew, _ := new(Player).GetById(objectId, userId)
		_ = history.PushHistory(&history.Template{
			Detail:    history.DetailTemplateFE,
			CreatorId: userId,
			RecordNew: recordNew.TableTemplate,
		})
		break
	case "channel_fe":
		recordNew, _ := new(Channels).GetById(objectId, userId)
		_ = history.PushHistory(&history.Channel{
			Detail:    history.DetailChannelFE,
			CreatorId: userId,
			RecordNew: recordNew.TableChannels,
		})
		break
	case "content_fe":
		recordNew, _ := new(Content).GetById(objectId, userId)
		_ = history.PushHistory(&history.Content{
			Detail:    history.DetailContentFE,
			CreatorId: userId,
			RecordNew: recordNew.TableContent,
		})
		break
	case "playlist_fe":
		recordNew, _ := new(Playlist).GetById(objectId, userId)
		_ = history.PushHistory(&history.Playlist{
			Detail:    history.DetailPlaylistFE,
			CreatorId: userId,
			RecordNew: recordNew.TablePlaylist,
		})
		break
	case "identity_fe":
		recordNew, _ := new(Identity).GetById(objectId, userId)
		_ = history.PushHistory(&history.Identity{
			Detail:    history.DetailIdentityFE,
			CreatorId: userId,
			RecordNew: recordNew.TableIdentity,
		})
		break
	case "floor_fe":
		recordNew, _ := new(Floor).GetById(objectId, userId)
		_ = history.PushHistory(&history.Floor{
			Detail:    history.DetailFloorFE,
			CreatorId: userId,
			RecordNew: recordNew.TableFloor,
		})
		break
	case "gam_setup_fe":
		recordNews := new(GamNetwork).GetByGamId(objectId, userId)
		for _, recordNew := range recordNews {
			recordNew.Id = recordNew.GamId
			_ = history.PushHistory(&history.GAM{
				Detail:    history.DetailGAMSetupFE,
				CreatorId: userId,
				RecordNew: recordNew.TableGamNetwork,
			})
		}
		break
	case "bidder_fe":
		recordNew := new(Bidder).GetById(objectId, userId)
		_ = history.PushHistory(&history.Bidder{
			Detail:    history.DetailBidderFE,
			CreatorId: userId,
			RecordNew: recordNew.TableBidder,
		})
		break
	}

	return
}

func (t *History) GetHistories(object string, id, userId int64) (History []mysql.TableHistory, err error) {
	err = mysql.Client.
		Scopes(
			t.SetFilterDetailType(object),
		).
		Where("object_id = ? and creator_id = ? and app = ?", id, userId, "FE").
		Order("created_at DESC").
		Find(&History).Error
	return
}

func (t *History) SetFilterDetailType(object string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if object != "" {
			switch object {
			case "gam_fe":
				objects := []string{history.DetailGAMSignInFE.String(), history.DetailGAMSetupFE.String()}
				return db.Where("detail_type IN ?", objects)
			default:
				return db.Where("detail_type = ?", object)
			}
		}
		return db
	}
}

func (t *History) GetHost(histories []ResponseHistory) (host string) {
	if len(histories) == 0 {
		return
	}
	for _, History := range histories {
		switch History.Row.DetailType {
		case history.DetailABTestFE.String():
			host = ""
			break
		case history.DetailAdTagFE.String():
			adtag := new(InventoryAdTag).GetById(History.Row.ObjectId)
			host = adtag.Name
			break
		case history.DetailBidderFE.String():
			Bidder := new(Bidder).GetById(History.Row.ObjectId, History.Row.CreatorId)
			host = Bidder.DisplayName
			break
		case history.DetailBidderTemplateBE.String():
			BidderTemplate := new(BidderTemplate).GetById(History.Row.ObjectId)
			host = BidderTemplate.DisplayName
			break
		case history.DetailBlockingFE.String():
			Blocking, _ := new(Blocking).GetById(History.Row.ObjectId, History.Row.CreatorId)
			host = Blocking.RestrictionName
			break
		case history.DetailChannelFE.String():
			Channels, _ := new(Channels).GetById(History.Row.ObjectId, History.Row.CreatorId)
			host = Channels.Name
			break
		case history.DetailContentFE.String():
			Content, _ := new(Content).GetById(History.Row.ObjectId, History.Row.CreatorId)
			host = Content.Title
			break
		case history.DetailFloorFE.String():
			Floor, _ := new(Floor).GetById(History.Row.ObjectId, History.Row.CreatorId)
			host = Floor.Name
			break
		case history.DetailGAMSignInFE.String():
			host = ""
			break
		case history.DetailGAMSetupFE.String():
			host = ""
			break
		case history.DetailIdentityFE.String():
			Identity, _ := new(Identity).GetById(History.Row.ObjectId, History.Row.CreatorId)
			host = Identity.Name
			break
		case history.DetailInventorySubmitFE.String(), history.DetailInventoryConfigFE.String(), history.DetailInventoryConsentFE.String(), history.DetailInventoryAdsTxtFE.String(), history.DetailInventoryConnectionFE.String():
			Inventory, _ := new(Inventory).GetById(History.Row.ObjectId, History.Row.CreatorId)
			host = Inventory.Domain
			break
		case history.DetailLineItemFE.String():
			LineItem, _ := new(LineItem).GetById(History.Row.ObjectId, History.Row.CreatorId)
			host = LineItem.Name
			break
		case history.DetailPlaylistFE.String():
			Playlist, _ := new(Playlist).GetById(History.Row.ObjectId, History.Row.CreatorId)
			host = Playlist.Name
			break
		case history.DetailTemplateFE.String():
			Playlist, _ := new(Template).GetById(History.Row.ObjectId, History.Row.CreatorId)
			host = Playlist.Name
			break

		case history.DetailUserProfileFE.String(), history.DetailUserChangePassFE.String(), history.DetailUserBillingFE.String():
			User := new(User).GetById(History.Row.ObjectId)
			host = User.Email
			break

		default:
			host = ""
			break
		}
	}
	return
}

func (t *History) GetAllHistoryPage() (pages []string) {
	var histories []HistoryRecord
	mysql.Client.Where("app = 'FE'").Group("page").Find(&histories)
	if len(histories) == 0 {
		return
	}

	for _, history := range histories {
		pages = append(pages, history.Page)
	}
	return
}

func (t *History) LoadObjectsByPage(inputs payload.LoadObjectByPage, userId int64) (objects []payload.ObjectPage, err error) {
	switch inputs.ObjectPage {
	case "Bidder":
		bidders := new(Bidder).GetAllBidderSystemByUser(userId)
		for _, value := range bidders {
			object := payload.ObjectPage{
				ID:   value.Id,
				Name: value.BidderCode,
			}
			if value.BidderCode == "google" {
				object.Name = value.DisplayName
			}
			objects = append(objects, object)
		}
		break
	case "Channels":
		channels := new(Channels).GetByUser(userId)
		for _, value := range channels {
			object := payload.ObjectPage{
				ID:   value.Id,
				Name: value.Name,
			}
			objects = append(objects, object)
		}
		break
	case "Content":
		contents := new(Content).GetByUser(userId)
		for _, value := range contents {
			object := payload.ObjectPage{
				ID:   value.Id,
				Name: value.Title,
			}
			objects = append(objects, object)
		}
		break
	case "Demand":
		lineItems := new(LineItem).GetByUser(userId)
		for _, value := range lineItems {
			object := payload.ObjectPage{
				ID:   value.Id,
				Name: value.Name,
			}
			objects = append(objects, object)
		}
		break
	case "GAM":

		break
	case "Identity":
		Identitys := new(Identity).GetByUser(userId)
		for _, value := range Identitys {
			object := payload.ObjectPage{
				ID:   value.Id,
				Name: value.Name,
			}
			objects = append(objects, object)
		}
		break
	case "Playlist":
		Playlists := new(Playlist).GetByUser(userId)
		for _, value := range Playlists {
			object := payload.ObjectPage{
				ID:   value.Id,
				Name: value.Name,
			}
			objects = append(objects, object)
		}
		break
	case "Rule":
		Blockings := new(Blocking).GetByUser(userId)
		for _, value := range Blockings {
			object := payload.ObjectPage{
				ID:   value.Id,
				Name: value.RestrictionName,
			}
			objects = append(objects, object)
		}
		Floors := new(Floor).GetByUser(userId)
		for _, value := range Floors {
			object := payload.ObjectPage{
				ID:   value.Id,
				Name: value.Name,
			}
			objects = append(objects, object)
		}
		BlockedPages := new(BlockedPage).GetByUser(userId)
		for _, value := range BlockedPages {
			object := payload.ObjectPage{
				ID:   value.ID,
				Name: value.Name,
			}
			objects = append(objects, object)
		}
		break
	case "Supply":
		inventories := new(Inventory).GetByUser(userId)
		for _, value := range inventories {
			object := payload.ObjectPage{
				ID:   value.Id,
				Name: value.Name,
			}
			objects = append(objects, object)
		}
		break
	case "Template":
		templates := new(Inventory).GetByUser(userId)
		for _, value := range templates {
			object := payload.ObjectPage{
				ID:   value.Id,
				Name: value.Name,
			}
			objects = append(objects, object)
		}
		break
	case "User":

		break
	}
	return
}
