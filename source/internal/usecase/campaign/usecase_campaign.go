package campaign

import (
	// "fmt"
	"source/internal/entity/dto"
	"source/internal/entity/model"
	"source/internal/errors"
	"source/internal/lang"
	"source/internal/repo"
	"source/internal/repo/campaign"
	"strconv"
	"strings"
)

type UsecaseCampaign interface {
	// Cronjob() (err error)
	Filter(payload *dto.PayloadCampaignFilter, user model.User) (response []model.CampaignModel, total int64, err error)
	GetCampaignById(id int64) (row model.CampaignModel)
	AddCampaign(payload *dto.PayloadCampaignAdd) (err error)
	EditCampaign(payload *dto.PayloadCampaignAdd) (err error)
	Delete(id int64) (err error)
	DisableNotificationCreativeSubmit(id int64) (err error)
	UpdateParamsForCampaign(campaign model.CampaignModel, params string) (err error)
	UpdateUrlTrack(campaign model.CampaignModel, UrlTrackImp, urlTrackClick string) (err error)

	// tungdt
	AddCampaign2(record *model.CampaignModel) (err error)
	EditCampaign2(record *model.CampaignModel) (err error)
	GetNotificationCreativeSubmit(campaignID int64) (notify bool)
	CountriesOutbrain() (records []model.CountryOutbrainModel, err error)
	GetNewCampaignID() (newID int64)
	AccountByObject(object, object_type string) (records []model.AccountModel, err error)
	AllAcounts() (records []model.AccountModel, err error)
	AddLandingPagesForDemandSource(dataPost *dto.PayloadLandingPages, userID int64) (err error)
	LoadLandingPagesByDemand(demandSource string) (records []model.LandingPagesDemand, err error)
}

type campaignUsecase struct {
	repos *repo.Repositories
	Trans *lang.Translation
}

func NewCampaignUsecase(repos *repo.Repositories, trans *lang.Translation) *campaignUsecase {
	return &campaignUsecase{repos: repos, Trans: trans}
}

// func (t *campaignUsecase) Cronjob() (err error) {
// 	loc, _ := time.LoadLocation("America/New_York")
// 	return t.repos.ReportBodis.CronjobReport(&reportbodisRepo.DateRange{
// 		StartDate: time.Now().AddDate(0, 0, -1).In(loc),
// 		EndDate:   time.Now().In(loc),
// 	})
// }

func (t *campaignUsecase) AddCampaign2(record *model.CampaignModel) (err error) {
	// validate
	if err = record.Validate(); err != nil {
		return
	}
	// check exist trong db
	if t.repos.Campaign.IsExist(record.Name) {
		return errors.NewWithID("name already exists", "name")
	}
	// insert to db
	if err = t.repos.Campaign.Save(record); err != nil {
		return errors.New("could not write to database")
	}
	return
}

func (t *campaignUsecase) EditCampaign2(record *model.CampaignModel) (err error) {
	if err = record.Validate(); err != nil {
		return
	}
	if err = t.repos.Campaign.Save(record); err != nil {
		return errors.New("could not write to database")
	}
	return
}

// 	return
// func (t *campaignUsecase) AddCampaign(payload *dto.PayloadCampaignAdd) (row model.CampaignModel, err error) {
// 	var landingPagesGeo []campaign.LandingPageGeo
// 	for _, landingPageGeo := range payload.LandingPagesGeo {
// 		landingPagesGeo = append(landingPagesGeo, campaign.LandingPageGeo{
// 			ID:          landingPageGeo.ID,
// 			Geo:         landingPageGeo.Geo,
// 			LandingPage: landingPageGeo.LandingPage,
// 		})
// 	}
// 	return t.repos.Campaign.AddCampaign(&campaign.AddInput{
// 		Name:            payload.Name,
// 		TrafficSource:   payload.TrafficSource,
// 		DemandSource:    payload.DemandSource,
// 		PixelId:         payload.PixelId,
// 		MainKeyword:     payload.MainKeyword,
// 		Channel:         payload.Channel,
// 		GD:              payload.GD,
// 		LandingPages:    payload.LandingPages,
// 		Keywords:        payload.Keywords,
// 		LandingPagesGEO: landingPagesGeo,
// // 	})
// }

func (t *campaignUsecase) AddCampaign(payload *dto.PayloadCampaignAdd) (err error) {
	var campaignGroup = payload.ToCampaignGroupModel()
	if err = t.repos.Campaign.SaveCampaignGroup(&campaignGroup); err != nil {
		return
	}
	var newCampaignGroup model.CampaignGroupModel
	newCampaignGroup, err = t.repos.Campaign.GetNewCampaignGroup()

	var record = payload.ToModel()
	record.CampaignGroupID = newCampaignGroup.ID
	if err = record.Validate(); err != nil {
		return
	}
	// if err = t.repos.Campaign.Save(&record); err != nil {
	// 	return
	// }
	// Update name = ID: TrafficSource -> DemandSource - Name
	// var newCampaign model.CampaignModel
	// newCampaign, err = t.repos.Campaign.GetNewCampaign()
	// if !strings.Contains(newCampaign.Name, strings.Title(newCampaign.TrafficSource)+"->"+strings.Title(newCampaign.DemandSource)) {
	// 	newCampaign.Name = strconv.FormatInt(newCampaign.ID, 10) + ": " + strings.Title(newCampaign.TrafficSource) + " -> " + strings.Title(newCampaign.DemandSource) + " - " + newCampaign.Name
	// }
	// if err = t.repos.Campaign.Save(&newCampaign); err != nil {
	// 	return
	// }
	// err = t.UpdateCreativeSubmitByCampaign(&newCampaign)
	err = t.CreateCampaignsByCreative(payload, newCampaignGroup.ID)
	return
}

func (t *campaignUsecase) EditCampaign(payload *dto.PayloadCampaignAdd) (err error) {
	var record = payload.ToModel()
	if err = record.Validate(); err != nil {
		return
	}
	var oldRecord model.CampaignModel
	if err = t.repos.Campaign.FindByID(&oldRecord, record.ID); err != nil {
		return
	}
	if err = t.repos.Campaign.RemoveRelationship(&oldRecord); err != nil {
		return
	}
	record.CampaignGroupID = oldRecord.CampaignGroupID
	record.TrafficSourceID = oldRecord.TrafficSourceID
	record.Flag = oldRecord.Flag
	if err = t.repos.Campaign.Save(&record); err != nil {
		return
	}
	// err = t.UpdateCreativeSubmitByCampaign(&record)
	return
}

func (t *campaignUsecase) CreateCampaignsByCreative(payload *dto.PayloadCampaignAdd, CampaignGroupID int64) (err error) {
	var name = payload.Name
	var newCampaignID int64
	// tạo campagin dựa vào creative (số lượng title và image trong creative)
	// số lượng row creative tương ứng totalRow = total(title) * total(image)
	// số lượng campaign tương ứng => mỗi campaign sẽ có tối đa 10 row creative

	if payload.CreativeID == 0 {
		var record = payload.ToModel()
		// save campaign không có creative
		if err = t.repos.Campaign.Save(&record); err != nil {
			return
		}
		// update campaign name với campaignID
		record.Name = strconv.FormatInt(record.ID, 10) + ": " + name
		if err = t.repos.Campaign.UpdateName(record.ID, record.Name); err != nil {
			return
		}

		return
	} else {
		var numberRowCreative = 0
		var firstCampaign = true
		// tạo nhiều campaign phụ thuộc vào creative
		var creative model.CreativeModel
		if err = t.repos.Creative.FindByID(&creative, int64(payload.CreativeID)); err != nil {
			return
		}

		// if strings.ToLower(payload.TrafficSource) != "outbrain" && strings.ToLower(payload.TrafficSource) != "revcontent" && strings.ToLower(payload.TrafficSource) != "pocpoc" {
		// 	if len(creative.Titles) > 0 {
		// 		for _, title := range creative.Titles {
		// 			// create campaign
		// 			if firstCampaign || numberRowCreative == 10 {
		// 				record := payload.ToModel()
		// 				record.Flag = "new"
		// 				record.ID = 0
		// 				if err = t.repos.Campaign.Save(&record); err != nil {
		// 					return
		// 				}
		// 				newCampaignID = record.ID
		// 				// update campaign name với campaignID
		// 				record.Name = strconv.FormatInt(record.ID, 10) + ": " + name
		// 				if err = t.repos.Campaign.UpdateName(record.ID, record.Name); err != nil {
		// 					return
		// 				}
		//
		// 				firstCampaign = false
		// 				numberRowCreative = 0
		// 			}
		//
		// 			// newCampaign, _ := t.repos.Campaign.GetNewCampaign()
		// 			// create creative submit
		// 			var rec = model.CreativeSubmitModel{
		// 				CampaignID: newCampaignID,
		// 				CreativeID: creative.ID,
		// 				SiteName:   "",
		// 				Title:      t.repos.Creative.HandleTitleCreative(title.Title, payload.TrafficSource),
		// 				Image:      "",
		// 				New:        "new",
		// 				Flag:       1,
		// 			}
		// 			if err = t.repos.Creative.SaveCreativeSubmit(&rec); err != nil {
		// 				return
		// 			}
		// 			numberRowCreative = numberRowCreative + 1
		// 		}
		// 	}
		// 	return
		// }
		// update/add creativeSubmit
		if len(creative.Titles) > 0 && len(creative.Images) > 0 {
			for _, title := range creative.Titles {
				for _, image := range creative.Images {
					if firstCampaign || numberRowCreative == 10 {
						record := payload.ToModel()
						record.Flag = "new"
						record.ID = 0
						if err = t.repos.Campaign.Save(&record); err != nil {
							return
						}
						newCampaignID = record.ID
						// update campaign name với campaignID
						record.Name = strconv.FormatInt(record.ID, 10) + ": " + name
						if err = t.repos.Campaign.UpdateName(record.ID, record.Name); err != nil {
							return
						}

						firstCampaign = false
						numberRowCreative = 0
					}
					// newCampaign, _ := t.repos.Campaign.GetNewCampaign()
					// create creative submit
					var rec = model.CreativeSubmitModel{
						CampaignID: newCampaignID,
						CreativeID: creative.ID,
						SiteName:   creative.SiteName,
						Title:      t.repos.Creative.HandleTitleCreative(title.Title, payload.TrafficSource),
						Image:      "https://ul.pubpowerplatform.io" + image.Image,
						New:        "new",
						Flag:       1,
					}
					if err = t.repos.Creative.SaveCreativeSubmit(&rec); err != nil {
						return
					}
					numberRowCreative = numberRowCreative + 1
				}
			}
		} else if len(creative.Titles) > 0 && len(creative.Images) == 0 {
			for _, title := range creative.Titles {
				if firstCampaign || numberRowCreative == 10 {
					record := payload.ToModel()
					record.Flag = "new"
					record.ID = 0
					if err = t.repos.Campaign.Save(&record); err != nil {
						return
					}
					newCampaignID = record.ID
					// update campaign name với campaignID
					record.Name = strconv.FormatInt(record.ID, 10) + ": " + name
					if err = t.repos.Campaign.UpdateName(record.ID, record.Name); err != nil {
						return
					}

					firstCampaign = false
					numberRowCreative = 0
				}

				// newCampaign, _ := t.repos.Campaign.GetNewCampaign()
				// create creative submit
				var rec = model.CreativeSubmitModel{
					CampaignID: newCampaignID,
					CreativeID: creative.ID,
					SiteName:   creative.SiteName,
					Title:      t.repos.Creative.HandleTitleCreative(title.Title, payload.TrafficSource),
					Image:      "",
					New:        "new",
					Flag:       1,
				}
				if err = t.repos.Creative.SaveCreativeSubmit(&rec); err != nil {
					return
				}
				numberRowCreative = numberRowCreative + 1
			}
		} else if len(creative.Titles) == 0 && len(creative.Images) > 0 {
			for _, image := range creative.Images {
				if firstCampaign || numberRowCreative == 10 {
					record := payload.ToModel()
					record.Flag = "new"
					record.ID = 0
					if err = t.repos.Campaign.Save(&record); err != nil {
						return
					}
					newCampaignID = record.ID
					// update campaign name với campaignID
					record.Name = strconv.FormatInt(record.ID, 10) + ": " + name
					if err = t.repos.Campaign.UpdateName(record.ID, record.Name); err != nil {
						return
					}

					firstCampaign = false
					numberRowCreative = 0
				}

				// newCampaign, _ := t.repos.Campaign.GetNewCampaign()
				// create creative submit
				var rec = model.CreativeSubmitModel{
					CampaignID: newCampaignID,
					CreativeID: creative.ID,
					SiteName:   creative.SiteName,
					Title:      "",
					Image:      "https://ul.pubpowerplatform.io" + image.Image,
					New:        "new",
					Flag:       1,
				}
				if err = t.repos.Creative.SaveCreativeSubmit(&rec); err != nil {
					return
				}
				numberRowCreative = numberRowCreative + 1
			}
		}
	}
	return
}

func (t *campaignUsecase) UpdateCreativeSubmitByCampaign(campaign *model.CampaignModel) (err error) {
	if campaign.CreativeID == 0 {
		// remove all creative submit for camp
		err = t.repos.Creative.RemoveAllCreativeSubmitByCampaign(campaign.ID)
		return
	}
	err = t.UpdateCreativeSubmitForCampaign(campaign)
	return
}

func (t *campaignUsecase) UpdateCreativeSubmitForCampaign(campaign *model.CampaignModel) (err error) {
	// ** trong table creative_submit có 2 field cần lưu ý
	// ** creative.New = "new" => mới đc add => hiển thị cho Thu Copy ở campaign
	// ** khi update hay add: thì set creative.Flag = 2, sau khi update hay add xong thì xóa toàn bộ creative.Flag = 1, và update lại toàn bộ creative.Flag = 2 về creative.Flag = 1, đây là bước xóa những thứ (title, site_name, image) không còn trong creative
	var creative model.CreativeModel
	if err = t.repos.Creative.FindByID(&creative, campaign.CreativeID); err != nil {
		return
	}

	if strings.ToLower(campaign.TrafficSource) != "outbrain" && strings.ToLower(campaign.TrafficSource) != "revcontent" && strings.ToLower(campaign.TrafficSource) != "pocpoc" {
		if len(creative.Titles) > 0 {
			for _, title := range creative.Titles {
				var rec = model.CreativeSubmitModel{
					CampaignID: campaign.ID,
					CreativeID: creative.ID,
					SiteName:   "",
					Title:      t.repos.Creative.HandleTitleCreative(title.Title, campaign.TrafficSource),
					Image:      "",
				}
				// check exist
				if record, isExist := t.repos.Creative.IsExistCreativeSubmit(rec); isExist {
					rec.ID = record.ID
					rec.Flag = 2
					rec.New = record.New
					if err = t.repos.Creative.SaveCreativeSubmit(&rec); err != nil {
						return
					}
					continue
				}
				// create creative submit
				rec.New = "new"
				rec.Flag = 2
				if err = t.repos.Creative.SaveCreativeSubmit(&rec); err != nil {
					return
				}
			}
		}

		// remove creativeSubmit
		if err = t.repos.Creative.RemoveCreativeSubmitForCampaign(campaign.ID, creative.ID); err != nil {
			return
		}
		// update flag = 1 sau khi hoàn thành các bước update/add/remove
		err = t.repos.Creative.UpdateFlagCreativeSubmit(campaign.ID, creative.ID)
		return
	}

	// update/add creativeSubmit
	if len(creative.Titles) > 0 && len(creative.Images) > 0 {
		for _, title := range creative.Titles {
			for _, image := range creative.Images {
				var rec = model.CreativeSubmitModel{
					CampaignID: campaign.ID,
					CreativeID: creative.ID,
					SiteName:   creative.SiteName,
					Title:      t.repos.Creative.HandleTitleCreative(title.Title, campaign.TrafficSource),
					Image:      "https://ul.pubpowerplatform.io" + image.Image,
				}
				// check exist
				if record, isExist := t.repos.Creative.IsExistCreativeSubmit(rec); isExist {
					rec.ID = record.ID
					rec.Flag = 2
					rec.New = record.New
					if err = t.repos.Creative.SaveCreativeSubmit(&rec); err != nil {
						return
					}
					continue
				}
				// create creative submit
				rec.New = "new"
				rec.Flag = 2
				if err = t.repos.Creative.SaveCreativeSubmit(&rec); err != nil {
					return
				}
			}
		}
	} else if len(creative.Titles) > 0 && len(creative.Images) == 0 {
		for _, title := range creative.Titles {
			var rec = model.CreativeSubmitModel{
				CampaignID: campaign.ID,
				CreativeID: creative.ID,
				SiteName:   creative.SiteName,
				Title:      t.repos.Creative.HandleTitleCreative(title.Title, campaign.TrafficSource),
				Image:      "",
			}
			// check exist
			if record, isExist := t.repos.Creative.IsExistCreativeSubmit(rec); isExist {
				rec.ID = record.ID
				rec.Flag = 2
				rec.New = record.New
				if err = t.repos.Creative.SaveCreativeSubmit(&rec); err != nil {
					return
				}
				continue
			}
			// create creative submit
			rec.New = "new"
			rec.Flag = 2
			if err = t.repos.Creative.SaveCreativeSubmit(&rec); err != nil {
				return
			}
		}
	} else if len(creative.Titles) == 0 && len(creative.Images) > 0 {
		for _, image := range creative.Images {
			var rec = model.CreativeSubmitModel{
				CampaignID: campaign.ID,
				CreativeID: creative.ID,
				SiteName:   creative.SiteName,
				Title:      "",
				Image:      "https://ul.pubpowerplatform.io" + image.Image,
			}
			// check exist
			if record, isExist := t.repos.Creative.IsExistCreativeSubmit(rec); isExist {
				rec.ID = record.ID
				rec.Flag = 2
				rec.New = record.New
				if err = t.repos.Creative.SaveCreativeSubmit(&rec); err != nil {
					return
				}
				continue
			}
			// create creative submit
			rec.New = "new"
			rec.Flag = 2
			if err = t.repos.Creative.SaveCreativeSubmit(&rec); err != nil {
				return
			}
		}
	}

	// remove creativeSubmit
	if err = t.repos.Creative.RemoveCreativeSubmitForCampaign(campaign.ID, creative.ID); err != nil {
		return
	}

	// update flag = 1 sau khi hoàn thành các bước update/add/remove
	err = t.repos.Creative.UpdateFlagCreativeSubmit(campaign.ID, creative.ID)
	return
}

func (t *campaignUsecase) Filter(payload *dto.PayloadCampaignFilter, user model.User) (response []model.CampaignModel, total int64, err error) {
	if err = payload.Validate(); err != nil {
		return
	}
	// chuyển đổi và validate dữ liệu time
	// gọi vào repo để thực hiện filter
	response, total, err = t.repos.Campaign.FindByFilter(&campaign.FilterInput{
		TrafficSource: payload.PostData.TrafficSource,
		DemandSource:  payload.PostData.DemandSource,
		Search:        payload.PostData.QuerySearch,
		Offset:        payload.Start,
		Limit:         payload.Length,
		Order:         payload.OrderString(),
	}, user)
	return
}

func (t *campaignUsecase) AddLandingPagesForDemandSource(dataPost *dto.PayloadLandingPages, userID int64) (err error) {
	if err = dataPost.Validate(); err != nil {
		return
	}
	// remove all landing pages by demand
	err = t.repos.Campaign.RemoveLandingPagesByDemand(dataPost.DemandSource)
	if err != nil {
		return
	}
	// add landing pages
	for index, value := range dataPost.LandingPages {
		if strings.TrimSpace(value) == "" {
			continue
		}
		err = t.repos.Campaign.AddLandingPages(&model.LandingPagesDemand{
			DemandSource: dataPost.DemandSource,
			UserID:       userID,
			Name:         dataPost.Name[index],
			LandingPage:  value,
			MainKeyword:  dataPost.MainKeyword[index],
		})
		if err != nil {
			return
		}
	}
	return
}

func (t *campaignUsecase) LoadLandingPagesByDemand(demandSource string) (records []model.LandingPagesDemand, err error) {
	if demandSource == "" {
		err = errors.New("Demand Source is empty")
	}
	records, err = t.repos.Campaign.LoadLandingPagesByDemand(demandSource)
	return
}

func (t *campaignUsecase) GetCampaignById(idSearch int64) (row model.CampaignModel) {
	t.repos.Campaign.FindByID(&row, idSearch)
	return row
}

func (t *campaignUsecase) Delete(id int64) (err error) {
	return t.repos.Campaign.Remove(id)
}

func (t *campaignUsecase) UpdateUrlTrack(campaign model.CampaignModel, UrlTrackImp, urlTrackClick string) (err error) {
	return t.repos.Campaign.UpdateUrlTrack(campaign, UrlTrackImp, urlTrackClick)
}

func (t *campaignUsecase) UpdateParamsForCampaign(campaign model.CampaignModel, params string) (err error) {
	return t.repos.Campaign.UpdateParamsForCampaignId(campaign, params)
}

func (t *campaignUsecase) DisableNotificationCreativeSubmit(id int64) (err error) {
	return t.repos.Creative.DisableNotificationCreativeSubmitByCampaign(id)
}

func (t *campaignUsecase) GetNotificationCreativeSubmit(campaignID int64) (notify bool) {
	return t.repos.Creative.GetNotificationCreativeSubmit(campaignID)
}

func (t *campaignUsecase) CountriesOutbrain() (records []model.CountryOutbrainModel, err error) {
	return t.repos.Country.GetAllCountryOutbrain()
}

func (t *campaignUsecase) GetNewCampaignID() (newID int64) {
	return t.repos.Campaign.GetNewCampaignID()
}

func (t *campaignUsecase) AccountByObject(object, object_type string) (records []model.AccountModel, err error) {
	if records, err = t.repos.Campaign.AccountByObject(object, object_type); err != nil {
		return
	}
	if len(records) == 0 {
		return t.repos.Campaign.AccountByObject(object, "")
	}
	return
}

func (t *campaignUsecase) AllAcounts() (records []model.AccountModel, err error) {
	return t.repos.Campaign.AllAcounts()
}
