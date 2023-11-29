package caching

import (
	"encoding/json"
	"errors"
	as "github.com/aerospike/aerospike-client-go"
	"log"
	"math"
)

type aerospikeClient struct {
	client    *as.Client
	namespace string
}

type Config struct {
	Host      string
	Port      int
	Namespace string
}

type Bin struct {
	Name  string
	Value interface{}
}

func NewAerospike(config Config) (client Cache, err error) {
	var aeClient *as.Client
	if aeClient, err = connect(config); err != nil {
		return nil, err
	}
	return &aerospikeClient{
		client:    aeClient,
		namespace: config.Namespace,
	}, nil
}

// Connect : đây là hàm connect Aerospike
func connect(config Config) (client *as.Client, err error) {
	client, err = as.NewClient(config.Host, config.Port)
	return
}

// Set : set cache vào aerospike
// lưu ý: object truyền vào phải là 1 struct có các value public (viết in hoa), không được truyền vào các kiểu dữ liệu int/string/float/...
// ví dụ:
// hợp lệ
//	type user struct{
//		Name string
//		Age int
//	}
// không hợp lệ:
// 	type user struct{
//		name string
//		age int
//	}
func (t *aerospikeClient) Set(key interface{}, value interface{}, table ...string) (err error) {
	if len(table) < 1 {
		return errors.New("missing set name")
	}
	aeKey, err := t.makeKey(table[0], key)
	if err != nil {
		log.Println(err)
		return
	}
	err = t.client.PutObject(t.makeDefaultWritePolicy(), aeKey, value)
	if err != nil {
		log.Println(err)
	}
	return
}

func (t *aerospikeClient) SetWithTTL(key interface{}, value interface{}, ttl uint32, table ...string) (err error) {
	if len(table) < 1 {
		return errors.New("missing set name")
	}
	aeKey, err := t.makeKey(table[0], key)
	if err != nil {
		log.Println(err)
		return
	}
	writePolicy := as.NewWritePolicy(0, ttl)
	writePolicy.SendKey = true
	err = t.client.PutObject(writePolicy, aeKey, value)
	if err != nil {
		log.Println(err)
	}
	return
}

// Get : get cache từ aerospike và truyền vào object
func (t *aerospikeClient) Get(key interface{}, record interface{}, table ...string) (err error) {
	if len(table) < 1 {
		return errors.New("missing set name")
	}
	aeKey, err := t.makeKey(table[0], key)

	if err != nil {
		return
	}
	err = t.client.GetObject(nil, aeKey, record)
	return err
}

// Get : get all cache từ aerospike, records nhận vào là một mảng object(struct)
func (t *aerospikeClient) GetAll(records interface{}, table ...string) (err error) {
	if len(table) < 1 {
		return errors.New("missing set name")
	}

	recs, err := t.client.ScanAll(nil, t.namespace, table[0])
	if err != nil {
		return
	}
	var results []interface{}
	// deal with the error here
	for res := range recs.Results() {
		if res.Err != nil {
			// handle error here
			// if you want to exit, cancel the recordset to release the resources
		} else {
			res.Record.Bins["key"] = res.Record.Key.Value().String()
			results = append(results, res.Record.Bins)
		}
	}

	// parser data to records
	b, err := json.Marshal(results)
	//fmt.Println(string(b))
	err = json.Unmarshal(b, records)
	//fmt.Println(records)
	return err
}

// Get : delete cache
func (t *aerospikeClient) Delete(key interface{}, table ...string) (err error) {
	if len(table) < 1 {
		return errors.New("missing set name")
	}
	aeKey, err := t.makeKey(table[0], key)

	if err != nil {
		return
	}
	_, err = t.client.Delete(nil, aeKey)
	return err
}

// Truncate : xóa toàn bộ dữ liệu cache trong 1 bảng
func (t *aerospikeClient) Truncate(table ...string) (err error) {
	if len(table) < 1 {
		return errors.New("missing set name")
	}
	err = t.client.Truncate(nil, t.namespace, table[0], nil)
	return
}

func (t *aerospikeClient) SetBins(setName string, key string, bins ...*Bin) (err error) {
	writePolicy := t.makeDefaultWritePolicy()
	aeKey, err := t.makeKey(setName, key)

	if err != nil {
		return
	}

	var binAEs []*as.Bin
	for _, bin := range bins {
		binAE := as.NewBin(bin.Name, bin.Value)
		binAEs = append(binAEs, binAE)
	}
	err = t.client.PutBins(writePolicy, aeKey, binAEs...)
	return
}

func (t *aerospikeClient) SetBinsWithTTL(setName string, key string, ttl uint32, bins ...*Bin) (err error) {
	writePolicy := as.NewWritePolicy(0, ttl)
	writePolicy.SendKey = true

	aeKey, err := t.makeKey(setName, key)

	if err != nil {
		return
	}

	var binAEs []*as.Bin
	for _, bin := range bins {
		binAE := as.NewBin(bin.Name, bin.Value)
		binAEs = append(binAEs, binAE)
	}
	err = t.client.PutBins(writePolicy, aeKey, binAEs...)
	return
}

func (t *aerospikeClient) makeKey(table string, key interface{}) (aeKey *as.Key, err error) {
	aeKey, err = as.NewKey(t.namespace, table, key)
	return
}

func (t *aerospikeClient) makeDefaultWritePolicy() (policy *as.WritePolicy) {
	//policy = as.NewWritePolicy(0, math.MaxUint32)
	policy = as.NewWritePolicy(0, 0)
	policy.SendKey = true
	policy.Expiration = math.MaxUint32
	return
}
