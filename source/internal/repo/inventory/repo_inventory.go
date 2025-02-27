package inventory

import (
	"gorm.io/gorm"
	"source/infrastructure/caching"
	"source/infrastructure/mysql"
	"source/internal/entity/model"
	// "source/core/technology/mysql"
)

type RepoInventory interface {
	Filter(input *InputFilter) (totalRecord int64, records []*model.InventoryModel, err error)
	FindByID(ID int64) (record *model.InventoryModel, err error)
	IsExists(input *InputIsExists, IDs ...int64) (exists bool)
	Save(record *model.InventoryModel) (err error)
	ResetCacheByID(ID int64) (err error)
	ResetCacheByUser(userID int64) (err error)
	DeleteByID(ID int64) error
	FindByBidderStatus(bidderID int64, status []model.TYPERlsBidderSystemInventoryStatus) (records []model.InventoryModel, err error)
	GetInventoriesByName(name string) (records []*model.InventoryModel, err error)
	UpdateAdsTxtDomain(id int64, adsTxt string) (err error)
	ListTagIdByDomainName(domain string) (records []model.InventoryAdTagModel, err error)
	GetByPublisher(userID int64) (records []model.InventoryModel, err error)
}

type inventoryRepo struct {
	Db    *gorm.DB
	Cache caching.Cache
}

func NewInventoryRepo(db *gorm.DB, cache caching.Cache) *inventoryRepo {
	return &inventoryRepo{Db: db, Cache: cache}
}

type InputIsExists struct {
	Name string
}

func (t *inventoryRepo) IsExists(input *InputIsExists, IDs ...int64) (exists bool) {
	tx := t.Db.
		Select("ID").
		Where("name = ?", input.Name)
	if len(IDs) > 0 {
		tx.Where("id != ?", IDs[0])
	}
	var record model.AdScheduleModel
	tx.Last(&record)
	if record.IsFound() {
		exists = true
	}
	return
}

type InputFilter struct {
	UserID     int64
	Status     interface{}
	SupplyType interface{}
	SyncAdsTxt interface{}
	Search     string
	Offset     int
	Limit      int
	Order      string
}

func (t *inventoryRepo) Filter(input *InputFilter) (totalRecord int64, records []*model.InventoryModel, err error) {
	var table = model.InventoryModel{}.TableName()
	err = t.Db.
		Scopes(
			t.setFilterCondition(input.UserID, input.Status, input.SupplyType, input.SyncAdsTxt),
			t.setFilterQuerySearch(input.Search),
		).
		Model(&records).Distinct(table + `.id`).Count(&totalRecord).
		Order(input.Order).
		Scopes(mysql.Paginate(mysql.Deps{Offset: input.Offset, Limit: input.Limit})).
		Select("*").Group(table + ".id").Find(&records).Error
	return
}

func (t *inventoryRepo) setFilterCondition(userID int64, listStatus, listSupplyType, listSyncAdsTxt interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var condition = make(map[string]interface{})
		// => convert input sang dữ liệu condition để query
		condition["sub_pub_id"] = userID
		if listStatus != nil {
			condition["status"] = listStatus
		}
		if listSupplyType != nil {
			condition["type"] = listSupplyType
		}
		if listSyncAdsTxt != nil {
			condition["sync_ads_txt"] = listSyncAdsTxt
		}
		return db.Where(condition)
	}
}

func (t *inventoryRepo) setFilterQuerySearch(search string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if search == "" {
			return db
		}
		return db.Where("name LIKE ?", "%"+search+"%")
	}
}

func (t *inventoryRepo) FindByID(ID int64) (record *model.InventoryModel, err error) {
	err = t.Db.
		// Debug().
		// Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload("AdBreakConfigs").
		Preload("AdBreakConfigs.AdTagUrls").
		First(&record, ID).Error
	return
}

func (t *inventoryRepo) GetInventoriesByName(name string) (records []*model.InventoryModel, err error) {
	err = t.Db.
		Where("name = ? AND status = 1", name).
		Find(&records).Error
	return
}

func (t *inventoryRepo) Save(record *model.InventoryModel) (err error) {
	err = t.Db.
		Save(record).Error
	return
}

func (t *inventoryRepo) DeleteByID(ID int64) error {
	return t.Db.Delete(&model.InventoryModel{}, ID).Error
}

func (t *inventoryRepo) ResetCacheByUser(userID int64) error {
	return t.Db.Model(&model.InventoryModel{}).Where("user_id = ?", userID).Update("render_cache", 1).Error
}

func (t *inventoryRepo) ResetCacheByID(ID int64) error {
	return t.Db.Model(&model.InventoryModel{}).Where("id = ?", ID).Update("render_cache", 1).Error
}

func (t *inventoryRepo) FindByBidderStatus(bidderID int64, status []model.TYPERlsBidderSystemInventoryStatus) (records []model.InventoryModel, err error) {
	t.Db.
		Joins("left join rls_bidder_system_inventory as rls on rls.inventory_name = inventory.name").
		Where("rls.bidder_id = ? and rls.status in ? and inventory.status = ?", bidderID, status, model.StatusApproved).
		Order("rls.id desc").
		Find(&records)
	return
}

func (t *inventoryRepo) UpdateAdsTxtDomain(id int64, adsTxt string) (err error) {
	var rec = model.InventoryModel{
		AdsTxtCustomByAdmin: model.TYPEAdsTxtCustom(adsTxt),
		RenderCache:         1,
	}
	return t.Db.Model(&model.InventoryModel{}).Where("id = ?", id).Updates(&rec).Error
}

func (t *inventoryRepo) ListTagIdByDomainName(domain string) (records []model.InventoryAdTagModel, err error) {
	err = t.Db.
		Model(&model.InventoryAdTagModel{}).
		Where("status = ?", model.TypeStatusAdTagLive).
		Scopes(
			t.setFilterDomainName(domain),
		).
		Find(&records).Error
	return
}

func (t *inventoryRepo) setFilterDomainName(domainName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		domains, err := t.GetInventoriesByName(domainName)
		if err != nil || len(domains) == 0 {
			return db.Where("inventory_id IN ?", []int{-1})
		}
		var domainIDs []int64
		for _, domain := range domains {
			domainIDs = append(domainIDs, domain.ID)
		}
		return db.Where("inventory_id IN ?", domainIDs)
	}
}

func (t *inventoryRepo) GetByPublisher(userID int64) (records []model.InventoryModel, err error) {
	err = t.Db.
		Model(&model.InventoryModel{}).
		Where("status = ? AND user_id = ?", model.StatusApproved, userID).
		// Debug().
		Find(&records).Error
	return
}
