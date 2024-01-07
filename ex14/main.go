package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

func isFibonacci(num int) bool {
	a, b := 0, 1
	for b < num {
		a, b = b, a+b
	}
	return b == num
}

func countFibonacciNumbers(numbers []string, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	count := 0
	for _, numStr := range numbers {
		num, err := strconv.Atoi(numStr)
		if err == nil && isFibonacci(num) {
			count++
		}
	}

	results <- count
}

func main() {
	content, err := os.ReadFile("ex14/input.txt")
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
			{"12", "5", "6", "13", "7"},
			{"21", "20", "42", "43", "8", "38"},
			{"54", "55", "34", "100"},
		}
	}

	results := make(chan int, len(input))
	var wg sync.WaitGroup

	// Map Step - Concurrently process each inner array of numbers
	for _, numbers := range input {
		wg.Add(1)
		go countFibonacciNumbers(numbers, results, &wg)
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
	fmt.Printf("The average number of Fibonacci numbers is: %.2f\n", average)
}
