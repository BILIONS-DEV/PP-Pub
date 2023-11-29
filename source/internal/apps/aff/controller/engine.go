package controller

import (
	"github.com/gofiber/template/html"
	"html/template"
	"source/pkg/utility"
)

func engine() *html.Engine {
	engine := html.New("./view/pages", ".gohtml")
	engine.Layout("embed")    // Optional. Default: "embed"
	engine.Delims("{{", "}}") // Optional. Default: engine delimiters
	if utility.IsWindow() {
		engine.Reload(true) // Optional. Default: false ... Bật cái này lên là bị race condition
		engine.Debug(false) // Optional. Default: false
	}
	addFunc(engine)
	return engine
}

func addFunc(engine *html.Engine) {
	//engine.AddFunc("IsActiveSidebar", view.IsActiveSidebar)
	//engine.AddFunc("IsActiveSidebarWithGroup", view.IsActiveSidebarWithGroup)

	//engine.AddFunc("EmbedCssInline", theme.EmbedCssInline)
	//engine.AddFunc("EmbedCSS", theme.EmbedCSS)
	//engine.AddFunc("EmbedJS", theme.EmbedJS)
	//engine.AddFunc("safeHTML", theme.SafeHTML)
	//engine.AddFunc("FirstCharacter", utility.FirstCharacter)
	//engine.AddFunc("FormatFloat", utility.FormatFloat)
	//engine.AddFunc("StrToTime", theme.StrToTime)
	//engine.AddFunc("IncIndex", theme.IncIndex)
	//engine.AddFunc("FormatHtml", theme.FormatHtml)
	//engine.AddFunc("Nl2br", theme.Nl2br)
	//engine.AddFunc("InArray", utility.InArray)
	//engine.AddFunc("MkArray", theme.MkSlice)
	//engine.AddFunc("MkSlice", theme.MkSlice)
	//engine.AddFunc("InnerString", theme.InnerString)
	engine.AddFunc("InSlice", func(search interface{}, array interface{}) bool {
		return utility.InArray(search, array, true)
	})
	// AddFunc adds a function to the template's global function map.
	engine.AddFunc("greet", func(name string) string {
		return "Hello, greetFunc: " + name + "!"
	})
	engine.AddFunc("IsStringEmpty", func(name string) bool {
		if name == "" {
			return false
		} else {
			return true
		}
	})
	engine.AddFunc("unescape", func(s string) template.HTML {
		return template.HTML(s)
	})
}
