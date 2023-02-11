package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/bkdshi/json2csv"
)

func main() {
	log.SetFlags(log.Llongfile)
	path := "sample.json"
	j, _ := os.ReadFile(path)
	var x interface{}

	err := json.Unmarshal(j, &x)
	_ = err

	results := json2csv.Json2Csv(x)
	f, err := os.Create("result.csv")
	json2csv.WriteCSV(f, results)
}
