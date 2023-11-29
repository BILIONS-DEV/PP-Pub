package utility

import (
	"bufio"
	"crypto/aes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gocolly/colly/v2"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var GOOS = runtime.GOOS

var ServiceLog = "https://be.pubpower.io/api/v1/log/worker"

var dev = false
var mod = false

//var ServiceLog = "http://127.0.0.1:8542/api/v1/log/worker"

type LogCreate struct {
	Function string `form:"column:function" json:"function"`
	Path     string `form:"column:" json:"path"`
	Line     string `form:"column:" json:"line"`
	Message  string `form:"column:" json:"message"`
}

func SetDev(b bool) {
	dev = b
}

func IsDev() bool {
	return dev
}

func SetDemo(b bool) {
	mod = b
}

func IsDemo() bool {
	return mod
}

func IsWindow() bool {
	if GOOS == "windows" {
		return true
	}
	return false
}

func WalkMatchOne(root, pattern string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		fileName := file.Name()
		if matched, err := filepath.Match(pattern, filepath.Base(fileName)); err != nil {
			return []string{}, err
		} else if matched {
			files = append(files, fileName)
		}
	}
	return files, nil
}

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// SplitLines Split string to lines
// https://stackoverflow.com/questions/8757389/reading-a-file-line-by-line-in-go
//
// param: s
// return:
func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}

// InArray Checks if the element exists in the array
//
// param: search
// param: array
// param: deep
// return:
func InArray(search interface{}, array interface{}, deep bool) bool {
	val := reflect.ValueOf(array)
	val = val.Convert(val.Type())

	typ := reflect.TypeOf(array).Kind()

	switch typ {
	case reflect.Map:
		s := val.MapRange()

		for s.Next() {
			s.Value().Convert(s.Value().Type())
			for i := 0; i < s.Value().Len(); i++ {
				if deep {
					if reflect.DeepEqual(search, s.Value().Index(i).Interface()) {
						return true
					}
				} else {
					str := s.Value().Index(i).String()
					if strings.Contains(str, search.(string)) {
						return true
					}
				}
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(search, val.Index(i).Interface()) {
				return true
			}
		}
	}

	return false
}

// InInt64Array Checks if the element exists in the string array
//
// param: val
// param: array
// return:
func InInt64Array(target int64, arr []int64) bool {
	for _, value := range arr {
		if value == target {
			return true
		}
	}
	return false
}

// InStringArray Checks if the element exists in the string array
//
// param: val
// param: array
// return:
func InStringArray(val string, array []string) (exists bool, index int) {
	exists = false
	index = -1

	for i, v := range array {
		if v == val {
			index = i
			exists = true
			break
		}
	}

	return
}

// IsNumeric Checks if the string is numeric
//
// param: s
// return:
func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// GetMD5Hash Convert string to md5
//
// param: string
// return: string md5
func GetMD5Hash(text string) string {
	harsher := md5.New()
	harsher.Write([]byte(text))
	return hex.EncodeToString(harsher.Sum(nil))
}

// DecodeUrl Decode url string
//
// param: string
// return: string decodedValue decode,error
func DecodeUrl(encodedValue string) (decodedValue string, err error) {
	decodedValue, err = url.QueryUnescape(encodedValue)
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}

// EncodeUrl Encode url string
//
// param: name=thai%nd&phone=%2B9199999999&phone=%2B628888888888
// return: name = [thai nd] phone = [+9199999999 +628888888888]
func EncodeUrl(decodedValue string) (encodedValue url.Values, err error) {
	encodedValue, err = url.ParseQuery(decodedValue)
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}

// Difference
//
// param: slice1
// param: slice2
// return:
func Difference(slice1 []string, slice2 []string) []string {
	var diff []string

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

// ConvertNumberToString
//
// param: input
// return:
func ConvertNumberToString(input interface{}) string {
	switch value := input.(type) {
	case int:
		if value > 0 {
			return strconv.Itoa(value)
		} else {
			return "0"
		}
	case float64:
		if value > 0 {
			return fmt.Sprintf("%f", value)
		} else {
			return "0"
		}
	default:
		return "0"
	}
}

// ConvertTime
//
// param: date
// return:
func ConvertTime(date time.Time) (ymd string, err error) {
	LayoutIsoYmd := "20060102"
	//timeUTC := time.Unix(date.Unix(), 0).UTC()
	//timezone := "America/New_York"
	//location, err := time.LoadLocation(timezone)
	//time := timeUTC.In(location)
	if err != nil {
		return
	}
	ymd = date.Format(LayoutIsoYmd)
	return
}

// FirstCharacter
//
// param: string
// return:
func FirstCharacter(character string) string {
	return character[0:1]
}

// FormatFloat
//
// param: float64
// return:
func FormatFloat(num float64, decimal int) float64 {
	output := math.Pow(10, float64(decimal))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// IsUrl
// Kiểm tra có phải là một url không
// param: string url
// return: bool
func IsUrl(str string) (isUrl bool) {
	u, err := url.Parse(str)
	if err != nil {
		return
	}
	if u.Scheme == "" && u.Host == "" {
		return
	}
	if !govalidator.IsURL(str) {
		return
	}
	return true
}

// Convert string list number to array int64
func ConvertStringListToArrayInt64(str string) (output []int64, err error) {
	arrStr := strings.Split(str, ",")
	for _, v := range arrStr {
		i, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return output, err
		}
		output = append(output, i)
	}
	return
}

// Standardize the string
func StandardizeTheString(str string) string {
	str = strings.TrimSpace(str)
	str = strings.ToLower(str)
	str = strings.ReplaceAll(str, " ", "_")
	return str
}

func ValidateString(s string) (valid string) {
	/*TrimSpace*/
	valid = strings.TrimSpace(s)
	return
}

func ValidateEmail(email string) bool {
	var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegexp.MatchString(email) {
		return false
	}
	return true
}

func ValidateDomainName(domain string) bool {
	RegExp := regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z]{2,3})$`)
	return RegExp.MatchString(domain)
}

func PushLogError(data LogCreate) {
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "application/json")
	})
	byteSql, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error Marshal data")
	}
	// attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		//fmt.Println(r.Body)
		//fmt.Printf("\n resp: %+v \n", resp)
	})

	err = c.PostRaw(ServiceLog, byteSql)
	if err != nil {
		fmt.Println("Error send api")
	}
}

func HandleError(err error, function string) (logError LogCreate) {
	if err != nil {
		_, path, line, _ := runtime.Caller(1)
		logError.Path = path
		logError.Function = function
		logError.Line = strconv.Itoa(line)
		logError.Message = err.Error()
	}
	return
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func ArrayToString(array interface{}, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(array), " ", delim, -1), "[]")
}

func Base64Decode(str string) (string, error) {
	switch len(str) % 4 {
	case 2:
		str += "=="
	case 3:
		str += "="
	}

	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func DecryptAes128Ecb(data, key []byte) (decrypted []byte, err error) {
	cipher, _ := aes.NewCipher([]byte(key))
	decrypted = make([]byte, len(data))
	size := 16

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		if bs > len(data) || be > len(data) {
			err = errors.New("decrypted fail")
			return
		}
		cipher.Decrypt(decrypted[bs:be], data[bs:be])
	}

	return
}

func SumFloat64(values ...float64) (output float64) {
	if len(values) == 0 {
		return 0.0
	}
	for _, value := range values {
		revOld := decimal.NewFromFloat(output)
		revAdd := decimal.NewFromFloat(value)

		sumRev := revOld.Add(revAdd)
		output, _ = sumRev.Float64()
	}
	return
}
