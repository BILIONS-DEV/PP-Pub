package model

import (
	"fmt"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"source/apps/frontend/ggapi"
	"source/apps/frontend/lang"
	"source/apps/frontend/payload"
	"source/apps/history"
	"source/core/technology/mysql"
	"source/pkg/datatable"
	"source/pkg/htmlblock"
	"source/pkg/pagination"
	"source/pkg/utility"
	"strings"
	"time"
)

type GamNetwork struct{}

type GamNetworkRecord struct {
	mysql.TableGamNetwork
}

func (GamNetworkRecord) TableName() string {
	return mysql.Tables.GamNetwork
}

type GamRecordDatatable struct {
	NetworkId   int64  `json:"network_id"`
	NetworkName string `json:"network_name"`
	ApiAccess   string `json:"api_access"`
	SetUp       string `json:"set_up"`
	CreateAt    string `json:"create_at"`
}

func (t *GamNetwork) GetByGam(userId int64, gamId int64) (recs []GamNetworkRecord) {
	mysql.Client.Where("user_id = ? AND gam_id = ?", userId, gamId).Find(&recs)
	return
}

func (t *GamNetwork) SelectByUser(inputs payload.GamSelectGam, user UserRecord, lang lang.Translation) (err error) {
	var network GamNetworkRecord
	if err = mysql.Client.First(&network, inputs.NetworkId).Error; err != nil {
		if !utility.IsWindow() {
			err = fmt.Errorf(lang.Errors.GamError.SelectNetwork.ToString())
			return
		}
		return err
	}
	//if err = mysql.Client.Model(&GamNetworkRecord{}).
	//	Where("user_id = ? AND gam_id = ?", user.Id, inputs.GamId).
	//	Update("status", mysql.StatusGamPending).Error; err != nil {
	//	return err
	//}
	var status mysql.TYPEStatusGam
	if inputs.Select {
		status = mysql.StatusGamSelected
	} else {
		status = mysql.StatusGamPending
	}
	if err = mysql.Client.Model(&network).Update("status", status).Error; err != nil {
		if !utility.IsWindow() {
			err = fmt.Errorf(lang.Errors.GamError.SelectNetwork.ToString())
			return
		}
		return err
	}
	return
}

func (t *GamNetwork) GetByGamId(gamId int64, userId int64) (recs []GamNetworkRecord) {
	mysql.Client.Where("gam_id = ? AND user_id = ?", gamId, userId).Find(&recs)
	return
}

func (t *GamNetwork) GetById(id int64, userId int64) (rec GamNetworkRecord) {
	mysql.Client.Where("id = ? AND user_id = ?", id, userId).Find(&rec)
	return
}

func (t *GamNetwork) GetByUser(userId int64) (records []GamNetworkRecord) {
	mysql.Client.Where("user_id = ? and api_access != 2", userId).Find(&records)
	return
}

func (t *GamNetwork) GetForSelect(userId int64) (records []GamNetworkRecord) {
	mysql.Client.Where("user_id = ? and api_access != 2 and network_id != '' and network_name != '' and gam_id != 0", userId).Find(&records)
	return
}

func (t *GamNetwork) GetListMCMAccept(userId int64) (records []GamNetworkRecord) {
	mysql.Client.Where("user_id = ? and connection_mcm = 1", userId).Find(&records)
	return
}

func (t *GamNetwork) GetByNetworkId(networkId int64, userId int64) (rec GamNetworkRecord) {
	mysql.Client.Where("user_id = ? AND network_id = ?", userId, networkId).Last(&rec)
	return
}

func (t *GamNetwork) CheckNetwork(networkId int64, gamId int64, userId int64) (rec GamNetworkRecord) {
	mysql.Client.Where("user_id = ? AND network_id = ? and gam_id = ?", userId, networkId, gamId).Last(&rec)
	return
}

func (t *GamNetwork) PushLineItem(networkId []int64, userId int64, userAdmin UserRecord, lang lang.Translation) (err error) {
	for _, v := range networkId {
		// Get recordOld
		recordOld := new(GamNetwork).GetById(v, userId)

		err = mysql.Client.Where("id = ? AND user_id = ? AND status = 2 AND api_access = 1 AND push_line_item != 1", v, userId).Updates(GamNetworkRecord{mysql.TableGamNetwork{
			PushLineItem:     3,
			DatePushLineItem: time.Now(),
		}}).Error
		if err != nil {
			if !utility.IsWindow() {
				err = fmt.Errorf(lang.Errors.GamError.PushLineItem.ToString())
				return
			}
		}
		// Get recordNew
		recordNew := new(GamNetwork).GetById(v, userId)
		// Push History
		var creatorId int64
		if userAdmin.Id != 0 {
			creatorId = userAdmin.Id
		} else {
			creatorId = userId
		}
		_ = history.PushHistory(&history.GAM{
			Detail:    history.DetailGAMSetupFE,
			CreatorId: creatorId,
			RecordOld: recordOld.TableGamNetwork,
			RecordNew: recordNew.TableGamNetwork,
		})
	}
	return
}

func (t *GamNetwork) Push(userId int64, userAdmin UserRecord, resp ggapi.Response, token *oauth2.Token) (err error) {
	gam, err := new(Gam).Push(userId, token, resp.User)
	if err != nil {
		return err
	}
	for _, network := range resp.Networks {
		isNetwork := t.CheckNetwork(network.Id, gam.Id, userId)
		if isNetwork.Id > 0 {
			continue
		}

		// Check api access
		var apiAccess mysql.TYPEApiAccess
		isEnable, err := ggapi.CheckAccessApi(gam.RefreshToken, network.Id, network.Name)
		if err != nil {
			continue
		}
		if !isEnable { // Nếu enable = false
			apiAccess = mysql.ApiAccessDisable
		} else {
			apiAccess = 0
		}

		rec := GamNetworkRecord{mysql.TableGamNetwork{
			UserId:           userId,
			GamId:            gam.Id,
			NetworkId:        network.Id,
			NetworkName:      network.Name,
			ApiAccess:        apiAccess,
			CurrencyCode:     network.CurrencyCode,
			TimeZone:         network.TimeZone,
			Status:           mysql.StatusGamPending,
			DatePushLineItem: time.Now(),
		}}
		mysql.Client.Create(&rec)

		// Push History
		var creatorId int64
		if userAdmin.Id != 0 {
			creatorId = userAdmin.Id
		} else {
			creatorId = userId
		}
		_ = history.PushHistory(&history.GAM{
			Detail:    history.DetailGAMSignInFE,
			CreatorId: creatorId,
			RecordOld: mysql.TableGamNetwork{},
			RecordNew: rec.TableGamNetwork,
		})
	}
	return
}

func (t *GamNetwork) GetByFilters(inputs *payload.GamFilterPayload, userId int64, lang lang.Translation) (response datatable.Response, err error) {
	var records []GamNetworkRecord
	var total int64
	err = mysql.Client.
		Where("user_id = ?", userId).
		Where("gam_id != 0").
		Scopes(
			t.setFilterSearch(inputs),
		).
		Model(&records).Count(&total).
		Scopes(
			t.setOrder(inputs),
			pagination.Paginate(pagination.Params{
				Limit:  inputs.Length,
				Offset: inputs.Start,
			}),
		).
		Find(&records).Error
	if err != nil {
		if !utility.IsWindow() {
			err = fmt.Errorf(lang.Errors.GamError.List.ToString())
		}
		return datatable.Response{}, err
	}
	response.Draw = inputs.Draw
	response.RecordsFiltered = total
	response.RecordsTotal = total
	response.Data = t.MakeResponseDatatable(records)
	return
}

func (t *GamNetwork) setFilterSearch(inputs *payload.GamFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		var flag bool
		// Search from form of datatable <- not use
		if inputs.Search != nil && inputs.Search.Value != "" {
			flag = true
		}
		// Search from form filter
		if inputs.PostData.QuerySearch != "" {
			flag = true
		}
		if !flag {
			return db
		}
		return db.Where("network_name LIKE ?", "%"+inputs.PostData.QuerySearch+"%")
	}
}

func (t *GamNetwork) setOrder(inputs *payload.GamFilterPayload) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(inputs.Order) > 0 {
			var orders []string
			for _, order := range inputs.Order {
				column := inputs.Columns[order.Column]
				orders = append(orders, fmt.Sprintf("%s %s", column.Data, order.Dir))
			}
			orderString := strings.Join(orders, ", ")
			return db.Order(orderString)
		}
		return db.Order("id desc")
	}
}

func (t *GamNetwork) MakeResponseDatatable(gams []GamNetworkRecord) (records []GamRecordDatatable) {
	for _, gam := range gams {
		record := new(Gam).GetById(gam.GamId)
		rec := GamRecordDatatable{
			NetworkId:   gam.NetworkId,
			NetworkName: gam.NetworkName,
			ApiAccess:   htmlblock.Render("gam/block/api.gohtml", gam).String(),
			SetUp:       htmlblock.Render("gam/block/setup.gohtml", gam).String(),
			CreateAt:    record.CreatedAt.Format("2006/01/02"),
		}
		records = append(records, rec)
	}
	return
}
