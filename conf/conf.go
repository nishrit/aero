package conf

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/jacobstr/confer"
	"github.com/serenize/snaker"
)

// TODO:
// - allow command line arguments to be passed
// - eg. --APP_PORT=1234, should be bubbled to the right place

var configuration *confer.Config

func init() {
	loadDefaultConfig()
}

func loadDefaultConfig() {
	configuration = confer.NewConfig()
	seek := []string{
		// least priority
		"./config/config.yaml",
		"config.yaml",

		// next priority
		"./config/dev.yaml",
		"dev.yaml",

		// highest priority
		"./config/production.yaml",
		"production.yaml",
	}
	var err error
	var files []string = make([]string, 0)
	for _, f := range seek {
		tmp := confer.NewConfig()
		err = tmp.ReadPaths(f)
		if err == nil {
			abs, _ := filepath.Abs(f)
			files = append(files, abs)
		}
	}

	if len(files) == 0 {
		fmt.Println("No yaml configuration file found.")
	} else {
		configuration.ReadPaths(files...)
		fmt.Println("Loading configurations:", len(files), "file(s)")
		for i := 0; i < len(files)-1; i++ {
			fmt.Print(files[i], " → ")
		}
		fmt.Print(files[len(files)-1], "\n")
	}
}

func Get(defaultVal interface{}, keys ...string) interface{} {
	key := strings.Join(keys, ".")
	if configuration.IsSet(key) {
		return configuration.Get(key)
	} else {
		return defaultVal
	}
}

func Int(defaultVal int, keys ...string) int {
	key := strings.Join(keys, ".")
	if Exists(key) {
		return configuration.GetInt(key)
	} else {
		return defaultVal
	}
}

func IntMust(keys ...string) int {
	key := strings.Join(keys, ".")
	if !Exists(key) {
		panic("Int key missing: " + key)
	}

	return configuration.GetInt(key)
}

func String(defaultVal string, keys ...string) string {
	key := strings.Join(keys, ".")
	if Exists(key) {
		return configuration.GetString(key)
	} else {
		return defaultVal
	}
}

func StringMust(keys ...string) string {
	key := strings.Join(keys, ".")
	if !Exists(key) {
		panic("String key missing: " + key)
	}
	return configuration.GetString(key)
}

func StringSlice(defaultVal []string, keys ...string) []string {
	key := strings.Join(keys, ".")
	if Exists(key) {
		return configuration.GetStringSlice(key)
	} else {
		return defaultVal
	}
}

func Bool(defaultVal bool, keys ...string) bool {
	key := strings.Join(keys, ".")
	if Exists(key) {
		return configuration.GetBool(key)
	} else {
		return defaultVal
	}
}

func BoolMust(keys ...string) bool {
	key := strings.Join(keys, ".")
	if !Exists(key) {
		panic("Bool key missing: " + key)
	}

	return configuration.GetBool(key)
}

func Exists(keys ...string) bool {
	key := strings.Join(keys, ".")
	return configuration.IsSet(key)
}

func Struct(addr interface{}, keys ...string) {
	container := strings.Join(keys, ".")

	// addr should be an address
	s := fmt.Sprintf("%s", reflect.TypeOf(addr))
	if !strings.HasPrefix(s, "*") {
		panic("conf.Read() expects address of struct")
	}

	rt := reflect.TypeOf(addr).Elem()
	rv := reflect.ValueOf(addr).Elem()
	for i := 0; i < rt.NumField(); i++ {
		ft := rt.Field(i)
		key := ft.Name

		// read key in any of thiese forms:
		// numItems, numitems, num_items, num-items
		keys := []string{key, strings.ToLower(key), snaker.CamelToSnake(key), strings.Replace(snaker.CamelToSnake(key), "_", "-", -1)}
		found := false
		for _, k := range keys {
			if Exists(container, k) {
				found = true
				switch fmt.Sprintf("%s", ft.Type) {
				case "string":
					rv.Field(i).SetString(String("", container, k))
				case "int":
					rv.Field(i).SetInt(int64(Int(0, container, k)))
				default:
					panic(fmt.Sprintf("conf.Read() found '%s' (must be string|int)", ft.Type))
				}
			}
		}

		if !found { // if it was required then error
			if !strings.Contains(ft.Tag.Get("conf"), "optional") {
				panic(fmt.Sprintf("Config '%s' missing in reading %s", key, rt))
			}
		}
	}
}
