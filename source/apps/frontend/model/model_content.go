package model

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
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

type Content struct{}

type ContentRecord struct {
	mysql.TableContent
}

func (ContentRecord) TableName() string {
	return mysql.Tables.Content
}

func (t *Content) Create(inputs payload.ContentCreate, user UserRecord, userAdmin UserRecord, lang lang.Translation) (record ContentRecord, errs []ajax.Error) {
	errs = t.ValidateCreate(inputs)
	if len(errs) > 0 {
		return
	}
	record = t.makeRowCreate(inputs)
	record.UserId = user.Id
	err := mysql.Client.Create(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.ContentError.Add.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}

	// Tạo tag trong bảng content_keyword
	for _, v := range inputs.Keyword {
		mysql.Client.Create(&ContentKeywordRecord{mysql.TableContentKeyword{
			ContentId: record.Id,
			UserId:    user.Id,
			Keyword:   v,
		}})
	}
	// Tạo ad break trong bảng content_ad_break
	for _, v := range inputs.AdBreaks {
		mysql.Client.Create(&ContentAdBreakRecord{mysql.TableContentAdBreak{
			ContentId:   record.Id,
			Type:        v.Type,
			TimeAdBreak: v.TimeBreak,
			BreakMode:   v.BreakMode,
		}})
	}

	list := t.GetPlaylistByContent(record.Id)
	for _, playlist := range list {
		if playlist.Type != mysql.TYPEDisplay {
			listInventory := new(Player).GetListInventoryByType(playlist.Type, user.Id)
			t.UpdateInventory(listInventory)
		}
	}

	listInventoryRecord := new(Inventory).GetByUser(user.Id)
	for _, inventory := range listInventoryRecord {
		new(Inventory).ResetCacheWorker(inventory.Id)
	}

	// Push History
	recordNew, _ := new(Content).GetById(record.Id, user.Id)
	// History
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = user.Id
	}
	_ = history.PushHistory(&history.Content{
		Detail:    history.DetailContentFE,
		CreatorId: creatorId,
		RecordOld: mysql.TableContent{},
		RecordNew: recordNew.TableContent,
	})
	return
}

func (t *Content) makeRowCreate(row payload.ContentCreate) (record ContentRecord) {
	record.Id = row.Id
	record.Uuid = utility.RandSeq(9)
	record.Title = row.Title
	record.Type = 1
	record.ContentDesc = row.ContentDesc
	record.Thumb = row.Thumb
	record.VideoUrl = row.VideoUrl
	record.VideoType = row.VideoType
	record.Duration = row.Duration
	// record.Category = row.Category
	record.Channels = row.Channels
	record.NameFile = row.NameFile
	record.Status = mysql.StatusApproved
	record.VideoType = 1 // public
	record.ConfigAdBreak = row.ConfigAdBreak
	return
}

func (t *Content) makeRowUpdate(row payload.ContentCreate, oldRecord ContentRecord) (record ContentRecord) {
	record.Id = row.Id
	record.Title = row.Title
	record.Type = 1
	record.ContentDesc = row.ContentDesc
	record.Thumb = row.Thumb
	record.VideoUrl = row.VideoUrl
	record.Duration = row.Duration
	// record.VideoType = row.VideoType
	// record.Category = row.Category
	record.Channels = row.Channels
	record.NameFile = row.NameFile
	record.Status = oldRecord.Status
	record.VideoType = 1 // public
	record.ConfigAdBreak = row.ConfigAdBreak
	return
}

func (t *Content) ValidateCreate(inputs payload.ContentCreate) (errs []ajax.Error) {
	if utility.ValidateString(inputs.Title) == "" {
		errs = append(errs, ajax.Error{
			Id:      "title",
			Message: "Title is required",
		})
	}

	if inputs.Channels == 0 {
		errs = append(errs, ajax.Error{
			Id:      "channels",
			Message: "Channels is required",
		})
	}

	if utility.ValidateString(inputs.VideoUrl) == "" {
		errs = append(errs, ajax.Error{
			Id:      "box-upload-video",
			Message: "Video Url is required",
		})
	}

	if utility.ValidateString(inputs.Thumb) == "" {
		errs = append(errs, ajax.Error{
			Id:      "box-upload-thumb",
			Message: "Thumb is required",
		})
	}
	if len(inputs.AdStartTime) > 0 {
		for i, adStartTime := range inputs.AdStartTime {
			id := "ad_start_time_" + strconv.Itoa(i+1)
			if govalidator.IsNull(adStartTime) {
				errs = append(errs, ajax.Error{
					Id:      id,
					Message: "(*) required",
				})
			}
			// else {
			//	r := regexp.MustCompile(`^(([0-9])|([0-5][0-9])):([0-9]|[0-5][0-9])$`)
			//	if !r.MatchString(adStartTime) {
			//		errs = append(errs, ajax.Error{
			//			Id:      id,
			//			Message: `Ad Start Time format is mm:ss - ex: 12:50`,
			//		})
			//	}
			// }
		}
	}
	return
}

func (t *Content) GetByFilters(inputs *payload.ContentFilterPayload, userId int64, lang lang.Translation) (response datatable.Response, err error) {
	var records []ContentRecord
	var total int64
	err = mysql.Client.Where("user_id = ?", userId).
		Scopes(
			// t.setFilterVideoType(inputs),
			t.setFilterType(inputs),
			t.setFilterSearch(inputs),
			// t.setFilterCategory(inputs),
		).
		Model(&records).Count(&total).
		Scopes(
			t.setOrder(inputs),
			pagination.Paginate(pagination.Params{
				Limit:  inputs.Length,
				Offset: inputs.Start,
			}),
		).Find(&records).Error
	if err != nil {
		if !utility.IsWindow() {
			err = fmt.Errorf(lang.Errors.ContentError.List.ToString())
		}
		return datatable.Response{}, err
	}
	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.MakeResponseDatatable(records)
	return
}

func (t *Content) setFilterVideoType(inputs *payload.ContentFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.VideoType != nil {
			switch inputs.PostData.VideoType.(type) {
			case string, int:
				if inputs.PostData.VideoType != "" {
					return db.Where("video_type = ?", inputs.PostData.VideoType)
				}
			case []string, []interface{}:
				return db.Where("video_type IN ?", inputs.PostData.VideoType)
			}
		}
		return db
	}
}

func (t *Content) setFilterType(inputs *payload.ContentFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.Type != nil {
			switch inputs.PostData.Type.(type) {
			case string, int:
				if inputs.PostData.VideoType != "" {
					return db.Where("type = ?", inputs.PostData.Type)
				}
			case []string, []interface{}:
				return db.Where("type IN ?", inputs.PostData.Type)
			}
		}
		return db
	}
}

func (t *Content) setFilterCategory(inputs *payload.ContentFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.Category != nil {
			switch inputs.PostData.Category.(type) {
			case string, int:
				if inputs.PostData.Category != "" {
					return db.Where("category = ?", inputs.PostData.Category)
				}
			case []string, []interface{}:
				return db.Where("category IN ?", inputs.PostData.Category)
			}
		}
		return db
	}
}

func (t *Content) setFilterSearch(inputs *payload.ContentFilterPayload) func(db *gorm.DB) *gorm.DB {
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
		return db.Where("title LIKE ?", "%"+inputs.PostData.QuerySearch+"%")
	}
}

func (t *Content) setOrder(inputs *payload.ContentFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(inputs.Order) > 0 {
			var orders []string
			for _, order := range inputs.Order {
				column := inputs.Columns[order.Column]
				orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
			}
			orderString := strings.Join(orders, ", ")
			return db.Order(orderString)
		}
		return db.Order("id DESC")
	}
}

type ContentRecordDatatable struct {
	ContentRecord
	RowId       string `json:"DT_RowId"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Description string `json:"content_desc"`
	Thumb       string `json:"thumb"`
	VideoUrl    string `json:"video_url"`
	VideoType   string `json:"video_type"`
	Category    string `json:"category"`
	Languague   string `json:"language"`
	Tags        string `json:"tags"`
	Action      string `json:"action"`
}

func (t *Content) MakeResponseDatatable(contents []ContentRecord) (records []ContentRecordDatatable) {
	for _, content := range contents {
		content.Title = truncateText(content.Title, 60)
		content.ContentDesc = truncateText(content.ContentDesc, 20)
		if strings.Contains(content.Thumb, "http") == false && strings.Contains(content.Thumb, "ul.pubpowerplatform") == false && strings.Contains(content.Thumb, "//") == false {
			content.Thumb = "//ul.pubpowerplatform.io" + content.Thumb
		}
		rec := ContentRecordDatatable{
			ContentRecord: content,
			RowId:         "bidder_" + strconv.FormatInt(content.Id, 10),
			Title:         htmlblock.Render("content/index/block.title.gohtml", content).String(),
			Description:   htmlblock.Render("content/index/block.description.gohtml", content).String(),
			Thumb:         "<a href=\"" + content.Thumb + "\" title=\"" + content.Thumb + "\" target=\"_blank\">Click here </a>",
			VideoUrl:      "<a href=\"" + content.VideoUrl + "\" title=\"" + content.VideoUrl + "\" target=\"_blank\">Click here </a>",
			// Category:      htmlblock.Render("content/index/block.category.gohtml", content).String(),
			// Tags:          new(InventoryAdTag).GetById(content.Tag).Name,
			// VideoType:     htmlblock.Render("content/index/block.video_type.gohtml", content).String(),
			Action: htmlblock.Render("content/index/block.action.gohtml", content).String(),
		}
		if content.Type == 1 {
			rec.Type = "Video"
		} else {
			rec.Type = "Quiz"
		}
		// category := new(Category).GetById(content.Category)
		// rec.Category = category.Name

		// language := new(Language).GetByLanguageCode(content.Language)
		// rec.Language = language.LanguageName

		// tags := new(ContentTag).GetByContent(content.Id)
		// for _, tag := range tags {
		// 	if rec.Tags == "" {
		// 		rec.Tags = tag.Tag
		// 	}else{
		// 		rec.Tags = rec.Tags + ", " + tag.Tag
		// 	}
		// }
		records = append(records, rec)
	}
	return
}

func (t *Content) HandleContents(contents []ContentRecord) (records []ContentRecord) {
	for _, content := range contents {
		content.ContentDesc = truncateText(content.ContentDesc, 20)
		if strings.Contains(content.Thumb, "http") == false && strings.Contains(content.Thumb, "ul.pubpowerplatform") == false && strings.Contains(content.Thumb, "//") == false {
			content.Thumb = "//ul.pubpowerplatform.io" + content.Thumb
		}

		records = append(records, content)
	}
	return
}

func (t *Content) GetById(id, userId int64) (record ContentRecord, err error) {
	err = mysql.Client.Where("id = ? and user_id = ?", id, userId).Find(&record).Error
	// Get rls
	record.GetRls()
	return
}

func (t *Content) GetByUser(userId int64) (records []ContentRecord) {
	mysql.Client.Where("status = 1 and user_id = ?", userId).Find(&records)
	return
}

func (t *Content) Update(inputs payload.ContentCreate, user UserRecord, userAdmin UserRecord, lang lang.Translation) (record ContentRecord, errs []ajax.Error) {
	oldContent, flag := t.VerificationRecord(inputs.Id, user.Id)
	if !flag {
		errs = append(errs, ajax.Error{
			Id:      "id",
			Message: "You don't own this content",
		})
		return
	}
	errs = t.ValidateCreate(inputs)
	if len(errs) > 0 {
		return
	}
	record = t.makeRowUpdate(inputs, oldContent)
	record.UserId = user.Id
	err := mysql.Client.Save(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.ContentError.Edit.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
	}

	// Xóa các Keyword cũ để add lại tag mới
	new(ContentKeyword).DeleteKeywordByContent(record.Id)
	// Tạo Keyword trong bảng content_keyword
	for _, v := range inputs.Keyword {
		mysql.Client.Create(&ContentKeywordRecord{mysql.TableContentKeyword{
			ContentId: record.Id,
			UserId:    user.Id,
			Keyword:   v,
		}})
	}

	// Xóa các ad break cũ để add lại ad break mới
	new(ContentAdBreak).DeleteAllAdBreakByContent(record.Id)
	// Tạo ad break trong bảng content_ad_break
	for _, v := range inputs.AdBreaks {
		mysql.Client.Create(&ContentAdBreakRecord{mysql.TableContentAdBreak{
			ContentId:   record.Id,
			Type:        v.Type,
			TimeAdBreak: v.TimeBreak,
			BreakMode:   v.BreakMode,
		}})
	}

	list := t.GetPlaylistByContent(record.Id)
	for _, playlist := range list {
		if playlist.Type != mysql.TYPEDisplay {
			listInventory := new(Player).GetListInventoryByType(playlist.Type, user.Id)
			t.UpdateInventory(listInventory)
		}
	}

	listInventoryRecord := new(Inventory).GetByUser(user.Id)
	for _, inventory := range listInventoryRecord {
		new(Inventory).ResetCacheWorker(inventory.Id)
	}
	// Push History
	recordNew, _ := new(Content).GetById(record.Id, user.Id)
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = user.Id
	}
	_ = history.PushHistory(&history.Content{
		Detail:    history.DetailContentFE,
		CreatorId: creatorId,
		RecordOld: oldContent.TableContent,
		RecordNew: recordNew.TableContent,
	})
	return
}

func (t *Content) CreateQuiz(inputs payload.ContentQuizCreate, user UserRecord, lang lang.Translation) (record ContentRecord, errs []ajax.Error) {
	errs = t.ValidateCreateQuiz(inputs)
	if len(errs) > 0 {
		return
	}
	record = t.makeRowQuizCreate(inputs)
	record.UserId = user.Id
	err := mysql.Client.Create(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.ContentError.Add.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}

	// Tạo tag trong bảng content_tag
	for _, v := range inputs.Tag {
		mysql.Client.Create(&ContentTagRecord{mysql.TableContentTag{
			ContentId: record.Id,
			Tag:       v,
		}})
	}
	// Tạo quiz trong bảng content_quiz
	for _, question := range inputs.Questions {
		answers, _ := json.Marshal(question.Answer)
		mysql.Client.Create(&ContentQuestionRecord{mysql.TableContentQuestion{
			ContentId:      record.Id,
			Title:          question.Title,
			Type:           question.Type,
			BackgroundType: question.BackgroundType,
			Background:     question.Background,
			PictureType:    question.PictureType,
			Picture:        question.Picture,
			Answers:        string(answers),
		}})
	}

	list := t.GetPlaylistByContent(record.Id)
	for _, playlist := range list {
		if playlist.Type != mysql.TYPEDisplay {
			listInventory := new(Player).GetListInventoryByType(playlist.Type, user.Id)
			t.UpdateInventory(listInventory)
		}
	}
	return
}

func (t *Content) UpdateQuiz(inputs payload.ContentQuizCreate, user UserRecord, lang lang.Translation) (record ContentRecord, errs []ajax.Error) {
	oldContent, flag := t.VerificationRecord(inputs.Id, user.Id)
	if !flag {
		errs = append(errs, ajax.Error{
			Id:      "id",
			Message: "You don't own this content",
		})
		return
	}
	errs = t.ValidateCreateQuiz(inputs)
	if len(errs) > 0 {
		return
	}
	record = t.makeRowQuizUpdate(inputs, oldContent)
	record.UserId = user.Id
	err := mysql.Client.Updates(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.ContentError.Edit.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
	}

	// Xóa các tag cũ để add lại tag mới
	new(ContentTag).DeleteTagByContent(record.Id)
	// Tạo tag trong bảng content_tag
	for _, v := range inputs.Tag {
		mysql.Client.Create(&ContentTagRecord{mysql.TableContentTag{
			ContentId: record.Id,
			Tag:       v,
		}})
	}

	// remove quiz trong bảng content_quiz
	deleteQuestion := []int64{}
	for _, question := range inputs.Questions {
		deleteQuestion = append(deleteQuestion, question.Id)
	}
	if deleteQuestion != nil {
		mysql.Client.Where("content_id = ? AND id not in ?", inputs.Id, deleteQuestion).Delete(ContentQuestionRecord{})
	}

	// update/create quiz trong bảng content_quiz
	for _, question := range inputs.Questions {
		if question.Id != 0 {
			answers, _ := json.Marshal(question.Answer)
			mysql.Client.Where("id = ?", question.Id).Updates(&ContentQuestionRecord{mysql.TableContentQuestion{
				ContentId:      record.Id,
				Title:          question.Title,
				Type:           question.Type,
				BackgroundType: question.BackgroundType,
				Background:     question.Background,
				PictureType:    question.PictureType,
				Picture:        question.Picture,
				Answers:        string(answers),
			}})
		} else {
			answers, _ := json.Marshal(question.Answer)
			mysql.Client.Create(&ContentQuestionRecord{mysql.TableContentQuestion{
				ContentId:      record.Id,
				Title:          question.Title,
				Type:           question.Type,
				BackgroundType: question.BackgroundType,
				Background:     question.Background,
				PictureType:    question.PictureType,
				Picture:        question.Picture,
				Answers:        string(answers),
			}})
		}
	}

	list := t.GetPlaylistByContent(record.Id)
	for _, playlist := range list {
		if playlist.Type != mysql.TYPEDisplay {
			listInventory := new(Player).GetListInventoryByType(playlist.Type, user.Id)
			t.UpdateInventory(listInventory)
		}
	}
	return
}

// func (t *Content) UpdateQuiz(inputs payload.ContentQuizCreate, user UserRecord, lang lang.Translation) (record ContentRecord, errs []ajax.Error) {
// 	flag := t.VerificationRecord(inputs.Id, user.Id)
// 	if !flag {
// 		errs = append(errs, ajax.Error{
// 			Id:      "id",
// 			Message: "You don't own this content",
// 		})
// 		return
// 	}
// 	errs = t.ValidateCreateQuiz(inputs)
// 	if len(errs) > 0 {
// 		return
// 	}
// 	record = t.makeRow(inputs)
// 	record.UserId = user.Id
// 	err := mysql.Client.Updates(&record).Error
// 	if err != nil {
// 		if !utility.IsWindow() {
// 			errs = append(errs, ajax.Error{
// 				Id:      "",
// 				Message: lang.Errors.ContentError.Edit.ToString(),
// 			})
// 		} else {
// 			errs = append(errs, ajax.Error{
// 				Id:      "",
// 				Message: err.Error(),
// 			})
// 		}
// 	}
//
// 	//Xóa các tag cũ để add lại tag mới
// 	new(ContentTag).DeleteTagByContent(record.Id)
// 	//Tạo tag trong bảng content_tag
// 	for _, v := range inputs.Tag {
// 		mysql.Client.Create(&ContentTagRecord{mysql.TableContentTag{
// 			ContentId: record.Id,
// 			Tag:       v,
// 		}})
// 	}
//
// 	//Xóa các ad break cũ để add lại ad break mới
// 	new(ContentAdBreak).DeleteAllAdBreakByContent(record.Id)
// 	//Tạo ad break trong bảng content_ad_break
// 	for _, v := range inputs.AdStartTime {
// 		mysql.Client.Create(&ContentAdBreakRecord{mysql.TableContentAdBreak{
// 			ContentId:   record.Id,
// 			TimeAdBreak: v,
// 		}})
// 	}
//
// 	list := t.GetPlaylistByContent(record.Id)
// 	for _, playlist := range list {
// 		if playlist.Type != mysql.TYPEDisplay {
// 			listInventory := new(Player).GetListInventoryByType(playlist.Type, user.Id)
// 			t.UpdateInventory(listInventory)
// 		}
// 	}
// 	return
// }

func (t *Content) ValidateCreateQuiz(inputs payload.ContentQuizCreate) (errs []ajax.Error) {
	if utility.ValidateString(inputs.Title) == "" {
		errs = append(errs, ajax.Error{
			Id:      "title",
			Message: "Title is required",
		})
	}

	if len(inputs.Questions) == 0 {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: "Question is required",
		})
	}

	// var question = payload.Questions{}
	for i, question := range inputs.Questions {
		number := strconv.Itoa(i + 1)
		if utility.ValidateString(question.Title) == "" {
			errs = append(errs, ajax.Error{
				Id:      "question" + number + "_title",
				Message: "Title question is required",
			})
		}
		if utility.ValidateString(question.Background) == "" {
			id := ""
			switch question.BackgroundType {
			case 1:
				id = "question" + number + "_background_color"
				break
			case 2:
				id = "question" + number + "_background_upload"
				break
			case 3:
				id = "question" + number + "_background_link"
				break
			}
			errs = append(errs, ajax.Error{
				Id:      id,
				Message: "Backround is required",
			})
		}
		if utility.ValidateString(question.Picture) == "" && question.Type == 2 {
			id := ""
			switch question.PictureType {
			case 1:
				id = "question" + number + "_picture_color"
				break
			case 2:
				id = "question" + number + "_picture_upload"
				break
			case 3:
				id = "question" + number + "_picture_link"
				break
			}
			errs = append(errs, ajax.Error{
				Id:      id,
				Message: "Picture is required",
			})
		}
		if question.BackgroundType == 3 {
			isUrl := utility.IsUrl(question.Background)
			if !isUrl {
				errs = append(errs, ajax.Error{
					Id:      "question" + number + "_background_link",
					Message: "Requested URL is not valid",
				})
			}
		}
		if question.Type == 2 && question.PictureType == 3 {
			isUrl := utility.IsUrl(question.Picture)
			if !isUrl {
				errs = append(errs, ajax.Error{
					Id:      "question" + number + "_picture_link",
					Message: "Requested URL is not valid",
				})
			}
		}
		if len(question.Answer) < 2 || len(question.Answer) > 4 {
			errs = append(errs, ajax.Error{
				Id:      "question" + number + "_answer",
				Message: "Please create 2-4 options of answer",
			})
		}
	}

	// if len(inputs.AdStartTime) > 0 {
	// 	for i, adStartTime := range inputs.AdStartTime {
	// 		id := "ad_start_time_" + strconv.Itoa(i+1)
	// 		if govalidator.IsNull(adStartTime) {
	// 			errs = append(errs, ajax.Error{
	// 				Id:      id,
	// 				Message: "(*) required",
	// 			})
	// 		}
	// 	}
	// }
	return
}

func (t *Content) makeRowQuizCreate(row payload.ContentQuizCreate) (record ContentRecord) {
	record.Id = row.Id
	record.Title = row.Title
	record.Type = 2
	record.Category = row.Category
	record.Status = mysql.StatusApproved
	// record.Language = row.Language
	// record.Tag = row.Tag
	return
}
func (t *Content) makeRowQuizUpdate(row payload.ContentQuizCreate, oldRecord ContentRecord) (record ContentRecord) {
	record.Id = row.Id
	record.Title = row.Title
	record.Type = 2
	record.Category = row.Category
	record.Status = oldRecord.Status
	// record.Language = row.Language
	// record.Tag = row.Tag
	return
}

func (this *Content) Delete(id, userId int64, userAdmin UserRecord, lang lang.Translation) fiber.Map {
	record, _ := new(Content).GetById(id, userId)
	err := mysql.Client.Model(&ContentRecord{}).Delete(&ContentRecord{}, "id = ? and user_id = ?", id, userId).Error
	if err != nil {
		if !utility.IsWindow() {
			return fiber.Map{
				"status":  "err",
				"message": lang.Errors.ContentError.Delete.ToString(),
				"id":      id,
			}
		}
		return fiber.Map{
			"status":  "err",
			"message": err.Error(),
			"id":      id,
		}
	} else {
		list := this.GetPlaylistByContent(id)
		for _, playlist := range list {
			if playlist.Type != mysql.TYPEDisplay {
				listInventory := new(Player).GetListInventoryByType(playlist.Type, userId)
				this.UpdateInventory(listInventory)
			}
		}

		listInventoryRecord := new(Inventory).GetByUser(userId)
		for _, inventory := range listInventoryRecord {
			new(Inventory).ResetCacheWorker(inventory.Id)
		}
		// History
		var creatorId int64
		if userAdmin.Id != 0 {
			creatorId = userAdmin.Id
		} else {
			creatorId = userId
		}
		_ = history.PushHistory(&history.Content{
			Detail:    history.DetailContentFE,
			CreatorId: creatorId,
			RecordOld: record.TableContent,
			RecordNew: mysql.TableContent{},
		})

		new(RlPlaylistContent).DeleteAllByContentId(id)
		return fiber.Map{
			"status":  "success",
			"message": "done",
			"id":      id,
		}
	}
}

func (t *Content) GetAll(userId int64) (records []ContentRecord) {
	mysql.Client.Where("user_id = ?", userId).Find(&records)
	return
}

func (t *Content) GetAllVideo(userId int64) (records []ContentRecord) {
	mysql.Client.Where("type = 1 and user_id = ?", userId).Find(&records)
	return
}

func truncateText(str string, max int) string {
	length := 0
	for i, _ := range str {
		length++
		if length >= max {
			return str[:i] + "..."
		}
	}
	return str
}

func (t *Content) VerificationRecord(id, userId int64) (row ContentRecord, flag bool) {
	err := mysql.Client.Model(&ContentRecord{}).Where("id = ? and user_id = ?", id, userId).Find(&row).Error
	if err != nil || row.Id == 0 {
		flag = false
		return
	}

	flag = true
	// Get rls
	row.GetRls()
	return
}

func (t *Content) GetPlaylistByContent(contentId int64) (rows []PlayerRecord) {
	var playlist []int64
	mysql.Client.Model(&RlsPlaylistContentRecord{}).Select("playlist_id").Where("content_id = ?", contentId).Find(&playlist)
	for _, id := range playlist {
		row := new(Player).GetDetail(id)
		if row.Id != 0 {
			rows = append(rows, row)
		}
	}
	return
}

func (t *Content) UpdateInventory(listInventory []int64) {
	for _, id := range listInventory {
		new(Inventory).ResetCacheWorker(id)
	}
}

func (t *Content) LoadMoreData(key, value string, userID int64, filterTarget payload.FilterTarget, listSelected []int64) (rows []ContentRecord, total int64, isMoreData, lastPage bool) {
	limit := 50
	page, offset := pagination.Pagination(key, limit)
	Query := mysql.Client.Limit(limit).Where("video_type = 1")
	if len(listSelected) > 0 {
		Query.Where("type = 1 and title like ? and user_id = ? and id not in ?", "%"+value+"%", userID, listSelected)
	} else {
		Query.Where("type = 1 and title like ? and user_id = ?", "%"+value+"%", userID)
	}
	if len(filterTarget.Language) > 0 {
		channels := []ChannelsRecord{}
		mysql.Client.Model(&ChannelsRecord{}).Where("user_id = ? and language in ?", userID, filterTarget.Language).Find(&channels)
		if len(channels) == 0 {
			return
		}
		channelsID := []int64{}
		for _, value := range channels {
			channelsID = append(channelsID, value.Id)
		}
		Query.Where("channels in ?", channelsID)
	}
	if len(filterTarget.ExcludeLanguage) > 0 {
		channels := []ChannelsRecord{}
		mysql.Client.Model(&ChannelsRecord{}).Where("user_id = ? and language not in ?", userID, filterTarget.ExcludeLanguage).Find(&channels)
		if len(channels) == 0 {
			return
		}
		channelsID := []int64{}
		for _, value := range channels {
			channelsID = append(channelsID, value.Id)
		}
		Query.Where("channels in ?", channelsID)
	}
	if len(filterTarget.Channels) > 0 {
		Query.Where("channels in ?", filterTarget.Channels)
	}
	if len(filterTarget.ExcludeChannels) > 0 {
		Query.Where("channels not in ?", filterTarget.ExcludeChannels)
	}
	if len(filterTarget.Keywords) > 0 {
		contentKeywords := []ContentKeywordRecord{}
		mysql.Client.Model(&ContentKeywordRecord{}).Where("user_id = ? and id in ?", userID, filterTarget.Keywords).Find(&contentKeywords)
		if len(contentKeywords) == 0 {
			return
		}
		contentsID := []int64{}
		for _, value := range contentKeywords {
			contentsID = append(contentsID, value.ContentId)
		}
		Query.Where("id in ?", contentsID)
	}
	if len(filterTarget.ExcludeKeywords) > 0 {
		contentKeywords := []ContentKeywordRecord{}
		mysql.Client.Model(&ContentKeywordRecord{}).Where("user_id = ? and id not in ?", userID, filterTarget.ExcludeKeywords).Find(&contentKeywords)
		if len(contentKeywords) == 0 {
			return
		}
		contentsID := []int64{}
		for _, value := range contentKeywords {
			contentsID = append(contentsID, value.ContentId)
		}
		Query.Where("id not in ?", contentsID)
	}

	Query.Offset(offset).Find(&rows)

	Query.Offset(0).Limit(10000000).Count(&total)
	// total = t.CountData(value, userID)
	totalPages := int(total) / limit
	if (int(total) % limit) != 0 {
		totalPages++
	}
	if page < totalPages {
		isMoreData = true
	}
	if page >= totalPages || len(rows) == 0 {
		isMoreData = false
		lastPage = true
	}
	return
}

func (t *Content) CountData(value string, userID int64) (count int64) {
	mysql.Client.Model(&ContentRecord{}).Where("type = 1 and title like ? and user_id = ? ", "%"+value+"%", userID).Count(&count)
	return
}

func (t *Content) GetContentByChannels(channelsID int64) (contents []ContentRecord) {
	mysql.Client.Model(&ContentRecord{}).Where("status = 1 and channels = ? ", channelsID).Find(&contents)
	return
}
