package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
)

// Function to run a Python script concurrently and stream its output
func runPythonScript(script string, wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.Command("python", script)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("[Error]: Failed to create stdout pipe for", script, ":", err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("[Error]: Failed to create stderr pipe for", script, ":", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("[Error]: Failed to start", script, ":", err)
		return
	}

	// Ensure goroutines for reading stdout and stderr start before waiting
	var outputWg sync.WaitGroup
	outputWg.Add(2)

	go func() {
		defer outputWg.Done()
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println("[Output "+script+"]:", scanner.Text())
		}
	}()

	go func() {
		defer outputWg.Done()
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println("[Error "+script+"]:", scanner.Text())
		}
	}()

	// Wait for both output goroutines to finish reading before waiting for process exit
	outputWg.Wait()

	// Now wait for the script to fully exit
	if err := cmd.Wait(); err != nil {
		fmt.Println("[Error]:", script, "exited with error:", err)
	}
}

func main() {
	var wg sync.WaitGroup

	// Start both scripts concurrently
	wg.Add(2)
	go runPythonScript("script1.py", &wg)
	go runPythonScript("script2.py", &wg)

	// Wait for both scripts to finish
	wg.Wait()
}
