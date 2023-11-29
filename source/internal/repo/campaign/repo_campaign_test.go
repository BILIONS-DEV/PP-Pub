package campaign

import (
	"gorm.io/gorm"
	"log"
	"source/config"
	"source/infrastructure/caching"
	"source/infrastructure/mysql"
	// "source/infrastructure/fakedb"
	"source/internal/entity/dto"
	"testing"
)

func TestNewCampaignRepo(t *testing.T) {
	var (
		DB *gorm.DB
		// Cache caching.Cache
		err error
	)
	configs := config.NewConfig()

	// => init Aerospike
	if configs.Aerospike != nil {
		log.Println("init aerospike...", true)
		if _, err = caching.NewAerospike(caching.Config{
			Host:      configs.Aerospike.Host,
			Port:      configs.Aerospike.Port,
			Namespace: configs.Aerospike.Namespace,
		}); err != nil {
			log.Fatalln("Error connect Cache: ", err)
			return
		}
	}

	// => init Mysql
	if configs.Mysql != nil {
		log.Println("init mysql....", true)
		if DB, err = mysql.Connect(mysql.Config{
			Username: configs.Mysql.Username,
			Password: configs.Mysql.Password,
			Host:     configs.Mysql.Host,
			Port:     configs.Mysql.Port,
			Database: "apac_aff",
			Encoding: configs.Mysql.Encoding,
		}); err != nil {
			log.Fatalln("Error connect mysql: ", err)
			return
		}
	}

	// db, _ := fakedb.NewMysql()
	// ca, _ := fakedb.NewCache()
	repoCampaign := NewCampaignRepo(DB)
	payload := dto.PayloadCampaignAdd{
		Id:            81,
		Name:          "TTTTTT",
		TrafficSource: "Taboola",
		DemandSource:  "codefuel",
		Vertical:      "abc",
		UserFlow:      "3_click",
		LandingPages:  "https://linksvip.net/",
		MainKeyword:   "aa",
		Channel:       "",
		GD:            "",
		Keywords: []dto.CampaignKeywords{
			{
				// Id:         541,
				CampaignId: 75,
				Keyword:    "keyword 1",
			},
			{
				// Id:         541,
				CampaignId: 75,
				Keyword:    "keyword 2",
			},
		},
		// CreativeId: 0,
	}
	rec := payload.ToModel()
	repoCampaign.AddCampaign(&rec)

}
