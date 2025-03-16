package main

import (
	"flag"
	"fmt"
)

func main() {
	cliMode := flag.Bool("cli", false, "Run in CLI mode")
	filePath := flag.String("file", "", "Path to the Input File")
	filter := flag.String("filter", "", "Filter condition (e.g., Age>30)")
	sortBy := flag.String("sort", "", "Column to sort by (e.g., Age)")
	desc := flag.Bool("desc", false, "Sort in descending order")
	interactive := flag.Bool("interactive", false, "Launch interactive mode")
	flag.Parse()

	if *cliMode {
		fmt.Println("Running in CLI mode...")
		runCLIMode(filePath, filter, sortBy, desc, interactive)
		return
	}

	hostSiteWithFiber()
}
