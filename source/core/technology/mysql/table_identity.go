package mysql

import (
	"gorm.io/gorm"
	"time"
)

type TableIdentity struct {
	Id            int64                     `gorm:"column:id" json:"id"`
	UserId        int64                     `gorm:"column:user_id" json:"user_id"`
	Name          string                    `gorm:"column:name" json:"name"`
	Description   string                    `gorm:"column:description" json:"description"`
	AuctionDelay  int                       `gorm:"column:auction_delay" json:"auction_delay"`
	SyncDelay     int                       `gorm:"column:sync_delay" json:"sync_delay"`
	Status        TypeOnOff                 `gorm:"column:status" json:"status"`
	Priority      int                       `gorm:"column:priority" json:"priority"`
	IsDefault     TypeOnOff                 `gorm:"column:is_default" json:"is_default"`
	CreatedAt     time.Time                 `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time                 `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     gorm.DeletedAt            `gorm:"column:deleted_at" json:"deleted_at"`
	UserIdModules []TableIdentityModuleInfo `gorm:"-"`
	Targets       []TableTarget             `gorm:"-"`
}

func (TableIdentity) TableName() string {
	return Tables.Identity
}

func (rec *TableIdentity) GetRls() {
	var userIdModules []TableIdentityModuleInfo
	Client.Where("identity_id = ?", rec.Id).Order("module_id").Find(&userIdModules)
	rec.UserIdModules = userIdModules
	var targets []TableTarget
	Client.Where(TableTarget{IdentityId: rec.Id}).Find(&targets)
	rec.Targets = targets
	return
}
