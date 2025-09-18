package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func push_to_db() {
	// --- 1. Connect to the Database ---
	// This part remains the same.
	dsn := "host=localhost user=postgres password=test dbname=oj port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database.")

	// --- 2. Prepare Files and Scanners ---
	// Open both input and output files.
	inputFile, err := os.Open("../input.txt")
	if err != nil {
		log.Fatalf("Failed to read input.txt: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Open("../output.txt")
	if err != nil {
		log.Fatalf("Failed to read output.txt: %v", err)
	}
	defer outputFile.Close()

	// Create a scanner for each file to read them line by line.
	inputScanner := bufio.NewScanner(inputFile)
	outputScanner := bufio.NewScanner(outputFile)

	// --- 3. Loop and Insert Each Test Case ---
	problemID, err := uuid.Parse("750e8400-e29b-41d4-a716-446655440001")
	if err != nil {
		log.Fatalf("Failed to parse UUID: %v", err)
	}

	query := `INSERT INTO test_cases (id, problem_id, input, output, created_at) VALUES ($1, $2, $3, $4, $5)`
	testCaseCount := 0

	// Loop through the input file line by line.
	for inputScanner.Scan() {
		// For each line in the input, we must read a corresponding line from the output.
		if !outputScanner.Scan() {
			log.Fatalf("Error: output.txt has fewer lines than input.txt. Files are mismatched.")
		}

		inputLine := inputScanner.Text()
		outputLine := outputScanner.Text()
		testCaseID := uuid.New() // Generate a new ID for each test case.

		// Execute the insert query for the current line pair.
		_, err = db.Exec(query, testCaseID, problemID, inputLine, outputLine, time.Now())
		if err != nil {
			log.Fatalf("Failed to push test case #%d to DB: %v", testCaseCount+1, err)
		}

		testCaseCount++
		fmt.Printf("Successfully pushed test case #%d with ID: %s\n", testCaseCount, testCaseID)
	}

	// --- 4. Final Checks and Summary ---
	// Check for any errors that might have occurred during scanning.
	if err := inputScanner.Err(); err != nil {
		log.Fatalf("Error reading from input.txt: %v", err)
	}
	if err := outputScanner.Err(); err != nil {
		log.Fatalf("Error reading from output.txt: %v", err)
	}

	log.Printf("\n🎉 All done! Successfully pushed a total of %d test cases.", testCaseCount)
}

// Added a main function to make the file runnable
