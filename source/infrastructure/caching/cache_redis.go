package caching

type redisClient struct{}

func NewRedis() Cache {
	return &redisClient{}
}

func (t *redisClient) Set(key interface{}, value interface{}, table ...string) (err error) { return }
func (t *redisClient) SetWithTTL(key interface{}, value interface{}, ttl uint32, table ...string) (err error) {
	return
}
func (t *redisClient) Get(key interface{}, record interface{}, table ...string) (err error) { return }

func (t *redisClient) GetAll(records interface{}, table ...string) (err error) {
	return err
}

func (t *redisClient) Delete(key interface{}, table ...string) (err error) { return }
func (t *redisClient) Truncate(table ...string) (err error) { return }
func (t *redisClient) SetBins(setName string, key string, bins ...*Bin) (err error) { return }
func (t *redisClient) SetBinsWithTTL(setName string, key string, ttl uint32, bins ...*Bin) (err error) { return }
