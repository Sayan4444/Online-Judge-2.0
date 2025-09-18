package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func generate_input() {
	// Seed the random number generator to get different values each time
	rand.Seed(time.Now().UnixNano())

	// Create a file named input.txt
	file, err := os.Create("../input.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Number of test cases
	t := 1000
	// The number of test cases 't' is no longer written to the file.

	totalN := 100000
	remainingN := totalN

	for i := 0; i < t; i++ {
		// Add a newline to separate test cases, except before the very first one.
		if i > 0 {
			fmt.Fprint(file, "\n")
		}

		var n int
		// For the last test case, use all remaining N
		if i == t-1 {
			n = remainingN
		} else {
			// We need to leave at least 2 elements for each of the future (t-1-i) test cases.
			minForFuture := 2 * (t - 1 - i)

			// The maximum value for n in this iteration.
			maxAllowedN := remainingN - minForFuture

			// To prevent test cases from being too skewed, let's set a soft upper bound
			// based on the average size of remaining test cases.
			avgN := remainingN / (t - i)
			upperBound := avgN * 2
			// Ensure the upper bound is at least 2.
			if upperBound < 2 {
				upperBound = 2
			}

			// Use the smaller of the two upper bounds.
			if maxAllowedN > upperBound {
				maxAllowedN = upperBound
			}

			// n must be at least 2, unless we don't have enough elements left.
			minAllowedN := 2

			if maxAllowedN < minAllowedN {
				// Not enough elements to guarantee n=2, take what's available.
				n = maxAllowedN
			} else {
				// Generate n in the range [minAllowedN, maxAllowedN]
				n = rand.Intn(maxAllowedN-minAllowedN+1) + minAllowedN
			}
		}

		if n < 0 {
			n = 0
		}

		remainingN -= n

		// Print n without a leading space.
		fmt.Fprint(file, n)

		if n == 0 {
			// For a 0-size array, the target is 0. Print with a leading space.
			fmt.Fprint(file, " 0")
			continue
		}

		nums := make([]int, n)
		for j := 0; j < n; j++ {
			nums[j] = rand.Intn(100001) // Random numbers between 0 and 100000
		}

		// To guarantee a solution exists, pick two random indices and set the target
		// as the sum of the elements at those indices.
		idx1 := rand.Intn(n)
		idx2 := rand.Intn(n)
		// Ensure indices are different if possible (for n > 1)
		if n > 1 {
			for idx1 == idx2 {
				idx2 = rand.Intn(n)
			}
		}

		target := nums[idx1] + nums[idx2]

		// Write the array elements to the file, each with a leading space.
		for _, num := range nums {
			fmt.Fprint(file, " ", num)
		}

		// Write the target with a leading space.
		fmt.Fprint(file, " ", target)
	}

	fmt.Println("Successfully generated input.txt")
}
