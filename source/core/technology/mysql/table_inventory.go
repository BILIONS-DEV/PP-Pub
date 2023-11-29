package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"source/pkg/adstxt"
	"source/pkg/utility"
	"strings"
	"time"
)

type TableInventory struct {
	Id                  int64                          `gorm:"column:id" json:"id"`
	UserId              int64                          `gorm:"column:user_id" json:"user_id"`
	Presenter           int64                          `gorm:"column:presenter" json:"presenter"`
	Name                string                         `gorm:"column:name" json:"name"`
	Domain              string                         `gorm:"column:domain" json:"domain"`
	Type                TYPEInventoryType              `gorm:"column:type" json:"type"`
	Status              TYPEStatus                     `gorm:"column:status" json:"status"`
	JsMode              string                         `gorm:"column:js_mode" json:"js_mode"`
	VpaidMode           string                         `gorm:"column:vpaid_mode" json:"vpaid_mode"`
	PrebidJs            string                         `gorm:"column:prebid_js" json:"prebid_js"`
	SyncAdsTxt          TYPEInventorySyncAdsTxt        `gorm:"column:sync_ads_txt" json:"sync_ads_txt"`
	LastScanAdsTxt      sql.NullTime                   `gorm:"column:last_scan_ads_txt" json:"last_scan_ads_txt"`
	AdsTxtUrl           string                         `gorm:"column:ads_txt_url" json:"ads_txt_url"`
	AdsTxtCustom        TYPEAdsTxtCustom               `gorm:"column:ads_txt_custom" json:"ads_txt_custom"`
	AdsTxtCustomByAdmin TYPEAdsTxtCustom               `gorm:"column:ads_txt_custom_by_admin" json:"ads_txt_custom_by_admin"`
	CreatedAt           time.Time                      `gorm:"column:created_at" json:"created_at"`
	DeletedAt           gorm.DeletedAt                 `gorm:"column:deleted_at" json:"deleted_at"`
	IabCategories       string                         `gorm:"column:iab_categories" json:"iab_categories"`
	Uuid                string                         `gorm:"column:uuid" json:"uuid"`
	Requests            int64                          `gorm:"column:requests" json:"requests"`
	Impressions         int64                          `gorm:"column:impressions" json:"impressions"`
	Revenue             float64                        `gorm:"column:revenue" json:"revenue"`
	ApacSiteId          int64                          `gorm:"column:apac_siteid" json:"apac_siteid"`
	CachedAt            time.Time                      `gorm:"column:cached_at"  json:"cached_at"`
	RenderCache         int                            `gorm:"column:render_cache" json:"render_cache"`
	SyncPocPoc          string                         `gorm:"default:pending;column:sync_pocpoc" json:"sync_pocpoc"`
	Config              TableInventoryConfig           `gorm:"-"`
	AdTag               []TableInventoryAdTag          `gorm:"-"`
	UserIdModule        []TableIdentityModuleInfo      `gorm:"-"`
	RlsBidderSystem     TableRlsBidderSystemInventory  `gorm:"-"` // Rls để check chang status trong BE
	RateSharing         TableRateSharingInventory      `gorm:"-"`
	ConnectionDemand    TableInventoryConnectionDemand `gorm:"-"`
}

func (TableInventory) TableName() string {
	return Tables.Inventory
}

func (rec *TableInventory) GetFullData() {
	if !rec.IsFound() {
		return
	}
	rec.GetConfig()
	rec.GetAdTag(true)
	rec.GetUserIdModules()
	rec.GetRateSharing()
}

func (rec *TableInventory) GetConfig() {
	if !rec.IsFound() {
		return
	}
	var config TableInventoryConfig
	Client.Where("inventory_id = ?", rec.Id).Find(&config)
	rec.Config = config
}

func (rec *TableInventory) GetAdTag(fullData bool) {
	if !rec.IsFound() {
		return
	}
	var adTag []TableInventoryAdTag
	Client.Where("inventory_id = ?", rec.Id).Find(&adTag)
	rec.AdTag = adTag
	if fullData {
		for k, adTag := range rec.AdTag {
			if adTag.Id == 57 {
				adTag.GetFullData()
			}
			rec.AdTag[k] = adTag
		}
	}
}

func (rec *TableInventory) GetUserIdModules() {
	if !rec.IsFound() {
		return
	}
	var targets []TableTarget
	var listIdentityId []int64
	Client.Where("(inventory_id = -1 OR inventory_id = ?) AND identity_id != 0 and user_id = ?", rec.Id, rec.UserId).Find(&targets)
	for _, target := range targets {
		listIdentityId = append(listIdentityId, target.IdentityId)
	}
	var identities []TableIdentity
	var identityChoose TableIdentity
	var priority = 17
	Client.Where("id IN ?", listIdentityId).Find(&identities)
	for _, identity := range identities {
		if identity.Priority < priority {
			identityChoose = identity
		}
	}

	var userIdModule []TableIdentityModuleInfo
	Client.Where("identity_id = ?", identityChoose.Id).Find(&userIdModule)
	rec.UserIdModule = userIdModule
}

func (rec *TableInventory) GetRateSharing() {
	if !rec.IsFound() {
		return
	}
	var rateSharing TableRateSharingInventory
	Client.Where("inventory_id = ?", rec.Id).Last(&rateSharing)
	rec.RateSharing = rateSharing
}

func (rec *TableInventory) GetAdsTxt() (adsTxt []string, err error) {
	if !rec.IsFound() {
		return
	}

	// Case Ads Txt đặc biệt
	//adsTxt = append(adsTxt, "\n")
	adsTxt = append(adsTxt, "OWNERDOMAIN="+rec.Name)
	adsTxt = append(adsTxt, "MANAGERDOMAIN=pubpower.io")

	// Get ads txt từ bidder của pub
	var bidderGoogles []TableBidder
	var bidderPrebids []TableBidder
	Client.Where("user_id = ? AND ads_txt != ? and bidder_template_id != 1", rec.UserId, "").Find(&bidderPrebids) // Lấy list bidder prebid
	if len(bidderPrebids) > 0 {
		for _, bidder := range bidderPrebids {
			lines := utility.SplitLines(bidder.AdsTxt)
			if len(lines) > 0 {
				adsTxt = append(adsTxt, "\n#"+cases.Title(language.Und, cases.NoLower).String(bidder.BidderCode))
			}
			for _, line := range lines {
				adsTxt = append(adsTxt, line)
			}
		}
	}
	Client.Where("user_id = ? AND ads_txt != ? and bidder_template_id = 1", rec.UserId, "").Find(&bidderGoogles) // Lấy list bidder google
	if len(bidderGoogles) > 0 {
		adsTxt = append(adsTxt, "\n#Google")
		for _, bidder := range bidderGoogles {
			lines := utility.SplitLines(bidder.AdsTxt)
			for _, line := range lines {
				adsTxt = append(adsTxt, line)
			}
		}
	}

	// Get ads txt từ bidder của admin
	Client.Where("user_id = 0 and status = 1 and bidder_template_id != 1").Find(&bidderPrebids) // Lấy list bidder prebid
	if len(bidderPrebids) > 0 {
		for _, bidder := range bidderPrebids {
			adsTxtBidder, err := rec.getAdsTxtFromBidderSystem(bidder)
			if err != nil {
				continue
			}
			if len(adsTxtBidder) > 0 {
				adsTxt = append(adsTxt, "\n#"+cases.Title(language.Und, cases.NoLower).String(bidder.BidderCode))
				adsTxt = append(adsTxt, adsTxtBidder...)
			}
		}
	}
	Client.Where("user_id = 0 and status = 1 and bidder_template_id = 1").Find(&bidderGoogles) // Lấy list bidder google
	if len(bidderGoogles) > 0 {
		adsTxt = append(adsTxt, "\n#Google")
		for _, bidder := range bidderGoogles {
			adsTxtBidder, err := rec.getAdsTxtFromBidderSystem(bidder)
			if err != nil {
				continue
			}
			if len(adsTxtBidder) > 0 {
				adsTxt = append(adsTxt, adsTxtBidder...)
			}
		}
	}

	// Add ads txt từ inventory
	adsTxt = append(adsTxt, "\n")
	adsTxt = append(adsTxt, rec.AdsTxtCustomByAdmin.ToArray()...) // Admin set
	adsTxt = append(adsTxt, rec.AdsTxtCustom.ToArray()...)        // Pub set

	// Add ads txt từ setup của user
	var user TableUser
	Client.Where("id = ?", rec.UserId).Find(&user)
	adsTxt = append(adsTxt, "\n")
	adsTxt = append(adsTxt, user.AdsTxtCustomByAdmin.ToArray()...)

	// ads.txt apacdex
	if user.SellerID > 0 && user.DisabledSeller == "" {
		adsTxt = append(adsTxt, "\n")
		adsTxt = append(adsTxt, fmt.Sprintf("quantumdex.io, %d, DIRECT", user.SellerID))
		adsTxt = append(adsTxt, fmt.Sprintf("interdogmedia.com, %d, DIRECT", user.SellerID))
		adsTxt = append(adsTxt, fmt.Sprintf("apacdex.com, %d, DIRECT", user.SellerID))
		adsTxt = append(adsTxt, fmt.Sprintf("pubpower.io , %d, DIRECT", user.SellerID))
	}
	// ads.txt từ setting
	var setting TableSetting
	Client.Where("meta_key = 'ads_txt'").Find(&setting)
	settingAdsTxt := TYPEAdsTxtCustom(setting.MetaValue)
	adsTxt = append(adsTxt, "\n")
	adsTxt = append(adsTxt, settingAdsTxt.ToArray()...)

	// Chuẩn hóa ads txt
	adsTxt = adstxt.Standardized(adsTxt)
	return
}

func (rec *TableInventory) getAdsTxtFromBidderSystem(bidder TableBidder) (adsTxt []string, err error) {
	if govalidator.IsNull(bidder.AdsTxt) {
		return
	}
	// Từ bidder system lấy ra các status accept để thêm adstxt
	var listStatusAdsTxt []string
	_ = json.Unmarshal([]byte(bidder.StatusAdsTxt), &listStatusAdsTxt)

	// Lấy ra các status accept của bidder để in ra Ads Txt
	var listStatusBidderAccept []TYPERlsBidderSystemInventoryStatus
	for _, statusAdsTxt := range listStatusAdsTxt {
		// Từ các status đã chọn này đổi sang TYPE của status RlsBidderSystemInventory
		statusBidderAccept := TYPEStatusAdsTxt(statusAdsTxt).StatusBidder()
		listStatusBidderAccept = append(listStatusBidderAccept, statusBidderAccept...)
	}

	// Từ list bidder status accept check trong db xem status có đang đúng không
	if len(listStatusBidderAccept) > 0 {
		var rlsBidderStatus TableRlsBidderSystemInventory
		Client.
			//Debug().
			Where("bidder_id = ? AND status in ? AND inventory_name = ?", bidder.Id, listStatusBidderAccept, rec.Name).
			Last(&rlsBidderStatus)

		// Nếu như không tồn tại status được accept bỏ qua bidder này
		if rlsBidderStatus.Id == 0 {
			return
		}
	}

	// Nếu như bidder được accept thì thêm ads txt vào
	lines := utility.SplitLines(bidder.AdsTxt)
	//fmt.Printf("%+v \n", lines)
	for _, line := range lines {
		adsTxt = append(adsTxt, line)
	}

	return
}

func (rec *TableInventory) IsFound() bool {
	if rec.Id > 0 {
		return true
	}
	return false
}

type TYPEAdsTxtCustom string

func (t TYPEAdsTxtCustom) ToArray() (array []string) {
	array = strings.Split(string(t), "\n")
	return
}

type TYPEInventoryType int

const (
	InventoryTypeWeb = iota + 1
	InventoryTypeApp
)

func (t TYPEInventoryType) String() string {
	switch t {
	case InventoryTypeWeb:
		return "web"
	case InventoryTypeApp:
		return "app"
	default:
		return ""
	}
}

type TYPEInventorySyncAdsTxt int

const (
	InventorySyncAdsTxt = iota + 1
	InventorySyncAdsTxtNotIn
	InventorySyncAdsTxtError
)

func (t TYPEInventorySyncAdsTxt) String() string {
	switch t {
	case InventorySyncAdsTxt:
		return "In Sync"
	case InventorySyncAdsTxtNotIn:
		return "Missing Line"
	case InventorySyncAdsTxtError:
		return "Does Not Exist"
	default:
		return ""
	}
}

func (rec TableInventory) GetLinkJS(user TableUser) (URLAsynchronous string, URLNormal string, URLAutoAds string) {
	tag := "powerTag"
	tagDomain := "nc.pubpowerplatform.io"
	if user.Id != 0 && user.SystemSync == 1 { // system VLI
		tagDomain = "cdn.vlitag.com"
		tag = "vitag"
	}
	URLNormal = "<script type=\"text/javascript\" src=\"//" + tagDomain + "/w/" + rec.Uuid + ".js\" async defer></script><script>var " + tag + " = " + tag + " || {};" + tag + ".gdprShowConsentToolButton = false;</script>"
	URLAsynchronous = "<script type=\"text/javascript\" src=\"//" + tagDomain + "/w/" + rec.Uuid + ".js\" async defer></script><script>var " + tag + " = " + tag + " || {};" + tag + ".gdprShowConsentToolButton = false;</script>"

	URLAutoAds = "<script type=\"text/javascript\" src=\"//" + tagDomain + "/ata/adv/" + rec.Uuid + ".js\" async defer></script>"
	return
}
