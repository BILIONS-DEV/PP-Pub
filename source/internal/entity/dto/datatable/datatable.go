package datatable

import (
	"fmt"
	"strings"
)

type Request struct {
	Draw    int      `query:"draw" json:"draw" form:"draw"`
	Start   int      `query:"start" json:"start" form:"start"`
	Length  int      `query:"length" json:"length" form:"length"`
	Search  *Search  `query:"search" json:"search" form:"search"`
	Columns []Column `query:"columns" json:"columns" form:"columns"`
	Order   []Order  `query:"order" json:"order" form:"order"`
}

func (t *Request) OrderString() (orderString string) {
	if len(t.Order) > 0 {
		var orders []string
		for _, order := range t.Order {
			column := t.Columns[order.Column]
			orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
		}
		orderString = strings.Join(orders, ", ")
	}
	return
}

type Search struct {
	Value string `query:"value" json:"value" form:"value"`
	Regex bool   `query:"regex" json:"regex" form:"regex"`
}
type Order struct {
	Column int    `query:"column" json:"column" form:"column"`
	Dir    string `query:"dir" json:"dir" form:"dir"`
}
type Column struct {
	Data       string `query:"data" json:"data" form:"data" schema:"data"`
	Name       string `query:"name" json:"name" form:"name" schema:"name"`
	Searchable bool   `query:"searchable" json:"searchable" form:"searchable"`
	OrderAble  bool   `query:"orderable" json:"order_able" form:"order_able"`
	Search     Search `query:"search" json:"search" form:"search"`
}

type Response struct {
	Draw            int         `json:"draw"`
	RecordsTotal    int64       `json:"recordsTotal"`
	RecordsFiltered int64       `json:"recordsFiltered"`
	Data            interface{} `json:"data"`
}

type PaymentTerm struct {
	Name  string
	Value string
}

type PaymentNet struct {
	Name  string
	Value int
}

type Status struct {
	Name  string
	Value int
}

type Permission struct {
	Name  string
	Value string
}

type Type struct {
	Name  string
	Value int
}
