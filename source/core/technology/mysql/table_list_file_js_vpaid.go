package mysql

type TableListFileJsVpaid struct {
	Id   int64  `json:"id"`
	Path string `json:"path"`
	Name string `json:"name"`
}

func (TableListFileJsVpaid) TableName() string {
	return Tables.ListFileJsVpaid
}
