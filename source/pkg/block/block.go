package block

import (
	"bytes"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
	"os"
	"path/filepath"
	"text/template"
)

var (
	path       = ""
	minifyHtml = false
)

type Option struct {
	Path       string
	MinifyHtml bool
}

var m *minify.M

func New(options Option) {
	path = options.Path
	minifyHtml = options.MinifyHtml
	m = minify.New()
}

func RenderToString(templateFileName string, data interface{}) (output string) {
	outBuf, err := Render(templateFileName, data)
	if err != nil {
		return ""
	}
	return outBuf.String()
}

func Render(templateFileName string, data interface{}) (output *bytes.Buffer, err error) {
	templateFilePath := path + "/" + templateFileName
	fileName, templateByte, err := ReadFileOS(templateFilePath)
	if err != nil {
		return
	}
	templateString := string(templateByte)
	if minifyHtml {
		//m.AddFunc("text/html", html.Minify)
		m.Add("text/html", &html.Minifier{
			KeepConditionalComments: false,
			KeepDefaultAttrVals:     false,
			KeepWhitespace:          false,
		})
		templateString, err = m.String("text/html", templateString)
		if err != nil {
			return
		}
	}
	output = new(bytes.Buffer)
	err = Parse(templateString, output, data, fileName)
	return
}

func Parse(templateText string, output *bytes.Buffer, data interface{}, fileName string) (err error) {
	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New(fileName).Funcs(FunMaps()).Parse(templateText)
	if err != nil {
		//log.Fatalf("parsing: %s", err)
		return err
	}
	// Run the template to verify the output.
	err = tmpl.Execute(output, data)
	if err != nil {
		//log.Fatalf("execution: %s", err)
		return err
	}
	return
}

func ReadFileOS(file string) (name string, b []byte, err error) {
	name = filepath.Base(file)
	b, err = os.ReadFile(file)
	return
}

// ReadContent opens a named file and read content from it
func ReadContent(rf *bytes.Buffer, name string) (n int64, err error) {
	// Read file
	f, err := os.Open(filepath.Clean(name))
	if err != nil {
		return 0, err
	}
	defer func() {
		err = f.Close()
	}()
	return rf.ReadFrom(f)
}
