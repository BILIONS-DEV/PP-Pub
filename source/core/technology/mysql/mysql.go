package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"source/pkg/helpers"
	"source/pkg/utility"
	"time"

	"github.com/joho/godotenv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	RevenueShareDefault = 90
)

type config struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	Encoding string
}

type tables struct {
	User                      string
	UserInfo                  string
	Inventory                 string
	InventoryAdTag            string
	TagRewardedAdConfig       string
	RlAdTagSizeAdditional     string
	TagInfo                   string
	Country                   string
	Channels                  string
	ChannelsKeyword           string
	Content                   string
	Device                    string
	AdSize                    string
	Bidder                    string
	RlBidderMediaType         string
	LineItem                  string
	LineItemV2                string
	BidderAssign              string
	AdTag                     string
	Rule                      string
	RuleBlockedPage           string
	Floor                     string
	PricingRule               string
	PricingRuleJobs           string
	Target                    string
	TargetV2                  string
	Playlist                  string
	PlaylistConfig            string
	RlPlaylistContent         string
	Template                  string
	TestAdTag                 string
	Config                    string
	Currency                  string
	AdType                    string
	MediaType                 string
	BidderTemplate            string
	LineItemBidderInfo        string
	LineItemBidderInfoV2      string
	LineItemAdsenseAdSlot     string
	UserBilling               string
	UserForgetPassword        string
	Gam                       string
	GamNetwork                string
	ModuleUserId              string
	InventoryConfig           string
	IdentityModuleInfo        string
	MissingAdsTxt             string
	PageCollapse              string
	Blocking                  string
	BlockingRestrictions      string
	RlBlockingInventory       string
	LineItemAccount           string
	PaymentInvoice            string
	PaymentRequest            string
	PaymentSubPub             string
	LogErrorPushLineItem      string
	LogErrorPushLineItemDfp   string
	LogErrorWorker            string
	SystemSetting             string
	RlsBidderSystemInventory  string
	UserPermission            string
	RlUserPermission          string
	BidderTemplateParams      string
	BidderParams              string
	LineItemBidderParams      string
	LineItemBidderParamsV2    string
	ListParamValidate         string
	PrebidModule              string
	Notification              string
	DfpLineItem               string
	DfpAdUnit                 string
	Category                  string
	ContentKeyword            string
	ContentTag                string
	ContentAdBreak            string
	ContentQuestion           string
	DfpCreative               string
	DfpLica                   string
	PbBidder                  string
	PbBidderParam             string
	Identity                  string
	ListFileJs                string
	ListFileJsVpaid           string
	History                   string
	HistoryDetail             string
	Language                  string
	AbTesting                 string
	RlsConnectionMCM          string
	RateSharingInventory      string
	RateSharing               string
	RevenueShare              string
	RequestApiCloudflare      string
	InventoryConnectionDemand string
	DynamicFloor              string
	DynamicFloorTag           string
	LogCpmAmz                 string
	KeyValueGam               string
	BlockRPM                  string
	TestCaseGroup             string
	TestCaseItem              string
	TestCaseProcess           string
	ContentTopArticles        string
	CronjobDomain             string
	Setting                   string
	LogIpRange                string
	ReportMonitor             string
	QizPosts                  string
	QizUsers                  string
	PricingRules              string
	InventoryAd               string
	CronJobBlockedPage        string
	FloorTargeting            string
	Cronjob                   string
	ManagerSub               string
}

var Tables tables
var Client *gorm.DB
var Crawler *gorm.DB
var DBQuiz *gorm.DB

func init() {
	/**
	Mysql Table
	*/
	Tables = tables{
		User:                      "user",
		UserInfo:                  "user_info",
		Inventory:                 "inventory",
		InventoryAdTag:            "inventory_ad_tag",
		TagRewardedAdConfig:       "tag_rewarded_ad_config",
		RlAdTagSizeAdditional:     "rls_adtag_size_additional",
		TagInfo:                   "tag_info",
		Country:                   "country",
		Channels:                  "channels",
		ChannelsKeyword:           "channels_keyword",
		Content:                   "content",
		Device:                    "devices",
		AdSize:                    "ad_size",
		Bidder:                    "bidder",
		RlBidderMediaType:         "rls_bidder_media_type",
		LineItem:                  "line_item",
		LineItemV2:                "line_item_v2",
		BidderAssign:              "bidder_assign",
		AdTag:                     "ad_tag",
		Rule:                      "rule",
		RuleBlockedPage:           "rule_blocked-page",
		Floor:                     "floor",
		PricingRule:               "new_pricing_rules",
		PricingRuleJobs:           "new_pricing_rule_jobs",
		Target:                    "target",
		TargetV2:                  "target_v2",
		Playlist:                  "playlist",
		PlaylistConfig:            "playlist_config",
		RlPlaylistContent:         "rls_playlist_content",
		Template:                  "template",
		TestAdTag:                 "test_adtag",
		Config:                    "config",
		Currency:                  "currency",
		AdType:                    "ad_type",
		MediaType:                 "media_type",
		BidderTemplate:            "bidder_template",
		LineItemBidderInfo:        "line_item_bidder_info",
		LineItemBidderInfoV2:      "line_item_bidder_info_v2",
		LineItemAdsenseAdSlot:     "line_item_adsense_ad_slot",
		UserBilling:               "user_billing",
		UserForgetPassword:        "user_forget_password",
		Gam:                       "gam",
		GamNetwork:                "gam_network",
		ModuleUserId:              "module_userid",
		InventoryConfig:           "inventory_config",
		IdentityModuleInfo:        "identity_module_info",
		MissingAdsTxt:             "missing_ads_txt",
		PageCollapse:              "page_collapse",
		Blocking:                  "blocking",
		BlockingRestrictions:      "blocking_restrictions",
		RlBlockingInventory:       "rls_blocking_inventory",
		LineItemAccount:           "line_item_account",
		PaymentInvoice:            "payment_invoice",
		PaymentRequest:            "payment_request",
		PaymentSubPub:             "payment_sub_pub",
		LogErrorPushLineItem:      "log_error_push_line_item",
		LogErrorPushLineItemDfp:   "log_error_push_line_item_dfp",
		LogErrorWorker:            "log_error_worker",
		SystemSetting:             "system_setting",
		RlsBidderSystemInventory:  "rls_bidder_system_inventory",
		UserPermission:            "user_permission",
		RlUserPermission:          "rls_user_permission",
		BidderTemplateParams:      "bidder_template_params",
		BidderParams:              "bidder_params",
		LineItemBidderParams:      "line_item_bidder_params",
		LineItemBidderParamsV2:    "line_item_bidder_params_v2",
		ListParamValidate:         "list_param_validate",
		PrebidModule:              "prebid_module",
		Notification:              "notification",
		DfpLineItem:               "dfp_line_item",
		DfpAdUnit:                 "dfp_ad_unit",
		Category:                  "category",
		ContentKeyword:            "content_keyword",
		ContentTag:                "content_tag",
		ContentAdBreak:            "content_ad_break",
		ContentQuestion:           "content_question",
		DfpCreative:               "dfp_creative",
		DfpLica:                   "dfp_lica",
		PbBidder:                  "pb_bidder",
		PbBidderParam:             "pb_bidder_param",
		Identity:                  "identity",
		ListFileJs:                "list_file_js",
		ListFileJsVpaid:           "list_file_js_vpaid",
		History:                   "history",
		HistoryDetail:             "history_detail",
		Language:                  "language",
		AbTesting:                 "ab_testing",
		RlsConnectionMCM:          "rls_connection_mcm",
		RateSharing:               "ratesharing",
		RateSharingInventory:      "rate_sharing",
		RevenueShare:              "revenue_share",
		RequestApiCloudflare:      "request_api_cloudflare",
		InventoryConnectionDemand: "inventory_connection_demand",
		DynamicFloor:              "dynamic_floor",
		DynamicFloorTag:           "dynamic_floor_tag",
		LogCpmAmz:                 "log_cpm_amz",
		KeyValueGam:               "key_value_gam",
		BlockRPM:                  "block_rpm",
		TestCaseGroup:             "test_case_group",
		TestCaseItem:              "test_case_item",
		TestCaseProcess:           "test_case_process",
		ContentTopArticles:        "content_top_articles",
		CronjobDomain:             "cronjob_domain",
		Setting:                   "setting",
		LogIpRange:                "log_ip_range",
		ReportMonitor:             "report_monitor",
		QizPosts:                  "qiz_posts",
		QizUsers:                  "qiz_users",
		PricingRules:              "pricing_rules",
		InventoryAd:               "inventory_ad",
		CronJobBlockedPage:        "cronjob_blocked_page",
		FloorTargeting:            "floor_targeting",
		Cronjob:                   "cronjob",
		ManagerSub:                "manager_sub",
	}
	Client = Connect()
	Crawler = ConnectCrawler()
	DBQuiz = ConnectDBQuiz()
}

// Connect connect mysql
//
// return:
func Connect() (DB *gorm.DB) {
	/**
	Mysql Config
	*/
	var Config config
	if utility.IsWindow() {
		Config = config{
			Username: "apacadmin",
			Password: "iK29&6%!9XKjs@",
			Database: "apac_ss",
			Host:     "localhost",
			Port:     "3306",
			Encoding: "utf8mb4",
		}
	} else {
		Config = config{
			Username: "apacadmin",
			Password: "iK29&6%!9XKjs@",
			Database: "apac_ss",
			// Host:     "localhost",
			// Host:     "23.92.74.98", //=> PublicIP server backend
			Host:     "192.168.9.10", // => PrivateIP server backend
			Port:     "9682",
			Encoding: "utf8mb4",
		}
	}
	if !utility.IsWindow() && utility.IsDev() {
		Config.Database = "apac_ss_dev"
	}
	if utility.IsDemo() {
		if !utility.IsWindow() {
			Config.Database = "apac_ss_demo"
		} else {
			Config.Database = "apac_ss"
		}
	}

	envFilename := os.Getenv("MODE") + ".env"
	_, err := os.Stat(envFilename)
	if err == nil {
		envs, err := godotenv.Read(envFilename)
		if err != nil {
			panic(errors.New("Error loading env file" + err.Error()))
		}
		Config = config{
			Username: envs["MYSQL_USERNAME"],
			Password: envs["MYSQL_PASSWORD"],
			Host:     envs["MYSQL_HOST"],
			Port:     envs["MYSQL_PORT"],
			Encoding: envs["MYSQL_ENCODING"],
			Database: envs["MYSQL_DATABASE_MAIN"],
		}
	}
	if helpers.IsHaiMode() {
		Config = haiSqlConfig
	}

	var ErrorConnect error
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", Config.Username, Config.Password, Config.Host, Config.Port, Config.Database, Config.Encoding)

	DB, ErrorConnect = gorm.Open(mysql.New(mysql.Config{
		DSN: dns, // data source name
		// DefaultStringSize:         256,                                                                        // default size for string fields
		// DisableDatetimePrecision:  true,                                                                       // disable datetime precision, which not supported before MySQL 5.6
		// DontSupportRenameIndex:    true,                                                                       // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		// DontSupportRenameColumn:   true,                                                                       // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		// SkipInitializeWithVersion: false,                                                                      // auto configure based on currently MySQL version
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Off log
		// Logger: logger.Default.LogMode(logger.Info), // On Notice for DEBUG
		// Logger: logger.Default.LogMode(logger.Warn), // On Notice log
	})
	if ErrorConnect != nil {
		panic(ErrorConnect)
	}

	sqlDB, _ := DB.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(1)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	// sqlDB.SetMaxOpenConns(1500)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// sqlDB.Stats()

	return
}

func ConnectCrawler() (DB *gorm.DB) {
	var crawlerConfig config
	if utility.IsWindow() {
		crawlerConfig = config{
			Username: "apacadmin",
			Password: "iK29&6%!9XKjs@",
			Database: "apac_ss",
			Host:     "localhost",
			Port:     "3306",
			Encoding: "utf8mb4",
		}
	} else {
		crawlerConfig = config{
			Username: "crawler",
			Password: "2VNHZGCftQRY2gV6",
			Database: "crawler",
			Host:     "192.168.9.10", // => PrivateIP server backend
			Port:     "9682",
			Encoding: "utf8mb4",
		}
	}

	envFilename := os.Getenv("MODE") + ".env"
	_, err := os.Stat(envFilename)
	if err == nil {
		envs, err := godotenv.Read(envFilename)
		if err != nil {
			panic(errors.New("Error loading env file" + err.Error()))
		}
		crawlerConfig = config{
			Username: envs["MYSQL_USERNAME"],
			Password: envs["MYSQL_PASSWORD"],
			Host:     envs["MYSQL_HOST"],
			Port:     envs["MYSQL_PORT"],
			Encoding: envs["MYSQL_ENCODING"],
			Database: envs["MYSQL_DATABASE_CRAWLER"],
		}
	}

	if helpers.IsHaiMode() {
		crawlerConfig = haiSqlConfig
	}

	var ErrorConnect error
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", crawlerConfig.Username, crawlerConfig.Password, crawlerConfig.Host, crawlerConfig.Port, crawlerConfig.Database, crawlerConfig.Encoding)
	DB, ErrorConnect = gorm.Open(mysql.New(mysql.Config{
		DSN: dns, // data source name
		// DefaultStringSize:         256,                                                                        // default size for string fields
		// DisableDatetimePrecision:  true,                                                                       // disable datetime precision, which not supported before MySQL 5.6
		// DontSupportRenameIndex:    true,                                                                       // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		// DontSupportRenameColumn:   true,                                                                       // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		// SkipInitializeWithVersion: false,                                                                      // auto configure based on currently MySQL version
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Off log
		// Logger: logger.Default.LogMode(logger.Info), // On Notice for DEBUG
		// Logger: logger.Default.LogMode(logger.Warn), // On Notice log
	})
	if ErrorConnect != nil {
		panic(ErrorConnect)
	}

	sqlDB, _ := DB.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(1)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	// sqlDB.SetMaxOpenConns(1500)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// sqlDB.Stats()

	return
}

func ConnectDBQuiz() (DB *gorm.DB) {
	var crawlerConfig config
	if utility.IsWindow() {
		crawlerConfig = config{
			Username: "apacadmin",
			Password: "iK29&6%!9XKjs@",
			Database: "apac_ss",
			Host:     "localhost",
			Port:     "3306",
			Encoding: "utf8mb4",
		}
	} else {
		crawlerConfig = config{
			Username: "wpadmin",
			Password: "XAIIDS8*@*#(!!009s",
			Database: "wp_quiz",
			Host:     "192.168.9.10", // => PrivateIP server backend
			Port:     "9682",
			Encoding: "utf8mb4",
		}
	}

	envFilename := os.Getenv("MODE") + ".env"
	_, err := os.Stat(envFilename)
	if err == nil {
		envs, err := godotenv.Read(envFilename)
		if err != nil {
			panic(errors.New("Error loading env file" + err.Error()))
		}
		crawlerConfig = config{
			Username: envs["MYSQL_USERNAME"],
			Password: envs["MYSQL_PASSWORD"],
			Host:     envs["MYSQL_HOST"],
			Port:     envs["MYSQL_PORT"],
			Encoding: envs["MYSQL_ENCODING"],
			Database: envs["MYSQL_DATABASE_QUIZ"],
		}
	}

	if helpers.IsHaiMode() {
		crawlerConfig = haiSqlConfig
	}
	var ErrorConnect error
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", crawlerConfig.Username, crawlerConfig.Password, crawlerConfig.Host, crawlerConfig.Port, crawlerConfig.Database, crawlerConfig.Encoding)
	DB, ErrorConnect = gorm.Open(mysql.New(mysql.Config{
		DSN: dns, // data source name
		// DefaultStringSize:         256,                                                                        // default size for string fields
		// DisableDatetimePrecision:  true,                                                                       // disable datetime precision, which not supported before MySQL 5.6
		// DontSupportRenameIndex:    true,                                                                       // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		// DontSupportRenameColumn:   true,                                                                       // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		// SkipInitializeWithVersion: false,                                                                      // auto configure based on currently MySQL version
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Off log
		// Logger: logger.Default.LogMode(logger.Info), // On Notice for DEBUG
		// Logger: logger.Default.LogMode(logger.Warn), // On Notice log
	})
	if ErrorConnect != nil {
		panic(ErrorConnect)
	}

	sqlDB, _ := DB.DB()
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(1)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	// sqlDB.SetMaxOpenConns(1500)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// sqlDB.Stats()

	return
}

type ArrayInt64 []int64

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *ArrayInt64) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var result []int64
	err := json.Unmarshal(bytes, &result)
	*j = ArrayInt64(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (j ArrayInt64) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

type ArrayString []string

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *ArrayString) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	var result []string
	err := json.Unmarshal(bytes, &result)
	*j = ArrayString(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (j ArrayString) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

type TYPEOption int

const (
	OptionInclude = iota + 1
	OptionExclude
)

func (t TYPEOption) String() string {
	switch t {
	case OptionInclude:
		return "include"
	case OptionExclude:
		return "exclude"
	default:
		return ""
	}
}

func (t TYPEOption) Standardized() TYPEOption {
	if t == 1 {
		return OptionInclude
	}
	return OptionExclude
}

type TYPEStatus int

const (
	StatusApproved TYPEStatus = iota + 1
	StatusPending
	StatusReject
)

func (s TYPEStatus) String() string {
	switch s {
	case StatusApproved:
		return "approved"
	case StatusPending:
		return "pending"
	case StatusReject:
		return "rejected"
	default:
		return ""
	}
}

func (s TYPEStatus) Int() int {
	switch s {
	case StatusApproved:
		return 1
	case StatusPending:
		return 2
	case StatusReject:
		return 3
	default:
		return 0
	}
}

type TYPEOnOff int

const (
	On TYPEOnOff = iota + 1
	Off
)

func (s TYPEOnOff) String() string {
	switch s {
	case On:
		return "on"
	case Off:
		return "off"
	default:
		return ""
	}
}

func (s TYPEOnOff) Boolean() bool {
	switch s {
	case On:
		return true
	case Off:
		return false
	default:
		return false
	}
}

func (s TYPEOnOff) StringRunPause() string {
	switch s {
	case On:
		return "running"
	case Off:
		return "paused"
	default:
		return ""
	}
}

func OffLog() {
	Client.Config.Logger = logger.Default.LogMode(logger.Silent)
}
