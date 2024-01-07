package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

func countAbsoluteAndRelativePaths(paths []string, results chan<- [2]int, wg *sync.WaitGroup) {
	defer wg.Done()

	absoluteCount := 0
	relativeCount := 0

	for _, path := range paths {
		if strings.Contains(path, "/") {
			if strings.HasPrefix(path, "/") {
				absoluteCount++
			} else {
				relativeCount++
			}
		}
	}

	results <- [2]int{absoluteCount, relativeCount}
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
			{"/dev/null", "/bin", "saar", "teme/scoala/2020", ""},
			{"proiect/tema", "/dev", "ddanube", "jahfjksgfjhs", "ajsdas", "urs"},
			{"scoica", "/teme/repos/git", "arac", "karnak"},
		}
	}

	results := make(chan [2]int, len(input))
	var wg sync.WaitGroup

	// Map Step - Concurrently process each inner array of paths
	for _, paths := range input {
		wg.Add(1)
		go countAbsoluteAndRelativePaths(paths, results, &wg) // Map: Count absolute and relative paths
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	totalAbsoluteCount := 0
	totalRelativeCount := 0

	// Reduce Step - Aggregate results from map step
	for counts := range results {
		totalAbsoluteCount += counts[0] // Reduce: Aggregate counts of absolute paths
		totalRelativeCount += counts[1] // Reduce: Aggregate counts of relative paths
	}

	averageAbsolute := float64(totalAbsoluteCount) / float64(len(input))
	averageRelative := float64(totalRelativeCount) / float64(len(input))

	fmt.Printf("The average number of absolute paths is: %.2f\n", averageAbsolute)
	fmt.Printf("The average number of relative paths is: %.2f\n", averageRelative)
}
