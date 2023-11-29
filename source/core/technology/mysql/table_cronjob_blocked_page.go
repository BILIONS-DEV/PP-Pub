package mysql

func (TableCronJobBlockedPage) TableName() string {
	return "cronjob_blocked_page"
}

type TableCronJobBlockedPage struct {
	MD5    string                       `gorm:"column:md5;primaryKey" json:"md5"`
	Type   TYPECronjobBlockedPage       `gorm:"column:type" json:"type"`
	UserID int64                        `gorm:"column:user_id" json:"user_id"`
	Page   string                       `gorm:"column:page" json:"page"`
	Status TYPEStatusCronjobBlockedPage `gorm:"column:status" json:"status"`
}

type TYPECronjobBlockedPage string

const (
	TYPECronjobBlockedPageBlock   TYPECronjobBlockedPage = "block"
	TYPECronjobBlockedPageUnBlock TYPECronjobBlockedPage = "unblock"
)

type TYPEStatusCronjobBlockedPage string

const (
	TYPEStatusCronjobBlockedPagePending    TYPEStatusCronjobBlockedPage = "pending"
	TYPEStatusCronjobBlockedPageProcessing TYPEStatusCronjobBlockedPage = "processing"
	TYPEStatusCronjobBlockedPageSuccess    TYPEStatusCronjobBlockedPage = "success"
	TYPEStatusCronjobBlockedPageError      TYPEStatusCronjobBlockedPage = "error"
)
