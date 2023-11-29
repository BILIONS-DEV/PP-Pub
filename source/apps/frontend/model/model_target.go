package model

import (
	"gorm.io/gorm"
	"source/apps/frontend/payload"
	"source/core/technology/mysql"
	"source/pkg/utility"
)

type Target struct{}

type TargetRecord struct {
	mysql.TableTarget
}

func (TargetRecord) TableName() string {
	return mysql.Tables.Target
}

func (t *Target) CreateTarget(target TargetRecord) (err error) {
	err = mysql.Client.Create(&target).Error
	return
}

func (t *Target) DeleteTarget(target TargetRecord) (err error) {
	err = mysql.Client.Where(target).Delete(TargetRecord{}).Error
	return
}

func (t *Target) DeleteTargetInventory(target TargetRecord) (err error) {
	err = mysql.Client.Where(target).Where("inventory_id != 0").Delete(TargetRecord{}).Error
	return
}

func (t *Target) GetTargetLineItem(lineItemId int64) (records []TargetRecord) {
	mysql.Client.Where(TargetRecord{mysql.TableTarget{LineItemId: lineItemId}}).Find(&records)
	return
}

func (t *Target) GetTargetFloor(floorId int64) (records []TargetRecord) {
	mysql.Client.Where(TargetRecord{mysql.TableTarget{FloorId: floorId}}).Find(&records)
	return
}

func (t *Target) GetTargetIdentity(identityId int64) (records []TargetRecord) {
	mysql.Client.Where(TargetRecord{mysql.TableTarget{IdentityId: identityId}}).Find(&records)
	return
}

func (t *Target) GetTargetAbTesting(abTestingId int64) (records []TargetRecord) {
	mysql.Client.Where(TargetRecord{mysql.TableTarget{AbTestingId: abTestingId}}).Find(&records)
	return
}

func (t *Target) GetAllTargetIdentityValidate(userId int64, identityId int64) (records []TargetRecord) {
	//mysql.Client.Where("identity_id != 0 and user_id = ? and identity_id != ?", userId, identityId).Find(&records)
	mysql.Client.
		Select("target.id", "target.inventory_id").
		Joins("LEFT JOIN identity ON target.identity_id = identity.id").
		Where("target.identity_id != 0 AND target.identity_id != ? AND target.user_id = ?", identityId, userId).
		Where("identity.status = ? AND identity.deleted_at is NULL", mysql.TypeOn).
		Find(&records)
	return
}

func (t *Target) GetAllTargetIdentity(userId int64) (records []TargetRecord) {
	mysql.Client.
		Select("target.id", "target.inventory_id").
		Joins("LEFT JOIN identity ON target.identity_id = identity.id").
		Where("target.identity_id != 0 AND target.user_id = ?", userId).
		Where("identity.status = ? AND identity.deleted_at is NULL", mysql.TypeOn).
		Find(&records)
	return
}

func (t *Target) GetAllByUser(userId int64, option string, target payload.FilterTarget) (listId []payload.ListTarget, listAdTagFilter []int64) {
	switch option {
	case "domain":
		mysql.Client.Select("id", "name").Model(InventoryRecord{}).Where("user_id = ?", userId).Find(&listId)
		return
	case "adformat":
		mysql.Client.Select("id", "name").Model(AdTypeRecord{}).Find(&listId)
		return
	case "adsize":
		mysql.Client.Select("id", "name").Model(AdSizeRecord{}).Find(&listId)
		return
	case "adtag":
		// Tạo list id type cho 2 định loại banner và video
		var listIdTypeDisplay []int64
		var listIdTypeVideo []int64
		var listIdTypeSticky []int64

		for _, format := range target.Format {
			if format == 1 {
				listIdTypeDisplay = append(listIdTypeDisplay, format)
			} else if format == 5 {
				listIdTypeSticky = append(listIdTypeSticky, format)
			} else {
				listIdTypeVideo = append(listIdTypeVideo, format)
			}
		}
		// Filter for display trường hợp này sử dụng cho cả filter format là display hoặc format all
		if len(listIdTypeDisplay) > 0 || len(target.Format) == 0 {
			var listAdTagFilterDisplay []int64
			mysql.Client.Model(&AdTagRecord{}).Select("inventory_ad_tag.id").
				Where("user_id = ?", userId).
				Scopes(
					setFilterInventory(target),
					setFilterFormat(listIdTypeDisplay),
					setFilterSize(target),
				).
				Group("inventory_ad_tag.id").
				Find(&listAdTagFilterDisplay)
			listAdTagFilter = append(listAdTagFilter, listAdTagFilterDisplay...)
		}

		// Filter for video, nếu trong trường hợp filter format != all và có format cho video thì filter theo dạng video không cần tính đến size
		if len(listIdTypeVideo) > 0 {
			var listAdTagFilterVideo []int64
			mysql.Client.Model(&AdTagRecord{}).Select("inventory_ad_tag.id").
				Where("user_id = ?", userId).
				Scopes(
					setFilterInventory(target),
					setFilterFormat(listIdTypeVideo),
				).
				Group("inventory_ad_tag.id").
				Find(&listAdTagFilterVideo)
			listAdTagFilter = append(listAdTagFilter, listAdTagFilterVideo...)
		}
		if len(listIdTypeSticky) > 0 || len(target.Format) == 0 {
			var listAdTagFilterSticky []int64
			mysql.Client.Model(&AdTagRecord{}).Select("inventory_ad_tag.id").
				Where("user_id = ?", userId).
				Scopes(
					setFilterInventory(target),
					setFilterFormat(listIdTypeDisplay),
					setFilterSizeSticky(target),
				).
				Group("inventory_ad_tag.id").
				Find(&listAdTagFilterSticky)
			listAdTagFilter = append(listAdTagFilter, listAdTagFilterSticky...)
		}
		mysql.Client.Select("id", "name").Model(AdTagRecord{}).Scopes(
			func(db *gorm.DB) *gorm.DB {
				if len(target.Inventory) > 0 || len(target.Format) > 0 || len(target.Size) > 0 {
					return db.Where("id IN ?", listAdTagFilter)
				}
				return db
			},
		).Where("user_id = ?", userId).Find(&listId)
		return
	case "device":
		mysql.Client.Select("id", "name").Model(DeviceRecord{}).Find(&listId)
		return
	case "geography":
		mysql.Client.Select("id", "name").Model(CountryRecord{}).Find(&listId)
		return
	}
	return
}

func (t *Target) GetTargetByFilterFloor(inputs *payload.FloorFilterPayload, userId int64) (records []TargetRecord, err error) {

	err = mysql.Client.Select("floor_id").Where("user_id = ? AND floor_id != 0", userId).
		Scopes(
			t.setFilterDomain(inputs),
			t.setFilterAdFormat(inputs),
			t.setFilterAdSize(inputs),
			t.setFilterAdTag(inputs),
			t.setFilterDevice(inputs),
			t.setFilterCountry(inputs),
		).
		Group("floor_id").
		Find(&records).Error
	return
}

func (t *Target) GetTargetByFilterAbTesting(inputs *payload.AbTestingFilterPayload, userId int64) (records []TargetRecord, err error) {
	err = mysql.Client.Select("ab_testing_id").Where("user_id = ? AND ab_testing_id != 0", userId).
		Scopes(
			t.setFilterDomainAbTesting(inputs),
			t.setFilterAdFormatAbTesting(inputs),
			t.setFilterAdSizeAbTesting(inputs),
			t.setFilterAdTagAbTesting(inputs),
			t.setFilterDeviceAbTesting(inputs),
			t.setFilterCountryAbTesting(inputs),
		).
		Group("ab_testing_id").
		Find(&records).Error
	return
}

func (t *Target) GetTargetByFilterIdentity(inputs *payload.IdentityFilterPayload, userId int64) (records []TargetRecord, err error) {
	err = mysql.Client.Select("identity_id").Where("user_id = ? AND identity_id != 0", userId).
		Scopes(
			t.setFilterDomainIdentity(inputs),
		).
		Group("identity_id").
		Find(&records).Error
	return
}

func (t *Target) setFilterDomainIdentity(inputs *payload.IdentityFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.Domain != nil {
			switch inputs.PostData.Domain.(type) {
			case string, int:
				if inputs.PostData.Domain != "" {
					return db.Where("target.inventory_id = ? or target.inventory_id = -1", inputs.PostData.Domain)
				}
			case []string, []interface{}:
				return db.Where("target.inventory_id IN ? or target.inventory_id = -1", inputs.PostData.Domain)
			}
		}
		return db
	}
}

func (t *Target) setFilterDomain(inputs *payload.FloorFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.Domain != nil {
			switch inputs.PostData.Domain.(type) {
			case string, int:
				if inputs.PostData.Domain != "" {
					return db.Where("target.inventory_id = ? or target.inventory_id = -1", inputs.PostData.Domain)
				}
			case []string, []interface{}:
				return db.Where("target.inventory_id IN ? or target.inventory_id = -1", inputs.PostData.Domain)
			}
		}
		return db
	}
}

func (t *Target) setFilterAdFormat(inputs *payload.FloorFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.AdFormat != nil {
			switch inputs.PostData.AdFormat.(type) {
			case string, int:
				if inputs.PostData.AdFormat != "" {
					return db.Where("target.ad_format_id = ? or target.ad_format_id = -1", inputs.PostData.AdFormat)
				}
			case []string, []interface{}:
				return db.Where("target.ad_format_id IN ? or target.ad_format_id = -1", inputs.PostData.AdFormat)
			}
		}
		return db
	}
}

func (t *Target) setFilterAdSize(inputs *payload.FloorFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.AdSize != nil {
			switch inputs.PostData.AdSize.(type) {
			case string, int:
				if inputs.PostData.AdSize != "" {
					return db.Where("target.ad_size_id = ? or target.ad_size_id = -1", inputs.PostData.AdSize)
				}
			case []string, []interface{}:
				return db.Where("target.ad_size_id IN ? or target.ad_size_id = -1", inputs.PostData.AdSize)
			}
		}
		return db
	}
}

func (t *Target) setFilterAdTag(inputs *payload.FloorFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.AdTag != nil {
			switch inputs.PostData.AdTag.(type) {
			case string, int:
				if inputs.PostData.AdTag != "" {
					return db.Where("target.tag_id = ? or target.tag_id = -1", inputs.PostData.AdTag)
				}
			case []string, []interface{}:
				return db.Where("target.tag_id IN ? or target.tag_id = -1", inputs.PostData.AdTag)
			}
		}
		return db
	}
}

func (t *Target) setFilterDevice(inputs *payload.FloorFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.Device != nil {
			switch inputs.PostData.Device.(type) {
			case string, int:
				if inputs.PostData.Device != "" {
					return db.Where("target.device_id = ? or target.device_id = -1", inputs.PostData.Device)
				}
			case []string, []interface{}:
				return db.Where("target.device_id IN ? or target.device_id = -1", inputs.PostData.Device)
			}
		}
		return db
	}
}

func (t *Target) setFilterCountry(inputs *payload.FloorFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.Country != nil {
			switch inputs.PostData.Country.(type) {
			case string, int:
				if inputs.PostData.Country != "" {
					return db.Where("target.geo_id = ? or target.geo_id = -1", inputs.PostData.Country)
				}
			case []string, []interface{}:
				return db.Where("target.geo_id IN ? or target.geo_id = -1", inputs.PostData.Country)
			}
		}
		return db
	}
}

func (t *Target) HandleForIdentity(recOld []InventoryRecord, user UserRecord) (records []InventoryRecord) {
	var listInventoryTargeted []int64
	//targets := new(Target).GetAllTargetIdentity(user.Id)
	targets := new(Target).GetAllTargetIdentity(user.Id)
	for _, target := range targets {
		listInventoryTargeted = append(listInventoryTargeted, target.InventoryId)
	}
	for _, inventory := range recOld {
		if utility.InArray(inventory.Id, listInventoryTargeted, false) {
			inventory.Status = 3
		}
		records = append(records, inventory)
	}
	return
}

func (t *Target) setFilterDomainAbTesting(inputs *payload.AbTestingFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.Domain != nil {
			switch inputs.PostData.Domain.(type) {
			case string, int:
				if inputs.PostData.Domain != "" {
					return db.Where("target.inventory_id = ? or target.inventory_id = -1", inputs.PostData.Domain)
				}
			case []string, []interface{}:
				return db.Where("target.inventory_id IN ? or target.inventory_id = -1", inputs.PostData.Domain)
			}
		}
		return db
	}
}

func (t *Target) setFilterAdFormatAbTesting(inputs *payload.AbTestingFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.AdFormat != nil {
			switch inputs.PostData.AdFormat.(type) {
			case string, int:
				if inputs.PostData.AdFormat != "" {
					return db.Where("target.ad_format_id = ? or target.ad_format_id = -1", inputs.PostData.AdFormat)
				}
			case []string, []interface{}:
				return db.Where("target.ad_format_id IN ? or target.ad_format_id = -1", inputs.PostData.AdFormat)
			}
		}
		return db
	}
}

func (t *Target) setFilterAdSizeAbTesting(inputs *payload.AbTestingFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.AdSize != nil {
			switch inputs.PostData.AdSize.(type) {
			case string, int:
				if inputs.PostData.AdSize != "" {
					return db.Where("target.ad_size_id = ? or target.ad_size_id = -1", inputs.PostData.AdSize)
				}
			case []string, []interface{}:
				return db.Where("target.ad_size_id IN ? or target.ad_size_id = -1", inputs.PostData.AdSize)
			}
		}
		return db
	}
}

func (t *Target) setFilterAdTagAbTesting(inputs *payload.AbTestingFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.AdTag != nil {
			switch inputs.PostData.AdTag.(type) {
			case string, int:
				if inputs.PostData.AdTag != "" {
					return db.Where("target.tag_id = ? or target.tag_id = -1", inputs.PostData.AdTag)
				}
			case []string, []interface{}:
				return db.Where("target.tag_id IN ? or target.tag_id = -1", inputs.PostData.AdTag)
			}
		}
		return db
	}
}

func (t *Target) setFilterDeviceAbTesting(inputs *payload.AbTestingFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.Device != nil {
			switch inputs.PostData.Device.(type) {
			case string, int:
				if inputs.PostData.Device != "" {
					return db.Where("target.device_id = ? or target.device_id = -1", inputs.PostData.Device)
				}
			case []string, []interface{}:
				return db.Where("target.device_id IN ? or target.device_id = -1", inputs.PostData.Device)
			}
		}
		return db
	}
}

func (t *Target) setFilterCountryAbTesting(inputs *payload.AbTestingFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.Country != nil {
			switch inputs.PostData.Country.(type) {
			case string, int:
				if inputs.PostData.Country != "" {
					return db.Where("target.geo_id = ? or target.geo_id = -1", inputs.PostData.Country)
				}
			case []string, []interface{}:
				return db.Where("target.geo_id IN ? or target.geo_id = -1", inputs.PostData.Country)
			}
		}
		return db
	}
}
