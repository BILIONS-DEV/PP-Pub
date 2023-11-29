package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"source/apps/frontend/lang"
	"source/apps/frontend/payload"
	//"source/core/technology/mysql"
	"time"
)

type AdBlock struct{}

type AdBlockRecord struct {
	//mysql.TableAdBlock
}

//func (AdBlockRecord) TableName() string {
//	return mysql.Tables.AdblockAnalytics
//}

type AdBlockRecordDatatable struct {
	AdBlockRecord
	Domain string `json:"domain"`
}

type ResponseFilters struct {
	Found int    `json:"total"`
	Date  string `json:"__time"`
}

func (t *AdBlock) GetByFilters(inputs *payload.AdBlockFilterPayload, userID int64, lang lang.Translation) (adblocks []ResponseFilters, err error) {
	var domains []InventoryRecord
	var domainsUUID string
	if inputs.InventoryId != 0 {
		var domain InventoryRecord
		domain, err = new(Inventory).GetByIdSystem(inputs.InventoryId)
		domainsUUID = domain.Uuid
	} else {
		domains = new(Inventory).GetByUser(userID)
		for _, value := range domains {
			if domainsUUID == "" {
				domainsUUID = "'" + value.Uuid + "'"
			} else {
				domainsUUID = domainsUUID + ",'" + value.Uuid + "'"
			}
		}
	}
	if domainsUUID == "" {
		return
	}
	//domainsUUID = "'15a7c41b3c16f0e37e72cee0c13ea9cc'"

	url := "https://query.vliplatform.com/druid/v2/sql/"
	client := &http.Client{}

	params := make(map[string]interface{})
	params["query"] = "SELECT SUM(totalRequest) as total, TIME_FORMAT(FLOOR(__time TO DAY), 'yyyy-MM-dd') as __time FROM \"vli-adblock-analytics\" WHERE FLOOR(__time TO DAY) >= TIMESTAMP '" + inputs.StartDate + "' AND FLOOR(__time TO DAY) <= TIMESTAMP '" + inputs.EndDate + "' AND detectType = 'found' AND domainUUID in (" + domainsUUID + ") group by FLOOR(__time TO DAY) ORDER BY __time desc"
	// params["context"] = map[string]string{
	// 	"sqlTimeZone": "America/New_York",
	// }
	//fmt.Printf("%+v\n", params["query"])

	payload, _ := json.Marshal(params)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic cmVhZG1hbjpMb3hpVVNKISEoODIhQA==")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.Unmarshal(body, &adblocks)
	return
}

func (t *AdBlock) setFilterCondition(inputs *payload.AdBlockFilterPayload, userID int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if inputs.InventoryId != 0 {
			db.Where("inventory_id  = ?", inputs.InventoryId)
		}
		return db.Where("date >= ? AND date <= ? AND user_id = ?", inputs.StartDate, inputs.EndDate, userID)
	}
}

func (t *AdBlock) GetPreriod(startDate, endDate string) (dates []string) {
	// Chuyển đổi chuỗi ngày thành kiểu time.Time
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Tạo một slice chứa các ngày trong khoảng thời gian từ startDate đến endDate
	times := make([]time.Time, 0)
	for current := start; !current.After(end); current = current.AddDate(0, 0, 1) {
		times = append(times, current)
	}

	// In ra các ngày trong slice
	for _, date := range times {
		dates = append(dates, date.Format("2006-01-02"))
	}
	return
}
