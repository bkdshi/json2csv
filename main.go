package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

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

func cloneKey(keys []string) []string {
	if len(keys) == 0 {
		return keys
	}

	cloneKey := make([]string, len(keys))
	copy(cloneKey, keys)
	return cloneKey
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

func main() {
	log.SetFlags(log.Llongfile)
	path := "sample.json"
	j, _ := os.ReadFile(path)
	var x interface{}

	err := json.Unmarshal(j, &x)
	_ = err

	results := json2csv(x)
	writeCSV(results)
}

func writeCSV(results []map[string]interface{}) {
	f, err := os.Create("result.csv")
	if err != nil {
		fmt.Println(err)
	}

	w := csv.NewWriter(f)

	headers := getHeader(results)
	sort.SliceStable(headers, func(i, j int) bool { return headers[i] < headers[j] })
	fmt.Println(headers)

	if err := w.Write(headers); err != nil {
		fmt.Println(err)
	}

	for i, result := range results {
		record := make([]string, 0, len(result))
		for _, key := range headers {
			if len(results) > 1 {
				key = strconv.Itoa(i) + "/" + key
			}
			// fmt.Println(key)
			value := result[key]
			if value == nil {
				value = ""
			}
			record = append(record, fmt.Sprintf("%v", value))

		}
		if err := w.Write(record); err != nil {
			fmt.Println(err)
		}
	}

	w.Flush() // バッファに残っているデータを書き込む
}
