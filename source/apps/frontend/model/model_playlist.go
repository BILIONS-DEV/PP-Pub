package model

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"source/apps/frontend/lang"
	"source/apps/frontend/payload"
	"source/apps/history"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"source/pkg/datatable"
	"source/pkg/htmlblock"
	"source/pkg/pagination"
	"source/pkg/utility"
	"strconv"
	"strings"
)

type Playlist struct{}

type PlaylistRecord struct {
	mysql.TablePlaylist
}

func (PlaylistRecord) TableName() string {
	return mysql.Tables.Playlist
}

func (t *Playlist) GetById(id, userId int64) (record PlaylistRecord, err error) {
	if id == 13 {
		err = mysql.Client.Where("id = ?", id).Find(&record).Error
	} else {
		err = mysql.Client.Where("id = ? and user_id = ?", id, userId).Find(&record).Error
	}
	if record.Id == 0 {
		err = errors.New("Record not found")
		return
	}
	// Get rls
	record.GetRls()
	return
}

func (t *Playlist) GetByUser(userId int64) (records []PlaylistRecord) {
	mysql.Client.Where("user_id = ?", userId).Find(&records)
	return
}

func (t *Playlist) Create(inputs payload.PlaylistCreate, userId int64, userAdmin UserRecord) (record PlaylistRecord, errs []ajax.Error) {
	errs = t.ValidateCreate(inputs, userId)
	if len(errs) > 0 {
		return
	}
	record = t.makeRow(inputs, userId)
	err := mysql.Client.Create(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Translate.Errors.PlaylistError.Add.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}

	// for _, id := range inputs.Content {
	// 	idContent, err := strconv.ParseInt(id, 10, 64)
	// 	if err != nil {
	// 		//errs = append(errs, ajax.Error{
	// 		//	Id:      "",
	// 		//	Message: err.Error(),
	// 		//})
	// 		continue
	// 	}
	// 	err = new(RlPlaylistContent).Create(RlsPlaylistContentRecord{mysql.TableRlsPlaylistContent{
	// 		PlaylistId: record.Id,
	// 		ContentId:  idContent,
	// 	}})
	// 	if err != nil {
	// 		if !utility.IsWindow() {
	// 			errs = append(errs, ajax.Error{
	// 				Id:      "",
	// 				Message: lang.Errors.PlaylistError.RlContentPlaylist.ToString(),
	// 			})
	// 		} else {
	// 			errs = append(errs, ajax.Error{
	// 				Id:      "",
	// 				Message: err.Error(),
	// 			})
	// 		}
	// 	}
	// }

	// Tạo target cho line item vào bảng playlist_config
	err = t.CreatePlaylistConfig(record.Id, userId, inputs)
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

	listInventory := t.GetInventoryByPlaylist(record.Id)
	for _, id := range listInventory {
		new(Inventory).ResetCacheWorker(id)
	}
	//Push History
	recordNew, _ := new(Playlist).GetById(record.Id, userId)
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userId
	}
	_ = history.PushHistory(&history.Playlist{
		Detail:    history.DetailPlaylistFE,
		CreatorId: creatorId,
		RecordOld: mysql.TablePlaylist{},
		RecordNew: recordNew.TablePlaylist,
	})
	return
}

func (t *Playlist) Update(inputs payload.PlaylistCreate, userId int64, userAdmin UserRecord) (record PlaylistRecord, errs []ajax.Error) {
	flag := t.VerificationRecord(inputs.Id, userId)
	if !flag {
		errs = append(errs, ajax.Error{
			Id:      "id",
			Message: "You don't own this playlist",
		})
		return
	}
	oldRecord, _ := t.GetById(inputs.Id, userId)
	errs = t.ValidateEdit(inputs, oldRecord)
	if len(errs) > 0 {
		return
	}
	record = t.makeRowEdit(inputs, oldRecord)
	err := mysql.Client.Save(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Translate.Errors.PlaylistError.Edit.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}

	// Delete all table rl_playlist_content theo playlistId
	// new(RlPlaylistContent).DeleteAllByPlaylistId(record.Id)
	//
	// for _, id := range inputs.Content {
	// 	idContent, err := strconv.ParseInt(id, 10, 64)
	// 	if err != nil {
	// 		//errs = append(errs, ajax.Error{
	// 		//	Id:      "",
	// 		//	Message: err.Error(),
	// 		//})
	// 		continue
	// 	}
	// 	err = new(RlPlaylistContent).Create(RlsPlaylistContentRecord{mysql.TableRlsPlaylistContent{
	// 		PlaylistId: record.Id,
	// 		ContentId:  idContent,
	// 	}})
	// 	if err != nil {
	// 		if !utility.IsWindow() {
	// 			errs = append(errs, ajax.Error{
	// 				Id:      "",
	// 				Message: lang.Errors.PlaylistError.RlContentPlaylist.ToString(),
	// 			})
	// 		} else {
	// 			errs = append(errs, ajax.Error{
	// 				Id:      "",
	// 				Message: err.Error(),
	// 			})
	// 		}
	// 	}
	// }

	//Xóa toàn bộ target cũ để tạo mới list target nhận đc
	err = new(PlaylistConfig).DeletePlaylistConfig(PlaylistConfigRecord{mysql.TablePlaylistConfig{
		PlaylistId: record.Id,
	}})
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Translate.Errors.PlaylistError.PlaylistConfig.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}
	// Tạo target cho line item vào bảng playlist_config
	err = t.CreatePlaylistConfig(record.Id, userId, inputs)
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Translate.Errors.PlaylistError.PlaylistConfig.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}

	listInventory := t.GetInventoryByPlaylist(record.Id)
	for _, id := range listInventory {
		new(Inventory).ResetCacheWorker(id)
	}
	//Push History
	recordNew, _ := new(Playlist).GetById(record.Id, userId)
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userId
	}
	_ = history.PushHistory(&history.Playlist{
		Detail:    history.DetailPlaylistFE,
		CreatorId: creatorId,
		RecordOld: oldRecord.TablePlaylist,
		RecordNew: recordNew.TablePlaylist,
	})
	return
}

func (t *Playlist) makeRow(input payload.PlaylistCreate, userId int64) (record PlaylistRecord) {
	record.Id = input.Id
	record.Name = input.Name
	record.Description = strings.TrimSpace(input.Description)
	record.OrderingMethod = input.OrderingMethod
	record.VideosLimit = input.VideosLimit
	record.UserId = userId
	record.IsDefault = mysql.TypeOff
	return
}

func (t *Playlist) makeRowEdit(input payload.PlaylistCreate, oldRecord PlaylistRecord) (record PlaylistRecord) {
	record.Id = input.Id
	record.Name = input.Name
	record.Description = strings.TrimSpace(input.Description)
	record.OrderingMethod = input.OrderingMethod
	record.VideosLimit = input.VideosLimit
	record.UserId = oldRecord.UserId
	record.IsDefault = oldRecord.IsDefault
	return
}

func (t *Playlist) ValidateCreate(inputs payload.PlaylistCreate, userId int64) (errs []ajax.Error) {
	if utility.ValidateString(inputs.Name) == "" {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Name is required",
		})
	}
	if t.ValidateName(inputs, userId) == false {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Name already exist!",
		})
	}
	return
}

func (t *Playlist) ValidateEdit(inputs payload.PlaylistCreate, oldRecord PlaylistRecord) (errs []ajax.Error) {
	if oldRecord.IsDefault == mysql.TypeOn {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: "You can not edit this playlist",
		})
	}
	if utility.ValidateString(inputs.Name) == "" {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Name is required",
		})
	}
	if t.ValidateName(inputs, oldRecord.UserId) == false {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Name already exist!",
		})
	}
	return
}

func (t *Playlist) ValidateName(inputs payload.PlaylistCreate, userId int64) bool {
	var record PlaylistRecord
	query := mysql.Client.Where("name = ? and user_id = ?", inputs.Name, userId)
	if inputs.Id != 0 {
		query.Where("id != ?", inputs.Id)
	}
	err := query.Find(&record).Error
	if err != nil || record.Id != 0 {
		return false
	}
	return true
}

func (t *Playlist) GetByFilters(inputs *payload.PlaylistFilterPayload, userId int64) (response datatable.Response, err error) {
	lang := lang.Translation{}
	var records []PlaylistRecord
	var total int64
	err = mysql.Client.Where("id = 13 OR user_id = ?", userId).
		Scopes(
			// t.SetFilterStatus(inputs),
			t.setFilterSearch(inputs),
		).
		Model(&records).Count(&total).
		Scopes(
			t.setOrder(inputs),
			pagination.Paginate(pagination.Params{
				Limit:  inputs.Length,
				Offset: inputs.Start,
			}),
		).
		Find(&records).Error
	if err != nil {
		if !utility.IsWindow() {
			err = fmt.Errorf(lang.Errors.PlaylistError.List.ToString())
		}
		return datatable.Response{}, err
	}
	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data, err = t.MakeResponseDatatable(records, userId)
	return
}

func (t *Playlist) SetFilterStatus(inputs *payload.PlaylistFilterPayload) func(db *gorm.DB) *gorm.DB {
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

func (t *Playlist) setFilterSearch(inputs *payload.PlaylistFilterPayload) func(db *gorm.DB) *gorm.DB {
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

func (t *Playlist) setOrder(inputs *payload.PlaylistFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var orders []string
		orders = append(orders, "is_default ASC")
		if len(inputs.Order) > 0 {
			for _, order := range inputs.Order {
				column := inputs.Columns[order.Column]
				orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
			}
		}
		orderString := strings.Join(orders, ", ")
		return db.Order(orderString)
	}
}

type PlaylistRecordDatatable struct {
	PlaylistRecord
	RowId       string `json:"DT_RowId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Action      string `json:"action"`
}

func (t *Playlist) MakeResponseDatatable(playlists []PlaylistRecord, userId int64) (records []PlaylistRecordDatatable, err error) {
	for _, playlist := range playlists {
		rec := PlaylistRecordDatatable{
			PlaylistRecord: playlist,
			RowId:          "bidder_" + strconv.FormatInt(playlist.Id, 10),
			Name:           htmlblock.Render("playlist/index/block.name.gohtml", playlist).String(),
			Description:    htmlblock.Render("playlist/index/description.gohtml", playlist).String(),
			Action:         htmlblock.Render("playlist/index/block.action.gohtml", playlist).String(),
		}
		records = append(records, rec)
	}
	return
}

func (this *Playlist) Delete(id, userId int64, userAdmin UserRecord) fiber.Map {
	record, _ := new(Playlist).GetById(id, userId)
	err := mysql.Client.Model(&PlaylistRecord{}).Delete(&PlaylistRecord{}, "id = ? and user_id", id, userId).Error
	if err != nil {
		if !utility.IsWindow() {
			lang := lang.Translation{}
			return fiber.Map{
				"status":  "err",
				"message": lang.Errors.PlaylistError.Delete.ToString(),
				"id":      id,
			}
		}
		return fiber.Map{
			"status":  "err",
			"message": err.Error(),
			"id":      id,
		}
	} else {
		listInventory := this.GetInventoryByPlaylist(id)
		for _, inventoryId := range listInventory {
			new(Inventory).ResetCacheWorker(inventoryId)
		}

		// History
		var creatorId int64
		if userAdmin.Id != 0 {
			creatorId = userAdmin.Id
		} else {
			creatorId = userId
		}
		_ = history.PushHistory(&history.Playlist{
			Detail:    history.DetailPlaylistFE,
			CreatorId: creatorId,
			RecordOld: record.TablePlaylist,
			RecordNew: mysql.TablePlaylist{},
		})

		return fiber.Map{
			"status":  "success",
			"message": "done",
			"id":      id,
		}
	}
}

func (t *Playlist) GetAll(userId int64) (records []PlaylistRecord, err error) {
	err = mysql.Client.Where("is_default = 1 OR user_id = ?", userId).Find(&records).Error
	return
}

func (t *Playlist) VerificationRecord(id, userId int64) bool {
	var row PlaylistRecord
	err := mysql.Client.Model(&PlaylistRecord{}).Where("id = ? and user_id = ?", id, userId).Find(&row).Error
	if err != nil || row.Id == 0 {
		return false
	}
	return true
}

func (t *Playlist) GetInventoryByPlaylist(playlistId int64) (listInventoryId []int64) {
	mysql.Client.Model(&InventoryAdTagRecord{}).Select("inventory_id").Where("content_source = ? and playlist_id = ?", 1, playlistId).Find(&listInventoryId)
	return
}

func (t *Playlist) GetListBoxCollapse(userId, playlistId int64, page, typ string) (list []string) {
	switch typ {
	case "add":
		mysql.Client.Select("box_collapse").Model(PageCollapseRecord{}).Where("user_id = ? and page_collapse = ? and is_collapse = ? and page_type = ?", userId, page, 1, typ).Find(&list)
		return
	case "edit":
		mysql.Client.Select("box_collapse").Model(PageCollapseRecord{}).Where("user_id = ? and page_collapse = ? and is_collapse = ? and page_type = ? and page_id = ?", userId, page, 1, typ, playlistId).Find(&list)
		return
	}
	return
}

func (t *Playlist) CreatePlaylistConfig(PlaylistId, userId int64, inputs payload.PlaylistCreate) (err error) {
	all := int64(-1)
	// Kiểm tra nếu đầu vào input list target = 0 thì thêm một target = 0 thể hiện select all
	if len(inputs.ListCategory) == 0 {
		inputs.ListCategory = []payload.ListPlaylistConfig{
			{
				Id: all,
			},
		}
	}
	if len(inputs.ListChannels) == 0 {
		inputs.ListChannels = []payload.ListPlaylistConfig{
			{
				Id: all,
			},
		}
	}
	if len(inputs.ListKeywords) == 0 {
		inputs.ListKeywords = []payload.ListPlaylistConfig{
			{
				Id: all,
			},
		}
	}
	if len(inputs.ListLanguage) == 0 {
		inputs.ListLanguage = []payload.ListPlaylistConfig{
			{
				Id: all,
			},
		}
	}
	if len(inputs.ListVideos) == 0 {
		inputs.ListVideos = []payload.ListPlaylistConfig{
			{
				Id: all,
			},
		}
	}

	for _, category := range inputs.ListCategory {
		var recordPlaylistConfig PlaylistConfigRecord
		err := mysql.Client.Model(&PlaylistConfigRecord{}).FirstOrCreate(&recordPlaylistConfig, PlaylistConfigRecord{mysql.TablePlaylistConfig{
			UserId:     userId,
			PlaylistId: PlaylistId,
			CategoryId: category.Id,
			Type:       inputs.TypeCategory,
		}}).Error
		if err != nil {
			return err
		}
	}
	for _, channels := range inputs.ListChannels {
		var recordPlaylistConfig PlaylistConfigRecord
		err := mysql.Client.Model(&PlaylistConfigRecord{}).FirstOrCreate(&recordPlaylistConfig, PlaylistConfigRecord{mysql.TablePlaylistConfig{
			UserId:     userId,
			PlaylistId: PlaylistId,
			ChannelsId: channels.Id,
			Type:       inputs.TypeChannels,
		}}).Error
		if err != nil {
			return err
		}
	}
	for _, keyword := range inputs.ListKeywords {
		var recordPlaylistConfig PlaylistConfigRecord
		err := mysql.Client.Model(&PlaylistConfigRecord{}).FirstOrCreate(&recordPlaylistConfig, PlaylistConfigRecord{mysql.TablePlaylistConfig{
			UserId:           userId,
			PlaylistId:       PlaylistId,
			ContentKeywordId: keyword.Id,
			Type:             inputs.TypeKeywords,
		}}).Error
		if err != nil {
			return err
		}
	}

	for _, language := range inputs.ListLanguage {
		var recordPlaylistConfig PlaylistConfigRecord
		err := mysql.Client.Model(&PlaylistConfigRecord{}).FirstOrCreate(&recordPlaylistConfig, PlaylistConfigRecord{mysql.TablePlaylistConfig{
			UserId:     userId,
			PlaylistId: PlaylistId,
			LanguageId: language.Id,
			Type:       inputs.TypeLanguage,
		}}).Error
		if err != nil {
			return err
		}
	}
	for _, video := range inputs.ListVideos {
		var recordPlaylistConfig PlaylistConfigRecord
		err := mysql.Client.Model(&PlaylistConfigRecord{}).FirstOrCreate(&recordPlaylistConfig, PlaylistConfigRecord{mysql.TablePlaylistConfig{
			UserId:     userId,
			PlaylistId: PlaylistId,
			ContentId:  video.Id,
			Type:       inputs.TypeVideos,
		}}).Error
		if err != nil {
			return err
		}
	}
	return
}

func (rec *PlaylistRecord) GetVideoContentByPlaylist() (videos []ContentRecord, err error) {
	// Language
	var listIdLanguageTargetInclude []int64
	var listIdLanguageTargetExclude []int64
	// Channels
	var listIdChannelsTargetInclude []int64
	var listIdChannelsTargetExclude []int64
	// Category
	var listIdCategoryTargetInclude []int64
	var listIdCategoryTargetExclude []int64
	// Keyword
	var listIdKeywordsTargetInclude []int64
	var listIdKeywordsTargetExclude []int64
	// Video
	var listIdVideosTargetInclude []int64
	var listIdVideosTargetExclude []int64
	listPlaylistConfig := new(PlaylistConfig).GetByPlaylistId(rec.Id)
	for _, playlistConfig := range listPlaylistConfig {
		if playlistConfig.LanguageId != 0 {
			if playlistConfig.Type == 1 {
				listIdLanguageTargetInclude = append(listIdLanguageTargetInclude, playlistConfig.LanguageId)
			} else if playlistConfig.Type == 2 {
				listIdLanguageTargetExclude = append(listIdLanguageTargetExclude, playlistConfig.LanguageId)
			}
		}
		if playlistConfig.ChannelsId != 0 {
			if playlistConfig.Type == 1 {
				listIdChannelsTargetInclude = append(listIdChannelsTargetInclude, playlistConfig.ChannelsId)
			} else if playlistConfig.Type == 2 {
				listIdChannelsTargetExclude = append(listIdChannelsTargetExclude, playlistConfig.ChannelsId)
			}
		}
		if playlistConfig.CategoryId != 0 {
			if playlistConfig.Type == 1 {
				listIdCategoryTargetInclude = append(listIdCategoryTargetInclude, playlistConfig.CategoryId)
			} else if playlistConfig.Type == 2 {
				listIdCategoryTargetExclude = append(listIdCategoryTargetExclude, playlistConfig.CategoryId)
			}
		}
		if playlistConfig.ContentKeywordId != 0 {
			if playlistConfig.Type == 1 {
				listIdKeywordsTargetInclude = append(listIdKeywordsTargetInclude, playlistConfig.ContentKeywordId)
			} else if playlistConfig.Type == 2 {
				listIdKeywordsTargetExclude = append(listIdKeywordsTargetExclude, playlistConfig.ContentKeywordId)
			}
		}
		if playlistConfig.ContentId != 0 {
			if playlistConfig.Type == 1 {
				listIdVideosTargetInclude = append(listIdVideosTargetInclude, playlistConfig.ContentId)
			} else if playlistConfig.Type == 2 {
				listIdVideosTargetExclude = append(listIdVideosTargetExclude, playlistConfig.ContentId)
			}
		}
	}

	// Trong trường hợp exclude all thì mặc định của các target là k có video nào phù hợp
	if len(listIdLanguageTargetExclude) == 1 {
		if listIdLanguageTargetExclude[0] == -1 { // -1 là đánh dấu all
			return
		}
	}
	if len(listIdChannelsTargetExclude) == 1 {
		if listIdChannelsTargetExclude[0] == -1 { // -1 là đánh dấu all
			return
		}
	}
	if len(listIdCategoryTargetExclude) == 1 {
		if listIdCategoryTargetExclude[0] == -1 { // -1 là đánh dấu all
			return
		}
	}
	if len(listIdKeywordsTargetExclude) == 1 {
		if listIdKeywordsTargetExclude[0] == -1 { // -1 là đánh dấu all
			return
		}
	}
	if len(listIdVideosTargetExclude) == 1 {
		if listIdVideosTargetExclude[0] == -1 { // -1 là đánh dấu all
			return
		}
	}
	// Từ Language + Channels + Category tìm ra các channel phù hợp
	var listChanelId []int64
	err = mysql.Client.Select("id").Model(&ChannelsRecord{}).Where("user_id = ?", rec.UserId).
		Scopes(
			// Xử lý cho language
			rec.filterIn("language", listIdLanguageTargetInclude),
			rec.filterNotIn("language", listIdLanguageTargetExclude),
			// Xử lý cho channels
			rec.filterIn("id", listIdChannelsTargetInclude),
			rec.filterNotIn("id", listIdChannelsTargetExclude),
			// Xử lý cho category
			rec.filterIn("category", listIdCategoryTargetInclude),
			rec.filterNotIn("category", listIdCategoryTargetExclude),
		).
		Find(&listChanelId).
		Error
	if err != nil {
		return
	}
	//fmt.Printf("%+v \n", listChanelId)
	// Từ list channelId tìm ra listVideoId phù hợp
	var listVideoIdFromChannel []int64
	mysql.Client.Select("id").Model(&ContentRecord{}).Where("channels in ?", listChanelId).Find(&listVideoIdFromChannel)
	if len(listVideoIdFromChannel) == 0 { // Nếu listId video rỗng tức không có video nào phù hợp target return luôn
		return
	}
	//fmt.Printf("%+v \n", listVideoIdFromChannel)

	// Tìm tất cả các video không có keyword mặc định hiểu là keyword all
	// Tìm các id video của user
	var listVideoIdByUser []int64
	mysql.Client.Model(&ContentRecord{}).Select("id").Where("user_id = ?", rec.UserId).Find(&listVideoIdByUser)
	//fmt.Printf("%+v \n", listVideoIdByUser)
	// Tìm các id video có keyword
	var listVideoIdHaveKeyword []int64
	mysql.Client.Select("content_id").Model(&ContentKeywordRecord{}).Where("user_id = ?", rec.UserId).
		Group("content_id").
		Find(&listVideoIdHaveKeyword)
	//fmt.Printf("%+v \n", listVideoIdHaveKeyword)
	var listVideoAllKeyword []int64
	for _, videoId := range listVideoIdByUser {
		if utility.InArray(videoId, listVideoIdHaveKeyword, false) { // Loại các video có keyword
			continue
		}
		listVideoAllKeyword = append(listVideoAllKeyword, videoId)
	}
	//fmt.Printf("%+v \n", listVideoAllKeyword)

	var listVideoIdFromKeyword []int64
	// Trong trường hợp target exclude lớn hơn không thì tiến hành tìm các videoId đã loại từ keyword
	if len(listIdKeywordsTargetExclude) > 0 {
		var listVideoIdFromKeyWordExclude []int64
		var listKeywords []string
		// Từ listContentKeywordId trong target tìm ra các keyword
		err = mysql.Client.Select("keyword").Model(&ContentKeywordRecord{}).Where("user_id = ?", rec.UserId).
			Scopes(
				// Xử lý cho keyword
				rec.filterIn("id", listIdKeywordsTargetExclude),
			).
			Group("keyword").
			Find(&listKeywords).
			Error
		if err != nil {
			return
		}
		//fmt.Println(listKeywords)
		// Từ các keyword tìm ra các video exclude tương ứng
		err = mysql.Client.Select("content_id").Model(&ContentKeywordRecord{}).Where("user_id = ?", rec.UserId).
			Scopes(
				// Xử lý cho keyword
				rec.filterInListKeyword("keyword", listKeywords),
			).
			Group("content_id").
			Find(&listVideoIdFromKeyWordExclude).
			Error
		if err != nil {
			return
		}
		// Từ các video exclude tìm các video có keyword còn lại
		var listVideoId []int64
		err = mysql.Client.Select("content_id").Model(&ContentKeywordRecord{}).Where("user_id = ?", rec.UserId).
			Scopes(
				// Xử lý cho keyword
				rec.filterNotIn("content_id", listVideoIdFromKeyWordExclude),
			).
			Group("content_id").
			Find(&listVideoId).
			Error
		if err != nil {
			return
		}
		//fmt.Println(listVideoId)
		listVideoIdFromKeyword = append(listVideoIdFromKeyword, listVideoId...)
	}
	//fmt.Printf("List video from key word after exclude : %+v \n", listVideoIdFromKeyword)

	// Từ list keywordId trong target này tìm ra các keywords => rồi từ các keyword tìm ra các video phù hợp
	if len(listIdKeywordsTargetInclude) > 0 {
		var listVideoIdFromKeyWordInclude []int64
		var listKeywords []string
		// Từ list keywordId trong target này tìm ra các keywords => rồi từ các keyword tìm ra các video phù hợp
		err = mysql.Client.Select("keyword").Model(&ContentKeywordRecord{}).Where("user_id = ?", rec.UserId).
			Scopes(
				// Xử lý cho keyword
				rec.filterIn("id", listIdKeywordsTargetInclude),
			).
			Group("keyword").
			Find(&listKeywords).
			Error
		//fmt.Println(listKeywords)
		err = mysql.Client.Select("content_id").Model(&ContentKeywordRecord{}).Where("user_id = ?", rec.UserId).
			Scopes(
				// Xử lý cho keyword
				rec.filterInListKeyword("keyword", listKeywords),
			).
			Group("content_id").
			Find(&listVideoIdFromKeyWordInclude).
			Error
		if err != nil {
			return
		}
		listVideoIdFromKeyword = append(listVideoIdFromKeyword, listVideoIdFromKeyWordInclude...)
	}
	//fmt.Printf("List video from key word after Include : %+v \n", listVideoIdFromKeyword)

	// Nối listVideoKeyword với listVideo không có keyword mặc định là all
	listVideoIdFromKeyword = append(listVideoIdFromKeyword, listVideoAllKeyword...)
	//fmt.Printf("ListVideoIdFromKeyword: %+v \n", listVideoIdFromKeyword)

	// Từ các list video trên tìm ra các video
	err = mysql.Client.Model(&ContentRecord{}).Where("user_id = ?", rec.UserId).
		Scopes(
			//Channel
			rec.filterIn("id", listVideoIdFromChannel),

			//Keyword
			rec.filterIn("id", listVideoIdFromKeyword),

			//Video
			rec.filterIn("id", listIdVideosTargetInclude),
			rec.filterNotIn("id", listIdVideosTargetExclude),
		).
		Find(&videos).
		Error
	if err != nil {
		return
	}
	//fmt.Printf("%+v \n", videos)
	return
}

func (rec *PlaylistRecord) filterIn(columnName string, listIn []int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(listIn) == 0 {
			return db
		}
		if len(listIn) == 1 {
			if listIn[0] == -1 {
				return db
			}
		}
		return db.Where(columnName+" in ?", listIn)
	}
}

func (rec *PlaylistRecord) filterNotIn(columnName string, listIn []int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(listIn) == 0 {
			return db
		}
		return db.Where(columnName+" not in ?", listIn)
	}
}

func (rec *PlaylistRecord) filterInListKeyword(columnName string, keywords []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(keywords) == 0 {
			return db
		}
		return db.Where(columnName+" in ?", keywords)
	}
}

func (rec *PlaylistRecord) filterNotInListKeyword(columnName string, keywords []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(keywords) == 0 {
			return db
		}
		return db.Where(columnName+" not in ?", keywords)
	}
}
