package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log"
	"source/config"
	"source/infrastructure/caching"
	"source/infrastructure/kafka"
	"source/infrastructure/mysql"
	"source/internal/entity/dto"
	"source/internal/entity/model"
	"source/internal/repo"
	"source/internal/usecase"
	reportBodisUC "source/internal/usecase/report-bodis"
	"source/pkg/logger"
	"source/pkg/utility"
	"strconv"
	"strings"
	"sync"
	"time"
)

type tonic struct {
	token   string
	expires int64
}

const (
	AccountTaboola1 = "account_1"
	AccountTaboola2 = "account_2"
)

var listAccountTaboola []string

type taboola struct {
	clientID     string
	clientSecret string
	token        string
	expires      int64
}

type handler struct {
	UseCases *usecase.UseCases
	tonic    tonic
	taboola  map[string]taboola
}

func (t *handler) initTaboola() {
	listAccountTaboola = append(listAccountTaboola, AccountTaboola1)
	listAccountTaboola = append(listAccountTaboola, AccountTaboola2)
	mapTaboola := make(map[string]taboola)
	mapTaboola[AccountTaboola1] = taboola{
		clientID:     "8411367616a8424c92a6010e5427a778",
		clientSecret: "dd3a45b631434ace905b9df599da2a02",
	}
	mapTaboola[AccountTaboola2] = taboola{
		clientID:     "98135e1de19940c8ae0d55dfa74aeaa4",
		clientSecret: "35db2d0be8ef44bb9bbec8a0e315e48e",
	}
	t.taboola = mapTaboola
}

func main() {
	var (
		DB    *gorm.DB
		Cache caching.Cache
		Kafka *kafka.Client
		err   error
	)
	configs := config.NewConfig()
	configs.Mysql.Database = "apac_aff"
	// => init Mysql
	if configs.Mysql != nil {
		fmt.Println("init mysql....", true)
		host := configs.Mysql.Host
		if !utility.IsWindow() {
			host = "192.168.9.13"
			configs.Mysql.Username = "aff_admin"
			configs.Mysql.Password = "PO()(@032KK!NMA!@a"
			configs.Mysql.Port = "15869"
		}
		if DB, err = mysql.Connect(mysql.Config{
			Username: configs.Mysql.Username,
			Password: configs.Mysql.Password,
			Host:     host,
			Port:     configs.Mysql.Port,
			Database: configs.Mysql.Database,
			Encoding: configs.Mysql.Encoding,
		}); err != nil {
			log.Fatalln("Error connect mysql: ", err)
			return
		}
	}

	// => init Aerospike
	if configs.Aerospike != nil {
		fmt.Println("init aerospike...", true)
		if Cache, err = caching.NewAerospike(caching.Config{
			Host:      configs.Aerospike.Host,
			Port:      configs.Aerospike.Port,
			Namespace: configs.Aerospike.Namespace,
		}); err != nil {
			log.Fatalln("Error connect Cache: ", err)
			return
		}
	}

	// => init Kafka
	if configs.Kafka != nil {
		fmt.Println("init kafka...", true)
		var brokers []string
		if utility.IsWindow() {
			brokers = []string{"127.0.0.1:9092"}
		} else {
			brokers = []string{"192.168.9.11:9092"}
		}
		Kafka = kafka.NewKafka(kafka.Config{
			Brokers: brokers,
			Topic:   "pw-report-tracking-ads",
		})
		if err != nil {
			log.Fatalln("Error connect Kafka: ", err)
			return
		}
	}

	// => init Repository
	repos := repo.NewRepositories(&repo.Deps{Db: DB, Cache: Cache, Kafka: Kafka})
	// repos.ReportOutBrain.Migrate()
	// repos.ReportBodisTraffic.Migrate()
	// repos.ReportOpenMailSubID.Migrate()
	// repos.ReportAff.Migrate()
	// repos.ReportAffSelection.Migrate()
	// repos.ReportAffPixel.Migrate()
	// repos.ReportTonic.Migrate()
	// repos.ReportAdsense.Migrate()
	// repos.ReportTaboola.Migrate()
	// return

	// => init UseCases
	useCases := usecase.NewUseCases(&usecase.Deps{Repos: repos})
	t := &handler{
		UseCases: useCases,
		tonic: tonic{
			token:   "",
			expires: 0,
		},
	}

	// => init taboola
	t.initTaboola()
	// Lấy tất cả report của Aff của ngày hôm nay và hôm qua

	// Xử lý và get report cho pixel
	go func() {
		for {
			t.startGetReportPixelAff()
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		for {
			var dates []string
			layout := "2006-01-02"
			today := time.Now().UTC()
			dates = append(dates, today.Format(layout))
			yesterday := today.AddDate(0, 0, -1)
			dates = append(dates, yesterday.Format(layout))
			timeStartUpdate := time.Now().UTC()
			err := t.UseCases.ReportAff.UpdateDataPixelForReportAff(dates)
			if err != nil {
				logger.Error(err.Error())
			}
			fmt.Println("Done Update From Pixel --Start: ", timeStartUpdate, " --End: ", time.Now().UTC(), "Total Time:", time.Now().UTC().Sub(timeStartUpdate).Milliseconds(), "ms")
			time.Sleep(60 * time.Second)
		}
	}()

	// Xử lý và get report cho outBrain
	go func() {
		for {
			t.startGetReportOutBrain()
			time.Sleep(120 * time.Second)
		}
	}()

	//Xử lý và get report cho taboola
	go func() {
		for {
			t.startGetReportTaboola()
			time.Sleep(120 * time.Second)
		}
	}()

	//Xử lý và get report cho mgid
	go func() {
		for {
			t.startGetReportMgid()
			time.Sleep(120 * time.Second)
		}
	}()

	//Xử lý và get report cho pocpoc
	go func() {
		for {
			t.startGetReportPocPoc()
			time.Sleep(120 * time.Second)
		}
	}()

	//Xử lý và get và update report cho codefuel
	go func() {
		for {
			// Get Report
			t.startGetReportCodeFuel()

			// Update report CodeFuel vào report aff
			layout := "2006-01-02"
			endDate := time.Now().UTC().Format(layout)
			startDate := time.Now().UTC().AddDate(0, 0, -2).Format(layout)
			// Update Report
			err = t.UseCases.ReportCodeFuel.UpdateReportAFF(startDate, endDate)
			if err != nil {
				logger.Error(err.Error())
			}
			time.Sleep(5 * time.Minute)
		}
	}()

	//Xử lý và get report cho adsense
	go func() {
		for {
			var dates []string
			layout := "2006-01-02"
			today := time.Now().UTC()
			dates = append(dates, today.Format(layout))
			yesterday := today.AddDate(0, 0, -1)
			dates = append(dates, yesterday.Format(layout))
			t.startGetReportAdsense()
			timeStartUpdate := time.Now().UTC()
			err = t.UseCases.ReportAff.UpdateDataFromReportAdsense(dates)
			if err != nil {
				fmt.Println("Err Update Adsense: ", err.Error())
			}
			fmt.Println("Done Update From Adsense --Start: ", timeStartUpdate, " --End: ", time.Now().UTC(), "Total Time:", time.Now().UTC().Sub(timeStartUpdate).Milliseconds(), "ms")
			time.Sleep(120 * time.Second)
		}
	}()

	// Xử lý tổng hợp các data cho report_aff
	t.updateDataForReportAff()

	for {

		var wg sync.WaitGroup
		// Xử lý goroutine
		wg.Add(2) // Số lượng đợi goroutine
		// Xử lý và get report cho bodis
		// t.startGetReportBodis() // Bỏ report Bodis bị bem

		// Xử lý và get report cho system1
		go func() {
			t.startGetReportOpenMail()
			wg.Done()
		}()

		// Xử lý và get report cho tonic
		go func() {
			t.startGetReportTonic()
			wg.Done()
		}()

		wg.Wait()

		if !utility.IsWindow() {
			// Xử lý tìm các section không hiệu quả block đi
			//t.blockSection()
		}

		fmt.Println("End round")
		time.Sleep(120 * time.Second)
	}
}

func (t *handler) startGetReportBodis() {
	fmt.Println("====== Start Get Report Bodis ======")
	// => Get
	err := t.UseCases.ReportBodis.Cronjob()
	if err != nil {
		fmt.Println(err)
	}

	// Date là mảng các ngày lấy report
	var dates []string
	// dates = append(dates, "2022-08-11")
	date := time.Now().UTC().AddDate(0, 0, -1)
	dates = append(dates, date.Format("2006-01-02"))
	date = time.Now().UTC()
	dates = append(dates, date.Format("2006-01-02"))
	fmt.Println(dates)
	total := 0.0
	// var listVisit []string
	for _, date := range dates {
		// Inputs để đưa vào usecase xử lý
		var inputs []reportBodisUC.InputHandlerReportTraffic
		page := 0
		for {
			page++
			// => Get
			var from, to string
			if date == time.Now().UTC().Format("2006-01-02") {
				from = date + " 00:00:00"
				to = time.Now().Add(-3660 * time.Second).UTC().Format("2006-01-02 15:04:05") // Report realtime chậm hơn 1 tiếng
			} else {
				from = date + " 00:00:00"
				to = date + " 23:59:59"
			}
			output, err := t.UseCases.ReportBodis.GetReportTraffic(from, to, strconv.Itoa(page))
			if err != nil {
				fmt.Println(err)
			}
			// => Parse to struct
			var response dto.ReportBodisParkingSearch
			_ = json.Unmarshal(output, &response)
			for _, data := range response.Data {
				// fmt.Println(data)
				// if data.EstimatedRevenue < 0 {
				//	fmt.Printf("%+v \n", data)
				//
				// }
				total += data.EstimatedRevenue

				// if utility.InArray(data.VisitID, listVisit, true) {
				//	fmt.Println("đã có visit:", data.VisitID)
				//	continue
				// } else {
				//	listVisit = append(listVisit, data.VisitID)
				// }
				// => Tạo input đưa vào usecase để xử lý
				inputs = append(inputs, reportBodisUC.InputHandlerReportTraffic{
					Time:             strings.Replace(data.ServerDatetime, " ", "T", 1),
					VisitID:          data.VisitID,
					DomainName:       data.DomainName,
					IpAddress:        data.IPAddress,
					Type:             data.Type,
					CountryID:        data.CountryID,
					PageQuery:        data.PageQuery,
					SubIds:           data.Subids,
					Clicks:           data.Clicks,
					EstimatedRevenue: data.EstimatedRevenue,
				})
			}
			if response.NextPageURL == nil {
				break
			}
		}
		err = t.UseCases.ReportBodis.HandlerReportTraffic(inputs)
	}
	fmt.Println(total)
	return
}

func (t *handler) startGetReportOpenMail() {
	fmt.Println("====== Start Get Report OpenMail ======")
	// Date là mảng các ngày lấy report
	var dates []string
	// dates = append(dates, "2022-08-04")
	date := time.Now().UTC().AddDate(0, 0, -1)
	dates = append(dates, date.Format("2006-01-02"))
	date = time.Now().UTC()
	dates = append(dates, date.Format("2006-01-02"))

	// Get reports
	query := "days=" + strings.Join(dates, ",")
	fmt.Println(query)
	records, err := t.UseCases.ReportOpenMail.GetReportHourly(query)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Truyền records vào usecase xử lý report
	err = t.UseCases.ReportOpenMail.HandlerReportSubID(records)
}

type SectionByHour struct {
	Breakdown       string      `json:"breakdown"`
	Dimensions      []Dimension `json:"dimensions"`
	TotalDimensions int         `json:"totalDimensions"`
}

type Dimension struct {
	ID          string       `json:"id"`
	DataPoints  []DataPoints `json:"dataPoints"`
	TotalPoints int          `json:"totalPoints"`
}

type DataPoints struct {
	Timestamp      string  `json:"timestamp"`
	Impressions    int64   `json:"impressions"`
	Clicks         int64   `json:"clicks"`
	Spend          float64 `json:"spend"`
	Ctr            float64 `json:"ctr"`
	Ecpc           float64 `json:"ecpc"`
	Conversions    int64   `json:"conversions"`
	ConversionRate float64 `json:"conversionRate"`
	Cpa            float64 `json:"cpa"`
}

func (t *handler) startGetReportOutBrain() {
	startTime := time.Now()
	fmt.Println("====== Start Get Report OutBrain --Time:", startTime, " ======")
	// => Get ra toàn bộ các marketer của report
	marketerIDS, err := t.UseCases.ReportOutBrain.GetMarketerIDs()
	if err != nil {
		fmt.Println(err)
		return
	}

	// => Từ marketerID xử lý lấy data cho mỗi marketer
	// limit := make(chan bool, len(marketerIDS))
	for _, marketerID := range marketerIDS {
		// limit <- true
		// go func(marketerID string) {
		// Khi return khỏi func truyền vào để chanel biết func hoàn thành
		// defer func() {
		//	<-limit
		// }()
		// Lấy ra toàn bộ Conversions của Marketer từ report
		conversionIDs, _ := t.UseCases.ReportOutBrain.GetConversionIDs(marketerID)

		// Lấy ra toàn bộ Campaign của Marketer từ report
		campaignIDs, err := t.UseCases.ReportAff.GetForReportAllTrafficSourceIDs("outbrain", marketerID)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// fmt.Println(campaigns)
		var records []model.ReportOutBrainModel
		for _, campaignID := range campaignIDs {
			// Từ campaignID truyền vào và lấy ra report section
			output, err := t.UseCases.ReportOutBrain.GetSectionByHour(marketerID, campaignID, conversionIDs)
			if err != nil {
				logger.Error(err.Error())
				continue
			}

			// Từ output unmarshal vào struct SectionByHour
			var sectionByHour SectionByHour
			err = json.Unmarshal(output, &sectionByHour)
			if err != nil {
				logger.Error(err.Error())
				continue
			}
			// fmt.Printf("%+v \n", sectionByHour)
			// Build ra inputs cho usecase HandlerReport
			for _, dimension := range sectionByHour.Dimensions {
				for _, dataPoint := range dimension.DataPoints {
					// Chỉ cần lấy data có spend
					if dataPoint.Spend > 0.0 {
						loc, _ := time.LoadLocation("America/New_York")
						layout := "2006-01-02 15:04:05"
						date, err := time.ParseInLocation(layout, dataPoint.Timestamp, loc)
						if err != nil {
							fmt.Println(err)
							return
						}
						layout = "2006-01-02 15:04:05"
						records = append(records, model.ReportOutBrainModel{
							MarketerID:  marketerID,
							CampaignID:  campaignID,
							SectionID:   dimension.ID,
							Time:        date.UTC().Format(layout),
							Spend:       dataPoint.Spend,
							Clicks:      dataPoint.Clicks,
							Conversions: dataPoint.Conversions,
						})
					}
				}
			}
		}

		// Từ inputs tiến hành xử lý report
		if err := t.UseCases.ReportOutBrain.HandlerReport(records); err != nil {
			logger.Error(err.Error())
			continue
		}
		// }(marketerID)
	}

	fmt.Println("====== End Get Report OutBrain --Total:", time.Now().Sub(startTime).Milliseconds(), "ms ======")
}

func (t *handler) startGetReportPixelAff() {
	//fmt.Println("====== Start Get Report Pixel Aff ======")

	// Xử lý get data pixel tạo report
	err := t.UseCases.ReportAff.HandlerReportAffPixel()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	//fmt.Println("====== End Get Report Pixel Aff ======")
}

func (t *handler) updateDataForReportAff() {
	fmt.Println("====== Update Data For Report Aff ======")

	// Sau khi xử lý report xử lý update các data từ các report khác
	err := t.UseCases.ReportAff.UpdateDataReportAff()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func (t *handler) startGetReportTonic() {
	fmt.Println("====== Start Get Report Tonic ======")

	if (t.tonic.expires-time.Now().UTC().Unix()) < 500 || t.tonic.token == "" {
		fmt.Println("Refresh Token")
		token, expiresNew, err := t.UseCases.ReportTonic.GetToken()
		if err != nil {
			return
		}
		t.tonic.token = token
		t.tonic.expires = expiresNew
	}

	// Date là mảng các ngày lấy report
	var dates []string
	var date time.Time
	// dates = append(dates, "2022-08-11")
	date = time.Now().UTC().AddDate(0, 0, -3)
	dates = append(dates, date.Format("2006-01-02"))
	date = time.Now().UTC().AddDate(0, 0, -2)
	dates = append(dates, date.Format("2006-01-02"))
	date = time.Now().UTC().AddDate(0, 0, -1)
	dates = append(dates, date.Format("2006-01-02"))
	date = time.Now().UTC()
	dates = append(dates, date.Format("2006-01-02"))
	fmt.Println(dates)

	for _, date := range dates {
		reports, err := t.UseCases.ReportTonic.GetReportFinal(t.tonic.token, date)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = t.UseCases.ReportTonic.SaveReport(reports)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (t *handler) startGetReportAdsense() {
	fmt.Println("====== Start Get Report Adsense ======")
	// Date là mảng các ngày lấy report
	var startDate, endDate string
	var date time.Time
	date = time.Now().UTC().AddDate(0, 0, -2)
	startDate = date.Format("2006-01-02")
	date = time.Now().UTC()
	endDate = date.Format("2006-01-02")

	listAdsense, err := t.UseCases.ReportAdsense.GetAllAdsense()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, adsense := range listAdsense {
		err := t.UseCases.ReportAdsense.HandlerReportAdsenseSubDomain(adsense.Token, string(adsense.Object), startDate, endDate)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		//if adsense.Object.IsAdsense3() {
		//	err = t.UseCases.ReportAdsense.HandlerReportAdsenseChannel(adsense.Token, string(adsense.Object), startDate, endDate)
		//	if err != nil {
		//		logger.Error(err.Error())
		//		continue
		//	}
		//}
	}
	fmt.Println("====== End Get Report Adsense ======")
}

func (t *handler) startGetReportTaboola() {
	fmt.Println("====== Start Get Report Taboola ======")
	// Hiện tại có 2 tài khoản taboola nên cần xử lý lấy report cho nhiều tài khoản
	for _, accountTaboola := range listAccountTaboola {
		if val, ok := t.taboola[accountTaboola]; ok {
			if (val.expires-time.Now().UTC().Unix()) < 500 || val.token == "" {
				fmt.Println("Refresh Token")
				token, expiresNew, err := t.UseCases.ReportTaboola.GetToken(val.clientID, val.clientSecret)
				if err != nil {
					fmt.Println(err)
					continue
				}
				val.token = token
				val.expires = expiresNew
				t.taboola[accountTaboola] = val
			}

			// Date là mảng các ngày lấy report
			startDate := time.Now().UTC().AddDate(0, 0, -2)
			endDate := time.Now().UTC()

			// Get All Account của tài khoản
			listAccountID, err := t.UseCases.ReportTaboola.GetAllAccountID(val.token)
			if err != nil {
				fmt.Println(err)
			}

			// Từ list Account ID get Report Campaign Sumary của Taboola
			for _, accountID := range listAccountID {
				reports, err := t.UseCases.ReportTaboola.GetReportCampaignSummary(accountTaboola, val.token, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), accountID)
				if err != nil {
					fmt.Println(err)
					continue
				}
				size := 200
				var j int
				for i := 0; i < len(reports); i += size {
					j += size
					if j > len(reports) {
						j = len(reports)
					}
					err = t.UseCases.ReportTaboola.SaveReport(reports[i:j])
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
			}
		}
	}
}

func (t *handler) startGetReportMgid() {
	startTime := time.Now()
	fmt.Println("====== Start Get Report Mgid --Time:", startTime, " ======")

	// Lấy ra toàn bộ Campaign của Marketer từ report
	campaignIDs, err := t.UseCases.ReportAff.GetForReportAllTrafficSourceIDs("mgid", "")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Lấy report cho các ngày
	var dates []string
	dates = append(dates, time.Now().UTC().Format("2006-01-02"))
	dates = append(dates, time.Now().AddDate(0, 0, -1).UTC().Format("2006-01-02"))
	// fmt.Println(campaigns)
	for _, date := range dates {
		for _, campaignID := range campaignIDs {
			// Từ campaignID truyền vào và lấy ra report section
			records, err := t.UseCases.ReportMgid.GetReport(campaignID, date)
			if err != nil {
				continue
			}
			// Từ inputs tiến hành xử lý lưu lại report
			if err := t.UseCases.ReportMgid.HandlerReport(records); err != nil {
				logger.Error(err.Error())
				continue
			}
		}
	}

	fmt.Println("====== End Get Report Mgid --Total:", time.Now().Sub(startTime).Milliseconds(), "ms ======")

}

func (t *handler) startGetReportPocPoc() {
	startTime := time.Now().UTC()
	fmt.Println("====== Start Get Report PocPoc --Time:", startTime, " ======")

	// Lấy report cho các ngày
	endDate := time.Now().UTC().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -1).UTC().Format("2006-01-02")

	// Từ campaignID truyền vào và lấy ra report section
	records, err := t.UseCases.ReportPocPoc.GetReport(startDate, endDate)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	// Từ inputs tiến hành xử lý lưu lại report
	if err := t.UseCases.ReportPocPoc.HandlerReport(records); err != nil {
		logger.Error(err.Error())
		return
	}

	fmt.Println("====== End Get Report PocPoc --Total:", time.Now().Sub(startTime).Milliseconds(), "ms ======")

}

func (t *handler) startGetReportCodeFuel() {
	startTime := time.Now()
	fmt.Println("====== Start Get Report CodeFuel --Time:", startTime, " ======")

	records, err := t.UseCases.ReportCodeFuel.GetReport()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, record := range records {
		fmt.Printf("%+v \n", record)
	}
	//if err := t.UseCases.ReportPocPoc.HandlerReport(records); err != nil {
	//	logger.Error(err.Error())
	//	return
	//}

	fmt.Println("====== End Get Report Code Fuel --Total:", time.Now().Sub(startTime).Milliseconds(), "ms ======")

}
