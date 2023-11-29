package report_mgid

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"reflect"
	"source/config"
	"source/infrastructure/caching"
	"source/infrastructure/kafka"
	"source/infrastructure/mysql"
	"source/internal/entity/model"
	"source/internal/repo"
	"source/pkg/utility"
	"testing"
)

var (
	Repos *repo.Repositories
	DB    *gorm.DB
	Cache caching.Cache
	Kafka *kafka.Client
)

func init() {
	var (
		err error
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
	Repos = repo.NewRepositories(&repo.Deps{Db: DB, Cache: Cache, Kafka: Kafka})
}
func Test_mgidUC_GetReport(t1 *testing.T) {

	type fields struct {
		repos *repo.Repositories
	}
	type args struct {
		campaignID string
		date       string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantReports model.AccountModel
		wantErr     bool
	}{
		{
			name:        "Test1",
			fields:      fields{Repos},
			args:        args{"11408990", "2023-01-16"},
			wantReports: model.AccountModel{},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &mgidUC{
				repos: tt.fields.repos,
			}
			gotReports, err := t.GetReport(tt.args.campaignID, tt.args.date)
			if (err != nil) != tt.wantErr {
				t1.Errorf("GetReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotReports, tt.wantReports) {
				t1.Errorf("GetReport() gotReports = %v, want %v", gotReports, tt.wantReports)
			}
		})
	}
}
