package engine

import (
	"time"

	"github.com/kklis/gomemcache"
	"github.com/nishrit/aero/key"
)

type Memcache struct {
	key.NoSpacesFormat
	mc *gomemcache.Memcache
}

func NewMemcache(host string, port int) Memcache {
	// connection check
	serv, err := gomemcache.Connect(host, port)
	if err != nil {
		panic(err)
	}

	return Memcache{
		mc: serv,
	}
}

func (c Memcache) Set(key string, data []byte, expireIn time.Duration) {
	key = c.Format(key)
	c.mc.Set(key, data, 0, int64(expireIn.Seconds()))
}

func (c Memcache) Get(key string) ([]byte, error) {
	var data []byte
	var err error

	key = c.Format(key)
	data, _, err = c.mc.Get(key)

	if err != nil {
		return nil, err
	} else {
		return data, nil
	}
}

func (c Memcache) Delete(key string) error {
	key = c.Format(key)
	return c.mc.Delete(key)
}

func (c Memcache) Close() {
	err := c.mc.Close()
	if err != nil {
		panic(err)
	}
}
