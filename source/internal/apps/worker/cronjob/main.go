package main

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"source/config"
	"source/infrastructure/caching"
	"source/infrastructure/mysql"
	"source/internal/entity/model"
	"source/internal/repo"
	"source/internal/usecase"
	"time"
)

type Channel struct {
	Response chan Response
	Limit    chan bool
}

type Response struct {
	Output *model.CronjobModel
	Error  []error
}

type handler struct {
	UseCases *usecase.UseCases
}

func main() {
	var (
		DB    *gorm.DB
		Cache caching.Cache
		err   error
	)
	configs := config.NewConfig()
	//=> init Mysql
	if configs.Mysql != nil {
		fmt.Println("init mysql....", true)
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
	//=> init Aerospike
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
	//=> init Repository
	repos := repo.NewRepositories(&repo.Deps{Db: DB, Cache: Cache})

	//=> init UseCases
	useCases := usecase.NewUseCases(&usecase.Deps{Repos: repos})

	t := &handler{
		UseCases: useCases,
	}

	//=> Start worker
	ticker := time.NewTicker(2 * time.Second)
	fmt.Println("Start worker cronjob...")
	done := make(chan bool)
	for {
		select {
		case <-ticker.C:
			t.StartWorker()
		case <-done:
			return
		}
	}
}

func (t *handler) StartWorker() {
	//=> Get cronjob
	cronjobs, _ := t.UseCases.CronJob.GetQueues(
		10,
	)

	total := len(cronjobs)
	// Make channel
	channel := Channel{
		Response: make(chan Response, total),
		Limit:    make(chan bool, 20),
	}
	// Producer
	for i, _ := range cronjobs {
		//_ = inventory.Lock()  //=> Lock row if in handle
		channel.Limit <- true //=> Set limit go routine
		// Go Routine
		go func(cronjob *model.CronjobModel, channel Channel) {
			err := t.UseCases.CronJob.Handler(cronjob)
			channel.Response <- Response{
				Output: cronjob,
				Error:  err,
			}
			<-channel.Limit
		}(&cronjobs[i], channel)
	}

	// Consume
	for i := 0; i < total; i++ {
		resp := <-channel.Response
		fmt.Printf("%+v \n", resp)
	}
}
