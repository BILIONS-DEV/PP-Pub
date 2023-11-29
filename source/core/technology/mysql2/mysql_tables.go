package mysql2

var Tables TableList

func init() {
	Tables = TableList{
		Inventory:           "inventory",
		InventoryAdTag:      "inventory_ad_tag",
		Blocking:            "blocking",
		BlockingRestriction: "blocking_restriction_2",
		LineItem:            "line_item",
		Target:              "target",
	}

}

type TableList struct {
	Inventory           string
	InventoryAdTag      string
	Blocking            string
	BlockingRestriction string
	LineItem            string
	Target              string
}
