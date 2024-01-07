package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode"
)

func countStrongPasswords(words []string, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	count := 0
	for _, word := range words {
		if isStrongPassword(word) {
			count++
		}
	}

	results <- count
}

func isStrongPassword(word string) bool {
	var (
		hasLowercase bool
		hasUppercase bool
		hasSymbol    bool
		hasNumber    bool
	)

	for _, char := range word {
		switch {
		case unicode.IsLower(char):
			hasLowercase = true
		case unicode.IsUpper(char):
			hasUppercase = true
		case unicode.IsDigit(char):
			hasNumber = true
		default:
			hasSymbol = true
		}
	}

	return hasLowercase && hasUppercase && hasSymbol && hasNumber
}

func main() {
	content, err := os.ReadFile("ex11/input.txt")
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
			{"sadsa1@A", "cevaA!4", "saar", "aaastrfb", ""},
			{"aaabbbccc", "!Caporal1", "ddanube", "jahfjksgfjhs", "ajsdas", "urs"},
			{"scoica", "Coral!@12", "arac", "karnak"},
		}
	}

	results := make(chan int, len(input))
	var wg sync.WaitGroup

	// Map Step - Concurrently process each inner array of words
	for _, words := range input {
		wg.Add(1)
		go countStrongPasswords(words, results, &wg) // Map: Count strong passwords
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	totalCount := 0

	// Reduce Step - Aggregate results from map step
	for count := range results {
		totalCount += count // Reduce: Aggregate counts of strong passwords
	}

	average := float64(totalCount) / float64(len(input))
	fmt.Printf("The average number of strong passwords is: %.2f\n", average)
}
