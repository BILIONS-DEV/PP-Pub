package scanAds

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"net"
	"net/http"
	"net/url"
	"source/pkg/utility"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	ErrorStatus      bool
	Error            error
	AdsTxtUrl        string
	HeaderStatusCode int
	Lines            []LineTracking
}

type LineTracking struct {
	Text  string
	Match bool
}

type LineInfo struct {
	Domain          string
	Id              string
	Type            string
	CertificationID string
	Comment         string
	Line            string
}

var c = colly.NewCollector(
	colly.UserAgent("scan ads.txt by APD"),
	colly.AllowURLRevisit(),
)

func init() {
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   15 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})
	if utility.IsWindow() {
		c.SetRequestTimeout(1 * time.Second)
	} else {
		c.SetRequestTimeout(15 * time.Second)
	}
}

func Scan(compares []string, samples []string) (resp Response) {
	if len(samples) == 0 {
		resp.ErrorStatus = true
		resp.Error = errors.New("samples contains no data")
		return
	}
	var results []LineTracking
	for _, sampleLine := range samples {
		// Với các dòng bắt đầu bằng # hoặc rỗng thì không cần phải so sánh mặc định là true
		if strings.HasPrefix(sampleLine, "#") || sampleLine == "" {
			results = append(results, LineTracking{
				Text:  sampleLine,
				Match: true,
			})
			continue
		}
		sample, err := ParseLine(sampleLine)
		if err != nil {
			continue
		}
		isMatch := false
		for _, compareLine := range compares {
			compare, err := ParseLine(compareLine)
			if err != nil {
				continue
			}
			if sample.Id == compare.Id && sample.Domain == compare.Domain && sample.Type == compare.Type && sample.CertificationID == compare.CertificationID {
				isMatch = true
			}
		}
		results = append(results, LineTracking{
			Text:  sample.Line,
			Match: isMatch,
		})
	}
	resp.Lines = results
	return
}

func ScanUrl(adsTxtUrl string, samples []string) (resp Response) {
	compares, err, statusCode := ReadURL(adsTxtUrl)
	if err != nil {
		resp.ErrorStatus = true
		resp.Error = err
		resp.HeaderStatusCode = statusCode
		resp.AdsTxtUrl = adsTxtUrl
		return
	}
	resp = Scan(compares, samples)
	resp.AdsTxtUrl = adsTxtUrl
	return
}

func makeRequest(c *colly.Collector) {
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	Transport := &http.Transport{
		Proxy: http.ProxyURL(&url.URL{
			Scheme: "http",
			User:   url.UserPassword("dgjgqumv-rotate", "djav3i3rt11z"),
			Host:   "p.webshare.io:80",
		}),
	}
	c.WithTransport(Transport)
}

func ReadURL(url string) (compares []string, err error, statusCode int) {
	//makeRequest(c)
	//fmt.Printf("%+v \n", url)
	var linkInvalid bool
	collect := c.Clone()
	collect.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})
	collect.SetRequestTimeout(30 * time.Second)
	collect.OnResponse(func(r *colly.Response) {
		if !strings.Contains(strings.ToLower(r.Headers.Get("Content-Type")), "text/plain") || !strings.Contains(strings.ToLower(r.Headers.Get("Cache-Control")), "no-cache") {
			linkInvalid = true
		} else {
			linkInvalid = false
		}
		body := string(r.Body)
		//fmt.Printf("%+v \n", body)
		fmt.Printf("Header: %+v \n", r.Headers)
		fmt.Printf("Header: %+v \n", r.Request.URL)
		compares = utility.SplitLines(body)
		statusCode = r.StatusCode
	})
	// Set error handler
	collect.OnError(func(r *colly.Response, errHandle error) {
		//fmt.Println("ScanAds --Url: ", url, " --errHandle: ", errHandle, " --Body: ", string(r.Body))
		statusCode = r.StatusCode
		if errHandle != nil {
			err = errHandle
		} else {
			if r.StatusCode != 200 {
				err = errors.New(fmt.Sprintf("status code error: %d", r.StatusCode))
			}
		}
		if err != nil {
			linkInvalid = true
		}
	})
	err = collect.Visit(url)
	collect.Wait()
	if linkInvalid {
		linkInvalid = false
		err = collect.Visit(url + "?t=" + strconv.FormatInt(time.Now().Unix(), 10))
		collect.Wait()

		if linkInvalid {
			err = collect.Visit(url)
			collect.Wait()
		}
	}
	return
}

func ParseLine(line string) (parse LineInfo, err error) {
	// Tách phần comment nếu có
	parts := strings.SplitN(line, "#", 2)
	var comment string
	if len(parts) == 2 {
		comment = strings.TrimSpace(parts[1])
	}
	mainPart := strings.TrimSpace(parts[0])

	parser := strings.Split(mainPart, ",")
	if len(parser) < 3 || len(parser) > 4 {
		err = errors.New("line is not in the correct format")
		return
	}

	typeText := strings.ToLower(strings.TrimSpace(parser[2]))
	if typeText != "reseller" && typeText != "direct" {
		err = errors.New("type of line is not in the correct format")
		return
	}

	parse.Domain = strings.ToLower(strings.TrimSpace(parser[0]))
	parse.Id = strings.ToLower(strings.TrimSpace(parser[1]))
	parse.Type = typeText
	parse.Comment = comment
	parse.Line = line
	if len(parser) == 4 {
		parse.CertificationID = strings.TrimSpace(parser[3])
	}

	return
}
