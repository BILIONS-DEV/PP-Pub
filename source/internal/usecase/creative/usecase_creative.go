package creative

import (
	"source/internal/entity/dto"
	"source/internal/entity/model"
	"source/internal/repo"
	"source/internal/repo/creative"
)

type UsecaseCreative interface {
	// Cronjob() (err error)
	Filter(payload *dto.PayloadCreativeFilter, user model.User) (reportBodis []model.CreativeModel, err error)
	GetCreativeById(id int64) (row model.CreativeModel)
	AddCreative(payload *dto.PayloadCreativeAdd, userID int64) (err error)
	EditCreative(payload *dto.PayloadCreativeAdd, userID int64) (err error)
	Delete(id int64) (err error)
	GetAll() (records []model.CreativeModel, err error)
	LoadCreativeSubmit(dataPost *dto.PayloadCreativeSubmit) (records []model.CreativeSubmitModel, totalRecord int, err error)
	UpdateCreativeSubmitForCampaigns(creativeID int64) (err error)
}

type creativeUsecase struct {
	repos *repo.Repositories
	// Trans *lang.Translation
}

func NewCreativeUsecase(repos *repo.Repositories) *creativeUsecase {
	return &creativeUsecase{repos: repos}
}

func (t *creativeUsecase) AddCreative(payload *dto.PayloadCreativeAdd, userID int64) (err error) {
	var record = payload.ToModel()
	record.UserID = userID
	if err = record.Validate(); err != nil {
		return
	}
	if err = t.repos.Creative.Save(&record); err != nil {
		return
	}
	return
}

func (t *creativeUsecase) EditCreative(payload *dto.PayloadCreativeAdd, userID int64) (err error) {
	var record = payload.ToModel()
	record.UserID = userID
	if err = record.Validate(); err != nil {
		return
	}

	var oldRecord model.CreativeModel
	if err = t.repos.Creative.FindByID(&oldRecord, record.ID); err != nil {
		return
	}

	if err = t.repos.Creative.RemoveRelationship(&oldRecord); err != nil {
		return
	}
	if err = t.repos.Creative.Save(&record); err != nil {
		return
	}
	// => update creative submit (table creative_submit) cho campaign đang select creative này
	err = t.UpdateCreativeSubmitForCampaigns(oldRecord.ID)
	return
}

func (t *creativeUsecase) Filter(payload *dto.PayloadCreativeFilter, user model.User) (reportBodis []model.CreativeModel, err error) {
	if err = payload.Validate(); err != nil {
		return
	}
	// chuyển đổi và validate dữ liệu time
	// gọi vào repo để thực hiện filter
	reportBodis, err = t.repos.Creative.FindByFilter(&creative.FilterInput{
		// Condition: payload.ToCondition(),
		// OrderBy:   payload.OrderBy,
		// GroupBy:   payload.GroupBy,
		// SubID:     payload.SubID,
		// StartDate: startDate,
		// EndDate:   endDate,
		// Order:     "id DESC",
	}, user)
	return
}

func (t *creativeUsecase) GetCreativeById(idSearch int64) (row model.CreativeModel) {
	// var CreativeSubmits = t.repos.Creative.ErrorTitle()
	// for _, value := range CreativeSubmits {
	// 	value.Title = strings.Replace(value.Title, "$$", "$", -1)
	// 	value.New = "new"
	// 	if err := t.repos.Creative.SaveCreativeSubmit(&value); err != nil {
	// 		return
	// 	}
	// }

	t.repos.Creative.FindByID(&row, idSearch)
	return row
}

func (t *creativeUsecase) Delete(id int64) (err error) {
	return t.repos.Creative.Remove(id)
}

func (t *creativeUsecase) GetAll() (records []model.CreativeModel, err error) {
	return t.repos.Creative.GetAll()
}

func (t *creativeUsecase) LoadCreativeSubmit(dataPost *dto.PayloadCreativeSubmit) (records []model.CreativeSubmitModel, totalRecord int, err error) {
	var record = model.CreativeSubmitFilter{
		CampaignID: dataPost.CampaignID,
		CreativeID: dataPost.CreativeID,
		Page:       dataPost.Page,
		Limit:      dataPost.Limit,
		Start:      dataPost.Start,
		Length:     dataPost.Length,
	}
	return t.repos.Creative.GetCreativeSubmitForCampaign(record)
}

func (t *creativeUsecase) UpdateCreativeSubmitForCampaigns(creativeID int64) (err error) {
	var creative model.CreativeModel
	if err = t.repos.Creative.FindByID(&creative, creativeID); err != nil {
		return
	}
	var campaigns []model.CampaignModel
	campaigns, err = t.repos.Campaign.GetCampaignsByCreative(creative.ID)
	if len(campaigns) == 0 || err != nil {
		return
	}
	for _, value := range campaigns {
		if err = t.UpdateCreativeSubmitForCampaign(creative, value.ID); err != nil {
			return
		}
	}
	return
}

func (t *creativeUsecase) UpdateCreativeSubmitForCampaign(creative model.CreativeModel, campaignID int64) (err error) {
	// ** trong table creative_submit có 2 field cần lưu ý
	// ** creative.New = "new" => mới đc add => hiển thị cho Thu Copy ở campaign
	// ** khi update hay add: thì set creative.Flag = 2, sau khi update hay add xong thì xóa toàn bộ creative.Flag = 1, và update lại toàn bộ creative.Flag = 2 về creative.Flag = 1, đây là bước xóa những thứ (title, site_name, image) không còn trong creative

	campaign := t.repos.Campaign.GetById(campaignID)
	if campaign.TrafficSource != "Outbrain" {
		if len(creative.Titles) > 0 {
			for _, title := range creative.Titles {
				var rec = model.CreativeSubmitModel{
					CampaignID: campaignID,
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
					CampaignID: campaignID,
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
				CampaignID: campaignID,
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
				CampaignID: campaignID,
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
	if err = t.repos.Creative.RemoveCreativeSubmitForCampaign(campaignID, creative.ID); err != nil {
		return
	}

	// update flag = 1 sau khi hoàn thành các bước update/add/remove
	err = t.repos.Creative.UpdateFlagCreativeSubmit(campaignID, creative.ID)
	return
}
