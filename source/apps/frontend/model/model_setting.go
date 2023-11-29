package model

import "source/core/technology/mysql"

type Setting struct{}

type SettingRecord struct {
	mysql.TableSetting
}

func (SettingRecord) TableName() string {
	return mysql.Tables.Setting
}

func (t *Setting) GetAdsTxtApd() (adsTxt mysql.TYPEAdsTxtCustom) {
	var setting SettingRecord
	mysql.Client.Where("meta_key = 'ads_txt'").Find(&setting)
	adsTxt = mysql.TYPEAdsTxtCustom(setting.MetaValue)
	return
}
