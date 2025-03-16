package main

import (
	"encoding/csv"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var originalData [][]string // Holds the original uploaded data
var storedData [][]string   // Holds current working data (filtered/sorted)

func hostSiteWithFiber() {
	app := fiber.New()

	// Serve static files (CSS, JS, etc.)
	app.Static("/static", "../frontend/static")
	app.Static("/css", "../frontend/css")

	// Serve HTML template
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("../frontend/templates/index.html")
	})

	// Upload CSV
	app.Post("/upload", handleUpload)

	// Filter CSV
	app.Get("/filter", handleFilter)

	// Sort CSV
	app.Get("/sort", handleSort)

	// Reset CSV
	app.Get("/reset", handleReset)

	// Start server
	fmt.Println("Server running on http://localhost:8505")
	err := app.Listen(":8505")
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func handleUpload(c *fiber.Ctx) error {
	file, err := c.FormFile("csvfile")
	if err != nil {
		return c.Status(400).SendString("Error: No file uploaded")
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(500).SendString("Failed to read file")
	}
	defer src.Close()

	reader := csv.NewReader(src)
	data, err := reader.ReadAll()
	if err != nil {
		return c.Status(500).SendString("Failed to parse CSV")
	}

	// Store both the original and working datasets
	originalData = make([][]string, len(data))
	copy(originalData, data)

	storedData = make([][]string, len(data))
	copy(storedData, data)

	return c.JSON(fiber.Map{"data": storedData})
}

func handleReset(c *fiber.Ctx) error {
	if originalData == nil {
		return c.Status(400).SendString("No data uploaded yet")
	}

	// Reset storedData to the original data
	storedData = make([][]string, len(originalData))
	copy(storedData, originalData)

	return c.JSON(fiber.Map{"data": storedData})
}

func handleFilter(c *fiber.Ctx) error {
	column := c.Query("column")
	condition := c.Query("condition")

	if storedData == nil || len(storedData) < 2 {
		return c.Status(400).SendString("No data uploaded yet")
	}

	filteredData, err := applyFilter(storedData, column+condition)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// Store the filtered data so sorting works on it
	storedData = filteredData

	return c.JSON(fiber.Map{"data": filteredData})
}

func handleSort(c *fiber.Ctx) error {
	column := c.Query("column")
	desc, _ := strconv.ParseBool(c.Query("desc"))

	if storedData == nil || len(storedData) < 2 {
		return c.Status(400).SendString("No data available to sort")
	}

	err := sortCSV(storedData, column, desc)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	return c.JSON(fiber.Map{"data": storedData})
}
