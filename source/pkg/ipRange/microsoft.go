package ipRange

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
)

type microsoft struct{}

type microsoftStruct struct {
	ChangeNumber int              `json:"changeNumber"`
	Cloud        string           `json:"cloud"`
	Values       []microsoftValue `json:"values"`
}

type microsoftValue struct {
	Name       string            `json:"name"`
	ID         string            `json:"id"`
	Properties microsoftProperty `json:"properties"`
}

type microsoftProperty struct {
	ChangeNumber    int      `json:"changeNumber"`
	Region          string   `json:"region"`
	RegionID        int      `json:"regionId"`
	Platform        string   `json:"platform"`
	SystemService   string   `json:"systemService"`
	AddressPrefixes []string `json:"addressPrefixes"`
	NetworkFeatures []string `json:"networkFeatures"`
}

func (t *microsoft) GetListIpMicrosoft() (listIp []string, err error) {
	cGetLink := c.Clone()

	linkDowload := ""
	cGetLink.OnHTML("a[data-bi-id=downloadretry]", func(e *colly.HTMLElement) {
		linkDowload = e.Attr("href")
	})

	err = cGetLink.Visit("https://www.microsoft.com/en-us/download/confirmation.aspx?id=56519")
	if err != nil {
		fmt.Println(err)
	}
	cGetLink.Wait()

	cGetJson := c.Clone()
	microsoft := microsoftStruct{}
	cGetJson.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &microsoft)
		if err != nil {
			return
		}
	})

	err = cGetJson.Visit(linkDowload)
	if err != nil {
		return
	}
	cGetJson.Wait()
	// Từ struct json lấy ra list ip
	for _, value := range microsoft.Values {
		for _, ip := range value.Properties.AddressPrefixes {
			listIp = append(listIp, ip)
		}
	}
	return
}
