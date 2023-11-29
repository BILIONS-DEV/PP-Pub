package payload

import "source/pkg/datatable"

type Login struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Remember int    `json:"remember" form:"remember"`
}

type Register struct {
	Email           string `json:"email" form:"email"`
	Password        string `json:"password" form:"password"`
	ConfirmPassWord string `form:"confirm_password" json:"confirm_password"`
	FirstName       string `json:"first_name" form:"first_name"`
	LastName        string `json:"last_name" form:"last_name"`
	Agree           int    `json:"agree" form:"agree"`
}

type Delete struct {
	Id int64
}

type UserIndex struct {
	UserFilterPostData
	QuerySearch string   `query:"f_q" json:"f_q" form:"f_q"`
	Status      []string `query:"f_status" form:"f_status" json:"f_status"`
	Permission  []string `query:"f_permission" form:"f_permission" json:"f_permission"`
}
type UserFilterPayload struct {
	datatable.Request
	PostData *UserFilterPostData `query:"postData"`
}
type UserFilterPostData struct {
	QuerySearch string      `query:"f_q" json:"f_q" form:"f_q"`
	Status      interface{} `query:"f_status[]" json:"f_status" form:"f_status[]"`
	Permission  interface{} `query:"f_permission[]" json:"f_permission" form:"f_permission[]"`
	// datatable input
	Page   int `query:"page" json:"page" form:"page"`
	Limit  int `query:"limit" json:"limit" form:"limit"`
	Start  int `query:"start" json:"start" form:"start"`
	Length int `query:"length" json:"length" form:"length"`
}

type UpdateBilling struct {
	Id                int64  `form:"column:id" json:"id"`
	UserId            int64  `form:"column:user_id" json:"user_id"`
	Method            string `form:"column:method" json:"method"`
	PaypalEmail       string `form:"column:paypal_email" json:"paypal_email"`
	PayOneerEmail     string `form:"column:payoneer_email" json:"payoneer_email"`
	Cryptocurrency    string `form:"column:cryptocurrency" json:"cryptocurrency"`
	WalletId          string `form:"column:wallet_id" json:"wallet_id"`
	BeneficiaryName   string `form:"column:beneficiary_name" json:"beneficiary_name"`
	BankName          string `form:"column:bank_name" json:"bank_name"`
	BankAddress       string `form:"column:bank_address" json:"bank_address"`
	BankAccountNumber string `form:"column:bank_account_number" json:"bank_account_number"`
	BankRoutingNumber string `form:"column:bank_routing_number" json:"bank_routing_number"`
	BankIbanNumber    string `form:"column:bank_iban_number" json:"bank_iban_number"`
	SwiftCode         string `form:"column:swift_code" json:"swift_code"`
	IFSCCode          string `form:"column:ifsc_code" json:"ifsc_code"`
}

type UpdateAccount struct {
	Email string `form:"column:email" json:"email"`
	// PassWord        string `form:"column:password" json:"password"`
	// ConfirmPassWord string `form:"column:confirm_password" json:"confirm_password"`
	LastName    string `form:"column:last_name" json:"last_name"`
	FirstName   string `form:"column:first_name" json:"first_name"`
	Address     string `form:"column:address" json:"address"`
	City        string `form:"column:city" json:"city"`
	State       string `form:"column:state" json:"state"`
	ZipCode     string `form:"column:zip_code" json:"zip_code"`
	Country     string `form:"column:country" json:"country"`
	PhoneNumber string `form:"column:phone_number" json:"phone_number"`
}

type NewPassWord struct {
	OldPassWord     string `form:"column:old_password" json:"old_password"`
	NewPassWord     string `form:"column:new_password" json:"new_password"`
	ConfirmPassWord string `form:"column:confirm_password" json:"confirm_password"`
	Uuid            string `form:"column:uuid" json:"uuid"`
	Email           string `form:"column:email" json:"email"`
}
