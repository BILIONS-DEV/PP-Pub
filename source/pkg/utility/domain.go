package utility

import (
	UAX "github.com/mileusna/useragent"
	"github.com/oschwald/geoip2-golang"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"io"
	"net"
	"net/http"
	"os"
	"reflect"
	"source/pkg/logger"
	"strconv"
	"strings"
	"sync"
	"time"
)

var GeoIPMutex sync.Mutex
var GeoIPData *geoip2.Reader
var PublicSuffixList *publicsuffix.List

func init() {
	var err error
	// Begin GeoIP process
	remoteGeoIPURL := "https://static.vliplatform.com/maxmind-auto-update/GeoLite2-City.mmdb"

	GeoIPDBPath := "GeoLite2-City.mmdb"
	if _, err := os.Stat(GeoIPDBPath); os.IsNotExist(err) {
		err := DownloadFile(GeoIPDBPath, remoteGeoIPURL)
		if err != nil {
			logger.Log.Fatal(err.Error())
		}
	}

	GeoIPMutex.Lock()
	GeoIPData, err = geoip2.Open(GeoIPDBPath)
	if err != nil {
		logger.Log.Fatal(err.Error())
	}
	GeoIPMutex.Unlock()

	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				resp, err := http.Get(remoteGeoIPURL)
				if err != nil {
					logger.Log.Info("[GeoIP] " + err.Error())
					continue
				}

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					logger.Log.Info("[GeoIP] " + err.Error())
					continue
				}
				GeoIPMutex.Lock()
				GeoIPDataTemp, err := geoip2.FromBytes(body)
				if err != nil {
					logger.Log.Info("[GeoIP] " + err.Error())
					GeoIPMutex.Unlock()
					continue
				} else {
					GeoIPData = GeoIPDataTemp
					logger.Log.Info("[GeoIP] GEO DB updated successfully.")
				}
				GeoIPMutex.Unlock()
			}
		}
	}()
	// End GeoIP process

	// Config PublicSuffixList to get root domain
	publicSuffixListPath := "public_suffix_list.dat"
	if _, err := os.Stat(publicSuffixListPath); os.IsNotExist(err) {
		err := DownloadFile(publicSuffixListPath, "https://publicsuffix.org/list/public_suffix_list.dat")
		if err != nil {
			logger.Log.Fatal(err.Error())
		}
	}
	PublicSuffixList, err = publicsuffix.NewListFromFile(publicSuffixListPath, &publicsuffix.ParserOption{PrivateDomains: true, ASCIIEncoded: false})
	if err != nil {
		logger.Log.Fatal(err.Error())
	}

	go func() {
		for {
			time.Sleep(time.Second * 21600)
			err := DownloadFile(publicSuffixListPath, "https://publicsuffix.org/list/public_suffix_list.dat")
			if err != nil {
				logger.Log.Info("[PublicSuffixList] " + err.Error())
			} else {
				PublicSuffixList, err = publicsuffix.NewListFromFile(publicSuffixListPath, &publicsuffix.ParserOption{PrivateDomains: true, ASCIIEncoded: false})
				if err != nil {
					logger.Log.Info("[PublicSuffixList] " + err.Error())
				}
			}
		}
	}()
}

func GetRootDomain(dm string) (string, error) {
	domainName, err := publicsuffix.DomainFromListWithOptions(PublicSuffixList, dm, &publicsuffix.FindOptions{})
	domainName = strings.ReplaceAll(domainName, "https://", "")
	domainName = strings.ReplaceAll(domainName, "http://", "")
	return domainName, err
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func FindIndex(array interface{}, value interface{}) int {
	switch v := array.(type) {
	case []string:
		for i, item := range v {
			if item == value {
				return i
			}
		}
		break
	case []int64:
		for i, item := range v {
			if item == value {
				return i
			}
		}
		break
	}
	return -1
}

func RemoveIndexSliceInt64(array []int64, index int) []int64 {
	return append(array[:index], array[index+1:]...)

}

func RemoveIndexSliceString(array []string, index int) []string {
	return append(array[:index], array[index+1:]...)

}

func ConvertArrayToString(array interface{}) string {
	switch v := array.(type) {
	case []int:
		var a []string
		for _, row := range v {
			s := strconv.Itoa(row)
			a = append(a, s)
		}
		return strings.Join(a, ",")
	case []string:
		return strings.Join(v, ",")
		//case []float64:
		//	var a []string
		//	for _, row := range v {
		//		s := fmt.Sprintf("%.2f", row)
		//		a = append(a, s)
		//	}
		//	return strings.Join(a, ",")
	}
	return ""
}

func ConvertStringToArray(string2 string) (arr []string) {
	arr = strings.Split(string2, ",")
	return
}

func ToInterface(obj interface{}) interface{} {
	vp := reflect.New(reflect.TypeOf(obj))
	vp.Elem().Set(reflect.ValueOf(obj))
	return vp.Interface()
}

func GetDeviceFromUA(userAgent string) string {
	ua := UAX.Parse(userAgent)
	if ua.Mobile {
		return "mobile"
	}
	if ua.Tablet {
		return "tablet"
	}
	if ua.Desktop {
		return "desktop"
	}
	return ""
}

func GetGeoDbFromIP(ip string) *geoip2.City {
	var record *geoip2.City
	ipAddr := net.ParseIP(ip)
	GeoIPMutex.Lock()
	record, err := GeoIPData.City(ipAddr)
	if err != nil {
		logger.Log.Info(err.Error())
	}
	GeoIPMutex.Unlock()
	return record
}

func GetGeoDbFromIP2(ip string) (record *geoip2.City, err error) {
	ipAddr := net.ParseIP(ip)
	GeoIPMutex.Lock()
	record, err = GeoIPData.City(ipAddr)
	if err != nil {
		return nil, err
	}
	GeoIPMutex.Unlock()
	return record, nil
}

func GetCountryCodeFromIP(ip string) string {
	ipAddr := net.ParseIP(ip)
	GeoIPMutex.Lock()
	record, err := GeoIPData.City(ipAddr)
	if err != nil {
		logger.Log.Info(err.Error())
		GeoIPMutex.Unlock()
		return "N/A"
	}
	GeoIPMutex.Unlock()
	return record.Country.IsoCode
}
