package main

import (
	"encoding/json"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Llongfile)
	path := "sample.json"
	j, _ := os.ReadFile(path)
	var x interface{}

	err := json.Unmarshal(j, &x)
	_ = err

	results := json2csv.json2csv(x)
	json2csv.writeCSV(results)
}
