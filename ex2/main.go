package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

func countPalindromes(words []string, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	count := 0
	for _, word := range words {
		if isPalindrome(word) {
			count++
		}
	}

	results <- count
}

func isPalindrome(word string) bool {
	runes := []rune(word)
	for i := 0; i < len(runes)/2; i++ {
		if runes[i] != runes[len(runes)-1-i] {
			return false
		}
	}
	return true
}

func main() {
	content, err := os.ReadFile("ex2/input.txt")
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
			{"a1551a", "parc", "ana", "minim", "1pcl3"},
			{"calabalac", "tivit", "leu", "zece10", "ploaie", "9ana9"},
			{"lalalal", "tema", "papa", "ger"},
		}
	}

	results := make(chan int, len(input))
	var wg sync.WaitGroup

	// Map Step - Concurrently process each inner array of words
	for _, words := range input {
		wg.Add(1)
		go countPalindromes(words, results, &wg)
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
	fmt.Printf("The average number of palindromes is: %.2f\n", average)
}
