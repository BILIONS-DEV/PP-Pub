package model

import (
	"fmt"
	"source/core/technology/mysql"
	"source/pkg/scanAds"
)

type MissingAdsTxt struct{}

func (t *MissingAdsTxt) GetByInventory(inventory InventoryRecord) (records []mysql.TableMissingAdsTxt) {
	mysql.Client.Where("inventory_id = ? AND user_id = ?", inventory.Id, inventory.UserId).Find(&records)
	return
}

func (t *MissingAdsTxt) ScanSuccessSaved(inventory InventoryRecord, scanResponse scanAds.Response) (err error, syncStatus mysql.TYPEInventorySyncAdsTxt) {
	// Remove all row in mysql
	err = t.RemoveAllByInventoryAndUser(inventory.UserId, inventory.Id)
	if err != nil {
		return
	}
	// Create row missing
	var records []mysql.TableMissingAdsTxt
	for _, line := range scanResponse.Lines {
		if !line.Match {
			records = append(records, mysql.TableMissingAdsTxt{
				Line:        line.Text,
				UserId:      inventory.UserId,
				InventoryId: inventory.Id,
				Domain:      inventory.Domain,
				AdsTxtUrl:   scanResponse.AdsTxtUrl,
			})
		}
	}
	syncStatus = mysql.InventorySyncAdsTxt
	if len(records) > 0 {
		syncStatus = mysql.InventorySyncAdsTxtNotIn
		err = mysql.Client.Create(&records).Error
	}
	return
}

func (t *MissingAdsTxt) ScanErrorSaved(inventory InventoryRecord, scanResponse scanAds.Response) (err error) {
	// Remove all row in mysql
	err = t.RemoveAllByInventoryAndUser(inventory.UserId, inventory.Id)
	if err != nil {
		return err
	}

	record := mysql.TableMissingAdsTxt{
		UserId:       inventory.UserId,
		InventoryId:  inventory.Id,
		Domain:       inventory.Domain,
		AdsTxtUrl:    scanResponse.AdsTxtUrl,
		ErrorMessage: fmt.Sprintf("%s >>> {statusCode}%d{/statusCode}", scanResponse.Error, scanResponse.HeaderStatusCode),
	}
	err = mysql.Client.Create(&record).Error
	return err
}

func (t *MissingAdsTxt) RemoveAllByInventoryAndUser(userId, inventoryId int64) (err error) {
	return mysql.Client.Where("user_id = ? AND inventory_id = ?", userId, inventoryId).Delete(&mysql.TableMissingAdsTxt{}).Error
}
