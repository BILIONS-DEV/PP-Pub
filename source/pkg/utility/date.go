package utility

import (
	"fmt"
	"strings"
	"time"
)

// param: s
// return:
func Date(s string) time.Time {
	var date time.Time
	var err error

	if strings.Index(s, "-") > -1 {
		date, err = time.Parse("2006-01-02", s)
	} else if strings.Index(s, "/") > -1 {
		date, err = time.Parse("2006/01/02", s)
	} else {
		date, err = time.Parse("20060102", s)
	}

	if err != nil {
		panic(err)
	}

	return date
}

// param: a, b
// return:
func Intervals(a, b time.Time) []time.Time {
	if a.After(b) {
		a, b = b, a
	}

	var days []time.Time
	for i := 0; i <= int(b.Unix()-a.Unix())/3600/24; i++ {
		day := a.Add((time.Hour * 24 * time.Duration(i)))

		day.Format("2006-01-02")
		days = append(days, day)
	}

	return days
}

func GetAllDates(startDate, endDate string) []time.Time {
	layout := "2006-01-02" // Định dạng ngày trong chuỗi

	startTime, _ := time.Parse(layout, startDate)
	endTime, _ := time.Parse(layout, endDate)

	var dates []time.Time

	// Thêm ngày bắt đầu vào danh sách
	dates = append(dates, startTime.UTC())

	// Tạo một biến để lưu trữ ngày hiện tại, khởi đầu bằng ngày bắt đầu
	currentDate := startTime.UTC()

	// Lặp qua từng ngày cho đến khi ngày hiện tại lớn hơn hoặc bằng ngày kết thúc
	for {
		// Thêm một ngày vào ngày hiện tại
		currentDate = currentDate.AddDate(0, 0, 1)

		// Thêm ngày hiện tại vào danh sách
		dates = append(dates, currentDate)

		// Nếu ngày hiện tại đã bằng ngày cuối thì dừng lại
		if currentDate.Format(layout) == endTime.Format(layout) {
			break
		}
	}
	fmt.Println(dates)
	return dates
}
