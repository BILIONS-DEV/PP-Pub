package dto

import (
	"encoding/json"
	"log"
	model "source/internal/entity/model"
	"source/internal/errors"
	"source/pkg/utility"
	"strconv"
)

type PayloadAdSchedule struct {
	ID             int64                     `json:"id"`
	Name           string                    `json:"name"`
	ClientType     string                    `json:"client_type"`
	AdBreakType    string                    `json:"ad_break_type"`
	VpaidMode      string                    `json:"vpaid_mode"`
	AdBreakConfigs []PayloadAdScheduleConfig `json:"ad_break_configs"`
}
type PayloadAdScheduleConfig struct {
	ConfigType      string   `json:"config_type"`
	SkipSecond      int      `json:"skip_second"`
	OverlayAd       bool     `json:"overlay_ad"`
	BreakTiming     int      `json:"break_timing"`
	BreakTimingType string   `json:"break_timing_type"`
	AdTagUrls       []string `json:"ad_tag_urls"`
}

func (p *PayloadAdSchedule) ToJson() string {
	b, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(b)
}
func (p *PayloadAdSchedule) FakeData() *PayloadAdSchedule {
	return &PayloadAdSchedule{
		ID:          0,
		Name:        "my ad schedule test 3",
		ClientType:  string(model.ClientTypeVast),
		VpaidMode:   string(model.VpaidModeDisabled),
		AdBreakType: string(model.AdBreakVmap),
		AdBreakConfigs: []PayloadAdScheduleConfig{
			//{
			//	ConfigType:      string(model.ConfigTypeVmap),
			//	SkipSecond:      10,
			//	OverlayAd:       false,
			//	BreakTiming:     20,
			//	BreakTimingType: string(model.BreakTimingSecondIntoVideo),
			//	AdTagUrls: []string{
			//		"Preroll_test1_url1",
			//		"Preroll_test1_url2",
			//		"Preroll_test1_url3",
			//		"Preroll_test1_url4",
			//		"Preroll_test1_url5",
			//	},
			//},
			{
				ConfigType:      string(model.ConfigTypePreroll),
				SkipSecond:      12,
				OverlayAd:       false,
				BreakTiming:     20,
				BreakTimingType: string(model.BreakTimingSecondIntoVideo),
				AdTagUrls: []string{
					"Preroll_test1_url1",
					"Preroll_test1_url2",
					//"Preroll_test1_url3",
					//"Preroll_test1_url4",
					//"Preroll_test1_url5",
				},
			},
			{
				ConfigType:      string(model.ConfigTypeMidroll),
				SkipSecond:      10,
				OverlayAd:       false,
				BreakTiming:     20,
				BreakTimingType: string(model.BreakTimingSecondIntoVideo),
				AdTagUrls: []string{
					"Midroll_test1_url1",
					"Midroll_test1_url2",
					//"Midroll_test1_url3",
					//"Midroll_test1_url4",
					//"Midroll_test1_url5",
				},
			},
			{
				ConfigType:      string(model.ConfigTypePostroll),
				SkipSecond:      10,
				OverlayAd:       false,
				BreakTiming:     20,
				BreakTimingType: string(model.BreakTimingSecondIntoVideo),
				AdTagUrls: []string{
					"Postroll_test1_url1",
					"Postroll_test1_url2",
					//"Postroll_test1_url3",
					//"Postroll_test1_url4",
					//"Postroll_test1_url5",
				},
			},
		},
	}
}
func (p *PayloadAdSchedule) Validate() (errs []error) {
	if p.Name == "" {
		errs = append(errs, errors.NewWithID(`"Name" is required`, "name"))
	}
	if p.ClientType == "" {
		errs = append(errs, errors.NewWithID(`"Ad Client" is required`, "client_type"))
	}
	if p.AdBreakType == "" {
		errs = append(errs, errors.NewWithID(`"Configure Ad Breaks" is required`, "ad_break_type"))
	}
	switch p.ClientType {
	case string(model.ClientTypeVast):
		if p.VpaidMode != "" {
			errs = append(errs, errors.NewWithID(`"VPAID Mode" is not used in case Client Type is "Google IMA"`, "vpaid_mode"))
		}
	case string(model.ClientTypeGoogleIma):
		if p.VpaidMode == "" {
			errs = append(errs, errors.NewWithID(`"VPAID Mode" is required with Client Type is "Google IMA""`, "vpaid_mode"))
		}

	}
	listConfigTypeOfManually := []string{
		string(model.ConfigTypeMidroll),
		string(model.ConfigTypePostroll),
		string(model.ConfigTypePreroll),
	}
	for k, config := range p.AdBreakConfigs {
		keyString := strconv.Itoa(k)
		if config.ConfigType == "" {
			errs = append(errs, errors.NewWithID(`"Ad Breaks" is required`, "config_type"))
		}
		switch p.AdBreakType {
		case string(model.AdBreakVmap):
			if config.ConfigType != string(model.ConfigTypeVmap) {
				errs = append(errs, errors.NewWithID(`"Ad Breaks `+config.ConfigType+`" is not used in case Ad Break Type is "VMAP"`, "config_type_"+keyString))
			}
		case string(model.AdBreakManually):
			if !utility.InArray(config.ConfigType, listConfigTypeOfManually, true) {
				errs = append(errs, errors.NewWithID(`"Ad Breaks `+config.ConfigType+`" is not used in case Ad Break Type is "Manually"`, "config_type_"+keyString))
			}
		}
	}
	return
}
func (p *PayloadAdSchedule) ToModel() model.AdScheduleModel {
	record := model.AdScheduleModel{}
	record.ID = p.ID
	record.Name = p.Name
	record.ClientType = model.AdScheduleClientTYPE(p.ClientType)
	record.AdBreakType = model.AdScheduleAdBreakTYPE(p.AdBreakType)
	if p.ClientType == string(model.ClientTypeGoogleIma) {
		record.VpaidMode = model.AdScheduleVpaidModeTYPE(p.VpaidMode)
	} else {
		record.VpaidMode = model.VpaidModeNull
	}
	for _, config := range p.AdBreakConfigs {
		var listAdTagUrl []model.AdScheduleConfigAdTagUrlModel
		for _, adTagUrl := range config.AdTagUrls {
			listAdTagUrl = append(listAdTagUrl, model.AdScheduleConfigAdTagUrlModel{
				AdTagUrl: adTagUrl,
			})
		}
		newConfig := model.AdScheduleConfigModel{
			ConfigType:      model.AdScheduleAdBreakConfigTYPE(config.ConfigType),
			SkipSecond:      config.SkipSecond,
			OverlayAd:       config.OverlayAd,
			BreakTiming:     config.BreakTiming,
			BreakTimingType: model.AdScheduleBreakTimingTYPE(config.BreakTimingType),
			AdTagUrls:       listAdTagUrl,
		}
		record.AdBreakConfigs = append(record.AdBreakConfigs, newConfig)
	}
	return record
}

type PayloadAdScheduleFilter struct {
	UserID            int64    `json:"user_id"`
	Page              int      `json:"page"`
	Search            string   `json:"search"`
	ClientTypes       []string `json:"client_types"`
	AdBreakTypes      []string `json:"ad_break_types"`
	AdBreakConfigType []string `json:"ad_break_config_type"`
	Order             string   `json:"order"`
}

func (p *PayloadAdScheduleFilter) Validate() (errs []Error) {
	if p.UserID == 0 {
		errs = append(errs, Error{Message: "permission denied"})
	}
	// kiểm tra xem các mảng dữ liệu đầu vào có nằm trong các field được quy định ở model không...
	return
}
func (p *PayloadAdScheduleFilter) FakeData() *PayloadAdScheduleFilter {
	return &PayloadAdScheduleFilter{
		Page: 1,
		ClientTypes: []string{
			string(model.ClientTypeVast),
			string(model.ClientTypeGoogleIma),
		},
		AdBreakTypes: []string{
			string(model.AdBreakManually),
			string(model.AdBreakVmap),
		},
		AdBreakConfigType: []string{
			string(model.ConfigTypePreroll),
			string(model.ConfigTypeMidroll),
			string(model.ConfigTypePostroll),
			string(model.ConfigTypeVmap),
		},
	}
}
func (p *PayloadAdScheduleFilter) ToJson() string {
	b, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(b)
}
