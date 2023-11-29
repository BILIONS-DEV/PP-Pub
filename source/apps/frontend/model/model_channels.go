package model

import (
	"encoding/csv"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"io"
	"log"
	"os"
	"source/apps/frontend/lang"
	"source/apps/frontend/payload"
	"source/apps/history"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"source/pkg/htmlblock"
	"source/pkg/pagination"
	"source/pkg/utility"
	"strconv"
	"strings"
	"time"
)

type Channels struct{}

type ChannelsRecord struct {
	mysql.TableChannels
}

type ChannelsIndex struct {
	Channel  ChannelsRecord
	Videos   []ContentRecord
	Category string
	Language string
	Duration string
	Count    int
	Views    int
}

func (ChannelsRecord) TableName() string {
	return mysql.Tables.Channels
}

func (t *Channels) Create(inputs payload.ChannelsCreate, user UserRecord, userAdmin UserRecord, lang lang.Translation) (record ChannelsRecord, errs []ajax.Error) {
	errs = t.ValidateCreate(inputs, user)
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
				Message: lang.Errors.ChannelsError.Add.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}

	// Tạo tag trong bảng channels_keyword
	for _, v := range inputs.Keyword {
		mysql.Client.Create(&ChannelsKeywordRecord{mysql.TableChannelsKeyword{
			ChannelsId: record.Id,
			Keyword:    v,
		}})
	}

	// list := t.GetPlaylistByChannels(record.Id)
	// for _, playlist := range list {
	// 	if playlist.Type != mysql.TYPEDisplay {
	// 		listInventory := new(Player).GetListInventoryByType(playlist.Type, user.Id)
	// 		t.UpdateInventory(listInventory)
	// 	}
	// }

	// Push History
	recordNew, _ := new(Channels).GetById(record.Id, user.Id)
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = user.Id
	}
	_ = history.PushHistory(&history.Channel{
		Detail:    history.DetailChannelFE,
		CreatorId: creatorId,
		RecordOld: mysql.TableChannels{},
		RecordNew: recordNew.TableChannels,
	})
	return
}

func (t *Channels) makeRowCreate(row payload.ChannelsCreate) (record ChannelsRecord) {
	record.Id = row.Id
	record.Name = row.Name
	record.Description = strings.TrimSpace(row.Description)
	record.Category = row.Category
	record.Language = row.Language
	record.Status = mysql.StatusApproved
	return
}

func (t *Channels) makeRowUpdate(row payload.ChannelsCreate, oldChannels ChannelsRecord) (record ChannelsRecord) {
	record.Id = row.Id
	record.Name = row.Name
	record.Description = strings.TrimSpace(row.Description)
	record.Category = row.Category
	record.Language = row.Language
	record.Status = oldChannels.Status
	return
}

func (t *Channels) ValidateCreate(inputs payload.ChannelsCreate, user UserRecord) (errs []ajax.Error) {
	if utility.ValidateString(inputs.Name) == "" {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Name is required",
		})
	}
	flagName := t.VerificationRecordName(inputs.Name, inputs.Id, user.Id)
	if !flagName {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Channels Name already exist",
		})
	}
	if inputs.Language == 0 {
		errs = append(errs, ajax.Error{
			Id:      "box_language",
			Message: "Language is required",
		})
	}
	if inputs.Category == 0 {
		errs = append(errs, ajax.Error{
			Id:      "box_category",
			Message: "Category is required",
		})
	}

	return
}

// func (t *Channels) GetByFilters(inputs *payload.ChannelsFilterPayload, userId int64) (response datatable.Response, err error) {
// 	var records []ChannelsRecord
// 	var total int64
// 	err = mysql.Client.Where("user_id = ?", userId).
// 		Scopes(
// 			// t.setFilterType(inputs),
// 			t.setFilterSearch(inputs),
// 			t.setFilterCategory(inputs),
// 		).
// 		Model(&records).Count(&total).
// 		Scopes(
// 			t.setOrder(inputs),
// 			pagination.Paginate(pagination.Params{
// 				Limit:  inputs.Length,
// 				Offset: inputs.Start,
// 			}),
// 		).Find(&records).Error
// 	if err != nil {
// 		lang := lang.Translation{}
// 		if !utility.IsWindow() {
// 			err = fmt.Errorf(lang.Errors.ChannelsError.List.ToString())
// 		}
// 		return datatable.Response{}, err
// 	}
// 	response.Draw = inputs.Draw
// 	response.RecordsFiltered = total
// 	response.RecordsTotal = total
// 	response.Data = t.MakeResponseDatatable(records)
// 	return
// }

func (t *Channels) GetByFilters(inputs *payload.ChannelsFilterPayload, userId int64) (rows []ChannelsIndex, err error) {
	var records []ChannelsRecord
	var total int64
	err = mysql.Client.Where("user_id = ?", userId).
		Scopes(
			// t.setFilterType(inputs),
			// t.setFilterSearch(inputs),
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
		lang := lang.Translation{}
		if !utility.IsWindow() {
			err = fmt.Errorf(lang.Errors.ChannelsError.List.ToString())
		}
		return
	}
	if len(records) == 0 {
		return
	}

	for _, record := range records {
		var row ChannelsIndex
		row.Channel = record
		// Get video content
		row.Videos = new(Content).GetContentByChannels(record.Id)

		// Langauge
		row.Language = new(Language).GetLanguageNameById(record.Language)
		// Category
		row.Category = new(Category).GetCategoryNameById(record.Category)
		// Count content
		row.Count = len(row.Videos)
		// Views
		row.Views = 0
		// Duration
		row.Duration = "00:00:00"
		if len(row.Videos) > 0 {
			var s = time.Second
			for _, video := range row.Videos {
				s = s + time.Duration(int64(video.Duration))*time.Second
			}
			row.Duration = Duration(s).Format("15:04:05")
		}

		rows = append(rows, row)
	}
	// response.Draw = inputs.Draw
	// response.RecordsFiltered = total
	// response.RecordsTotal = total
	// response.Data = t.MakeResponseDatatable(records)
	return
}

type Duration time.Duration

func (t Duration) Format(format string) string {
	z := time.Unix(0, 0).UTC()
	return z.Add(time.Duration(t)).Format(format)
}

// func (t *Channels) setFilterType(inputs *payload.ChannelsFilterPayload) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		if inputs.PostData.Type != nil {
// 			switch inputs.PostData.Type.(type) {
// 			case string, int:
// 				if inputs.PostData.VideoType != "" {
// 					return db.Where("type = ?", inputs.PostData.Type)
// 				}
// 			case []string, []interface{}:
// 				return db.Where("type IN ?", inputs.PostData.Type)
// 			}
// 		}
// 		return db
// 	}
// }

func (t *Channels) setFilterCategory(inputs *payload.ChannelsFilterPayload) func(db *gorm.DB) *gorm.DB {
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

func (t *Channels) setFilterSearch(inputs *payload.ChannelsFilterPayload) func(db *gorm.DB) *gorm.DB {
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

func (t *Channels) setOrder(inputs *payload.ChannelsFilterPayload) func(db *gorm.DB) *gorm.DB {
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

type ChannelsRecordDatatable struct {
	ChannelsRecord
	RowId        string `json:"DT_RowId"`
	Name         string `json:"name"`
	Description  string `json:"Description"`
	LanguageName string `json:"language"`
	Category     string `json:"category"`
	Keyword      string `json:"keyword"`
	Action       string `json:"action"`
}

func (t *Channels) MakeResponseDatatable(Channels []ChannelsRecord) (records []ChannelsRecordDatatable) {
	for _, channels := range Channels {
		channels.Description = truncateText(channels.Description, 20)
		language := new(Language).GetById(channels.Language)
		rec := ChannelsRecordDatatable{
			ChannelsRecord: channels,
			RowId:          "channels_" + strconv.FormatInt(channels.Id, 10),
			Name:           htmlblock.Render("channels/index/block.name.gohtml", channels).String(),
			Description:    htmlblock.Render("channels/index/block.description.gohtml", channels).String(),
			LanguageName:   htmlblock.Render("channels/index/block.language.gohtml", language).String(),
			// Category:      htmlblock.Render("content/index/block.category.gohtml", content).String(),
			// Tags:          new(InventoryAdTag).GetById(content.Tag).Name,
			Action: htmlblock.Render("channels/index/block.action.gohtml", channels).String(),
		}
		category := new(Category).GetById(channels.Category)
		rec.Category = category.Name

		records = append(records, rec)
	}
	return
}

func (t *Channels) HandleChannels(Channels []ChannelsRecord) (records []ChannelsRecord) {
	for _, Channel := range Channels {
		Channel.Description = truncateText(Channel.Description, 20)
		records = append(records, Channel)
	}
	return
}

func (t *Channels) GetById(id, userId int64) (record ChannelsRecord, err error) {
	err = mysql.Client.Where("id = ? and user_id = ?", id, userId).Find(&record).Error
	// Get rls
	record.TableChannels.GetRls()
	return
}

func (t *Channels) GetByUser(userId int64) (records []ChannelsRecord) {
	mysql.Client.Where("user_id = ? AND status = 1", userId).Find(&records)
	return
}

func (t *Channels) Update(inputs payload.ChannelsCreate, user UserRecord, userAdmin UserRecord, lang lang.Translation) (record ChannelsRecord, errs []ajax.Error) {
	oldChannels, flag := t.VerificationRecord(inputs.Id, user.Id)
	if !flag {
		errs = append(errs, ajax.Error{
			Id:      "id",
			Message: "You don't own this channels",
		})
		return
	}
	errs = t.ValidateCreate(inputs, user)
	if len(errs) > 0 {
		return
	}
	record = t.makeRowUpdate(inputs, oldChannels)
	record.UserId = user.Id
	err := mysql.Client.Save(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.ChannelsError.Edit.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
	}

	// Xóa các tag cũ để add lại tag mới
	new(ChannelsKeyword).DeleteTagByChannels(record.Id)
	// Tạo tag trong bảng channels_keyword
	for _, v := range inputs.Keyword {
		mysql.Client.Create(&ChannelsKeywordRecord{mysql.TableChannelsKeyword{
			ChannelsId: record.Id,
			Keyword:    v,
		}})
	}

	// Push History
	recordNew, _ := new(Channels).GetById(record.Id, user.Id)
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = user.Id
	}
	_ = history.PushHistory(&history.Channel{
		Detail:    history.DetailChannelFE,
		CreatorId: creatorId,
		RecordOld: oldChannels.TableChannels,
		RecordNew: recordNew.TableChannels,
	})
	return
}

func (this *Channels) Delete(id, userId int64, userAdmin UserRecord, lang lang.Translation) fiber.Map {
	record, _ := new(Channels).GetById(id, userId)
	err := mysql.Client.Model(&ChannelsRecord{}).Delete(&ChannelsRecord{}, "id = ? and user_id = ?", id, userId).Error
	if err != nil {
		if !utility.IsWindow() {
			return fiber.Map{
				"status":  "err",
				"message": lang.Errors.ChannelsError.Delete.ToString(),
				"id":      id,
			}
		}
		return fiber.Map{
			"status":  "err",
			"message": err.Error(),
			"id":      id,
		}
	}

	// History
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userId
	}
	_ = history.PushHistory(&history.Channel{
		Detail:    history.DetailChannelFE,
		CreatorId: creatorId,
		RecordOld: record.TableChannels,
		RecordNew: mysql.TableChannels{},
	})

	return fiber.Map{
		"status":  "success",
		"message": "done",
		"id":      id,
	}
}

func (t *Channels) GetAll(userId int64) (records []ChannelsRecord) {
	mysql.Client.Where("user_id = ?", userId).Find(&records)
	return
}

func (t *Channels) VerificationRecord(id, userId int64) (row ChannelsRecord, flag bool) {

	err := mysql.Client.Model(&ChannelsRecord{}).Where("id = ? and user_id = ?", id, userId).Find(&row).Error
	if err != nil || row.Id == 0 {
		flag = false
		return
	}

	flag = true
	// Get rls
	row.GetRls()
	return
}

// func (t *Channels) GetPlaylistByChannels(ChannelsId int64) (rows []PlayerRecord) {
// 	var playlist []int64
// 	mysql.Client.Model(&RlsPlaylistChannelsRecord{}).Select("playlist_id").Where("Channels_id = ?", ChannelsId).Find(&playlist)
// 	for _, id := range playlist {
// 		row := new(Player).GetDetail(id)
// 		if row.Id != 0 {
// 			rows = append(rows, row)
// 		}
// 	}
// 	return
// }

func (t *Channels) UpdateInventory(listInventory []int64) {
	for _, id := range listInventory {
		new(Inventory).ResetCacheWorker(id)
	}
}

func (t *Channels) ParserCSV() {
	// open file
	f, err := os.Open("C:/ProjectGo/selfserve.project/source/www/themes/muze/assets/language.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		type TableChannels struct {
			Code         string `gorm:"column:code" json:"code"`
			LanguageName string `gorm:"column:language_name" json:"language_name"`
		}

		language := TableChannels{
			Code:         rec[0],
			LanguageName: rec[1],
		}

		// insert language to mysql
		err = mysql.Client.Table("language").Create(language).Error
	}
}

func (t *Channels) LoadMoreData(key, value string, userID int64, filterTarget payload.FilterTarget, listSelected []int64) (rows []ChannelsRecord, isMoreData, lastPage bool) {
	limit := 30
	page, offset := pagination.Pagination(key, limit)
	Query := mysql.Client.Limit(limit)
	if len(listSelected) > 0 {
		Query.Where("name like ? and user_id = ? and id not in ?", "%"+value+"%", userID, listSelected)
	} else {
		Query.Where("name like ? and user_id = ?", "%"+value+"%", userID)
	}
	if len(filterTarget.Language) > 0 {
		Query.Where("language in ?", filterTarget.Language)
	}
	if len(filterTarget.ExcludeLanguage) > 0 {
		Query.Where("language not in ?", filterTarget.ExcludeLanguage)
	}

	Query.Offset(offset).Find(&rows)

	total := int64(0)
	Query.Offset(0).Limit(10000000).Count(&total)
	// total := t.CountData(value, userID)
	totalPages := int(total) / 30
	if (int(total) % 30) != 0 {
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

func (t *Channels) CountData(value string, userID int64) (count int64) {
	mysql.Client.Model(&ChannelsRecord{}).Where("name like ? and user_id = ? ", "%"+value+"%", userID).Count(&count)
	return
}

func (t *Channels) VerificationRecordName(name string, id, userId int64) bool {
	var row ChannelsRecord
	err := mysql.Client.Model(&ChannelsRecord{}).Where("name = ? and id != ? and user_id = ?", name, id, userId).Find(&row).Error
	if err != nil || row.Id != 0 {
		return false
	}
	return true
}
