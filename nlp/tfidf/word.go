package tfidf

import (
	"math"
	"sort"
)

type Word struct {
	Word  string  `json:"word"`
	Count float64 `json:"count,omitempty"`
	Score float64 `json:"score,omitempty"`
	meta  *Meta
}

func (w *Word) SetScore(
	allDocs, // 全部文档数量
	allWordCount float64, // 单篇文档中所有词的出现次数
) *Word {
	tf, idf := GetTF(w.Count, allWordCount), GetIDF(w.meta.Docs, allDocs)
	w.Score = tf * idf
	debug("%s\tTF = %.0f/%.0f = %f, IDF = lg(%.0f/%.0f) = %f, @weight = %f\n",
		w.Word, w.Count, allWordCount, tf, w.meta.Docs, allDocs, idf, w.Score,
	)
	return w
}

func (w *Word) GetDocs() (d int64) {

	if w.meta != nil {
		d = int64(w.meta.Docs)
	}

	return
}

func GetIDF(
	d, // 出现该词的文档数
	D float64, // 文档总数
) float64 {
	if d <= 0 {
		return 1
	}
	return math.Log10(D / d)
}

func GetTF(
	n, // 文档中出现该词次数
	N float64, // 文档中所有词出现次数总和
) float64 {
	if N <= 0 {
		return 1
	}
	return n / N
}

func desc(words []*Word) {
	sort.Slice(words, func(i, j int) bool {
		return words[i].Score > words[j].Score
	})
}

func topN(n int, words []*Word) []*Word {

	desc(words)

	if n < 0 {
		return words
	}

	if n >= len(words) {
		return words
	}

	return words[:n]
}
