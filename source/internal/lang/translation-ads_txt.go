package lang

type AdsTxt struct {
	Title        TYPETranslation `json:"title"`
	SearchDomain TYPETranslation `json:"search_domain"`
	Run          TYPETranslation `json:"run"`
	TitleDetail  TYPETranslation `json:"title_detail"`
	Setup1       TYPETranslation `json:"setup_1"`
	Content1     TYPETranslation `json:"content_1"`
	Setup2       TYPETranslation `json:"setup_2"`
	Content2     TYPETranslation `json:"content_2"`
	Save         TYPETranslation `json:"save"`
	Setup3       TYPETranslation `json:"setup_3"`
	Content3     TYPETranslation `json:"content_3"`
	RedirectUrl  TYPETranslation `json:"redirect_url"`
	Different    TYPETranslation `json:"different"`
	Setup4       TYPETranslation `json:"setup_4"`
	Content4     TYPETranslation `json:"content_4"`
	ButtonCheck  TYPETranslation `json:"button_check"`
}

type AdsTxtError struct {
	NotFound TYPETranslation `json:"not_found"`
}