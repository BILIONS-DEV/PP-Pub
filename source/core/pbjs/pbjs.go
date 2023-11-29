package pbjs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

type Bundle struct {
	PBJSPath string
}

// NewPBJSBundle Khởi tạo trình bundle
func NewPBJSBundle(pbjsPath string) *Bundle {
	return &Bundle{
		PBJSPath: pbjsPath,
	}
}

// Build
// Khởi tạo trình bundle
// Arguments
// @input listModules (array string)
// @input filePath (string) - đường dẫn lưu file output
// @input fileName (string) - tên file output
// Return
// đường dẫn đầy đủ của file được build (string)
// lỗi (error)
func (p *Bundle) Build(listModules []string, fileName string) (fileLocation string, err error) {
	ti := time.Now()
	if len(listModules) == 0 {
		return "", errors.New("the list of modules cannot be empty")
	}

	if len(fileName) == 0 {
		return "", errors.New("file name cannot be empty")
	}

	listModules = append(listModules, "bidtestBidAdapter")

	listModulesString := strings.Join(listModules, ",")

	cmdBuild := exec.Command(p.PBJSPath+"/node_modules/.bin/gulp", "bundle", `--modules=`+listModulesString, "--bundleName="+fileName)
	cmdBuild.Dir = p.PBJSPath
	var stderr bytes.Buffer
	cmdBuild.Stderr = &stderr

	out, err := cmdBuild.Output()
	if err != nil {
		return "", errors.New(fmt.Sprint(err) + ": " + stderr.String())
	}
	fmt.Println("Command Result: " + string(out))
	fmt.Println("Build", fileName, "success in", time.Since(ti))
	fileLocation = p.PBJSPath + "/build/dist/" + fileName
	return fileLocation, nil
}

func (p *Bundle) GetVersionPrebid() (currentPbJsVersion string, err error) {
	packageJSONFile, err := ioutil.ReadFile(p.PBJSPath + "/package.json")
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	err = json.Unmarshal(packageJSONFile, &data)
	if err != nil {
		return "", err
	}

	if val, ok := data["version"]; ok {
		currentPbJsVersion = val.(string)
	}
	return
}
