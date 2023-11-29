package payload

type AdBlockIndex struct {
	InventoryAdTagFilterPostData
}

type AdBlockFilterPayload struct {
	StartDate   string `query:"startDate" json:"startDate" form:"startDate"`
	EndDate     string `query:"endDate" json:"endDate" form:"endDate"`
	InventoryId int64  `query:"inventory_id" json:"inventory_id" form:"inventory_id"`
}

type RenderGeneratorPost struct {
	DdbDomain          string `query:"adb_domain" json:"adb_domain" form:"adb_domain"`
	AdbContent         string `query:"adb_content" json:"adb_content" form:"adb_content"`
	AdbModalFixed      string `query:"adb_modal_fixed" json:"adb_modal_fixed" form:"adb_modal_fixed"`
	AdbHideCloseButton string `query:"adb_hide_close_button" json:"adb_hide_close_button" form:"adb_hide_close_button"`
	AdbDisplayTime     string `query:"adb_display_time" json:"adb_display_time" form:"adb_display_time"`
	AdbCloseBackground string `query:"adb_close_background" json:"adb_close_background" form:"adb_close_background"`
}
