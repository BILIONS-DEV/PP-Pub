package mysql

type TableSystemSetting struct {
	Id            int    `gorm:"column:id"`
	JsProcessPath string `gorm:"column:js_process_path"`
}

func (TableSystemSetting) GetJsProcessPath() string {
	rec := TableSystemSetting{}
	Client.Table(Tables.SystemSetting).Select("js_process_path").Where("id = ?", 1).Last(&rec)
	return rec.JsProcessPath
}

func (TableSystemSetting) GetJsVpaidPath() string {
	rec := TableSystemSetting{}
	Client.Table(Tables.SystemSetting).Select("js_process_path").Where("id = ?", 3).Last(&rec)
	return rec.JsProcessPath
}