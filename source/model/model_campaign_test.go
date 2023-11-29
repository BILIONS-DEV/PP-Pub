package model

import (
	"fmt"
	"source/infrastructure/fakedb"
	"testing"
)

func TestCampaignModel_ValidateModel(t *testing.T) {
	db, _ := fakedb.NewMysql()
	//db.AutoMigrate(&CampaignModel{}, &CampaignKeywordModel{})
	//return
	record := CampaignModel{
		ID:   1,
		Name: "campaign 3",
		//TrafficSource: "TrafficSource",
		//DemandSource:  "DemandSource",
		//PixelId:       "PixelId",
		//LandingPages:  "LandingPages",
		//MainKeyword:   "MainKeyword",
		//Channel:       "Channel",
		//GD:            "GD",
		//Params:        "Params",
		//UrlTrackImp:   "UrlTrackImp",
		//URLTrackClick: "URLTrackClick",
		Keywords: []CampaignKeywordModel{
			{ID: 1, Keyword: "keyword 1.0.2"},
			//{Keyword: "keyword 1.0.0"},
			//{Keyword: "keyword 2"},
			//{Keyword: "keyword 3"},
			//{Keyword: "keyword 5"},
		},
	}

	db.Debug().Save(&record)
	//db.Debug().Model(&record).Association("Keywords").Clear()

	//db.Debug().Select("Keywords").Delete(&CampaignModel{ID: 3}).Save(&record)
	//db.Debug().Select("Keywords").Delete(&CampaignModel{ID: 2})

	fmt.Println("record: ", record)
}
