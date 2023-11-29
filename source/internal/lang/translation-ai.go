package lang

type Ai struct {
	Title            TYPETranslation `json:"title"`
	AddDynamicFloor  TYPETranslation `json:"add_floor"`
	EditDynamicFloor TYPETranslation `json:"edit_floor"`
	Main             TYPETranslation `json:"main"`
	Run              TYPETranslation `json:"run"`
	Status           TYPETranslation `json:"status"`
	Clear            TYPETranslation `json:"clear"`
	Button           TYPETranslation `json:"button"`
}

type AiError struct {
	Add      TYPETranslation `json:"add"`
	Edit     TYPETranslation `json:"edit"`
	Delete   TYPETranslation `json:"delete"`
	List     TYPETranslation `json:"list"`
	NotFound TYPETranslation `json:"not_found"`
}
