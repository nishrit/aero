package orm

import (
	"github.com/jinzhu/gorm"
	"github.com/nishrit/aero/cstr"
)

var engines map[string]*gorm.DB

func init() {
	engines = make(map[string]*gorm.DB)

	// default to table singular naming convention
	Initialize(func(o *gorm.DB) {
		o.SingularTable(true)
	})
}

func Get(useMaster bool) *gorm.DB {
	s := cstr.Get(useMaster)
	return GetConn(s.Engine, s.Conn)
}

func GetConfig(container string) *gorm.DB {
	s := cstr.ReadConfig(container)
	return GetConn(s.Engine, s.Conn)
}

func GetConn(engine string, conn string) *gorm.DB {
	var ormCurr *gorm.DB
	var ormObj *gorm.DB
	var ok bool
	var err error

	if ormCurr, ok = engines[conn]; ok {
		return ormCurr.Unscoped()
	}

	// http://go-database-sql.org/accessing.html
	// the sql.DB object is designed to be long-lived
	if ormObj, err = gorm.Open(engine, conn); err == nil {
		if ormInit != nil {
			for _, fn := range ormInit {
				fn(ormObj)
			}
		}
		engines[conn] = ormObj
		return ormObj.Unscoped()
	}
	panic(err)
}

var ormInit []func(*gorm.DB) = make([]func(*gorm.DB), 0)

func Initialize(fn func(*gorm.DB)) {
	// TODO: use mutex
	ormInit = append(ormInit, fn)
}
