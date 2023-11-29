package mysql

type TableStickyBanner struct {
	Id                 int64  `gorm:"column:id" json:"id"`
	Name               string `gorm:"column:name" json:"name"`
	UserId             int64  `gorm:"column:user_id" json:"user_id"`
	Type               string `gorm:"type" json:"type"`
	InventoryId        int64  `gorm:"inventory_id" json:"inventory_id"`
	GamDesktop         string `gorm:"column:gam_desktop" json:"gam_desktop"`
	SizeDesktop        int64  `gorm:"size_desktop" json:"size_desktop"`
	PositionDesktop    int    `gorm:"position_desktop" json:"position_desktop"`
	CloseButtonDesktop int    `gorm:"close_button_desktop" json:"close_button_desktop"`
	ShowOnMobile       int    `gorm:"show_on_mobile" json:"show_on_mobile"`
	GamMobile          string `gorm:"gam_mobile" json:"gam_mobile"`
	SizeMobile         int    `gorm:"size_mobile" json:"size_mobile"`
	PositionMobile     int    `gorm:"position_mobile" json:"position_mobile"`
	CloseButtonMobile  int    `gorm:"close_button_mobile" json:"close_button_mobile"`
}
