package model

import (
	"github.com/badoux/checkmail"
	"source/apps/frontend/payload"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"source/pkg/utility"
	"time"
)

type UserBilling struct{}

type UserBillingRecord struct {
	mysql.TableUserBilling
}

func (UserBillingRecord) TableName() string {
	return mysql.Tables.UserBilling
}

func (t *UserBilling) makeRowBilling(row *payload.UpdateBilling) (record UserBillingRecord) {
	record.Id = row.Id
	record.UserId = row.UserId
	switch row.Method {
	case "bank":
		record.Method = "bank"
		record.BeneficiaryName = row.BeneficiaryName
		record.BankName = row.BankName
		record.BankAddress = row.BankAddress
		record.BankAccountNumber = row.BankAccountNumber
		record.BankRoutingNumber = row.BankRoutingNumber
		record.BankIbanNumber = row.BankIbanNumber
		record.SwiftCode = row.SwiftCode
		record.IFSCCode = row.IFSCCode
		record.PaypalEmail = ""
		record.PayOneerEmail = ""
		record.Cryptocurrency = ""
		record.WalletId = ""
		break
	case "paypal":
		record.Method = "paypal"
		record.BeneficiaryName = ""
		record.BankName = ""
		record.BankAddress = ""
		record.BankAccountNumber = ""
		record.BankRoutingNumber = ""
		record.BankIbanNumber = ""
		record.SwiftCode = ""
		record.IFSCCode = ""
		record.PaypalEmail = row.PaypalEmail
		record.Cryptocurrency = ""
		record.WalletId = ""
		break
	case "payoneer":
		record.Method = "payoneer"
		record.BeneficiaryName = ""
		record.BankName = ""
		record.BankAddress = ""
		record.BankAccountNumber = ""
		record.BankRoutingNumber = ""
		record.BankIbanNumber = ""
		record.SwiftCode = ""
		record.PaypalEmail = ""
		record.IFSCCode = ""
		record.PayOneerEmail = row.PayOneerEmail
		record.Cryptocurrency = ""
		record.WalletId = ""
		break
	case "currency":
		record.Method = "currency"
		record.BeneficiaryName = ""
		record.BankName = ""
		record.BankAddress = ""
		record.BankAccountNumber = ""
		record.BankRoutingNumber = ""
		record.BankIbanNumber = ""
		record.SwiftCode = ""
		record.PaypalEmail = ""
		record.PayOneerEmail = ""
		record.IFSCCode = ""
		record.Cryptocurrency = row.Cryptocurrency
		record.WalletId = row.WalletId
		break
	}
	record.UpdatedAt = time.Now()
	if record.Id != 0 {
		data := t.GetByUserId(row.UserId)
		record.CreatedAt = data.CreatedAt
	}
	return
}

func (t *UserBilling) validateBilling(row *payload.UpdateBilling) (ajaxErrors []ajax.Error) {
	if utility.ValidateString(row.Method) == "" {
		ajaxErrors = append(ajaxErrors, ajax.Error{
			Id:      "method",
			Message: "Method is required",
		})
		return
	}
	switch row.Method {
	case "bank":
		if utility.ValidateString(row.BeneficiaryName) == "" {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "beneficiary_name",
				Message: "Beneficiary Name is required",
			})
		}
		if utility.ValidateString(row.BankName) == "" {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "bank_name",
				Message: "Bank Name is required",
			})
		}
		if utility.ValidateString(row.BankAddress) == "" {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "bank_address",
				Message: "Bank Address is required",
			})
		}
		if utility.ValidateString(row.BankAccountNumber) == "" {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "bank_account_number",
				Message: "Bank Account Number is required",
			})
		}
		break
	case "paypal":
		if utility.ValidateString(row.PaypalEmail) == "" {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "paypal_email",
				Message: "Paypal Email is required",
			})
		}
		err := checkmail.ValidateFormat(row.PaypalEmail)
		if err != nil {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "paypal_email",
				Message: err.Error(),
			})
		}
		break
	case "payoneer":
		if utility.ValidateString(row.PayOneerEmail) == "" {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "payoneer_email",
				Message: "PayOneer Email is required",
			})
		}
		err := checkmail.ValidateFormat(row.PayOneerEmail)
		if err != nil {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "payoneer_email",
				Message: err.Error(),
			})
		}
		break
	case "currency":
		if utility.ValidateString(row.Cryptocurrency) == "" {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "cryptocurrency",
				Message: "Crypto Currency is required",
			})
		}
		if utility.ValidateString(row.WalletId) == "" {
			ajaxErrors = append(ajaxErrors, ajax.Error{
				Id:      "wallet_id",
				Message: "WalletID is required",
			})
		}
		break
	}
	return
}

func (t *UserBilling) Create(row UserBillingRecord) (err error) {
	err = mysql.Client.Create(&row).Error
	return
}

func (t *UserBilling) Update(row UserBillingRecord) (err error) {
	err = mysql.Client.Save(&row).Error
	return
}

func (t *UserBilling) GetByUserId(userId int64) (record UserBillingRecord) {
	mysql.Client.Where("user_id = ?", userId).Find(&record)
	return
}
