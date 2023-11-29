package ggapi

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
)

var pkgPath string

type Response struct {
	Status   bool      `json:"status"`
	Message  string    `json:"message"`
	User     User      `json:"user"`
	Networks []Network `json:"networks"`
}

type User struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type Network struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	CurrencyCode string `json:"currencyCode"`
	TimeZone     string `json:"timeZone"`
}

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalln("No caller information")
	}
	pkgPath = path.Dir(filename)
}

func Test() {
	c := exec.Command("php", pkgPath+"/test.php", "--name=tungdt", "--age=30")
	out, err := c.Output()
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n out: %+v \n", string(out))
}

func GetNetworks(refreshToken string) ([]byte, error) {
	c := exec.Command("php", pkgPath+"/network.php", "--refreshToken="+refreshToken)
	fmt.Printf("network: %+v \n", c.String())
	out, err := c.Output()
	if err != nil {
		return []byte{}, err
	}
	fmt.Printf("network: %+v \n", string(out))
	return out, err
}

func CheckAccessApi(refreshToken string, networkId int64, networkName string) (isEnable bool, err error) {
	c := exec.Command("php", pkgPath+"/check_api_access.php", "--refreshToken="+refreshToken, "--networkId="+strconv.FormatInt(networkId, 10), "--networkName="+networkName)
	out, err := c.Output()
	fmt.Println(c.String())
	if err != nil {
		isEnable = false
		return
	}
	fmt.Println(string(out))
	var response Response
	err = json.Unmarshal(out, &response)
	if err != nil {
		isEnable = false
		return
	}
	isEnable = true
	// Nếu có lỗi xảy ra check lại message lỗi
	if !response.Status {
		// Parse lỗi
		response.Message = strings.ReplaceAll(response.Message, "[", "")
		response.Message = strings.ReplaceAll(response.Message, "]", "")
		messages := strings.Split(response.Message, "@")
		var logError string
		if len(messages) > 0 {
			logError = messages[0]
		}
		// Nếu message lỗi là lỗi api access disable thì return lại thông báo
		if strings.TrimSpace(logError) == "AuthenticationError.NETWORK_API_ACCESS_DISABLED" || strings.TrimSpace(logError) == "refresh token is required" {
			isEnable = false
			return
		}
	}
	return
}

func CheckAdUnitByName(refreshToken string, networkId int64, networkName string, adUnitName string) (check bool, err error) {
	c := exec.Command("php", pkgPath+"/check_adunit.php",
		"--refreshToken="+refreshToken,
		"--networkId="+strconv.FormatInt(networkId, 10),
		"--networkName="+networkName,
		"--name="+adUnitName,
	)
	//fmt.Println(c.String())
	out, err := c.Output()
	if err != nil {
		return
	}
	//fmt.Println(string(out))
	var response Response
	err = json.Unmarshal(out, &response)
	if err != nil {
		return
	}
	check = response.Status
	return
}
