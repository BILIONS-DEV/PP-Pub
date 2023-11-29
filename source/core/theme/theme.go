package theme

import (
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"github.com/valyala/fasttemplate"
	"github.com/yosssi/gohtml"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"source/pkg/utility"
	"strings"
	"time"
)

var (
	MapThemes            = make(MapTheme)
	version              = uuid.New().String()
	defaultKeyForAllPage = `{{allPage}}`
	Reload               = false
	ModuleAssetsPath     = "./view"
	StaticAssetsPath     = "../../www/themes/muze/assets"
)

type MapTheme map[string]MapPage
type MapPage map[string]*PageMap

type PageMap struct {
	Path            string `xml:"path,attr"`
	IncludePath     string `xml:"includePath,attr"`
	ListCss         []Css  `xml:"css"`
	ListCssOnTop    []Css  `xml:"cssOnTop"`
	ListCssOnBottom []Css  `xml:"cssOnBottom"`
	ListJs          []Js   `xml:"js"`
	ListJsOnTop     []Js   `xml:"jsOnTop"`
	ListJsOnBottom  []Js   `xml:"jsOnBottom"`

	ListJsForRender  []Js
	ListCssForRender []Css
	CssInline        template.HTML
	CssTag           template.HTML
	JsTag            template.HTML
	Theme            string
	IsRenderList     bool
}

type Theme struct {
	//XMLName xml.Name `xml:"theme"`
	Name   string    `xml:"name,attr"`
	Reload bool      `xml:"reload,attr"`
	Pages  []PageMap `xml:"page"`
}

type Css struct {
	Href   string `xml:"css"`
	Link   string `xml:"link,attr"`
	Inline bool   `xml:"inline,attr"`
	Media  string `xml:"media,attr"`
	Rel    string `xml:"rel,attr"`
	Module bool   `xml:"module,attr"`
}

type Js struct {
	Src    string `xml:"js"`
	Link   string `xml:"link,attr"`
	Async  bool   `xml:"async,attr"`
	Defer  bool   `xml:"defer,attr"`
	Module bool   `xml:"module,attr"`
}

func init() {
	readConfig()
	if !Reload {
		build()
	}
}

// build Based on data declared from init, render jsTag and cssTag to make data ready for EmbedCSS & EmbedJS function
func build() {
	for theme, mapPages := range MapThemes {
		for _, page := range mapPages {
			if page.Path == defaultKeyForAllPage {
				continue
			}
			page.Theme = theme
			page.buildTag()
		}
	}
}
func (page *PageMap) buildTag() *PageMap {
	if !page.IsRenderList {
		// Make listCSS & ListJs from defaultPage & includePage
		defaultPage, haveDefaultPage := MapThemes[page.Theme][defaultKeyForAllPage]
		includePage, haveIncludePage := MapThemes[page.Theme][page.IncludePath]
		//=> include Default Page
		if haveDefaultPage {
			page.ListCssForRender = append(page.ListCssForRender, defaultPage.ListCssOnTop...) // append
			page.ListCssForRender = append(page.ListCssForRender, defaultPage.ListCss...)      // append
			page.ListJsForRender = append(page.ListJsForRender, defaultPage.ListJsOnTop...)    // append
			page.ListJsForRender = append(page.ListJsForRender, defaultPage.ListJs...)         // append
			//page.ListJsForRender = append(defaultPage.ListJs, page.ListJsForRender...)      // prepend
		}
		//=> include Include Page
		if haveIncludePage {
			page.ListCssForRender = append(page.ListCssForRender, includePage.ListCssOnTop...) // append
			page.ListCssForRender = append(page.ListCssForRender, includePage.ListCss...)      // append
			page.ListJsForRender = append(page.ListJsForRender, includePage.ListJsOnTop...)    // append
			page.ListJsForRender = append(page.ListJsForRender, includePage.ListJs...)         // append
		}
		//=> Current Page
		page.ListCssForRender = append(page.ListCssForRender, page.ListCssOnTop...) // append
		page.ListCssForRender = append(page.ListCssForRender, page.ListCss...)      // append
		page.ListJsForRender = append(page.ListJsForRender, page.ListJsOnTop...)    // append
		page.ListJsForRender = append(page.ListJsForRender, page.ListJs...)         // append
		//=> include in bottom Default Page
		if haveDefaultPage {
			page.ListCssForRender = append(page.ListCssForRender, defaultPage.ListCssOnBottom...) // append
			page.ListJsForRender = append(page.ListJsForRender, defaultPage.ListJsOnBottom...)    // append
		}
		//=> include in bottom Include Page
		if haveIncludePage {
			page.ListCssForRender = append(page.ListCssForRender, includePage.ListCssOnBottom...) // append
			page.ListJsForRender = append(page.ListJsForRender, includePage.ListJsOnBottom...)    // append
		}
		//=> include in bottom Current Page
		page.ListCssForRender = append(page.ListCssForRender, page.ListCssOnBottom...) // append
		page.ListJsForRender = append(page.ListJsForRender, page.ListJsOnBottom...)    // append
		page.IsRenderList = true
	}
	// Reset for Reload
	if Reload {
		page.CssTag = ""
		page.JsTag = ""
		page.CssInline = ""
	}
	page.CssInline = page.buildCssInline()
	page.CssTag = page.buildCssTag()
	page.JsTag = page.buildJsTag()
	return page
}

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// readConfig Read config from xml and handle assign for MapThemes
func readConfig() {
	files, err := WalkMatch("./view", "*.xml")
	if err != nil {
		panic(err)
	}
	for i, file := range files {
		theme := ParseXML(".", file)
		if theme.Reload && utility.IsWindow() && i == 0 {
			Reload = true
		}
		if MapThemes[theme.Name] == nil {
			MapThemes[theme.Name] = make(MapPage)
		}
		for _, pages := range theme.Pages {
			MapThemes[theme.Name][pages.Path] = &PageMap{
				Path:            pages.Path,
				IncludePath:     pages.IncludePath,
				ListCss:         pages.ListCss,
				ListCssOnTop:    pages.ListCssOnTop,
				ListCssOnBottom: pages.ListCssOnBottom,
				ListJs:          pages.ListJs,
				ListJsOnTop:     pages.ListJsOnTop,
				ListJsOnBottom:  pages.ListJsOnBottom,
			}
		}
	}
}

// ParseXML Read config from file xml
func ParseXML(themeFolder string, file string) Theme {
	var (
		xmlFile *os.File
		err     error
		theme   Theme
	)
	if xmlFile, err = os.Open(themeFolder + "/" + file); err != nil {
		panic(err)
	}
	defer func(xmlFile *os.File) {
		if err := xmlFile.Close(); err != nil {
			panic(err)
		}
	}(xmlFile)
	byteValue, _ := io.ReadAll(xmlFile)
	if err = xml.Unmarshal(byteValue, &theme); err != nil {
		panic(err)
	}
	return theme
}

// readFile read content file
func readFile(filePath string) string {
	var file *os.File
	var err error
	if file, err = os.Open(filePath); err != nil {
		panic(err)
	}
	defer func(cssFile *os.File) {
		err := cssFile.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	byteValue, _ := io.ReadAll(file)
	return string(byteValue)
}

// buildCssTag
func (page *PageMap) buildCssInline() template.HTML {
	var cssInline, filePath string
	for _, css := range page.ListCssForRender {
		if !css.Inline {
			continue
		}
		if css.Module {
			filePath = strings.Replace(css.Link, "/assets", ModuleAssetsPath, -1)
		} else {
			filePath = strings.Replace(css.Link, "/assets", StaticAssetsPath, -1)
		}
		cssInline += readFile(filePath)
	}
	return template.HTML(cssInline)
}

// buildCssTag render <link> tag from Page.ListCss and save to Page.CssTag and response for EmbedCSS
func (page *PageMap) buildCssTag() template.HTML {
	temp := fasttemplate.New(`<link rel="stylesheet" type="text/css" href="{{href}}?v={{version}}"{{more}} />`, "{{", "}}")
	var cssTag string
	for _, css := range page.ListCssForRender {
		if css.Link == "" {
			continue
		}
		if css.Inline {
			continue
		}
		cssTag += temp.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
			switch tag {
			case "href":
				return w.Write([]byte(css.Link))
			case "version":
				return w.Write([]byte(version))
			case "more":
				var options []string
				if css.Rel != "" {
					options = append(options, `rel="`+css.Rel+`"`)
				}
				if css.Media != "" {
					options = append(options, `media="`+css.Media+`"`)
				}
				more := " " + strings.Join(options, " ")
				return w.Write([]byte(more))
			default:
				return w.Write([]byte(""))
			}
		})
	}
	return template.HTML(gohtml.Format(cssTag))
}

// buildJsTag render <link> tag from Page.ListJs and save to Page.JsTag and response for EmbedJS
func (page *PageMap) buildJsTag() template.HTML {
	temp := fasttemplate.New(`<script src="{{src}}?v={{version}}"{{more}}></script>`, "{{", "}}")
	var jsTag string
	for _, js := range page.ListJsForRender {
		link := js.Src
		if link == "" {
			link = js.Link
		}
		if link == "" {
			continue
		}
		jsTag += temp.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
			switch tag {
			case "src":
				return w.Write([]byte(link))
			case "version":
				return w.Write([]byte(version))
			case "more":
				var options []string
				if js.Module {
					options = append(options, `type="module"`)
				}
				if js.Async {
					options = append(options, "async")
				}
				if js.Defer {
					options = append(options, "defer")
				}
				more := " " + strings.Join(options, " ")
				return w.Write([]byte(more))
			default:
				return w.Write([]byte(""))
			}
		})
	}
	return template.HTML(gohtml.Format(jsTag))
}

// EmbedCSS Generate Page.CssTag tags to embed in html/template
func EmbedCSS(theme string, templatePath string) (tag template.HTML) {
	if Reload {
		version = time.Now().String()
		build()
	}
	if page, ok := MapThemes[theme][templatePath]; ok {
		return page.CssTag
	}
	return
}

// EmbedCssInline Generate Page.CssInline tags to embed in html/template
func EmbedCssInline(theme string, templatePath string) (cssInline template.HTML) {
	if page, ok := MapThemes[theme][templatePath]; ok {
		return template.HTML(fmt.Sprintf(`<style>%s</style>`, page.CssInline))
	}
	return
}

// EmbedJS Generate Page.JsTag tags to embed in html/template
func EmbedJS(theme string, templatePath string) (tag template.HTML) {
	if page, ok := MapThemes[theme][templatePath]; ok {
		return page.JsTag
	}
	return tag
}
