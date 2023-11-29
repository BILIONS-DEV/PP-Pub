package lang

type Payment struct {
	Title         TYPETranslation `json:"title"`
	TitleRequest  TYPETranslation `json:"title_request"`
	TitleInvoice  TYPETranslation `json:"title_invoice"`
	Add           TYPETranslation `json:"add"`
	Edit          TYPETranslation `json:"edit"`
	UpdateInvoice TYPETranslation `json:"update_invoice"`
	Run           TYPETranslation `json:"run"`
	Top           TYPETranslation `json:"top"`
	Main          TYPETranslation `json:"main"`
	Name          TYPETranslation `json:"name"`
	StartDate     TYPETranslation `json:"start_date"`
	EndDate       TYPETranslation `json:"end_date"`
	Status        TYPETranslation `json:"status"`
	Note          TYPETranslation `json:"note"`
	User          TYPETranslation `json:"user"`
	Button        TYPETranslation `json:"button"`
}

type PaymentError struct {
	Add      TYPETranslation `json:"add"`
	Edit     TYPETranslation `json:"edit"`
	Delete   TYPETranslation `json:"delete"`
	List     TYPETranslation `json:"list"`
	NotFound TYPETranslation `json:"not_found"`
}
