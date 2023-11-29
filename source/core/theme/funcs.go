package theme

import (
	"fmt"
	"github.com/hermanschaaf/prettyprint"
	"github.com/leekchan/accounting"
	"html/template"
	"strings"
	"time"
)

func MkSlice(args ...interface{}) []interface{} {
	return args
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

func IncIndex(number int) int {
	return number + 1
}

func SafeHTML(s string) template.HTML {
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

func InnerString(inputs ...interface{}) (output string) {
	for _, input := range inputs {
		output += fmt.Sprintf("%v", input)
	}
	return
}
