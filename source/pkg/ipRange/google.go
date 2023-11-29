package ipRange

import (
	"encoding/json"
	"github.com/gocolly/colly/v2"
)

var c = colly.NewCollector()

type google struct{}

type googleStruct struct {
	SyncToken    string     `json:"syncToken"`
	CreationTime string     `json:"creationTime"`
	Prefixes     []prefixes `json:"prefixes"`
}

type prefixes struct {
	Ipv4Prefix string `json:"ipv4Prefix"`
	Ipv6Prefix string `json:"ipv6Prefix"`
	Service    string `json:"service"`
	Scope      string `json:"scope"`
}

func (t *google) GetListIpGoogle() (listIp []string, err error) {
	// Get từ link goog.json
	cLinkGoog := c.Clone()
	goog := googleStruct{}
	cLinkGoog.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &goog)
		if err != nil {
			return
		}
	})
	err = cLinkGoog.Visit("https://www.gstatic.com/ipranges/goog.json")
	if err != nil {
		return
	}
	cLinkGoog.Wait()

	// Get từ link cloud.json
	cLinkClood := c.Clone()
	cloud := googleStruct{}
	cLinkClood.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &cloud)
		if err != nil {
			return
		}
	})
	err = cLinkClood.Visit("https://www.gstatic.com/ipranges/cloud.json")
	if err != nil {
		return
	}
	cLinkClood.Wait()

	// Từ 2 struct json lấy từ 2 link tạo list ip
	for _, prefix := range goog.Prefixes {
		if prefix.Ipv4Prefix != "" {
			listIp = append(listIp, prefix.Ipv4Prefix)
		}
		if prefix.Ipv6Prefix != "" {
			listIp = append(listIp, prefix.Ipv6Prefix)
		}
	}
	for _, prefix := range cloud.Prefixes {
		if prefix.Ipv4Prefix != "" {
			listIp = append(listIp, prefix.Ipv4Prefix)
		}
		if prefix.Ipv6Prefix != "" {
			listIp = append(listIp, prefix.Ipv6Prefix)
		}
	}
	return
}
