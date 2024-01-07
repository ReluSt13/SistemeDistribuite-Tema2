package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode"
)

func countWordsWithConditions(words []string, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	count := 0
	for _, word := range words {
		if startsAndEndsWithUppercase(word) && hasEvenLowercaseLetters(word) {
			count++
		}
	}

	results <- count
}

func startsAndEndsWithUppercase(word string) bool {
	if len(word) > 0 {
		firstChar, lastChar := rune(word[0]), rune(word[len(word)-1])
		return unicode.IsUpper(firstChar) && unicode.IsUpper(lastChar)
	}
	return false
}

func hasEvenLowercaseLetters(word string) bool {
	lowercaseCount := 0
	for _, char := range word {
		if unicode.IsLower(char) {
			lowercaseCount++
		}
	}
	return lowercaseCount%2 == 0
}

func main() {
	content, err := os.ReadFile("ex6/input.txt")
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
			{"AcasA", "CasA", "FacultatE", "SisTemE", "distribuite"},
			{"GolanG", "map", "reduce", "Problema", "TemA", "ProieCt"},
			{"LicentA", "semestru", "ALGORitM", "StuDent"},
		}
	}

	results := make(chan int, len(input))
	var wg sync.WaitGroup

	// Map Step - Concurrently process each inner array of words
	for _, words := range input {
		wg.Add(1)
		go countWordsWithConditions(words, results, &wg) // Map: Count words meeting conditions
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	totalCount := 0

	// Reduce Step - Aggregate results from map step
	for count := range results {
		totalCount += count // Reduce: Aggregate counts of words meeting conditions
	}

	average := float64(totalCount) / float64(len(input))
	fmt.Printf("The average number of words meeting the conditions is: %.2f\n", average)
}
