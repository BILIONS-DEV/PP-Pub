package mysql

type TableListFileJs struct {
	Id   int64  `json:"id"`
	Path string `json:"path"`
	Name string `json:"name"`
}

func (TableListFileJs) TableName() string {
	return Tables.ListFileJs
}
