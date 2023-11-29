package model

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"source/pkg/pagination"
	"source/pkg/utility"
	"strconv"
	"strings"
	"time"
)

type Bidder struct{}

type BidderRecord struct {
	mysql.TableBidder
}

func (BidderRecord) TableName() string {
	return mysql.Tables.Bidder
}

func validateLineAdsTxt(adsTxtLine string) (flag bool) {
	flag = true
	return
}

func (t *Bidder) GetAllUser(userId int64) (records []BidderRecord) {
	query := BidderRecord{mysql.TableBidder{
		UserId: userId,
	}}
	mysql.Client.Where(query).Find(&records)
	return
}

func (t *Bidder) GetAllBidderSystem() (records []BidderRecord) {
	query := BidderRecord{mysql.TableBidder{
		BidderType: 2,
	}}
	mysql.Client.Where(query).Find(&records)
	return
}

func (t *Bidder) GetById(id int64, userId int64) (record BidderRecord) {
	mysql.Client.Where("id = ? and user_id = ?", id, userId).First(&record, id)
	return
}

func (t *Bidder) GetByIdNoCheckUser(id int64) (record BidderRecord) {
	mysql.Client.Where("id = ?", id).First(&record, id)
	return
}

func (t *Bidder) CheckDefault(id int64) bool {
	var record BidderRecord
	mysql.Client.Select("is_default").Where("id = ?", id).First(&record, id)
	if record.IsDefault == mysql.TypeOn {
		return true
	}
	return false
}

func (t *Bidder) IsDisplayNameGoogleUnique(displayName string, userId int64) (flag bool) {
	var record BidderRecord
	mysql.Client.Select("id").Where("user_id = ? AND display_name = ?", userId, displayName).Last(&record)
	if record.Id > 0 {
		flag = true
	}
	return
}

func (t *Bidder) IsExists(bidderTemplateId int64, userId int64, id int64) bool {
	var record BidderRecord
	mysql.Client.Where("bidder_template_id = ? and user_id = ? and id != ?", bidderTemplateId, userId, id).Find(&record)
	if record.Id > 0 {
		return true
	}
	return false
}

func (t *Bidder) CountData(value string) (count int64) {
	mysql.Client.Model(&BidderRecord{}).Where("name like ?", "%"+value+"%").Count(&count)
	return
}

func (t *Bidder) LoadMoreData(key, value string) (rows []BidderRecord, isMoreData bool) {
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

func (t *Bidder) CheckUniqueBidderAlias(bidderId int64, bidderAlias string) (flag bool) {
	var rec BidderRecord
	mysql.Client.Select("id").Where("bidder_alias = ? and id != ?", bidderAlias, bidderId).Last(&rec)
	if rec.Id > 0 {
		flag = true
	} else {
		flag = false
	}
	return
}

var RunningPath string
var AssetsPath string

var maxSize = 15 * 1024 * 1024

var validExtensionXlsx = []string{".xlsx", ".xlsm", ".xls"}

func init() {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	RunningPath = filepath.Dir(d)
	AssetsPath = RunningPath + "/../../www/themes/muze/assets"
}

func (t *Bidder) SaveFileXlsx(ctx *fiber.Ctx) (errs []ajax.Error, data fiber.Map) {
	// rootDomain := ctx.Protocol() + "://" + ctx.Hostname()
	file, err := ctx.FormFile("file")
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "upload_failed",
			Message: "file error",
		})
		return
	}
	if file != nil {
		if file.Size > int64(maxSize) {
			errs = append(errs, ajax.Error{
				Id:      "upload_failed",
				Message: "You uploaded file over 10mb, please choose another file!",
			})
		}
		extension := filepath.Ext(file.Filename)
		extension = strings.ToLower(extension)
		flag := utility.InArray(extension, validExtensionXlsx, false)
		if flag {
			timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
			err := os.Mkdir(AssetsPath+"/xlsx/"+timeStamp, 0755)
			if err != nil {
				log.Fatal(err)
			}
			href := fmt.Sprintf(AssetsPath+"/xlsx/"+timeStamp+"/%s", file.Filename)
			err = ctx.SaveFile(file, href)
			if err != nil {
				errs = append(errs, ajax.Error{
					Id:      "upload_failed",
					Message: "save file error",
				})
				return
			}
			objectCpm, _ := new(System).BuildCpmFromXlsxAmz("/assets/xlsx/" + timeStamp + "/" + file.Filename)
			if len(objectCpm) == 0 {
				errs = append(errs, ajax.Error{
					Id:      "upload_failed",
					Message: "Not found CPM in xlsx",
				})
				return
			}
			data = fiber.Map{
				"xlsx_url": "/assets/xlsx/" + timeStamp + "/" + file.Filename,
			}
		} else {
			errs = append(errs, ajax.Error{
				Id:      "upload_failed",
				Message: "File upload is invalid, allowed extensions are: " + strings.Join(validExtensionXlsx, ","),
			})
		}
	}
	return
}

func (t *Bidder) GetAllBidderSystemByUser(userId int64) (records []BidderRecord) {
	query := BidderRecord{mysql.TableBidder{
		BidderType: 1,
		Status:     1,
		UserId:     userId,
	}}
	mysql.Client.Where(query).Find(&records)
	return
}
