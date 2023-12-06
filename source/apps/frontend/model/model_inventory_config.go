package model

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"source/apps/frontend/lang"
	"source/apps/frontend/payload"
	"source/apps/history"
	"source/core/technology/mysql"
	"source/pkg/ajax"
	"source/pkg/utility"
	"strconv"
	// "fmt"
)

type InventoryConfig struct{}

type InventoryConfigRecord struct {
	mysql.TableInventoryConfig
}

func (InventoryConfigRecord) TableName() string {
	return mysql.Tables.InventoryConfig
}

func (t *InventoryConfig) SetupInventoryConfig(inputs payload.GeneralInventory, value string, userId int64, userAdmin UserRecord) (errs []ajax.Error) {
	lang := lang.Translate
	inventory, _ := new(Inventory).GetById(inputs.InventoryId, userId)
	recordOld := t.VerificationRecord(inputs.InventoryId)
	inventory.Config = recordOld.TableInventoryConfig
	if recordOld.Id == 0 {
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: "You don't own this inventory",
		})
	}
	errs = t.ValidateConfig(inputs, value, userId)
	if len(errs) > 0 {
		return
	}
	inventoryNew := inventory
	var recordNew InventoryConfigRecord
	recordNew = recordOld
	updateValue := recordNew.makeRowConfig(inputs, value)
	err := t.UpdateRow(recordNew, updateValue)
	// fmt.Printf("%+v\n", err)
	if err != nil {
		if !utility.IsWindow() {
			if value == "consent" {
				errs = append(errs, ajax.Error{
					Id:      "",
					Message: lang.Errors.InventoryError.SetupConsent.ToString(),
				})
			} else {
				errs = append(errs, ajax.Error{
					Id:      "",
					Message: lang.Errors.InventoryError.SetupConfig.ToString(),
				})
			}
			return
		}
		errs = append(errs, ajax.Error{
			Id:      "",
			Message: err.Error(),
		})
		return
	}
	// Push history
	inventoryNew.Config = recordNew.TableInventoryConfig
	var creatorId int64
	if userAdmin.Id != 0 {
		creatorId = userAdmin.Id
	} else {
		creatorId = userId
	}
	switch value {
	case "config":
		err = history.PushHistory(&history.Inventory{
			Detail:    history.DetailInventoryConfigFE,
			CreatorId: creatorId,
			RecordOld: inventory.TableInventory,
			RecordNew: inventoryNew.TableInventory,
		})
		break
	case "consent":
		err = history.PushHistory(&history.Inventory{
			Detail:    history.DetailInventoryConsentFE,
			CreatorId: creatorId,
			RecordOld: inventory.TableInventory,
			RecordNew: inventoryNew.TableInventory,
		})
		break
	}
	new(Inventory).ResetCacheWorker(inputs.InventoryId)
	return
}

func (t *InventoryConfig) ValidateModule(inputs payload.GeneralInventory) (errs []ajax.Error) {
	// var listIdTypeClient []int64
	for _, ModuleInfo := range inputs.ModuleParams {
		module := new(ModuleUserId).GetById(ModuleInfo.ModuleId)
		var Param []payload.ParamModuleUserId
		var Storage []payload.StorageModuleUserId
		json.Unmarshal([]byte(module.Params), &Param)
		json.Unmarshal([]byte(module.Storage), &Storage)

		for _, ModuleParam := range ModuleInfo.Params {
			idParam := strconv.FormatInt(ModuleInfo.ModuleId, 10) + "-" + ModuleParam.Name
			value := ModuleParam.Template

			for _, _Param := range Param {
				switch _Param.Type {
				case "int":
					if value != "" && ModuleParam.Name == _Param.Name {
						if !govalidator.IsInt(ModuleParam.Template) {
							errs = append(errs, ajax.Error{
								Id:      idParam,
								Message: "param " + ModuleParam.Name + " value is int",
							})
						}
					}
					break
				case "float":
					if value != "" && ModuleParam.Name == _Param.Name {
						if !govalidator.IsFloat(ModuleParam.Template) {
							errs = append(errs, ajax.Error{
								Id:      idParam,
								Message: "param " + ModuleParam.Name + " value is float",
							})
						}
					}
					break
				case "json":
					if value != "" && ModuleParam.Name == _Param.Name {
						if !govalidator.IsJSON(ModuleParam.Template) {
							errs = append(errs, ajax.Error{
								Id:      idParam,
								Message: "param " + ModuleParam.Name + " value is json",
							})
						}
					}
				case "boolean":
					if value != "" && ModuleParam.Name == _Param.Name {
						errs = append(errs, ajax.Error{
							Id:      idParam,
							Message: "param " + ModuleParam.Name + " value is true and false",
						})
					}
					break
				}
			}
		}

		for _, ModuleStorage := range ModuleInfo.Storage {
			idStorage := strconv.FormatInt(ModuleInfo.ModuleId, 10) + "-" + ModuleStorage.Name
			value := ModuleStorage.Template

			for _, _Storage := range Storage {
				switch _Storage.Type {
				case "int":
					if value != "" && ModuleStorage.Name == _Storage.Name {
						if !govalidator.IsInt(ModuleStorage.Template) {
							errs = append(errs, ajax.Error{
								Id:      idStorage,
								Message: "storage " + ModuleStorage.Name + " value is int",
							})
						}
					}
					break
				case "float":
					if value != "" && ModuleStorage.Name == _Storage.Name {
						if !govalidator.IsFloat(ModuleStorage.Template) {
							errs = append(errs, ajax.Error{
								Id:      idStorage,
								Message: "storage " + ModuleStorage.Name + " value is float",
							})
						}
					}
					break
				case "json":
					if value != "" && ModuleStorage.Name == _Storage.Name {
						if !govalidator.IsJSON(ModuleStorage.Template) {
							errs = append(errs, ajax.Error{
								Id:      idStorage,
								Message: "storage " + ModuleStorage.Name + " value is json",
							})
						}
					}
				case "boolean":
					if value != "" && ModuleStorage.Name == _Storage.Name {
						errs = append(errs, ajax.Error{
							Id:      idStorage,
							Message: "storage " + ModuleStorage.Name + " value is true and false",
						})
					}
					break
				}
			}
		}
	}
	return
}

func (t *InventoryConfig) ValidateConfig(inputs payload.GeneralInventory, value string, userId int64) (errs []ajax.Error) {
	switch value {
	case "config":
		if inputs.PrebidTimeout == 0 || inputs.PrebidTimeout < 0 {
			errs = append(errs, ajax.Error{
				Id:      "prebid_timeout",
				Message: "PrebidTimeout is required",
			})
		} else if inputs.PrebidTimeout < 100 || inputs.PrebidTimeout > 5000 {
			errs = append(errs, ajax.Error{
				Id:      "prebid_timeout",
				Message: "min = 100, max 5000",
			})
		}

		if inputs.AdRefresh == 1 {
			if inputs.AdRefreshTime == 0 || inputs.AdRefreshTime < 0 {
				errs = append(errs, ajax.Error{
					Id:      "ad_refresh_time",
					Message: "Ad Refresh Time is required",
				})
			} else if inputs.AdRefreshTime < 20 {
				errs = append(errs, ajax.Error{
					Id:      "ad_refresh_time",
					Message: "min = 20",
				})
			}
		}
		if inputs.LoadAdType == "lazyload" {
			if inputs.RenderMarginPercent > inputs.FetchMarginPercent {
				errs = append(errs, ajax.Error{
					Id:      "box_render_margin_percent",
					Message: "Render Margin Percent <= Fetch Margin Percent",
				})
			}
		}
		break
	case "consent":
		if inputs.Gdpr == 1 {
			if inputs.GdprTimeout == 0 || inputs.GdprTimeout < 0 {
				errs = append(errs, ajax.Error{
					Id:      "gdpr_timeout",
					Message: "Gdpr Timeout is required",
				})
			}

			if inputs.GdprTimeout > 15000 {
				errs = append(errs, ajax.Error{
					Id:      "gdpr_timeout",
					Message: "Maximum Gdpr Timeout is 15000ms",
				})
			}
		}

		if inputs.Ccpa == 1 {
			if inputs.CcpaTimeout == 0 || inputs.CcpaTimeout < 0 {
				errs = append(errs, ajax.Error{
					Id:      "ccpa_timeout",
					Message: "Ccpa Timeout is required",
				})
			}

			if inputs.CcpaTimeout > 5000 {
				errs = append(errs, ajax.Error{
					Id:      "ccpa_timeout",
					Message: "Maximum Ccpa Timeout is 5000ms",
				})
			}
		}

		if inputs.CustomBrand == 1 {
			if utility.ValidateString(inputs.Logo) == "" {
				errs = append(errs, ajax.Error{
					Id:      "custom_brand_logo",
					Message: "Logo is required",
				})
			}
			if !utility.IsUrl(inputs.Logo) {
				errs = append(errs, ajax.Error{
					Id:      "custom_brand_logo",
					Message: "Logo must be a url",
				})
			}
			if utility.ValidateString(inputs.Title) == "" {
				errs = append(errs, ajax.Error{
					Id:      "custom_brand_title",
					Message: "Title is required",
				})
			}
			if utility.ValidateString(inputs.Content) == "" {
				errs = append(errs, ajax.Error{
					Id:      "custom_brand_content",
					Message: "Content is required",
				})
			}
		}
		break
	case "userid":
		if inputs.SyncDelay < 1 {
			errs = append(errs, ajax.Error{
				Id:      "sync_delay",
				Message: "Sync Delay is required",
			})
		}

		if inputs.AuctionDelay < 0 {
			errs = append(errs, ajax.Error{
				Id:      "auction_delay",
				Message: "Auction Delay must be positive",
			})
		}
		break
	}
	return
}

func (rec *InventoryConfigRecord) makeRowConfig(record payload.GeneralInventory, value string) (listUpdates map[string]interface{}) {
	var m = make(map[string]interface{})
	switch value {
	case "config":
		rec.PrebidTimeOut = record.PrebidTimeout
		rec.AdRefresh = record.AdRefresh
		rec.DirectSales = record.DirectSales
		rec.GamAccount = record.GamAccount
		rec.SafeFrame = record.SafeFrame
		rec.GamAutoCreate = record.GamAutoCreate
		rec.PbRenderMode = record.PbRenderMode
		rec.LoadAdType = record.LoadAdType
		rec.FetchMarginPercent = record.FetchMarginPercent
		rec.RenderMarginPercent = record.RenderMarginPercent
		rec.MobileScaling = record.MobileScaling
		if rec.AdRefresh == 1 {
			rec.AdRefreshTime = record.AdRefreshTime
			rec.LoadAdRefresh = record.LoadAdRefresh
		} else {
			rec.AdRefreshTime = 0
			// rec.LoadAdType = ""
			rec.LoadAdRefresh = ""
			m["ad_refresh_time"] = rec.AdRefreshTime
			// m["load_ad_type"] = rec.AdRefreshTime
			m["load_ad_refresh"] = rec.LoadAdRefresh
		}
		break
	case "consent":
		rec.Gdpr = record.Gdpr
		if rec.Gdpr == 1 {
			rec.GdprTimeout = record.GdprTimeout
		} else {
			rec.GdprTimeout = 0
			m["gdpr_timeout"] = rec.GdprTimeout
		}
		rec.Ccpa = record.Ccpa
		if rec.Ccpa == 1 {
			rec.CcpaTimeout = record.CcpaTimeout
		} else {
			rec.CcpaTimeout = 0
			m["ccpa_timeout"] = rec.CcpaTimeout
		}
		rec.CustomBrand = record.CustomBrand
		if rec.CustomBrand == 1 {
			rec.Logo = record.Logo
			rec.Title = record.Title
			rec.Content = record.Content
		} else {
			rec.Logo = ""
			rec.Title = ""
			rec.Content = ""
			m["logo"] = rec.Title
			m["title"] = rec.Logo
			m["content"] = rec.Content
		}
		break
	case "userid":
		rec.AuctionDelay = record.AuctionDelay
		rec.SyncDelay = record.SyncDelay
		m["auction_delay"] = rec.AuctionDelay
		break
	}
	listUpdates = m
	return
}

func (t *InventoryConfig) GetByInventoryId(id int64) (row InventoryConfigRecord) {
	mysql.Client.Where("inventory_id = ?", id).Last(&row)
	return
}

func (t *InventoryConfig) UpdateRow(record InventoryConfigRecord, valueUpdate map[string]interface{}) (err error) {
	if record.Id == 0 {
		err = mysql.Client.
			//Debug().
			Create(&record).Error
		return
	} else {
		err = mysql.Client.
			//Debug().
			Save(&record).
			Updates(valueUpdate).
			Where("id = ?", record.Id).Error
	}
	return
}

func (t *InventoryConfig) VerificationRecord(inventoryId int64) (row InventoryConfigRecord) {
	mysql.Client.
		//Debug().
		Model(&InventoryConfigRecord{}).Where("inventory_id = ?", inventoryId).Find(&row)
	return
}

func (t *InventoryConfig) MakeRowDefault(inventoryId int64) (cf InventoryConfigRecord) {
	var syncDelayDefaultValue = 3000
	var auctionDelayDefaultValue = 0
	cf = InventoryConfigRecord{mysql.TableInventoryConfig{
		InventoryId:         inventoryId,
		AdRefresh:           mysql.TypeOn,
		DirectSales:         mysql.TypeOn,
		AdRefreshTime:       30,
		PrebidTimeOut:       1000,
		LoadAdType:          "lazyload",
		LoadAdRefresh:       "signal reload",
		Gdpr:                1,
		GdprTimeout:         8000,
		Ccpa:                1,
		CcpaTimeout:         3000,
		CustomBrand:         2,
		SafeFrame:           1,
		Logo:                "",
		Title:               "We Value Your Privacy",
		Content:             "We and our partners store or access information on devices, such as cookies and process personal data, such as unique identifiers and standard information sent by a device for the purposes described below.You may click to consent to our and our partners’ processing for such purposes. Alternatively, you may click to refuse to consent, or access more detailed information and change your preferences before consenting.Your preferences will apply to this website only. Please note that some processing of your personal data may not require your consent, but you have a right to object to such processing. You can change your preferences at any time by returning to this site or visit our privacy policy.",
		AuctionDelay:        auctionDelayDefaultValue,
		SyncDelay:           syncDelayDefaultValue,
		GamAutoCreate:       1,
		PbRenderMode:        mysql.TYPEPbRenderModeInIframe,
		FetchMarginPercent:  500,
		RenderMarginPercent: 200,
		MobileScaling:       2,
	}}
	mysql.Client.
		Debug().
		Model(&InventoryConfigRecord{}).Create(&cf)
	return
}

func (t *InventoryConfig) DeleteConfigAfterDeleteDomain(domainId int64) {
	mysql.Client.Model(&InventoryConfigRecord{}).Delete(&InventoryConfigRecord{}, "inventory_id = ?", domainId)
}
