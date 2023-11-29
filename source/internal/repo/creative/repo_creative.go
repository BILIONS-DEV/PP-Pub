package creative

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/infrastructure/caching"
	"source/internal/entity/model"
	"source/pkg/pagination"
	"strconv"
	"strings"
	// "time"
)

type RepoCreative interface {
	Save(record *model.CreativeModel) error
	// UpdateCreative(input *model.CreativeModel) (err error)
	FindByFilter(input *FilterInput, user model.User) (creatives []model.CreativeModel, err error)
	GetById(idSearch int64) (row model.CreativeModel)
	// FindByID(idSearch int64) (row model.CreativeModel, err error)
	FindByID(record *model.CreativeModel, ID int64) (err error)
	Remove(id int64) (err error)
	RemoveRelationship(record *model.CreativeModel) (err error)
	GetNewCreative() (record model.CreativeModel, err error)
	IsExist(name string, id ...int64) bool

	GetAll() (record []model.CreativeModel, err error)
	IsExistCreativeSubmit(record model.CreativeSubmitModel) (model.CreativeSubmitModel, bool)
	SaveCreativeSubmit(record *model.CreativeSubmitModel) (err error)
	GetCreativeSubmitForCampaign(inputs model.CreativeSubmitFilter) (records []model.CreativeSubmitModel, total int, err error)
	RemoveCreativeSubmitForCampaign(campaignID, creativeID int64) (err error)
	UpdateFlagCreativeSubmit(campaignID, creativeID int64) (err error)
	RemoveAllCreativeSubmitByCampaign(campaignID int64) (err error)
	DisableNotificationCreativeSubmitByCampaign(campaignID int64) (err error)
	GetNotificationCreativeSubmit(campaignID int64) bool
	HandleTitleCreative(title, traffic_source string) string
	// ErrorTitle() (creativeSubmit []model.CreativeSubmitModel)
}

type creativeRepo struct {
	Db    *gorm.DB
	Cache caching.Cache
}

func (t *creativeRepo) Save(record *model.CreativeModel) error {
	return t.Db.Save(record).Error
}

func (t *creativeRepo) Validate(record model.CreativeModel) (err error) {
	return
}

func (t *creativeRepo) IsExist(name string, id ...int64) bool {
	var rec model.CreativeModel
	query := t.Db.Select("id").
		Where("name = ?", name)
	if len(id) > 0 {
		query.Where("id = ?", id[0])
	}
	query.Take(&rec)
	if rec.ID > 0 {
		return true
	}
	return false
}

func NewCreativeRepo(db *gorm.DB) *creativeRepo {
	return &creativeRepo{Db: db}
}

type AddInput struct {
	Id       string                     `json:"id"`
	SiteName string                     `json:"site_name"`
	Titles   []model.CreativeTitleModel `json:"titles"`
	Images   []model.CreativeImageModel `json:"images"`
}

func (t *creativeRepo) AddCreative(record *model.CreativeModel) (err error) {
	// var record model.CreativeModel
	// err = t.Validate(record)
	// if err != nil {
	// 	return
	// }
	// record = t.makeRow(input)
	err = t.Db.Save(record).Error
	// err = t.UpdateKeywordCreative(record, input)
	// //
	// err = t.UpdateLandingPageCreative(record, input)
	return
}

func (t *creativeRepo) makeRow(input *AddInput) (row model.CreativeModel) {
	if input.Id != "" && input.Id != "0" {
		row.ID, _ = strconv.ParseInt(input.Id, 10, 64)
	}
	row.SiteName = input.SiteName
	row.Titles = input.Titles
	row.Images = input.Images
	return
}

func (t *creativeRepo) UpdateCreative(record *model.CreativeModel) (err error) {
	// var record model.CreativeModel
	// err = t.Validate(record)
	// if err != nil {
	// 	return
	// }
	// record = t.makeRow(input)

	// err = t.Db.
	// 	Debug().
	// 	Save(record).Error
	// if err != nil {
	// 	return
	// }
	// err = t.UpdateKeywordCreative(record, input)

	// t.Db.Select("Titles").Delete(oldRecord.Creative.Titles)
	// t.Db.Select("Images").Delete(oldRecord.Creative.Images)
	err = t.Db.Save(&record).Error
	// /
	// err = t.UpdateLandingPageCreative(record, input)
	return
}

func (t *creativeRepo) RemoveRelationship(oldRecord *model.CreativeModel) (err error) {
	err = t.Db.Select(clause.Associations).Delete(oldRecord.Images, "creative_id = ?", oldRecord.ID).Error
	err = t.Db.Select(clause.Associations).Delete(oldRecord.Titles, "creative_id = ?", oldRecord.ID).Error
	if err != nil {
		return
	}

	return
}

type FilterInput struct {
	// filter
	// Condition map[string]interface{}
	// StartDate time.Time
	// EndDate   time.Time
	// Search    string
	// GroupBy   string
	// OrderBy   string
	// SubID     string
	// // order and pagination
	// Order string
	// Page  int
}

func (t *creativeRepo) FindByFilter(input *FilterInput, user model.User) (creatives []model.CreativeModel, err error) {
	t.Db.
		Scopes(
			// t.setFilterQueryOrderBy(input.OrderBy),
			// t.setFilterQueryGroupBy(input.GroupBy),
			// t.setFilterQueryDate(input.StartDate, input.EndDate),
			t.setFilterUser(user),
		).
		Model(&creatives).
		Order("id DESC").
		Select("*").Find(&creatives)
	return
}

func (t *creativeRepo) setFilterUser(user model.User) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if user.Permission == model.UserPermissionAdmin {
			return db
		}
		return db.Where("user_id = ?", user.ID)
	}
}

func (t *creativeRepo) GetById(id int64) (creative model.CreativeModel) {
	t.Db.
		Preload(clause.Associations).
		Where("id = ?", id).
		Find(&creative)
	return
}

func (t *creativeRepo) FindByID(record *model.CreativeModel, ID int64) (err error) {
	err = t.Db.
		// Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		First(&record, ID).Error
	return
}

func (t *creativeRepo) GetAll() (record []model.CreativeModel, err error) {
	err = t.Db.Preload(clause.Associations).Find(&record).Error
	return
}

func (t *creativeRepo) Remove(id int64) (err error) {
	err = t.Db.
		Preload("Keywords").
		Model(&model.CreativeModel{}).
		Delete(&model.CreativeModel{}, "id = ?", id).Error
	return
}

func (t *creativeRepo) GetNewCreative() (record model.CreativeModel, err error) {
	err = t.Db.Model(&model.CreativeModel{}).Order("id DESC").Find(&record).Error
	return
}

func (t *creativeRepo) IsExistCreativeSubmit(record model.CreativeSubmitModel) (model.CreativeSubmitModel, bool) {
	t.Db.Model(&model.CreativeSubmitModel{}).
		Where("campaign_id = ? AND creative_id = ? AND site_name = ? AND title = ? AND image = ?", record.CampaignID, record.CreativeID, record.SiteName, record.Title, record.Image).
		Find(&record)
	if record.ID > 0 {
		return record, true
	}
	return record, false
}

func (t *creativeRepo) SaveCreativeSubmit(record *model.CreativeSubmitModel) (err error) {
	err = t.Db.Save(record).Error
	return
}

func (t *creativeRepo) GetCreativeSubmitForCampaign(inputs model.CreativeSubmitFilter) (records []model.CreativeSubmitModel, total int, err error) {
	var totalRecords []model.CreativeSubmitModel
	err = t.Db.Model(&model.CreativeSubmitModel{}).
		Where("creative_id = ? AND campaign_id = ?", inputs.CreativeID, inputs.CampaignID).
		Find(&totalRecords).
		Scopes(
			pagination.Paginate(pagination.Params{
				Limit:  inputs.Length,
				Offset: inputs.Start,
			}),
		).
		Find(&records).Error
	total = len(totalRecords)
	return
}

func (t *creativeRepo) RemoveCreativeSubmitForCampaign(campaignID, creativeID int64) (err error) {
	err = t.Db.Model(&model.CreativeSubmitModel{}).
		Where("(creative_id = ? AND campaign_id = ? AND flag = 1) OR (creative_id != ? AND campaign_id = ?)", creativeID, campaignID, creativeID, campaignID).
		Delete(&model.CreativeSubmitModel{}).Error
	return
}

func (t *creativeRepo) UpdateFlagCreativeSubmit(campaignID, creativeID int64) (err error) {
	err = t.Db.Model(&model.CreativeSubmitModel{}).
		Where("creative_id = ? AND campaign_id = ?", creativeID, campaignID).
		Update("flag", 1).Error
	return
}

func (t *creativeRepo) RemoveAllCreativeSubmitByCampaign(campaignID int64) (err error) {
	err = t.Db.Where("campaign_id = ?", campaignID).Delete(&model.CreativeSubmitModel{}).Error
	return
}

func (t *creativeRepo) DisableNotificationCreativeSubmitByCampaign(campaignID int64) (err error) {
	err = t.Db.Model(&model.CreativeSubmitModel{}).Where("campaign_id = ?", campaignID).Update("new", "").Error
	return
}

func (t *creativeRepo) GetNotificationCreativeSubmit(campaignID int64) bool {
	var creativeSubmit model.CreativeSubmitModel
	t.Db.Where("campaign_id = ? AND new = ?", campaignID, "new").Find(&creativeSubmit)
	if creativeSubmit.ID > 0 {
		return true
	}
	return false
}

// func (t *creativeRepo) ErrorTitle() (creativeSubmit []model.CreativeSubmitModel) {
// 	t.Db.Debug().Where("title like ?", "%$$%").Find(&creativeSubmit)
// 	return
// }

func (t *creativeRepo) HandleTitleCreative(title, traffic_source string) string {
	switch traffic_source {
	case "Outbrain":
		if strings.Contains(title, "${city:capitalized}$") {
			title = strings.Replace(title, "${city:capitalized}$", "${city}$", -1)
		} else if strings.Contains(title, "${city}$") {
			title = strings.Replace(title, "${city}$", "${city}$", -1)
		} else if strings.Contains(title, "${city}") {
			title = strings.Replace(title, "${city}", "${city}$", -1)
		} else if strings.Contains(title, "{city}") {
			title = strings.Replace(title, "{city}", "${city}$", -1)
		}

		if strings.Contains(title, "${country:capitalized}$") {
			title = strings.Replace(title, "${country:capitalized}$", "${country}$", -1)
		} else if strings.Contains(title, "${country}$") {
			title = strings.Replace(title, "${country}$", "${country}$", -1)
		} else if strings.Contains(title, "${country}") {
			title = strings.Replace(title, "${country}", "${country}$", -1)
		} else if strings.Contains(title, "{country}") {
			title = strings.Replace(title, "{country}", "${country}$", -1)
		}
		break
	case "Taboola":
		if strings.Contains(title, "${city}$") {
			title = strings.Replace(title, "${city}$", "${city:capitalized}$", -1)
		} else if strings.Contains(title, "${city}") {
			title = strings.Replace(title, "${city}", "${city:capitalized}$", -1)
		} else if strings.Contains(title, "{city}") {
			title = strings.Replace(title, "{city}", "${city:capitalized}$", -1)
		}

		if strings.Contains(title, "${country}$") {
			title = strings.Replace(title, "${country}$", "${country:capitalized}$", -1)
		} else if strings.Contains(title, "${country}") {
			title = strings.Replace(title, "${country}", "${country:capitalized}$", -1)
		} else if strings.Contains(title, "{country}") {
			title = strings.Replace(title, "{country}", "${country:capitalized}$", -1)
		}
		break
	case "PocPoc":
		if strings.Contains(title, "${city:capitalized}$") {
			title = strings.Replace(title, "${city:capitalized}$", "${city}", -1)
		} else if strings.Contains(title, "${city}$") {
			title = strings.Replace(title, "${city}$", "${city}", -1)
		} else if strings.Contains(title, "${city}") {
			title = strings.Replace(title, "${city}", "${city}", -1)
		} else if strings.Contains(title, "{city}") {
			title = strings.Replace(title, "{city}", "${city}", -1)
		}

		if strings.Contains(title, "${country:capitalized}$") {
			title = strings.Replace(title, "${country:capitalized}$", "${country}", -1)
		} else if strings.Contains(title, "${country}$") {
			title = strings.Replace(title, "${country}$", "${country}", -1)
		} else if strings.Contains(title, "${country}") {
			title = strings.Replace(title, "${country}", "${country}", -1)
		} else if strings.Contains(title, "{country}") {
			title = strings.Replace(title, "{country}", "${country}", -1)
		}
		break
	case "Mgid":
		if strings.Contains(title, "${city:capitalized}$") {
			title = strings.Replace(title, "${city:capitalized}$", "{city}", -1)
		} else if strings.Contains(title, "${city}$") {
			title = strings.Replace(title, "${city}$", "{city}", -1)
		} else if strings.Contains(title, "${city}") {
			title = strings.Replace(title, "${city}", "{city}", -1)
		}

		if strings.Contains(title, "${country:capitalized}$") {
			title = strings.Replace(title, "${country:capitalized}$", "{country}", -1)
		} else if strings.Contains(title, "${country}$") {
			title = strings.Replace(title, "${country}$", "{country}", -1)
		} else if strings.Contains(title, "${country}") {
			title = strings.Replace(title, "${country}", "{country}", -1)
		}
		break
	default:
		return title
		break
	}

	return title
}
