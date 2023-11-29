package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"log"
	"net/http"
	"source/apps/frontend/payload"
	"source/core/technology/mysql"
	"source/pkg/utility"
	"strconv"
)

type ApiAdsense struct{}

func (t *ApiAdsense) HandlerAdTagPostAPIAdSlot(adTag AdTagRecord) (err error) {
	// Get size từ adTag
	primarySize := new(AdSize).GetById(adTag.PrimaryAdSize)

	// Trong trường hợp sticky banner sẽ check thêm cả size cho mobile
	var mobileSize AdSizeRecord
	if adTag.Type == mysql.TYPEStickyBanner {
		mobileSize = new(AdSize).GetById(adTag.PrimaryAdSizeMobile)
	}

	// Get inventory
	inventory, _ := new(Inventory).GetByIdSystem(adTag.InventoryId)

	// Get các bidder adsense đã approve
	listBidderId := new(RlsBidderSystemInventory).GetListIdBidderApprove(inventory.Name)
	for _, bidderId := range listBidderId {
		var bidder BidderRecord
		bidder.GetById(bidderId)
		if bidder.BidderTemplateId == 1 && bidder.AccountType == mysql.TYPEAccountTypeAdsense {
			// Xử lý theo media type
			for _, mediaType := range bidder.MediaTypes {
				if mediaType.MediaTypeId == 1 { // Display
					// build line adsense name
					lineAdsenseName := new(LineItem).BuildLineNameAdsenseDisplay(inventory.Name, bidder.DisplayName)

					// get line item adsense đang dùng cho inventory + adsense approve này
					lineItemAdsense := new(LineItem).GetByName(lineAdsenseName)
					lineItemAdsense.GetRls()

					// Check size của adTag đã có adSlot cho line item chưa
					var isSize bool
					for _, adSlot := range lineItemAdsense.AdsenseAdSlots {
						if adSlot.Size == primarySize.Name {
							isSize = true
						}
					}

					// Nếu size chưa có trong list adsense slot post lên API của a Tuấn để tạo
					if !isSize {
						err = t.PostAPIAdAdsenseSlot(lineItemAdsense, primarySize.Name)
					}

					// Xử lý thêm trường hợp size mobile cho sticky banner
					if adTag.Type == mysql.TYPEStickyBanner {
						// Check size của adTag đã có adSlot cho line item chưa
						var isSizeMobile bool
						for _, adSlot := range lineItemAdsense.AdsenseAdSlots {
							if adSlot.Size == mobileSize.Name {
								isSizeMobile = true
							}
						}

						// Nếu size chưa có trong list adsense slot post lên API của a Tuấn để tạo
						if !isSizeMobile {
							err = t.PostAPIAdAdsenseSlot(lineItemAdsense, mobileSize.Name)
						}
					}
				}
			}
		}
	}
	return
}

func (t *ApiAdsense) PostAPISubmitAdsense(inventoryName string, accountName string) (err error) {
	if utility.IsWindow() || utility.IsDev() || utility.IsDemo() {
		return
	}
	linkAPI := "http://vli.earns.io/api/api_missing_adsense.php"

	type detail struct {
		Id          string `url:"id"`
		DomainName  string `url:"domain_name"`
		AccountName string `url:"account_name"`
	}
	type request struct {
		Data detail `url:"data"`
	}

	reqs := request{
		Data: detail{
			Id:          "0",
			DomainName:  inventoryName,
			AccountName: accountName,
		},
	}

	params, _ := query.Values(reqs)
	c := &http.Client{}

	req, _ := http.NewRequest("POST", linkAPI, bytes.NewBufferString(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))
	// Submit the request
	res, err := c.Do(req)
	if err != nil {
		return
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	resp := payload.ResponseApiAdsense{}
	_ = json.Unmarshal(bodyBytes, &resp)
	if len(resp.InsertedID) == 0 {
		err = errors.New("submit fail")
	}
	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
	}
	return
}

//func (t *ApiAdsense) HandlerAdsenseAdslot(inventoryName string, bidder BidderRecord) (err error) {
//	// từ inventory lấy ra listInventoryId và listSize từ các tag
//	inventories, _ := new(Inventory).GetByName(inventoryName)
//	var listSize []string
//	var listInventoryId []int64
//	for _, inventory := range inventories {
//		listInventoryId = append(listInventoryId, inventory.Id)
//
//		adTags := new(InventoryAdTag).GetByInventory(inventory.Id)
//		for _, adTag := range adTags {
//			// Bỏ qua các type không phải banner
//			if adTag.Type != mysql.TYPEDisplay && adTag.Type != mysql.TYPEStickyBanner {
//				continue
//			}
//			primarySize := new(AdSize).GetById(adTag.PrimaryAdSize)
//			if !utility.InArray(primarySize.Name, listSize, false) {
//				listSize = append(listSize, primarySize.Name)
//			}
//		}
//	}
//
//	// Từ list size đẩy lên API a Tuấn để tạo creative
//	for _, size := range listSize {
//		err = t.PostAPIAdAdsenseSlot(inventoryName, bidder, size)
//	}
//
//	return
//}

func (t *ApiAdsense) PostAPIAdAdsenseSlot(lineItemAdsense LineItemRecord, size string) (err error) {
	//if utility.IsWindow() || utility.IsDev() || utility.IsDemo() {
	//	return
	//}
	//
	//linkAPI := "http://vli.earns.io/api/api_missing_line_item.php"
	//type detail struct {
	//	DomainName  string `url:"domain_name"`
	//	AccountName string `url:"account"`
	//	Size        string `url:"size"`
	//}
	//type request struct {
	//	Data detail `url:"data[0]"`
	//}
	//
	//reqs := request{}
	//reqs.Data = detail{
	//	DomainName:  inventoryName,
	//	AccountName: bidder.DisplayName,
	//	Size:        size,
	//}
	//
	//params, _ := query.Values(reqs)
	//c := &http.Client{}
	//
	//fmt.Printf("PostAPIAdAdsenseSlot: %+v \n", params)
	//req, _ := http.NewRequest("POST", linkAPI, bytes.NewBufferString(params.Encode()))
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))
	//// Submit the request
	//res, err := c.Do(req)
	//if err != nil {
	//	return
	//}
	//bodyBytes, err := io.ReadAll(res.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Printf("%+v \n", string(bodyBytes))
	//resp := input.ResponseApiAdsense{}
	//_ = json.Unmarshal(bodyBytes, &resp)
	////if len(resp.InsertedID) == 0 {
	////	err = errors.New("submit fail")
	////}
	//// Check the response
	//if res.StatusCode != http.StatusOK {
	//	err = fmt.Errorf("bad status: %s", res.Status)
	//}

	_ = new(LineItemAdsenseAdSlot).Push(lineItemAdsense.Id, size, "")
	return
}
