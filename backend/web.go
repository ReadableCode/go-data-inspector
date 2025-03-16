package main

import (
	"encoding/csv"

	"github.com/gofiber/fiber/v2"
)

var templatesDir = "../frontend/templates"

func hostSiteWithFiber() {
	app := fiber.New()

	// Serve frontend files
	app.Static("/", "../frontend/templates/index.html")
	app.Static("/css", "../frontend/css")

	// Upload route
	app.Post("/upload", handleUpload)

	app.Listen(":8505")
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

	return c.JSON(fiber.Map{"data": data})
}
