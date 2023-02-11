package json2csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
)

// Write CSV into file
func WriteCSV(f *os.File, results []map[string]interface{}) {
	w := csv.NewWriter(f)

	headers := getHeader(results)
	sort.SliceStable(headers, func(i, j int) bool { return headers[i] < headers[j] })

	if err := w.Write(headers); err != nil {
		fmt.Println(err)
	}

	for _, result := range results {
		record := make([]string, 0, len(result))
		for _, key := range headers {
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

// Get Header of CSV
func getHeader(results []map[string]interface{}) []string {
	headers := make([]string, 0)
	for _, result := range results {
		for k := range result {
			flag := true
			for _, head := range headers {
				if k == head {
					flag = false
					break
				}
			}
			if flag {
				headers = append(headers, k)
			}
		}
	}

	return headers
}
