package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os"
	"os/signal"
	"reflect"
	"source/config"
	"source/internal/usecase"
	"syscall"
)

type Deps struct {
	UseCases  *usecase.UseCases
	AppConfig *config.Config
}

type handler struct {
	Deps
}

// NewHandler : khởi tạo service server với các tiêm phụ thuộc là userCases & translation dùng cho multiple language
//
// param: useCases *usecase.UseCases
// param: translation *lang.Translation
// return: handler *handler
func NewHandler(deps Deps) *handler {
	return &handler{deps}
}

// StartServer : start server với package fiber.
// Đây sẽ là nơi khai báo & cấu hình fiber
//
// param: port
func (h *handler) StartServer() {
	app := fiber.New(fiber.Config{
		BodyLimit:      20 * 1024 * 1024,
		ReadBufferSize: 128 * 1024,
		Views:          engine(),
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Origin, Content-Type, Accept, X-Token",
		AllowCredentials: true,
	}))
	h.setupViews(app)
	h.setupRouters(app)
	h.runService(app)
}

// setupViews : bổ sung các cấu hình khi sử dụng view của golang
func (h *handler) setupViews(app fiber.Router) {
	app.Static("/assets", "../../www/themes/muze/assets", fiber.Static{MaxAge: 9999999})
	app.Static("/assets", "./view/.assets", fiber.Static{MaxAge: 9999999, CacheDuration: 9999999})
}

// setupRouters : tự động nhận diện các hàm init routes thông qua các method bắt đầu bằng `InitRoutes` và setup các handler trong đó
func (h *handler) setupRouters(app fiber.Router) {
	handlerType := reflect.TypeOf(h)
	handlerValue := reflect.ValueOf(h)
	for i := 0; i < handlerType.NumMethod(); i++ {
		method := handlerType.Method(i)
		//=> tất cả các hàm initRoutes của controller đều phải bắt đầu bằng `InitRoutes`. Ví dụ: controller_user thì là InitRoutesUser
		methodPrefixRequired := string([]rune(method.Name)[0:10])
		if methodPrefixRequired == "InitRoutes" {
			method.Func.Call([]reflect.Value{handlerValue, reflect.ValueOf(app)})
		}
	}
}

// runService : khởi động ứng dụng
func (h *handler) runService(app *fiber.App) {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	serverShutdown := make(chan bool)
	go func(app *fiber.App) {
		<-sigterm
		log.Println("Gracefully shutting down...")
		err := app.Shutdown()
		if err != nil {
			log.Panic(err)
		}
		serverShutdown <- true
	}(app)
	if err := app.Listen(fmt.Sprintf("%s:%d", "127.0.0.1", 3000)); err != nil {
		log.Panic(err)
	}
}
