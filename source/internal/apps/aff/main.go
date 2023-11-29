package main

import (
	"gorm.io/gorm"
	"log"
	"source/config"
	"source/infrastructure/caching"
	"source/infrastructure/mysql"
	"source/internal/apps/aff/controller"
	"source/internal/repo"
	"source/internal/usecase"
)

func main() {

	var (
		DB    *gorm.DB
		Cache caching.Cache
		err   error
	)
	configs := config.NewConfig()

	//=> init Aerospike
	if configs.Aerospike != nil {
		log.Println("init aerospike...", true)
		if Cache, err = caching.NewAerospike(caching.Config{
			Host:      configs.Aerospike.Host,
			Port:      configs.Aerospike.Port,
			Namespace: configs.Aerospike.Namespace,
		}); err != nil {
			log.Fatalln("Error connect Cache: ", err)
			return
		}
	}

	//=> init Mysql
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

	//=> init Repository
	repos := repo.NewRepositories(&repo.Deps{Db: DB, Cache: Cache})

	//=> init UseCases
	useCases := usecase.NewUseCases(&usecase.Deps{Repos: repos})

	//=> init Controller
	handle := controller.NewHandler(controller.Deps{
		UseCases:  useCases,
		AppConfig: configs,
	})
	handle.StartServer()
	log.Println("Running cleanup tasks...")

}
