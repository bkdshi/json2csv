package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
)

func flatten(obj interface{}) map[string]interface{} {
	f := make(map[string]interface{}, 0)
	// fmt.Println(obj)

	_flatten(f, "", obj)

	return f
}

func _flatten(f map[string]interface{}, key string, obj interface{}) {
	// interface to reflect.Value
	// value := reflect.ValueOf(obj)
	value, ok := obj.(reflect.Value)
	fmt.Println("ORI", value)
	fmt.Println("ORI", value.Kind())
	if !ok {
		value = reflect.ValueOf(obj)
		fmt.Println("!OK", value)
	}

	if value.Kind() == reflect.Interface {
		value = value.Elem()
		fmt.Println("ELE", value)
	}

	fmt.Println(key, value, value.Kind())

	switch value.Kind() {
	case reflect.Map:
		_flattenMap(f, value)
	case reflect.Slice:
		_flattenSlice(f, key, value)
	case reflect.Float64:
		f[key] = value.Float()
	case reflect.String:
		f[key] = value.String()
	}
}

func _flattenMap(f map[string]interface{}, value reflect.Value) {

	for _, key := range value.MapKeys() {
		// fmt.Println(key.String(), value.MapIndex(key))
		_flatten(f, key.String(), value.MapIndex(key))
	}
}

func _flattenSlice(f map[string]interface{}, key string, value reflect.Value) {
	if value.Len() == 0 {
		_flatten(f, key, "")
	}
	for i := 0; i < value.Len(); i++ {
		fmt.Println(value.Index(i))
		_flatten(f, key+"/"+strconv.Itoa(i), value.Index(i))
	}
}

func json2csv(data interface{}) []map[string]interface{} {
	var results []map[string]interface{}
	v := reflect.ValueOf(data)
	// fmt.Println(v.Kind())

	switch v.Kind() {
	case reflect.Map:
		// fmt.Println(v.Len())
		if v.Len() > 0 {
			results = append(results, flatten(v))
		}
	case reflect.Slice:
		// fmt.Println(v.Len())

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
	var x map[string]interface{}

	err := json.Unmarshal(j, &x)

	if err != nil {
		log.Fatal()
	}

	results := json2csv(x)
	writeCSV(results)

	// f, _ := os.Create("result.csv")

	// for _, result := range results {
	// 	fmt.Println(result)
	// 	// bytes, _ := json.Marshal(result)
	// 	// ioutil.WriteFile("result.csv", bytes, os.FileMode(0600))
	// }

}

func writeCSV(results []map[string]interface{}) {
	f, err := os.Create("result.csv")
	if err != nil {
		fmt.Println(err)
	}

	w := csv.NewWriter(f)

	headers := getHeader(results)
	fmt.Println(headers)
	// for _, result := range results {
	// 	fmt.Println(result)
	// }

	if err := w.Write(headers); err != nil {
		fmt.Println(err)
	}

	for _, result := range results {
		record := make([]string, 0, len(result))
		// _ = record
		for _, key := range headers {
			record = append(record, fmt.Sprintf("%v", result[key]))
			// fmt.Println(key)
		}
		if err := w.Write(record); err != nil {
			fmt.Println(err)
		}
	}

	w.Flush() // バッファに残っているデータを書き込む
}

func getHeader(results []map[string]interface{}) []string {
	headers := make([]string, 0, len(results))
	for _, result := range results {
		for k := range result {
			headers = append(headers, k)
		}
	}
	return headers
}
