package mysql

type TableLanguage struct {
	Id   int64  `json:"id"`
	Code string `json:"code"`
	LanguageName string `json:"language_name"`
}

func (TableLanguage) TableName() string {
	return Tables.Language
}