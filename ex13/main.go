package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

func countNamesEndingWithEscu(names []string, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	count := 0
	for _, name := range names {
		if strings.HasSuffix(name, "escu") {
			count++
		}
	}

	results <- count
}

func main() {
	content, err := os.ReadFile("ex13/input.txt")
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
			{"Popescu", "Ionescu", "Pop", "aaastrfb", ""},
			{"Nicolae", "Dumitrescu", "ddanube", "jahfjksgfjhs", "ajsdas", "urs"},
			{"Dumitru", "Angelescu", "arac", "karnak"},
		}
	}

	results := make(chan int, len(input))
	var wg sync.WaitGroup

	// Map Step - Concurrently process each inner array of names
	for _, names := range input {
		wg.Add(1)
		go countNamesEndingWithEscu(names, results, &wg)
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
	fmt.Printf("The average number of names ending with 'escu' is: %.2f\n", average)
}
