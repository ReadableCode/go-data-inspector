package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"flag"
	"os"
	"strconv"
	"strings"
	
	"github.com/olekukonko/tablewriter"
)

func main() {
	filePath := flag.String("file","","Path to the Input File")
	filter := flag.String("filter", "", "Filter condition (e.g., Age>30)")
	flag.Parse()

	// Check if the file path is provided
	if *filePath == "" {
		fmt.Println("Usage: goinspector --file <path to the input file>")
		os.Exit(1)
	}
	
	// Read the csv file with the csv reader function
	data, err := readCSV(*filePath)
	if err != nil {
		fmt.Println("error reading CSV:", err)
		os.Exit(1)
	}
	
	if *filter != "" {
		data, err = applyFilter(data, *filter)
		if err != nil {
			fmt.Println("Error applying filter:", err)
			os.Exit(1)
		}
	}
	
	printTable(data)
	
	fmt.Println("Go Data Inspector - CLI Tool")
	fmt.Println("Processing file:", *filePath)
}

func readCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	return reader.ReadAll() // reads the entire file into memory
}

func applyFilter(data [][]string, filter string) ([][]string, error) {
    if len(data) < 2 {
        return nil, errors.New("no data to filter")
    }

    // Parse the filter condition
    parts := strings.FieldsFunc(filter, func(r rune) bool {
        return r == '>' || r == '<' || r == '='
    })

    if len(parts) != 2 {
        return nil, errors.New("invalid filter format, use column>value | column<value | column=value")
    }

    column := strings.TrimSpace(parts[0])
    value := strings.TrimSpace(parts[1])

    // Extract the operator (supporting >= and <=)
    operator := filter[len(column) : len(column)+1]
    if len(filter) > len(column)+1 && (filter[len(column)+1] == '=' || filter[len(column)+1] == '>') {
        operator += string(filter[len(column)+1])
    }

    // Print parsed values for debugging
    fmt.Println("Column:", column)
    fmt.Println("Operator:", operator)
    fmt.Println("Value:", value)
	
	// Find the index of the column
	colIndex := -1
	for i, colName := range data[0] {  // for each column name in the first row
		if strings.EqualFold(colName, column) {
			colIndex = i
			break
		}
	}
	
	if colIndex == -1 {
		return nil, fmt.Errorf("column %s not found", column)
	}
	
	// Convert value to a number (if applicable)
	numValue, err := strconv.ParseFloat(value, 64)
	isNumeric := (err == nil)
	
	// Filter the data
	filteredData := [][]string{data[0]} // keep header
	
	for _, row := range data[1:] {
		if colIndex >= len(row) {  // handle the case where the column index is out of range
			continue
		}
		cellValue := row[colIndex]
		
		// Convert cell value to number (if applicable)
		cellNum, err := strconv.ParseFloat(cellValue, 64)
		cellIsNumeric := (err == nil)
		
		// Apply the filter condition
		match := false
		switch operator {
			case ">":
				match = isNumeric && cellIsNumeric && cellNum > numValue
			case "<":
				match = isNumeric && cellIsNumeric && cellNum < numValue
			case "=":
				match = strings.EqualFold(cellValue, value)
			case ">=":
				match = isNumeric && cellIsNumeric && cellNum >= numValue
			case "<=":
				match = isNumeric && cellIsNumeric && cellNum <= numValue
			default:
				return nil, errors.New("invalid operator, use >, <, =, >=, <=")
			}
		
		if match {
			filteredData = append(filteredData, row)
		}
	}
	
	return filteredData, nil

}

func printTable(data [][]string) {
	if len(data) == 0 {
		fmt.Println("No data to display")
		return
	}
	
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(data[0]) // first row is the header
	
	for _, row := range data[1:] {
		table.Append(row)
	}
	
	table.Render()
}