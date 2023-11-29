package lang

type User struct {
	TopBilling        TYPETranslation `json:"top_billing"`
	MainBilling       TYPETranslation `json:"main_billing"`
	Method            TYPETranslation `json:"method"`
	PaypalEmail       TYPETranslation `json:"paypal_email"`
	PayOneerEmail     TYPETranslation `json:"payoneer_email"`
	BeneficiaryName   TYPETranslation `json:"beneficiary_name"`
	BankName          TYPETranslation `json:"bank_name"`
	BankAddress       TYPETranslation `json:"bank_address"`
	BankAccountNumber TYPETranslation `json:"bank_account_number"`
	BankRoutingNumber TYPETranslation `json:"bank_routing_number"`
	BankIbanNumber    TYPETranslation `json:"bank_iban_number"`
	SwiftCode         TYPETranslation `json:"swift_code"`
	TopAccount        TYPETranslation `json:"top_account"`
	MainAccount       TYPETranslation `json:"main_account"`
	Profile           TYPETranslation `json:"profile"`
	ChangePassword    TYPETranslation `json:"change_password"`
	Email             TYPETranslation `json:"email"`
	FirstName         TYPETranslation `json:"first_name"`
	LastName          TYPETranslation `json:"last_name"`
	OldPassword       TYPETranslation `json:"old_password"`
	NewPassword       TYPETranslation `json:"new_password"`
	ConfirmPassword   TYPETranslation `json:"confirm_password"`
	Button            TYPETranslation `json:"button"`
}

type UserError struct {
	Login            TYPETranslation `json:"login"`
	PermissionDenied TYPETranslation `json:"permission_denied"`
	Email            TYPETranslation `json:"email"`
	CreateBilling    TYPETranslation `json:"create_billing"`
	UpdateBilling    TYPETranslation `json:"update_billing"`
	UpdateAccount    TYPETranslation `json:"update_account"`
	ChangePassword   TYPETranslation `json:"change_password"`
	GetByEmail       TYPETranslation `json:"get_by_email"`
	LinkValid        TYPETranslation `json:"link_valid"`
	LinkOutDate      TYPETranslation `json:"link_out_date"`
	SendMail         TYPETranslation `json:"send_mail"`
	List             TYPETranslation `json:"list"`
}
