package main

import (
	"encoding/csv"
	"fmt"
	"flag"
	"os"
	
	"github.com/olekukonko/tablewriter"
)

func main() {
	filePath := flag.String("file","","Path to the Input File")
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
	
	fmt.Println("CSV Contents:")
	// for _, row := range data {
	// 	fmt.Println(row)
	// }
	
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