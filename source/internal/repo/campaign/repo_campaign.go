package campaign

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/infrastructure/caching"
	"source/infrastructure/mysql"
	"source/internal/entity/model"
	"source/pkg/datatable"
	"strconv"
	"strings"
	// "time"
)

type RepoCampaign interface {
	Save(record *model.CampaignModel) error
	SaveCampaignGroup(record *model.CampaignGroupModel) error
	// UpdateCampaign(input *model.CampaignModel) (err error)
	Migrate()
	FindByFilter(inputs *FilterInput, user model.User) (response []model.CampaignModel, total int64, err error)
	GetById(id int64, unscoped ...bool) (row model.CampaignModel)
	GetByChannel(channel string) (row model.CampaignModel)
	// FindByID(idSearch int64) (row model.CampaignModel, err error)
	FindByID(record *model.CampaignModel, ID int64) (err error)
	Remove(id int64) (err error)
	UpdateName(campaignID int64, name string) (err error)
	UpdateUrlTrack(campaign model.CampaignModel, UrlTrackImp, urlTrackClick string) (err error)
	UpdateParamsForCampaignId(campaign model.CampaignModel, params string) (err error)
	RemoveRelationship(record *model.CampaignModel) (err error)
	GetNewCampaign() (record model.CampaignModel, err error)
	GetNewCampaignGroup() (record model.CampaignGroupModel, err error)
	IsExist(name string, id ...int64) bool
	GetCampaignsByCreative(creativeID int64) (records []model.CampaignModel, err error)
	GetNewCampaignID() (newID int64)
	AccountByObject(object, object_type string) (records []model.AccountModel, err error)
	AllAcounts() (records []model.AccountModel, err error)
	AddLandingPages(record *model.LandingPagesDemand) (err error)
	LoadLandingPagesByDemand(demandSource string) (records []model.LandingPagesDemand, err error)
	RemoveLandingPagesByDemand(demandSource string) (err error)
	FindAllTrafficSourceID(trafficSource string, accountID string) (campaigns []*model.CampaignModel, err error)
	GetByTrafficSourceID(trafficSourceID string) (campaign model.CampaignModel, err error)
	FindCacheSubDomainParams(key string) (record model.SubDomainParamsAE, err error)
}

type campaignRepo struct {
	Db    *gorm.DB
	Cache caching.Cache
}

func NewCampaignRepo(db *gorm.DB, cache caching.Cache) *campaignRepo {
	return &campaignRepo{Db: db, Cache: cache}
}

func (t *campaignRepo) Save(record *model.CampaignModel) error {
	return t.Db.Save(record).Error
}

func (t *campaignRepo) SaveCampaignGroup(record *model.CampaignGroupModel) error {
	return t.Db.Save(record).Error
}

func (t *campaignRepo) IsExist(name string, id ...int64) bool {
	var rec model.CampaignModel
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

type AddInput struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	TrafficSource string `json:"traffic_source"`
	DemandSource  string `json:"demand_source"`
	// PixelId         string           `json:"pixel_id"`
	Vertival        string           `json:"vertical"`
	MainKeyword     string           `json:"main_keyword"`
	Channel         string           `json:"channel"`
	GD              string           `json:"gd"`
	Keywords        []string         `json:"keywords"`
	LandingPages    string           `json:"landing_pages"`
	LandingPagesGEO []LandingPageGeo `json:"landing_pages_geo"`
}

type LandingPageGeo struct {
	ID          int64  `json:"id"`
	Geo         string `json:"geo"`
	LandingPage string `json:"landing_page"`
}

func (t *campaignRepo) AddCampaign(record *model.CampaignModel) (err error) {
	// var record model.CampaignModel
	// err = t.Validate(record)
	// if err != nil {
	// 	return
	// }
	// record = t.makeRow(input)
	err = t.Db.Save(record).Error
	// err = t.UpdateKeywordCampaign(record, input)
	// //
	// err = t.UpdateLandingPageCampaign(record, input)
	return
}

func (t *campaignRepo) makeRow(input *AddInput) (row model.CampaignModel) {
	if input.Id != "" && input.Id != "0" {
		row.ID, _ = strconv.ParseInt(input.Id, 10, 64)
	}
	row.Name = input.Name
	row.TrafficSource = input.TrafficSource
	row.DemandSource = input.DemandSource
	row.LandingPages = input.LandingPages
	row.Channel = input.Channel
	row.GD = input.GD
	row.MainKeyword = input.MainKeyword
	// row.PixelId = input.PixelId
	row.Vertical = input.Vertival
	// for _, landingPage := range input.LandingPagesGEO {
	//	row.LandingPage = append(row.LandingPage, model.CampaignLandingPageModel{
	//		ID:          landingPage.ID,
	//		CampaignID:  row.ID,
	//		Country:     landingPage.Geo,
	//		LandingPage: landingPage.LandingPage,
	//	})
	// }
	return
}

func (t *campaignRepo) UpdateCampaign(record *model.CampaignModel) (err error) {
	// var record model.CampaignModel
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
	// err = t.UpdateKeywordCampaign(record, input)

	// t.Db.Select("Titles").Delete(oldRecord.Creative.Titles)
	// t.Db.Select("Images").Delete(oldRecord.Creative.Images)
	err = t.Db.Save(&record).Error
	// /
	// err = t.UpdateLandingPageCampaign(record, input)
	return
}

func (t *campaignRepo) RemoveRelationship(oldRecord *model.CampaignModel) (err error) {

	if err = t.Db.Select(clause.Associations).Delete(oldRecord.Keywords, "campaign_id = ?", oldRecord.ID).Error; err != nil {
		return
	}
	if err = t.Db.Select(clause.Associations).Delete(oldRecord.LandingPageGeo, "campaign_id = ?", oldRecord.ID).Error; err != nil {
		return
	}
	if err = t.Db.Select(clause.Associations).Delete(oldRecord.Device, "campaign_id = ?", oldRecord.ID).Error; err != nil {
		return
	}
	if err = t.Db.Select(clause.Associations).Delete(oldRecord.Location, "campaign_id = ?", oldRecord.ID).Error; err != nil {
		return
	}
	if err = t.Db.Select(clause.Associations).Delete(oldRecord.Country, "campaign_id = ?", oldRecord.ID).Error; err != nil {
		return
	}

	return
}

func (t *campaignRepo) UpdateKeywordCampaign(campaign model.CampaignModel, input *AddInput) (err error) {
	// remove keyword for campaign
	err = t.RemoveAllKeywordForCampaign(campaign.ID)
	if err != nil {
		return
	}
	if len(input.Keywords) == 0 {
		return
	}
	for _, keyword := range input.Keywords {
		err = t.AddKeywordForCampaign(campaign.ID, keyword)
		if err != nil {
			return
		}
	}
	return
}

func (t *campaignRepo) RemoveAllKeywordForCampaign(campaignID int64) (err error) {
	err = t.Db.
		Model(&model.CampaignKeywordModel{}).
		Delete(&model.CampaignKeywordModel{}, "campaign_id = ?", campaignID).Error
	return
}

func (t *campaignRepo) AddKeywordForCampaign(campaignID int64, keyword string) (err error) {
	if keyword == "" {
		return
	}
	row := model.CampaignKeywordModel{
		CampaignID: campaignID,
		Keyword:    keyword,
	}
	err = t.Db.
		Save(&row).Error
	return
}

func (t *campaignRepo) UpdateLandingPageCampaign(campaign model.CampaignModel, input *AddInput) (err error) {
	// remove keyword for campaign
	err = t.RemoveAllLandingPageForCampaign(campaign.ID)
	if err != nil {
		return
	}
	if len(input.LandingPagesGEO) == 0 {
		return
	}
	for _, landingPage := range input.LandingPagesGEO {
		err = t.AddLandingPageForCampaign(campaign.ID, landingPage.Geo, landingPage.LandingPage)
		if err != nil {
			return
		}
	}
	return
}

func (t *campaignRepo) RemoveAllLandingPageForCampaign(campaignID int64) (err error) {
	err = t.Db.
		// Debug().
		Model(&model.CampaignLandingPageModel{}).
		Delete(&model.CampaignLandingPageModel{}, "campaign_id = ?", campaignID).Error
	return
}

func (t *campaignRepo) AddLandingPageForCampaign(campaignID int64, country string, landingPage string) (err error) {
	if landingPage == "" {
		return
	}
	row := model.CampaignLandingPageModel{
		CampaignID:  campaignID,
		Country:     country,
		LandingPage: landingPage,
	}
	err = t.Db.
		// Debug().
		Save(&row).Error
	return
}

func (t *campaignRepo) Validate(record model.CampaignModel) (err error) {
	return
}

type FilterInput struct {
	datatable.Request
	// filter
	// Condition map[string]interface{}
	// StartDate time.Time
	// EndDate   time.Time
	TrafficSource string
	DemandSource  string
	Search        string
	Offset        int
	Limit         int
	Order         string
	// SubID     string
	// // order and pagination

}

func (t *campaignRepo) FindByFilter(inputs *FilterInput, user model.User) (campaigns []model.CampaignModel, total int64, err error) {
	var allCampaigns []model.CampaignModel
	err = t.Db.Scopes(
		t.setFilterCampaign(inputs),
		t.setFilterUser(user),
	).
		Model(&campaigns).
		Order("id DESC").
		Select("*").
		Find(&allCampaigns).
		Scopes(
			mysql.Paginate(
				mysql.Deps{
					Offset: inputs.Offset,
					Limit:  inputs.Limit,
				},
			),
		).
		Find(&campaigns).Error
	if err != nil {
		return
	}
	total = int64(len(allCampaigns))
	return
}

func (t *campaignRepo) setFilterCampaign(inputs *FilterInput) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.TrafficSource != "" {
			db.Where("traffic_source = ?", inputs.TrafficSource)
		}
		if inputs.DemandSource != "" {
			db.Where("demand_source = ?", inputs.DemandSource)
		}
		if inputs.Search != "" {
			db.Where("Name like ?", "%"+inputs.Search+"%")
		}
		return db
	}
}

func (t *campaignRepo) setFilterUser(user model.User) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if user.Permission == model.UserPermissionAdmin {
			return db
		}
		return db.Where("user_id = ?", user.ID)
	}
}

func (t *campaignRepo) GetById(id int64, unscoped ...bool) (campaign model.CampaignModel) {
	t.Db.
		//Debug().
		Preload(clause.Associations).
		Scopes(
			func(db *gorm.DB) *gorm.DB {
				if len(unscoped) > 0 {
					if unscoped[0] {
						return db.Unscoped()
					}
				}
				return db
			},
		).
		Where("id = ?", id).
		Find(&campaign)
	return
}

func (t *campaignRepo) GetByChannel(channel string) (campaign model.CampaignModel) {
	t.Db.
		Preload(clause.Associations).
		Where("channel = ?", channel).
		Where("status = 'on'").
		Find(&campaign)
	return
}

func (t *campaignRepo) FindByID(record *model.CampaignModel, ID int64) (err error) {
	err = t.Db.
		// Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		First(&record, ID).Error
	return
}

func (t *campaignRepo) Remove(id int64) (err error) {
	err = t.Db.
		Preload("Keywords").
		Model(&model.CampaignModel{}).
		Delete(&model.CampaignModel{}, "id = ?", id).Error
	return
}

func (t *campaignRepo) UpdateName(campaignID int64, name string) (err error) {
	err = t.Db.
		Where("id = ?", campaignID).
		Model(&model.CampaignModel{}).
		Update("name", name).Error
	return
}

func (t *campaignRepo) UpdateUrlTrack(campaign model.CampaignModel, UrlTrackImp, urlTrackClick string) (err error) {
	campaign.UrlTrackImp = UrlTrackImp
	campaign.URLTrackClick = urlTrackClick
	err = t.Db.
		Where("id = ?", campaign.ID).
		Model(&model.CampaignModel{}).
		Updates(&campaign).Error
	return
}

func (t *campaignRepo) UpdateParamsForCampaignId(campaign model.CampaignModel, params string) (err error) {
	campaign.Params = params
	err = t.Db.
		Where("id = ?", campaign.ID).
		Model(&model.CampaignModel{}).
		Updates(&campaign).Error
	return
}

func (t *campaignRepo) Migrate() {
	err := t.Db.AutoMigrate(
		&model.CampaignModel{},
		&model.CampaignKeywordModel{},
		&model.CampaignLandingPageModel{},
	)
	if err != nil {
		panic(err)
	}
	return
}

func (t *campaignRepo) GetNewCampaign() (record model.CampaignModel, err error) {
	err = t.Db.Preload(clause.Associations).Model(&model.CampaignModel{}).Order("id DESC").Find(&record).Error
	return
}

func (t *campaignRepo) GetNewCampaignGroup() (record model.CampaignGroupModel, err error) {
	err = t.Db.Model(&model.CampaignGroupModel{}).Order("id DESC").Find(&record).Error
	return
}

func (t *campaignRepo) GetCampaignsByCreative(creativeID int64) (records []model.CampaignModel, err error) {
	err = t.Db.Model(&model.CampaignModel{}).Where("creative_id = ?", creativeID).Find(&records).Error
	return
}

func (t *campaignRepo) GetNewCampaignID() (newID int64) {
	var record model.CampaignModel
	t.Db.Model(&model.CampaignModel{}).Order("id DESC").Find(&record)
	newID = record.ID + 1
	return
}

func (t *campaignRepo) AccountByObject(object, object_type string) (records []model.AccountModel, err error) {
	err = t.Db.Where("object = ? AND object_type = ?", object, object_type).Find(&records).Error
	return
}

func (t *campaignRepo) AllAcounts() (records []model.AccountModel, err error) {
	err = t.Db.Order("object ASC").Find(&records).Error
	return
}

type LandingPageInput struct {
	DemandSource string
	LandingPage  string
	UserID       int64
}

func (t *campaignRepo) AddLandingPages(record *model.LandingPagesDemand) (err error) {
	err = t.Db.Save(record).Error
	return
}

func (t *campaignRepo) RemoveLandingPagesByDemand(demandSource string) (err error) {
	err = t.Db.Delete(&model.LandingPagesDemand{}, "demand_source = ?", demandSource).Error
	return
}

func (t *campaignRepo) LoadLandingPagesByDemand(demandSource string) (records []model.LandingPagesDemand, err error) {
	err = t.Db.
		Where("demand_source = ?", demandSource).
		Find(&records).Error
	return
}

func (t *campaignRepo) FindAllTrafficSourceID(trafficSource string, accountID string) (campaigns []*model.CampaignModel, err error) {
	err = t.Db.
		//Debug().
		Model(&model.CampaignModel{}).
		Select("traffic_source_id, name").
		Where("traffic_source = ? and account = ?", strings.Title(trafficSource), accountID).
		Where("status = 'on'").
		Find(&campaigns).Error
	return
}

func (t *campaignRepo) GetByTrafficSourceID(trafficSourceID string) (campaign model.CampaignModel, err error) {
	err = t.Db.
		//Debug().
		Unscoped().
		Model(&model.CampaignModel{}).
		Where(model.CampaignModel{TrafficSourceID: trafficSourceID}).
		Find(&campaign).Error
	return
}

func (t *campaignRepo) FindCacheSubDomainParams(key string) (record model.SubDomainParamsAE, err error) {
	err = t.Cache.Get(key, &record, record.SetName())
	return
}
