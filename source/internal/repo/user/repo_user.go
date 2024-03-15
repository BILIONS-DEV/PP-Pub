package user

import (
	"fmt"
	"gorm.io/gorm"
	"source/infrastructure/caching"
	"source/infrastructure/mysql"
	"source/internal/entity/dto"
	"source/internal/entity/model"
	"strings"
)

type RepoUser interface {
	FindByID(ID int64, fields ...string) (record model.User)
	FindByLogin(email, passwordHash string, fields ...string) (record model.User, err error)
	FindByLoginToken(loginToken string, fields ...string) (record model.User, err error)
	FindByFilter(inputs *dto.UserFilterPayload) (totalRecord int64, records []model.User, err error)
	CheckEmail(email string) bool
	Save(record model.User) (err error)
	GetNewUser() (record model.User)
	GetByLoginToken(loginToken string) (user model.User)
	GetBillingByUser(userID int64) (billing model.UserBilling, err error)
	SaveBilling(billing *model.UserBilling) (err error)
	GetByEmail(email string) (record model.User)
	GetByEmails(emails []string) (records []model.User)
	UpdateAdsTxtPublisher(id int64, adsTxt string) (err error)
	ResetCacheAll(userId int64) (err error)
	GetUsersByPermission(permissions []int64) (records []model.User, err error)
	GetBySearch(keyword string) (users []model.User, err error)
	GetInfoByUserID(ID int64) (record model.TableUserInfo)
}

type user struct {
	DB    *gorm.DB
	Cache caching.Cache
}

func NewUserRepo(DB *gorm.DB, cache caching.Cache) *user {
	return &user{DB: DB, Cache: cache}
}

func (t *user) FindByFilter(inputs *dto.UserFilterPayload) (totalRecord int64, records []model.User, err error) {
	err = t.DB.
		Scopes(
			t.setFilterStatus(inputs),
			t.setFilterSearch(inputs),
		).
		Model(&records).Count(&totalRecord).
		Scopes(
			t.setOrder(inputs),
			mysql.Paginate(mysql.Deps{
				Limit:  inputs.Length,
				Offset: inputs.Start,
			}),
		).
		Find(&records).Error
	return
}

func (t *user) setFilterStatus(inputs *dto.UserFilterPayload) func(db *gorm.DB) *gorm.DB {
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

func (t *user) setFilterSearch(inputs *dto.UserFilterPayload) func(db *gorm.DB) *gorm.DB {
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
		return db.Where("email LIKE ?", "%"+inputs.PostData.QuerySearch+"%").
			Or("first_name LIKE ?", "%"+inputs.PostData.QuerySearch+"%").
			Or("last_name LIKE ?", "%"+inputs.PostData.QuerySearch+"%").
			Or("address LIKE ?", "%"+inputs.PostData.QuerySearch+"%")
	}
}

func (t *user) setOrder(inputs *dto.UserFilterPayload) func(db *gorm.DB) *gorm.DB {
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
		return db.Order("id desc")
	}
}

func (t *user) FindByID(ID int64, fields ...string) (record model.User) {
	if ID < 1 {
		return
	}
	query := t.DB
	if len(fields) > 0 {
		query = query.Select(fields)
	}
	query.Take(&record, ID)
	return
}

func (t *user) GetByEmail(email string) (record model.User) {
	t.DB.Model(&model.User{}).Where("email = ?", email).Find(&record)
	return
}

func (t *user) GetInfoByUserID(ID int64) (record model.TableUserInfo) {
	t.DB.Model(&model.TableUserInfo{}).Where("user_id = ?", ID).Find(&record)
	return
}

func (t *user) GetByEmails(emails []string) (records []model.User) {
	t.DB.Model(&model.User{}).Where("email in ?", emails).Find(&records)
	return
}

func (t *user) UpdateAdsTxtPublisher(id int64, adsTxt string) (err error) {
	var rec = model.User{
		AdsTxtCustomByAdmin: model.TYPEAdsTxtCustom(adsTxt),
	}
	return t.DB.Model(&model.User{}).Where("id = ?", id).Updates(&rec).Error
}

func (t *user) ResetCacheAll(userId int64) (err error) {
	return t.DB.Model(&model.InventoryModel{}).Where("user_id = ?", userId).Debug().Update("render_cache", 1).Error
}

func (t *user) FindByLogin(email, password string, fields ...string) (record model.User, err error) {
	query := t.DB.Where("email = ? AND password = ?", email, model.HashPassword(password))
	if len(fields) > 0 {
		query = query.Select(fields)
	}
	err = query.Take(&record).Error
	return
}

func (t *user) FindByLoginToken(loginToken string, fields ...string) (record model.User, err error) {
	query := t.DB.Where("login_token = ?", loginToken)
	if len(fields) > 0 {
		query = query.Select(fields)
	}
	err = query.Find(&record).Error
	return
}

func (t *user) Save(record model.User) (err error) {
	err = t.DB.Save(&record).Error
	return
}

func (t *user) CheckEmail(email string) bool {
	var rec model.User
	t.DB.
		Model(&model.User{}).
		Where("email = ?", email).
		Find(&rec)
	if rec.ID == 0 {
		return false
	}
	return true
}

func (t *user) GetNewUser() (record model.User) {
	t.DB.
		Model(&model.User{}).
		Order("id DESC").
		Find(&record)
	return
}

func (t *user) GetByLoginToken(loginToken string) (record model.User) {
	t.DB.
		Model(&model.User{}).
		Where("login_token = ?", loginToken).
		Find(&record)
	return
}

func (t *user) GetBillingByUser(userID int64) (billing model.UserBilling, err error) {
	err = t.DB.
		Model(&model.UserBilling{}).
		Where("user_id = ?", userID).
		Find(&billing).Error
	return
}

func (t *user) SaveBilling(billing *model.UserBilling) (err error) {
	err = t.DB.Save(billing).Error
	return
}

func (t *user) GetUsersByPermission(permissions []int64) (records []model.User, err error) {
	err = t.DB.Model(&model.User{}).
		Where("permission in ?", permissions).
		Find(&records).Error
	return
}

func (t *user) GetBySearch(keyword string) (users []model.User, err error) {
	err = t.DB.Model(&model.User{}).Where("email like ?", "%"+keyword+"%").Find(&users).Error
	return
}
