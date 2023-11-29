package mysql

type TableLogCpmAmz struct {
	Id        int64   `json:"id"`
	BidderId  int64   `json:"bidder_id"`
	NetworkId int64     `json:"network_id"`
	XlsxPath  string  `json:"xlsx_path"`
	Target    string  `json:"target"`
	Value     float64 `json:"value"`
}

func (TableLogCpmAmz) TableName() string {
	return Tables.ListFileJs
}
