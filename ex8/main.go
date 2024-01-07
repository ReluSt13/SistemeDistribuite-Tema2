package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

func countWordsWithSubstitutionPair(words []string, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	count := 0
	for _, word := range words {
		if hasSubstitutionPair(word, words) {
			count++
		}
	}

	results <- count
}

func hasSubstitutionPair(word string, words []string) bool {
	// Create a map for substitution cipher
	cipherMap := make(map[rune]rune)
	alphabet := "abcdefghijklmnopqrstuvwxyz"
	for i, r := range alphabet {
		cipherMap[r] = rune(alphabet[len(alphabet)-1-i])
	}

	substitution := ""
	for _, r := range word {
		if val, ok := cipherMap[r]; ok {
			substitution += string(val)
		}
	}

	// Check if the substitution exists in the word list
	for _, w := range words {
		if w == substitution {
			return true
		}
	}

	return false
}

func main() {
	content, err := os.ReadFile("ex8/input.txt")
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
			{"prajitura", "camion", "foaie", "liliac", "uezrv"},
			{"carte", "trofeu", "xzigw", "laptop", "scris", "muzica"},
			{"pictura", "telefon", "parapanta", "catel"},
		}
	}

	results := make(chan int, len(input))
	var wg sync.WaitGroup

	// Map Step - Concurrently process each inner array of words
	for _, words := range input {
		wg.Add(1)
		go countWordsWithSubstitutionPair(words, results, &wg)
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

	average := (float64(totalCount) / 2) / float64(len(input))
	fmt.Printf("The average number of words with substitution pair is: %.2f\n", average)
}
