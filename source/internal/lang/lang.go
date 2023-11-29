package lang

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"source/pkg/utility"
)

type Schema struct {
	Lang        string      `json:"lang"`
	Translation Translation `json:"translation"`
}

type MapLang map[string]Translation

var (
	LANG      = "EN"
	Map       = make(MapLang)
	Translate Translation
)

func init() {
	readConfig()
}

func Register(lang string) *Translation {
	if lang != "" {
		LANG = lang
	}
	Translate = Map[LANG]
	return &Translate
}

func readConfig() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	langPath := path.Dir(filename) + "/data"
	files, err := utility.WalkMatch(langPath, "*.json")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		var config Schema
		if config, err = parseLang(file); err != nil {
			continue
		}
		Map[config.Lang] = config.Translation
	}
}

func parseLang(file string) (config Schema, err error) {
	var jsonFile *os.File
	if jsonFile, err = os.Open(file); err != nil {
		return
	}
	defer func(jsonFile *os.File) {
		if err := jsonFile.Close(); err != nil {
			return
		}
	}(jsonFile)
	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err = json.Unmarshal(byteValue, &config); err != nil {
		return
	}
	return
}
