package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type wordFrequency struct {
	word      string
	frequency int
}

type wordsFrequencies []*wordFrequency

func (w wordsFrequencies) Len() int      { return len(w) }
func (w wordsFrequencies) Swap(i, j int) { w[i], w[j] = w[j], w[i] }
func (w wordsFrequencies) Less(i, j int) bool {
	if w[i].frequency > w[j].frequency {
		return true
	}

	if w[i].frequency == w[j].frequency {
		return w[i].word < w[j].word
	}

	return false
}

var reWord = regexp.MustCompile(`(\p{L}(-+\p{L})?)+`)

func Top10(text string) []string {
	if text == "" {
		return nil
	}

	words := strings.Fields(text)
	hits := make(map[string]*wordFrequency, len(words))

	for _, w := range words {
		w = strings.ToLower(reWord.FindString(w))
		if w == "" {
			continue
		}

		if _, exist := hits[w]; exist {
			hits[w].frequency++

			continue
		}

		hits[w] = &wordFrequency{word: w, frequency: 1}
	}

	frequencies := make(wordsFrequencies, 0, len(hits))
	for _, hit := range hits {
		frequencies = append(frequencies, hit)
	}

	sort.Sort(frequencies)

	resCount := len(frequencies)
	if resCount > 10 {
		resCount = 10
	}

	res := make([]string, 0, resCount)
	for _, wf := range frequencies[:resCount] {
		res = append(res, wf.word)
	}

	return res
}
