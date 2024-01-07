package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

func countVowelStartEndWords(words []string, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	count := 0
	for _, word := range words {
		if startsAndEndsWithVowel(word) {
			count++
		}
	}

	results <- count
}

func startsAndEndsWithVowel(word string) bool {
	vowels := map[rune]bool{'a': true, 'e': true, 'i': true, 'o': true, 'u': true}
	if len(word) > 0 {
		firstChar, lastChar := rune(word[0]), rune(word[len(word)-1])
		return vowels[firstChar] && vowels[lastChar]
	}
	return false
}

func main() {
	content, err := os.ReadFile("ex4/input.txt")
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
			{"ana", "parc", "impare", "era", "copil"},
			{"cer", "program", "leu", "alee", "golang", "info"},
			{"inima", "impar", "apa", "eleve"},
		}
	}

	results := make(chan int, len(input))
	var wg sync.WaitGroup

	// Map Step - Concurrently process each inner array of words
	for _, words := range input {
		wg.Add(1)
		go countVowelStartEndWords(words, results, &wg) // Map: Count words starting and ending with vowels
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	totalCount := 0

	// Reduce Step - Aggregate results from map step
	for count := range results {
		totalCount += count // Reduce: Aggregate counts of words starting and ending with vowels
	}

	average := float64(totalCount) / float64(len(input))
	fmt.Printf("The average number of words starting and ending with vowels is: %.2f\n", average)
}
