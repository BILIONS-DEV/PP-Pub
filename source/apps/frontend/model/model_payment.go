package model

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/syyongx/php2go"
	"gorm.io/gorm"
	"source/apps/frontend/lang"
	"source/apps/frontend/payload"
	"source/core/technology/mysql"
	"source/pkg/block"
	"source/pkg/datatable"
	"source/pkg/pagination"
	"source/pkg/utility"
	"strings"
	"time"
)

type Payment struct{}

type PaymentRecord struct {
	mysql.TablePaymentSubPub
}

func (PaymentRecord) TableName() string {
	return mysql.Tables.PaymentSubPub
}

func (t *Payment) GetByFilters(inputs *payload.InputPaymentFilter, user UserRecord) (response datatable.Response, err error) {
	var Payments []PaymentRecord
	var total int64
	err = mysql.Client.Where("user_id = ?", user.Id).
		Scopes(
			t.SetFilterStatus(inputs),
			t.setFilterSearch(inputs),
		).
		Model(&Payments).Count(&total).
		Scopes(
			t.setOrder(inputs),
			pagination.Paginate(pagination.Params{
				Limit:  inputs.Length,
				Offset: inputs.Start,
			}),
		).
		Find(&Payments).Error
	if err != nil {
		if !utility.IsWindow() {
			err = fmt.Errorf(lang.Translate.Errors.InventoryError.List.ToString())
		}
		return datatable.Response{}, err
	}
	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.MakeResponseDatatable(Payments)
	return
}

func (t *Payment) SetFilterStatus(inputs *payload.InputPaymentFilter) func(db *gorm.DB) *gorm.DB {
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

func (t *Payment) setFilterSearch(inputs *payload.InputPaymentFilter) func(db *gorm.DB) *gorm.DB {
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
		return db.Where("note LIKE ? OR amount LIKE ? OR billing LIKE ?", "%"+inputs.PostData.QuerySearch+"%", "%"+inputs.PostData.QuerySearch+"%", "%"+inputs.PostData.QuerySearch+"%")
	}
}

func (t *Payment) setOrder(inputs *payload.InputPaymentFilter) func(db *gorm.DB) *gorm.DB {
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

type PaymentRecordDatatable struct {
	PaymentRecord
	Period string `json:"period"`
	Amout  string `json:"amount"`
	// Method string `json:"method"`
	Info   string `json:"info"`
	Status string `json:"status"`
	Action string `json:"action"`
}

func (t *Payment) MakeResponseDatatable(Invoices []PaymentRecord) (records []PaymentRecordDatatable) {
	for _, Invoice := range Invoices {
		user := new(User).GetById(Invoice.UserId)
		if user.Id == 0 {
			continue
		}

		PaymentDateNow := false
		if Invoice.PaymentDate.Unix() < time.Now().Unix() {
			PaymentDateNow = true
		}
		HaveBeenPaid := false
		if Invoice.PaidDate.Unix() > 0 {
			HaveBeenPaid = true
		}
		RevenueShare, err := new(RevenueShare).GetByUserId(user.Id)
		if err != nil {
			continue
		}
		// startDate, _ := time.Parse("20060102", strconv.Itoa(Invoice.StartDate))
		// endDate, _ := time.Parse("20060102", strconv.Itoa(Invoice.EndDate))
		rec := PaymentRecordDatatable{
			Period: block.RenderToString("payment/index/block.period.gohtml", fiber.Map{
				"Month": Invoice.StartDate.Format("January 2006"),
			}),
			Amout: block.RenderToString("payment/index/block.amount.gohtml", fiber.Map{
				"Amount":       php2go.NumberFormat(Invoice.Amount, 2, ".", ","),
				"RevenueShare": RevenueShare.Rate,
			}),
			// Method: block.RenderToString("Payment/invoice/block.method.gohtml", Invoice),
			Info: block.RenderToString("payment/index/block.info.gohtml", fiber.Map{
				"PaymentDate":    Invoice.PaymentDate.Format("01/02/2006"),
				"PaymentDateNow": PaymentDateNow,
				"HaveBeenPaid":   HaveBeenPaid,
				"PaidDate":       Invoice.PaidDate.Format("01/02/2006"),
			}),
			Status: block.RenderToString("payment/index/block.status.gohtml", Invoice),
			Action: block.RenderToString("payment/index/block.action.gohtml", Invoice),
		}
		records = append(records, rec)
	}
	return
}

func (t *Payment) GetAll() (Payments []PaymentRecord) {
	mysql.Client.Find(&Payments)
	return
}

func (t *Payment) GetById(id int64) (record PaymentRecord, err error) {
	err = mysql.Client.Where("id = ?", id).Find(&record).Error
	if record.Id == 0 {
		err = errors.New("Record not found")
	}
	return
}

func (t *Payment) CheckRecord(id, userId int64) (row PaymentRecord, err error) {
	err = mysql.Client.Model(&PaymentRecord{}).Where("id = ? and user_id = ?", id, userId).Find(&row).Error
	return
}
