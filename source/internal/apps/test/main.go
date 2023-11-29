package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/golang-jwt/jwt/v4"
	url2 "net/url"
	"os"
	"time"
)

var SecretKey = []byte("this_is_my_key")

func main() {
	listVideoIDs, err := GetListVideoID()
	if err != nil {
		return
	}

	fmt.Printf("%+v \n", listVideoIDs)
	err = SaveListVideoInfo(listVideoIDs)
	if err != nil {
		return
	}
	//tokenString, err := generateJWT()
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("\n tokenString: %+v \n", tokenString)
	//tokenString += "12345"
	//
	//tokenParse, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	//	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	//		return nil, fmt.Errorf("there was an error in parsing")
	//	}
	//	return SecretKey, nil
	//})
	//if err != nil {
	//	panic(err)
	//}
	//if claims, ok := tokenParse.Claims.(jwt.MapClaims); ok && tokenParse.Valid {
	//	fmt.Printf("\n claims: %+v \n", claims)
	//}
}

func generateJWT() (tokenString string, err error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"email":      "tungdt2@83.com.vn",
		"role":       "role",
		"exp":        time.Now().Add(time.Minute * 30).Unix(),
	}
	tokenString, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(SecretKey)
	return
}

type ResponseAds struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestId string `json:"request_id"`
	Data      struct {
		List []struct {
			AdName                    string      `json:"ad_name"`
			IdentityType              string      `json:"identity_type"`
			AdgroupName               string      `json:"adgroup_name"`
			Deeplink                  string      `json:"deeplink"`
			VideoId                   string      `json:"video_id"`
			CreateTime                string      `json:"create_time"`
			LandingPageUrl            string      `json:"landing_page_url"`
			IsNewStructure            bool        `json:"is_new_structure"`
			ImpressionTrackingUrl     interface{} `json:"impression_tracking_url"`
			AdgroupId                 string      `json:"adgroup_id"`
			LandingPageUrls           interface{} `json:"landing_page_urls"`
			ViewabilityVastUrl        interface{} `json:"viewability_vast_url"`
			AdRefPixelId              int64       `json:"ad_ref_pixel_id"`
			AvatarIconWebUri          string      `json:"avatar_icon_web_uri"`
			IdentityId                string      `json:"identity_id"`
			CallToActionId            string      `json:"call_to_action_id"`
			IsAco                     bool        `json:"is_aco"`
			ClickTrackingUrl          interface{} `json:"click_tracking_url"`
			OptimizationEvent         string      `json:"optimization_event"`
			PageId                    interface{} `json:"page_id"`
			ModifyTime                string      `json:"modify_time"`
			BrandSafetyVastUrl        interface{} `json:"brand_safety_vast_url"`
			OperationStatus           string      `json:"operation_status"`
			DisplayName               string      `json:"display_name"`
			AppName                   string      `json:"app_name"`
			DeeplinkType              string      `json:"deeplink_type"`
			CampaignName              string      `json:"campaign_name"`
			AdFormat                  string      `json:"ad_format"`
			CreativeType              interface{} `json:"creative_type"`
			AdTexts                   interface{} `json:"ad_texts"`
			CreativeAuthorized        bool        `json:"creative_authorized"`
			PlayableUrl               string      `json:"playable_url"`
			ImageIds                  []string    `json:"image_ids"`
			CampaignId                string      `json:"campaign_id"`
			MusicId                   interface{} `json:"music_id"`
			CardId                    interface{} `json:"card_id"`
			SecondaryStatus           string      `json:"secondary_status"`
			BrandSafetyPostbidPartner string      `json:"brand_safety_postbid_partner"`
			ProfileImageUrl           string      `json:"profile_image_url"`
			VastMoatEnabled           bool        `json:"vast_moat_enabled"`
			ViewabilityPostbidPartner string      `json:"viewability_postbid_partner"`
			AdvertiserId              string      `json:"advertiser_id"`
			AdId                      string      `json:"ad_id"`
			FallbackType              string      `json:"fallback_type"`
			TrackingPixelId           int64       `json:"tracking_pixel_id"`
			AdText                    string      `json:"ad_text"`
		} `json:"list"`
		PageInfo struct {
			Page        int `json:"page"`
			PageSize    int `json:"page_size"`
			TotalPage   int `json:"total_page"`
			TotalNumber int `json:"total_number"`
		} `json:"page_info"`
	} `json:"data"`
}

const (
	AdvertiserID = "7229589319848116226"
	CampaignID   = "1765744184237073"
	RedirectID   = "35939"
)

func GetListVideoID() (listVideoIDs []string, err error) {
	time.Sleep(3 * time.Second)
	url := "https://business-api.tiktok.com/open_api/v1.3/ad/get/"
	//=> Thêm Auth Key của để check report của mình
	values := url2.Values{}
	values.Set("advertiser_id", AdvertiserID)
	values.Set("filtering", "{\"campaign_ids\":[\""+CampaignID+"\"]}")

	//=> Tạo link get Reports
	urlReport := url + "?" + values.Encode()
	fmt.Println(urlReport)

	//=> Dùng Colly request lên link để lấy về response Reports
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	c.MaxBodySize = 100 * 1024 * 1024
	// Request
	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Access-Token", "12f193559ebd286a751748d14faffd0c041ff8f1")
		request.Headers.Set("Content-Type", "application/json")
	})
	var response ResponseAds
	// Check Response
	c.OnResponse(func(r *colly.Response) {
		fmt.Printf("%+v \n", string(r.Body))
		err = json.Unmarshal(r.Body, &response)
	})
	// Check Error
	c.OnError(func(r *colly.Response, errHandle error) {
		if errHandle != nil {
			err = errHandle
		} else {
			if r.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("status code error: %d", r.StatusCode))
			}
		}
	})
	err = c.Visit(urlReport)
	c.Wait()
	// Sau khi đợi xử lý xong nếu có lỗi thì return luôn
	if err != nil {
		return
	}
	for _, data := range response.Data.List {
		listVideoIDs = append(listVideoIDs, data.VideoId)
	}
	return
}

func SaveListVideoInfo(listVideoIDs []string) (err error) {
	bListVideoIDs, err := json.Marshal(&listVideoIDs)
	if err != nil {
		return
	}

	time.Sleep(3 * time.Second)
	url := "https://business-api.tiktok.com/open_api/v1.3/file/video/ad/search/"
	//=> Thêm Auth Key của để check report của mình
	values := url2.Values{}
	values.Set("advertiser_id", AdvertiserID)
	values.Set("filtering", "{\"video_ids\":"+string(bListVideoIDs)+"}")

	//=> Tạo link get Reports
	urlReport := url + "?" + values.Encode()
	fmt.Println(urlReport)

	//=> Dùng Colly request lên link để lấy về response Reports
	var c = colly.NewCollector(
		colly.AllowURLRevisit(),
	)
	c.MaxBodySize = 100 * 1024 * 1024
	// Request
	c.OnRequest(func(request *colly.Request) {
		request.Headers.Set("Access-Token", "12f193559ebd286a751748d14faffd0c041ff8f1")
		request.Headers.Set("Content-Type", "application/json")
	})
	// Check Response
	c.OnResponse(func(r *colly.Response) {
		err = os.WriteFile(RedirectID+".json", r.Body, 0644)
	})
	// Check Error
	c.OnError(func(r *colly.Response, errHandle error) {
		if errHandle != nil {
			err = errHandle
		} else {
			if r.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("status code error: %d", r.StatusCode))
			}
		}
	})
	err = c.Visit(urlReport)
	c.Wait()
	// Sau khi đợi xử lý xong nếu có lỗi thì return luôn
	if err != nil {
		return
	}
	return
}
