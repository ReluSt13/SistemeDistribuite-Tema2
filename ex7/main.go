package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

func countWordsWithDiacritice(words []string, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	count := 0
	for _, word := range words {
		if countRomanianDiacritice(word) >= 2 {
			count++
		}
	}

	results <- count
}

func countRomanianDiacritice(word string) int {
	romanianDiacritice := []rune{'Â', 'Ă', 'Î', 'Ș', 'Ț', 'â', 'ă', 'î', 'ș', 'ț'}
	diacritice := 0
	for _, char := range word {
		if contains(romanianDiacritice, char) {
			diacritice++
		}
	}
	return diacritice
}

func contains(slice []rune, char rune) bool {
	for _, c := range slice {
		if c == char {
			return true
		}
	}
	return false
}

func main() {
	content, err := os.ReadFile("ex7/input.txt")
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
			{"țânțar", "carte", "ulcior", "copac", "plante"},
			{"beci", "", "mlăștinos", "astronaut", "stele", "planete"},
			{"floare", "somn", "șosetă", "scârțar"},
		}
	}

	results := make(chan int, len(input))
	var wg sync.WaitGroup

	// Map Step - Concurrently process each inner array of words
	for _, words := range input {
		wg.Add(1)
		go countWordsWithDiacritice(words, results, &wg)
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
	fmt.Printf("The average number of words with at least two diacritice is: %.2f\n", average)
}
