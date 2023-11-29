package mysql

type TableCategory struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (TableCategory) TableName() string {
	return Tables.Category
}