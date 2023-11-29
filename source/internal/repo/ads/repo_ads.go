package ads

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"source/internal/entity/model"
	"source/pkg/telegram"
	"source/pkg/utility"
	"strings"
	// "source/core/technology/mysql"
)

type RepoAds interface {
	Filter(input *InputFilter) (records []model.AdsModel, total int64, err error)
	GetToken(input *InputGetToken) (response ResponseGetToken)
	ChangeActionAd(input ChangeActionAdRP) (err error)
	// Save(record *model.AdsModel) (err error)
}

type adsRepo struct{}

func NewAdsRepo() *adsRepo {
	return &adsRepo{}
}

const (
	UrlApiPocPoc = "https://api.pocpoc.io"
)

type InputFilter struct {
	User      model.User
	Status    interface{}
	Inventory string
	Search    interface{}
	Page      int
	Limit     int
	Order     string
}

type PayloadFilterPocPoc struct {
	Email     string      `json:"email,omitempty"`
	Inventory string      `json:"inventory,omitempty"`
	Placement string      `json:"placement,omitempty"`
	Status    interface{} `json:"status,omitempty"`
	Search    interface{} `json:"search,omitempty"`
	Limit     int         `json:"limit"`
	Page      int         `json:"page"`
	OrderBy   string      `json:"order_by,omitempty"`
	Sort      string      `json:"sort,omitempty"`
}

type ResponseAds struct {
	Status  bool            `json:"status"`
	Message string          `json:"message,omitempty"`
	Errors  []Error         `json:"errors,omitempty"`
	Data    ResponseDataAds `json:"data"`
}
type ResponseDataAds struct {
	Limit      int              `json:"limit"`
	Page       int              `json:"page"`
	Sort       int              `json:"sort"`
	TotalRows  int              `json:"total_rows"`
	TotalPages int              `json:"total_pages"`
	Rows       []model.AdsModel `json:"rows"`
}

func (t *adsRepo) Filter(input *InputFilter) (records []model.AdsModel, total int64, err error) {
	responseLogin := t.GetTokenAdmin()
	if !responseLogin.Status || responseLogin.Data.AccessToken == "" {
		telegram.SendErrorTuan(responseLogin.Message, "Get token PocPoc prom Pubpower fail")
	}

	params := PayloadFilterPocPoc{
		Inventory: input.Inventory,
		Email:     input.User.Email,
		// Placement: "",
		Status:  input.Status,
		Search:  input.Search,
		Limit:   input.Limit,
		Page:    input.Page,
		OrderBy: input.Order,
	}
	jsonEncode, errb := json.Marshal(params)
	if errb != nil {
		return
	}
	url := UrlApiPocPoc + "/ad/filters2"
	if utility.IsWindow() {
		url = "http://127.0.0.1:9191/ad/filters2"
	}
	method := "POST"
	payload := strings.NewReader(string(jsonEncode))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Printf("%+v\n", responseLogin.Data.AccessToken)
	req.Header.Add("X-Token", responseLogin.Data.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		return
	}
	var response ResponseAds
	json.Unmarshal(body, &response)
	// fmt.Println("==============================================")
	//fmt.Println("body: ", string(body))
	//fmt.Printf("%+v\n", response)
	if response.Status {
		records = response.Data.Rows
		total = int64(response.Data.TotalRows)
	}
	return
}

type InputGetToken struct {
	Email    string
	Password string
}
type ResponseGetToken struct {
	Status  bool          `json:"status"`
	Message string        `json:"message,omitempty"`
	Errors  []Error       `json:"errors,omitempty"`
	Data    TokenResponse `json:"data,omitempty"`
}
type TokenResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
type Error struct {
	ID        string `json:"id,omitempty"`
	Message   string `json:"message"`
	ErrorCode string `json:"error_code,omitempty"`
}

func (t *adsRepo) GetToken(input *InputGetToken) (response ResponseGetToken) {
	input.Password = "PocPoc@123"
	// input.Password = input.Email + "1#"
	if input.Email == "thulinh0705@gmail.com" {
		input.Password = "Linh25101997@"
	}
	if input.Email == "xuantm.ad@gmail.com" {
		input.Password = "Xuantm8228@@"
	}
	//fmt.Printf("%+v\n", input)
	url := UrlApiPocPoc + "/account/login/publisher"
	if utility.IsWindow() {
		url = "http://127.0.0.1:9191/account/login/publisher"
	}
	method := "POST"
	payloadByte, err := json.Marshal(input)
	if err != nil {
		response.Status = false
		response.Message = "Get token PocPoc fail!"
		response.Errors = append(response.Errors, Error{
			ID:        "",
			Message:   err.Error(),
			ErrorCode: "400",
		})
		return
	}
	payload := strings.NewReader(string(payloadByte))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		response.Status = false
		response.Message = "Get token PocPoc fail!"
		response.Errors = append(response.Errors, Error{
			ID:        "",
			Message:   err.Error(),
			ErrorCode: "400",
		})
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	//fmt.Printf("%+v\n", string(body))
	if err != nil {
		fmt.Println(err)
		response.Status = false
		response.Message = "Get token PocPoc fail!"
		response.Errors = append(response.Errors, Error{
			ID:        "",
			Message:   err.Error(),
			ErrorCode: "400",
		})
		return
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		response.Status = false
		response.Message = "Get token PocPoc fail!"
		response.Errors = append(response.Errors, Error{
			ID:        "",
			Message:   err.Error(),
			ErrorCode: "400",
		})
		return
	}
	return
}

func (t *adsRepo) GetTokenAdmin() (response ResponseGetToken) {
	var input = InputGetToken{
		Email:    "tuantt@bil.vn",
		Password: "159357555",
	}

	url := UrlApiPocPoc + "/account/login/admin"
	if utility.IsWindow() {
		url = "http://127.0.0.1:9191/account/login/admin"
	}
	method := "POST"
	payloadByte, err := json.Marshal(input)
	if err != nil {
		response.Status = false
		response.Message = "Get token PocPoc fail!"
		response.Errors = append(response.Errors, Error{
			ID:        "",
			Message:   err.Error(),
			ErrorCode: "400",
		})
		return
	}
	payload := strings.NewReader(string(payloadByte))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		response.Status = false
		response.Message = "Get token PocPoc fail!"
		response.Errors = append(response.Errors, Error{
			ID:        "",
			Message:   err.Error(),
			ErrorCode: "400",
		})
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		response.Status = false
		response.Message = "Get token PocPoc fail!"
		response.Errors = append(response.Errors, Error{
			ID:        "",
			Message:   err.Error(),
			ErrorCode: "400",
		})
		return
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(err)
		response.Status = false
		response.Message = "Get token PocPoc fail!"
		response.Errors = append(response.Errors, Error{
			ID:        "",
			Message:   err.Error(),
			ErrorCode: "400",
		})
		return
	}
	return
}

type ChangeActionAdRP struct {
	UserLogin model.User
	Email     string `json:"email"`
	AdID      int64  `json:"ad_id"`
	Action    string `json:"action"`
	Inventory string `json:"inventory"`
	Placement string `json:"placement"`
}

func (t *adsRepo) ChangeActionAd(input ChangeActionAdRP) (err error) {
	//responseLogin := t.GetToken(&InputGetToken{
	//	Email:    input.UserLogin.Email,
	//	Password: "pocpoc@" + strconv.FormatInt(input.UserLogin.ID, 10),
	//})
	responseLogin := t.GetTokenAdmin()
	//fmt.Printf("%+v\n", responseLogin)
	if !responseLogin.Status || responseLogin.Data.AccessToken == "" {
		telegram.SendErrorTuan(responseLogin.Message, "Get token PocPoc prom Pubpower fail")
		err = errors.New("Get token fail")
		return
	}

	jsonEncode, errb := json.Marshal(input)
	if errb != nil {
		return
	}

	url := UrlApiPocPoc + "/ad/block_by_publisher"
	if utility.IsWindow() {
		url = "http://127.0.0.1:9191/ad/block_by_publisher"
	}
	//fmt.Printf("%+v\n", url)
	method := "POST"
	// fmt.Printf("%+v\n", url)
	// fmt.Printf("%+v\n", string(jsonEncode))
	payload := strings.NewReader(string(jsonEncode))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("X-Token", responseLogin.Data.AccessToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var response ResponseAds
	json.Unmarshal(body, &response)
	//fmt.Println("body: ", string(body))
	// fmt.Printf("%+v\n", response)
	if !response.Status {
		if response.Message != "" {
			return errors.New(response.Message)
		} else if len(response.Errors) > 0 {
			return errors.New(response.Errors[0].Message)
		} else {
			return errors.New("Change Action Fail!")
		}
	}
	return
}
