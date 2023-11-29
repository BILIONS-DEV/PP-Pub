package dto

import (
	"github.com/badoux/checkmail"
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
	"source/internal/errors"
	"source/pkg/ajax"
	"source/pkg/utility"
	"strings"
)

type Login struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Remember bool   `json:"remember" form:"remember"`
}

func (t *Login) Validate() (err error) {
	if t.Email == "" {
		err = errors.NewWithID("Email is required", "email")
	}
	if t.Password == "" {
		err = errors.NewWithID("Password is required", "password")
	}
	return
}

type UserFilterParams struct {
	QuerySearch string   `query:"f_q" json:"f_q" form:"f_q"`
	Status      []string `query:"f_status" form:"f_status" json:"f_status"`
	UserFilterPostData
}

type UserFilterPayload struct {
	datatable.Request
	PostData *UserFilterPostData
}

type UserFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Status      interface{} `query:"f_status[]" json:"f_status" form:"f_status[]"`
	// datatable input
	Page   int `query:"page" json:"page" form:"page"`
	Limit  int `query:"limit" json:"limit" form:"limit"`
	Start  int `query:"start" json:"start" form:"start"`
	Length int `query:"length" json:"length" form:"length"`
}

type PayloadUser struct {
	ID              int64  `json:"id" form:"id"`
	Email           string `json:"email" form:"email"`
	Password        string `json:"password" form:"password"`
	ConfirmPassWord string `form:"confirm_password" json:"confirm_password"`
	FirstName       string `json:"first_name" form:"first_name"`
	LastName        string `json:"last_name" form:"last_name"`
	Agree           int    `json:"agree" form:"agree"`
}

func (t *PayloadUser) Validate() (ajaxErrors []ajax.Error) {

	if utility.ValidateString(t.FirstName) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "first_name",
			Message: "first name is required",
		})
	}
	if utility.ValidateString(t.LastName) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "last_name",
			Message: "last name is required",
		})
	}
	if utility.ValidateString(t.Email) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "email",
			Message: "email is required",
		})
	} else {
		err := checkmail.ValidateFormat(t.Email)
		if err != nil {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "email",
				Message: err.Error(),
			})
		}
	}
	if utility.ValidateString(t.Password) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "password",
			Message: "password is required",
		})
	} else {
		if len(t.Password) < 8 {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "password",
				Message: "password must be longer than 8 characters",
			})
		}
	}
	if utility.ValidateString(t.ConfirmPassWord) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "confirm_password",
			Message: "confirm password is required",
		})
	} else {
		if len(t.ConfirmPassWord) < 8 {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "confirm_password",
				Message: "confirm password must be longer than 8 characters",
			})
		}
	}
	if utility.ValidateString(t.Password) != utility.ValidateString(t.ConfirmPassWord) {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "",
			Message: "Password and Confirm Password don't match",
		})
	}
	if t.Agree != 1 {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "agree",
			Message: "please read carefully and accept our terms of use",
		})
	}
	return
}

func (t *PayloadUser) ToModel() (record model.User) {
	record.ID = t.ID
	record.Email = strings.TrimSpace(t.Email)
	record.Password = record.MakePassword(strings.TrimSpace(t.Password))
	record.FirstName = strings.TrimSpace(t.FirstName)
	record.LastName = strings.TrimSpace(t.LastName)
	record.LoginToken = record.MakeLoginToken()
	record.Permission = model.UserPermissionMember
	record.Status = model.StatusApproved
	record.PaymentNet = 30
	record.PaymentTerm = 1
	return
}

type PayloadProfile struct {
	ID          int64  `json:"id" form:"id"`
	Email       string `json:"email" form:"email"`
	FirstName   string `json:"first_name" form:"first_name"`
	LastName    string `json:"last_name" form:"last_name"`
	Address     string `json:"address" form:"address"`
	City        string `json:"city" form:"city"`
	State       string `json:"state" form:"state"`
	ZipCode     string `json:"zip_code" form:"zip_code"`
	Country     string `json:"country" form:"country"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
}

type PayloadBilling struct {
	ID                int64  `form:"column:id" json:"id"`
	UserID            int64  `form:"column:user_id" json:"user_id"`
	Method            string `form:"column:method" json:"method"`
	PaypalEmail       string `form:"column:paypal_email" json:"paypal_email"`
	PayOneerEmail     string `form:"column:payoneer_email" json:"payoneer_email"`
	Cryptocurrency    string `form:"column:cryptocurrency" json:"cryptocurrency"`
	WalletID          string `form:"column:wallet_id" json:"wallet_id"`
	BeneficiaryName   string `form:"column:beneficiary_name" json:"beneficiary_name"`
	BankName          string `form:"column:bank_name" json:"bank_name"`
	BankAddress       string `form:"column:bank_address" json:"bank_address"`
	BankAccountNumber string `form:"column:bank_account_number" json:"bank_account_number"`
	BankRoutingNumber string `form:"column:bank_routing_number" json:"bank_routing_number"`
	BankIbanNumber    string `form:"column:bank_iban_number" json:"bank_iban_number"`
	SwiftCode         string `form:"column:swift_code" json:"swift_code"`
	IFSCCode          string `form:"column:ifsc_code" json:"ifsc_code"`
}

func (t *PayloadBilling) Validate() (ajaxErrors []ajax.Error) {
	if t.UserID == 0 {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "user_id",
			Message: "UserID is required",
		})
	}
	if utility.ValidateString(t.Method) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "first_name",
			Message: "Method is required",
		})
	}

	switch t.Method {
	case "paypal":
		if utility.ValidateString(t.PaypalEmail) == "" {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "first_name",
				Message: "Paypal email is required",
			})
		}
	case "payoneer":
		if utility.ValidateString(t.PayOneerEmail) == "" {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "first_name",
				Message: "PayOneer email is required",
			})
		}
	case "bank":
		if utility.ValidateString(t.BeneficiaryName) == "" {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "beneficiary_name",
				Message: "Beneficiary name is required",
			})
		}
		if utility.ValidateString(t.BankName) == "" {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "bank_name",
				Message: "Bank name is required",
			})
		}
		if utility.ValidateString(t.BankAddress) == "" {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "bank_address",
				Message: "Bank address is required",
			})
		}
		if utility.ValidateString(t.BankAccountNumber) == "" {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "bank_account_number",
				Message: "Bank account number is required",
			})
		}
	}
	return
}

func (t *PayloadBilling) ToModel() (billing model.UserBilling) {
	billing.Id = t.ID
	billing.UserId = t.UserID
	billing.Method = t.Method
	switch billing.Method {
	case "paypal":
		billing.PaypalEmail = t.PaypalEmail
	case "payoneer":
		billing.PayOneerEmail = t.PayOneerEmail
	case "currency":
		billing.Cryptocurrency = t.Cryptocurrency
		billing.WalletId = t.WalletID
	case "bank":
		billing.BeneficiaryName = t.BeneficiaryName
		billing.BankName = t.BankName
		billing.BankAddress = t.BankAddress
		billing.BankAccountNumber = t.BankAccountNumber
		billing.BankRoutingNumber = t.BankRoutingNumber
		billing.BankIbanNumber = t.BankIbanNumber
		billing.SwiftCode = t.SwiftCode
		billing.IFSCCode = t.IFSCCode
	}
	return
}
