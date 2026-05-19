package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	anagramGroups := searchAnagrams(words)
	fmt.Print(anagramGroups)
}

func searchAnagrams(words []string) map[string][]string {
	sortedToFirst := make(map[string]string)
	sortedToGroup := make(map[string][]string)

	for _, word := range words {
		lower := strings.ToLower(word)
		key := sortWord(lower)

		if _, exists := sortedToFirst[key]; !exists {
			sortedToFirst[key] = lower
		}
		sortedToGroup[key] = append(sortedToGroup[key], lower)
	}

	result := make(map[string][]string)
	for key, group := range sortedToGroup {
		if len(group) < 2 {
			continue
		}
		firstWord := sortedToFirst[key]
		sort.Strings(group)
		result[firstWord] = group
	}
	return result
}

func sortWord(word string) string {
	chars := strings.Split(word, "")
	sort.Strings(chars)
	sortedWord := strings.Join(chars, "")
	return sortedWord
}
