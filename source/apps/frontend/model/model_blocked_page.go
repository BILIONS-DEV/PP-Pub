package model

import (
	"crypto/md5"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gocarina/gocsv"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"source/apps/frontend/lang"
	"source/apps/frontend/payload"
	"source/apps/history"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"source/pkg/logger"
	"source/pkg/utility"
	"strconv"
	"strings"
	"time"
)

var runningPath string
var assetsPath string

var maxSizeCsvBlockedPage = 15 * 1024 * 1024

var validExtensionCSV = []string{".csv"}

func init() {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	runningPath = filepath.Dir(d)
	assetsPath = runningPath + "/../../www/themes/muze/assets"
}

type BlockedPage struct {
}

type RuleBlockedPageRecord struct {
	mysql.TableRuleBlockedPage
}

func (RuleBlockedPageRecord) TableName() string {
	return mysql.Tables.RuleBlockedPage
}

func (t *BlockedPage) formatURL(pageURL string) string {
	parse, err := url.Parse(pageURL)
	if err != nil {
		return pageURL
	}
	newPageURL := parse.Host
	if parse.Path != "" {
		newPageURL += strings.TrimSuffix(parse.Path, "/")
	}
	if parse.RawQuery != "" {
		newPageURL += "?" + parse.RawQuery
	}
	if parse.Fragment != "" {
		newPageURL += "#" + parse.Fragment
	}
	return newPageURL
}

func (t *BlockedPage) GetById(id, userId int64) (record RuleRecord, err error) {
	err = record.GetById(id, userId)
	return
}

func (t *BlockedPage) GetByUser(userId int64) (records []RuleRecord) {
	records = new(Rule).GetByUser(userId)
	return
}

func (t *BlockedPage) SaveFileCSV(ctx *fiber.Ctx, file *multipart.FileHeader) (pathFile string, errs []ajax.Error) {
	// Validate file csv
	errs = t.validateCSV(file)
	if len(errs) > 0 {
		return
	}

	// Tạo thư mục timestamp chứa file
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	err := os.Mkdir(assetsPath+"/csv/"+timeStamp, 0755)
	if err != nil {
		fmt.Println(err)
		errs = append(errs, ajax.Error{
			Id:      "upload_failed",
			Message: "error",
		})
		return
	}

	// Tạo path file
	pathFile = fmt.Sprintf(assetsPath+"/csv/"+timeStamp+"/%s", file.Filename)

	// Save file
	err = ctx.SaveFile(file, pathFile)
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "upload_failed",
			Message: "error",
		})
		return
	}
	return
}

func (t *BlockedPage) validateCSV(file *multipart.FileHeader) (errs []ajax.Error) {
	if file.Size > int64(maxSizeCsvBlockedPage) {
		errs = append(errs, ajax.Error{
			Id:      "upload_failed",
			Message: "You uploaded file over 15mb, please choose another file!",
		})
		return
	}
	extension := filepath.Ext(file.Filename)
	extension = strings.ToLower(extension)
	if !utility.InArray(extension, validExtensionCSV, false) {
		errs = append(errs, ajax.Error{
			Id:      "upload_failed",
			Message: "Format file not approved!",
		})
		return
	}
	return
}

func (t *BlockedPage) ParserCSV(pathFile string) (outputs []payload.CSV, errs []ajax.Error) {
	f, err := os.Open(pathFile)
	if err != nil {
		return
	}
	defer f.Close()

	if err = gocsv.UnmarshalFile(f, &outputs); err != nil { // Load clients from file
		fmt.Println(err)
		return
	}
	return
}

func (t *BlockedPage) AddPost(inputs payload.RuleSubmit, user UserRecord, userAdmin UserRecord) (record RuleRecord, errs []ajax.Error) {
	// Validate
	errs = t.validateCreate(inputs)
	if len(errs) > 0 {
		return
	}

	// Check exists
	exists := t.CheckExist(user.Id, inputs.Name)
	if exists {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "name already exists",
		})
		return
	}

	// Make record
	record, err := t.makeRecordCreate(inputs, user)
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: "error",
		})
		return
	}

	// Create
	if err = new(Rule).Create(&record); err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: "error",
		})
		return
	}

	// Create rls blocked page
	var ruleBlockedPages []RuleBlockedPageRecord
	for _, page := range inputs.Pages {
		u, err := url.Parse("https://" + page)
		if err != nil {
			continue
		}
		ruleBlockedPages = append(ruleBlockedPages, RuleBlockedPageRecord{
			mysql.TableRuleBlockedPage{
				RuleID: record.ID,
				Page:   u.Hostname() + u.Path,
			},
		})
	}
	_ = t.CreateRls(ruleBlockedPages)

	_ = t.SetCacheBlocked(record.UserID)

	recordNew, _ := new(BlockedPage).GetById(record.ID, user.Id)
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = user.Id
	}
	_ = history.PushHistory(&history.BlockedPage{
		Detail:    history.DetailBlockedPageFE,
		CreatorId: creatorId,
		RecordOld: mysql.TableRule{},
		RecordNew: recordNew.TableRule,
	})
	return
}

func (t *BlockedPage) EditPost(inputs payload.RuleSubmit, user UserRecord, userAdmin UserRecord) (record RuleRecord, errs []ajax.Error) {
	// Validate
	errs = t.validateCreate(inputs)
	if len(errs) > 0 {
		return
	}

	// Check exists
	exists := t.CheckExist(user.Id, inputs.Name, inputs.ID)
	if exists {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "name already exists",
		})
		return
	}

	// Get recordOld
	recordOld, _ := new(BlockedPage).GetById(inputs.ID, user.Id)

	// Đánh dấu unblock cho tất cả unblock cũ của user
	err := new(CronJobBlockedPage).UnBlockAllByUser(user.Id)
	if err != nil {
		return
	}

	// Make record
	record, err = t.makeRecordUpdate(inputs, user)
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: "error",
		})
		return
	}

	// Delete rls
	if err = t.DeleteRls(record.ID); err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: "error",
		})
		return
	}

	// Create
	if err = new(Rule).Update(&record); err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: "error",
		})
		return
	}

	// Create rls blocked page
	var ruleBlockedPages []RuleBlockedPageRecord
	for _, page := range inputs.Pages {
		u, err := url.Parse("https://" + page)
		if err != nil {
			continue
		}
		ruleBlockedPages = append(ruleBlockedPages, RuleBlockedPageRecord{
			mysql.TableRuleBlockedPage{
				RuleID: record.ID,
				Page:   u.Hostname() + u.Path,
			},
		})
	}
	_ = t.CreateRls(ruleBlockedPages)

	// Cache
	_ = t.SetCacheBlocked(user.Id)

	recordNew, _ := new(BlockedPage).GetById(record.ID, user.Id)
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = user.Id
	}
	_ = history.PushHistory(&history.BlockedPage{
		Detail:    history.DetailBlockedPageFE,
		CreatorId: creatorId,
		RecordOld: recordOld.TableRule,
		RecordNew: recordNew.TableRule,
	})
	return
}

func (t *BlockedPage) Delete(id, userId int64, userAdmin UserRecord) fiber.Map {
	record, _ := new(BlockedPage).GetById(id, userId)
	// Get listKeyOld

	// Đánh dấu unblock cho tất cả unblock cũ của user
	err := new(CronJobBlockedPage).UnBlockAllByUser(userId)
	if err != nil {
		logger.Error(err.Error())
		return fiber.Map{
			"status":  "err",
			"message": "error",
			"id":      id,
		}
	}

	err = mysql.Client.Model(&RuleRecord{}).Delete(&RuleRecord{}, "id = ? and user_id = ? and type = 3", id, userId).Error
	if err != nil {
		if !utility.IsWindow() {
			return fiber.Map{
				"status":  "err",
				"message": "error",
				"id":      id,
			}
		}
		return fiber.Map{
			"status":  "err",
			"message": err.Error(),
			"id":      id,
		}
	}
	// Cache
	_ = t.SetCacheBlocked(userId)

	// History
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userId
	}
	_ = history.PushHistory(&history.BlockedPage{
		Detail:    history.DetailBlockedPageFE,
		CreatorId: creatorId,
		RecordOld: record.TableRule,
		RecordNew: mysql.TableRule{},
	})

	return fiber.Map{
		"status":  "success",
		"message": "done",
		"id":      id,
	}
}

func (t *BlockedPage) validateCreate(inputs payload.RuleSubmit) (errs []ajax.Error) {
	if govalidator.IsNull(inputs.Name) {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: lang.Translate.ErrorRequired.ToString(),
		})
	}
	if len(inputs.Pages) == 0 {
		errs = append(errs, ajax.Error{
			Id:      "box-upload-csv",
			Message: lang.Translate.ErrorRequired.ToString(),
		})
	}
	return
}

func (t *BlockedPage) validateUpdate(inputs payload.RuleSubmit) (errs []ajax.Error) {
	if govalidator.IsNull(inputs.Name) {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: lang.Translate.ErrorRequired.ToString(),
		})
	}
	if len(inputs.Pages) == 0 {
		errs = append(errs, ajax.Error{
			Id:      "box-upload-csv",
			Message: lang.Translate.ErrorRequired.ToString(),
		})
	}
	return
}

func (t *BlockedPage) makeRecordCreate(inputs payload.RuleSubmit, user UserRecord) (record RuleRecord, err error) {
	record.UserID = user.Id
	record.Name = inputs.Name
	record.Status = mysql.TYPEStatusOn
	record.Type = mysql.TYPERuleTypeBlockedPage
	return
}

func (t *BlockedPage) makeRecordUpdate(inputs payload.RuleSubmit, user UserRecord) (record RuleRecord, err error) {
	record.ID = inputs.ID
	record.UserID = user.Id
	record.Name = inputs.Name
	record.Status = mysql.TYPEStatusOn
	record.Type = mysql.TYPERuleTypeBlockedPage
	return
}

func (t *BlockedPage) MakeKeyCache(userID int64, pageURL string) string {
	// Loại dấu / ở đầu cuối PageURL
	pageURL = strings.TrimRight(pageURL, "/")
	pageURL = strings.TrimLeft(pageURL, "/")
	data := []byte(strconv.FormatInt(userID, 10) + "_" + pageURL)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func (t *BlockedPage) getCronJobBlock(userId int64) (listPageBlock []string, cronJobBlockedPages []*CronJobBlockedPageRecord, err error) {
	var rules []RuleRecord
	err = mysql.Client.Where("user_id = ?", userId).Find(&rules).Error
	for _, rule := range rules {
		var ruleBlockedPages []RuleBlockedPageRecord
		err = mysql.Client.Where("rule_id = ?", rule.ID).Find(&ruleBlockedPages).Error
		for _, ruleBlockedPage := range ruleBlockedPages {
			if !utility.InArray(ruleBlockedPage.Page, listPageBlock, false) {
				md := t.MakeKeyCache(userId, ruleBlockedPage.Page)
				cronJobBlockedPages = append(cronJobBlockedPages, &CronJobBlockedPageRecord{
					mysql.TableCronJobBlockedPage{
						MD5:    md,
						Type:   mysql.TYPECronjobBlockedPageBlock,
						UserID: userId,
						Page:   ruleBlockedPage.Page,
						Status: mysql.TYPEStatusCronjobBlockedPagePending,
					},
				})
			}
		}
	}
	return
}

func (t *BlockedPage) CheckExist(userId int64, name string, id ...int64) (exist bool) {
	var record RuleRecord
	db := mysql.Client.Select("id").Where("user_id = ? and name = ?", userId, name)
	if len(id) > 0 {
		db.Where("id != ?", id[0])
	}
	db.Find(&record)
	if record.ID > 0 {
		exist = true
	} else {
		exist = false
	}
	return
}

func (t *BlockedPage) CreateRls(rls []RuleBlockedPageRecord) (err error) {
	err = mysql.Client.Create(&rls).Error
	return
}

func (t *BlockedPage) DeleteRls(ruleId int64) (err error) {
	err = mysql.Client.Where(RuleBlockedPageRecord{mysql.TableRuleBlockedPage{
		RuleID: ruleId,
	}}).Delete(RuleBlockedPageRecord{}).Error
	return
}

func (t *BlockedPage) SetCacheBlocked(userId int64) (err error) {
	_, cronJobBlockedPages, err := t.getCronJobBlock(userId)
	if err != nil {
		logger.Error("bl-509: " + err.Error())
		return
	}
	for _, cronJobBlockedPage := range cronJobBlockedPages {
		err = new(CronJobBlockedPage).Save(cronJobBlockedPage)
		if err != nil {
			logger.Error("bl-515: " + err.Error())
			continue
		}
	}
	return err
}
