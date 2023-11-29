package helper

import (
	"fmt"
	"github.com/gofiber/template/html"
	"github.com/hermanschaaf/prettyprint"
	"github.com/leekchan/accounting"
	"html/template"
	"source/apps/frontend/view"
	"source/core/theme"
	"source/pkg/block"
	"source/pkg/utility"
	"strings"
	"time"
)

func Engine() *html.Engine {
	// Create a new engine by passing the template folder
	// and template extension using <engine>.New(dir, ext string)
	//engine := html.New("apps/endpoint/views", ".html")
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
	engine.AddFunc("GetEmailById", block.GetEmailById)
	engine.AddFunc("EmbedCssInline", theme.EmbedCssInline)
	engine.AddFunc("EmbedCSS", theme.EmbedCSS)
	engine.AddFunc("EmbedJS", theme.EmbedJS)
	engine.AddFunc("IsActiveSidebar", view.IsActiveSidebar)
	engine.AddFunc("IsActiveSidebarWithGroup", view.IsActiveSidebarWithGroup)
	engine.AddFunc("safeHTML", safeHTML)
	engine.AddFunc("InArray", utility.InArray)
	engine.AddFunc("FirstCharacter", utility.FirstCharacter)
	engine.AddFunc("FormatFloat", utility.FormatFloat)
	engine.AddFunc("StrToTime", StrToTime)
	engine.AddFunc("IncIndex", incIndex)
	engine.AddFunc("FormatHtml", FormatHtml)
	engine.AddFunc("Nl2br", Nl2br)
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
}

func Nl2br(text string) template.HTML {
	return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
}

func FormatHtml(html string) template.HTML {
	//html := `<html><head></head><body><a href="http://test.com">Test link</a><p><br/></p></body></html>`
	// second parameter is the indent to use
	// (can be anything: spaces, a tab, etc.)
	pretty, err := prettyprint.Prettify(html, "    ")
	if err != nil {
		fmt.Printf("\n err: %+v \n", err)
		// do something in case of error
	}
	// print the prettified html:
	fmt.Println(pretty)
	return template.HTML(pretty)
}

func incIndex(number int) int {
	return number + 1
}

func safeHTML(s string) template.HTML {
	return template.HTML(s)
}

var acc = accounting.Accounting{Symbol: "$", Precision: 2}

func USDFormat(value interface{}) string {
	return acc.FormatMoney(value)
}

func TimeToClockInSidebar() int64 {
	return time.Now().Unix()
}

func StrToTime(timeStr string) int64 {
	LayoutISO := "Mon, 02 Jan 2006 15:04:05 MST"
	t, err := time.Parse(LayoutISO, timeStr)
	if err != nil {
		fmt.Println(err)
	}
	return t.Unix()
}
