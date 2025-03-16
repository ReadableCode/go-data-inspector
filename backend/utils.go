package main

import (
	"encoding/csv"
	"errors"
	"os"
	"sort"
	"strconv"
	"strings"
)

func readCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func applyFilter(data [][]string, filter string) ([][]string, error) {
	if len(data) < 2 {
		return nil, errors.New("no data to filter")
	}

	parts := strings.FieldsFunc(filter, func(r rune) bool {
		return r == '>' || r == '<' || r == '='
	})

	if len(parts) != 2 {
		return nil, errors.New("invalid filter format, use column>value | column<value | column=value")
	}

	column := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	operator := filter[len(column) : len(column)+1]

	if len(filter) > len(column)+1 && (filter[len(column)+1] == '=' || filter[len(column)+1] == '>') {
		operator += string(filter[len(column)+1])
	}

	colIndex := -1
	for i, colName := range data[0] {
		if strings.EqualFold(colName, column) {
			colIndex = i
			break
		}
	}
	if colIndex == -1 {
		return nil, errors.New("column not found")
	}

	numValue, err := strconv.ParseFloat(value, 64)
	isNumeric := (err == nil)

	filteredData := [][]string{data[0]}

	for _, row := range data[1:] {
		if colIndex >= len(row) {
			continue
		}
		cellValue := row[colIndex]
		cellNum, err := strconv.ParseFloat(cellValue, 64)
		cellIsNumeric := (err == nil)

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
			return nil, errors.New("invalid operator")
		}

		if match {
			filteredData = append(filteredData, row)
		}
	}

	return filteredData, nil
}

func sortCSV(data [][]string, sortBy string, descending bool) error {
	if len(data) < 2 {
		return errors.New("no data to sort")
	}

	colIndex := -1
	for i, colName := range data[0] {
		if strings.EqualFold(colName, sortBy) {
			colIndex = i
			break
		}
	}
	if colIndex == -1 {
		return errors.New("column not found")
	}

	isNumeric := true
	for _, row := range data[1:] {
		if _, err := strconv.ParseFloat(row[colIndex], 64); err != nil {
			isNumeric = false
			break
		}
	}

	sort.SliceStable(data[1:], func(i, j int) bool {
		if isNumeric {
			a, _ := strconv.ParseFloat(data[i+1][colIndex], 64)
			b, _ := strconv.ParseFloat(data[j+1][colIndex], 64)
			if descending {
				return a > b
			}
			return a < b
		} else {
			if descending {
				return data[i+1][colIndex] > data[j+1][colIndex]
			}
			return data[i+1][colIndex] < data[j+1][colIndex]
		}
	})

	return nil
}
