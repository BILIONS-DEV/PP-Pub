package htmlblock

import (
	"bytes"
	"log"
	"os"
	"text/template"
)

var FilePath string

func New(filePath string) {
	FilePath = filePath
}

func ParseTemplate(templateFileName string, data interface{}) (buf *bytes.Buffer, err error) {
	buf = new(bytes.Buffer)
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return
	}
	err = t.Execute(buf, data)
	return buf, err
}

func Render(templateFileName string, data interface{}) (buf *bytes.Buffer) {
	var templateFilePath string
	var err error
	if FilePath != "" {
		templateFilePath = FilePath + "/" + templateFileName
	} else {
		templateFilePath = templateFileName
	}
	if _, err := os.Stat(templateFilePath); os.IsNotExist(err) {
		log.Println("file path not have: ", templateFilePath)
		return
	}
	buf = new(bytes.Buffer)
	buf, err = ParseTemplate(templateFilePath, data)
	if os.IsNotExist(err) {
		log.Println(err)
	}
	return
}
