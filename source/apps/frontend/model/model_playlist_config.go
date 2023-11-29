package model

import (
	// "gorm.io/gorm"
	// "source/apps/frontend/payload"
	"source/core/technology/mysql"
	// "source/pkg/utility"
)

type PlaylistConfig struct{}

type PlaylistConfigRecord struct {
	mysql.TablePlaylistConfig
}

func (PlaylistConfigRecord) TableName() string {
	return mysql.Tables.PlaylistConfig
}

func (t *PlaylistConfig) CreatePlaylistConfig(PlaylistConfig PlaylistConfigRecord) (err error) {
	err = mysql.Client.Create(&PlaylistConfig).Error
	return
}

func (t *PlaylistConfig) DeletePlaylistConfig(PlaylistConfig PlaylistConfigRecord) (err error) {
	err = mysql.Client.Where(PlaylistConfig).Delete(PlaylistConfigRecord{}).Error
	return
}

func (t *PlaylistConfig) GetPlaylistConfigPlaylist(PlaylistId int64) (records []PlaylistConfigRecord) {
	mysql.Client.Where(PlaylistConfigRecord{mysql.TablePlaylistConfig{PlaylistId: PlaylistId}}).Find(&records)
	return
}

func (t *PlaylistConfig) GetPlaylistConfigLanguage(LanguageId int64) (records []PlaylistConfigRecord) {
	mysql.Client.Where(PlaylistConfigRecord{mysql.TablePlaylistConfig{LanguageId: LanguageId}}).Find(&records)
	return
}

func (t *PlaylistConfig) GetPlaylistConfigChannels(ChannelsId int64) (records []PlaylistConfigRecord) {
	mysql.Client.Where(PlaylistConfigRecord{mysql.TablePlaylistConfig{ChannelsId: ChannelsId}}).Find(&records)
	return
}

func (t *PlaylistConfig) GetByPlaylistId(PlaylistId int64) (records []PlaylistConfigRecord) {
	mysql.Client.Where(PlaylistConfigRecord{mysql.TablePlaylistConfig{PlaylistId: PlaylistId}}).Find(&records)
	return
}

// 
// func (t *PlaylistConfig) GetAllPlaylistConfigChannelsValidate(userId int64, ChannelsId int64) (records []PlaylistConfigRecord) {
// 	mysql.Client.Where("Channels_id != 0 and user_id = ? and Channels_id != ?", userId, ChannelsId).Find(&records)
// 	return
// }
// 
// func (t *PlaylistConfig) GetAllPlaylistConfigChannels(userId int64) (records []PlaylistConfigRecord) {
// 	mysql.Client.Where("Channels_id != 0 and user_id = ?", userId).Find(&records)
// 	return
// }

// func (t *PlaylistConfig) GetAllByUser(userId int64, option string, PlaylistConfig payload.FilterPlaylistConfig) (listId []payload.ListPlaylistConfig, listKeywordsFilter []int64) {
// 	switch option {
// 	case "Channels":
// 		mysql.Client.Select("id", "name").Model(InventoryRecord{}).Where("user_id = ?", userId).Find(&listId)
// 		return
// 	case "Language":
// 		mysql.Client.Select("id", "name").Model(AdTypeRecord{}).Find(&listId)
// 		return
// 	case "Category":
// 		mysql.Client.Select("id", "name").Model(CategoryRecord{}).Find(&listId)
// 		return
// 	case "Keywords":
// 		// Tạo list id type cho 2 định loại banner và video
// 		var listIdTypeDisplay []int64
// 		var listIdTypeVideo []int64
// 		var listIdTypeSticky []int64
// 
// 		for _, format := range PlaylistConfig.Format {
// 			if format == 1 {
// 				listIdTypeDisplay = append(listIdTypeDisplay, format)
// 			} else if format == 5 {
// 				listIdTypeSticky = append(listIdTypeSticky, format)
// 			} else {
// 				listIdTypeVideo = append(listIdTypeVideo, format)
// 			}
// 		}
// 		// Filter for display trường hợp này sử dụng cho cả filter format là display hoặc format all
// 		if len(listIdTypeDisplay) > 0 || len(PlaylistConfig.Format) == 0 {
// 			var listKeywordsFilterDisplay []int64
// 			mysql.Client.Model(&KeywordsRecord{}).Select("inventory_ad_tag.id").
// 				Where("user_id = ?", userId).
// 				Scopes(
// 					setFilterInventory(PlaylistConfig),
// 					setFilterFormat(listIdTypeDisplay),
// 					setFilterSize(PlaylistConfig),
// 				).
// 				Group("inventory_ad_tag.id").
// 				Find(&listKeywordsFilterDisplay)
// 			listKeywordsFilter = append(listKeywordsFilter, listKeywordsFilterDisplay...)
// 		}
// 
// 		// Filter for video, nếu trong trường hợp filter format != all và có format cho video thì filter theo dạng video không cần tính đến size
// 		if len(listIdTypeVideo) > 0 {
// 			var listKeywordsFilterVideo []int64
// 			mysql.Client.Model(&KeywordsRecord{}).Select("inventory_ad_tag.id").
// 				Where("user_id = ?", userId).
// 				Scopes(
// 					setFilterInventory(PlaylistConfig),
// 					setFilterFormat(listIdTypeVideo),
// 				).
// 				Group("inventory_ad_tag.id").
// 				Find(&listKeywordsFilterVideo)
// 			listKeywordsFilter = append(listKeywordsFilter, listKeywordsFilterVideo...)
// 		}
// 		if len(listIdTypeSticky) > 0 || len(PlaylistConfig.Format) == 0 {
// 			var listKeywordsFilterSticky []int64
// 			mysql.Client.Model(&KeywordsRecord{}).Select("inventory_ad_tag.id").
// 				Where("user_id = ?", userId).
// 				Scopes(
// 					setFilterInventory(PlaylistConfig),
// 					setFilterFormat(listIdTypeDisplay),
// 					setFilterSizeSticky(PlaylistConfig),
// 				).
// 				Group("inventory_ad_tag.id").
// 				Find(&listKeywordsFilterSticky)
// 			listKeywordsFilter = append(listKeywordsFilter, listKeywordsFilterSticky...)
// 		}
// 		mysql.Client.Select("id", "name").Model(KeywordsRecord{}).Scopes(
// 			func(db *gorm.DB) *gorm.DB {
// 				if len(PlaylistConfig.Inventory) > 0 || len(PlaylistConfig.Format) > 0 || len(PlaylistConfig.Size) > 0 {
// 					return db.Where("id IN ?", listKeywordsFilter)
// 				}
// 				return db
// 			},
// 		).Where("user_id = ?", userId).Find(&listId)
// 		return
// 	case "Videos":
// 		mysql.Client.Select("id", "name").Model(VideosRecord{}).Find(&listId)
// 		return
// 	case "geography":
// 		mysql.Client.Select("id", "name").Model(CountryRecord{}).Find(&listId)
// 		return
// 	}
// 	return
// }

// func (t *PlaylistConfig) GetPlaylistConfigByFilterLanguage(inputs *payload.LanguageFilterPayload, userId int64) (records []PlaylistConfigRecord, err error) {
//
// 	err = mysql.Client.Select("language_id").Where("user_id = ? AND language_id != 0", userId).
// 		Scopes(
// 			t.setFilterChannels(inputs),
// 			t.setFilterLanguage(inputs),
// 			t.setFilterCategory(inputs),
// 			t.setFilterKeywords(inputs),
// 			t.setFilterVideos(inputs),
// 		).Debug().
// 		Group("language_id").
// 		Find(&records).Error
// 	return
// }

// func (t *PlaylistConfig) GetPlaylistConfigByFilterChannels(inputs *payload.ChannelsFilterPayload, userId int64) (records []PlaylistConfigRecord, err error) {
// 	err = mysql.Client.Select("Channels_id").Where("user_id = ? AND Channels_id != 0", userId).
// 		Scopes(
// 			t.setFilterChannelsChannels(inputs),
// 		).Debug().
// 		Group("Channels_id").
// 		Find(&records).Error
// 	return
// }

// func (t *PlaylistConfig) setFilterChannelsChannels(inputs *payload.ChannelsFilterPayload) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		if inputs.PostData.Channels != nil {
// 			switch inputs.PostData.Channels.(type) {
// 			case string, int:
// 				if inputs.PostData.Channels != "" {
// 					return db.Where("PlaylistConfig.inventory_id = ? or PlaylistConfig.inventory_id = -1", inputs.PostData.Channels)
// 				}
// 			case []string, []interface{}:
// 				return db.Where("PlaylistConfig.inventory_id IN ? or PlaylistConfig.inventory_id = -1", inputs.PostData.Channels)
// 			}
// 		}
// 		return db
// 	}
// }

// func (t *PlaylistConfig) setFilterChannels(inputs *payload.LanguageFilterPayload) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		if inputs.PostData.Channels != nil {
// 			switch inputs.PostData.Channels.(type) {
// 			case string, int:
// 				if inputs.PostData.Channels != "" {
// 					return db.Where("PlaylistConfig.inventory_id = ? or PlaylistConfig.inventory_id = -1", inputs.PostData.Channels)
// 				}
// 			case []string, []interface{}:
// 				return db.Where("PlaylistConfig.inventory_id IN ? or PlaylistConfig.inventory_id = -1", inputs.PostData.Channels)
// 			}
// 		}
// 		return db
// 	}
// }

// func (t *PlaylistConfig) setFilterLanguage(inputs *payload.LanguageFilterPayload) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		if inputs.PostData.Language != nil {
// 			switch inputs.PostData.Language.(type) {
// 			case string, int:
// 				if inputs.PostData.Language != "" {
// 					return db.Where("PlaylistConfig.ad_format_id = ? or PlaylistConfig.ad_format_id = -1", inputs.PostData.Language)
// 				}
// 			case []string, []interface{}:
// 				return db.Where("PlaylistConfig.ad_format_id IN ? or PlaylistConfig.ad_format_id = -1", inputs.PostData.Language)
// 			}
// 		}
// 		return db
// 	}
// }

// func (t *PlaylistConfig) setFilterCategory(inputs *payload.LanguageFilterPayload) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		if inputs.PostData.Category != nil {
// 			switch inputs.PostData.Category.(type) {
// 			case string, int:
// 				if inputs.PostData.Category != "" {
// 					return db.Where("PlaylistConfig.ad_size_id = ? or PlaylistConfig.ad_size_id = -1", inputs.PostData.Category)
// 				}
// 			case []string, []interface{}:
// 				return db.Where("PlaylistConfig.ad_size_id IN ? or PlaylistConfig.ad_size_id = -1", inputs.PostData.Category)
// 			}
// 		}
// 		return db
// 	}
// }

// func (t *PlaylistConfig) setFilterKeywords(inputs *payload.LanguageFilterPayload) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		if inputs.PostData.Keywords != nil {
// 			switch inputs.PostData.Keywords.(type) {
// 			case string, int:
// 				if inputs.PostData.Keywords != "" {
// 					return db.Where("PlaylistConfig.tag_id = ? or PlaylistConfig.tag_id = -1", inputs.PostData.Keywords)
// 				}
// 			case []string, []interface{}:
// 				return db.Where("PlaylistConfig.tag_id IN ? or PlaylistConfig.tag_id = -1", inputs.PostData.Keywords)
// 			}
// 		}
// 		return db
// 	}
// }

// func (t *PlaylistConfig) setFilterVideos(inputs *payload.LanguageFilterPayload) func(db *gorm.DB) *gorm.DB {
// 	return func(db *gorm.DB) *gorm.DB {
// 		if inputs.PostData.Videos != nil {
// 			switch inputs.PostData.Videos.(type) {
// 			case string, int:
// 				if inputs.PostData.Videos != "" {
// 					return db.Where("PlaylistConfig.Videos_id = ? or PlaylistConfig.Videos_id = -1", inputs.PostData.Videos)
// 				}
// 			case []string, []interface{}:
// 				return db.Where("PlaylistConfig.Videos_id IN ? or PlaylistConfig.Videos_id = -1", inputs.PostData.Videos)
// 			}
// 		}
// 		return db
// 	}
// }

// func (t *PlaylistConfig) HandleForChannels(recOld []InventoryRecord, user UserRecord) (records []InventoryRecord) {
// 	var listInventoryPlaylistConfiged []int64
// 	PlaylistConfigs := new(PlaylistConfig).GetAllPlaylistConfigChannels(user.Id)
// 	for _,PlaylistConfig := range PlaylistConfigs{
// 		listInventoryPlaylistConfiged = append(listInventoryPlaylistConfiged, PlaylistConfig.InventoryId)
// 	}
// 	for _, inventory := range recOld {
// 		if utility.InArray(inventory.Id, listInventoryPlaylistConfiged, false) {
// 			inventory.Status = 3
// 		}
// 		records = append(records, inventory)
// 	}
// 	return
// }