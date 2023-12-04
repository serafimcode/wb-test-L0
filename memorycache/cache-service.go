package memorycache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/serafimcode/wb-test-L0/model"
	"sync"
	"time"
)

type Cache struct {
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]Item
	sync.RWMutex
}

func (c *Cache) Log() {
	fmt.Println("CACHE DATA:")

	for _, item := range c.items {
		data := item.Value.(json.RawMessage)
		var dataInfo model.OrderInfoDTO

		json.Unmarshal(data, &dataInfo)

		fmt.Printf("%v: \n", dataInfo)
	}
}

type Item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

func New(expiration, cleanupInterval time.Duration) *Cache {
	cache := Cache{
		defaultExpiration: expiration,
		cleanupInterval:   cleanupInterval,
		items:             map[string]Item{},
	}

	if cleanupInterval > 0 {
		go cache.startCleanUp()
	}

	return &cache
}

func (cache *Cache) startCleanUp() {
	ticker := time.NewTicker(cache.cleanupInterval)

	for {
		<-ticker.C

		if cache.items == nil {
			return
		}

		for key, item := range cache.items {
			if time.Now().UnixNano() > item.Expiration && item.Expiration > 0 {
				delete(cache.items, key)
			}
		}
	}
}

func (cache *Cache) Set(key string, value interface{}, ttl time.Duration) {
	var expiration int64

	if ttl == 0 {
		expiration = int64(cache.defaultExpiration)
	}

	if ttl > 0 {
		expiration = time.Now().Add(ttl).UnixNano()
	}

	cache.Lock()
	defer cache.Unlock()

	cache.items[key] = Item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

func (cache *Cache) Get(key string) (interface{}, bool) {
	cache.RLock()
	defer cache.RUnlock()

	item, found := cache.items[key]

	if !found {
		return nil, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}

	return item.Value, true
}

func (cache *Cache) Delete(key string) error {
	cache.RLock()
	defer cache.RUnlock()

	if _, found := cache.items[key]; !found {
		return errors.New("key not found")
	}

	delete(cache.items, key)

	return nil
}
