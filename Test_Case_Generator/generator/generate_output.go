package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
)

func generate_output() {
	inputFile := "../input.txt"
	outputFile := "../output.txt"
	executablePath := "../two_sum"

	// 1. Open the input file for reading.
	in, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Failed to open input file %s: %v", inputFile, err)
	}
	defer in.Close()

	// 2. Create/overwrite the output file for writing the results.
	out, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create output file %s: %v", outputFile, err)
	}
	defer out.Close()

	// 3. Create a scanner to read the input file line by line.
	scanner := bufio.NewScanner(in)
	lineNum := 0

	log.Println("Starting to process test cases...")

	// 4. Loop through each line (each test case) in the input file.
	for scanner.Scan() {
		lineNum++
		line := scanner.Text() // Get the current test case as a string.

		// Create a new command for the executable for each line.
		cmd := exec.Command(executablePath)

		// Pipe the current line as the standard input for the command.
		cmd.Stdin = strings.NewReader(line)

		// Run the command and capture its standard output.
		outputBytes, err := cmd.Output()
		if err != nil {
			// If the command fails, log the error and the line that caused it.
			log.Fatalf("Error executing on line %d ('%s'): %v", lineNum, line, err)
		}

		// Write the captured output to our output file.
		// We assume your C++ program already adds a newline to its output.
		if _, err := out.Write(outputBytes); err != nil {
			log.Fatalf("Failed to write result to output file: %v", err)
		}
	}

	// Check for any errors during the scanning process.
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	log.Printf("Successfully processed %d test cases into %s.", lineNum, outputFile)
}