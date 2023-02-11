package json2csv

import (
	"log"
	"reflect"
	"strconv"
	"strings"
)

// Return true if v is object
func isObject(v interface{}) bool {
	value := getValue(v)
	for i := 0; i < value.Len(); i++ {
		if getValue(value.Index(i)).Kind() != reflect.Map {
			return false
		}
	}
	return true
}

// Get reflect.Value from interface
func getValue(obj interface{}) reflect.Value {
	// interface to reflect.Value
	value, ok := obj.(reflect.Value)

	// if invalid reflect value, make new value initialized
	if !ok {
		value = reflect.ValueOf(obj)
	}

	// if it interface, get value
	if value.Kind() == reflect.Interface {
		value = value.Elem()
	}
	return value
}

// Convert JSON to CSV
func Json2Csv(data interface{}) []map[string]interface{} {
	var results []map[string]interface{}

	v := getValue(data)

	switch v.Kind() {
	case reflect.Map:
		if v.Len() > 0 {
			results = append(results, flatten(v))
		}
	case reflect.Slice:
		if isObject(v) {
			for i := 0; i < v.Len(); i++ {
				results = append(results, flatten(v.Index(i)))
			}
		} else if v.Len() > 0 {
			results = append(results, flatten(v))
		}
	default:
		log.Fatal("Unsupported")
	}
	return results
}

// Flatten JSON structure into one-dimensional structure
func flatten(obj interface{}) map[string]interface{} {
	f := make(map[string]interface{}, 0)
	keys := make([]string, 0)

	_flatten(f, keys, obj)
	return f
}

// Flatten structure recursively
func _flatten(f map[string]interface{}, keys []string, obj interface{}) {
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
		log.Println(keys, value, value.Kind())
		log.Fatal("SWITCH DEFAULT UNKOWN KIND")
	}
}

// Flatten map
func _flattenMap(f map[string]interface{}, keys []string, value reflect.Value) {
	for _, key := range value.MapKeys() {
		cloneKey := cloneKey(keys)
		cloneKey = append(cloneKey, key.String())
		_flatten(f, cloneKey, value.MapIndex(key))
	}
}

// Flatten Slice
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

// Clone Key
func cloneKey(keys []string) []string {
	if len(keys) == 0 {
		return keys
	}

	cloneKey := make([]string, len(keys))
	copy(cloneKey, keys)
	return cloneKey
}
