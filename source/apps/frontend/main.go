package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"source/apps/frontend/controllerv2"
	"source/apps/frontend/helper"
	"source/apps/frontend/helper/router"
	"source/apps/frontend/lang"
	"source/config"
	"source/infrastructure/caching"
	"source/infrastructure/mysql"
	lang2 "source/internal/lang"
	"source/internal/repo"
	"source/internal/usecase"
	"source/pkg/block"
	"source/pkg/htmlblock"
	"source/pkg/utility"
)

func main() {

	lang.Register("EN")

	engine := helper.Engine()
	app := fiber.New(fiber.Config{
		BodyLimit:      20 * 1024 * 1024,
		ReadBufferSize: 128 * 1024,
		Views:          engine,
	})
	app.Static("/jsx", "./view/ads_txt_react", fiber.Static{MaxAge: 1})
	app.Static("/assets", "../../www/themes/muze/assets", fiber.Static{MaxAge: 9999999})
	app.Static("/assets", "./view/.assets", fiber.Static{MaxAge: 9999999, CacheDuration: 9999999})
	//app.Static("/assets", "./view/pages")
	htmlblock.New("./view/pages")
	block.New(block.Option{
		Path:       "./view/pages",
		MinifyHtml: true,
	})

	initForV2(app) // init cho clean architecture

	//app.Use(helper.Bootstrap)
	router.Register(app)
	run(app)
}

func initForV2(app fiber.Router) {
	var (
		DB    *gorm.DB
		Cache caching.Cache
		err   error
	)
	configs := config.NewConfig()
	fmt.Printf("%+v\n", configs.Mysql)
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
	// Register language for multiple site
	translation := lang2.Register("EN")
	//=> init Repository
	repos := repo.NewRepositories(&repo.Deps{Db: DB, Cache: Cache})

	//=> init UseCases
	useCases := usecase.NewUseCases(&usecase.Deps{Repos: repos, Translation: translation})

	//=> init Controller
	handle := controllerv2.NewHandler(controllerv2.Deps{
		UseCases:    useCases,
		Translation: translation,
		Block: block.NewBlock(block.Option{
			Path:       "./view/pages",
			MinifyHtml: true,
		}),
	})
	handle.StartServer(8550)
}

func run(app *fiber.App) {
	if utility.IsWindow() {
		// log.Fatal(app.ListenTLS(":8540", "../../localhost.pem", "../../localhost-key.pem"))
		log.Fatal(app.Listen(":8550"))
	} else {
		log.Fatal(app.Listen(":8550"))
	}
}
