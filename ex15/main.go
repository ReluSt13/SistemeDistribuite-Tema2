package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

func countNumbersWithThreeOnes(numbers []string, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	count := 0
	for _, numStr := range numbers {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			continue
		}

		// Count the number of set bits (1s) in the binary representation
		bitCount := 0
		for num > 0 {
			if num&1 == 1 {
				bitCount++
			}
			num >>= 1 // Right shift to check each bit
		}

		// Check if the count of set bits is exactly 3
		if bitCount == 3 {
			count++
		}
	}

	results <- count
}

func main() {
	content, err := os.ReadFile("ex12/input.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	lines := strings.Split(string(content), "\n")

	var input [][]string

	for _, line := range lines {
		if line == "" {
			continue
		}
		words := strings.Fields(line)
		input = append(input, words)
	}

	if len(input) == 0 {
		// Use default array if no valid input found
		input = [][]string{
			{"1", "13", "6", "7", "9"},
			{"19", "20", "43", "43", "21", "53"},
			{"54", "55", "28", "101"},
		}
	}

	results := make(chan int, len(input))
	var wg sync.WaitGroup

	// Map Step - Concurrently process each inner array of numbers
	for _, numbers := range input {
		wg.Add(1)
		go countNumbersWithThreeOnes(numbers, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	totalCount := 0

	// Reduce Step - Aggregate results from map step
	for count := range results {
		totalCount += count
	}

	average := float64(totalCount) / float64(len(input))
	fmt.Printf("The average number of numbers with exactly 3 ones in binary representation is: %.2f\n", average)
}
