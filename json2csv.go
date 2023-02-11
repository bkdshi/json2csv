package json2csv

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func json2csv(data interface{}) []map[string]interface{} {
	var results []map[string]interface{}

	v := getValue(data)

	switch v.Kind() {
	case reflect.Map:
		if v.Len() > 0 {
			results = append(results, flatten(v))
		}
	case reflect.Slice:

		for i := 0; i < v.Len(); i++ {
			results = append(results, flatten(v))
		}

		// if v.Len() > 0 {
		// 	results = append(results, flatten(v))
		// }
	}
	return results
}

func flatten(obj interface{}) map[string]interface{} {
	f := make(map[string]interface{}, 0)
	keys := make([]string, 0)

	_flatten(f, keys, obj)

	return f
}

func _flatten(f map[string]interface{}, keys []string, obj interface{}) {
	// get reflect.Value from interface
	value := getValue(obj)

	switch value.Kind() {
	case reflect.Map:
		_flattenMap(f, keys, value)
	case reflect.Slice:
		_flattenSlice(f, keys, value)
	case reflect.Float64:
		key := strings.Join(keys, "/")
		f[key] = value.Float()
	case reflect.String:
		key := strings.Join(keys, "/")
		f[key] = value.String()
	default:
		fmt.Println(keys, value, value.Kind())
		log.Fatal("SWITCH DEFAULT UNKOWN KIND")
	}
}

func _flattenMap(f map[string]interface{}, keys []string, value reflect.Value) {
	// keys := sortedMapKeys(value)
	mapKeys := value.MapKeys()
	// sort.Sort(mapKeys)
	for _, key := range mapKeys {
		cloneKey := cloneKey(keys)
		cloneKey = append(cloneKey, key.String())
		_flatten(f, cloneKey, value.MapIndex(key))
	}
}

func _flattenSlice(f map[string]interface{}, keys []string, value reflect.Value) {
	if value.Len() == 0 {
		_flatten(f, keys, "")
	}
	for i := 0; i < value.Len(); i++ {
		cloneKey := cloneKey(keys)
		cloneKey = append(cloneKey, strconv.Itoa(i))
		_flatten(f, cloneKey, value.Index(i))
	}
}
