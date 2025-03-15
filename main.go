package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"flag"
	"os"
	"strconv"
	"strings"
	"sort"
	
	"github.com/gdamore/tcell/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/rivo/tview"
)

func main() {
	// Collect command-line arguments
	filePath := flag.String("file","","Path to the Input File")
	filter := flag.String("filter", "", "Filter condition (e.g., Age>30)")
	sortBy := flag.String("sort", "", "Column to sort by (e.g., Age)")
	desc := flag.Bool("desc", false, "Sort in descending order")
	interactive := flag.Bool("interactive", false, "Launch interactive mode")
	flag.Parse()

	// Check if the file path is provided
	if *filePath == "" {
        fmt.Println("Usage: goinspector --file <path-to-csv> [--filter column>value] [--sort column] [--desc] [--interactive]")
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
	
	if *sortBy != "" {
		err = sortCSV(data, *sortBy, *desc)
		if err != nil {
			fmt.Println("Error sorting CSV:", err)
			os.Exit(1)
		}
	}
	
	if *interactive {
		runInteractiveTable(data)
		return
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

func sortCSV(data [][]string, sortBy string, descending bool) error {
    if len(data) < 2 {
        return errors.New("no data to sort")
    }

    // Find the column index
    colIndex := -1
    for i, colName := range data[0] {
        if strings.EqualFold(colName, sortBy) {
            colIndex = i
            break
        }
    }
    if colIndex == -1 {
        return fmt.Errorf("column '%s' not found", sortBy)
    }

    // Check if column is numeric
    isNumeric := true
    for _, row := range data[1:] {
        if _, err := strconv.ParseFloat(row[colIndex], 64); err != nil {
            isNumeric = false
            break
        }
    }

    // Sorting function
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

func runInteractiveTable(data [][]string) {
    app := tview.NewApplication()

    table := tview.NewTable().
        SetBorders(true)

    // Set column headers
    for col, header := range data[0] {
        table.SetCell(0, col,
            tview.NewTableCell(header).
                SetTextColor(tcell.ColorYellow).
                SetSelectable(false).
                SetAlign(tview.AlignCenter))
    }

    // Add rows to the table
    for rowIndex, row := range data[1:] {
        for colIndex, cell := range row {
            table.SetCell(rowIndex+1, colIndex,
                tview.NewTableCell(cell).
                    SetTextColor(tcell.ColorWhite).
                    SetAlign(tview.AlignLeft))
        }
    }

    // Allow arrow key navigation
    table.SetSelectable(true, false)

    // Set up UI layout
    flex := tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(tview.NewTextView().
            SetText("Press ESC to exit").
            SetTextColor(tcell.ColorGreen), 1, 1, false).
        AddItem(table, 0, 1, true)

    // Handle key events
    table.SetDoneFunc(func(key tcell.Key) {
        if key == tcell.KeyEscape {
            app.Stop()
        }
    })

    // Run the app
    if err := app.SetRoot(flex, true).Run(); err != nil {
        fmt.Println("Error running interactive mode:", err)
    }
}
