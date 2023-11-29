package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/syyongx/php2go"
	"gorm.io/gorm"
	"source/infrastructure/caching"
	"source/internal/entity/model"
	"strconv"
	"strings"
	"time"
)

const (
	TIME_ZONE_DEFAULT = "America/New_York"
)

type RepoPayment interface {
	FiltersInvoice(inputs *model.ParamPaymentIndex, user model.User) (invoices []model.PaymentInvoice, total int64, err error)
	GetTotalPagesInvoice(inputs *model.ParamPaymentIndex, user model.User) (total int64)
	GetTotalInfoInvoicesByPermission(permission string) (invoices []model.PaymentInvoice)
	GetTotalAmountByFilter(inputs *model.ParamPaymentIndex, user model.User) (invoices []model.PaymentInvoice)
	GetByID(id int64) (record model.PaymentInvoice, err error)
	GetRequestByInvoice(requestID []int64) (record []model.PaymemtRequest, err error)
	GetHistoryInvoiceByInvoice(invoiceID int64) (records []model.PaymentInvoice, err error)
	GetInvoiceByPublisher(userID int64) (Payments []model.PaymentInvoice, err error)
	ReportAgencyInTime(agency model.User, month string, numberTimeOut int) (report model.ReportAgency, err error)
	GetInvoiceByUsers(userIDs []int64) (Payments []model.PaymentInvoice, err error)
}

type paymentRepo struct {
	Db    *gorm.DB
	Cache caching.Cache
}

func NewPaymentRepo(db *gorm.DB) *paymentRepo {
	return &paymentRepo{Db: db}
}

func (t *paymentRepo) FiltersInvoice(inputs *model.ParamPaymentIndex, user model.User) (invoices []model.PaymentInvoice, total int64, err error) {
	// var Payments []model.PaymentInvoice
	// var Payments []PaymentInvoiceRecord
	err = t.Db.
		Model(&model.PaymentInvoice{}).
		Where("amount >= 0.01").
		Scopes(
			t.SetFilterStatusInvoice(inputs),
			t.SetFilterMoneyInvoice(inputs),
			t.setFilterSearchInvoice(inputs),
			t.setFilterPaidDateInvoice(inputs),
			t.setFilterMonthInvoice(inputs),
			t.setFilterPermission(inputs, user),
		).
		Count(&total).
		Scopes(
			t.setLimitInvoice(inputs),
			t.setOrderInvoice(inputs),
		).
		Find(&invoices).Error
	return
}

func (t *paymentRepo) SetFilterStatusInvoice(inputs *model.ParamPaymentIndex) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.Status != nil && len(inputs.Status) > 0 {
			return db.Where("payment_invoice.status IN ?", inputs.Status)
		}
		return db
	}
}
func (t *paymentRepo) setFilterPaidDateInvoice(inputs *model.ParamPaymentIndex) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PaidDate != "" {
			t, _ := time.Parse("01/02/2006", inputs.PaidDate)
			paid := t.Format("2006-01-02")
			return db.Where("paid_date = ?", paid)
		}
		return db
	}
}

func (t *paymentRepo) setFilterMonthInvoice(inputs *model.ParamPaymentIndex) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.Month != "" {
			t, _ := time.Parse("Jan 2006", inputs.Month)
			month := t.Format("2006-01")
			if inputs.Sort != "" {
				switch inputs.Sort {
				case "period":
					return db.Where("start_date LIKE ? OR end_date LIKE ?", month+"%", month+"%")
				case "payment_date":
					return db.Where("payment_date LIKE  ?", month+"%")
				case "paid_date":
					return db.Where("paid_date LIKE  ?", month+"%")
				}
			} else {
				return db.Where("payment_date LIKE  ?", month+"%")
			}
		}
		return db
	}
}

func (t *paymentRepo) setFilterPermission(inputs *model.ParamPaymentIndex, user model.User) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.Permission != "" {
			switch inputs.Permission {
			case "sale":
				if user.Permission.String() == "sale" {
					return db.Joins("INNER JOIN user ON user.id = payment_invoice.user_id").Where("user.permission IN (?) AND user.email = ?", model.UserPermissionSale, user.Email)
					break
				}
				return db.Joins("INNER JOIN user ON user.id = payment_invoice.user_id").Where("user.permission IN (?)", model.UserPermissionSale)
				break
			case "publisher":
				if user.Permission.String() == "sale" {
					return db.Joins("INNER JOIN user ON user.id = payment_invoice.user_id").Where("user.permission IN (?,?,?) AND user.presenter = ?", model.UserPermissionMember, model.UserPermissionManagedService, model.UserPermissionPublisher, user.ID)
					break
				}
				return db.Joins("INNER JOIN user ON user.id = payment_invoice.user_id").Where("user.permission IN (?,?,?)", model.UserPermissionMember, model.UserPermissionManagedService, model.UserPermissionPublisher)
				break
			case "affiliate":
				return db
				break
			}
		}
		return db
	}
}

func (t *paymentRepo) SetFilterMoneyInvoice(inputs *model.ParamPaymentIndex) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if inputs.FromMoney != "" && inputs.ToMoney != "" {
			return db.Where("payment_invoice.amount >= ? and payment_invoice.amount <= ?", inputs.FromMoney, inputs.ToMoney)
		} else {
			if inputs.FromMoney != "" {
				return db.Where("payment_invoice.amount >= ?", inputs.FromMoney)
			}
			if inputs.ToMoney != "" {
				return db.Where("payment_invoice.amount <= ?", inputs.ToMoney)
			}
		}

		return db
	}
}

func (t *paymentRepo) setFilterSearchInvoice(inputs *model.ParamPaymentIndex) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var flag bool
		// Search from form of datatable <- not use
		if inputs.Search != nil && inputs.Search.Value != "" {
			flag = true
		}
		// Search from form filter
		if inputs.QuerySearch != "" {
			flag = true
		}
		if !flag {
			return db
		}
		if inputs.QuerySearch != "" && inputs.Permission != "" {
			keyword := php2go.Trim(inputs.QuerySearch)
			return db.Where("user.email like ?", "%"+keyword+"%").Where("payment_invoice.note LIKE ? OR payment_invoice.amount LIKE ? OR payment_invoice.billing LIKE ? OR user.email like ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
		} else if inputs.QuerySearch != "" {
			keyword := php2go.Trim(inputs.QuerySearch)
			return db.Joins("inner join user on user.id = payment_invoice.user_id", t.Db.Where("email like ?", "%"+keyword+"%")).Where("payment_invoice.note LIKE ? OR payment_invoice.amount LIKE ? OR payment_invoice.billing LIKE ? OR user.email like ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
		}
		return db
	}
}

func (t *paymentRepo) setLimitInvoice(inputs *model.ParamPaymentIndex) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var size int
		if inputs.Size != "" {
			size, _ = strconv.Atoi(inputs.Size)
		} else {
			size = model.PageSize
		}
		if inputs.Page != "" {
			p, _ := strconv.Atoi(inputs.Page)
			if p < 1 {
				p = 1
			}
			skip := (p - 1) * size
			limitpage := size
			return db.Limit(limitpage).Offset(skip)
		} else {
			p := 1
			skip := (p - 1) * size
			limitpage := size
			return db.Limit(limitpage).Offset(skip)
		}
		return db
	}
}

func (t *paymentRepo) setOrderInvoice(inputs *model.ParamPaymentIndex) func(db *gorm.DB) *gorm.DB {
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
		return db.Order("end_date DESC, amount DESC")
	}
}

func (t *paymentRepo) GetTotalPagesInvoice(inputs *model.ParamPaymentIndex, user model.User) (total int64) {
	var invoices []model.PaymentInvoice
	re := t.Db.Model(&invoices).
		Where("amount >= 0.01").
		Scopes(
			t.SetFilterStatusInvoice(inputs),
			t.SetFilterMoneyInvoice(inputs),
			t.setFilterSearchInvoice(inputs),
			t.setFilterPaidDateInvoice(inputs),
			t.setFilterMonthInvoice(inputs),
			t.setFilterPermission(inputs, user),
		).
		Find(&invoices)

	total = re.RowsAffected
	return
}

func (t *paymentRepo) GetTotalInfoInvoicesByPermission(permission string) (invoices []model.PaymentInvoice) {

	query := t.Db.
		Model(model.PaymentInvoice{}).
		Joins("INNER JOIN user ON payment_invoice.user_id = user.id")
	// Where("payment_invoice.status = ?", model.StatusPaymentPending)
	switch permission {
	case "publisher":
		query.Where("permission in (?,?,?)", model.UserPermissionPublisher, model.UserPermissionMember, model.UserPermissionManagedService)
		break
	case "sale":
		query.Where("permission in (?)", model.UserPermissionSale)
		break
	default:
		query.Where("permission in (?)", -1)
		break
	}
	query.Find(&invoices)
	return
}

func (t *paymentRepo) GetTotalAmountByFilter(inputs *model.ParamPaymentIndex, user model.User) (Payments []model.PaymentInvoice) {
	if inputs.QuerySearch == "" {
		return
	}
	t.Db.Where("amount >= 0.01").
		Scopes(
			t.SetFilterStatusInvoice(inputs),
			t.SetFilterMoneyInvoice(inputs),
			t.setFilterSearchInvoice(inputs),
			t.setFilterPaidDateInvoice(inputs),
			t.setFilterMonthInvoice(inputs),
			t.setFilterPermission(inputs, user),
		).
		Model(&Payments).
		Find(&Payments)
	return
}

func (t *paymentRepo) GetByID(id int64) (record model.PaymentInvoice, err error) {
	err = t.Db.Model(&model.PaymentInvoice{}).
		Where("id = ?", id).
		Find(&record).Error
	return
}

func (t *paymentRepo) GetRequestByInvoice(requestID []int64) (records []model.PaymemtRequest, err error) {
	err = t.Db.Model(&model.PaymemtRequest{}).
		Where("id IN ?", requestID).
		Find(&records).Error
	return
}

func (t *paymentRepo) GetHistoryInvoiceByInvoice(invoiceID int64) (records []model.PaymentInvoice, err error) {
	invoice := model.PaymentInvoice{}
	err = t.Db.Model(&invoice).
		Where("id = ?", invoiceID).
		Find(&invoice).Error
	if err != nil || invoice.ID == 0 {
		return
	}
	err = t.Db.
		Where("status = ? and user_id = ? and id != ?", model.StatusPaymentDone, invoice.UserID, invoiceID).
		Model(&model.PaymentInvoice{}).Order("id DESC").Limit(10).
		Find(&records).Error
	return
}

func (t *paymentRepo) GetInvoiceByPublisher(userID int64) (Payments []model.PaymentInvoice, err error) {
	err = t.Db.Where("amount >= 0.01 AND user_id = ?", userID).
		Model(&Payments).
		Find(&Payments).Error
	return
}

func (t *paymentRepo) ReportAgencyInTime(agency model.User, month string, numberTimeOut int) (report model.ReportAgency, err error) {
	query := "SELECT SUM(grossRevenue) as grossRevenue, SUM(netRevenue) as netRevenue FROM \"pw_bidder\" WHERE owner = 'apac' AND accManageId = " + strconv.FormatInt(agency.ID, 10) + " and __time Like '" + month + "%'"
	body, err := druid(query)
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout") {
			numberTimeOut += 1
		}
		if numberTimeOut < 5 {
			numberTimeOut += 1
			report, err = t.ReportAgencyInTime(agency, month, numberTimeOut)
		}
		return
	}
	var reports []model.ReportAgency
	json.Unmarshal([]byte(body), &reports)
	if len(reports) == 0 {
		return
	}
	report = reports[0]
	return
}

// func druid(query string, output BidderStatus) (Output BidderStatus, errs []ajax.Error) {
func druid(query string) (body []byte, err error) {
	URL := "https://query.vliplatform.com/druid/v2/sql/"
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("AUTHORIZATION", "Basic cmVhZG1hbjpMb3hpVVNKISEoODIhQA==")
		r.Headers.Set("Content-Type", "application/json;charset=UTF-8")
	})

	// On every a element which has href attribute call callback
	c.OnResponse(func(r *colly.Response) {
		body = r.Body
	})
	c.OnError(func(r *colly.Response, errHandle error) {
		if errHandle != nil {
			err = errHandle
			return
		} else {
			if r.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("status code error: %d", r.StatusCode))
				return
			}
		}
	})
	message := map[string]interface{}{
		"query": query,
		"context": map[string]interface{}{
			"useCache":      "true",
			"populateCache": "true",
			"sqlTimeZone":   TIME_ZONE_DEFAULT,
		},
	}
	payload, err := json.Marshal(message)
	if err != nil {
		return
	}
	errVisit := c.PostRaw(URL, payload)
	if err == nil {
		err = errVisit
	}
	c.Wait()

	return
}

func (t *paymentRepo) GetInvoiceByUsers(userIDs []int64) (Payments []model.PaymentInvoice, err error) {
	err = t.Db.
		Model(&model.PaymentInvoice{}).
		Where("user_id in ?", userIDs).
		Find(&Payments).Error
	return
}
