package orm

import (
	"fmt"
	"reflect"

	"github.com/nishrit/aero/ds"
	"github.com/nishrit/aero/refl"
	"github.com/nishrit/aero/str"
)

func Insertable(modl interface{}, data map[string]string) (bool, []error) {
	return validate(modl, data, "insert")
}

func Updatable(modl interface{}, data map[string]string) (bool, []error) {
	return validate(modl, data, "update")
}

func validate(modl interface{}, data map[string]string, lookupTag string) (bool, []error) {
	var errs []error = make([]error, 0)

	input := clone(data)

	obj := modl
	if reflect.TypeOf(modl).Kind() == reflect.Ptr {
		obj = reflect.ValueOf(modl).Elem()
	}

	for _, fld := range refl.NestedFields(obj) {
		name := fld.Name
		sql := str.SnakeCase(name)
		_, ok := input[sql]

		// are "must" fields missing?
		if ok == false && fld.Tag.Get(lookupTag) == "must" {
			errs = append(errs, fmt.Errorf("Compulsory field missing: %s", sql))
		}

		// are "no" fields present?
		if ok == true && fld.Tag.Get(lookupTag) == "no" {
			errs = append(errs, fmt.Errorf("Extra field found: %s", sql))
		}

		// json_array, json_map and json validations
		sgnt := refl.TypeSignature(fld.Type)
		if ok {
			//fmt.Println(sql, sgnt)
			if sgnt == "sl:." || sgnt == "*sl:." {
				var test []interface{}
				if err := ds.Load(&test, []byte(data[sql])); err != nil {
					errs = append(errs, fmt.Errorf("Field must be json array: %s", sql))
				}
			} else if sgnt == "map" || sgnt == "*map" {
				var test map[string]interface{}
				if ds.Load(&test, []byte(data[sql])) != nil {
					errs = append(errs, fmt.Errorf("Field must be json document: %s", sql))
				}
			} else if sgnt == "*sl:.uint8" || sgnt == "sl:.uint8" {
				var test interface{}
				if ds.Load(&test, []byte(data[sql])) != nil {
					errs = append(errs, fmt.Errorf("Field must be valid json: %s", sql))
				}
			} else if sgnt == "*sl:.string" || sgnt == "sl:.string" {
				var test []string
				if err := ds.Load(&test, []byte(data[sql])); err != nil {
					errs = append(errs, fmt.Errorf("Field must be json string array: %s", sql))
				}
			}
		}
	}

	if len(errs) == 0 {
		return true, nil
	}
	return false, errs
}

func clone(data map[string]string) map[string]interface{} {
	out := make(map[string]interface{})
	for key, val := range data {
		out[key] = val
	}
	return out
}
