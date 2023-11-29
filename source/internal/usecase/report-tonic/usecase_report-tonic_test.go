package report_tonic

import (
	"fmt"
	"source/infrastructure/fakedb"
	"source/internal/repo"
	"strconv"
	"testing"
)

func Test_tonicUC_GetToken(t1 *testing.T) {
	db, err := fakedb.NewMysql()
	if err != nil {
		return
	}
	repos := repo.NewRepositories(&repo.Deps{
		Db:    db,
		Cache: nil,
		Kafka: nil,
	})

	reportTonicUC := NewReportTonicUC(repos)
	token, expires, err := reportTonicUC.GetToken()
	if err != nil {
		return
	}

	fmt.Println(token)
	fmt.Println(expires)
}

func Test_tonicUC_GetReportTracking(t1 *testing.T) {
	db, err := fakedb.NewMysql()
	if err != nil {
		return
	}
	repos := repo.NewRepositories(&repo.Deps{
		Db:    db,
		Cache: nil,
		Kafka: nil,
	})

	reportTonicUC := NewReportTonicUC(repos)

	reports, err := reportTonicUC.GetReportFinal(
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJkbG4iOiJxb1hCU0dnMXdJT25pR2ZiNFl2WFdzQzhrek1JYi1aSVVCbUpaWDhCLWJlOGl3NVpKYndCZUcyMWFoYmppSWpaQWd3Z0ZCVEY2TVQwNDBmQ1RnalRyaGhsUWVIQU9oclRpVkhqRUxYcDdIdWVhemxkZlM0M3NyT0JkUzlKUVZpeTBuQTlyTE9KejNER2lVdXcxNFRYbVBMNUx0Yy1ORW1HZ3JwUzlDQkJWUkUiLCJleHAiOjE2NjI1NDgwMjZ9.j47OPqAOtv1cgfr708eyEwa3EpDNGwsXQpddFfqCIEo",
		"2022-09-06")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(reports))
	total := 0.0
	for _, report := range reports {
		revenue, _ := strconv.ParseFloat(report.RevenueUsd, 64)
		total += revenue
		fmt.Println(report)
	}
	fmt.Println(total)
}
