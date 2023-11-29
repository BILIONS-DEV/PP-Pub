package ipRange

import (
	"github.com/gocolly/colly/v2"
	"strings"
)

type cloudFlare struct{}

func (t *cloudFlare) GetListIpCloudFlare() (listIp []string, err error) {
	cLinkIpRanges := c.Clone()
	cLinkIpRanges.OnResponse(func(r *colly.Response) {
		listIpResponse := strings.Split(string(r.Body), "\n")
		for _, ip := range listIpResponse {
			listIp = append(listIp, strings.TrimSpace(ip))
		}
	})
	err = cLinkIpRanges.Visit("https://www.cloudflare.com/ips-v4")
	if err != nil {
		return
	}
	err = cLinkIpRanges.Visit("https://www.cloudflare.com/ips-v6")
	if err != nil {
		return
	}
	cLinkIpRanges.Wait()
	// Từ struct amazon json lấy ra ipv6 và ipv4
	//for _, prefix := range amazon.Prefixes {
	//	if prefix.IPPrefix != "" {
	//		listIp = append(listIp, prefix.IPPrefix)
	//	}
	//}
	//for _, prefix := range amazon.Ipv6Prefixes {
	//	if prefix.Ipv6Prefix != "" {
	//		listIp = append(listIp, prefix.Ipv6Prefix)
	//	}
	//}
	return
}
