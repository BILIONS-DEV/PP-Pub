package main

import (
	"encoding/xml"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log"
	"source/apps/frontend/soap/response"
	"source/pkg/htmlblock"
	"source/core/technology/mysql"
)

func main() {
	URL := "https://ads.google.com/apis/ads/publisher/v202105/NetworkService?wsdl"
	AccessToken := GetAccessToken()
	if AccessToken == "" {
		panic("not have AccessToken")
	}

	payload := htmlblock.Render("network.xml", fiber.Map{})

	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Authorization", "Bearer "+AccessToken)
		r.Headers.Set("Content-Type", "application/xml;charset=UTF-8")
	})

	// attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		//fmt.Println(string(r.Body))
		resp := response.Envelope2{}
		err := xml.Unmarshal(r.Body, &resp)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", resp)
		//fmt.Printf("\n resp: %+v \n", resp)
	})

	err := c.PostRaw(URL, payload.Bytes())
	if err != nil {
		panic(err)
	}
}

var conf = &oauth2.Config{
	ClientID:     "697051439959-ola1rso2vsjf3cqhvmnoh3lnmkd5mn1t.apps.googleusercontent.com",
	ClientSecret: "HcItQm1nUvIIWgp1BCbWI07d",
	RedirectURL:  "https://self-serve.interdogmedia.com/gam/callback",
	Scopes: []string{
		"https://www.googleapis.com/auth/dfp",
	},
	Endpoint: google.Endpoint,
}

func GetAccessToken() (accessToken string) {
	//return ""
	gam := mysql.TableGam{}
	mysql.Client.Table("gam").Last(&gam, 4)
	token := &oauth2.Token{RefreshToken: gam.RefreshToken}
	tokenSource := conf.TokenSource(oauth2.NoContext, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	accessToken = newToken.AccessToken
	return
}
