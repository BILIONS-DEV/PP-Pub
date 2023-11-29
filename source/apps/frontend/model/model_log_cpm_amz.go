package model

import "source/core/technology/mysql"

type LogCpmAmz struct{}

type LogCpmAmzRecord struct {
	mysql.TableLogCpmAmz
}

func (LogCpmAmzRecord) TableName() string {
	return mysql.Tables.LogCpmAmz
}

func (t *LogCpmAmz) DeleteByBidderId(bidderId int64) {
	mysql.Client.Model(&LogCpmAmzRecord{}).Where("bidder_id = ?", bidderId).Delete(&LogCpmAmzRecord{})
}
