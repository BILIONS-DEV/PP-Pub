package model

import (
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

type InventoryAdTag struct{}

type InventoryAdTagRecord struct {
	mysql.TableInventoryAdTag
}

func (InventoryAdTagRecord) TableName() string {
	return mysql.Tables.InventoryAdTag
}

func (t *InventoryAdTag) GetByFilters(inputs *payload.InventoryAdTagFilterPayload, userLogin UserRecord, userAdmin UserRecord, lang lang.Translation) (response datatable.Response, err error) {
	var inventoryAdTags []InventoryAdTagRecord
	var total int64
	err = mysql.Client.Where("user_id = ?", userLogin.Id).
		Scopes(
			t.SetFilterStatus(inputs),
			t.setFilterSearch(inputs),
			t.setFilterInventory(inputs),
			t.setFilterUser(userLogin.Id),
			t.setFilterType(inputs),
		).
		Model(&inventoryAdTags).Count(&total).
		Scopes(
			t.setOrder(inputs),
			pagination.Paginate(pagination.Params{
				Limit:  inputs.Length,
				Offset: inputs.Start,
			}),
		).
		Find(&inventoryAdTags).Error
	if err != nil {
		if !utility.IsWindow() {
			err = fmt.Errorf(lang.Errors.AdTagError.List.ToString())
		}
		return datatable.Response{}, err
	}
	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.MakeResponseDatatable(inventoryAdTags, userLogin, userAdmin)
	return
}

type InventoryAdTagRecordDatatable struct {
	InventoryAdTagRecord
	RowId      string `json:"DT_RowId"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	Size       string `json:"size"`
	FloorPrice string `json:"floor_price"`
	Action     string `json:"action"`
}

func (t *InventoryAdTag) MakeResponseDatatable(inventoryAdTags []InventoryAdTagRecord, userLogin UserRecord, userAdmin UserRecord) (records []InventoryAdTagRecordDatatable) {
	for _, inventoryAdTag := range inventoryAdTags {
		var adSize AdSizeRecord
		adSize = new(AdSize).GetById(inventoryAdTag.PrimaryAdSize)
		isDisable := false
		if userAdmin.Permission == mysql.UserPermissionAdmin || userAdmin.Permission == mysql.UserPermissionSale || inventoryAdTag.Type == mysql.TYPEDisplay {
			isDisable = true
		} else if userLogin.Permission != mysql.UserPermissionManagedService {
			isDisable = true
		}
		rec := InventoryAdTagRecordDatatable{
			InventoryAdTagRecord: inventoryAdTag,
			RowId:                strconv.FormatInt(inventoryAdTag.Id, 10),
			Name:                 htmlblock.Render("supply/adtag/block.name.gohtml", inventoryAdTag).String(),
			Status:               htmlblock.Render("supply/adtag/block.status.gohtml", inventoryAdTag).String(),
			Type:                 inventoryAdTag.TableInventoryAdTag.Type.String(),
			Size:                 htmlblock.Render("supply/adtag/block.size.gohtml", adSize).String(),
			//FloorPrice:           "$" + fmt.Sprintf("%.2f", inventoryAdTag.TableInventoryAdTag.FloorPrice),
			Action: htmlblock.Render("supply/adtag/block.action.gohtml", fiber.Map{
				"adTag":     inventoryAdTag,
				"isDisable": isDisable,
			}).String(),
		}
		records = append(records, rec)
	}
	return
}

func (t *InventoryAdTag) SetFilterStatus(inputs *payload.InventoryAdTagFilterPayload) func(db *gorm.DB) *gorm.DB {
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

func (t *InventoryAdTag) setFilterType(inputs *payload.InventoryAdTagFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.Type != nil {
			switch inputs.PostData.Type.(type) {
			case string, int:
				if inputs.PostData.Type != "" {
					return db.Where("type = ?", inputs.PostData.Type)
				}
			case []string, []interface{}:
				return db.Where("type IN ?", inputs.PostData.Type)
			}
		}
		return db
	}
}

func (t *InventoryAdTag) setFilterSearch(inputs *payload.InventoryAdTagFilterPayload) func(db *gorm.DB) *gorm.DB {
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

func (t *InventoryAdTag) setFilterInventory(inputs *payload.InventoryAdTagFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var flag bool
		if inputs.PostData.InventoryId != 0 {
			flag = true
		}
		if !flag {
			return db
		}
		return db.Where("inventory_id = ?", inputs.PostData.InventoryId)
	}
}

func (t *InventoryAdTag) setFilterUser(userId int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		return db.Where("user_id = ?", userId)
	}
}

func (t *InventoryAdTag) setOrder(inputs *payload.InventoryAdTagFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(inputs.Order) > 0 {
			var orders []string
			for _, order := range inputs.Order {
				column := inputs.Columns[order.Column]
				if column.Data == "size" {
					orders = append(orders, fmt.Sprintf("primary_ad_size %s", order.Dir))
				} else {
					orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
				}
			}
			orderString := strings.Join(orders, ", ")
			return db.Order(orderString)
		}
		return db
	}
}

func (t *InventoryAdTag) Create(inputs payload.InventoryAdTagSubmit, user UserRecord) (record InventoryAdTagRecord, errs []ajax.Error) {
	// Validate inputs
	errs = t.Validate(inputs)
	if len(errs) > 0 {
		return
	}
	// Insert to database
	record = t.makeRecord(inputs)
	record.UserId = user.Id
	err := mysql.Client.Create(&record).Error
	if err != nil {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
	}
	return
}

func (t *InventoryAdTag) Validate(inputs payload.InventoryAdTagSubmit) (errs []ajax.Error) {
	if utility.ValidateString(inputs.Name) == "" {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Line Item Name is required",
		})
	}
	if inputs.FloorPrice <= 0 {
		errs = append(errs, ajax.Error{
			Id:      "floor_price",
			Message: "Floor Price is required",
		})
	}
	if utility.ValidateString(inputs.Gam) == "" {
		errs = append(errs, ajax.Error{
			Id:      "gam",
			Message: "GAM is required",
		})
	}
	return
}

func (t *InventoryAdTag) makeRecord(inputs payload.InventoryAdTagSubmit) (rec InventoryAdTagRecord) {
	rec = InventoryAdTagRecord{TableInventoryAdTag: mysql.TableInventoryAdTag{
		InventoryId: inputs.InventoryId,
		Name:        inputs.Name,
		Status:      inputs.Status,
		Type:        inputs.AdTagType,
		//FloorPrice:  inputs.FloorPrice,
		//SizeId:      inputs.AdTagSize,
		Gam:      inputs.Gam,
		PassBack: inputs.PassBack,
	}}
	return
}

func (t *InventoryAdTag) CountData(value string, target payload.FilterTarget, listAdTagFilter []int64, userId int64) (count int64) {
	mysql.Client.Model(&InventoryAdTagRecord{}).
		Where("user_id = ?", userId).
		Scopes(
			func(db *gorm.DB) *gorm.DB {
				if len(target.Inventory) > 0 || len(target.Format) > 0 || len(target.Size) > 0 {
					return db.Where("id IN ?", listAdTagFilter)
				}
				return db
			},
		).
		Where("name like ?", "%"+value+"%").Count(&count)
	return
}

func (t *InventoryAdTag) CountDataSystem(value string, target payload.FilterTarget, listAdTagFilter []int64) (count int64) {
	mysql.Client.Model(&InventoryAdTagRecord{}).
		Scopes(
			func(db *gorm.DB) *gorm.DB {
				if len(target.Inventory) > 0 || len(target.Format) > 0 || len(target.Size) > 0 {
					return db.Where("id IN ?", listAdTagFilter)
				}
				return db
			},
		).
		Where("name like ?", "%"+value+"%").Count(&count)
	return
}

func (t *InventoryAdTag) CountDataPageEdit(target payload.FilterTarget, listAdTagFilter []int64, userId int64) (count int64) {
	mysql.Client.Model(&InventoryAdTagRecord{}).
		Where("user_id = ?", userId).
		Scopes(
			func(db *gorm.DB) *gorm.DB {
				if len(target.Inventory) > 0 || len(target.Format) > 0 || len(target.Size) > 0 {
					return db.Where("id IN ?", listAdTagFilter)
				}
				return db
			},
		).Count(&count)
	return
}

func (t *InventoryAdTag) CountDataPageEditSystem(target payload.FilterTarget, listAdTagFilter []int64) (count int64) {
	mysql.Client.Model(&InventoryAdTagRecord{}).
		Scopes(
			func(db *gorm.DB) *gorm.DB {
				if len(target.Inventory) > 0 || len(target.Format) > 0 || len(target.Size) > 0 {
					return db.Where("id IN ?", listAdTagFilter)
				}
				return db
			},
		).Count(&count)
	return
}

func (t *InventoryAdTag) LoadMoreData(key, value string, target payload.FilterTarget, userId int64, listSelected []int64) (rows []InventoryAdTagRecord, isMoreData bool, listAdTagFilter []int64, lastPage bool) {
	//Tạo list id type cho 2 định loại banner và video
	var listIdTypeDisplay []int64
	//var listIdTypeSticky []int64
	var listIdTypeVideo []int64
	for _, format := range target.Format {
		if format == 1 || format == 5 {
			listIdTypeDisplay = append(listIdTypeDisplay, format)
			//} else if format == 5 {
			//	listIdTypeSticky = append(listIdTypeSticky, format)
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

	//Filter for video, nếu trong trường hợp filter format != all và có format cho video thì filter theo dạng video không cần tính đến size
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

	//if len(listIdTypeSticky) > 0 || len(target.Format) == 0 {
	//	var listAdTagFilterSticky []int64
	//	mysql.Client.Model(&AdTagRecord{}).Select("inventory_ad_tag.id").
	//		Where("user_id = ?", userId).
	//		Scopes(
	//			setFilterInventory(target),
	//			setFilterFormat(listIdTypeDisplay),
	//			setFilterSizeSticky(target),
	//		).
	//		Group("inventory_ad_tag.id").
	//		Find(&listAdTagFilterSticky)
	//	listAdTagFilter = append(listAdTagFilter, listAdTagFilterSticky...)
	//}

	//Get adTag từ list ad tag đã filter
	limit := 10
	page, offset := pagination.Pagination(key, limit)
	if len(listSelected) > 0 {
		mysql.Client.
			Where("user_id = ?", userId).
			Scopes(
				func(db *gorm.DB) *gorm.DB {
					if len(target.Inventory) > 0 || len(target.Format) > 0 || len(target.Size) > 0 {
						return db.Where("id IN ? and id not in ?", listAdTagFilter, listSelected)
					}
					return db
				},
			).
			Where("name like ?", "%"+value+"%").Limit(limit).Offset(offset).Find(&rows)
	} else {
		mysql.Client.
			Where("user_id = ?", userId).
			Scopes(
				func(db *gorm.DB) *gorm.DB {
					if len(target.Inventory) > 0 || len(target.Format) > 0 || len(target.Size) > 0 {
						return db.Where("id IN ?", listAdTagFilter)
					}
					return db
				},
			).
			Where("name like ?", "%"+value+"%").Limit(limit).Offset(offset).Find(&rows)
	}
	total := t.CountData(value, target, listAdTagFilter, userId)
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

func (t *InventoryAdTag) LoadMoreDataSystem(key, value string, target payload.FilterTarget, listSelected []int64) (rows []InventoryAdTagRecord, isMoreData bool, listAdTagFilter []int64, lastPage bool) {
	//Tạo list id type cho 2 định loại banner và video
	var listIdTypeDisplay []int64
	//var listIdTypeSticky []int64
	var listIdTypeVideo []int64
	for _, format := range target.Format {
		if format == 1 || format == 5 {
			listIdTypeDisplay = append(listIdTypeDisplay, format)
			//} else if format == 5 {
			//	listIdTypeSticky = append(listIdTypeSticky, format)
		} else {
			listIdTypeVideo = append(listIdTypeVideo, format)
		}
	}
	// Filter for display trường hợp này sử dụng cho cả filter format là display hoặc format all
	if len(listIdTypeDisplay) > 0 || len(target.Format) == 0 {
		var listAdTagFilterDisplay []int64
		mysql.Client.Model(&AdTagRecord{}).Select("inventory_ad_tag.id").
			Scopes(
				setFilterInventory(target),
				setFilterFormat(listIdTypeDisplay),
				setFilterSize(target),
			).
			Group("inventory_ad_tag.id").
			Find(&listAdTagFilterDisplay)
		listAdTagFilter = append(listAdTagFilter, listAdTagFilterDisplay...)
	}

	//Filter for video, nếu trong trường hợp filter format != all và có format cho video thì filter theo dạng video không cần tính đến size
	if len(listIdTypeVideo) > 0 {
		var listAdTagFilterVideo []int64
		mysql.Client.Model(&AdTagRecord{}).Select("inventory_ad_tag.id").
			Scopes(
				setFilterInventory(target),
				setFilterFormat(listIdTypeVideo),
			).
			Group("inventory_ad_tag.id").
			Find(&listAdTagFilterVideo)
		listAdTagFilter = append(listAdTagFilter, listAdTagFilterVideo...)
	}

	//if len(listIdTypeSticky) > 0 || len(target.Format) == 0 {
	//	var listAdTagFilterSticky []int64
	//	mysql.Client.Model(&AdTagRecord{}).Select("inventory_ad_tag.id").
	//		Where("user_id = ?", userId).
	//		Scopes(
	//			setFilterInventory(target),
	//			setFilterFormat(listIdTypeDisplay),
	//			setFilterSizeSticky(target),
	//		).
	//		Group("inventory_ad_tag.id").
	//		Find(&listAdTagFilterSticky)
	//	listAdTagFilter = append(listAdTagFilter, listAdTagFilterSticky...)
	//}

	//Get adTag từ list ad tag đã filter
	limit := 10
	page, offset := pagination.Pagination(key, limit)
	if len(listSelected) > 0 {
		mysql.Client.
			Scopes(
				func(db *gorm.DB) *gorm.DB {
					if len(target.Inventory) > 0 || len(target.Format) > 0 || len(target.Size) > 0 {
						return db.Where("id IN ? and id not in ?", listAdTagFilter, listSelected)
					}
					return db
				},
			).
			Where("name like ?", "%"+value+"%").Limit(limit).Offset(offset).Find(&rows)
	} else {
		mysql.Client.
			Scopes(
				func(db *gorm.DB) *gorm.DB {
					if len(target.Inventory) > 0 || len(target.Format) > 0 || len(target.Size) > 0 {
						return db.Where("id IN ?", listAdTagFilter)
					}
					return db
				},
			).
			Where("name like ?", "%"+value+"%").Limit(limit).Offset(offset).Find(&rows)
	}
	total := t.CountDataSystem(value, target, listAdTagFilter)
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

func (t *InventoryAdTag) LoadMoreDataPageEdit(target payload.FilterTarget, userId int64, listSelected []int64) (rows []InventoryAdTagRecord, isMoreData bool, listAdTagFilter []int64, lastPage bool) {
	//Tạo list id type cho 2 định loại banner và video
	var listIdTypeDisplay []int64
	var listIdTypeVideo []int64
	//var listIdTypeSticky []int64
	for _, format := range target.Format {
		if format == 1 || format == 5 {
			listIdTypeDisplay = append(listIdTypeDisplay, format)
			//} else if format == 5 {
			//	listIdTypeSticky = append(listIdTypeSticky, format)
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

	//Filter for video, nếu trong trường hợp filter format != all và có format cho video thì filter theo dạng video không cần tính đến size
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

	//if len(listIdTypeSticky) > 0 || len(target.Format) == 0 {
	//	var listAdTagFilterSticky []int64
	//	mysql.Client.Model(&AdTagRecord{}).Select("inventory_ad_tag.id").
	//		Where("user_id = ?", userId).
	//		Scopes(
	//			setFilterInventory(target),
	//			setFilterFormat(listIdTypeDisplay),
	//			setFilterSizeSticky(target),
	//		).
	//		Group("inventory_ad_tag.id").
	//		Find(&listAdTagFilterSticky)
	//	listAdTagFilter = append(listAdTagFilter, listAdTagFilterSticky...)
	//}

	//Get adTag từ list ad tag đã filter
	limit := 10
	page, offset := pagination.Pagination("1", limit)
	if len(listSelected) > 0 {
		mysql.Client.
			Where("user_id = ?", userId).
			Scopes(
				func(db *gorm.DB) *gorm.DB {
					if len(target.Inventory) > 0 || len(target.Format) > 0 || len(target.Size) > 0 {
						return db.Where("id IN ? and id not in ?", listAdTagFilter, listSelected)
					}
					return db
				},
			).Limit(limit).Offset(offset).Find(&rows)
	} else {
		mysql.Client.
			Where("user_id = ?", userId).
			Scopes(
				func(db *gorm.DB) *gorm.DB {
					if len(target.Inventory) > 0 || len(target.Format) > 0 || len(target.Size) > 0 {
						return db.Where("id IN ?", listAdTagFilter)
					}
					return db
				},
			).Limit(limit).Offset(offset).Find(&rows)
	}
	if len(rows) > 10 {
		rows = rows[0:9]
	}
	total := t.CountDataPageEdit(target, listAdTagFilter, userId)
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

func (t *InventoryAdTag) LoadMoreDataPageEditSystem(target payload.FilterTarget, listSelected []int64) (rows []InventoryAdTagRecord, isMoreData bool, listAdTagFilter []int64, lastPage bool) {
	//Tạo list id type cho 2 định loại banner và video
	var listIdTypeDisplay []int64
	var listIdTypeVideo []int64
	//var listIdTypeSticky []int64
	for _, format := range target.Format {
		if format == 1 || format == 5 {
			listIdTypeDisplay = append(listIdTypeDisplay, format)
			//} else if format == 5 {
			//	listIdTypeSticky = append(listIdTypeSticky, format)
		} else {
			listIdTypeVideo = append(listIdTypeVideo, format)
		}
	}
	// Filter for display trường hợp này sử dụng cho cả filter format là display hoặc format all
	if len(listIdTypeDisplay) > 0 || len(target.Format) == 0 {
		var listAdTagFilterDisplay []int64
		mysql.Client.Model(&AdTagRecord{}).Select("inventory_ad_tag.id").
			Scopes(
				setFilterInventory(target),
				setFilterFormat(listIdTypeDisplay),
				setFilterSize(target),
			).
			Group("inventory_ad_tag.id").
			Find(&listAdTagFilterDisplay)
		listAdTagFilter = append(listAdTagFilter, listAdTagFilterDisplay...)
	}

	//Filter for video, nếu trong trường hợp filter format != all và có format cho video thì filter theo dạng video không cần tính đến size
	if len(listIdTypeVideo) > 0 {
		var listAdTagFilterVideo []int64
		mysql.Client.Model(&AdTagRecord{}).Select("inventory_ad_tag.id").
			Scopes(
				setFilterInventory(target),
				setFilterFormat(listIdTypeVideo),
			).
			Group("inventory_ad_tag.id").
			Find(&listAdTagFilterVideo)
		listAdTagFilter = append(listAdTagFilter, listAdTagFilterVideo...)
	}

	//if len(listIdTypeSticky) > 0 || len(target.Format) == 0 {
	//	var listAdTagFilterSticky []int64
	//	mysql.Client.Model(&AdTagRecord{}).Select("inventory_ad_tag.id").
	//		Where("user_id = ?", userId).
	//		Scopes(
	//			setFilterInventory(target),
	//			setFilterFormat(listIdTypeDisplay),
	//			setFilterSizeSticky(target),
	//		).
	//		Group("inventory_ad_tag.id").
	//		Find(&listAdTagFilterSticky)
	//	listAdTagFilter = append(listAdTagFilter, listAdTagFilterSticky...)
	//}

	//Get adTag từ list ad tag đã filter
	limit := 10
	page, offset := pagination.Pagination("1", limit)
	if len(listSelected) > 0 {
		mysql.Client.
			Scopes(
				func(db *gorm.DB) *gorm.DB {
					if len(target.Inventory) > 0 || len(target.Format) > 0 || len(target.Size) > 0 {
						return db.Where("id IN ? and id not in ?", listAdTagFilter, listSelected)
					}
					return db
				},
			).Limit(limit).Offset(offset).Find(&rows)
	} else {
		mysql.Client.
			Scopes(
				func(db *gorm.DB) *gorm.DB {
					if len(target.Inventory) > 0 || len(target.Format) > 0 || len(target.Size) > 0 {
						return db.Where("id IN ?", listAdTagFilter)
					}
					return db
				},
			).Limit(limit).Offset(offset).Find(&rows)
	}
	if len(rows) > 10 {
		rows = rows[0:9]
	}
	total := t.CountDataPageEditSystem(target, listAdTagFilter)
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

func setFilterInventory(target payload.FilterTarget) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(target.Inventory) > 0 {
			return db.Where("inventory_id IN ?", target.Inventory)
		}
		return db
	}
}

func setFilterFormat(target []int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(target) > 0 {
			return db.Where("type IN ?", target)
		}
		return db
	}
}

func setFilterSize(target payload.FilterTarget) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(target.Size) > 0 {
			return db.Where("primary_ad_size IN ?", target.Size)
		}
		return db
	}
}

func setFilterSizeSticky(target payload.FilterTarget) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(target.Size) > 0 {
			return db.Where("primary_ad_size IN ?", target.Size)
		}
		return db
	}
}

func (t *InventoryAdTag) GetAll(userId int64) (records []InventoryAdTagRecord) {
	mysql.Client.Where("user_id = ?", userId).Find(&records)
	return
}

func (t *InventoryAdTag) GetById(id int64) (record InventoryAdTagRecord) {
	if id < 1 {
		return
	}
	mysql.Client.First(&record, id)
	return
}

func (t *InventoryAdTag) GetByInventory(id int64) (records []InventoryAdTagRecord) {
	mysql.Client.Where("inventory_id = ? and status != 3", id).Order("id DESC").Find(&records)
	return
}

func (t *InventoryAdTag) GetByIdForFilter(id int64) (record InventoryAdTagRecord) {
	mysql.Client.Unscoped().First(&record, id)
	return
}

func (t *InventoryAdTag) Archived(id, userId int64, userAdmin UserRecord) fiber.Map {
	err := mysql.Client.Model(&InventoryAdTagRecord{}).
		Where("id = ? AND user_id = ?", id, userId).
		Update("status", mysql.TypeStatusAdTagArchived).Error
	if err != nil {
		if !utility.IsWindow() {
			return fiber.Map{
				"status":  "err",
				"message": lang.Translate.Errors.AdTagError.Delete.ToString(),
				"id":      id,
			}
		}
		return fiber.Map{
			"status":  "err",
			"message": err.Error(),
			"id":      id,
		}
	} else {
		row, _ := new(AdTag).GetDetail(id)
		// History
		var creatorId int64
		if userAdmin.Id != 0 {
			creatorId = userAdmin.Id
		} else {
			creatorId = userId
		}
		_ = history.PushHistory(&history.AdTag{
			Detail:    history.DetailAdTagFE,
			CreatorId: creatorId,
			RecordOld: row.TableInventoryAdTag,
			RecordNew: mysql.TableInventoryAdTag{},
		})
		new(Inventory).ResetCacheWorker(row.InventoryId)
		return fiber.Map{
			"status":  "success",
			"message": "done",
			"id":      id,
		}
	}
}

func (t *InventoryAdTag) Delete(id, userId int64, lang lang.Translation) fiber.Map {
	err := mysql.Client.Model(&InventoryAdTagRecord{}).Delete(&InventoryAdTagRecord{}, "id = ? and user_id = ?", id, userId).Error
	if err != nil {
		if !utility.IsWindow() {
			return fiber.Map{
				"status":  "err",
				"message": lang.Errors.AdTagError.Delete.ToString(),
				"id":      id,
			}
		}
		return fiber.Map{
			"status":  "err",
			"message": err.Error(),
			"id":      id,
		}
	} else {
		row, _ := new(AdTag).GetDetail(id)
		new(Inventory).ResetCacheWorker(row.InventoryId)
		return fiber.Map{
			"status":  "success",
			"message": "done",
			"id":      id,
		}
	}
}

func (t *InventoryAdTag) GetTagByFilter(target payload.FilterTarget, userId int64) (rows []InventoryAdTagRecord) {
	//Tạo list id type cho 2 định loại banner và video
	var listIdTypeDisplay []int64
	var listIdTypeVideo []int64
	var listAdTagFilter []int64
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

	//Filter for video, nếu trong trường hợp filter format != all và có format cho video thì filter theo dạng video không cần tính đến size
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

	mysql.Client.
		Where("user_id = ?", userId).
		Scopes(
			func(db *gorm.DB) *gorm.DB {
				if len(target.Inventory) > 0 || len(target.Format) > 0 || len(target.Size) > 0 {
					return db.Where("id IN ?", listAdTagFilter)
				}
				return db
			},
		).Find(&rows)

	return
}

func (t *InventoryAdTag) DeleteTagAfterDeleteDomain(domainId, userId int64) {
	//var listTag []int64
	//mysql.Client.Model(&InventoryAdTagRecord{}).Select("id").Where("inventory_id = ?").Find(&listTag)
	//for _, id := range listTag {
	//	new(AdTagSizeAdditional).DeleteAllByTagId(id)
	//}
	mysql.Client.Model(&InventoryAdTagRecord{}).Delete(&InventoryAdTagRecord{}, "inventory_id = ? and user_id = ?", domainId, userId)
}
