package model

import (
	"errors"
	"fmt"
	"source/apps/frontend/lang"
	"source/apps/frontend/payload"
	"source/apps/history"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"source/pkg/utility"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Player struct{}

type PlayerRecord struct {
	mysql.TableTemplate
}

func (PlayerRecord) TableName() string {
	return mysql.Tables.Template
}

func (t *Player) IsHave(name string, userId int64, templateId int64) bool {
	var record PlayerRecord
	// case check if create
	if templateId == 0 {
		mysql.Client.Where("user_id = ? AND name = ?", userId, name).Find(&record)
	} else {
		mysql.Client.Where("id != ? AND user_id = ? AND name = ?", templateId, userId, name).Find(&record)
	}
	if record.Id > 0 {
		return true
	}
	return false
}

func (t *Player) Create(inputs payload.TemplateCreate, user UserRecord, userAdmin UserRecord) (record PlayerRecord, errs []ajax.Error) {
	lang := lang.Translate
	errs = t.ValidateCreate(inputs)
	if len(errs) > 0 {
		return
	}
	if t.IsHave(inputs.Name, user.Id, 0) {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Name cannot be duplicated",
		})
		return
	}
	record.MakeRow(inputs, user)
	err := mysql.Client.Create(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.TemplateError.Add.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}
	if record.Type != mysql.TYPEDisplay {
		listInventoryId := t.GetListInventoryByType(record.Type, user.Id)
		for _, id := range listInventoryId {
			new(Inventory).ResetCacheWorker(id)
		}
	}

	// Push history
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = user.Id
	}
	_ = history.PushHistory(&history.Template{
		Detail:    history.DetailTemplateFE,
		CreatorId: creatorId,
		RecordOld: mysql.TableTemplate{},
		RecordNew: record.TableTemplate,
	})
	return
}

func (t *Player) Edit(inputs payload.TemplateCreate, user UserRecord, userAdmin UserRecord) (record PlayerRecord, errs []ajax.Error) {
	lang := lang.Translate
	recordOld := t.VerificationRecord(inputs.Id, user.Id)
	if recordOld.Id < 1 {
		errs = append(errs, ajax.Error{
			Id:      "id",
			Message: "You don't own this template",
		})
		return
	}
	if t.IsHave(inputs.Name, user.Id, recordOld.Id) {
		fmt.Println("help", true)
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Name cannot be duplicated",
		})
		return
	}
	// Lưu lại các giá trị cũ
	record.TableTemplate = recordOld.TableTemplate
	inputs.Type = recordOld.Type
	errs = t.ValidateCreate(inputs)
	if len(errs) > 0 {
		return
	}
	// Thêm các giá trị mới được update vào record
	record.MakeRow(inputs, user)
	// Không cho đổi type khi update template
	record.Type = recordOld.Type
	err := mysql.Client.Save(&record).Error
	if err != nil {
		if !utility.IsWindow() {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: lang.Errors.TemplateError.Edit.ToString(),
			})
		} else {
			errs = append(errs, ajax.Error{
				Id:      "",
				Message: err.Error(),
			})
		}
		return
	}
	if record.Type != mysql.TYPEDisplay {
		listInventoryId := t.GetListInventoryByType(record.Type, user.Id)
		for _, id := range listInventoryId {
			new(Inventory).ResetCacheWorker(id)
		}
	}

	// Push history
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = user.Id
	}
	_ = history.PushHistory(&history.Template{
		Detail:    history.DetailTemplateFE,
		CreatorId: creatorId,
		RecordOld: recordOld.TableTemplate,
		RecordNew: record.TableTemplate,
	})
	return
}

func (record *PlayerRecord) MakeRow(row payload.TemplateCreate, user UserRecord) {
	record.Id = row.Id
	record.UserId = user.Id
	record.Name = row.Name
	record.Type = row.Type

	// if row.Type == mysql.TYPEVideo {
	// 	temp := []mysql.TypePlayerLayout{mysql.TypePlayerLayoutBasic,
	// 		mysql.TypePlayerLayoutClassic,
	// 		mysql.TypePlayerLayoutSmall,
	// 		mysql.TypePlayerLayoutInContent,
	// 		mysql.TypePlayerLayoutSide,
	// 		mysql.TypePlayerLayoutInContentThumb,
	// 		mysql.TypePlayerLayoutInContentText,
	// 		mysql.TypePlayerLayoutTopArticle}

	// 	for _, v := range temp {
	// 		if row.PlayerLayout == v {
	// 			row.Type = mysql.TYPEInStream
	// 			break
	// 		}
	// 	}
	// }

	if record.Type == mysql.TYPEPNative {
		record.AdStyle = mysql.TYPETemplateAdStyle(row.AdStyle)
		if record.AdStyle == mysql.TYPETemplateAdStyleMultiple {
			record.Template = "grid"
		} else {
			record.Template = "standard"
		}
		record.AdSize = row.AdSize
		record.Columns = row.Columns
		record.Rows = row.Rows
		record.Mode = row.Mode
		record.ActionButtonColor = row.ActionButtonColor
		record.AdvertiserColor = row.AdvertiserColor
		record.BackgroundColor = row.BackgroundColor
		record.TitleColor = row.TitleColor
		record.SponsoredBrand = mysql.TYPESponsoredBrandFalse
		if row.IsDefault == "true" && user.Id == 2 && user.Email == "k.vision@valueimpression.com" {
			record.IsDefault = mysql.TypeOn
		} else {
			record.IsDefault = mysql.TypeOff
		}
		return
	}
	record.PlayerLayout = row.PlayerLayout
	record.Size = row.Size
	record.MaxWidth = row.MaxWidth
	record.Width = row.Width
	record.PlayMode = row.PlayMode
	record.CloseFloatingButtonDesktop = row.CloseFloatingButtonDesktop
	record.FloatOnBottom = row.FloatOnBottom
	record.FloatingOnView = row.FloatingOnView
	record.FloatingOnImpression = row.FloatingOnImpression
	record.FloatingOnAdFetched = row.FloatingOnAdFetched
	record.FloatingWidth = row.FloatingWidth
	record.FloatingPositionDesktop = row.FloatingPositionDesktop
	record.MarginTopDesktop = row.MarginTopDesktop
	record.MarginBottomDesktop = row.MarginBottomDesktop
	record.MarginRightDesktop = row.MarginRightDesktop
	record.MarginLeftDesktop = row.MarginLeftDesktop
	record.CloseFloatingButtonMobile = row.CloseFloatingButtonMobile
	record.FloatOnBottomMobile = row.FloatOnBottomMobile
	record.FloatingOnViewMobile = row.FloatingOnViewMobile
	record.FloatingOnAdFetchedMobile = row.FloatingOnAdFetchedMobile
	record.FloatingWidthMobile = row.FloatingWidthMobile
	record.FloatingPositionMobile = row.FloatingPositionMobile
	record.MarginBottomMobile = row.MarginBottomMobile
	record.MarginRightMobile = row.MarginRightMobile
	record.MarginLeftMobile = row.MarginLeftMobile
	record.ColumnsNumber = row.ColumnsNumber
	record.ColumnsPosition = row.ColumnsPosition
	record.MainTitle = row.MainTitle
	record.MainTitleText = row.MainTitleText
	record.SubTitle = row.SubTitle
	record.SubTitleText = row.SubTitleText
	record.ActionButton = row.ActionButton
	record.ActionButtonText = row.ActionButtonText
	record.TitleEnable = row.TitleEnable
	record.DescriptionEnable = row.DescriptionEnable
	record.ControlColor = row.ControlColor
	record.ThemeColor = row.ThemeColor
	record.BackgroundColor = row.BackgroundColor
	record.MainTitleBackgroundColor = row.MainTitleBackgroundColor
	record.TitleColor = row.TitleColor
	record.MainTitleColor = row.MainTitleColor
	record.TitleBackgroundColor = row.TitleBackgroundColor
	record.TitleHoverBackgroundColor = row.TitleHoverBackgroundColor
	record.ActionButtonColor = row.ActionButtonColor
	record.DescriptionColor = row.DescriptionColor
	record.DefaultSoundMode = row.DefaultSoundMode
	record.FullscreenButton = row.FullscreenButton
	record.NextPrevArrowsButton = row.NextPrevArrowsButton
	record.NextPrevTime = row.NextPrevTime
	record.VideoConfig = row.VideoConfig
	record.ShowStats = row.ShowStats
	record.ShareButton = row.ShareButton
	record.CustomLogo = row.CustomLogo
	record.Link = row.Link
	record.ClickThrough = row.ClickThrough
	record.WaitForAd = row.WaitForAd
	record.VastRetry = row.VastRetry
	record.Delay = row.Delay
	record.AutoSkip = row.AutoSkip
	record.TimeToSkip = row.TimeToSkip
	record.ShowAutoSkipButton = row.ShowAutoSkipButton
	record.NumberOfPreRollAds = row.NumberOfPreRollAds
	record.FloatingOnDesktop = row.FloatingOnDesktop
	record.FloatingOnMobile = row.FloatingOnMobile
	record.PoweredBy = row.PoweredBy
	record.EnableLogo = row.EnableLogo
	record.AdvertisementScenario = row.AdvertisementScenario
	record.PreRoll = row.PreRoll
	record.MidRoll = row.MidRoll
	record.PostRoll = row.PostRoll
	record.AutoStart = row.AutoStart
	record.ShowControls = row.ShowControls
	record.PubPowerLogo = row.PubPowerLogo
	//if row.Type == mysql.TYPETopArticles {
	//	record.MainTitleText = row.MainTitleTopArticle
	//	record.EnableLogo = row.EnableLogoTopArticle
	//	record.CustomLogo = row.CustomLogoTopArticle
	//	record.Link = row.LinkTopArticle
	//	record.ClickThrough = row.ClickThroughTopArticle
	//	record.PoweredBy = row.PoweredByTopArticle
	//}
	if row.IsDefault == "true" && user.Id == 2 && user.Email == "k.vision@valueimpression.com" {
		record.IsDefault = mysql.TypeOn
	} else {
		record.IsDefault = mysql.TypeOff
	}
	return
}

func (t *Player) ValidateCreate(inputs payload.TemplateCreate) (errs []ajax.Error) {
	if utility.ValidateString(inputs.Name) == "" {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Name is required",
		})
	}
	switch inputs.Type {
	case mysql.TYPEInStream:
		err := validateTypeInStream(inputs)
		errs = append(errs, err...)
		break
	case mysql.TYPEOutStream:
		err := validateTypeOutStream(inputs)
		errs = append(errs, err...)
		break
	case mysql.TYPETopArticles:
		err := validateTypeTopArticles(inputs)
		errs = append(errs, err...)
	case mysql.TYPEPNative:
		err := validateTypeNative(inputs)
		errs = append(errs, err...)
	case mysql.TYPEVideo:
		err := validateTypeVideo(inputs)
		errs = append(errs, err...)
		break
	}
	return
}

func (t *Player) GetById(id, userId int64) (record PlayerRecord, err error) {
	err = mysql.Client.Where("id = ? and user_id = ?", id, userId).Find(&record).Error
	if record.Id == 0 {
		err = errors.New("record not found")
		return
	}
	return
}

func (t *Player) GetDetail(id int64) (record PlayerRecord) {
	mysql.Client.Where("id = ?", id).Find(&record)
	return
}

func (t *Player) Delete(id, userId int64, userAdmin UserRecord, lang lang.Translation) fiber.Map {
	record := t.GetDetail(id)
	err := mysql.Client.Model(&PlayerRecord{}).Delete(&PlayerRecord{}, "id = ? and user_id = ?", id, userId).Error
	if err != nil {
		if !utility.IsWindow() {
			return fiber.Map{
				"status":  "err",
				"message": lang.Errors.TemplateError.Delete.ToString(),
				"id":      id,
			}
		}
		return fiber.Map{
			"status":  "err",
			"message": err.Error(),
			"id":      id,
		}
	} else {
		if record.Type != mysql.TYPEDisplay {
			listInventoryId := t.GetListInventoryByType(record.Type, userId)
			for _, id := range listInventoryId {
				new(Inventory).ResetCacheWorker(id)
			}
		}

		// History
		var creatorId int64
		if userAdmin.Id != 0 {
			creatorId = userAdmin.Id
		} else {
			creatorId = userId
		}
		fmt.Println(record)
		_ = history.PushHistory(&history.Template{
			Detail:    history.DetailTemplateFE,
			CreatorId: creatorId,
			RecordOld: record.TableTemplate,
			RecordNew: mysql.TableTemplate{},
		})

		return fiber.Map{
			"status":  "success",
			"message": "done",
			"id":      id,
		}
	}
}

func (t *Player) Duplicate(id int64, user UserRecord) fiber.Map {
	var record PlayerRecord
	err := mysql.Client.Where("id = ?", id).Find(&record).Error
	// if record.UserId != 2 {
	// 	return fiber.Map{
	// 		"status":  "err",
	// 		"message": err.Error(),
	// 	}
	// }

	// set name => name + " - copy 1"
	var records []PlayerRecord
	array := strings.Split(record.Name, " - copy ")
	err = mysql.Client.Where("name like ?", "%"+array[0]+" - copy%").Find(&records).Error
	if len(records) == 0 {
		record.Name = record.Name + " - copy 1"
	} else {
		number := 1
		for _, value := range records {
			name := value.Name
			Number := strings.Replace(name, array[0]+" - copy ", "", -1)
			intVar, _ := strconv.Atoi(Number)
			if intVar >= number {
				number = intVar + 1
			}
		}
		record.Name = array[0] + " - copy " + strconv.Itoa(number)
	}

	record.Id = 0
	record.UserId = user.Id
	record.IsDefault = mysql.TypeOff
	err = mysql.Client.Create(&record).Error
	if err != nil {
		return fiber.Map{
			"status":  "err",
			"message": err.Error(),
		}
	} else {
		return fiber.Map{
			"status":  "success",
			"message": "done",
		}
	}
}

func validateTypeVideo(inputs payload.TemplateCreate) (errs []ajax.Error) {
	if inputs.PlayerLayout == mysql.TypePlayerLayoutBasic ||
		inputs.PlayerLayout == mysql.TypePlayerLayoutClassic ||
		inputs.PlayerLayout == mysql.TypePlayerLayoutSmall ||
		inputs.PlayerLayout == mysql.TypePlayerLayoutInContent ||
		inputs.PlayerLayout == mysql.TypePlayerLayoutSide ||
		inputs.PlayerLayout == mysql.TypePlayerLayoutInContentThumb ||
		inputs.PlayerLayout == mysql.TypePlayerLayoutInContentText ||
		inputs.PlayerLayout == mysql.TypePlayerLayoutTopArticle {
		errs = append(errs, validateTypeInStream(inputs)...)
	} else if inputs.PlayerLayout == mysql.TypePlayerLayoutNone {
		errs = append(errs, validateTypeOutStream(inputs)...)
	} else {
		errs = append(errs, ajax.Error{
			Id:      "player_layout",
			Message: "Player layout incorrect",
		})
	}
	return errs
}

func validateTypeInStream(inputs payload.TemplateCreate) (errs []ajax.Error) {
	if inputs.PlayerLayout == 0 {
		errs = append(errs, ajax.Error{
			Id:      "player_layout",
			Message: "Player Layout is required",
		})
	}

	if inputs.Size == mysql.TypeSizeResponsive {
		// if inputs.MaxWidth == 0 {
		// 	errs = append(errs, ajax.Error{
		// 		Id:      "max_width",
		// 		Message: "Max Width is required",
		// 	})
		// }
	} else if inputs.Size == mysql.TypeSizeFixed {
		if inputs.Width == 0 {
			errs = append(errs, ajax.Error{
				Id:      "width",
				Message: "Width is required",
			})
		}
	}

	// if inputs.PlayerLayout == 4 {
	//	if inputs.MaxWidth < 500 {
	//		errs = append(errs, ajax.Error{
	//			Id:      "max_width",
	//			Message: "Max Width must be more than 500px",
	//		})
	//	}
	// }

	if inputs.FloatingOnDesktop == mysql.TypeOn {
		if inputs.FloatingWidth == 0 {
			errs = append(errs, ajax.Error{
				Id:      "floating_width",
				Message: "Floating Width is required",
			})
		}

		if inputs.FloatingWidth != 0 && inputs.FloatingWidth < 256 {
			errs = append(errs, ajax.Error{
				Id:      "floating_width",
				Message: "min 256px",
			})
		}
		// switch inputs.FloatingPositionDesktop {
		// case mysql.TypePositionDesktopBottomRight:
		//	if inputs.MarginBottomDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_bottom_desktop",
		//			Message: "Margin Bottom Desktop is required",
		//		})
		//	}
		//	if inputs.MarginRightDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_right_desktop",
		//			Message: "Margin Right Desktop is required",
		//		})
		//	}
		// case mysql.TypePositionDesktopBottomLeft:
		//	if inputs.MarginBottomDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_bottom_desktop",
		//			Message: "Margin Bottom Desktop is required",
		//		})
		//	}
		//	if inputs.MarginLeftDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_left_desktop",
		//			Message: "Margin Left Desktop is required",
		//		})
		//	}
		// case mysql.TypePositionDesktopTopRight:
		//	if inputs.MarginTopDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_top_desktop",
		//			Message: "Margin Top Desktop is required",
		//		})
		//	}
		//	if inputs.MarginRightDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_right_desktop",
		//			Message: "Margin Right Desktop is required",
		//		})
		//	}
		// case mysql.TypePositionDesktopTopLeft:
		//	if inputs.MarginTopDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_top_desktop",
		//			Message: "Margin Top Desktop is required",
		//		})
		//	}
		//	if inputs.MarginLeftDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_left_desktop",
		//			Message: "Margin Left Desktop is required",
		//		})
		//	}
		// }
	}
	if inputs.FloatingOnMobile == mysql.TypeOn {
		if inputs.FloatingWidthMobile == 0 {
			errs = append(errs, ajax.Error{
				Id:      "floating_width_mobile",
				Message: "Floating Width Mobile is required",
			})
		}

		if inputs.FloatingWidthMobile != 0 && inputs.FloatingWidthMobile < 256 {
			errs = append(errs, ajax.Error{
				Id:      "floating_width_mobile",
				Message: "min 256px",
			})
		}

		//switch inputs.FloatingPositionMobile {
		//case mysql.TypePositionMobileBottomRight:
		//	if inputs.MarginBottomMobile == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_bottom_mobile",
		//			Message: "Margin Bottom Mobile is required",
		//		})
		//	}
		//	if inputs.MarginRightMobile == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_right_mobile",
		//			Message: "Margin Left Mobile is required",
		//		})
		//	}
		//case mysql.TypePositionMobileBottomLeft:
		//	if inputs.MarginBottomMobile == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_bottom_mobile",
		//			Message: "Margin Bottom Mobile is required",
		//		})
		//	}
		//	if inputs.MarginLeftMobile == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_left_mobile",
		//			Message: "Margin Left Mobile is required",
		//		})
		//	}
		//}
	}

	if inputs.MainTitle == 1 {
		if utility.ValidateString(inputs.MainTitleText) == "" {
			errs = append(errs, ajax.Error{
				Id:      "main_title_text",
				Message: "Main Title Text is required",
			})
		}
	}

	if utility.ValidateString(inputs.ControlColor) == "" {
		errs = append(errs, ajax.Error{
			Id:      "control_color",
			Message: "Control Color is required",
		})
	}

	if utility.ValidateString(inputs.BackgroundColor) == "" {
		errs = append(errs, ajax.Error{
			Id:      "background_color",
			Message: "Background Color is required",
		})
	}

	if inputs.PlayerLayout != 1 && inputs.PlayerLayout != 2 && inputs.PlayerLayout != 6 {
		if utility.ValidateString(inputs.DescriptionColor) == "" {
			errs = append(errs, ajax.Error{
				Id:      "description_color",
				Message: "Description Color is required",
			})
		}

		if utility.ValidateString(inputs.TitleColor) == "" {
			errs = append(errs, ajax.Error{
				Id:      "title_color",
				Message: "Title Color is required",
			})
		}
	}

	if inputs.PlayerLayout == 6 {
		if inputs.ColumnsNumber < 1 || inputs.ColumnsNumber > 3 {
			errs = append(errs, ajax.Error{
				Id:      "columns_number",
				Message: "Min = 1, Max = 3",
			})
		}
	}

	if inputs.CustomLogo == mysql.TypeOn {
		// if utility.ValidateString(inputs.Link) == "" {
		// 	errs = append(errs, ajax.Error{
		// 		Id:      "link",
		// 		Message: "Link is required",
		// 	})
		// }
		// if utility.ValidateString(inputs.ClickThrough) == "" {
		// 	errs = append(errs, ajax.Error{
		// 		Id:      "click_through",
		// 		Message: "Click Through required",
		// 	})
		// }
	}

	if inputs.VastRetry == 0 {
		errs = append(errs, ajax.Error{
			Id:      "vast_retry",
			Message: "Vast Retry required",
		})
	}

	if inputs.AutoSkip == mysql.TypeOn {
		if inputs.TimeToSkip == 0 {
			errs = append(errs, ajax.Error{
				Id:      "time_to_skip",
				Message: "Time To Skip required",
			})
		}
		if inputs.ShowAutoSkipButton == 0 {
			errs = append(errs, ajax.Error{
				Id:      "show_auto_skip_button",
				Message: "Show Auto Skip Button required",
			})
		}
		if inputs.NumberOfPreRollAds < 1 || inputs.NumberOfPreRollAds > 2 {
			errs = append(errs, ajax.Error{
				Id:      "number_of_pre_roll_ads",
				Message: "Min = 1, Max = 2",
			})
		}
	}
	return
}

func validateTypeOutStream(inputs payload.TemplateCreate) (errs []ajax.Error) {
	if inputs.Size == mysql.TypeSizeResponsive {
		// if inputs.MaxWidth == 0 {
		// 	errs = append(errs, ajax.Error{
		// 		Id:      "max_width",
		// 		Message: "Max Width is required",
		// 	})
		// }
	} else if inputs.Size == mysql.TypeSizeFixed {
		if inputs.Width == 0 {
			errs = append(errs, ajax.Error{
				Id:      "width",
				Message: "Width is required",
			})
		}
	}

	if inputs.FloatingOnDesktop == mysql.TypeOn {
		if inputs.FloatingWidth == 0 {
			errs = append(errs, ajax.Error{
				Id:      "floating_width",
				Message: "Floating Width is required",
			})
		}

		if inputs.FloatingWidth != 0 && inputs.FloatingWidth < 256 {
			errs = append(errs, ajax.Error{
				Id:      "floating_width",
				Message: "min 256px",
			})
		}
		// switch inputs.FloatingPositionDesktop {
		// case mysql.TypePositionDesktopBottomRight:
		//	if inputs.MarginBottomDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_bottom_desktop",
		//			Message: "Margin Bottom Desktop is required",
		//		})
		//	}
		//	if inputs.MarginRightDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_right_desktop",
		//			Message: "Margin Right Desktop is required",
		//		})
		//	}
		// case mysql.TypePositionDesktopBottomLeft:
		//	if inputs.MarginBottomDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_bottom_desktop",
		//			Message: "Margin Bottom Desktop is required",
		//		})
		//	}
		//	if inputs.MarginLeftDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_left_desktop",
		//			Message: "Margin Left Desktop is required",
		//		})
		//	}
		// case mysql.TypePositionDesktopTopRight:
		//	if inputs.MarginTopDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_top_desktop",
		//			Message: "Margin Top Desktop is required",
		//		})
		//	}
		//	if inputs.MarginRightDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_right_desktop",
		//			Message: "Margin Right Desktop is required",
		//		})
		//	}
		// case mysql.TypePositionDesktopTopLeft:
		//	if inputs.MarginTopDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_top_desktop",
		//			Message: "Margin Top Desktop is required",
		//		})
		//	}
		//	if inputs.MarginLeftDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_left_desktop",
		//			Message: "Margin Left Desktop is required",
		//		})
		//	}
		// }
	}

	if inputs.FloatingOnMobile == mysql.TypeOn {
		if inputs.FloatingWidthMobile == 0 {
			errs = append(errs, ajax.Error{
				Id:      "floating_width_mobile",
				Message: "Floating Width Mobile is required",
			})
		}

		if inputs.FloatingWidthMobile != 0 && inputs.FloatingWidthMobile < 256 {
			errs = append(errs, ajax.Error{
				Id:      "floating_width_mobile",
				Message: "min 256px",
			})
		}
		// switch inputs.FloatingPositionMobile {
		// case mysql.TypePositionMobileBottomRight:
		//	if inputs.MarginBottomMobile == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_bottom_mobile",
		//			Message: "Margin Bottom Mobile is required",
		//		})
		//	}
		//	if inputs.MarginRightMobile == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_right_mobile",
		//			Message: "Margin Left Mobile is required",
		//		})
		//	}
		// case mysql.TypePositionMobileBottomLeft:
		//	if inputs.MarginBottomMobile == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_bottom_mobile",
		//			Message: "Margin Bottom Mobile is required",
		//		})
		//	}
		//	if inputs.MarginLeftMobile == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_left_mobile",
		//			Message: "Margin Left Mobile is required",
		//		})
		//	}
		// }
	}

	if inputs.VastRetry == 0 {
		errs = append(errs, ajax.Error{
			Id:      "vast_retry",
			Message: "Vast Retry required",
		})
	}
	return
}

func validateTypeTopArticles(inputs payload.TemplateCreate) (errs []ajax.Error) {
	if inputs.Size == mysql.TypeSizeResponsive {
		// if inputs.MaxWidth == 0 {
		// 	errs = append(errs, ajax.Error{
		// 		Id:      "max_width",
		// 		Message: "Max Width is required",
		// 	})
		// }
	} else if inputs.Size == mysql.TypeSizeFixed {
		if inputs.Width == 0 {
			errs = append(errs, ajax.Error{
				Id:      "width",
				Message: "Width is required",
			})
		}
	}

	if inputs.FloatingOnDesktop == mysql.TypeOn {
		if inputs.FloatingWidth == 0 {
			errs = append(errs, ajax.Error{
				Id:      "floating_width",
				Message: "Floating Width is required",
			})
		}

		if inputs.FloatingWidth != 0 && inputs.FloatingWidth < 256 {
			errs = append(errs, ajax.Error{
				Id:      "floating_width",
				Message: "min 256px",
			})
		}

		// switch inputs.FloatingPositionDesktop {
		// case mysql.TypePositionDesktopBottomRight:
		//	if inputs.MarginBottomDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_bottom_desktop",
		//			Message: "Margin Bottom Desktop is required",
		//		})
		//	}
		//	if inputs.MarginRightDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_right_desktop",
		//			Message: "Margin Right Desktop is required",
		//		})
		//	}
		// case mysql.TypePositionDesktopBottomLeft:
		//	if inputs.MarginBottomDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_bottom_desktop",
		//			Message: "Margin Bottom Desktop is required",
		//		})
		//	}
		//	if inputs.MarginLeftDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_left_desktop",
		//			Message: "Margin Left Desktop is required",
		//		})
		//	}
		// case mysql.TypePositionDesktopTopRight:
		//	if inputs.MarginTopDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_top_desktop",
		//			Message: "Margin Top Desktop is required",
		//		})
		//	}
		//	if inputs.MarginRightDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_right_desktop",
		//			Message: "Margin Right Desktop is required",
		//		})
		//	}
		// case mysql.TypePositionDesktopTopLeft:
		//	if inputs.MarginTopDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_top_desktop",
		//			Message: "Margin Top Desktop is required",
		//		})
		//	}
		//	if inputs.MarginLeftDesktop == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_left_desktop",
		//			Message: "Margin Left Desktop is required",
		//		})
		//	}
		// }
	}
	if inputs.FloatingOnMobile == mysql.TypeOn {
		if inputs.FloatingWidthMobile == 0 {
			errs = append(errs, ajax.Error{
				Id:      "floating_width_mobile",
				Message: "Floating Width Mobile is required",
			})
		}
		if inputs.FloatingWidthMobile != 0 && inputs.FloatingWidthMobile < 256 {
			errs = append(errs, ajax.Error{
				Id:      "floating_width_mobile",
				Message: "min 256px",
			})
		}

		// switch inputs.FloatingPositionMobile {
		// case mysql.TypePositionMobileBottomRight:
		//	if inputs.MarginBottomMobile == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_bottom_mobile",
		//			Message: "Margin Bottom Mobile is required",
		//		})
		//	}
		//	if inputs.MarginRightMobile == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_right_mobile",
		//			Message: "Margin Left Mobile is required",
		//		})
		//	}
		// case mysql.TypePositionMobileBottomLeft:
		//	if inputs.MarginBottomMobile == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_bottom_mobile",
		//			Message: "Margin Bottom Mobile is required",
		//		})
		//	}
		//	if inputs.MarginLeftMobile == 0 {
		//		errs = append(errs, ajax.Error{
		//			Id:      "margin_left_mobile",
		//			Message: "Margin Left Mobile is required",
		//		})
		//	}
		// }
	}

	if utility.ValidateString(inputs.ThemeColor) == "" {
		errs = append(errs, ajax.Error{
			Id:      "theme_color",
			Message: "Theme Color is required",
		})
	}

	if utility.ValidateString(inputs.BackgroundColor) == "" {
		errs = append(errs, ajax.Error{
			Id:      "background_color",
			Message: "Background Color is required",
		})
	}

	if utility.ValidateString(inputs.MainTitleBackgroundColor) == "" {
		errs = append(errs, ajax.Error{
			Id:      "main_title_background_color",
			Message: "Main Title Background Color is required",
		})
	}

	if utility.ValidateString(inputs.MainTitleColor) == "" {
		errs = append(errs, ajax.Error{
			Id:      "main_title_color",
			Message: "Main Title Color is required",
		})
	}

	if utility.ValidateString(inputs.TitleColor) == "" {
		errs = append(errs, ajax.Error{
			Id:      "title_color",
			Message: "Title Color is required",
		})
	}

	if utility.ValidateString(inputs.TitleBackgroundColor) == "" {
		errs = append(errs, ajax.Error{
			Id:      "title_background_color",
			Message: "Title Background Color is required",
		})
	}

	if utility.ValidateString(inputs.TitleHoverBackgroundColor) == "" {
		errs = append(errs, ajax.Error{
			Id:      "title_hover_background_color",
			Message: "Title Hover Background Color is required",
		})
	}

	if utility.ValidateString(inputs.ActionButtonColor) == "" {
		errs = append(errs, ajax.Error{
			Id:      "action_button_color",
			Message: "Action Button Color is required",
		})
	}

	// if utility.ValidateString(inputs.MainTitleTopArticle) == "" {
	// 	errs = append(errs, ajax.Error{
	// 		Id:      "main_title_top_article",
	// 		Message: "Main Title is required",
	// 	})
	// }

	if inputs.CustomLogoTopArticle == mysql.TypeOn {
		// if utility.ValidateString(inputs.LinkTopArticle) == "" {
		// 	errs = append(errs, ajax.Error{
		// 		Id:      "link_top_article",
		// 		Message: "Link is required",
		// 	})
		// }
		// if utility.ValidateString(inputs.ClickThroughTopArticle) == "" {
		// 	errs = append(errs, ajax.Error{
		// 		Id:      "click_through_top_article",
		// 		Message: "Click Through required",
		// 	})
		// }
	}

	if inputs.VastRetry == 0 {
		errs = append(errs, ajax.Error{
			Id:      "vast_retry",
			Message: "Vast Retry required",
		})
	}
	return
}

func validateTypeNative(inputs payload.TemplateCreate) (errs []ajax.Error) {
	if utility.ValidateString(inputs.Name) == "" {
		errs = append(errs, ajax.Error{
			Id:      "name",
			Message: "Name is required",
		})
	}
	if utility.ValidateString(inputs.AdStyle) == "" {
		errs = append(errs, ajax.Error{
			Id:      "ad_style",
			Message: "Ad Style is required",
		})
	}
	if utility.ValidateString(inputs.Mode) == "" {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: "mode is required",
		})
	}
	if inputs.AdSize == 0 {
		errs = append(errs, ajax.Error{
			Id:      "ad_size",
			Message: "Ad size is required",
		})
	}
	if utility.ValidateString(inputs.BackgroundColor) == "" {
		errs = append(errs, ajax.Error{
			Id:      "native_background",
			Message: "Background is required",
		})
	}
	if utility.ValidateString(inputs.TitleColor) == "" {
		errs = append(errs, ajax.Error{
			Id:      "native_title_color",
			Message: "Title is required",
		})
	}
	if utility.ValidateString(inputs.AdvertiserColor) == "" {
		errs = append(errs, ajax.Error{
			Id:      "native_advertiser_name",
			Message: "Advertiser name is required",
		})
	}
	if utility.ValidateString(inputs.ActionButtonColor) == "" {
		errs = append(errs, ajax.Error{
			Id:      "native_cta_button",
			Message: "CTA Button is required",
		})
	}
	return
}

func (t *Player) VerificationRecord(id, userId int64) (record TemplateRecord) {
	mysql.Client.Model(&TemplateRecord{}).Where("id = ? and user_id = ?", id, userId).Find(&record)
	return
}

func (t *Player) MakeContents(typ mysql.TYPEAdType) (contents []payload.Content) {
	if typ == mysql.TYPETopArticles {
		contents = append(contents, payload.Content{
			CreateTime: "",
			DeletedAt:  "",
			Des:        "",
			ID:         "",
			IsDefault:  "",
			Link:       "",
			Thumb:      "",
			Title:      "",
			UserID:     "",
			VideoURL: payload.VideoURL{
				M3U8: "https://ul.pubpowerplatform.io/video/16518330826278/output_69bceb9570ed563f029799ef9ef98c35.m3u8",
				Mp4:  "https://ms.pubpowerplatform.io/v2/videoplayback?ytid=GfFApfdMWGs",
				Ogg:  "",
			},
		})
	} else {
		var playlist PlaylistRecord
		mysql.Client.Find(&playlist, 13)
		videos, _ := playlist.GetVideoContentByPlaylist()
		domainUpload := "https://ul.pubpowerplatform.io"
		for _, content := range videos {
			contents = append(contents, payload.Content{
				CreateTime: content.CreatedAt.String(),
				DeletedAt:  content.DeletedAt.Time.String(),
				Des:        content.ContentDesc,
				ID:         strconv.FormatInt(content.Id, 10),
				IsDefault:  "0",
				Link:       domainUpload + content.VideoUrl,
				Thumb:      domainUpload + content.Thumb,
				Title:      content.Title,
				UserID:     strconv.FormatInt(content.UserId, 10),
				VideoURL: payload.VideoURL{
					M3U8: domainUpload + content.VideoUrl,
					Mp4:  "",
					Ogg:  "",
				},
			})
		}
	}
	return
}

func (t *Player) MakeTemplate(templateRecord PlayerRecord) (template payload.Template) {
	var maxWidth int
	switch templateRecord.Size {
	case mysql.TypeSizeFixed:
		maxWidth = templateRecord.Width
		break
	case mysql.TypeSizeResponsive:
		maxWidth = templateRecord.MaxWidth
		break
	}

	var marginTopBotDesktop int
	var marginLeftRightDesktop int
	switch templateRecord.FloatingPositionDesktop {
	case mysql.TypePositionDesktopBottomRight:
		marginTopBotDesktop = templateRecord.MarginBottomDesktop
		marginLeftRightDesktop = templateRecord.MarginRightDesktop
	case mysql.TypePositionDesktopBottomLeft:
		marginTopBotDesktop = templateRecord.MarginBottomDesktop
		marginLeftRightDesktop = templateRecord.MarginLeftDesktop
	case mysql.TypePositionDesktopTopRight:
		marginTopBotDesktop = templateRecord.MarginTopDesktop
		marginLeftRightDesktop = templateRecord.MarginRightDesktop
	case mysql.TypePositionDesktopTopLeft:
		marginTopBotDesktop = templateRecord.MarginTopDesktop
		marginLeftRightDesktop = templateRecord.MarginLeftDesktop
	}

	//Config mobile
	var marginBotMobile int
	var marginLeftRightMobile int
	switch templateRecord.FloatingPositionMobile {
	case mysql.TypePositionMobileBottomLeft:
		marginBotMobile = templateRecord.MarginBottomMobile
		marginLeftRightMobile = templateRecord.MarginLeftMobile
	case mysql.TypePositionMobileBottomRight:
		marginBotMobile = templateRecord.MarginBottomMobile
		marginLeftRightMobile = templateRecord.MarginRightMobile
	}

	//Default Column Setting
	columnSetting := payload.ColumnSetting{}

	if templateRecord.Type == mysql.TYPEVideo {
		temp := []mysql.TypePlayerLayout{mysql.TypePlayerLayoutBasic,
			mysql.TypePlayerLayoutClassic,
			mysql.TypePlayerLayoutSmall,
			mysql.TypePlayerLayoutInContent,
			mysql.TypePlayerLayoutSide,
			mysql.TypePlayerLayoutInContentThumb,
			mysql.TypePlayerLayoutInContentText,
			mysql.TypePlayerLayoutTopArticle}

		for _, v := range temp {
			if templateRecord.PlayerLayout == v {
				templateRecord.Type = mysql.TYPEInStream
				break
			}
		}
	}

	switch true {
	case templateRecord.Type == mysql.TYPEInStream:
		if templateRecord.PlayerLayout == mysql.TypePlayerLayoutInContentThumb {
			if templateRecord.ColumnsPosition == 2 {
				columnSetting.ColumnPosition = "left"
			} else {
				columnSetting.ColumnPosition = "right"
			}
			columnSetting.ColumnNumber = templateRecord.ColumnsNumber
		} else if templateRecord.PlayerLayout == mysql.TypePlayerLayoutInContentText {
			if templateRecord.ColumnsPosition == 2 {
				columnSetting.ColumnPosition = "left"
			} else {
				columnSetting.ColumnPosition = "right"
			}
		}
		template = payload.Template{
			VideoTempName: templateRecord.Name,
			AdType:        templateRecord.Type.String(),
			Appearance: &payload.Appearance{
				PlayerLayout: &payload.PlayerLayout{
					Type:      templateRecord.PlayerLayout,
					ListVideo: nil,
				},
				PlayerSize:    strings.ToLower(templateRecord.Size.String()),
				MaxWidth:      maxWidth,
				ColumnSetting: columnSetting,
			},
			Text: &payload.Text{
				MainTitle:     templateRecord.MainTitleText,
				TitleOn:       templateRecord.TitleEnable.Boolean(),
				DescriptionOn: templateRecord.DescriptionEnable.Boolean(),
			},
			Color: &payload.Color{
				Controls:    templateRecord.ControlColor,
				Background:  templateRecord.BackgroundColor,
				Title:       templateRecord.TitleColor,
				Description: templateRecord.DescriptionColor,
			},
			Controls: &payload.Controls{
				DefaultSoundMode: templateRecord.DefaultSoundMode.Boolean(),
				Fullscreen:       templateRecord.FullscreenButton.Boolean(),
				NextPrevArrow:    templateRecord.NextPrevArrowsButton.Boolean(),
				NextPrevSkip:     templateRecord.NextPrevTime.Boolean(),
				VideoConfig:      templateRecord.VideoConfig.Boolean(),
				ViewsLikes:       templateRecord.ShowStats.Boolean(),
				Share:            templateRecord.ShareButton.Boolean(),
			},
			LogoBand: &payload.LogoBand{
				EndableLogo:      templateRecord.EnableLogo.Boolean(),
				PoweredByApacdex: templateRecord.PoweredBy.Boolean(),
				LogoByPubPower:   templateRecord.PubPowerLogo.Boolean(),
				CustomLogo: &payload.CustomLogo{
					Link:         templateRecord.Link,
					ClickThrough: templateRecord.ClickThrough,
				},
			},
			AdConfig: &payload.AdConfig{
				VastRetry: templateRecord.VastRetry,
			},
		}
		template.AutoStart = templateRecord.AutoStart.Code()

		if templateRecord.AutoSkip == mysql.TypeOn {
			template.AdConfig.AutoSkip = &payload.AutoSkip{
				TimeToSkip:  templateRecord.TimeToSkip,
				AutoSkipBtn: templateRecord.ShowAutoSkipButton,
				AdsNums:     templateRecord.NumberOfPreRollAds,
			}
		} else {
			template.AdConfig.AutoSkip = nil
			template.AdConfig.Delay = templateRecord.Delay
		}

	// case isOutstream:
	// 	template = payload.Template{
	// 		VideoTempName: templateRecord.Name,
	// 		AdType:        templateRecord.Type.String(),
	// 		Appearance: &payload.Appearance{
	// 			PlayerLayout: nil,
	// 			PlayerSize:   strings.ToLower(templateRecord.Size.String()),
	// 			MaxWidth:     maxWidth,
	// 		},
	// 		Text:     nil,
	// 		Color:    nil,
	// 		Controls: nil,
	// 		LogoBand: nil,
	// 		AdConfig: &payload.AdConfig{
	// 			VastRetry: templateRecord.VastRetry,
	// 			AutoSkip:  nil,
	// 			Delay:     templateRecord.Delay,
	// 		},
	// 	}
	// case templateRecord.Type == mysql.TYPETopArticles:
	// 	template = payload.Template{
	// 		VideoTempName: templateRecord.Name,
	// 		AdType:        templateRecord.Type.String(),
	// 		Appearance: &payload.Appearance{
	// 			PlayerSize: strings.ToLower(templateRecord.Size.String()),
	// 			MaxWidth:   maxWidth,
	// 		},
	// 		Text: &payload.Text{
	// 			MainTitle:     templateRecord.MainTitleText,
	// 			TitleOn:       templateRecord.MainTitle.Boolean(),
	// 			DescriptionOn: templateRecord.DescriptionEnable.Boolean(),
	// 			ActionButton:  templateRecord.ActionButton.Boolean(),
	// 			ReadMore:      templateRecord.ActionButtonText,
	// 		},
	// 		Color: &payload.Color{
	// 			Theme:                    templateRecord.ThemeColor,
	// 			Background:               templateRecord.BackgroundColor,
	// 			Title:                    templateRecord.TitleColor,
	// 			TitleBackground:          templateRecord.TitleBackgroundColor,
	// 			MainTitle:                templateRecord.MainTitleColor,
	// 			MainTitleBackground:      templateRecord.MainTitleBackgroundColor,
	// 			TitleTextBackgroundHover: templateRecord.TitleHoverBackgroundColor,
	// 			ActionButtonColor:        templateRecord.ActionButtonColor,
	// 		},
	// 		Controls: &payload.Controls{
	// 			DefaultSoundMode: templateRecord.DefaultSoundMode.Boolean(),
	// 			Fullscreen:       templateRecord.FullscreenButton.Boolean(),
	// 			NextPrevArrow:    templateRecord.NextPrevArrowsButton.Boolean(),
	// 			NextPrevSkip:     templateRecord.NextPrevTime.Boolean(),
	// 			VideoConfig:      templateRecord.VideoConfig.Boolean(),
	// 			ViewsLikes:       templateRecord.ShowStats.Boolean(),
	// 			Share:            templateRecord.ShareButton.Boolean(),
	// 		},
	// 		LogoBand: &payload.LogoBand{
	// 			EndableLogo:      templateRecord.EnableLogo.Boolean(),
	// 			PoweredByApacdex: templateRecord.SubTitle.Boolean(),
	// 			PoweredText:      templateRecord.SubTitleText,
	// 			CustomLogo: &payload.CustomLogo{
	// 				Link:         templateRecord.Link,
	// 				ClickThrough: templateRecord.ClickThrough,
	// 			},
	// 		},
	// 		AdConfig: &payload.AdConfig{
	// 			VastRetry: templateRecord.VastRetry,
	// 			AutoSkip:  nil,
	// 			Delay:     templateRecord.Delay,
	// 		},
	// 	}
	// 	template.AutoStart = templateRecord.AutoStart.Code()
	default:
	}

	// Xử lý playMode
	if template.Appearance != nil {
		if templateRecord.FloatingOnDesktop == mysql.TypeOn {
			template.Appearance.FloatingSetting = &payload.FloatingSetting{
				CloseFloatingBtn: templateRecord.CloseFloatingButtonDesktop.Boolean(),
				FloatOnBottom:    templateRecord.FloatOnBottom.Boolean(),
				FloatingOnView:   templateRecord.FloatingOnView.Boolean(),
				Width:            templateRecord.FloatingWidth,
				Position:         templateRecord.FloatingPositionDesktop.Int(),
				MarginTopBot:     marginTopBotDesktop,
				MarginLeftRight:  marginLeftRightDesktop,
			}
		} else {
			template.Appearance.FloatingSetting = nil
		}
	}

	if templateRecord.FloatingOnMobile == mysql.TypeOn {
		template.MobileConfig = &payload.MobileConfig{
			CloseFloatingBtn: templateRecord.CloseFloatingButtonMobile.Boolean(),
			FloatOnBottom:    templateRecord.FloatOnBottomMobile.Boolean(),
			FloatingOnView:   templateRecord.FloatingOnViewMobile.Boolean(),
			Width:            templateRecord.FloatingWidthMobile,
			Position:         templateRecord.FloatingPositionMobile.Int(),
			MarginBot:        marginBotMobile,
			MarginLeftRight:  marginLeftRightMobile,
		}
	} else {
		template.MobileConfig = nil
	}

	if template.LogoBand != nil {
		// Xử lý check off custom logo
		if templateRecord.CustomLogo == mysql.TypeOff && templateRecord.Type != mysql.TYPEOutStream {
			template.LogoBand.CustomLogo = nil
		}
	}
	return
}

func (t *Player) MakeTopArticle() (topArticle []payload.TopArticle) {
	topArticle = []payload.TopArticle{
		{
			Title: "Bitcoin Tutorial",
			Image: "https://www.tutorialspoint.com/bitcoin/images/bitcoin-mini-logo.jpg",
			Link:  "https://www.tutorialspoint.com/bitcoin/index.htm",
		},
		{
			Title: "Blockchain Tutorial",
			Image: "https://www.tutorialspoint.com/blockchain/images/blockchain-mini-logo.jpg",
			Link:  "https://www.tutorialspoint.com/blockchain/index.htm",
		},
		{
			Title: "Blue Prism Tutorial",
			Image: "https://www.tutorialspoint.com/blue_prism/images/blueprism-logo.jpg",
			Link:  "https://www.tutorialspoint.com/blockchain/index.htm",
		},
		{
			Title: "Ethereum Tutorial",
			Image: "https://www.tutorialspoint.com/ethereum/images/ethereum-mini-logo.jpg",
			Link:  "https://www.tutorialspoint.com/ethereum/index.htm",
		},
		{
			Title: "OpenShift Tutorial",
			Image: "https://www.tutorialspoint.com/openshift/images/openshift-mini-logo.jpg",
			Link:  "https://www.tutorialspoint.com/openshift/index.htm",
		},
	}
	return
}

func (t *Player) MakeInfo() (infos []payload.Info) {
	infos = []payload.Info{
		{
			Text: "Why a 'Surprise' Kickstarter Hit Actually Took a Year of Intense Planning",
			Link: "https://www.inc.com/christine-lagorio-chafkin/kickstarter-campaign-marketing-boona-shower-head.html",
		},
		{
			Text: "4 To-Dos That Will Make You A Better Leader",
			Link: "https://www.inc.com/peter-cohan/4-to-dos-that-will-make-you-a-better-leader.html",
		},
		{
			Text: "Meta's New Small-Business Service Offers Greater Access to First-Party Customer Data--and It's Free",
			Link: "https://www.inc.com/rebecca-deczynski/meta-small-business-digital-tools-expansion-lead-conversion.html",
		},
		{
			Text: "6 of the Best Newsletters for Business Owners (and Busy People, Generally)",
			Link: "https://www.inc.com/brit-morse/best-newsletters-business-busy-people-email-inbox.html",
		},
		{
			Text: "Why Facial Recognition Technology Has an Uncertain Future With Small Business",
			Link: "https://www.inc.com/steven-i-weiss/clearview-ai-facial-recognition-technology.html",
		},
		{
			Text: "Social Media Marketing Is Costly and Far Less Effective Than It Once Was. Here's What to Do Instead",
			Link: "https://www.inc.com/jordan-hickey/social-media-marketing-alternatives-loyal-customers.html",
		},
		{
			Text: "Is the Pandemic Over? There's a Lot to Consider",
			Link: "https://www.inc.com/brit-morse/pandemic-over-fauci-businesses-cdc-data-mandates.html",
		},
		{
			Text: "3 Ways to Build a Crazy-Successful Sustainable Product Brand",
			Link: "https://www.inc.com/jordan-hickey/richard-branson-grove-collaborative-sustainable-products.html",
		},
		{
			Text: "How Policymakers--Beyond the Fed--Can Curb Inflation",
			Link: "https://www.inc.com/melissa-angell/high-inflation-record-monetary-policy-fed-tariffs-regulation-energy-production.html",
		},
		{
			Text: "How This Behavioral Analyst Who Helps Children With Autism Landed the SBA's Highest Honor",
			Link: "https://www.inc.com/rebecca-deczynski/small-business-administration-small-business-person-of-the-year-bright-futures-jill-scarbro.html",
		},
	}
	return
}

func (t *Player) GetListBoxCollapse(userId, recordId int64, page, typ string) (list []string) {
	switch typ {
	case "add":
		mysql.Client.Select("box_collapse").Model(PageCollapseRecord{}).Where("user_id = ? and page_collapse = ? and is_collapse = ? and page_type = ?", userId, page, 1, typ).Find(&list)
		return
	case "edit":
		mysql.Client.Select("box_collapse").Model(PageCollapseRecord{}).Where("user_id = ? and page_collapse = ? and is_collapse = ? and page_type = ? and page_id = ?", userId, page, 1, typ, recordId).Find(&list)
		return
	}
	return
}

func (t *Player) GetListInventoryByType(typ mysql.TYPEAdType, userId int64) (listInventoryId []int64) {
	mysql.Client.Model(&InventoryAdTagRecord{}).Select("inventory_id").Where("user_id = ? AND type = ?", userId, typ).Group("inventory_id").Find(&listInventoryId)
	return
}
