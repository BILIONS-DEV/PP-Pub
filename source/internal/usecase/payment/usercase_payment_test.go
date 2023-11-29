package payment

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"source/config"
	"source/infrastructure/caching"
	"source/infrastructure/mysql"
	"source/internal/entity/model"
	"source/internal/repo"
	"source/internal/repo/payment"
	"source/pkg/datatable"
	"testing"
)

func TestPaymentUsecase_FiltersInvoice(t *testing.T) {
	var (
		DB *gorm.DB
		// Cache caching.Cache
		err error
	)
	configs := config.NewConfig()

	// => init Aerospike
	if configs.Aerospike != nil {
		log.Println("init aerospike...", true)
		if _, err = caching.NewAerospike(caching.Config{
			Host:      configs.Aerospike.Host,
			Port:      configs.Aerospike.Port,
			Namespace: configs.Aerospike.Namespace,
		}); err != nil {
			log.Fatalln("Error connect Cache: ", err)
			return
		}
	}

	// => init Mysql
	if configs.Mysql != nil {
		log.Println("init mysql....", true)
		if DB, err = mysql.Connect(mysql.Config{
			Username: configs.Mysql.Username,
			Password: configs.Mysql.Password,
			Host:     configs.Mysql.Host,
			Port:     configs.Mysql.Port,
			Database: configs.Mysql.Database,
			Encoding: configs.Mysql.Encoding,
		}); err != nil {
			log.Fatalln("Error connect mysql: ", err)
			return
		}
	}

	paymentRepo := payment.NewPaymentRepo(DB)
	repos := &repo.Repositories{
		Payment: paymentRepo,
	}

	paymentUC := NewPaymentUsecase(repos, nil)

	var record = model.ParamPaymentIndex{
		Request:     datatable.Request{},
		QuerySearch: "",
		FromMoney:   "",
		ToMoney:     "",
		PaidDate:    "",
		Month:       "",
		Sort:        "",
		Permission:  "sale",
		Page:        "2",
		Status:      nil,
	}

	payments, _, err := paymentUC.GetByFiltersInvoice(&record, model.User{})
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", payments)
}
