package mysql

type TableMediaType struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (TableMediaType) TableName() string {
	return Tables.MediaType
}
