package router

import (
	"github.com/gofiber/fiber/v2"
	"source/apps/frontend/config"
	"source/apps/frontend/controller"
)

func Content(app *fiber.App) {
	Content := new(controller.Content)
	app.Get(config.URIContent, Content.Index)
	app.Post(config.URIContent, Content.Filter)
	app.Get(config.URIContentAdd, Content.Add)
	app.Get(config.URIContentAddVideo, Content.Add)
	app.Post(config.URIContentAdd, Content.AddPost)
	app.Get(config.URIContentEdit, Content.Edit)
	app.Get(config.URIContentEditVideo, Content.Edit)
	app.Post(config.URIContentEdit, Content.EditPost)
	app.Post(config.URIContentDel, Content.Delete)
	//app.Get(config.URIContentAddQuiz, Content.AddQuiz)
	//app.Post(config.URIContentAddQuiz, Content.AddQuizPost)
	//app.Get(config.URIContentEditQuiz, Content.EditQuiz)
	//app.Post(config.URIContentEditQuiz, Content.EditQuizPost)
}
