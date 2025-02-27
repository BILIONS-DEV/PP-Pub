package controller

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/config/assign"
	"source/apps/frontend/model"
	"source/apps/frontend/payload"
	"source/apps/frontend/view"
)

type Payment struct{}

type AssignPayment struct {
	assign.Schema
	Users  []model.UserRecord
	Params payload.ParamPaymentIndex
}

// func (t *Payment) Index(ctx *fiber.Ctx) error {
// 	assigns := AssignPayment{Schema: assign.Get(ctx)}
// 	assigns.Title = config.TitleWithPrefix("Payments")
// 	assigns.Users = new(model.User).GetAll()
// 	assigns.Inventories, _ = new(model.Inventory).GetByUserIdLimit(1, 10)
// 	//assigns.TemplatePath = "/payment"
// 	return ctx.Render("payment/index", assigns, view.LAYOUTMain)
// }

func (t *Payment) Index(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPayment)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	params := payload.ParamPaymentIndex{}
	if err := ctx.QueryParser(&params); err != nil {
		return err
	}
	assigns := AssignPayment{Schema: assign.Get(ctx)}

	assigns.Params = params
	assigns.Title = config.TitleWithPrefix("Payment")
	return ctx.Render("payment/index", assigns, view.LAYOUTMain)
}

func (t *Payment) IndexFilter(ctx *fiber.Ctx) error {
	userLogin := GetUserLogin(ctx)
	userAdmin := GetUserAdmin(ctx)
	isAccept := new(model.User).CheckUserLogin(userLogin, userAdmin, config.URIPayment)
	if !isAccept {
		return ctx.SendStatus(fiber.StatusNotFound)
	}
	// Get Input post
	var inputs payload.InputPaymentFilter
	if err := ctx.BodyParser(&inputs); err != nil {
		return err
	}
	// // Get data from model
	dataTable, err := new(model.Payment).GetByFilters(&inputs, userLogin)
	if err != nil {
		return err
	}
	// // Print Data to JSON
	return ctx.JSON(dataTable)
}

type AssignPaymentInvoice struct {
	assign.Schema
	Row                model.PaymentRecord
	Billing            model.UserBillingRecord
	User               model.UserRecord
	Today              string
	PaidDate           string
	TotalAmountReqeust string
	// RevenueShare int64
}

// func (t *Payment) Preivew(ctx *fiber.Ctx) error {
// 	user := GetUserLogin(ctx)
// 	if !user.IsFound() {
// 		return ctx.SendStatus(fiber.StatusNotFound)
// 	}
// 	assigns := AssignPaymentInvoice{Schema: assign.Get(ctx)}
// 	id := ctx.Query("invoiceId")
// 	invoiceId, err := strconv.ParseInt(id, 10, 64)
// 	if err != nil {
// 		return err
// 	}
// 	row, err := new(model.Payment).CheckRecord(invoiceId, user.Id)
// 	if err != nil || row.Id == 0 {
// 		return fmt.Errorf(GetLang(ctx).Errors.AdTagError.NotFound.ToString())
// 	}
// 	assigns.Row = row
// 	assigns.Billing = new(model.UserBilling).GetByUserId(row.UserId)
// 	assigns.User = new(model.User).GetById(row.UserId)
// 	assigns.Requests, _ = new(model.PaymentRequest).GetByInvoice(row)
// 	assigns.Today = time.Now().Format("2006-01-02")
// 	assigns.PaidDate = row.PaidDate.Format("2006-01-02")
// 	TotalAmountReqeust := float64(0.0)
// 	for _, value := range assigns.Requests {
// 		TotalAmountReqeust = TotalAmountReqeust + value.Amount
// 	}
// 	assigns.TotalAmountReqeust = "$" + humanize.Commaf(TotalAmountReqeust)
//
// 	return ctx.Render("payment/preview", assigns, view.LAYOUTEmpty)
// }

// func (t *Payment) ExportPDF(ctx *fiber.Ctx) error {
//
// 	return ctx.JSON("")
// 	// return ctx.Render("payment/preview", assigns, view.LAYOUTEmpty)
// }
