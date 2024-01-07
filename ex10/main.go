package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode"
)

func countWordsWithAlternatingPattern(words []string, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	count := 0
	for _, word := range words {
		if hasAlternatingPattern(word) {
			count++
		}
	}

	results <- count
}

func hasAlternatingPattern(word string) bool {
	if len(word) < 2 {
		return false
	}

	vowels := map[rune]bool{'a': true, 'e': true, 'i': true, 'o': true, 'u': true, 'ă': true, 'î': true, 'â': true}

	prevIsVowel := vowels[unicode.ToLower(rune(word[0]))]
	for i := 1; i < len(word); i++ {
		currIsVowel := vowels[unicode.ToLower(rune(word[i]))]
		if currIsVowel == prevIsVowel {
			return false
		}
		prevIsVowel = currIsVowel
	}

	return true
}

func main() {
	content, err := os.ReadFile("ex10/input.txt")
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
			{"caracatiță", "ceva", "saar", "aaastrfb", ""},
			{"aaabbbccc", "caporal", "ddanube", "jahfjksgfjhs", "ajsdas", "urs"},
			{"scoica", "coral", "arac", "karnak"},
		}
	}

	results := make(chan int, len(input))
	var wg sync.WaitGroup

	// Map Step - Concurrently process each inner array of words
	for _, words := range input {
		wg.Add(1)
		go countWordsWithAlternatingPattern(words, results, &wg) // Map: Count words with alternating pattern
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	totalCount := 0

	// Reduce Step - Aggregate results from map step
	for count := range results {
		totalCount += count // Reduce: Aggregate counts of words with alternating pattern
	}

	average := float64(totalCount) / float64(len(input))
	fmt.Printf("The average number of words with alternating pattern is: %.2f\n", average)
}
