package history

import (
	"encoding/json"
	"fmt"
	"source/core/technology/mysql"
)

type User struct {
	Detail    DetailUser
	CreatorId int64
	RecordOld mysql.TableUser
	RecordNew mysql.TableUser
}

func (t *User) Page() string {
	return "User"
}

type DetailUser int

const (
	DetailUserProfileFE DetailUser = iota + 1
	DetailUserChangePassFE
	DetailUserBillingFE
	DetailUserBE
)

func (t DetailUser) String() string {
	switch t {
	case DetailUserProfileFE:
		return "user_profile_fe"
	case DetailUserChangePassFE:
		return "user_change_password_fe"
	case DetailUserBillingFE:
		return "user_billing_fe"
	case DetailUserBE:
		return "user_be"
	}
	return ""
}

func (t DetailUser) App() string {
	switch t {
	case DetailUserProfileFE, DetailUserChangePassFE, DetailUserBillingFE:
		return "FE"
	case DetailUserBE:
		return "BE"
	}
	return ""
}

func (t *User) Type() TYPEHistory {
	return TYPEHistoryUser
}

func (t *User) Action() mysql.TYPEObjectType {
	if t.RecordOld.Id == 0 && t.RecordNew.Id != 0 {
		return mysql.TYPEObjectTypeAdd
	} else if t.RecordOld.Id != 0 && t.RecordNew.Id == 0 {
		return mysql.TYPEObjectTypeDel
	}
	return mysql.TYPEObjectTypeUpdate
}

func (t *User) Data() mysql.TableHistory {
	switch t.Detail {
	case DetailUserProfileFE:
		return t.getHistoryUserFEProfile()
	case DetailUserChangePassFE:
		return t.getHistoryUserChangePasswordFE()
	case DetailUserBillingFE:
		return t.getHistoryUserBillingFE()
	case DetailUserBE:
		return t.getHistoryUserBE()
	}
	return mysql.TableHistory{}
}

func (t *User) CompareData(history mysql.TableHistory) (res []ResponseCompare) {
	switch history.DetailType {
	case DetailUserProfileFE.String():
		return t.compareDataUserProfileFE(history)
	case DetailUserChangePassFE.String():
		return t.compareDataUserChangePassFE(history)
	case DetailUserBillingFE.String():
		return t.compareDataUserBillingFE(history)
	case DetailUserBE.String():
		return t.compareDataUserBE(history)
	}
	return []ResponseCompare{}
}

func (t *User) getRootRecord() (record mysql.TableUser) {
	switch t.Action() {
	case mysql.TYPEObjectTypeAdd:
		return t.RecordNew
	case mysql.TYPEObjectTypeUpdate:
		return t.RecordNew
	case mysql.TYPEObjectTypeDel:
		return t.RecordOld
	}
	return
}

type userFEProfile struct {
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
	ZipCode     string `json:"zip_code"`
	Country     string `json:"country"`
	PhoneNumber string `json:"phone_number"`
}

func (t *User) getHistoryUserFEProfile() (history mysql.TableHistory) {
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.User,
		ObjectId:   t.getRootRecord().Id,
		ObjectName: t.getRootRecord().Email,
		ObjectType: t.Action(),
		DetailType: t.Detail.String(),
		App:        t.Detail.App(),
		UserId:     t.getRootRecord().Id,
	}
	var newData, oldData []byte
	if t.RecordNew.Id != 0 {
		newData, _ = json.Marshal(userFEProfile{
			Email:       t.RecordNew.Email,
			FirstName:   t.RecordNew.FirstName,
			LastName:    t.RecordNew.LastName,
			Address:     t.RecordNew.Address,
			City:        t.RecordNew.City,
			State:       t.RecordNew.State,
			ZipCode:     t.RecordNew.ZipCode,
			Country:     t.RecordNew.Country,
			PhoneNumber: t.RecordNew.PhoneNumber,
		})
	}
	if t.RecordOld.Id != 0 {
		oldData, _ = json.Marshal(userFEProfile{
			Email:       t.RecordOld.Email,
			FirstName:   t.RecordOld.FirstName,
			LastName:    t.RecordOld.LastName,
			Address:     t.RecordOld.Address,
			City:        t.RecordOld.City,
			State:       t.RecordOld.State,
			ZipCode:     t.RecordOld.ZipCode,
			Country:     t.RecordOld.Country,
			PhoneNumber: t.RecordOld.PhoneNumber,
		})
	}
	if t.Action() == mysql.TYPEObjectTypeAdd {
		history.Title = "Add User"
		history.NewData = string(newData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update User"
		history.NewData = string(newData)
		history.OldData = string(oldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete User"
		history.OldData = string(oldData)
	}
	return
}

func (t *User) getHistoryUserChangePasswordFE() (history mysql.TableHistory) {
	history = mysql.TableHistory{
		CreatorId:  t.RecordNew.Id,
		Title:      "Change Password",
		Object:     mysql.Tables.User,
		ObjectId:   t.RecordNew.Id,
		ObjectName: t.RecordNew.Email,
		ObjectType: t.Action(),
		DetailType: t.Detail.String(),
		App:        t.Detail.App(),
		UserId:     t.RecordNew.Id,
	}
	return
}

func (t *User) compareDataUserProfileFE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew userFEProfile
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row

	// Email
	res, err := makeResponseCompare("Email", &recordOld.Email, &recordNew.Email, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// First Name
	res, err = makeResponseCompare("First Name", &recordOld.FirstName, &recordNew.FirstName, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Last Name
	res, err = makeResponseCompare("Last Name", &recordOld.LastName, &recordNew.LastName, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Address
	res, err = makeResponseCompare("Address", &recordOld.Address, &recordNew.Address, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// City
	res, err = makeResponseCompare("City", &recordOld.City, &recordNew.City, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// State
	res, err = makeResponseCompare("State", &recordOld.State, &recordNew.State, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Zip Code
	res, err = makeResponseCompare("Zip Code", &recordOld.ZipCode, &recordNew.ZipCode, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Country
	res, err = makeResponseCompare("Country", &recordOld.Country, &recordNew.Country, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Phone Number
	res, err = makeResponseCompare("Phone Number", &recordOld.PhoneNumber, &recordNew.PhoneNumber, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}

func (t *User) compareDataUserChangePassFE(history mysql.TableHistory) (responses []ResponseCompare) {
	// Xử lý compare từng row

	// Title
	responses = append(responses, ResponseCompare{
		Action:  "update",
		Text:    history.Title,
		OldData: "",
		NewData: "",
	})
	return
}

// User BE
type userBE struct {
	Email          string `json:"email"`
	Status         string `json:"status"`
	AccountManager string `json:"account_manager"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Permission     string `json:"permission"`
	Password       string `json:"password"`
	RevenueShare   string `json:"revenue_share"`
	PaymentNet     string `json:"payment_net"`
	PaymentTemp    string `json:"payment_temp"`
	AccountType    string `json:"account_type"`
	CustomAdsTxt   string `json:"custom_ads_txt"`
}

func (t *User) getHistoryUserBE() (history mysql.TableHistory) {
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.User,
		ObjectId:   t.getRootRecord().Id,
		ObjectName: t.getRootRecord().Email,
		ObjectType: t.Action(),
		DetailType: t.Detail.String(),
		App:        t.Detail.App(),
		UserId:     t.getRootRecord().Id,
	}
	var newData, oldData []byte
	if t.RecordNew.Id != 0 {
		newData, _ = json.Marshal(t.makeDataUserBE(t.RecordNew))
	}
	if t.RecordOld.Id != 0 {
		oldData, _ = json.Marshal(t.makeDataUserBE(t.RecordOld))
	}
	if t.Action() == mysql.TYPEObjectTypeAdd {
		history.Title = "Add User"
		history.NewData = string(newData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update User"
		history.NewData = string(newData)
		history.OldData = string(oldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete User"
		history.OldData = string(oldData)
	}
	return
}

func (t *User) makeDataUserBE(record mysql.TableUser) (data userBE) {
	data.Email = record.Email
	data.Status = record.Status.String()

	var presenter mysql.TableUser
	mysql.Client.Find(&presenter, record.Presenter)
	data.AccountManager = presenter.Email

	data.FirstName = record.FirstName
	data.LastName = record.LastName
	data.Permission = record.Permission.String()
	data.Password = record.Password
	data.RevenueShare = fmt.Sprintf("%.f", record.Revenue)
	data.PaymentNet = fmt.Sprintf("%d", record.PaymentNet)
	data.PaymentTemp = fmt.Sprintf("%d", record.PaymentTerm)
	data.AccountType = record.AccountType.String()
	data.CustomAdsTxt = string(record.AdsTxtCustomByAdmin)
	return
}

func (t *User) compareDataUserBE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew userBE
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row
	// Email
	res, err := makeResponseCompare("Email", &recordOld.Email, &recordNew.Email, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Status
	res, err = makeResponseCompare("Status", &recordOld.Status, &recordNew.Status, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Account Manage
	res, err = makeResponseCompare("Account Manage", &recordOld.AccountManager, &recordNew.AccountManager, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// First Name
	res, err = makeResponseCompare("First Name", &recordOld.FirstName, &recordNew.FirstName, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Last Name
	res, err = makeResponseCompare("Last Name", &recordOld.LastName, &recordNew.LastName, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Permission
	res, err = makeResponseCompare("Permission", &recordOld.Permission, &recordNew.Permission, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Revenue Share
	res, err = makeResponseCompare("Revenue Share", &recordOld.RevenueShare, &recordNew.RevenueShare, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Payment Net
	res, err = makeResponseCompare("Payment Net", &recordOld.PaymentNet, &recordNew.PaymentNet, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Payment Temp
	res, err = makeResponseCompare("Payment Temp", &recordOld.PaymentTemp, &recordNew.PaymentTemp, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Account Type
	res, err = makeResponseCompare("Account Type", &recordOld.AccountType, &recordNew.AccountType, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Custom Ads Txt
	res, err = makeResponseCompare("Custom Ads Txt", &recordOld.CustomAdsTxt, &recordNew.CustomAdsTxt, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}

type userBillingFE struct {
	Method            *string `json:"method,omitempty"`
	PayOneerEmail     *string `json:"payoneer_email,omitempty"`
	BeneficiaryName   *string `json:"beneficiary_name,omitempty"`
	BankName          *string `json:"bank_name,omitempty"`
	BankAddress       *string `json:"bank_address,omitempty"`
	BankAccountNumber *string `json:"bank_account_number,omitempty"`
	BankRoutingNumber *string `json:"bank_routing_number,omitempty"`
	BankIbanNumber    *string `json:"bank_iban_number,omitempty"`
	SwiftCode         *string `json:"swift_code,omitempty"`
}

func (t *User) getHistoryUserBillingFE() (history mysql.TableHistory) {
	history = mysql.TableHistory{
		CreatorId:  t.CreatorId,
		Object:     mysql.Tables.UserBilling,
		ObjectId:   t.getRootRecord().UserBilling.Id,
		ObjectName: t.getRootRecord().Email,
		ObjectType: t.Action(),
		DetailType: t.Detail.String(),
		App:        t.Detail.App(),
		UserId:     t.getRootRecord().UserBilling.UserId,
	}
	var newData, oldData []byte
	if t.RecordNew.UserBilling.Id != 0 {
		var record userBillingFE
		record.Method = &t.RecordNew.UserBilling.Method
		if t.RecordNew.UserBilling.Method == "payoneer" {
			record.PayOneerEmail = &t.RecordNew.UserBilling.PayOneerEmail
		} else if t.RecordNew.UserBilling.Method == "bank" {
			record.BeneficiaryName = &t.RecordNew.UserBilling.BeneficiaryName
			record.BankName = &t.RecordNew.UserBilling.BankName
			record.BankAddress = &t.RecordNew.UserBilling.BankAddress
			record.BankAccountNumber = &t.RecordNew.UserBilling.BankAccountNumber
			record.BankRoutingNumber = &t.RecordNew.UserBilling.BankRoutingNumber
			record.BankIbanNumber = &t.RecordNew.UserBilling.BankIbanNumber
			record.SwiftCode = &t.RecordNew.UserBilling.SwiftCode
		}
		newData, _ = json.Marshal(record)
	}
	if t.RecordOld.UserBilling.Id != 0 {
		var record userBillingFE
		record.Method = &t.RecordOld.UserBilling.Method
		if t.RecordOld.UserBilling.Method == "payoneer" {
			record.PayOneerEmail = &t.RecordOld.UserBilling.PayOneerEmail
		} else if t.RecordOld.UserBilling.Method == "bank" {
			record.BeneficiaryName = &t.RecordOld.UserBilling.BeneficiaryName
			record.BankName = &t.RecordOld.UserBilling.BankName
			record.BankAddress = &t.RecordOld.UserBilling.BankAddress
			record.BankAccountNumber = &t.RecordOld.UserBilling.BankAccountNumber
			record.BankRoutingNumber = &t.RecordOld.UserBilling.BankRoutingNumber
			record.BankIbanNumber = &t.RecordOld.UserBilling.BankIbanNumber
			record.SwiftCode = &t.RecordOld.UserBilling.SwiftCode
		}
		oldData, _ = json.Marshal(record)
	}
	if t.Action() == mysql.TYPEObjectTypeAdd {
		history.Title = "Add User Billing"
		history.NewData = string(newData)
		history.CreatedAt = t.RecordNew.CreatedAt
	} else if t.Action() == mysql.TYPEObjectTypeUpdate {
		history.Title = "Update User Billing"
		history.NewData = string(newData)
		history.OldData = string(oldData)
	} else if t.Action() == mysql.TYPEObjectTypeDel {
		history.Title = "Delete User Billing"
		history.OldData = string(oldData)
	}
	return
}

func (t *User) compareDataUserBillingFE(history mysql.TableHistory) (responses []ResponseCompare) {
	var recordOld, recordNew userBillingFE
	_ = json.Unmarshal([]byte(history.OldData), &recordOld)
	_ = json.Unmarshal([]byte(history.NewData), &recordNew)

	// Xử lý compare từng row
	// Method
	res, err := makeResponseCompare("Method", recordOld.Method, recordNew.Method, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Payoneer Email
	res, err = makeResponseCompare("Payoneer Email", recordOld.PayOneerEmail, recordNew.PayOneerEmail, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Beneficiary Name
	res, err = makeResponseCompare("Beneficiary Name", recordOld.BeneficiaryName, recordNew.BeneficiaryName, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Bank Name
	res, err = makeResponseCompare("Bank Name", recordOld.BankName, recordNew.BankName, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Bank Address
	res, err = makeResponseCompare("Bank Address", recordOld.BankAddress, recordNew.BankAddress, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Bank Account Number
	res, err = makeResponseCompare("Bank Account Number", recordOld.BankAccountNumber, recordNew.BankAccountNumber, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Bank Routing Number
	res, err = makeResponseCompare("Bank Routing Number", recordOld.BankRoutingNumber, recordNew.BankRoutingNumber, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Bank Iban Number
	res, err = makeResponseCompare("Bank Iban Number", recordOld.BankIbanNumber, recordNew.BankIbanNumber, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	// Swift Code
	res, err = makeResponseCompare("Swift Code", recordOld.SwiftCode, recordNew.SwiftCode, history.ObjectType)
	if err == nil {
		responses = append(responses, res)
	}
	return
}
