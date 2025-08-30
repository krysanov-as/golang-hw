package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	if text == "" {
		return nil
	}

	sliceUniqWord := strings.Fields(text)

	maps := make(map[string]int, len(text))

	initCounter := 1
	for _, v := range sliceUniqWord {
		_, ok := maps[v]

		if !ok {
			maps[v] = initCounter
		} else {
			maps[v]++
		}
	}

	resultSlice := make([]string, 0, len(maps))
	for w := range maps {
		resultSlice = append(resultSlice, w)
	}

	sort.Slice(resultSlice, func(i, j int) bool {
		iWords, jWords := resultSlice[i], resultSlice[j]
		if maps[iWords] != maps[jWords] {
			return maps[iWords] > maps[jWords]
		}
		return iWords < jWords
	})

	return resultSlice[:10]
}
