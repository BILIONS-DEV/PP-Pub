package caching

type Cache interface {
	Set(key interface{}, value interface{}, table ...string) error
	SetWithTTL(key interface{}, value interface{}, ttl uint32, table ...string) error
	Get(key interface{}, rec interface{}, table ...string) error
	GetAll(records interface{}, table ...string) (err error)
	Truncate(table ...string) error
	Delete(key interface{}, table ...string) (err error)
	SetBins(setName string, key string, bins ...*Bin) (err error)
	SetBinsWithTTL(setName string, key string, ttl uint32, bins ...*Bin) (err error)
}