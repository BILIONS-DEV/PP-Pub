package aerospike

import (
	"errors"
	as "github.com/aerospike/aerospike-client-go"
	"github.com/joho/godotenv"
	"log"
	"math"
	"os"
	"source/pkg/utility"
	"strconv"
	"strings"
)

type setting struct {
	Host      string
	Port      int
	Namespace string
}

type sets struct {
	Inventory    string
	DynamicFloor string
	DomainData   string
	Adset        string
	UserCookie   string
	Pixel        string
	JsCache      string
	BlockedPage  string
	BlockRPM     string
	IPRanges     string
	Quiz         string
	AdTag        string
	VldObject    string
	AutoAd       string
}

var Setting setting
var Sets sets
var Client *as.Client

func init() {
	/**
	Config Adset
	*/

	var setPrefix string
	_, err := os.Stat(os.Getenv("MODE") + ".env")
	if err == nil {
		envs, err := godotenv.Read(os.Getenv("MODE") + ".env")
		if err != nil {
			panic(errors.New("Error loading env file" + err.Error()))
		}
		setPrefix = envs["AE_SET_PREFIX"]
	}
	Sets = sets{
		Inventory:    setPrefix + "inventory",
		DynamicFloor: setPrefix + "dynamic_floor",
		Adset:        setPrefix + "adset",
		DomainData:   setPrefix + "domain_data",
		UserCookie:   setPrefix + "user_cookie",
		Pixel:        setPrefix + "pixel",
		JsCache:      setPrefix + "js_cache",
		BlockedPage:  setPrefix + "blocked_page",
		BlockRPM:     setPrefix + "block_rpm",
		IPRanges:     setPrefix + "ip_ranges",
		Quiz:         setPrefix + "quiz",
		AdTag:        setPrefix + "adtag",
		VldObject:    setPrefix + "vld_object",
		AutoAd:       setPrefix + "auto_ad",
	}

	/**
	Config Connect
	*/
	if utility.IsWindow() {
		Setting = setting{
			Host:      "127.0.0.1",
			Port:      3000,
			Namespace: "common",
		}
		Sets = sets{
			Inventory:    "inventory",
			DynamicFloor: "dynamic_floor",
			Adset:        "adset",
			DomainData:   "domain_data",
			UserCookie:   "user_cookie",
			Pixel:        "pixel",
			JsCache:      "js_cache",
			BlockedPage:  "blocked_page",
			BlockRPM:     "block_rpm",
			IPRanges:     "ip_ranges",
			Quiz:         "quiz",
			AdTag:        "adtag",
			AutoAd:       "auto_ad",
		}
	} else {
		Setting = setting{
			//Host:      "127.0.0.1", //=> aerospike trên backend
			//Host:      "23.92.69.154", //=> publicIP server service
			Host:      "192.168.9.13", //=> privateIP aerospike trên service
			Port:      3000,
			Namespace: "disk_cache",
		}
	}
	if !utility.IsWindow() {
		_, err = os.Stat(os.Getenv("MODE") + ".env")
		if err == nil {
			envs, err := godotenv.Read(os.Getenv("MODE") + ".env")
			if err != nil {
				panic(errors.New("Error loading env file" + err.Error()))
			}
			Setting.Namespace = envs["AE_NAMESPACE"]
		}
	}
	Client = Connect()
}

func Connect() (Client *as.Client) {
	var err error
	_, err = os.Stat(os.Getenv("MODE") + ".env")
	if err == nil {
		envs, err := godotenv.Read(os.Getenv("MODE") + ".env")
		if err != nil {
			panic(errors.New("Error loading env file" + err.Error()))
		}
		var aeHosts []*as.Host
		addresses := strings.Split(envs["AE_ADDRESS"], ",")
		for _, address := range addresses {
			addr := strings.Split(address, ":")
			if len(addr) != 2 {
				panic(errors.New("Env AE_ADDRESS malformed"))
			}
			h := addr[0]
			p, err := strconv.ParseInt(addr[1], 10, 64)
			if err != nil {
				panic(errors.New("Parser AE_ADDRESS port error: " + err.Error()))
			}
			aeHosts = append(aeHosts, &as.Host{
				Name: h,
				Port: int(p),
			})

			Client, err = as.NewClientWithPolicyAndHost(as.NewClientPolicy(), aeHosts...)
			if err != nil {
				panic("Couldn't connect to aerospike: " + err.Error())
			}
		}
	} else {
		Client, err = as.NewClient(Setting.Host, Setting.Port)
		if err != nil {
			panic("Couldn't connect to aerospike: " + err.Error())
		}
	}
	return
}

func MakeKey(setName string, key interface{}) (aeKey *as.Key, err error) {
	aeKey, err = as.NewKey(Setting.Namespace, setName, key)
	return
}

func MakeDefaultWritePolicy() (policy *as.WritePolicy) {
	policy = as.NewWritePolicy(0, math.MaxUint32)
	policy.SendKey = true
	policy.Expiration = math.MaxUint32
	return
}

func MakeDefaultWritePolicyWithTTL(ttl uint32) (policy *as.WritePolicy) {
	policy = as.NewWritePolicy(0, math.MaxUint32)
	policy.SendKey = true
	policy.Expiration = ttl
	return
}

func Truncate(setName string) error {
	err := Client.Truncate(nil, Setting.Namespace, setName, nil)
	return err
}

func Set(key interface{}, value interface{}, table ...string) (err error) {
	if len(table) < 1 {
		return errors.New("missing set name")
	}
	aeKey, err := MakeKey(table[0], key)
	if err != nil {
		log.Println(err)
		return
	}
	err = Client.PutObject(MakeDefaultWritePolicy(), aeKey, value)
	if err != nil {
		log.Println(err)
	}
	return
}

func Delete(key interface{}, table ...string) (err error) {
	if len(table) < 1 {
		return errors.New("missing set name")
	}
	aeKey, err := MakeKey(table[0], key)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = Client.Delete(MakeDefaultWritePolicy(), aeKey)
	if err != nil {
		log.Println(err)
	}
	return
}

func Exists(key interface{}, table ...string) (exists bool, err error) {
	if len(table) < 1 {
		err = errors.New("missing set name")
		return
	}
	aeKey, err := MakeKey(table[0], key)
	if err != nil {
		log.Println(err)
		return
	}
	exists, err = Client.Exists(nil, aeKey)
	if err != nil {
		log.Println(err)
	}
	return
}
