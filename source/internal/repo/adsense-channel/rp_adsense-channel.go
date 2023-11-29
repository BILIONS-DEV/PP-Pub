package adsense_channel

import (
	"gorm.io/gorm"
	"source/infrastructure/caching"
	"source/internal/entity/model"
	"time"
)

const KeyName = "RepoAdsenseChannel"

type RepoAdsenseChannel interface {
	AddAdsenseChannel(record *model.AdsenseChannelModel) (err error)
	GetAllAdsenseChannels() (records []model.AdsenseChannelModel)
	GetChannelsByAccount(account string) (channels []model.AdsenseChannelModel, err error)
	// GetChannelNotCampaign(account string) (records model.AdsenseChannelModel, err error)
	// RemoveAdsenseChannel(id int64) (err error)
	RemoveChannelByAccount(channels []string, account string) (err error)
	IsExistChannel(channel string, accountAdsense string) bool
	// UpdateCampaignForChannel(channel, country string, campaignID int64) (err error)
	// UpdateChannelForCampaign(channel string, campaignID int64) (err error)
	// RemoveCampaignForChannel(campaignID int64) (err error)
	GetChannelsUse(account string) (channelsUse, total int64)
	GetChannelsNotUse() (records []model.AdsenseChannelModel, total int64)
	GetChannelByChannel(channel, account string) (record model.AdsenseChannelModel, err error)
	ResetCampaignIDToChannel() (err error)
	CountChannelNotUse(account string) (total int64, err error)
	// CountChannelAvailableByAccount(account string) (total int64)
	CheckCampaignCanUseChannel(campaignID int64, channel string, locations []string) bool
	AddChannelCountry(record model.ChannelCountryModel) (err error)
	EditChannelCountry(record model.ChannelCountryModel) (err error)
	OffChannelCountryByCampaign(campaignID int64) (err error)
	RemoveChannelCountryByCampaign(campaignID int64, channel string) (err error)
	FindChannelCountry(channel, country string) (record model.ChannelCountryModel, err error)
	FindChannelOFF(channel string) (record model.ChannelCountryModel, err error)
}

type AdsenseChannelRP struct {
	Db    *gorm.DB
	Cache caching.Cache
}

func NewAdsenseChannelRP(db *gorm.DB) *AdsenseChannelRP {
	return &AdsenseChannelRP{Db: db}
}

func (t *AdsenseChannelRP) RemoveChannelByAccount(channels []string, account string) (err error) {
	err = t.Db.Delete(&model.AdsenseChannelModel{}, "channel NOT IN ? AND account = ?", channels, account).Error
	return
}

func (t *AdsenseChannelRP) AddAdsenseChannel(record *model.AdsenseChannelModel) (err error) {
	err = t.Db.Save(record).Error
	return
}

func (t *AdsenseChannelRP) GetAllAdsenseChannels() (records []model.AdsenseChannelModel) {
	t.Db.Model(&model.AdsenseChannelModel{}).Find(&records)
	return
}

func (t *AdsenseChannelRP) GetChannelsByAccount(account string) (channels []model.AdsenseChannelModel, err error) {
	t.Db.Model(&model.AdsenseChannelModel{}).
		// Debug().
		Where("account = ?", account).
		Order("id DESC").
		Find(&channels)
	return
}

// func (t *AdsenseChannelRP) GetChannelNotCampaign(account string) (record model.AdsenseChannelModel, err error) {
// 	// time_off = null hoặc time_off phải lớn hơn 10h từ lúc off
// 	time10 := time.Now().Add(time.Hour * time.Duration(-10))
// 	today := time.Now().Format("2006-01-02")
// 	err = t.Db.
// 		// Debug().
// 		Model(&model.AdsenseChannelModel{}).
// 		Where("account = ? AND campaign_id = 0 AND (time_off is null or (time_off < ? AND time_off < ?))", account, time10, today).
// 		First(&record).Error
// 	return
// }

// func (t *AdsenseChannelRP) RemoveAdsenseChannel(id int64) (err error) {
// 	err = t.Db.Delete(&model.AdsenseChannelModel{}, "id = ?", id).Error
// 	return
// }

func (t *AdsenseChannelRP) IsExistChannel(channel string, accountAdsense string) bool {
	var record model.AdsenseChannelModel
	t.Db.Model(&model.AdsenseChannelModel{}).Where("channel = ? AND account = ?", channel, accountAdsense).Find(&record)
	if record.Channel != "" && record.ID > 0 {
		return true
	}
	return false
}

// Update campaignID cho table adsense_channel
//
//	func (t *AdsenseChannelRP) UpdateCampaignForChannel(channel, account string, campaignID int64) (err error) {
//		if err = t.RemoveCampaignForChannel(campaignID); err != nil {
//			return
//		}
//		err = t.Db.
//			// Debug().
//			Model(&model.AdsenseChannelModel{}).
//			Where("channel = ? AND account = ?", channel, account).
//			Update("campaign_id", campaignID).Error
//		return
//	}

// update campaign_id = 0 của campaignID trong table AdsenseChannel nếu có
// func (t *AdsenseChannelRP) RemoveCampaignForChannel(campaignID int64) (err error) {
// 	err = t.Db.
// 		// Debug().
// 		Model(&model.AdsenseChannelModel{}).
// 		Where("campaign_id = ?", campaignID).
// 		Updates(map[string]interface{}{"campaign_id": 0, "time_off": time.Now()}).Error
// 	return
// }

func (t *AdsenseChannelRP) GetChannelsUse(account string) (channelsUse, total int64) {
	t.Db.
		// Debug().
		Model(&model.AdsenseChannelModel{}).
		Where("account = ?", account).
		Count(&total).
		Where("campaign_id != 0").
		Count(&channelsUse)
	return
}

func (t *AdsenseChannelRP) GetChannelsNotUse() (records []model.AdsenseChannelModel, total int64) {
	t.Db.
		Model(&model.AdsenseChannelModel{}).
		Where("campaign_id = 0").
		Count(&total).
		Find(&records)
	return
}

func (t *AdsenseChannelRP) GetChannelByChannel(channel, account string) (record model.AdsenseChannelModel, err error) {
	err = t.Db.
		Model(&model.AdsenseChannelModel{}).
		Where("channel = ? AND account = ?", channel, account).
		Find(&record).Error
	return
}

// update  toàn bộ campaignID = 0
func (t *AdsenseChannelRP) ResetCampaignIDToChannel() (err error) {
	err = t.Db.
		Model(&model.AdsenseChannelModel{}).
		Where("1 = 1").
		Update("campaign_id", 0).Error
	return
}

func (t *AdsenseChannelRP) CountChannelNotUse(account string) (total int64, err error) {
	err = t.Db.
		Model(&model.AdsenseChannelModel{}).
		Where("campaign_id = 0 AND account = ?", account).
		Count(&total).Error
	return
}

// count channel khả dụng
func (t *AdsenseChannelRP) CountChannelAvailableByAccount(account string) (total int64) {
	// time_off = null hoặc time_off phải lớn hơn 10h từ lúc off
	time10 := time.Now().Add(time.Hour * time.Duration(-10))
	today := time.Now().Format("2006-01-02")
	t.Db.
		// Debug().
		Model(&model.AdsenseChannelModel{}).
		Where("account = ? AND campaign_id = 0 AND (time_off is null or (time_off < ? AND time_off < ?))", account, time10, today).
		Count(&total)
	return
}

// check channel có bị trùng country với các campaign khác k?
func (t *AdsenseChannelRP) CheckCampaignCanUseChannel(campaignID int64, channel string, locations []string) bool {
	var record model.ChannelCountryModel
	t.Db.
		// Debug().
		Model(&model.ChannelCountryModel{}).
		Where("campaign_id != ? AND channel = ? AND country in ?", campaignID, channel, locations).
		Find(&record)
	if record.Channel != "" {
		return false
	}
	return true
}

func (t *AdsenseChannelRP) AddChannelCountry(record model.ChannelCountryModel) (err error) {
	// remove record old
	err = t.Db.
		Model(&model.ChannelCountryModel{}).
		Create(&record).Error
	return
}

func (t *AdsenseChannelRP) EditChannelCountry(record model.ChannelCountryModel) (err error) {
	err = t.Db.
		// Debug().
		Model(&model.ChannelCountryModel{}).
		Where("channel = ? AND country = ? AND time_off is null", record.Channel, record.Country).
		Updates(&record).Error
	return
}

func (t *AdsenseChannelRP) OffChannelCountryByCampaign(campaignID int64) (err error) {
	err = t.Db.
		// Debug().
		Model(&model.ChannelCountryModel{}).
		Where("campaign_id = ?", campaignID).
		Update("time_off", time.Now()).Error
	return
}

func (t *AdsenseChannelRP) RemoveChannelCountryByCampaign(campaignID int64, channel string) (err error) {
	err = t.Db.Delete(&model.ChannelCountryModel{}, "campaign_id = ? AND channel = ?", campaignID, channel).Error
	return
}

func (t *AdsenseChannelRP) FindChannelCountry(channel, country string) (record model.ChannelCountryModel, err error) {
	err = t.Db.
		//Debug().
		Model(&model.ChannelCountryModel{}).
		Where("channel = ? AND country = ?", channel, country).
		Find(&record).Error
	return
}

// kiểm tra channel có đang off strong 10 tiếng k?
func (t *AdsenseChannelRP) FindChannelOFF(channel string) (record model.ChannelCountryModel, err error) {
	err = t.Db.
		// Debug().
		Model(&model.ChannelCountryModel{}).
		Where("channel = ? AND time_off is not null", channel).
		Find(&record).Error
	return
}
