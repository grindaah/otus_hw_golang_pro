package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type WordCounter struct {
	Word  string
	Count int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Top10(s string) []string {
	if len(s) == 0 {
		return []string{}
	}
	wordsMap := make(map[string]int)
	wordsSlice := make([]WordCounter, 0)

	flds := strings.Fields(s)
	for _, fld := range flds {
		normalizedString := normalizeString(fld)
		if fld != "" {
			wordsMap[normalizedString]++
		}
	}
	for k, v := range wordsMap {
		wordsSlice = append(wordsSlice, WordCounter{Word: k, Count: v})
	}

	sort.Slice(wordsSlice, func(i, j int) bool {
		return wordsSlice[i].Count > wordsSlice[j].Count ||
			(wordsSlice[i].Count == wordsSlice[j].Count && wordsSlice[i].Word < wordsSlice[j].Word)
	})

	sz := min(10, len(wordsSlice))
	result := make([]string, sz)
	for i := 0; i < sz; i++ {
		result[i] = wordsSlice[i].Word
	}

	return result
}

func normalizeString(s string) string {
	// TODO find regexp for handling case-insensitive+ [,.]?
	// r := regexp.MustCompile("\\b[\.,]?\\b")
	// w := r.FindString(s)
	return s
}
