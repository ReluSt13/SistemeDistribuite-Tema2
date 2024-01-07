package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

func numarVocale(cuvant string) int {
	count := 0
	vocale := "aeiouAEIOU"
	for _, char := range cuvant {
		if strings.ContainsRune(vocale, char) {
			count++
		}
	}
	return count
}

func processArray(innerArray []string, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	count := 0
	for _, word := range innerArray {
		if (numarVocale(word)%2 == 0) && ((len(word)-numarVocale(word))%3 == 0) {
			count++
		}
	}
	// Send the result to the channel
	results <- count
}

func main() {
	content, err := os.ReadFile("ex1/input.txt")
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
			{"aabbb", "ebep", "blablablaa", "hijk", "wsww"},
			{"abba", "eeeppp", "cocor", "ppppppaa", "qwerty", "acasq"},
			{"lalala", "lalal", "papapa", "papap"},
		}
	}

	results := make(chan int, len(input))
	var wg sync.WaitGroup
	// Map step - process each subarray in a goroutine
	for _, subArray := range input {
		wg.Add(1)
		go processArray(subArray, results, &wg)
	}
	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(results)
	}()
	// Reduce step - sum the results
	finalCount := 0
	for result := range results {
		finalCount += result
	}
	fmt.Printf("Numarul mediu de cuvinte care indeplinesc conditiile este: %.2f\n", float64(finalCount)/float64(len(input)))
}
