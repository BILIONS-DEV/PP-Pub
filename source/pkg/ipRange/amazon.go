package ipRange

import (
	"encoding/json"
	"github.com/gocolly/colly/v2"
)

type amazon struct{}

type amazonStruct struct {
	SyncToken    string               `json:"syncToken"`
	CreateDate   string               `json:"createDate"`
	Prefixes     []amazonPrefixes     `json:"prefixes"`
	Ipv6Prefixes []amazonIpv6Prefixes `json:"ipv6_prefixes"`
}

type amazonPrefixes struct {
	IPPrefix           string `json:"ip_prefix"`
	Region             string `json:"region"`
	Service            string `json:"service"`
	NetworkBorderGroup string `json:"network_border_group"`
}

type amazonIpv6Prefixes struct {
	Ipv6Prefix         string `json:"ipv6_prefix"`
	Region             string `json:"region"`
	Service            string `json:"service"`
	NetworkBorderGroup string `json:"network_border_group"`
}

func (t *amazon) GetListIpAmazon() (listIp []string, err error) {
	cLinkIpRanges := c.Clone()
	amazon := amazonStruct{}
	cLinkIpRanges.OnResponse(func(r *colly.Response) {
		err = json.Unmarshal(r.Body, &amazon)
		if err != nil {
			return
		}
	})
	err = cLinkIpRanges.Visit("https://ip-ranges.amazonaws.com/ip-ranges.json")
	if err != nil {
		return
	}
	cLinkIpRanges.Wait()
	// Từ struct amazon json lấy ra ipv6 và ipv4
	for _, prefix := range amazon.Prefixes {
		if prefix.IPPrefix != "" {
			listIp = append(listIp, prefix.IPPrefix)
		}
	}
	for _, prefix := range amazon.Ipv6Prefixes {
		if prefix.Ipv6Prefix != "" {
			listIp = append(listIp, prefix.Ipv6Prefix)
		}
	}
	return
}
