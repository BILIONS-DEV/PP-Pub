package dto

import (
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
)

type PayloadCreativeFilter struct {
	datatable.Request
	OrderColumn int    `query:"order_column" json:"order_column" form:"order_column"`
	OrderDir    string `query:"order_dir" json:"order_dir" form:"order_dir"`
}

func (t *PayloadCreativeFilter) Validate() (err error) {
	return
}

type PayloadCreativeAdd struct {
	Id       int             `json:"id"`
	SiteName string          `json:"site_name"`
	Titles   []CreativeTitle `json:"titles"`
	Images   []CreativeImage `json:"images"`
}

type CreativeTitle struct {
	Id         int    `json:"id"`
	CreativeId int    `json:"creative_id"`
	Title      string `json:"title"`
}

type CreativeImage struct {
	Id         int    `json:"id"`
	CreativeId int    `json:"creative_id"`
	Image      string `json:"image"`
}

func (p *PayloadCreativeAdd) ToModel() model.CreativeModel {
	ID := int64(p.Id)
	var CreativeTitles []model.CreativeTitleModel
	for _, value := range p.Titles {
		if value.Title == "" {
			continue
		}
		var title model.CreativeTitleModel
		title.Title = value.Title
		CreativeTitles = append(CreativeTitles, title)
	}
	var CreativeImages []model.CreativeImageModel
	for _, value := range p.Images {
		if value.Image == "" {
			continue
		}
		var image model.CreativeImageModel
		image.Image = value.Image
		CreativeImages = append(CreativeImages, image)
	}

	record := model.CreativeModel{}
	record.ID = ID
	record.SiteName = p.SiteName
	record.Titles = CreativeTitles
	record.Images = CreativeImages

	return record
}

type PayloadCreativeSubmit struct {
	CampaignID int `json:"campaign_id" form:"campaign_id"`
	CreativeID int `json:"creative_id" form:"creative_id"`
	Page       int `query:"page" json:"page" form:"page"`
	Limit      int `query:"limit" json:"limit" form:"limit"`
	Start      int `query:"start" json:"start" form:"start"`
	Length     int `query:"length" json:"length" form:"length"`
}

type PayloadCreativeSubmitFilter struct {
	datatable.Request
	PostData PayloadCreativeSubmit `query:"postData"`
}

func (p *PayloadCreativeSubmit) ToModel() (record model.CreativeSubmitModel) {
	record.CampaignID = int64(p.CampaignID)
	record.CreativeID = int64(p.CreativeID)
	return
}

// func (t *PayloadCampaignFilter) ToCondition() map[string]interface{} {
// 	var condition = make(map[string]interface{})
// 	if t.SubID != "" {
// 		condition["subid"] = t.SubID
// 	}
//
// 	return condition
// }

// type PayloadGetHistory struct {
// 	Id     string  `json:"id" form:"id"`
// 	Object string `json:"object" form:"object"`
// }
