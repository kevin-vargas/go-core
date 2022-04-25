package cache

import (
	"sync"
	"time"
)

const NOT_EXPIRED = 0
const DEFAULT_TTL_CACHE = 3600000

type Cache[T any] interface {
	Set(key string, value T)
	SetWithTTL(key string, value T, ttl time.Duration)
	SetWithTTLNano(key string, value T, ttl time.Duration)
	Get(k string) (T, bool)
}

type data[T any] struct {
	value T
	exp   int64
}

func (data *data[T]) Expired() bool {
	exp := data.exp
	if exp == NOT_EXPIRED {
		return false
	}
	return time.Now().UnixNano() > exp
}

type cache[T any] struct {
	sync.RWMutex
	datas      map[string]*data[T]
	ttlDefault time.Duration
}

func (c *cache[T]) Set(key string, value T) {
	c.SetWithTTLNano(key, value, c.ttlDefault)
}

func (c *cache[T]) SetWithTTL(key string, value T, ttl time.Duration) {
	c.SetWithTTLNano(key, value, ttl*time.Millisecond)
}

func (c *cache[T]) SetWithTTLNano(key string, value T, ttl time.Duration) {
	var exp int64 = NOT_EXPIRED
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}
	var data = &data[T]{
		value,
		exp,
	}
	c.Lock()
	c.datas[key] = data
	c.Unlock()
}
func (c *cache[T]) Get(k string) (interface{}, bool) {
	c.RLock()
	data, found := c.datas[k]
	c.RUnlock()
	if !found || data.Expired() {
		return nil, false
	}
	return data.value, true
}

var once sync.Once
var instance Cache[any]

func GetInstance() Cache[any] {
	once.Do(func() {
		instance = &cache[any]{
			datas:      make(map[string]*data[any]),
			ttlDefault: DEFAULT_TTL_CACHE * time.Millisecond,
		}
	})
	return instance
}
