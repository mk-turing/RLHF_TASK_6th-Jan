package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func extractDataFromCSV(filename string) [][]string {
	var rows [][]string
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if err := reader.ReadHeader(); err != nil {
		log.Fatal(err)
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		rows = append(rows, row)
	}
	return rows
}

func transformData(rows [][]string) {
	transformedRows := make([][]string, 0, len(rows))

	// Sample transformation logic
	for _, row := range rows {
		key := row[0]                    // Assuming key is the first column
		value := strings.ToUpper(row[1]) // Sample transformation
		transformedRows = append(transformedRows, []string{key, value})
	}

	// Example aggregation using maps
	sums := make(map[string]int)
	for _, row := range rows {
		key := row[0]
		value, err := strconv.Atoi(row[2])
		if err != nil {
			continue // Skip non-integer values
		}
		sums[key] += value
	}

	fmt.Println("Aggregated sums:", sums)
}

func loadDataToCSV(filename string, transformedRows [][]string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.Write(transformedRows); err != nil {
		log.Fatal(err)
	}

	if err := writer.Flush(); err != nil {
		log.Fatal(err)
	}
}

func main() {
}
