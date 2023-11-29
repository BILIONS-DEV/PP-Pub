package ipRange

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"source/pkg/utility"
)

func SaveIpRange(path string) {
	var listIpRange []string
	listIpGoogle, err := new(google).GetListIpGoogle()
	if err != nil {
		fmt.Println(err)
	}
	//for _, ipRange := range listIpGoogle {
	//	new(model.LogIpRange).Create(model.LogIpRangeRecord{TableLogIpRange: mysql.TableLogIpRange{
	//		Type:    "google",
	//		IpRange: ipRange,
	//	}})
	//}

	listIpAmazon, err := new(amazon).GetListIpAmazon()
	if err != nil {
		fmt.Println(err)
	}
	for _, ipRange := range listIpAmazon {
		if !utility.InArray(ipRange, listIpRange, false) {
			listIpRange = append(listIpRange, ipRange)
		}
	}
	listIpRange = append(listIpRange, listIpGoogle...)

	listIpClourdFlare, err := new(cloudFlare).GetListIpCloudFlare()
	if err != nil {
		fmt.Println(err)
	}
	for _, ipRange := range listIpClourdFlare {
		if !utility.InArray(ipRange, listIpRange, false) {
			listIpRange = append(listIpRange, ipRange)
		}
	}

	listIpMicrosoft, err := new(microsoft).GetListIpMicrosoft()
	if err != nil {
		fmt.Println(err)
	}
	for _, ipRange := range listIpMicrosoft {
		if !utility.InArray(ipRange, listIpRange, false) {
			listIpRange = append(listIpRange, ipRange)
		}
	}

	data, _ := json.Marshal(listIpRange)
	f, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.Write(data)

	if err2 != nil {
		log.Fatal(err2)
	}

}
