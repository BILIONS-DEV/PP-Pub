package aerospike

import (
	"time"
)

type SetJsCache struct {
	LastTimeUpdate time.Time `as:"last_update" json:"last_update"`
}

type JsCache struct{}

var defaultCacheParam int64 = 91615211004

func (JsCache) Update(inventoryUuid string) error {
	key, err := MakeKey(Sets.JsCache, inventoryUuid)
	if err != nil {
		return err
	}
	obj := SetJsCache{
		LastTimeUpdate: time.Now(),
	}
	err = Client.PutObject(MakeDefaultWritePolicy(), key, &obj)
	return err
}

func (JsCache) GetCacheParam(inventoryUuid string) (cacheParams int64) {
	key, err := MakeKey(Sets.JsCache, inventoryUuid)
	if err != nil {
		return defaultCacheParam
	}
	obj := SetJsCache{}
	err = Client.GetObject(nil, key, &obj)
	if err != nil {
		_ = new(JsCache).Update(inventoryUuid)
		return time.Now().Unix()
	}
	cacheParams = obj.LastTimeUpdate.Unix()
	return
}
