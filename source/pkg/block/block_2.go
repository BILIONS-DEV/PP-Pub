package block

import (
	"bytes"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
	"os"
	"path/filepath"
	"text/template"
)

type Block struct {
	m          *minify.M
	path       string
	minifyHtml bool
}

func NewBlock(options Option) *Block {
	var block Block
	block.path = options.Path
	block.minifyHtml = options.MinifyHtml
	block.m = minify.New()
	return &block
}

func (t *Block) RenderToString(templateFileName string, data interface{}) (output string) {
	outBuf, err := t.Render(templateFileName, data)
	if err != nil {
		return ""
	}
	return outBuf.String()
}

func (t *Block) Render(templateFileName string, data interface{}) (output *bytes.Buffer, err error) {
	templateFilePath := path + "/" + templateFileName
	fileName, templateByte, err := t.ReadFileOS(templateFilePath)
	if err != nil {
		return
	}
	templateString := string(templateByte)
	if minifyHtml {
		//m.AddFunc("text/html", html.Minify)
		t.m.Add("text/html", &html.Minifier{
			KeepConditionalComments: false,
			KeepDefaultAttrVals:     false,
			KeepWhitespace:          false,
		})
		templateString, err = t.m.String("text/html", templateString)
		if err != nil {
			return
		}
	}
	output = new(bytes.Buffer)
	err = t.Parse(templateString, output, data, fileName)
	return
}

func (t *Block) Parse(templateText string, output *bytes.Buffer, data interface{}, fileName string) (err error) {
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

func (t *Block) ReadFileOS(file string) (name string, b []byte, err error) {
	name = filepath.Base(file)
	b, err = os.ReadFile(file)
	return
}

// ReadContent opens a named file and read content from it
func (t *Block) ReadContent(rf *bytes.Buffer, name string) (n int64, err error) {
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
