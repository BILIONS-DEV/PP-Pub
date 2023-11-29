package model

import "source/core/technology/mysql"

type ChannelsKeyword struct{}

type ChannelsKeywordRecord struct {
	mysql.TableChannelsKeyword
}

func (ChannelsKeywordRecord) TableName() string {
	return mysql.Tables.ChannelsKeyword
}

func (t *ChannelsKeyword) GetByChannels(channelsId int64) (records []ChannelsKeywordRecord) {
	mysql.Client.Where("channels_id = ?", channelsId).Find(&records)
	return
}

func (t *ChannelsKeyword) DeleteTagByChannels(channelsId int64) {
	mysql.Client.Where(ChannelsKeywordRecord{mysql.TableChannelsKeyword{ChannelsId: channelsId}}).Delete(ChannelsKeywordRecord{})
	return
}