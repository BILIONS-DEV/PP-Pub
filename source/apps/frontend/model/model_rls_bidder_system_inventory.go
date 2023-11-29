package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"source/core/technology/mysql"
)

type RlsBidderSystemInventory struct{}

type RlsBidderSystemRecord struct {
	mysql.TableRlsBidderSystemInventory
}

// GetAll get all BidderSystem for admin
func (t *RlsBidderSystemInventory) GetAll() (rows []BidderRecord) {
	mysql.Client.Where("bidder_type = ?", 2).Order("id DESC").Find(&rows)
	return
}

// GetAll get all BidderSystem for admin
func (t *RlsBidderSystemInventory) GetListIdBidderApprove(inventoryName string) (listBidderId []int64) {
	mysql.Client.Model(&RlsBidderSystemRecord{}).Select("bidder_id").Where("inventory_name = ? and (status = 3 or status = 7 or status = 8)", inventoryName).Find(&listBidderId)
	return
}

// GetAll get all inventory approve bidder
func (t *RlsBidderSystemInventory) GetListInventoryApprove(bidderId int64, listInventoryName []string) (outputs []string) {
	mysql.Client.Model(&RlsBidderSystemRecord{}).Select("inventory_name").Where("bidder_id = ? and (status = 3 or status = 7 or status = 8) and inventory_name in ?", bidderId, listInventoryName).Find(&outputs)
	return
}

// GetRlsOfInventory Get relationship of BidderSystemRecord with InventoryRecord
func (t *RlsBidderSystemInventory) GetRlsOfInventory(bidderSystemIds int64, inventoryName string) (listRls RlsBidderSystemRecord) {
	_ = mysql.Client.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Where("bidder_id = ? AND inventory_name = ?", bidderSystemIds, inventoryName).Last(&listRls).Error
	return
}

func (t *RlsBidderSystemInventory) CheckApproveBidderAdxByUser(bidderAdxId int64, userId int64) (check bool) {
	var checkApproveDomain, checkAcceptMCM bool
	var listNameDomain []string
	mysql.Client.Model(InventoryRecord{}).Select("name").Where("user_id = ? and status = ?", userId, mysql.StatusApproved).Find(&listNameDomain)
	var rls RlsBidderSystemRecord
	_ = mysql.Client.Where("bidder_id = ? AND inventory_name in ? and (status = 3 or status = 7 or status = 8)", bidderAdxId, listNameDomain).Last(&rls).Error
	if rls.Status.IsApproved() || rls.Status == mysql.RlsBidderSystemInventorySubmitted {
		checkApproveDomain = true
	}
	listGam := new(GamNetwork).GetByUser(userId)
	for _, gam := range listGam {
		rlsConnectionMCM := new(RlsConnectionMCM).GetStatus(bidderAdxId, gam.NetworkId, userId)
		if rlsConnectionMCM.Status == mysql.TYPEConnectionMCMTypeAccept {
			checkAcceptMCM = true
		}
	}
	if checkApproveDomain && checkAcceptMCM {
		return true
	}
	return
}
