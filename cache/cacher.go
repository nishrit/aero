package cache

import (
	"fmt"
	"strings"
	"time"

	"github.com/nishrit/aero/conf"
	"github.com/nishrit/aero/db/cstr"
	"github.com/nishrit/aero/engine"
	"github.com/nishrit/aero/key"
)

type Cacher interface {
	key.KeyFormatter
	Set(key string, data []byte, expireIn time.Duration)
	Get(key string) ([]byte, error)
	Delete(key string) error
	Close()
}

func NewCacher(container ...string) (out Cacher) {
	parent := strings.Join(container, ".")

	engn := conf.String("", parent, "engine")
	if engn == "" {
		panic("cache engine is not specified")
	}

	switch engn {
	case "memcache":
		{
			cnf := cstr.Memcache{}
			conf.Struct(&cnf, parent)
			out = engine.NewMemcache(cnf.Host, cnf.Port)
		}

	case "inproc", "inproccache":
		{
			expiry := conf.String("1h", parent, "lifetime")
			life, err := time.ParseDuration(expiry)
			if err != nil {
				panic(err)
			}
			out = engine.NewInProcCache(life)
		}

	case "log":
		{
			dir := conf.String("", parent, "dir")
			inner := NewCacher(parent, "inner")
			out = NewCacheLogger(dir, inner)
		}

	case "redis":
		{
			cnf := cstr.Redis{}
			conf.Struct(&cnf, parent)
			out = engine.NewRedis(cnf.Host, cnf.Port, cnf.Db)
		}

	default:
		panic(fmt.Errorf("Unknown cache engine: %s", engn))
	}

	return out
}
