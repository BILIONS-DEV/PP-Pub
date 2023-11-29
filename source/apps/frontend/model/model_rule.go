package model

import (
	"fmt"
	"source/apps/frontend/payload"
	"source/core/technology/mysql"
	"source/pkg/datatable"
	"source/pkg/htmlblock"
	"source/pkg/pagination"
	"source/pkg/utility"
	"strings"

	"gorm.io/gorm"
)

type Rule struct{}

type RuleRecord struct {
	mysql.TableRule
}

func (t *Rule) GetByFilters(inputs *payload.RuleFilterPayload, user UserRecord) (response datatable.Response, err error) {
	var array []string
	var records []payload.Rule
	var total int64
	array = append(array, "(SELECT t1.id, t1.user_id, t1.name as name, t1.created_at, t1.deleted_at,'floor' as type, status FROM floor t1 )")
	array = append(array, "(SELECT t2.id, t2.user_id, t2.restriction_name as name, t2.created_at, t2.deleted_at, 'blocking' as type, 1 as status FROM blocking t2 )")
	array = append(array, "(SELECT t3.id, t3.user_id, t3.name, t3.created_at, t3.deleted_at, 'blocked-page' as type, status FROM rule t3 )")
	err = mysql.Client.Select("*").
		Table("("+strings.Join(array, " UNION ALL ")+") t").
		Scopes(
			t.setFilterSearch(inputs),
			t.setFilterType(inputs),
		).
		Where("t.user_id = ?", user.Id).
		Where("t.deleted_at is null").
		Count(&total).
		Scopes(
			t.setOrder(inputs),
			pagination.Paginate(pagination.Params{
				Limit:  inputs.Length,
				Offset: inputs.Start,
			}),
		).
		Find(&records).Error
	if err != nil {
		if !utility.IsWindow() {
			err = fmt.Errorf("error")
		}
		return datatable.Response{}, err
	}

	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.MakeResponseDatatable(records)
	return
}

func (t *Rule) setFilterSearch(inputs *payload.RuleFilterPayload) func(db *gorm.DB) *gorm.DB {
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
		return db.Where("t.name LIKE ?", "%"+inputs.PostData.QuerySearch+"%")
	}
}

func (t *Rule) setFilterType(inputs *payload.RuleFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.Type != nil {
			switch inputs.PostData.Type.(type) {
			case string, int:
				if inputs.PostData.Type != "" {
					return db.Where("t.type = ?", inputs.PostData.Type)
				}
			case []string, []interface{}:
				return db.Where("t.type IN ?", inputs.PostData.Type)
			}
		}
		return db
	}
}

func (t *Rule) setOrder(inputs *payload.RuleFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(inputs.Order) > 0 {
			var orders []string
			for _, order := range inputs.Order {
				column := inputs.Columns[order.Column]
				orders = append(orders, fmt.Sprintf("t.%s %s", column.Data, order.Dir))
			}
			orderString := strings.Join(orders, ", ")
			return db.Order(orderString)
		}
		return db.Order("t.created_at desc")
	}
}

type RuleRecordDatatable struct {
	RowId  int64  `json:"DT_RowId"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Action string `json:"action"`
}

func (t *Rule) MakeResponseDatatable(records []payload.Rule) (responses []RuleRecordDatatable) {
	for _, rec := range records {
		resp := RuleRecordDatatable{
			RowId:  rec.ID,
			Name:   rec.Name,
			Type:   htmlblock.Render("rule/index/datatable/type.gohtml", rec).String(),
			Status: htmlblock.Render("rule/index/datatable/status.gohtml", rec).String(),
			Action: htmlblock.Render("rule/index/datatable/action.gohtml", rec).String(),
		}
		responses = append(responses, resp)
	}
	return responses
}

func (t *Rule) Create(record *RuleRecord) (err error) {
	err = mysql.Client.Create(record).Error
	return
}

func (t *Rule) Update(record *RuleRecord) (err error) {
	err = mysql.Client.Save(record).Error
	return
}

func (t *Rule) GetByUser(userId int64) (records []RuleRecord) {
	mysql.Client.Where("user_id = ?", userId).Find(records)
	return
}
