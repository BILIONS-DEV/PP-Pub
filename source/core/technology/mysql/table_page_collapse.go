package mysql

type TablePageCollapse struct {
	Id           int64  `gorm:"column:id" json:"id"`
	PageCollapse string `gorm:"column:page_collapse" json:"page_collapse"`
	BoxCollapse  string `gorm:"column:box_collapse" json:"box_collapse"`
	UserId       int64  `gorm:"column:user_id" json:"user_id"`
	IsCollapse   int    `gorm:"column:is_collapse" json:"is_collapse"`
	PageType     string `gorm:"column:page_type" json:"page_type"`
	PageId       int64  `gorm:"column:page_id" json:"page_id"`
}
