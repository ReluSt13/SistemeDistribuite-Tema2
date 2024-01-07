package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

func countAnagrams(words []string, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	count := 0
	targetWord := "facultate"
	sortedTarget := sortStringChars(targetWord)

	for _, word := range words {
		if len(word) == len(targetWord) && sortStringChars(word) == sortedTarget {
			count++
		}
	}

	results <- count
}

func sortStringChars(s string) string {
	chars := []rune(s)

	sort.Slice(chars, func(i, j int) bool {
		return chars[i] < chars[j]
	})

	return string(chars)
}

func main() {
	content, err := os.ReadFile("ex5/input.txt")
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
			{"acultatef", "parc", "cultateaf", "faculatet", "copil"},
			{"cer", "tatefacul", "leu", "alee", "golang", "ultatefac"},
			{"tefaculta", "impar", "apa", "eleve"},
		}
	}

	results := make(chan int, len(input))
	var wg sync.WaitGroup

	// Map Step - Concurrently process each inner array of words
	for _, words := range input {
		wg.Add(1)
		go countAnagrams(words, results, &wg)
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
	fmt.Printf("The average number of anagrams of 'facultate' is: %.2f\n", average)
}
