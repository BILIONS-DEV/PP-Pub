package mysql

type TableDynamicFloor struct {
	Id int64 `gorm:"column:id;primary_key" json:"id"`
	//InventoryId     int64         `gorm:"column:inventory_id" json:"inventory_id"`
	//Size            string        `gorm:"column:size" json:"size"`
	//Format          string        `gorm:"column:format" json:"format"`
	//TagId           pq.Int64Array `gorm:"column:tagId;type:integer[]" json:"tagId"`
	Device          string  `gorm:"column:device" json:"device"`
	CountryCode     string  `gorm:"column:countryCode" json:"countryCode"`
	FloorTest       string  `gorm:"column:floorTest" json:"floorTest"`
	Status          string  `gorm:"column:status" json:"status"`
	StartFloor      float64 `gorm:"column:startFloor" json:"startFloor"`
	Volume          int64   `gorm:"column:volume" json:"volume"`
	DesiredFillRate float64 `gorm:"column:desiredFillRate" json:"desiredFillRate"`
	CurrentFloor    float64 `gorm:"column:currentFloor" json:"currentFloor"`
	CurrentStep     float64 `gorm:"column:currentStep" json:"currentStep"`
	LastScanTime    int64   `gorm:"column:lastScanTime" json:"lastScanTime"`
	NumberTest      int64   `gorm:"column:numberTest" json:"numberTest"`
	FoundFloor      int64   `gorm:"column:foundFloor" json:"foundFloor"`
	FloorMatched    float64 `gorm:"column:floorMatched" json:"floorMatched"`
}

func (TableDynamicFloor) TableName() string {
	return Tables.DynamicFloor
}

type TableDynamicFloorTag struct {
	Id             int64 `gorm:"column:id" json:"id"`
	DynamicFloorId int64 `gorm:"column:dynamic_floor_id" json:"dynamic_floor_id"`
	TagId          int64 `gorm:"column:tag_id" json:"tag_id"`
}

func (TableDynamicFloorTag) TableName() string {
	return Tables.DynamicFloorTag
}
