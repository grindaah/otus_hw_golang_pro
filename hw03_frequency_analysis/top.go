package hw03frequencyanalysis

import (
	"regexp"
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
	wordsMap := make(map[string]int, 0)
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
		return wordsSlice[i].Count > wordsSlice[j].Count
	})
	sort.Slice(wordsSlice, func(i, j int) bool {
		return wordsSlice[i].Count == wordsSlice[j].Count && wordsSlice[i].Word < wordsSlice[j].Word
	})

	sz := min(10, len(wordsSlice))
	result := make([]string, sz)
	for i := 0; i < sz; i++ {
		result[i] = wordsSlice[i].Word
	}

	return result
}

func normalizeString(s string) string {
	r := regexp.MustCompile("\\b.*,?\\b")
	w := r.FindString(s)
	return w
}
