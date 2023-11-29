package report_codefuel

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
	"os/exec"
	"sort"
	"source/internal/entity/model"
	"source/internal/repo"
	"strconv"
	"strings"
	"time"
)

type UseCaseReportCodeFuel interface {
	GetReport() (records []*model.ReportCodeFuelModel, err error)
	UpdateReportAFF(startDate, endDate string) (err error)
}

type codeFuelUC struct {
	repos    *repo.Repositories
	user     string
	password string
	host     string
	port     int
}

func NewReportCodeFuelUC(repos *repo.Repositories) *codeFuelUC {
	return &codeFuelUC{
		repos:    repos,
		user:     "interdogMediaLimited",
		password: "hzoS567zL95f",
		host:     "ftp.toolmx.com",
		port:     990,
	}
}

// Xử lý report Mgid
func (t *codeFuelUC) HandlerReport(records []*model.ReportCodeFuelModel) (err error) {
	//err = t.repos.ReportCodeFuel.SaveSlice(records)
	//if err != nil {
	//	logger.Error(err.Error())
	//	return
	//}
	return
}

func (t *codeFuelUC) UpdateReportAFF(startDate, endDate string) (err error) {
	startTime := time.Now()
	fmt.Println("====== Update Report CodeFuel --Time:", startTime, " ======")
	defer func() {
		if err != nil {
			fmt.Println("====== Error Update Report CodeFuel", "--Error:", err, "--Total:", time.Now().Sub(startTime).Milliseconds(), "ms ======")
		} else {
			fmt.Println("====== Done Update Report CodeFuel --Total:", time.Now().Sub(startTime).Milliseconds(), "ms ======")
		}
	}()
	// Từ startDate và endDate lấy ra report của codefuel trong khoảng thời gian này
	records, err := t.repos.ReportCodeFuel.FindByDayForReportAff(startDate, endDate)
	if err != nil {
		return err
	}

	var reportAffs []*model.ReportAffModel
	for _, report := range records {
		if report.Time == "" || report.CampaignID == 0 {
			continue
		}
		campaignModel := t.repos.Campaign.GetById(report.CampaignID, true)
		if campaignModel.ID == 0 {
			// Nếu không tìm thấy campaign thì bỏ qua
			continue
		}

		// Tạo query để tìm report tương ứng từ reportAff
		var sectionID = "unknown"
		query := make(map[string]interface{})
		query["date"] = report.Time
		query["traffic_source"] = strings.ToLower(campaignModel.TrafficSource)
		query["partner"] = "codefuel"
		//query["campaign_id"] = campaignModel.TrafficSourceID
		query["redirect_id"] = campaignModel.ID
		query["section_id"] = sectionID
		query["style_id"] = ""
		query["layout_id"] = ""
		query["layout_version"] = ""
		query["device"] = ""
		query["geo"] = ""

		recordReportAff, _ := t.repos.ReportAff.FindOneByQuery(query)
		if recordReportAff.ID != 0 { // Nếu đã tồn tại thì update các value cần
			recordReportAff.Impressions = report.TotalMonetizedSearches
			recordReportAff.ClickAdsense = report.AdClicks
			recordReportAff.Revenue = report.Amount
		} else { // / Nếu chưa tồn tại tạo mới 1 record
			recordReportAff = &model.ReportAffModel{
				UserID:        campaignModel.UserID,
				Date:          report.Time,
				TrafficSource: strings.ToLower(campaignModel.TrafficSource),
				Partner:       "codefuel",
				SectionId:     sectionID,
				CampaignID:    campaignModel.TrafficSourceID,
				RedirectID:    campaignModel.ID,
				StyleID:       "",
				Impressions:   report.TotalMonetizedSearches,
				ClickAdsense:  report.AdClicks,
				Revenue:       report.Amount,
			}
		}
		reportAffs = append(reportAffs, recordReportAff)
	}

	size := 500
	var j int
	for i := 0; i < len(records); i += size {
		j += size
		if j > len(records) {
			j = len(records)
		}
		err = t.repos.ReportAff.SaveSlice(reportAffs[i:j])
		if err != nil {
			continue
		}
	}

	return
}

// GetReport lấy report từ ftp về lưu lại csv và tạo model
func (t *codeFuelUC) GetReport() (records []*model.ReportCodeFuelModel, err error) {
	c := exec.Command("curl",
		"--list-only",
		fmt.Sprintf("ftps://%s:%s@%s:%d", t.user, t.password, t.host, t.port),
	)
	fmt.Println(c.String())
	out, err := c.Output()
	if err != nil {
		return
	}
	// Tạo ra một map để sắp xếp các file theo time
	mapFile := make(map[int64]string)
	var listTime []int64

	// Tách list row theo từng dòng xử lý tách lấy file cho từng dòng
	listNameFile := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, nameFile := range listNameFile {
		// Từ name tách ra các phần để lấy time
		if nameFile == "" {
			continue
		}

		splName1 := strings.Split(nameFile, ".")
		if len(splName1) == 0 {
			continue
		}
		// Xử lý từ tên file lấy ra time
		splName2 := strings.Split(splName1[0], "-")
		if len(splName2) < 3 {
			continue
		}
		splName3 := strings.Split(splName2[2], "_")
		if len(splName3) < 4 {
			continue
		}
		date, err := time.Parse("2006_01_02 15:04:05", splName2[1]+" "+strings.Join(splName3[0:3], ":"))
		if err != nil {
			continue
		}

		// Đặt file name và date vào map
		_, exist := mapFile[date.Unix()]
		if exist {
			date.Add(1 * time.Second)
		}
		mapFile[date.Unix()] = nameFile
		listTime = append(listTime, date.Unix())

	}
	// Sắp xếp lại listTime theo thứ tự tăng dần để lấy report chuẩn nhất
	sort.Slice(listTime, func(i, j int) bool { return listTime[i] < listTime[j] })
	//mapReport := make(map[string]*model.ReportCodeFuelModel)
	// Xử lý report cho từ file
	for _, key := range listTime {
		nameFile, exist := mapFile[key]
		if !exist {
			continue
		}
		_, err = t.handlerReportCodeFuel(nameFile)
		if err != nil {
			continue
		}
	}
	return
}

func (t *codeFuelUC) handlerReportCodeFuel(nameFile string) (records []*model.ReportCodeFuelModel, err error) {
	// Xử lý download file về
	c := exec.Command("curl",
		"-Lk", "--ftp-ssl",
		fmt.Sprintf("ftps://%s:%s@%s:%d/%s", t.user, t.password, t.host, t.port, nameFile),
		"-o",
		"./codefuel/"+nameFile,
	)
	fmt.Println(c.String())
	_, err = c.Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	type ReportCodeFuel struct {
		Date                   string  `csv:"Date"`
		SearchChannel          string  `csv:"Channel"`
		TotalMonetizedSearches int64   `csv:"Monetized Searches"`
		AdClicks               int64   `csv:"Ad Clicks"`
		Amount                 float64 `csv:"Amount"`
	}

	var reports []ReportCodeFuel
	var fileCsv *os.File
	fileCsv, err = os.OpenFile("./codefuel/"+nameFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return
	}
	defer fileCsv.Close()
	if err = gocsv.UnmarshalFile(fileCsv, &reports); err != nil { // Load clients from file
		fmt.Println(err)
		return
	}
	for _, report := range reports {
		if report.Date == "" {
			continue
		}
		splitDate := strings.Split(report.Date, "/")
		if len(splitDate) < 3 {
			continue
		}
		day, _ := strconv.Atoi(splitDate[1])
		month, _ := strconv.Atoi(splitDate[0])
		year, _ := strconv.Atoi(splitDate[2])
		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		redirectID, _ := strconv.ParseInt(report.SearchChannel, 10, 64)
		campaignModel := t.repos.Campaign.GetById(redirectID, true)
		if campaignModel.ID == 0 {
			// Nếu không tìm thấy campaign thì bỏ qua
			continue
		}
		// Tạo list model
		records = append(records, &model.ReportCodeFuelModel{
			Time:                   date,
			CampaignID:             campaignModel.ID,
			TotalMonetizedSearches: report.TotalMonetizedSearches,
			AdClicks:               report.AdClicks,
			Amount:                 report.Amount,
		})
	}

	// Save lại các record
	err = t.repos.ReportCodeFuel.SaveSlice(records)
	if err != nil {
		return
	}
	// Sau khi save lại xử lý xóa file trên ftp đi
	cmdDeleteFile := exec.Command("curl",
		"-v",
		"-u",
		fmt.Sprintf("%s:%s", t.user, t.password),
		fmt.Sprintf("ftps://%s:%s@%s:%d/%s", t.user, t.password, t.host, t.port, nameFile),
		"-Q",
		fmt.Sprintf("DELE %s", nameFile),
	)
	fmt.Println(cmdDeleteFile.String())
	_, _ = cmdDeleteFile.Output()
	//if err != nil {
	//	//logger.Error(err.Error())
	//	return
	//}
	return
}
