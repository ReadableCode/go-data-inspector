package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Function to run Python scripts and capture output
func runPythonScript(script string, output *string, update chan<- struct{}) {
	cmd := exec.Command("python", script)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	scan := func(pipe *bufio.Reader) {
		for {
			line, err := pipe.ReadString('\n')
			if err != nil {
				break
			}
			*output += strings.TrimSpace(line) + "\n"
			update <- struct{}{} // Notify that new output is available
		}
	}

	go scan(bufio.NewReader(stdout))
	go scan(bufio.NewReader(stderr))
	cmd.Wait()
	close(update) // Signal that script execution is done
}

func main() {
	var output string
	update := make(chan struct{}, 1)

	// Run Python script in a goroutine
	go runPythonScript("script1.py", &output, update)

	// Wait for updates and print new output
	for range update {
		fmt.Println(output)
	}

	// Ensure all output is printed before exiting
	time.Sleep(1 * time.Second)
}
