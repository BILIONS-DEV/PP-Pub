package lang

type GAM struct {
	Step struct {
		Headline TYPETranslation `json:"headline"`
		Step1    GamStep         `json:"step_1"`
		Step2    GamStep         `json:"step_2"`
		Step3    GamStep         `json:"step_3"`
		Step4    GamStep         `json:"step_4"`
	} `json:"step"`
}

type GamStep struct {
	Headline    TYPETranslation `json:"headline"`
	Title       TYPETranslation `json:"title"`
	Description TYPETranslation `json:"description"`
}

type GamError struct {
	Add           TYPETranslation `json:"add"`
	Edit          TYPETranslation `json:"edit"`
	Delete        TYPETranslation `json:"delete"`
	List          TYPETranslation `json:"list"`
	NotFound      TYPETranslation `json:"not_found"`
	SelectNetwork TYPETranslation `json:"select_network"`
	PushLineItem  TYPETranslation `json:"push_line_item"`
}
