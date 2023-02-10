package main

import (
	"reflect"
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
