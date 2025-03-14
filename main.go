package main

import (
	"fmt"
	"flag"
	"os"
)

func main() {
	filePath := flag.String("file","","Path to the Input File")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Usage: goinspector --file <path to the input file>")
		os.Exit(1)
	}
	
	fmt.Println("Go Data Inspector - CLI Tool")
	fmt.Println("Processing file:", *filePath)
}