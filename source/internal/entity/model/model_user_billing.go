package model

import (
	"time"
)

func (UserBilling) TableName() string {
	return "user_billing"
}

type UserBilling struct {
	Id                int64     `gorm:"column:id" json:"id"`
	UserId            int64     `gorm:"column:user_id" json:"user_id"`
	Method            string    `gorm:"column:method" json:"method"`
	PaypalEmail       string    `gorm:"column:paypal_email" json:"paypal_email"`
	PayOneerEmail     string    `gorm:"column:payoneer_email" json:"payoneer_email"`
	Cryptocurrency    string    `gorm:"column:cryptocurrency" json:"cryptocurrency"`
	WalletId          string    `gorm:"column:wallet_id" json:"wallet_id"`
	BeneficiaryName   string    `gorm:"column:beneficiary_name" json:"beneficiary_name"`
	BankName          string    `gorm:"column:bank_name" json:"bank_name"`
	BankAddress       string    `gorm:"column:bank_address" json:"bank_address"`
	BankAccountNumber string    `gorm:"column:bank_account_number" json:"bank_account_number"`
	BankRoutingNumber string    `gorm:"column:bank_routing_number" json:"bank_routing_number"`
	BankIbanNumber    string    `gorm:"column:bank_iban_number" json:"bank_iban_number"`
	SwiftCode         string    `gorm:"column:swift_code" json:"swift_code"`
	IFSCCode          string    `gorm:"column:ifsc_code" json:"ifsc_code"`
	NewUpdate         string    `gorm:"column:new_update" json:"new_update"`
	UpdatedAt         time.Time `gorm:"column:updated_at" json:"updated_at"`
	CreatedAt         time.Time `gorm:"column:created_at" json:"created_at"`
}
