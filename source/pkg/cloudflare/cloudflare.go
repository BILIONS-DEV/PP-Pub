package cloudflare

func ResetLinkAds(uuids []string) (err error) {
	//return
	//if utility.IsWindow() {
	//	return
	//}
	//api, err := apiCloudflare.NewWithAPIToken("Xz3ZVdf00kp2x4VAxtdT6yq6mU-XAmBHsD_JkKJD")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//// Most API calls require a Context
	//ctx := context.Background()
	//var linkAdsTxts []string
	//for _, uuid := range uuids {
	//	linkAdsTxts = append(linkAdsTxts, "https://ms.pubpowerplatform.io/"+uuid+"/ads.txt")
	//}
	////fmt.Println(linkAdsTxts)
	//purgeCacheRequest := apiCloudflare.PurgeCacheRequest{
	//	Everything: false,
	//	Files:      linkAdsTxts,
	//	Tags:       nil,
	//	Hosts:      nil,
	//	Prefixes:   nil,
	//}
	////fmt.Printf("%+v \n", purgeCacheRequest)
	//// Purge Cache
	//_, err = api.PurgeCache(ctx, "27c3588ebe1fc66534e3b90cc7c6c115", purgeCacheRequest)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	////fmt.Printf("%+v \n", res)
	return
}

var accountID = "b72286bd38c4a67ed6345caf8572957f"
var namespaceID = "6ffd875c2768492f8eb9cb63a7481999"
var apiToken = "g2l5whYeUCQdtPpD9V7wxsrRKgskRYjJTVQdAfbM"

type APIResponse struct {
	Success bool `json:"success"`
}

func WriteKeyAndValue(key, value string) error {
	//if utility.IsWindow() {
	//	return nil
	//}
	//
	//api, err := apiCloudflare.NewWithAPIToken(apiToken, apiCloudflare.UsingAccount(accountID))
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	//
	//payload := []byte(value)
	//
	//resp, err := api.WriteWorkersKV(context.Background(), namespaceID, key, payload)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if !resp.Success {
	//	fmt.Println(resp)
	//	return errors.New("cache cloudflare fail")
	//}
	//fmt.Println(resp)
	return nil
}

type KeyAndValue struct {
	Key   string
	Value string
}

func WriteMultipleKeyAndValue(keyAndValues []KeyAndValue) error {
	//if utility.IsWindow() {
	//	return nil
	//}
	//api, err := apiCloudflare.NewWithAPIToken(apiToken, apiCloudflare.UsingAccount(accountID))
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	//ksv := apiCloudflare.WorkersKVBulkWriteRequest{}
	//for _, keyAndValue := range keyAndValues {
	//	kvPair := apiCloudflare.WorkersKVPair{
	//		Key:   keyAndValue.Key,
	//		Value: keyAndValue.Value,
	//	}
	//	ksv = append(ksv, &kvPair)
	//}
	//
	//resp, err := api.WriteWorkersKVBulk(context.Background(), namespaceID, ksv)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if !resp.Success {
	//	fmt.Println(resp)
	//	return errors.New("cache cloudflare fail")
	//}
	//fmt.Println(resp)
	return nil
}
