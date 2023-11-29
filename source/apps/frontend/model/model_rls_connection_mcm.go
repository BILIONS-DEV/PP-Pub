package model

import 	"source/core/technology/mysql"

type RlsConnectionMCM struct{}

type RlsConnectionMCMRecord struct {
	mysql.TableRlsConnectionMCM
}

func (RlsConnectionMCMRecord) TableName() string {
	return mysql.Tables.RlsConnectionMCM
}

func (RlsConnectionMCM) Create(rl RlsConnectionMCMRecord) (err error) {
	err = mysql.Client.Create(&rl).Error
	return err
}

func (RlsConnectionMCM) Delete(bidderSystemId int64) (err error) {
	err = mysql.Client.Where(RlsConnectionMCMRecord{mysql.TableRlsConnectionMCM{BidderId: bidderSystemId}}).Delete(RlsConnectionMCMRecord{}).Error
	return err
}

func (RlsConnectionMCM) GetByBidderId(bidderSystemId int64) (record []RlsConnectionMCMRecord) {
	mysql.Client.Where(RlsConnectionMCMRecord{mysql.TableRlsConnectionMCM{BidderId: bidderSystemId}}).Find(&record)
	return
}

func (RlsConnectionMCM) GetStatus(bidderSystemId int64, gamNetworkId int64, userId int64) (record RlsConnectionMCMRecord) {
	gamNetwork := new(GamNetwork).GetByNetworkId(gamNetworkId, userId)
	mysql.Client.Where(RlsConnectionMCMRecord{mysql.TableRlsConnectionMCM{BidderId: bidderSystemId, NetworkId: gamNetwork.NetworkId}}).Find(&record)
	return
}
