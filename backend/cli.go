package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/olekukonko/tablewriter"
	"github.com/rivo/tview"
)

func runCLIMode(filePath *string, filter *string, sortBy *string, desc *bool, interactive *bool) {
	if *filePath == "" {
		fmt.Println("Usage: goinspector --file <path-to-csv> [--filter column>value] [--sort column] [--desc] [--interactive]")
		os.Exit(1)
	}

	data, err := readCSV(*filePath)
	if err != nil {
		fmt.Println("Error reading CSV:", err)
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
	fmt.Println("Processing file:", *filePath)
}

func printTable(data [][]string) {
	if len(data) == 0 {
		fmt.Println("No data to display")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(data[0])

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
