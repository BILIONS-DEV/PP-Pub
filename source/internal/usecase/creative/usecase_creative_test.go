package creative

import (
	"gorm.io/gorm"
	"log"
	"source/internal/repo"
	// "source/internal/repo/campaign"

	// "source/apps/worker/log_ip_ranges/model"
	"source/config"
	"source/infrastructure/caching"
	"source/infrastructure/mysql"
	"source/internal/entity/dto"
	"source/internal/repo/creative"
	"testing"
)

// func TestCreativeUsecase_AddCreative2(t *testing.T) {
// 	recEx := model.CreativeModel{
// 		Name:          "BBBB",
// 		TrafficSource: "Mgid",
// 		DemandSource:  "system1",
// 	}
// 	campRepo := mocks.NewRepoCreative(t)
// 	campRepo.On("IsExist", recEx.Name).Return(false)
// 	campRepo.On("Save", &recEx).Return(nil)
// 	repos := &repo.Repositories{
// 		Creative: campRepo,
// 	}
// 	campUC := NewCreativeUsecase(repos, nil)
//
// 	rec := model.CreativeModel{
// 		Name:          "BBBB",
// 		TrafficSource: "Mgid",
// 		DemandSource:  "system1",
// 	}
// 	err := campUC.AddCreative2(&rec)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	// repos := repo.Repositories{Creative: }
// 	// campUC := NewCreativeUsecase()
// }
//
// func TestCreativeUsecase_EditCreative2(t *testing.T) {
// 	rec := model.CreativeModel{
// 		ID:            51,
// 		Name:          "BBBB",
// 		TrafficSource: "Mgid",
// 		DemandSource:  "system1",
// 	}
// 	campRepo := mocks.NewRepoCreative(t)
// 	// campRepo.On("IsExist", "BBBB").Return(false)
// 	campRepo.On("Save", rec).Return(nil)
// 	repos := &repo.Repositories{
// 		Creative: campRepo,
// 	}
// 	campUC := NewCreativeUsecase(repos, nil)
//
// 	err := campUC.EditCreative2(&rec)
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func TestCreativeUsecase_UpdateCreativeSubmit(t *testing.T) {
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
			Database: configs.Mysql.Database,
			Encoding: configs.Mysql.Encoding,
		}); err != nil {
			log.Fatalln("Error connect mysql: ", err)
			return
		}
	}

	creativeRepo := creative.NewCreativeRepo(DB)
	// campaignRepo := campaign.NewCampaignRepo(DB)
	repos := &repo.Repositories{
		Creative: creativeRepo,
	}
	creativeUC := NewCreativeUsecase(repos)

	// creatives, err := creativeRepo.GetAll()
	// for _, creaTive := range creatives {
	// 	campaign := campaignRepo.GetById(creaTive.CampaignID)
	// 	if campaign.ID == 0 {
	// 		continue
	// 	}
	//
	// 	if len(creaTive.Titles) > 0 && len(creaTive.Images) > 0 {
	// 		for _, title := range creaTive.Titles {
	// 			for _, image := range creaTive.Images {
	// 				var rec = model.CreativeSubmitModel{
	// 					CampaignID: creaTive.CampaignID,
	// 					CreativeID: creaTive.ID,
	// 					SiteName:   creaTive.SiteName,
	// 					Title:      title.Title,
	// 					Image:      image.Image,
	// 				}
	// 				// check exist
	// 				if _, isExist := creativeRepo.IsExistCreativeSubmit(rec); isExist {
	// 					continue
	// 				}
	// 				// create creative submit
	// 				if err = creativeRepo.SaveCreativeSubmit(&rec); err != nil {
	// 					fmt.Printf("%+v\n", err)
	// 				}
	// 			}
	// 		}
	// 	} else if len(creaTive.Titles) > 0 && len(creaTive.Images) == 0 {
	// 		for _, title := range creaTive.Titles {
	// 			var rec = model.CreativeSubmitModel{
	// 				CampaignID: creaTive.CampaignID,
	// 				CreativeID: creaTive.ID,
	// 				SiteName:   creaTive.SiteName,
	// 				Title:      title.Title,
	// 				Image:      "",
	// 			}
	// 			// check exist
	// 			if _, isExist := creativeRepo.IsExistCreativeSubmit(rec); isExist {
	// 				continue
	// 			}
	// 			// create creative submit
	// 			if err = creativeRepo.SaveCreativeSubmit(&rec); err != nil {
	// 				fmt.Printf("%+v\n", err)
	// 			}
	// 		}
	// 	} else if len(creaTive.Titles) == 0 && len(creaTive.Images) > 0 {
	// 		for _, image := range creaTive.Images {
	// 			var rec = model.CreativeSubmitModel{
	// 				CampaignID: creaTive.CampaignID,
	// 				CreativeID: creaTive.ID,
	// 				SiteName:   creaTive.SiteName,
	// 				Title:      "",
	// 				Image:      image.Image,
	// 			}
	// 			// check exist
	// 			if _, isExist := creativeRepo.IsExistCreativeSubmit(rec); isExist {
	// 				continue
	// 			}
	// 			// create creative submit
	// 			if err = creativeRepo.SaveCreativeSubmit(&rec); err != nil {
	// 				fmt.Printf("%+v\n", err)
	// 			}
	// 		}
	// 	}
	// }

	record := dto.PayloadCreativeAdd{
		Id:       0,
		SiteName: "",
		Titles:   nil,
		Images:   nil,
	}
	if err = creativeUC.EditCreative(&record); err != nil {
		t.Error(err)
	}
}
