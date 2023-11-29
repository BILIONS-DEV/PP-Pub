package bidder

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"source/core/technology/mysql"
	"source/internal/entity/model"
)

type RepoBidder interface {
	FindByID(id int64) (record *model.BidderModel, err error)
	FindAllByQuery(query model.BidderModel, queryRaw ...string) (record []*model.BidderModel, err error)
	FindByUser(userID int64) (records []model.BidderModel, err error)
	FindByName(name string, userID ...int64) (record *model.BidderModel, err error)
	ListDomainBidderStatusByBidder(bidderID int64, domains []string) (records []model.RlsBidderSystemInventoryModel, err error)
}

func NewBidderRepo(DB *gorm.DB) *bidderRepo {
	return &bidderRepo{Db: DB}
}

type bidderRepo struct {
	Db *gorm.DB
}

func (t *bidderRepo) FindByID(id int64) (record *model.BidderModel, err error) {
	err = t.Db.Find(&record, id).Error
	return
}

func (t *bidderRepo) FindByUser(userID int64) (records []model.BidderModel, err error) {
	err = t.Db.Where("user_id = ?", userID).Find(&records).Error
	return
}

func (t *bidderRepo) FindByName(name string, userID ...int64) (record *model.BidderModel, err error) {
	err = t.Db.Where("(bidder_code = ? or display_name = ?) and user_id in ?", name, name, userID).Find(&record).Error
	return
}

func (t *bidderRepo) ListDomainBidderStatusByBidder(bidderID int64, domains []string) (records []model.RlsBidderSystemInventoryModel, err error) {
	query := t.Db.Where("bidder_id = ?", bidderID).
		Joins("JOIN inventory ON inventory.name = rls_bidder_system_inventory.inventory_name AND inventory.status = ?", mysql.StatusApproved)
	if len(domains) > 1 || (len(domains) == 1 && domains[0] != "") {
		query.Where("inventory_name in ?", domains)
	}
	err = query.Find(&records).Error
	return
}

func (t *bidderRepo) FindAllByQuery(query model.BidderModel, queryRaw ...string) (record []*model.BidderModel, err error) {
	err = t.Db.
		//Debug().
		// Unscoped(). //=> để gọi ra các field deleted_at != NULL
		Preload(clause.Associations).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if len(queryRaw) > 0 {
				for _, raw := range queryRaw {
					db.Where(raw)
				}
			}
			return db
		}).
		Where(query).
		Find(&record).Error
	return
}
