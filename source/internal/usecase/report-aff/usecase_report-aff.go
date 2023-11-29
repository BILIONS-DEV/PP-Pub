package reportAff

import (
	"fmt"
	//crypt "source/apps/aff/helpers"
	"source/internal/entity/dto/datatable"
	"source/internal/entity/model"
	"source/internal/lang"
	"source/internal/repo"
	reportAff "source/internal/repo/report-aff"
	"source/pkg/logger"
	"source/pkg/utility"
	"strconv"
	"strings"
	"time"
)

type UsecaseReportAff interface {
	Filter(input *InputFilterReportAff) (totalRecords int64, records []*model.ReportAffModel, recordTotal *model.ReportAffModel, err error)
	HandlerReportAffPixel() (err error)
	UpdateDataReportAff() (err error)
	PushCachePixel(record model.PixelAff) (err error)
	GetAllCampaign() (records []*model.ReportAffModel, err error)
	GetAllTrafficSource() (records []*model.ReportAffModel, err error)
	GetAllDemandSource() (records []*model.ReportAffModel, err error)
	GetAllSection() (records []*model.ReportAffModel, err error)
	PushCacheCheckPixel(record model.PixelAffCheck) (err error)
	GetCacheCheckPixel(key string) (record model.PixelAffCheck, err error)
	MarkRequestPixelImpression(record model.PixelAffCheck) (err error)
	MarkRequestPixelSystemTraffic(record model.PixelAffCheck) (err error)
	MarkRequestPixelClick(record model.PixelAffCheck) (err error)
	MarkRequestPixelClick2(record model.PixelAffCheck) (err error)
	MarkRequestPixelClick3(record model.PixelAffCheck) (err error)
	MarkRequestPixelPreEstimatedRevenue(record model.PixelAffCheck, revenue float64) (err error)
	MarkRequestPixelEstimatedRevenue(record model.PixelAffCheck, revenue float64) (err error)
	GetForBlockSectionOutBrain(endDate string) (records []*model.ReportAffModel, err error)
	SaveBlockSectionOutBrain(campaignID string, sectionIDs []string) (err error)
	GetListSectionBlocked(trafficSource, campaignID string) (listSectionBlocked []string, err error)
	CheckSectionBlocked(trafficSource, campaignID, sectionID string) (isBlocked bool, err error)
	GetReportTotalForBlock(campaignID string, date string) (record *model.ReportAffModel, err error)
	GetForReportAllTrafficSourceIDs(trafficSource, account string) (IDs []string, err error)
	UpdateDataFromReportAdsense(dates []string) (err error)
	UpdateDataPixelForReportAff(dates []string) (err error)
	UpdateDataFromReportFacebook(dates []string) (err error)
	UpdateDataFromReportTaboola(dates []string) (err error)
	UpdateDataFromReportTonic(dates []string) (err error)
	UpdateDataFromReportSystem1(dates []string) (err error)
	UpdateDataFromReportDomainActive(dates []string) (err error)
}

type reportAffUsecase struct {
	repos *repo.Repositories
	Trans *lang.Translation
}

func NewReportAffUsecase(repos *repo.Repositories) *reportAffUsecase {
	return &reportAffUsecase{repos: repos}
}

type InputFilterReportAff struct {
	datatable.Request
	UserID         int64
	QuerySearch    string
	StartDate      string
	EndDate        string
	Campaigns      interface{}
	SectionID      interface{}
	TrafficSources interface{}
	DemandSources  interface{}
	GroupBy        interface{}
}

func (t *reportAffUsecase) Filter(input *InputFilterReportAff) (totalRecords int64, records []*model.ReportAffModel, recordTotal *model.ReportAffModel, err error) {
	// => khởi tạo input của Repo
	inputRepo := reportAff.InputFilter{
		UserID:         input.UserID,
		Search:         input.QuerySearch,
		Offset:         input.Start,
		Limit:          input.Length,
		Order:          input.OrderString(),
		StartDate:      input.StartDate,
		EndDate:        input.EndDate,
		Campaigns:      input.Campaigns,
		SectionID:      input.SectionID,
		TrafficSources: input.TrafficSources,
		DemandSources:  input.DemandSources,
		GroupBy:        input.GroupBy,
	}
	return t.repos.ReportAff.Filter(&inputRepo)
}

func (t *reportAffUsecase) HandlerReportAffPixel() (err error) {
	// Get all cache từ aerospike
	cacheRecords, err := t.repos.ReportAff.FindAllCache()
	if err != nil {
		return
	}

	// Tạo map để gom đếm impression + click theo key
	mapData := make(map[string]model.PixelAffForAll)
	// Xử lý cho từng cache
	for _, cacheRecord := range cacheRecords {
		key := cacheRecord.UID + "_" +
			cacheRecord.LayoutID + "_" +
			cacheRecord.LayoutVersion + "_" +
			cacheRecord.Device + "_" +
			cacheRecord.Geo
		// Key tạo từ uid
		if val, ok := mapData[key]; ok {
			val.Impressions += cacheRecord.Impressions
			val.Click += cacheRecord.Click
			val.Click2 += cacheRecord.Click2
			val.Click3 += cacheRecord.Click3
			val.Click4 += cacheRecord.Click4
			val.ImpQuantumdex += cacheRecord.ImpQuantumdex
			val.SystemTraffic += cacheRecord.SystemTraffic
			val.PreBotTraffic += cacheRecord.PreBotTraffic
			val.PreEstimateRevenue += cacheRecord.PreEstimateRevenue
			val.EstimateRevenue += cacheRecord.EstimateRevenue
			if cacheRecord.CPC != "" {
				if val.TrafficSource == "quantumdex" {
					cpcOld, _ := strconv.ParseFloat(val.CPC, 64)
					cpcAdd, _ := strconv.ParseFloat(cacheRecord.CPC, 64)
					val.CPC = fmt.Sprintf("%f", cpcOld+cpcAdd)
				} else {
					val.CPC = cacheRecord.CPC
				}
			}
			if cacheRecord.StyleID != "" {
				val.StyleID = cacheRecord.StyleID
			}
			if cacheRecord.Account != "" {
				val.Account = cacheRecord.Account
			}
			val.Uuid = append(val.Uuid, cacheRecord.Key)
			mapData[cacheRecord.UID] = val
		} else {
			cacheRecord.Uuid = append(cacheRecord.Uuid, cacheRecord.Key)
			mapData[cacheRecord.UID] = cacheRecord
		}
	}

	// Check exist
	// var uuid []string
	// var reportPixels []*model.ReportAffPixelModel
	for _, data := range mapData {
		reportPixel := data.MakeModelReportPixel()

		switch strings.ToLower(data.TrafficSource) {
		case "taboola":
			// Do CPC từ taboola đã được mã hóa phân tích đoạn mã để được cpc điền vào cost
			var account *model.AccountModel
			if data.Account != "" {
				query := make(map[string]interface{})
				query["name"] = data.Account
				account, _ = t.repos.Account.FindByQuery(query)
				reportPixel.Cost, _ = t.repos.ReportTaboola.Decrypt(account.KeyCpcTaboola, data.CPC)
			} else {
				reportPixel.Cost, _ = t.repos.ReportTaboola.Decrypt("6b94068a4c83c171873fa7f6f5171f7f", data.CPC)
			}
			break
		case "pocpoc":
			//redirectID, _ := strconv.ParseInt(data.RedirectID, 10, 64)
			//campaign := t.repos.Campaign.GetById(redirectID, true)
			//VliCrypt := crypt.NewCPCCrypt()
			//if data.CPC != "" && campaign.Bidding != "cpm" {
			//	reportPixel.Cost, _ = strconv.ParseFloat(VliCrypt.Decode(data.CPC), 64)
			//}
			break
		case "quantumdex", "outbrain":
			reportPixel.Cost, _ = strconv.ParseFloat(data.CPC, 64)
			break
		}

		// Save lại pixel đã xử lý
		err = t.repos.ReportAffPixel.Save(&reportPixel)
		if err != nil {
			logger.Error(err.Error())
			fmt.Printf("%+v \n", reportPixel)
			continue
		}

		// Save lại section
		if reportPixel.SectionId != "" {
			_ = t.repos.Section.Save(&model.Section{
				TrafficSource: reportPixel.TrafficSource,
				SectionID:     reportPixel.SectionId,
				SectionName:   reportPixel.SectionName,
				Referrer:      reportPixel.Referrer,
			})
		}

		// Sau khi đã xử lý lên DB xong tiến hành xóa các key đã xử lý xong
		for _, key := range data.Uuid {
			err = t.repos.ReportAff.DeleteCache(key)
			if err != nil {
				if err != nil {
					logger.Error(err.Error())
					fmt.Printf("%+v \n", key)
				}
				continue
			}
		}
	}

	return
}

func (t *reportAffUsecase) UpdateDataReportAff() (err error) {
	// Update data Campaign Name cho report Aff
	//err = t.updateDataCampaignName(dates)
	//if err != nil {
	//	logger.Error(err.Error())
	//	return
	//}

	// Update các data lấy từ OutBrain
	go func() {
		for {
			// Lấy tất cả report Pixel của Aff của ngày hôm nay và hôm qua
			var dates []string
			layout := "2006-01-02"
			today := time.Now().UTC()
			dates = append(dates, today.Format(layout))
			yesterday := today.AddDate(0, 0, -1)
			dates = append(dates, yesterday.Format(layout))
			timeStartUpdate := time.Now().UTC()
			err = t.updateDataFromReportOutBrain(dates)
			if err != nil {
				logger.Error(err.Error())
			}
			fmt.Println("Done Update From OutBrain --Start: ", timeStartUpdate, " --End: ", time.Now().UTC(), "Total Time:", time.Now().UTC().Sub(timeStartUpdate).Milliseconds(), "ms")
			time.Sleep(5 * time.Minute)
		}
	}()

	// Update các data lấy từ Taboola
	go func() {
		for {
			// Lấy tất cả report Pixel của Aff của ngày hôm nay và hôm qua
			var dates []string
			layout := "2006-01-02"
			today := time.Now().UTC()
			dates = append(dates, today.Format(layout))
			yesterday := today.AddDate(0, 0, -1)
			dates = append(dates, yesterday.Format(layout))
			timeStartUpdate := time.Now().UTC()
			err = t.UpdateDataFromReportTaboola(dates)
			if err != nil {
				logger.Error(err.Error())
			}

			fmt.Println("Done Update From Taboola --Start: ", timeStartUpdate, " --End: ", time.Now().UTC(), "Total Time:", time.Now().UTC().Sub(timeStartUpdate).Milliseconds(), "ms")
			time.Sleep(5 * time.Minute)
		}
	}()

	// Update các data lấy từ Mgid
	go func() {
		for {
			// Lấy tất cả report Pixel của Aff của ngày hôm nay và hôm qua
			var dates []string
			layout := "2006-01-02"
			today := time.Now().UTC()
			dates = append(dates, today.Format(layout))
			yesterday := today.AddDate(0, 0, -1)
			dates = append(dates, yesterday.Format(layout))
			timeStartUpdate := time.Now().UTC()
			err = t.updateDataFromReportMgid(dates)
			if err != nil {
				logger.Error(err.Error())
			}

			fmt.Println("Done Update From Mgid --Start: ", timeStartUpdate, " --End: ", time.Now().UTC(), "Total Time:", time.Now().UTC().Sub(timeStartUpdate).Milliseconds(), "ms")
			time.Sleep(5 * time.Minute)
		}
	}()

	// Update các data lấy từ PocPoc
	go func() {
		for {
			// Lấy tất cả report Pixel của Aff của ngày hôm nay và hôm qua
			var dates []string
			layout := "2006-01-02"
			today := time.Now().UTC()
			dates = append(dates, today.Format(layout))
			yesterday := today.AddDate(0, 0, -1)
			dates = append(dates, yesterday.Format(layout))
			timeStartUpdate := time.Now().UTC()
			err = t.updateDataFromReportPocPoc(dates)
			if err != nil {
				logger.Error(err.Error())
			}

			fmt.Println("Done Update From PocPoc --Start: ", timeStartUpdate, " --End: ", time.Now().UTC(), "Total Time:", time.Now().UTC().Sub(timeStartUpdate).Milliseconds(), "ms")
			time.Sleep(5 * time.Minute)
		}
	}()

	// Update các data lấy từ System1
	go func() {
		for {
			// Lấy tất cả report Pixel của Aff của ngày hôm nay và hôm qua
			var dates []string
			layout := "2006-01-02"
			today := time.Now().UTC()
			dates = append(dates, today.Format(layout))
			yesterday := today.AddDate(0, 0, -1)
			dates = append(dates, yesterday.Format(layout))
			timeStartUpdate := time.Now().UTC()
			err = t.UpdateDataFromReportSystem1(dates)
			if err != nil {
				logger.Error(err.Error())
			}
			fmt.Println("Done Update From System1 --Start: ", timeStartUpdate, " --End: ", time.Now().UTC(), "Total Time:", time.Now().UTC().Sub(timeStartUpdate).Milliseconds(), "ms")
			time.Sleep(5 * time.Minute)
		}
	}()

	// Update các data lấy từ Tonic
	go func() {
		for {
			// Lấy tất cả report Pixel của Aff của ngày hôm nay và hôm qua
			var dates []string
			layout := "2006-01-02"
			today := time.Now().UTC()
			dates = append(dates, today.Format(layout))
			yesterday := today.AddDate(0, 0, -2)
			dates = append(dates, yesterday.Format(layout))
			timeStartUpdate := time.Now().UTC()
			err = t.UpdateDataFromReportTonic(dates)
			if err != nil {
				logger.Error(err.Error())
			}
			fmt.Println("Done Update From Tonic --Start: ", timeStartUpdate, " --End: ", time.Now().UTC(), "Total Time:", time.Now().UTC().Sub(timeStartUpdate).Milliseconds(), "ms")
			time.Sleep(5 * time.Minute)
		}
	}()

	// Update các data lấy từ TikTok
	go func() {
		for {
			// Lấy tất cả report Pixel của Aff của ngày hôm nay và hôm qua
			var dates []string
			layout := "2006-01-02"
			today := time.Now().UTC()
			dates = append(dates, today.Format(layout))
			yesterday := today.AddDate(0, 0, -1)
			dates = append(dates, yesterday.Format(layout))
			timeStartUpdate := time.Now().UTC()
			err = t.updateDataFromReportTikTok(dates)
			if err != nil {
				logger.Error(err.Error())
			}
			fmt.Println("Done Update From TikTok --Start: ", timeStartUpdate, " --End: ", time.Now().UTC(), "Total Time:", time.Now().UTC().Sub(timeStartUpdate).Milliseconds(), "ms")
			time.Sleep(5 * time.Minute)
		}
	}()

	// Update các data lấy từ Facebook
	go func() {
		for {
			// Lấy tất cả report Pixel của Aff của ngày hôm nay và hôm qua
			var dates []string
			layout := "2006-01-02"
			today := time.Now().UTC()
			dates = append(dates, today.Format(layout))
			yesterday := today.AddDate(0, 0, -1)
			dates = append(dates, yesterday.Format(layout))
			timeStartUpdate := time.Now().UTC()
			err = t.UpdateDataFromReportFacebook(dates)
			if err != nil {
				logger.Error(err.Error())
			}
			fmt.Println("Done Update From Facebook --Start: ", timeStartUpdate, " --End: ", time.Now().UTC(), "Total Time:", time.Now().UTC().Sub(timeStartUpdate).Milliseconds(), "ms")
			time.Sleep(2 * time.Minute)
		}
	}()

	// Update các data lấy từ Facebook
	go func() {
		for {
			// Lấy tất cả report domain active của ngày hôm nay và hôm qua
			var dates []string
			layout := "2006-01-02"
			today := time.Now().UTC()
			dates = append(dates, today.Format(layout))
			yesterday := today.AddDate(0, 0, -1)
			dates = append(dates, yesterday.Format(layout))
			timeStartUpdate := time.Now().UTC()
			err = t.UpdateDataFromReportDomainActive(dates)
			if err != nil {
				logger.Error(err.Error())
			}
			fmt.Println("Done Update From Domain Active --Start: ", timeStartUpdate, " --End: ", time.Now().UTC(), "Total Time:", time.Now().UTC().Sub(timeStartUpdate).Milliseconds(), "ms")
			time.Sleep(5 * time.Minute)
		}
	}()
	return
}

func (t *reportAffUsecase) UpdateDataPixelForReportAff(dates []string) (err error) {
	for _, date := range dates {
		reportPixels, _ := t.repos.ReportAffPixel.FindForReportAff(date)
		err = t.updateDataPixelForReportAff(date, reportPixels)
		if err != nil {
			logger.Error(err.Error())
			return err
		}
	}
	return
}

func (t *reportAffUsecase) updateDataPixelForReportAff(date string, datas []*model.ReportAffPixelModel) (err error) {
	// fmt.Println(data.TimeConversion)
	// date, err := time.Parse("2006-01-02 15:04:05", data.TimeConversion)
	// if err != nil {
	//	fmt.Println(err)
	//	return
	// }
	// layoutDay := "2006-01-02"
	// day := date.Format(layoutDay)
	// Tạo query và get recordOld nếu có

	mapReportAff := make(map[string]*model.ReportAffModel)

	for _, data := range datas {
		if data.CampaignID == "" || data.CampaignID == "unknow" {
			continue
		}
		var campaignID string

		timeReport, err := time.Parse("2006-01-02", date)
		if err != nil {
			continue
		}
		timeStartNew := time.Date(2022, 12, 01, 0, 0, 0, 0, time.UTC)
		if timeReport.Sub(timeStartNew).Nanoseconds() >= 0 {
			campaignID = data.CampaignID
		} else {
			if strings.ToLower(data.TrafficSource) == "taboola" {
				campaignID = data.RedirectID
			} else {
				campaignID = data.CampaignID
			}
		}
		query := make(map[string]interface{})
		query["date"] = date
		query["traffic_source"] = data.TrafficSource
		query["partner"] = data.DemandSource
		//query["campaign_id"] = campaignID
		query["redirect_id"] = data.RedirectID
		query["section_id"] = data.SectionId
		query["style_id"] = data.StyleID
		query["layout_id"] = data.LayoutID
		query["layout_version"] = data.LayoutVersion
		query["device"] = data.Device
		query["geo"] = data.Geo
		recordOld, _ := t.repos.ReportAff.FindOneByQuery(query)
		redirectID, _ := strconv.ParseInt(data.RedirectID, 10, 64)

		key := date + "_" +
			data.TrafficSource + "_" +
			strconv.FormatInt(redirectID, 10) + "_" +
			data.SectionId + "_" +
			data.StyleID + "_" +
			strconv.FormatInt(data.LayoutID, 10) + "_" +
			data.LayoutVersion + "_" +
			data.Device + "_" +
			data.Geo
		if val, ok := mapReportAff[key]; ok {
			if data.DemandSource == "adsense" {
				val.Impressions += data.Click2
				val.Click += data.Click3
				val.Click4 += data.Click4
			} else if data.DemandSource == "system1" {
				//nothing
			} else {
				val.Impressions += data.Impressions
				val.Click += data.Click
			}
			val.ImpQuantumdex += data.ImpQuantumdex
			val.SystemTraffic += data.SystemTraffic
			val.PreBotTraffic += data.PreBotTraffic
			val.BotTraffic += data.BotTraffic
			if data.DemandSource == "hooligapps" || data.DemandSource == "pdl" {
				val.Revenue = utility.SumFloat64(val.Revenue, data.EstimateRevenue)
			} else {
				val.EstimatedRevenue += data.EstimateRevenue
			}
			val.PreEstimatedRevenue += data.PreEstimateRevenue
			if data.TrafficSource == "pocpoc" {
				//val.Cost += data.Cost
				//val.SupplyClick += data.PreBotTraffic
			} else if data.TrafficSource == "quantumdex" {
				val.Cost = utility.SumFloat64(val.Cost, data.Cost)
			}
			if data.TrafficSource != "tiktok" && data.TrafficSource != "facebook" {
				val.CostCPC = utility.SumFloat64(val.CostCPC, data.Cost)
			}
			mapReportAff[key] = val
		} else {
			// Xử lý tạo mới nếu record chưa tồn tại
			if recordOld.ID == 0 {
				record := model.ReportAffModel{
					Date:                date,
					UserID:              data.UserID,
					TrafficSource:       data.TrafficSource,
					SectionId:           data.SectionId,
					SectionName:         data.SectionName,
					CampaignID:          campaignID,
					RedirectID:          redirectID,
					StyleID:             data.StyleID,
					LayoutID:            data.LayoutID,
					LayoutVersion:       data.LayoutVersion,
					Device:              data.Device,
					Geo:                 data.Geo,
					Partner:             data.DemandSource,
					Impressions:         data.Impressions,
					Click:               data.Click,
					ImpQuantumdex:       data.ImpQuantumdex,
					SystemTraffic:       data.SystemTraffic,
					PreBotTraffic:       data.PreBotTraffic,
					BotTraffic:          data.BotTraffic,
					PreEstimatedRevenue: data.PreEstimateRevenue,
				}
				if data.DemandSource == "adsense" {
					record.Impressions = data.Click2
					record.Click = data.Click3
					record.Click4 = data.Click4
				} else if data.DemandSource == "system1" {
					//nothing
				} else {
					record.Impressions = data.Impressions
					record.Click = data.Click
				}
				if data.DemandSource == "hooligapps" || data.DemandSource == "pdl" {
					record.Revenue = data.EstimateRevenue
				} else {
					record.EstimatedRevenue = data.EstimateRevenue
				}
				if data.TrafficSource == "pocpoc" {
					//record.Cost = data.Cost
					//record.SupplyClick = data.PreBotTraffic
				} else if data.TrafficSource == "quantumdex" {
					record.Cost = data.Cost
				}
				if data.TrafficSource != "tiktok" && data.TrafficSource != "facebook" {
					record.CostCPC = data.Cost
				}
				mapReportAff[key] = &record
			} else {
				// Nếu như đã tồn tại record update những field từ recordOld
				recordOld.SectionName = data.SectionName
				recordOld.Partner = data.DemandSource
				recordOld.RedirectID = redirectID
				if data.DemandSource == "adsense" {
					recordOld.Impressions = data.Click2
					recordOld.Click = data.Click3
					recordOld.Click4 = data.Click4
				} else if data.DemandSource == "system1" {
					//nothing
				} else {
					recordOld.Impressions = data.Impressions
					recordOld.Click = data.Click
				}
				recordOld.ImpQuantumdex = data.ImpQuantumdex
				recordOld.SystemTraffic = data.SystemTraffic
				recordOld.PreBotTraffic = data.PreBotTraffic
				recordOld.BotTraffic = data.BotTraffic
				if data.DemandSource == "hooligapps" || data.DemandSource == "pdl" {
					recordOld.Revenue = data.EstimateRevenue
				} else {
					recordOld.EstimatedRevenue = data.EstimateRevenue
				}
				recordOld.PreEstimatedRevenue = data.PreEstimateRevenue

				if data.TrafficSource == "pocpoc" {
					//recordOld.Cost = data.Cost
					//recordOld.SupplyClick = data.PreBotTraffic
				} else if data.TrafficSource == "quantumdex" {
					recordOld.Cost = data.Cost
				}
				if data.TrafficSource != "tiktok" && data.TrafficSource != "facebook" {
					recordOld.CostCPC = data.Cost
				}
				mapReportAff[key] = recordOld
			}
		}
	}

	for _, record := range mapReportAff {
		//fmt.Printf("%+v \n", record)
		err = t.repos.ReportAff.Save(record)
		if err != nil {
			logger.Error(err.Error())
		}
	}
	return
}

func (t *reportAffUsecase) updateDataFromReportOutBrain(dates []string) (err error) {
	fmt.Println("====== Update Data From OutBrain For Report Aff ======")

	// Tạo reportAff để update hoặc create
	var recordReportAffs []*model.ReportAffModel

	for _, date := range dates {
		reportOutBrains, _ := t.repos.ReportOutBrain.FindByDayForReportAff(date)
		for _, reportOutBrain := range reportOutBrains {
			campaign, err := t.repos.Campaign.GetByTrafficSourceID(reportOutBrain.CampaignID)
			if err != nil {
				continue
			}
			query := make(map[string]interface{})
			query["date"] = date
			query["traffic_source"] = "outbrain"
			query["partner"] = strings.ToLower(campaign.DemandSource)
			//query["campaign_id"] = reportOutBrain.CampaignID
			query["redirect_id"] = campaign.ID
			query["section_id"] = reportOutBrain.SectionID
			query["style_id"] = ""
			query["layout_id"] = ""
			query["layout_version"] = ""
			query["device"] = ""
			query["geo"] = ""
			recordReportAff, _ := t.repos.ReportAff.FindOneByQuery(query)
			if recordReportAff.ID != 0 { // Nếu đã tồn tại thì update các value cần
				// Từ report OutBrain có thể update được Cost(spend), Supply Click(clicks), Supply Conversion(conversions) cho reportAff
				recordReportAff.Cost = reportOutBrain.Spend
				recordReportAff.SupplyClick = reportOutBrain.Clicks
				recordReportAff.SupplyConversions = reportOutBrain.Conversions

				// Update lại các row thiếu data
				recordReportAff.RedirectID = campaign.ID
				recordReportAff.Partner = strings.ToLower(campaign.DemandSource)
			} else { // / Nếu chưa tồn tại tạo mới 1 record
				recordReportAff = &model.ReportAffModel{
					Date:              date,
					UserID:            campaign.UserID,
					TrafficSource:     "outbrain",
					SectionId:         reportOutBrain.SectionID,
					CampaignID:        reportOutBrain.CampaignID,
					CampaignName:      reportOutBrain.CampaignName,
					StyleID:           "",
					Cost:              reportOutBrain.Spend,
					SupplyClick:       reportOutBrain.Clicks,
					SupplyConversions: reportOutBrain.Conversions,
					RedirectID:        campaign.ID,
					Partner:           strings.ToLower(campaign.DemandSource),
				}
			}
			recordReportAffs = append(recordReportAffs, recordReportAff)
		}
	}

	// Sau khi xử lý xong update lại các records
	for _, recordReportAff := range recordReportAffs {
		err = t.repos.ReportAff.Save(recordReportAff)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
	}
	return
}

func (t *reportAffUsecase) UpdateDataFromReportTaboola(dates []string) (err error) {
	fmt.Println("====== Update Data From Taboola For Report Aff ======")
	// Tạo reportAff để update hoặc create
	mapReportTaboola := make(map[string]*model.ReportTaboolaModel)
	for _, date := range dates {
		reportTaboolas, _ := t.repos.ReportTaboola.FindByDayForReportAff(date)
		for _, reportTaboola := range reportTaboolas {
			var campaignID string
			campaignID = reportTaboola.CampaignID
			key := date + "_" + "taboola" + "_" + campaignID + "_" + reportTaboola.SiteID
			if val, ok := mapReportTaboola[key]; ok {
				val.Spent += reportTaboola.Spent
				val.Clicks += reportTaboola.Clicks
				val.CpaActionsNum += reportTaboola.CpaActionsNum
			} else {
				mapReportTaboola[key] = reportTaboola
			}
		}
	}
	var recordReportAffs []*model.ReportAffModel
	for key, reportTaboola := range mapReportTaboola {
		if reportTaboola.Spent <= 0 {
			continue
		}
		splitKey := strings.Split(key, "_")
		date := splitKey[0]
		trafficSource := splitKey[1]
		campaignID := splitKey[2]
		sectionID := splitKey[3]

		var campaign model.CampaignModel
		var recordReportAff *model.ReportAffModel
		timeReport, err := time.Parse("2006-01-02 15:04:05.0", reportTaboola.Date)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		timeStartNew := time.Date(2022, 12, 01, 0, 0, 0, 0, time.UTC)
		if timeReport.Sub(timeStartNew).Nanoseconds() >= 0 {
			// new
			campaign, err = t.repos.Campaign.GetByTrafficSourceID(campaignID)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
		} else {
			// old
			rid, _ := strconv.ParseInt(campaignID, 10, 64)
			campaign = t.repos.Campaign.GetById(rid, true)
		}
		if campaign.ID == 0 {
			continue
		}
		query := make(map[string]interface{})
		query["date"] = date
		query["traffic_source"] = trafficSource
		query["partner"] = strings.ToLower(campaign.DemandSource)
		//query["campaign_id"] = campaignID
		query["redirect_id"] = campaign.ID
		query["section_id"] = sectionID
		query["style_id"] = ""
		query["layout_id"] = ""
		query["layout_version"] = ""
		query["device"] = ""
		query["geo"] = ""
		recordReportAff, _ = t.repos.ReportAff.FindOneByQuery(query)

		if recordReportAff.ID != 0 { // Nếu đã tồn tại thì update các value cần
			// Từ report OutBrain có thể update được Cost(spend), Supply Click(clicks), Supply Conversion(conversions) cho reportAff
			recordReportAff.Cost = reportTaboola.Spent
			recordReportAff.SupplyClick = reportTaboola.Clicks
			recordReportAff.SupplyConversions = reportTaboola.CpaActionsNum
		} else { // / Nếu chưa tồn tại tạo mới 1 record
			recordReportAff = &model.ReportAffModel{
				Date:              date,
				UserID:            campaign.UserID,
				TrafficSource:     "taboola",
				SectionId:         reportTaboola.SiteID,
				SectionName:       reportTaboola.Site,
				CampaignID:        campaignID,
				CampaignName:      campaign.Name,
				StyleID:           "",
				Cost:              reportTaboola.Spent,
				SupplyClick:       reportTaboola.Clicks,
				SupplyConversions: reportTaboola.CpaActionsNum,
				RedirectID:        campaign.ID,
				Partner:           strings.ToLower(campaign.DemandSource),
			}
		}
		recordReportAffs = append(recordReportAffs, recordReportAff)
	}

	size := 500
	var j int
	for i := 0; i < len(recordReportAffs); i += size {
		j += size
		if j > len(recordReportAffs) {
			j = len(recordReportAffs)
		}
		// Log lại toàn bộ report vào DB
		err = t.repos.ReportAff.SaveSlice(recordReportAffs[i:j])
		if err != nil {
			logger.Error(err.Error())
			continue
		}
	}
	return
}

func (t *reportAffUsecase) UpdateDataFromReportSystem1(dates []string) (err error) {
	fmt.Println("====== Update Data From System1 For Report Aff ======")

	// Tạo reportAff để update hoặc create
	mapReportAff := make(map[string]*model.ReportAffModel)
	for _, date := range dates {
		reportSystem1s, _ := t.repos.ReportOpenMailSubID.FindByDayForReportAff(date)
		for _, reportSystem1 := range reportSystem1s { // Từ subID lấy được trafficSource, campaign, sectionID, publisherID
			var day, trafficSource, campaignID, sectionID, styleID string
			splitSubID := strings.Split(reportSystem1.SubID, ":")
			if len(splitSubID) < 2 {
				// err = errors.New("subID malformed")
				continue
			}
			var campaign model.CampaignModel
			subID1 := splitSubID[0] // subID1 trafficSource_campaignID || uid
			if subID1 == "uid" {
				uid := splitSubID[1]

				// Tạm thời để để xử lý các subID bị lỗi sau có thể bỏ đi
				tempSplitSubID := strings.Split(splitSubID[1], "cache_buster")
				if len(tempSplitSubID) >= 2 {
					uid = tempSplitSubID[0]
				}

				// Nếu subID1 là clickID thì subID 2 chứa clickID để tìm map với bảng report_aff_pixel
				query := make(map[string]interface{})
				query["uid"] = uid
				record, err := t.repos.ReportAffPixel.FindOneByQuery(query)
				if err != nil {
					continue
				}
				redirectID, _ := strconv.ParseInt(record.RedirectID, 10, 64)
				campaign = t.repos.Campaign.GetById(redirectID, true)
				date, _ := time.Parse("2006-01-02 15:04:05", record.TimeConversion)
				layoutDay := "2006-01-02"
				day = date.Format(layoutDay)
				trafficSource = record.TrafficSource
				campaignID = record.CampaignID
				sectionID = record.SectionId
				styleID = record.StyleID
			}

			if campaign.ID == 0 {
				continue
			}
			// Tạo query để tìm report tương ứng từ reportAff
			query := make(map[string]interface{})
			query["date"] = day
			query["traffic_source"] = trafficSource
			query["partner"] = strings.ToLower(campaign.DemandSource)
			//query["campaign_id"] = campaignID
			query["redirect_id"] = campaign.ID
			query["section_id"] = sectionID
			query["style_id"] = styleID
			query["layout_id"] = ""
			query["layout_version"] = ""
			query["device"] = ""
			query["geo"] = ""
			key := day + "_" + trafficSource + "_" + strconv.FormatInt(campaign.ID, 10) + "_" + sectionID + "_" + styleID
			if val, ok := mapReportAff[key]; ok {
				val.Revenue = utility.SumFloat64(val.Revenue, reportSystem1.EstimatedRevenue)
				val.Impressions += reportSystem1.Searches
				//val.Click += reportSystem1.Clicks
				val.ClickAdsense += reportSystem1.Clicks
				mapReportAff[key] = val
			} else {
				recordReportAff, _ := t.repos.ReportAff.FindOneByQuery(query)
				if recordReportAff.ID != 0 { // Nếu đã tồn tại thì update các value cần
					// Từ report OutBrain có thể update được Revenue(estimated_revenue) cho reportAff
					recordReportAff.Revenue = reportSystem1.EstimatedRevenue
					recordReportAff.Impressions = reportSystem1.Searches
					//recordReportAff.Click = reportSystem1.Clicks
					recordReportAff.ClickAdsense = reportSystem1.Clicks
				} else { // / Nếu chưa tồn tại tạo mới 1 record
					recordReportAff = &model.ReportAffModel{
						UserID:        campaign.UserID,
						Date:          day,
						Partner:       "system1",
						TrafficSource: trafficSource,
						SectionId:     sectionID,
						CampaignID:    campaignID,
						Revenue:       reportSystem1.EstimatedRevenue,
						Impressions:   reportSystem1.Searches,
						//Click:         reportSystem1.Clicks,
						ClickAdsense: reportSystem1.Clicks,
						RedirectID:   campaign.ID,
					}
				}
				mapReportAff[key] = recordReportAff
			}

		}
	}

	// Sau khi xử lý xong update lại các records
	for _, record := range mapReportAff {
		_ = t.repos.ReportAff.Save(record)
	}

	return
}

func (t *reportAffUsecase) UpdateDataFromReportTonic(dates []string) (err error) {
	fmt.Println("====== Update Data From Tonic For Report Aff ======")

	dates = []string{}
	var date time.Time
	// dates = append(dates, "2022-08-11")
	date = time.Now().UTC().AddDate(0, 0, -2)
	dates = append(dates, date.Format("2006-01-02"))
	date = time.Now().UTC().AddDate(0, 0, -1)
	dates = append(dates, date.Format("2006-01-02"))
	date = time.Now().UTC()
	dates = append(dates, date.Format("2006-01-02"))
	// fmt.Println(dates)
	// Tạo reportAff để update hoặc create
	mapReportAff := make(map[string]*model.ReportAffModel)
	for _, date := range dates {
		reportTonic, _ := t.repos.ReportTonic.FindByDayForReportAff(date)
		for _, report := range reportTonic {
			var day, trafficSource, campaignID, sectionID, styleID string
			record := new(model.ReportAffPixelModel)
			query := make(map[string]interface{})
			if len(report.Subid1) > 99 {
				record, err = t.repos.ReportAffPixel.FindOneByQuery(query, "uid like '"+report.Subid1+"%'")
				if err != nil {
					continue
				}
			} else {
				query["uid"] = report.Subid1
				record, err = t.repos.ReportAffPixel.FindOneByQuery(query)
				if err != nil {
					continue
				}
			}

			redirectID, _ := strconv.ParseInt(record.RedirectID, 10, 64)
			campaign := t.repos.Campaign.GetById(redirectID, true)
			if campaign.ID == 0 {
				continue
			}
			loc, _ := time.LoadLocation("America/Los_Angeles")
			timeTonic, _ := time.ParseInLocation("2006-01-02 15:04:05", report.Timestamp, loc)
			layoutDay := "2006-01-02"
			day = timeTonic.UTC().Format(layoutDay)
			trafficSource = record.TrafficSource
			campaignID = record.CampaignID
			sectionID = record.SectionId
			styleID = record.StyleID

			revenue, _ := strconv.ParseFloat(report.RevenueUsd, 64)
			clicks, _ := strconv.ParseInt(report.Clicks, 10, 64)
			// Tạo query để tìm report tương ứng từ reportAff
			query = make(map[string]interface{})
			query["date"] = day
			query["traffic_source"] = trafficSource
			query["partner"] = "tonic"
			//query["campaign_id"] = campaignID
			query["redirect_id"] = campaign.ID
			query["section_id"] = sectionID
			query["style_id"] = styleID
			query["layout_id"] = ""
			query["layout_version"] = ""
			query["device"] = ""
			query["geo"] = ""
			key := day + "_" + trafficSource + "_" + strconv.FormatInt(campaign.ID, 10) + "_" + sectionID + "_" + styleID
			if val, ok := mapReportAff[key]; ok {
				val.Revenue = utility.SumFloat64(val.Revenue, revenue)
				mapReportAff[key] = val
			} else {
				recordReportAff, _ := t.repos.ReportAff.FindOneByQuery(query)
				if recordReportAff.ID != 0 { // Nếu đã tồn tại thì update các value cần
					recordReportAff.Revenue = revenue
					recordReportAff.ClickAdsense = clicks
				} else { // / Nếu chưa tồn tại tạo mới 1 record
					recordReportAff = &model.ReportAffModel{
						UserID:        campaign.UserID,
						Partner:       "tonic",
						Date:          day,
						TrafficSource: trafficSource,
						SectionId:     sectionID,
						CampaignID:    campaignID,
						StyleID:       styleID,
						Revenue:       revenue,
						ClickAdsense:  clicks,
						RedirectID:    campaign.ID,
					}
				}
				mapReportAff[key] = recordReportAff
			}
		}
	}

	// Sau khi xử lý xong update lại các records
	for _, record := range mapReportAff {
		_ = t.repos.ReportAff.Save(record)
	}
	return
}

func (t *reportAffUsecase) UpdateDataFromReportAdsense(dates []string) (err error) {
	fmt.Println("====== Update Data From Adsense For Report Aff ======")

	// fmt.Println(dates)
	// Tạo reportAff để update hoặc create
	mapReportAff := make(map[string]*model.ReportAffModel)
	for _, date := range dates {
		reportAdsense, _ := t.repos.ReportAdsense.FindByDayForReportAff(date)
		for _, report := range reportAdsense {
			var campaign model.CampaignModel
			var day, trafficSource, campaignID, sectionID, styleID string
			day = report.Date
			var redirectID int64
			// Xử lý report subDomain
			if report.Type == model.TYPEReportAdsenseSubDomain {
				splitHost := strings.Split(report.DomainCode, ".")
				if len(splitHost) < 3 {
					continue
				}
				subDomain := splitHost[0]
				if strings.Contains(subDomain, "_") {
					splitSubDomain := strings.Split(subDomain, "_")
					if len(splitSubDomain) < 2 {
						continue
					}
					redirectID, _ = strconv.ParseInt(splitSubDomain[0], 10, 64)
					sectionID = splitSubDomain[1]
				} else {
					params, err := t.repos.Campaign.FindCacheSubDomainParams(subDomain)
					if err != nil {
						continue
					}
					redirectID = params.CampaignID
					sectionID = params.SectionID
				}
			} else if report.Type == model.TYPEReportAdsenseChannel {
				// Xử lý report channel
				if report.CampaignID == 0 {
					continue
				}
				redirectID = report.CampaignID
				sectionID = ""
			} else {
				continue
			}
			campaign = t.repos.Campaign.GetById(redirectID, true)
			if campaign.ID == 0 {
				continue
			}
			campaignID = campaign.TrafficSourceID
			trafficSource = strings.ToLower(campaign.TrafficSource)
			// Tạo query để tìm report tương ứng từ reportAff
			query := make(map[string]interface{})
			query["date"] = day
			query["traffic_source"] = trafficSource
			query["partner"] = "adsense"
			//query["campaign_id"] = campaignID
			query["redirect_id"] = campaign.ID
			query["section_id"] = sectionID
			query["style_id"] = styleID
			query["layout_id"] = ""
			query["layout_version"] = ""
			query["device"] = ""
			query["geo"] = ""
			key := day + "_" + trafficSource + "_" + strconv.FormatInt(campaign.ID, 10) + "_" + sectionID + "_" + styleID
			if val, ok := mapReportAff[key]; ok {
				val.Revenue = utility.SumFloat64(val.Revenue, report.EstimatedEarnings)
				val.ClickAdsense += report.Clicks
				mapReportAff[key] = val
			} else {
				recordReportAff, _ := t.repos.ReportAff.FindOneByQuery(query)
				if recordReportAff.ID != 0 { // Nếu đã tồn tại thì update các value cần
					recordReportAff.Revenue = report.EstimatedEarnings
					recordReportAff.ClickAdsense = report.Clicks
				} else { // / Nếu chưa tồn tại tạo mới 1 record
					recordReportAff = &model.ReportAffModel{
						UserID:        campaign.UserID,
						Date:          day,
						TrafficSource: trafficSource,
						Partner:       "adsense",
						SectionId:     sectionID,
						CampaignID:    campaignID,
						RedirectID:    campaign.ID,
						StyleID:       styleID,
						ClickAdsense:  report.Clicks,
						Revenue:       report.EstimatedEarnings,
					}
				}
				mapReportAff[key] = recordReportAff
			}
		}
	}

	// Sau khi xử lý xong update lại các records
	for _, record := range mapReportAff {
		_ = t.repos.ReportAff.Save(record)
	}

	return
}

func (t *reportAffUsecase) GetAllCampaign() (records []*model.ReportAffModel, err error) {
	return t.repos.ReportAff.FindAllByGroup("campaign_id")
}

func (t *reportAffUsecase) GetAllTrafficSource() (records []*model.ReportAffModel, err error) {
	return t.repos.ReportAff.FindAllByGroup("traffic_source")
}

func (t *reportAffUsecase) GetAllDemandSource() (records []*model.ReportAffModel, err error) {
	return t.repos.ReportAff.FindAllByGroup("partner")
}

func (t *reportAffUsecase) GetAllSection() (records []*model.ReportAffModel, err error) {
	return t.repos.ReportAff.FindAllByGroup("section_id")
}

func (t *reportAffUsecase) PushCachePixel(record model.PixelAff) (err error) {
	return t.repos.ReportAff.PushCache(record)
}

func (t *reportAffUsecase) PushCacheCheckPixel(record model.PixelAffCheck) (err error) {
	return t.repos.ReportAff.PushCacheCheck(record.Key, record)
}

func (t *reportAffUsecase) MarkRequestPixelImpression(record model.PixelAffCheck) (err error) {
	return t.repos.ReportAff.PushBinCache(record.SetName(), record.Key, "impressions", 1)
}

func (t *reportAffUsecase) MarkRequestPixelSystemTraffic(record model.PixelAffCheck) (err error) {
	return t.repos.ReportAff.PushBinCache(record.SetName(), record.Key, "system_traffic", 1)
}

func (t *reportAffUsecase) MarkRequestPixelClick(record model.PixelAffCheck) (err error) {
	return t.repos.ReportAff.PushBinCache(record.SetName(), record.Key, "click", 1)
}

func (t *reportAffUsecase) MarkRequestPixelClick2(record model.PixelAffCheck) (err error) {
	return t.repos.ReportAff.PushBinCache(record.SetName(), record.Key, "click2", 1)
}

func (t *reportAffUsecase) MarkRequestPixelClick3(record model.PixelAffCheck) (err error) {
	return t.repos.ReportAff.PushBinCache(record.SetName(), record.Key, "click3", 1)
}

func (t *reportAffUsecase) MarkRequestPixelPreEstimatedRevenue(record model.PixelAffCheck, revenue float64) (err error) {
	return t.repos.ReportAff.PushBinCache(record.SetName(), record.Key, "pre_revenue", revenue)
}

func (t *reportAffUsecase) MarkRequestPixelEstimatedRevenue(record model.PixelAffCheck, revenue float64) (err error) {
	return t.repos.ReportAff.PushBinCache(record.SetName(), record.Key, "revenue", revenue)
}

func (t *reportAffUsecase) GetCacheCheckPixel(key string) (record model.PixelAffCheck, err error) {
	return t.repos.ReportAff.FindCacheCheck(key)
}

type InputBlockSectionOutBrain struct {
	StartDemandTrafficCheck int64
}

func (t *reportAffUsecase) GetForBlockSectionOutBrain(endDate string) (records []*model.ReportAffModel, err error) {
	return t.repos.ReportAff.FindForBlockSectionOutBrain(endDate)
}

func (t *reportAffUsecase) SaveBlockSectionOutBrain(campaignID string, sectionIDs []string) (err error) {
	err = t.repos.AffBlockSection.DeleteByCampaignID(campaignID)
	if err != nil {
		return err
	}

	var records []*model.AffBlockSectionModel
	for _, sectionID := range sectionIDs {
		records = append(records, &model.AffBlockSectionModel{
			TrafficSource: "outbrain",
			CampaignID:    campaignID,
			SectionID:     sectionID,
		})
	}
	err = t.repos.AffBlockSection.SaveSlice(records)
	if err != nil {
		return err
	}
	return
}

func (t *reportAffUsecase) GetListSectionBlocked(trafficSource, campaignID string) (listSectionBlocked []string, err error) {
	if campaignID == "" {
		return
	}
	query := make(map[string]interface{})
	query["traffic_source"] = trafficSource
	query["campaign_id"] = campaignID
	records, err := t.repos.AffBlockSection.FindAllByQuery(query)
	if err != nil {
		return
	}
	for _, record := range records {
		listSectionBlocked = append(listSectionBlocked, record.SectionID)
	}
	return
}

func (t *reportAffUsecase) CheckSectionBlocked(trafficSource, campaignID, sectionID string) (isBlocked bool, err error) {
	if trafficSource == "" || campaignID == "" || sectionID == "" {
		return
	}
	query := make(map[string]interface{})
	query["traffic_source"] = trafficSource
	query["campaign_id"] = campaignID
	query["section_id"] = sectionID
	record, err := t.repos.AffBlockSection.FindOneByQuery(query)
	if err != nil {
		// logger.Error(err.Error())
		return
	}
	if record.IsFound() {
		return true, nil
	}
	return
}

func (t *reportAffUsecase) GetReportTotalForBlock(campaignID string, date string) (record *model.ReportAffModel, err error) {
	return t.repos.ReportAff.FindReportTotalForBlock(campaignID, date)
}

func (t *reportAffUsecase) GetForReportAllTrafficSourceIDs(trafficSource, account string) (IDs []string, err error) {
	IDs, err = t.repos.ReportAffPixel.FindForReportAllTrafficSourceIDs(trafficSource, account)
	if err != nil {
		return
	}
	records, _ := t.repos.Campaign.FindAllTrafficSourceID(trafficSource, account)
	if len(records) > 0 {
		for _, record := range records {
			if record.TrafficSourceID != "" && !utility.InArray(record.TrafficSourceID, IDs, true) {
				IDs = append(IDs, record.TrafficSourceID)
			}
		}
	}
	return
}

func (t *reportAffUsecase) updateDataFromReportMgid(dates []string) (err error) {
	fmt.Println("====== Update Data From Mgid For Report Aff ======")

	// Tạo reportAff để update hoặc create
	var recordReportAffs []*model.ReportAffModel

	for _, date := range dates {
		reportMgids, _ := t.repos.ReportMgid.FindByDayForReportAff(date)
		for _, reportMgid := range reportMgids {
			campaign, err := t.repos.Campaign.GetByTrafficSourceID(reportMgid.CampaignID)
			if err != nil {
				continue
			}
			query := make(map[string]interface{})
			query["date"] = date
			query["traffic_source"] = "mgid"
			query["partner"] = strings.ToLower(campaign.DemandSource)
			//query["campaign_id"] = reportMgid.CampaignID
			query["redirect_id"] = campaign.ID
			query["section_id"] = reportMgid.SectionID
			query["style_id"] = ""
			query["layout_id"] = ""
			query["layout_version"] = ""
			query["device"] = ""
			query["geo"] = ""
			recordReportAff, _ := t.repos.ReportAff.FindOneByQuery(query)
			if recordReportAff.ID != 0 { // Nếu đã tồn tại thì update các value cần
				// Từ report Mgid có thể update được Cost(spend), Supply Click(clicks), Supply Conversion(conversions) cho reportAff
				recordReportAff.Cost = reportMgid.Spent
				recordReportAff.SupplyClick = reportMgid.Clicks
				//recordReportAff.SupplyConversions = reportMgid.Conversions

			} else { // / Nếu chưa tồn tại tạo mới 1 record
				recordReportAff = &model.ReportAffModel{
					UserID:        campaign.UserID,
					Date:          date,
					TrafficSource: "mgid",
					SectionId:     reportMgid.SectionID,
					CampaignID:    reportMgid.CampaignID,
					CampaignName:  campaign.Name,
					RedirectID:    campaign.ID,
					StyleID:       "",
					Partner:       strings.ToLower(campaign.DemandSource),
					Cost:          reportMgid.Spent,
					SupplyClick:   reportMgid.Clicks,
				}
			}
			recordReportAffs = append(recordReportAffs, recordReportAff)
		}
	}

	// Sau khi xử lý xong update lại các records
	for _, recordReportAff := range recordReportAffs {
		err = t.repos.ReportAff.Save(recordReportAff)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
	}
	return
}

func (t *reportAffUsecase) updateDataFromReportPocPoc(dates []string) (err error) {
	fmt.Println("====== Update Data From PocPoc For Report Aff ======")

	// Tạo reportAff để update hoặc create
	var recordReportAffs []*model.ReportAffModel

	for _, date := range dates {
		reportPocPocs, _ := t.repos.ReportPocPoc.FindByDayForReportAff(date)
		for _, reportPocPoc := range reportPocPocs {
			campaign, err := t.repos.Campaign.GetByTrafficSourceID(reportPocPoc.CampaignID)
			if err != nil {
				continue
			}
			if campaign.ID == 0 {
				continue
			}
			query := make(map[string]interface{})
			query["date"] = date
			query["traffic_source"] = "pocpoc"
			query["partner"] = strings.ToLower(campaign.DemandSource)
			//query["campaign_id"] = reportPocPoc.CampaignID
			query["redirect_id"] = campaign.ID
			query["section_id"] = reportPocPoc.SectionID
			query["style_id"] = ""
			query["layout_id"] = ""
			query["layout_version"] = ""
			query["device"] = ""
			query["geo"] = ""
			recordReportAff, _ := t.repos.ReportAff.FindOneByQuery(query)
			if recordReportAff.ID != 0 { // Nếu đã tồn tại thì update các value cần
				// Từ report Mgid có thể update được Cost(spend), Supply Click(clicks), Supply Conversion(conversions) cho reportAff
				recordReportAff.Cost = reportPocPoc.Spent
				recordReportAff.SupplyClick = reportPocPoc.Clicks
				if campaign.Bidding == "cpm" {
					recordReportAff.CostCPC = reportPocPoc.Spent
				}
				//recordReportAff.SupplyConversions = reportPocPoc.Conversions

			} else { // / Nếu chưa tồn tại tạo mới 1 record
				costCPC := 0.0
				if campaign.Bidding == "cpm" {
					costCPC = reportPocPoc.Spent
				}
				recordReportAff = &model.ReportAffModel{
					UserID:        campaign.UserID,
					Date:          date,
					TrafficSource: "pocpoc",
					SectionId:     reportPocPoc.SectionID,
					CampaignID:    reportPocPoc.CampaignID,
					CampaignName:  campaign.Name,
					RedirectID:    campaign.ID,
					StyleID:       "",
					Partner:       strings.ToLower(campaign.DemandSource),
					Cost:          reportPocPoc.Spent,
					CostCPC:       costCPC,
					SupplyClick:   reportPocPoc.Clicks,
				}
			}
			recordReportAffs = append(recordReportAffs, recordReportAff)
		}
	}

	// Sau khi xử lý xong update lại các records
	for _, recordReportAff := range recordReportAffs {
		err = t.repos.ReportAff.Save(recordReportAff)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
	}
	return
}

func (t *reportAffUsecase) updateDataFromReportTikTok(dates []string) (err error) {
	fmt.Println("====== Update Data From TikTok For Report Aff ======")
	// Tạo reportAff để update hoặc create
	mapReportTikTok := make(map[string]*model.ReportTikTokModel)
	for _, date := range dates {
		reportTikToks, _ := t.repos.ReportTikTok.FindByDayForReportAff(date)
		for _, reportTikTok := range reportTikToks {
			var campaignID string
			campaignID = reportTikTok.CampaignID
			key := date + "_" + "tiktok" + "_" + campaignID + "_" + reportTikTok.AdGroupID
			if val, ok := mapReportTikTok[key]; ok {
				val.Spend += reportTikTok.Spend
				val.Clicks += reportTikTok.Clicks
				val.Impressions += reportTikTok.Impressions
			} else {
				mapReportTikTok[key] = reportTikTok
			}
		}
	}
	var recordReportAffs []*model.ReportAffModel
	for key, reportTikTok := range mapReportTikTok {
		splitKey := strings.Split(key, "_")
		date := splitKey[0]
		trafficSource := splitKey[1]
		campaignID := splitKey[2]
		adGroupID := splitKey[3]

		var campaign model.CampaignModel
		var recordReportAff *model.ReportAffModel
		campaign = t.repos.Campaign.GetById(reportTikTok.RedirectID, true)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		if campaign.ID == 0 {
			logger.Error("campaign null")
			continue
		}

		query := make(map[string]interface{})
		query["date"] = date
		query["traffic_source"] = trafficSource
		query["partner"] = strings.ToLower(campaign.DemandSource)
		//query["campaign_id"] = reportPocPoc.CampaignID
		query["redirect_id"] = campaign.ID
		query["section_id"] = adGroupID
		query["style_id"] = ""
		query["layout_id"] = ""
		query["layout_version"] = ""
		query["device"] = ""
		query["geo"] = ""
		recordReportAff, _ = t.repos.ReportAff.FindOneByQuery(query)

		if recordReportAff.ID != 0 { // Nếu đã tồn tại thì update các value cần
			// Từ report OutBrain có thể update được Cost(spend), Supply Click(clicks), Supply Conversion(conversions) cho reportAff
			recordReportAff.Cost = reportTikTok.Spend
			recordReportAff.CostCPC = reportTikTok.Spend
			recordReportAff.SupplyClick = reportTikTok.Clicks
			//recordReportAff.SupplyConversions = reportTikTok.CpaActionsNum
		} else { // / Nếu chưa tồn tại tạo mới 1 record
			recordReportAff = &model.ReportAffModel{
				Date:          date,
				UserID:        campaign.UserID,
				TrafficSource: trafficSource,
				SectionId:     adGroupID,
				//SectionName:   reportTikTok.Site,
				CampaignID:   campaignID,
				CampaignName: campaign.Name,
				StyleID:      "",
				Cost:         reportTikTok.Spend,
				CostCPC:      reportTikTok.Spend,
				SupplyClick:  reportTikTok.Clicks,
				//SupplyConversions: reportTikTok.CpaActionsNum,
				RedirectID: campaign.ID,
				Partner:    strings.ToLower(campaign.DemandSource),
			}
		}
		recordReportAffs = append(recordReportAffs, recordReportAff)
	}
	// Sau khi xử lý xong update lại các records
	for _, recordReportAff := range recordReportAffs {
		err = t.repos.ReportAff.Save(recordReportAff)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
	}
	return
}

func (t *reportAffUsecase) UpdateDataFromReportFacebook(dates []string) (err error) {
	fmt.Println("====== Update Data From TikTok For Report Aff ======")
	// Tạo reportAff để update hoặc create
	mapReportFacebook := make(map[string]*model.ReportFacebookModel)
	for _, date := range dates {
		reportFacebooks, _ := t.repos.ReportFacebook.FindByDayForReportAff(date)
		for _, reportFacebook := range reportFacebooks {
			var campaignID string
			campaignID = reportFacebook.CampaignID
			key := date + "_" + "facebook" + "_" + campaignID
			if val, ok := mapReportFacebook[key]; ok {
				val.Spend += reportFacebook.Spend
				val.Clicks += reportFacebook.Clicks
				val.Impressions += reportFacebook.Impressions
			} else {
				mapReportFacebook[key] = reportFacebook
			}
		}
	}
	var recordReportAffs []*model.ReportAffModel
	for key, reportFacebook := range mapReportFacebook {
		splitKey := strings.Split(key, "_")
		date := splitKey[0]
		trafficSource := splitKey[1]
		campaignID := splitKey[2]

		var campaign model.CampaignModel
		var recordReportAff *model.ReportAffModel
		campaign = t.repos.Campaign.GetById(reportFacebook.RedirectID, true)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		if campaign.ID == 0 {
			logger.Error("campaign null")
			continue
		}

		query := make(map[string]interface{})
		query["date"] = date
		query["traffic_source"] = trafficSource
		query["partner"] = strings.ToLower(campaign.DemandSource)
		query["redirect_id"] = campaign.ID
		query["section_id"] = ""
		query["style_id"] = ""
		query["layout_id"] = ""
		query["layout_version"] = ""
		query["device"] = ""
		query["geo"] = ""
		recordReportAff, _ = t.repos.ReportAff.FindOneByQuery(query)

		if recordReportAff.ID != 0 { // Nếu đã tồn tại thì update các value cần
			// Từ report OutBrain có thể update được Cost(spend), Supply Click(clicks), Supply Conversion(conversions) cho reportAff
			recordReportAff.Cost = reportFacebook.Spend
			recordReportAff.CostCPC = reportFacebook.Spend
			recordReportAff.SupplyClick = reportFacebook.Clicks
			//recordReportAff.SupplyConversions = reportFacebook.CpaActionsNum
		} else { // / Nếu chưa tồn tại tạo mới 1 record
			recordReportAff = &model.ReportAffModel{
				Date:          date,
				UserID:        campaign.UserID,
				TrafficSource: trafficSource,
				//SectionId:     reportFacebook.SiteID,
				//SectionName:   reportFacebook.Site,
				CampaignID:   campaignID,
				CampaignName: campaign.Name,
				StyleID:      "",
				Cost:         reportFacebook.Spend,
				CostCPC:      reportFacebook.Spend,
				SupplyClick:  reportFacebook.Clicks,
				//SupplyConversions: reportFacebook.CpaActionsNum,
				RedirectID: campaign.ID,
				Partner:    strings.ToLower(campaign.DemandSource),
			}
		}
		recordReportAffs = append(recordReportAffs, recordReportAff)
	}
	// Sau khi xử lý xong update lại các records
	for _, recordReportAff := range recordReportAffs {
		err = t.repos.ReportAff.Save(recordReportAff)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
	}
	return
}

func (t *reportAffUsecase) UpdateDataFromReportDomainActive(dates []string) (err error) {
	fmt.Println("====== Update Data From Domain Active For Report Aff ======")
	// Tạo reportAff để update hoặc create
	mapReportDomainActive := make(map[string]*model.ReportDomainActiveModel)
	for _, date := range dates {
		reportDomainActives, _ := t.repos.ReportDomainActive.FindByDayForReportAff(date)
		for _, reportDomainActive := range reportDomainActives {
			key := date + "_" + strconv.FormatInt(reportDomainActive.RedirectID, 10) + "_" + reportDomainActive.SectionID
			if val, ok := mapReportDomainActive[key]; ok {
				val.RevenueClicks += reportDomainActive.RevenueClicks
				val.PublisherRevenueAmount += reportDomainActive.PublisherRevenueAmount
				val.TotalVisitors += reportDomainActive.TotalVisitors
				val.TrackedVisitors += reportDomainActive.TrackedVisitors
			} else {
				mapReportDomainActive[key] = reportDomainActive
			}
		}
	}
	var recordReportAffs []*model.ReportAffModel
	for key, reportDomainActive := range mapReportDomainActive {
		splitKey := strings.Split(key, "_")
		date := splitKey[0]

		var campaign model.CampaignModel
		var recordReportAff *model.ReportAffModel
		campaign = t.repos.Campaign.GetById(reportDomainActive.RedirectID, true)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		if campaign.ID == 0 {
			logger.Error("campaign null")
			continue
		}

		query := make(map[string]interface{})
		query["date"] = date
		query["traffic_source"] = strings.ToLower(campaign.TrafficSource)
		query["partner"] = strings.ToLower(campaign.DemandSource)
		query["redirect_id"] = campaign.ID
		query["section_id"] = reportDomainActive.SectionID
		query["style_id"] = ""
		query["layout_id"] = ""
		query["layout_version"] = ""
		query["device"] = ""
		query["geo"] = ""
		recordReportAff, _ = t.repos.ReportAff.FindOneByQuery(query)
		if recordReportAff.ID != 0 { // Nếu đã tồn tại thì update các value cần
			recordReportAff.Revenue = reportDomainActive.PublisherRevenueAmount
			recordReportAff.ClickAdsense = reportDomainActive.RevenueClicks
		} else { // / Nếu chưa tồn tại tạo mới 1 record
			recordReportAff = &model.ReportAffModel{
				Date:          date,
				UserID:        campaign.UserID,
				TrafficSource: strings.ToLower(campaign.TrafficSource),
				Partner:       strings.ToLower(campaign.DemandSource),
				SectionId:     reportDomainActive.SectionID,
				CampaignID:    campaign.TrafficSourceID,
				CampaignName:  campaign.Name,
				Revenue:       reportDomainActive.PublisherRevenueAmount,
				ClickAdsense:  reportDomainActive.RevenueClicks,
				RedirectID:    campaign.ID,
			}
		}
		recordReportAffs = append(recordReportAffs, recordReportAff)
	}
	// Sau khi xử lý xong update lại các records
	for _, recordReportAff := range recordReportAffs {
		err = t.repos.ReportAff.Save(recordReportAff)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
	}
	return
}
