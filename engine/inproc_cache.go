package engine

import (
	"errors"
	"time"

	"github.com/nishrit/aero/key"
	goc "github.com/pmylund/go-cache"
)

// Stores data in memory. Major use-case is for use on development machines
type InProcCache struct {
	key.AsIsFormat
	ram *goc.Cache
}

func NewInProcCache(expires time.Duration) InProcCache {
	return InProcCache{
		ram: goc.New(expires, 5*time.Minute),
	}
}

func (c InProcCache) Set(key string, data []byte, expireIn time.Duration) {
	key = c.Format(key)
	c.ram.Set(key, data, expireIn)
}

func (c InProcCache) Get(key string) ([]byte, error) {
	key = c.Format(key)
	if i, ok := c.ram.Get(key); ok {
		return i.([]byte), nil
	} else {
		return nil, errors.New("key not found in inmem cache: " + key)
	}
}

func (c InProcCache) Delete(key string) error {
	key = c.Format(key)
	c.ram.Delete(key)
	return nil
}

func (c InProcCache) Close() {

}
