package model

import (
	"fmt"
	"source/apps/frontend/lang"
	"source/apps/frontend/payload"
	"source/core/technology/mysql"
	"source/pkg/datatable"
	"source/pkg/htmlblock"
	"source/pkg/pagination"
	"source/pkg/utility"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type Template struct{}

type TemplateRecord struct {
	mysql.TableTemplate
}

func (TemplateRecord) TableName() string {
	return mysql.Tables.Template
}

func (t *Template) GetByFilters(inputs *payload.PlayerFilterPayload, user UserRecord, lang lang.Translation) (response datatable.Response, err error) {
	var templates []TemplateRecord
	var total int64
	err = mysql.Client.Where("user_id = ?", user.Id).
		Scopes(
			t.setFilterSearch(inputs),
			t.setFilterType(inputs),
			t.setFilterSize(inputs),
			t.setDefault(user),
		).
		Model(&templates).Count(&total).
		Scopes(
			t.setOrder(inputs),
			pagination.Paginate(pagination.Params{
				Limit:  inputs.Length,
				Offset: inputs.Start,
			}),
		).
		Find(&templates).Error
	if err != nil {
		if !utility.IsWindow() {
			err = fmt.Errorf(lang.Errors.TemplateError.List.ToString())
		}
		return datatable.Response{}, err
	}
	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.MakeResponseDatatable(templates, inputs, user)
	return
}

func (t *Template) GetByFiltersV2(inputs *payload.PlayerFilterPayload, user UserRecord, lang lang.Translation) (response datatable.Response, err error) {
	var templates []TemplateRecord
	var total int64
	err = mysql.Client.Where("user_id = ?", user.Id).
		Scopes(
			t.setFilterSearch(inputs),
			t.setFilterType(inputs),
			t.setFilterSize(inputs),
			t.setDefault(user),
		).
		Model(&templates).Count(&total).
		Scopes(
			t.setOrder(inputs),
			pagination.Paginate(pagination.Params{
				Limit:  inputs.Length,
				Offset: inputs.Start,
			}),
		).
		Find(&templates).Error
	if err != nil {
		if !utility.IsWindow() {
			err = fmt.Errorf(lang.Errors.TemplateError.List.ToString())
		}
		return datatable.Response{}, err
	}
	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.MakeResponseDatatableV2(templates, inputs, user)
	return
}

type TemplateRecordDatatable struct {
	TemplateRecord
	RowId              string `json:"DT_RowId"`
	Name               string `json:"name"`
	Type               string `json:"type"`
	Size               string `json:"size"`
	Color              string `json:"color"`
	RelatedVideo       string `json:"related_video"`
	Layout             string `json:"layout"`
	Sticky             string `json:"sticky"`
	FullScreenButton   string `json:"fullscreen_button"`
	NextPrevArrows     string `json:"next_prev_arrows_button"`
	NextPrevTime       string `json:"next_prev_time"`
	EnableLogo         string `json:"enable_logo"`
	CustomLogo         string `json:"custom_logo"`
	Link               string `json:"link"`
	ClickThrough       string `json:"click_through"`
	AutoSkip           string `json:"auto_skip"`
	TimeToSkip         string `json:"time_to_skip"`
	ShowAutoSkipButton string `json:"show_auto_skip_button"`
	Action             string `json:"action"`
}

func (t *Template) MakeResponseDatatable(templates []TemplateRecord, inputs *payload.PlayerFilterPayload, user UserRecord) (records []TemplateRecordDatatable) {
	if inputs.Start == 0 && user.Id != 2 && user.Email != "k.vision@valueimpression.com" {
		templateDefault, _ := t.GetTemplateDefault(inputs)
		for _, template := range templateDefault {
			rec := TemplateRecordDatatable{
				TemplateRecord: template,
				RowId:          "template_" + strconv.FormatInt(template.Id, 10),
				Name:           htmlblock.Render("player/index/name.block.gohtml", template).String(),
				Type:           htmlblock.Render("player/index/type.block.gohtml", template).String(),
				//Size:               htmlblock.Render("player/index/size.block.gohtml", template).String(),
				//RelatedVideo:       htmlblock.Render("player/index/related_video.block.gohtml", template).String(),
				//Layout:             htmlblock.Render("player/index/layout.block.gohtml", template).String(),
				//Sticky:             htmlblock.Render("player/index/sticky.block.gohtml", template).String(),
				//FullScreenButton:   htmlblock.Render("player/index/fullscreen_button.block.gohtml", template).String(),
				//NextPrevArrows:     htmlblock.Render("player/index/next_prev_arrows_button.block.gohtml", template).String(),
				//NextPrevTime:       htmlblock.Render("player/index/next_prev_time.block.gohtml", template).String(),
				//EnableLogo:         htmlblock.Render("player/index/enable_logo.block.gohtml", template).String(),
				//CustomLogo:         htmlblock.Render("player/index/custom_logo.block.gohtml", template).String(),
				//Link:               "<a href=\"" + template.Link + "\" title=\"" + template.Link + "\">Click here</a>",
				//ClickThrough:       "<a href=\"" + template.ClickThrough + "\" title=\"" + template.ClickThrough + "\">Click here</a>",
				//AutoSkip:           htmlblock.Render("player/index/auto_skip.block.gohtml", template).String(),
				//TimeToSkip:         strconv.Itoa(template.TimeToSkip) + "s",
				//ShowAutoSkipButton: htmlblock.Render("player/index/show_auto_skip_button.block.gohtml", template).String(),
				Action: htmlblock.Render("player/index/action_default.block.gohtml", template).String(),
			}
			records = append(records, rec)
		}
	}

	for _, template := range templates {
		rec := TemplateRecordDatatable{
			TemplateRecord: template,
			RowId:          "template_" + strconv.FormatInt(template.Id, 10),
			Name:           htmlblock.Render("player/index/name.block.gohtml", template).String(),
			Type:           htmlblock.Render("player/index/type.block.gohtml", template).String(),
			//Size:               htmlblock.Render("player/index/size.block.gohtml", template).String(),
			//RelatedVideo:       htmlblock.Render("player/index/related_video.block.gohtml", template).String(),
			//Layout:             htmlblock.Render("player/index/layout.block.gohtml", template).String(),
			//Sticky:             htmlblock.Render("player/index/sticky.block.gohtml", template).String(),
			//FullScreenButton:   htmlblock.Render("player/index/fullscreen_button.block.gohtml", template).String(),
			//NextPrevArrows:     htmlblock.Render("player/index/next_prev_arrows_button.block.gohtml", template).String(),
			//NextPrevTime:       htmlblock.Render("player/index/next_prev_time.block.gohtml", template).String(),
			//EnableLogo:         htmlblock.Render("player/index/enable_logo.block.gohtml", template).String(),
			//CustomLogo:         htmlblock.Render("player/index/custom_logo.block.gohtml", template).String(),
			//Link:               "<a href=\"" + template.Link + "\" title=\"" + template.Link + "\">Click here</a>",
			//ClickThrough:       "<a href=\"" + template.ClickThrough + "\" title=\"" + template.ClickThrough + "\">Click here</a>",
			//AutoSkip:           htmlblock.Render("player/index/auto_skip.block.gohtml", template).String(),
			//TimeToSkip:         strconv.Itoa(template.TimeToSkip) + "s",
			//ShowAutoSkipButton: htmlblock.Render("player/index/show_auto_skip_button.block.gohtml", template).String(),
			Action: htmlblock.Render("player/index/action.block.gohtml", template).String(),
		}
		records = append(records, rec)
	}
	return
}

func (t *Template) MakeResponseDatatableV2(templates []TemplateRecord, inputs *payload.PlayerFilterPayload, user UserRecord) (records []TemplateRecordDatatable) {
	if inputs.Start == 0 && user.Id != 2 && user.Email != "k.vision@valueimpression.com" {
		templateDefault, _ := t.GetTemplateDefault(inputs)
		for _, template := range templateDefault {
			rec := TemplateRecordDatatable{
				TemplateRecord: template,
				RowId:          "template_" + strconv.FormatInt(template.Id, 10),
				Name:           htmlblock.Render("player-v2/index/name.block.gohtml", template).String(),
				Type:           htmlblock.Render("player-v2/index/type.block.gohtml", template).String(),
				//Size:               htmlblock.Render("player/index/size.block.gohtml", template).String(),
				//RelatedVideo:       htmlblock.Render("player/index/related_video.block.gohtml", template).String(),
				//Layout:             htmlblock.Render("player/index/layout.block.gohtml", template).String(),
				//Sticky:             htmlblock.Render("player/index/sticky.block.gohtml", template).String(),
				//FullScreenButton:   htmlblock.Render("player/index/fullscreen_button.block.gohtml", template).String(),
				//NextPrevArrows:     htmlblock.Render("player/index/next_prev_arrows_button.block.gohtml", template).String(),
				//NextPrevTime:       htmlblock.Render("player/index/next_prev_time.block.gohtml", template).String(),
				//EnableLogo:         htmlblock.Render("player/index/enable_logo.block.gohtml", template).String(),
				//CustomLogo:         htmlblock.Render("player/index/custom_logo.block.gohtml", template).String(),
				//Link:               "<a href=\"" + template.Link + "\" title=\"" + template.Link + "\">Click here</a>",
				//ClickThrough:       "<a href=\"" + template.ClickThrough + "\" title=\"" + template.ClickThrough + "\">Click here</a>",
				//AutoSkip:           htmlblock.Render("player/index/auto_skip.block.gohtml", template).String(),
				//TimeToSkip:         strconv.Itoa(template.TimeToSkip) + "s",
				//ShowAutoSkipButton: htmlblock.Render("player/index/show_auto_skip_button.block.gohtml", template).String(),
				Action: htmlblock.Render("player-v2/index/action_default.block.gohtml", template).String(),
			}
			records = append(records, rec)
		}
	}

	for _, template := range templates {
		rec := TemplateRecordDatatable{
			TemplateRecord: template,
			RowId:          "template_" + strconv.FormatInt(template.Id, 10),
			Name:           htmlblock.Render("player-v2/index/name.block.gohtml", template).String(),
			Type:           htmlblock.Render("player-v2/index/type.block.gohtml", template).String(),
			//Size:               htmlblock.Render("player/index/size.block.gohtml", template).String(),
			//RelatedVideo:       htmlblock.Render("player/index/related_video.block.gohtml", template).String(),
			//Layout:             htmlblock.Render("player/index/layout.block.gohtml", template).String(),
			//Sticky:             htmlblock.Render("player/index/sticky.block.gohtml", template).String(),
			//FullScreenButton:   htmlblock.Render("player/index/fullscreen_button.block.gohtml", template).String(),
			//NextPrevArrows:     htmlblock.Render("player/index/next_prev_arrows_button.block.gohtml", template).String(),
			//NextPrevTime:       htmlblock.Render("player/index/next_prev_time.block.gohtml", template).String(),
			//EnableLogo:         htmlblock.Render("player/index/enable_logo.block.gohtml", template).String(),
			//CustomLogo:         htmlblock.Render("player/index/custom_logo.block.gohtml", template).String(),
			//Link:               "<a href=\"" + template.Link + "\" title=\"" + template.Link + "\">Click here</a>",
			//ClickThrough:       "<a href=\"" + template.ClickThrough + "\" title=\"" + template.ClickThrough + "\">Click here</a>",
			//AutoSkip:           htmlblock.Render("player/index/auto_skip.block.gohtml", template).String(),
			//TimeToSkip:         strconv.Itoa(template.TimeToSkip) + "s",
			//ShowAutoSkipButton: htmlblock.Render("player/index/show_auto_skip_button.block.gohtml", template).String(),
			Action: htmlblock.Render("player-v2/index/action.block.gohtml", template).String(),
		}
		records = append(records, rec)
	}
	return
}

func (t *Template) setFilterSearch(inputs *payload.PlayerFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var flag bool
		// Search from form of datatable <- not use
		if inputs.Search != nil && inputs.Search.Value != "" {
			flag = true
		}
		// Search from form filter
		if inputs.PostData.QuerySearch != "" {
			flag = true
		}
		if !flag {
			return db
		}
		return db.Where("name  LIKE ?", "%"+inputs.PostData.QuerySearch+"%")
	}
}

func (t *Template) setFilterType(inputs *payload.PlayerFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.Type != nil {
			switch inputs.PostData.Type.(type) {
			case string, int:
				if inputs.PostData.Type != "" {
					return db.Where("type = ?", inputs.PostData.Type)
				}
			case []string, []interface{}:
				return db.Where("type IN ?", inputs.PostData.Type)
			}
		}
		return db
	}
}

func (t *Template) setFilterSize(inputs *payload.PlayerFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.PostData.Size != nil {
			switch inputs.PostData.Size.(type) {
			case string, int:
				if inputs.PostData.Size != "" {
					return db.Where("size = ?", inputs.PostData.Size)
				}
			case []string, []interface{}:
				return db.Where("size IN ?", inputs.PostData.Size)
			}
		}
		return db
	}
}

func (t *Template) setDefault(user UserRecord) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if user.Id != 2 && user.Email != "k.vision@valueimpression.com" {
			return db.Where("is_default != 1")
		}
		return db
	}
}

func (t *Template) setOrder(inputs *payload.PlayerFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(inputs.Order) > 0 {
			var orders []string
			for _, order := range inputs.Order {
				column := inputs.Columns[order.Column]
				orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
			}
			orderString := strings.Join(orders, ", ")
			return db.Order(orderString)
		}
		return db.Order("id DESC")
	}
}

func (t *Template) GetAll(userId int64) (records []TemplateRecord, err error) {
	err = mysql.Client.Where("is_default = 1 or user_id = ?", userId).Find(&records).Error
	return
}

func (t *Template) GetTemplateDefault(inputs *payload.PlayerFilterPayload) (records []TemplateRecord, err error) {
	err = mysql.Client.Where("is_default = ?", mysql.TypeOn).
		Scopes(
			t.setFilterSearch(inputs),
			t.setFilterType(inputs),
			t.setFilterSize(inputs),
		).
		Find(&records).Error
	return
}

func (t *Template) GetTemplateDefaultById(id int64) (record PlayerRecord, err error) {
	err = mysql.Client.Where("is_default = ? and id = ?", mysql.TypeOn, id).Find(&record).Error
	return
}

func (t *Template) GetById(id, userId int64) (row PlayerRecord, err error) {
	err = mysql.Client.Where("id = ? and user_id = ?", id, userId).Find(&row).Error
	return
}
