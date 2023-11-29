package payment

import (
	"github.com/gofiber/fiber/v2"
	"github.com/syyongx/php2go"
	"source/internal/entity/model"
	"source/internal/lang"
	"source/internal/repo"
	"source/pkg/block"
	"strconv"
	"strings"
	"time"
)

type UsercasePayment interface {
	GetByFiltersInvoice(params *model.ParamPaymentIndex, user model.User) (invoices []model.PaymentInvoicesIndex, total int64, err error)
	Pagination(params *model.ParamPaymentIndex, user model.User, ctx *fiber.Ctx) (PageData model.PageData)
	GetTotalInfoByInvoices(permission string) (TotalInfo string)
	GetTotalAmountByFilter(params *model.ParamPaymentIndex, user model.User) (totalAmount string)
	APIFiltersInvoicePublisher(params *model.ParamPaymentIndex, userLogin model.User) (Payments []model.PaymentInvoice, total int64, err error)
	APIGetTotalInfoByInvoices(permission string) (response ResponseTotalInfoInoivces)
	GetByID(id int64) (record model.PaymentInvoice, err error)
	GetRequestByInvoice(requestID []int64) (record []model.PaymemtRequest, err error)
	GetHistoryInvoiceByInvoice(invoiceID int64) (records []model.PaymentInvoice, err error)
	GetInvoicesByPublishers(publishers []string) (response []ResponseGetInvoicesByPublishers, err error)
	ReportAgencyInTime(Input InputReportAgencyInTime) (report model.ReportAgency, err error)
	SearchInvoice(keyword string) (invoices []ResponseSearch, err error)
}

type paymentUsecase struct {
	repos *repo.Repositories
	Trans *lang.Translation
}

func NewPaymentUsecase(repos *repo.Repositories, trans *lang.Translation) *paymentUsecase {
	return &paymentUsecase{repos: repos, Trans: trans}
}

func (t *paymentUsecase) GetByFiltersInvoice(params *model.ParamPaymentIndex, userLogin model.User) (invoices []model.PaymentInvoicesIndex, total int64, err error) {
	var Payments []model.PaymentInvoice
	Payments, total, err = t.repos.Payment.FiltersInvoice(params, userLogin)
	if err != nil || len(Payments) == 0 {
		return
	}

	// handle payment
	for _, Payment := range Payments {
		var paymentInvoice model.PaymentInvoicesIndex
		paymentInvoice.Invoice = Payment

		var user model.User
		user = t.repos.User.FindByID(Payment.UserID)
		if err != nil || user.ID == 0 {
			continue
		}
		var agency model.User
		if user.Presenter > 0 {
			agency = t.repos.User.FindByID(user.Presenter)
			if err != nil {
				continue
			}
		}
		billing := t.repos.Billing.GetByUserID(user.ID)
		RevenueShare := t.repos.RevenueShare.GetByUserID(user.ID)
		if err != nil {
			continue
		}
		PaymentDateNow := false
		if Payment.PaymentDate.Unix() < time.Now().Unix() && Payment.Status != 2 {
			PaymentDateNow = true
		}
		HaveBeenPaid := false
		if Payment.PaidDate.Unix() > 0 {
			HaveBeenPaid = true
		}
		paymentInvoice.Id = strconv.FormatInt(Payment.ID, 10)
		paymentInvoice.Period = block.RenderToString("payment/invoice/block.period.gohtml", fiber.Map{
			"Period": Payment.StartDate.Format("January 2006"),
		})
		paymentInvoice.Amount = block.RenderToString("payment/invoice/block.amount.gohtml", fiber.Map{
			"Amount":       php2go.NumberFormat(Payment.Amount, 2, ".", ","),
			"RevenueShare": RevenueShare.Rate,
		})
		paymentInvoice.User = block.RenderToString("payment/invoice/block.user.gohtml", fiber.Map{
			"User":   user,
			"Agency": agency,
		})
		paymentInvoice.Method = block.RenderToString("payment/invoice/block.method.gohtml", fiber.Map{
			"Payment": Payment,
			"Billing": billing,
		})
		paymentInvoice.Status = block.RenderToString("payment/invoice/block.status.gohtml", Payment)
		// payment.Action = block.RenderToString("payment/invoice/block.action.gohtml", Payment)
		paymentInvoice.Action = block.RenderToString("payment/invoice/block.action.gohtml", fiber.Map{
			"Payment": Payment,
			"IdUser":  userLogin.ID,
		})
		paymentInvoice.Info = block.RenderToString("payment/invoice/block.info.gohtml", fiber.Map{
			"PaymentDate":    Payment.PaymentDate.Format("02/01/2006"),
			"PaymentDateNow": PaymentDateNow,
			"HaveBeenPaid":   HaveBeenPaid,
			"PaidDate":       Payment.PaidDate.Format("02/01/2006"),
		})

		invoices = append(invoices, paymentInvoice)
	}
	return
}

func (t *paymentUsecase) GetTotalInfoByInvoices(permission string) (TotalInfo string) {
	var total int
	var amount float64
	var pending int
	var amountPending float64

	invoices := t.repos.Payment.GetTotalInfoInvoicesByPermission(permission)
	for _, value := range invoices {
		total = total + 1
		amount = amount + value.Amount

		if value.Status == model.StatusPaymentPending {
			pending = pending + 1
			amountPending = amountPending + value.Amount
		}
	}

	TotalInfo = block.RenderToString("payment/invoice/block.total.gohtml", fiber.Map{
		"Total":         strconv.Itoa(total),
		"Amount":        php2go.NumberFormat(amount, 2, ".", ","),
		"Pending":       strconv.Itoa(pending),
		"AmountPending": php2go.NumberFormat(amountPending, 2, ".", ","),
	})
	return
}

type ResponseTotalInfoInoivces struct {
	Total         int     `json:"total"`
	Amount        float64 `json:"total_amount"`
	Pending       int     `json:"pending"`
	AmountPending float64 `json:"total_amount_pending"`
}

func (t *paymentUsecase) APIGetTotalInfoByInvoices(permission string) (response ResponseTotalInfoInoivces) {
	invoices := t.repos.Payment.GetTotalInfoInvoicesByPermission(permission)
	for _, value := range invoices {
		response.Total = response.Total + 1
		response.Amount = response.Amount + value.Amount

		if value.Status == model.StatusPaymentPending {
			response.Pending = response.Pending + 1
			response.AmountPending = response.AmountPending + value.Amount
		}
	}
	return
}

func (t *paymentUsecase) GetTotalAmountByFilter(params *model.ParamPaymentIndex, user model.User) (totalAmount string) {
	Payments := t.repos.Payment.GetTotalAmountByFilter(params, user)
	if len(Payments) == 0 {
		return
	}
	var total float64
	for _, value := range Payments {
		total = total + value.Amount
	}
	totalAmount = php2go.NumberFormat(total, 2, ".", ",")
	return
}

func (t *paymentUsecase) Pagination(params *model.ParamPaymentIndex, user model.User, ctx *fiber.Ctx) (PageData model.PageData) {
	page, _ := strconv.Atoi(params.Page)
	if page < 1 {
		page = 1
	}
	Totalpages := t.repos.Payment.GetTotalPagesInvoice(params, user)

	filter := ctx.OriginalURL()
	baseQueries := ctx.Context().QueryArgs()
	len := baseQueries.Len()
	pageNow := string(baseQueries.Peek("page"))
	if page != 0 {
		filter = strings.ReplaceAll(filter, "?page="+params.Page, "")
		filter = strings.ReplaceAll(filter, "&page="+params.Page, "")
	}

	var character string

	if len == 0 || len == 1 && pageNow != "" {
		character = "?"
		// Firtp = filter + "&page=1"
	} else {
		character = "&"
		// Firtp = filter + "?page=1"
	}

	listPage, lastPage := t.CalculatorPages(page, Totalpages, filter, character)
	PageData.FirstPage = filter + character + "page=1"
	PageData.ListPage = listPage
	PageData.LastPage = lastPage
	return
}

func (t *paymentUsecase) CalculatorPages(page int, totalRecords int64, filter string, character string) ([]model.Pages, string) {
	var (
		listPage        []model.Pages
		pages, min, max int
	)
	totalPagesInt := totalRecords / model.PageSize
	totalPages := float64(totalRecords) / float64(model.PageSize)

	if totalPages > float64(totalPagesInt) {
		pages = int(totalPagesInt) + 1
	} else {
		pages = int(totalPagesInt)
	}
	if model.PageShow >= pages {
		min = 1
		max = pages
	} else {
		if page <= model.PageShow/2 {
			min = 1
			max = model.PageShow
		} else {
			if (pages - page) <= model.PageShow/2 {
				min = pages - model.PageShow
				max = pages
			} else {
				min = page - model.PageShow/2
				max = page + model.PageShow/2
			}
		}
	}
	for i := min; i <= max; i++ {
		// iString := Convert(i, filter, character)
		iString := filter + character + "page=" + strconv.Itoa(i)
		tempPage := model.Pages{
			LinkPage: iString,
			Number:   i,
		}
		listPage = append(listPage, tempPage)
	}
	// lastPage := Convert(pages, filter, character)
	lastPage := filter + character + "page=" + strconv.Itoa(pages)
	return listPage, lastPage
}

func (t *paymentUsecase) APIFiltersInvoicePublisher(params *model.ParamPaymentIndex, userLogin model.User) (Payments []model.PaymentInvoice, total int64, err error) {
	Payments, total, err = t.repos.Payment.FiltersInvoice(params, userLogin)
	return
}

func (t *paymentUsecase) GetByID(id int64) (record model.PaymentInvoice, err error) {
	record, err = t.repos.Payment.GetByID(id)
	return
}

func (t *paymentUsecase) GetRequestByInvoice(requestID []int64) (record []model.PaymemtRequest, err error) {
	record, err = t.repos.Payment.GetRequestByInvoice(requestID)
	return
}

func (t *paymentUsecase) GetHistoryInvoiceByInvoice(invoiceID int64) (records []model.PaymentInvoice, err error) {
	records, err = t.repos.Payment.GetHistoryInvoiceByInvoice(invoiceID)
	return
}

type ResponseGetInvoicesByPublishers struct {
	Email            string  `json:"email"`
	PaymentDetail    string  `json:"payment_detail"`
	PaymentType      string  `json:"payment_type"`
	AvailablePayment float64 `json:"available_payment"`
}

func (t *paymentUsecase) GetInvoicesByPublishers(publishers []string) (response []ResponseGetInvoicesByPublishers, err error) {
	users := t.repos.User.GetByEmails(publishers)
	if len(users) == 0 {
		return
	}
	var userIDs []int64
	for _, value := range users {
		userIDs = append(userIDs, value.ID)
	}

	for _, user := range users {
		var invoicies []model.PaymentInvoice
		invoicies, err = t.repos.Payment.GetInvoiceByPublisher(user.ID)
		if err != nil {
			return
		}
		if len(invoicies) == 0 {
			continue
		}
		billing := t.repos.Billing.GetByUserID(user.ID)
		var paymentDetail string
		var paymentType string
		switch billing.Method {
		case "bank":
			paymentDetail = billing.BankAccountNumber
			paymentType = "Wire Transfer"
			break
		case "currency":
			paymentDetail = billing.WalletId
			paymentType = billing.Cryptocurrency
			break
		case "paypal":
			paymentDetail = billing.PaypalEmail
			paymentType = billing.Method
		case "payoneer":
			paymentDetail = billing.PayOneerEmail
			paymentType = billing.Method
			break
		}
		var availablePayment float64
		for _, invoice := range invoicies {
			if invoice.Status != model.StatusPaymentPending {
				continue
			}
			sub := invoice.PaymentDate.Sub(time.Now())
			if sub > 0 {
				continue
			}
			availablePayment += invoice.Amount
		}
		if availablePayment > 0 {
			response = append(response, ResponseGetInvoicesByPublishers{
				Email:            user.Email,
				PaymentDetail:    paymentDetail,
				PaymentType:      paymentType,
				AvailablePayment: availablePayment,
			})
		}
	}
	return
}

type InputReportAgencyInTime struct {
	Agency model.User
	Month  string
}

func (t *paymentUsecase) ReportAgencyInTime(Input InputReportAgencyInTime) (report model.ReportAgency, err error) {
	report, err = t.repos.Payment.ReportAgencyInTime(Input.Agency, Input.Month, 0)
	return
}

type ResponseSearch struct {
	ID             int64   `json:"id"`
	System         string  `json:"system"`
	Period         string  `json:"period"`
	User           string  `json:"user"`
	Permission     string  `json:"permission"`
	AccountManager string  `json:"account_manager"`
	Amount         float64 `json:"amount"`
	Status         string  `json:"status"`
	Method         string  `json:"method"`
	Month          string  `json:"month"`
	Note           string  `json:"note"`
	PaymentDate    string  `json:"payment_date"`
	PaidDate       string  `json:"paid_date"`
}

func (t *paymentUsecase) SearchInvoice(keyword string) (invoices []ResponseSearch, err error) {
	// get users by keyword
	var users []model.User
	users, err = t.repos.User.GetBySearch(keyword)
	var userIDs []int64
	for _, value := range users {
		userIDs = append(userIDs, value.ID)
	}

	// get invoice by user
	var Payments []model.PaymentInvoice
	Payments, err = t.repos.Payment.GetInvoiceByUsers(userIDs)
	if len(Payments) == 0 {
		return
	}
	for _, value := range Payments {
		user := t.repos.User.FindByID(value.UserID)
		accountManager := t.repos.User.FindByID(user.Presenter)
		var method string
		var billing model.UserBilling
		if value.Status == model.StatusPaymentDone {
			method = value.Method
		} else {
			billing = t.repos.Billing.GetByUserID(value.UserID)
			if billing.Id != 0 {
				method = billing.Method
			}
		}
		var status string
		if value.Status == model.StatusPaymentDone {
			status = "paid"
		}
		switch value.Status {
		case model.StatusPaymentDone:
			status = "paid"
			break
		case model.StatusPaymentPending:
			status = "pending"
			break
		case model.StatusPaymentRject:
			status = "rejected"
			break
		}
		var Month, PaymentDate, PaidDate string
		if !model.CheckTimeEmpty(value.EndDate) {
			Month = value.EndDate.Format("200601")
		}
		if !model.CheckTimeEmpty(value.PaymentDate) {
			PaymentDate = value.PaymentDate.Format("02/01/2006")
		}
		if !model.CheckTimeEmpty(value.PaidDate) {
			PaidDate = value.PaidDate.Format("02/01/2006")
		}
		var permission = "publisher"
		if user.Permission == model.UserPermissionSale {
			permission = "agency"
		}
		invoices = append(invoices, ResponseSearch{
			ID:             value.ID,
			System:         "Pubpower",
			Period:         value.StartDate.Format("02/01/2006") + " - " + value.EndDate.Format("02/01/2006"),
			User:           user.Email,
			Permission:     permission,
			AccountManager: accountManager.Email,
			Amount:         value.Amount,
			Status:         status,
			Method:         method,
			Month:          Month,
			Note:           value.Note,
			PaymentDate:    PaymentDate,
			PaidDate:       PaidDate,
		})
	}
	return
}
