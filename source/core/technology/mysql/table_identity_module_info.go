package mysql

type TableIdentityModuleInfo struct {
	Id          int64          `gorm:"column:id" json:"id"`
	IdentityId int64          `gorm:"column:identity_id" json:"identity_id"`
	ModuleId    int64          `gorm:"column:module_id" json:"module_id"`
	Name        string         `gorm:"column:name" json:"name"`
	Params      string         `gorm:"column:params" json:"params"`
	Storage     string         `gorm:"column:storage" json:"storage"`
	AbTesting   TYPEOnOff      `gorm:"column:ab_testing" json:"ab_testing"`
	Volume      int            `gorm:"column:volume" json:"volume"`
}

func (TableIdentityModuleInfo) TableName() string {
	return Tables.IdentityModuleInfo
}
