package json2csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

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

func getHeader(results []map[string]interface{}) []string {
	headers := make([]string, 0)
	for _, result := range results {
		for k := range result {
			if len(results) > 1 {
				tmp := strings.Join(strings.Split(k, "/")[1:], "/")
				flag := true
				for _, head := range headers {
					if tmp == head {
						flag = false
						break
					}
				}
				if flag {
					headers = append(headers, tmp)
				}
			} else {
				headers = append(headers, k)
			}

		}
	}
	rm_list := make([]string, 0)

	for i := 0; i < len(headers); i++ {
		tmp := headers[i] + "/0"
		// fmt.Println(tmp)
		for j := 0; j < len(headers); j++ {
			// fmt.Println("\t", headers[j])
			if strings.HasPrefix(headers[j], tmp) {
				rm_list = append(rm_list, headers[i])
			}
		}
	}
	// fmt.Println(rm_list)

	new_headers := make([]string, 0)
	for _, header := range headers {
		flag := true
		for _, rm := range rm_list {
			if header == rm {
				flag = false
				break
			}
		}
		if flag {
			new_headers = append(new_headers, header)
		}
	}

	return new_headers
}

func cloneKey(keys []string) []string {
	if len(keys) == 0 {
		return keys
	}

	cloneKey := make([]string, len(keys))
	copy(cloneKey, keys)
	return cloneKey
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
