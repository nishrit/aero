package my

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/nishrit/aero/str"
)

type table struct {
	name  string
	model interface{}
}

func NewTable(name string) *table {
	return &table{
		name: name,
	}
}

func NewTable2(model interface{}) *table {
	typ := reflect.TypeOf(model)
	if typ.Elem() == nil {
		panic("model must be an address")
	}

	return &table{
		name:  str.SnakeCase(typ.Elem().Name()),
		model: model,
	}
}

func (t *table) isHistory() bool {
	return strings.HasSuffix(t.name, "_history")
}

func (t *table) history() *table {
	if strings.HasSuffix(t.name, "_history") {
		panic("already a history table")
	}
	return &table{
		name: t.name + "_history",
	}
}

func (t *table) exists() bool {
	sql := fmt.Sprintf("show tables like '%s'", t.name)
	return sqlHasRows(sql)
}

func (t *table) hasData() bool {
	sql := fmt.Sprintf("select * from %s limit 1", t.name)
	return sqlHasRows(sql)
}

func (t *table) drop(force bool) {

	if !t.exists() {
		return
	}

	if t.hasData() && force == false {
		return
	}

	sql := fmt.Sprintf("drop table %s", t.name)
	sqlExec(sql)
}

func (t *table) fields() []field {
	var fields []field
	sql := fmt.Sprintf("desc %s", t.name)
	err := Dbo.Raw(sql).Find(&fields).Error
	if err != nil {
		panic(err)
	}
	return fields
}

func (t *table) field(name string) *field {
	flds := t.fields()
	for _, f := range flds {
		if f.Field == name {
			return &f
		}
	}
	return nil
}

func (t *table) autoIncrField() *field {
	flds := t.fields()
	for _, f := range flds {
		if strings.Contains(f.Extra, "auto_increment") {
			return &f
		}
	}
	return nil
}

func (t *table) primaryKeys() []field {
	pkeys := []field{}

	flds := t.fields()
	for _, f := range flds {
		if strings.Contains(f.Key, "PRI") {
			pkeys = append(pkeys, f)
		}
	}
	return pkeys
}
