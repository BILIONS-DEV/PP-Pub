package report_taboola

import (
	"encoding/hex"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"net/url"
	"path"
	"runtime"
	"source/infrastructure/caching"
	"source/infrastructure/kafka"
	"source/internal/entity/model"
	"source/pkg/utility"
	"strconv"
	"strings"
	"time"
)

type RepoReportTaboola interface {
	Decrypt(key string, value string, try ...int64) (cpc float64, err error)
	Migrate()
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	FindByQuery(query map[string]interface{}) (record *model.ReportTaboolaModel, err error)
	FindAllByQuery(query map[string]interface{}) (record []*model.ReportTaboolaModel, err error)
	Save(record *model.ReportTaboolaModel) (err error)
	SaveSlice(records []*model.ReportTaboolaModel) (err error)
	FindByDayForReportAff(day string) (records []*model.ReportTaboolaModel, err error)
}

type reportTaboolaRepo struct {
	Db         *gorm.DB
	Cache      caching.Cache
	Kafka      *kafka.Client
	pkgPath    string
	totalRetry int64
}

func NewReportTaboolaRepo(db *gorm.DB) *reportTaboolaRepo {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalln("No caller information")
	}
	pkgPath := path.Dir(filename)
	return &reportTaboolaRepo{
		Db:      db,
		pkgPath: pkgPath,
		//key:        "440040379f2d2a67c82220108635aba2", // key test
		totalRetry: 3,
	}
}

func (t *reportTaboolaRepo) Decrypt(key string, value string, try ...int64) (cpc float64, err error) {
	// Xử lý retry lại nếu có
	var Try int64
	if len(try) > 0 {
		Try = try[0] + 1
	}
	defer func() {
		if len(try) > 0 {
			if Try < t.totalRetry {
				time.Sleep(3 * time.Second)
				cpc, err = t.Decrypt(key, value, Try)
			}
		}
	}()
	return decodeCPCTaboola(key, value)
}

func decodeCPCTaboola(key, cpcEnc string) (float64, error) {
	decoded_value, _ := url.QueryUnescape(cpcEnc)
	decoded_value = strings.ReplaceAll(decoded_value, "-", "+")
	decoded_value = strings.ReplaceAll(decoded_value, "_", "/")
	b64, err := utility.Base64Decode(decoded_value)
	if err != nil {
		return 0, err
	}

	hexByte, err := hex.DecodeString(key)
	if err != nil {
		return 0, err
	}

	rs, err := utility.DecryptAes128Ecb([]byte(b64), hexByte)
	if err != nil {
		return 0, err
	}
	sp := strings.Split(string(rs), "_")
	if len(sp) > 0 {
		cpc, err := strconv.ParseFloat(sp[0], 64)
		if err != nil {
			return 0, err
		}
		return cpc, nil
	}

	return 0, nil
}

type InputIsExists struct {
	Account  string
	Date     string
	Campaign string
	SiteID   string
}

func (t *reportTaboolaRepo) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		//Debug().
		Where(input).
		Select("ID")
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.ReportTaboolaModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

func (t *reportTaboolaRepo) Migrate() {
	//os.Exit(1)
	err := t.Db.AutoMigrate(
		&model.ReportTaboolaModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *reportTaboolaRepo) FindByQuery(query map[string]interface{}) (record *model.ReportTaboolaModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		//Preload(clause.Associations).
		Where(query).
		Last(&record).Error
	return
}

func (t *reportTaboolaRepo) FindAllByQuery(query map[string]interface{}) (record []*model.ReportTaboolaModel, err error) {
	err = t.Db.
		//Debug().
		//Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Where(query).
		Find(&record).Error
	return
}

func (t *reportTaboolaRepo) Save(record *model.ReportTaboolaModel) (err error) {
	err = t.Db.
		//Debug().
		Save(record).Error
	return
}

func (t *reportTaboolaRepo) SaveSlice(records []*model.ReportTaboolaModel) (err error) {
	if len(records) == 0 {
		return
	}
	err = t.Db.
		//Debug().
		Save(&records).Error
	return
}

func (t *reportTaboolaRepo) FindByDayForReportAff(day string) (records []*model.ReportTaboolaModel, err error) {
	err = t.Db.
		//Debug().
		Model(model.ReportTaboolaModel{}).
		Where("date LIKE '%" + day + "%'").
		//Where("campaign IN ('22981627','24155854','24441474','24335252')").
		//Where("account != 'account_2'").
		Find(&records).Error
	return
}
