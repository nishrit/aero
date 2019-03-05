package db

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/nishrit/aero/ds"
)

// Add support for Json fields
// http://www.booneputney.com/2015-06-18-gorm-golang-jsonb-value-copy/

type JDoc map[string]interface{}

func NewJDoc(data map[string]interface{}) *JDoc {
	j := make(JDoc)
	for key := range data {
		j[key] = data[key]
	}
	return &j
}

func NewJDoc2() *JDoc {
	j := make(JDoc)
	return &j
}

func (j *JDoc) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	str, err := json.Marshal(j)
	return string(str), err
}

func (j *JDoc) Scan(value interface{}) error {
	if value == nil {
		j = nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source was not []bytes")
	}
	if err := json.Unmarshal(bytes, j); err != nil {
		return err
	}
	return nil
}

func (j *JDoc) Set(key string, val interface{}) *JDoc {
	(*j)[key] = val
	return j
}

type JArr []interface{}

func NewJArr(items ...interface{}) *JArr {
	//return new(JsonA)
	len := len(items)
	arr := make(JArr, len)
	for i := 0; i < len; i++ {
		arr[i] = items[i]
	}
	return &arr
}

func (j *JArr) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	str, err := json.Marshal(j)
	return string(str), err
}

func (j *JArr) Scan(value interface{}) error {
	if value == nil {
		j = nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source was not []bytes")
	}
	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}
	return nil
}

type JSArr []string

func NewJSArr(items ...string) *JSArr {
	len := len(items)
	arr := make(JSArr, len)
	for i := 0; i < len; i++ {
		arr[i] = items[i]
	}
	return &arr
}

func (j *JSArr) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	str, err := json.Marshal(j)
	return string(str), err
}

func (j *JSArr) Scan(value interface{}) error {
	if value == nil {
		j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source was not []bytes")
	}
	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}
	return nil
}

type JRaw []byte

func NewJRaw(ifc interface{}) *JRaw {
	bytes, err := ds.ToBytes(ifc, false)
	if err != nil {
		panic(err)
	}
	var j JRaw = bytes
	return &j
}

func NewJRaw2(str string) *JRaw {
	var j JRaw = []byte(str)
	return &j
}

func (j *JRaw) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return string(*j), nil
}

func (j *JRaw) Scan(value interface{}) error {
	if value == nil {
		*j = nil // todo check?
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source was not []byte")
	}

	*j = bytes
	return nil
}

func (j *JRaw) MarshalJSON() ([]byte, error) {
	return *j, nil
}

func (j *JRaw) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("JRaw: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

func (j *JRaw) Obtain() interface{} {
	if j == nil {
		return nil
	}
	var ifc interface{}
	err := ds.Load(&ifc, *j)
	if err != nil {
		panic(err)
	}
	return ifc
}
