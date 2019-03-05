package que

import (
	"fmt"
	"strings"
	"time"

	"github.com/nishrit/aero/conf"
	"github.com/nishrit/aero/db/cstr"
	"github.com/nishrit/aero/engine"
)

type Queue interface {
	Push(data []byte) error
	Pop() ([]byte, error)
	PopWait(time.Duration) ([]byte, error)
	Len() (int, error)
	Close()
}

func NewQueue(container ...string) (out Queue) {
	parent := strings.Join(container, ".")

	engn := conf.String("", parent, "engine")
	if engn == "" {
		panic(fmt.Sprintf("queue engine is not specified under %s", parent))
	}

	switch engn {
	case "redis":
		{
			cnf := cstr.Redis{}
			conf.Struct(&cnf, parent)
			if cnf.Name == "" {
				panic(fmt.Sprintf("Redis queue name missing"))
			}
			out = engine.NewRedis2(cnf.Host, cnf.Port, cnf.Db, cnf.Name)
		}

	default:
		panic(fmt.Sprintf("Unknown queue engine: %s", engn))
	}

	return out
}
