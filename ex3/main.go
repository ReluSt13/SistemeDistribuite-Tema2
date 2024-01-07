package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

func countBirdLanguageWords(words []string, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	count := 0
	for _, word := range words {
		if isBirdLanguageWord(word) {
			count++
		}
	}

	results <- count
}

func isBirdLanguageWord(word string) bool {
	vowels := map[rune]bool{'a': true, 'e': true, 'i': true, 'o': true, 'u': true}
	for i, char := range word {
		if vowels[char] {
			if i < len(word)-1 && word[i+1] != 'p' {
				return false
			}
		}
	}
	return true
}

func main() {
	content, err := os.ReadFile("ex3/input.txt")
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
			{"apap", "paprc", "apnap", "mipnipm", "copil"},
			{"cepr", "program", "lepu", "zepcep", "golang", "tema"},
			{"par", "impar", "papap", "gepr"},
		}
	}

	results := make(chan int, len(input))
	var wg sync.WaitGroup

	// Map Step - Concurrently process each inner array of words
	for _, words := range input {
		wg.Add(1)
		go countBirdLanguageWords(words, results, &wg)
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
	fmt.Printf("The average number of bird language words is: %.1f\n", average)
}
