package ngram

import (
	"github.com/jack0829/letsgo/nlp/tfidf"
	"golang.org/x/exp/utf8string"
	"sort"
)

type (
	ScanResult struct {
		data  map[string]*Meta
		freq  int // 所有词出现次数总和
		metas []*Meta
		words []*Word
	}
	Word struct {
		Word  string  `json:"word"`
		Score float64 `json:"score,omitempty"`
		Count int     `json:"count,omitempty"`
	}
)

func (sr *ScanResult) add(meta *Meta) {

	sr.freq++

	if w, ok := sr.data[meta.Word]; ok {
		w.Freq++
		return
	}

	meta.Freq = 1
	sr.data[meta.Word] = meta
	sr.metas = append(sr.metas, meta)
}

func (sr *ScanResult) GetTotalFrequency() int {
	return sr.freq
}

func (sr *ScanResult) dataToWords(totalDocs int) {

	sr.words = make([]*Word, 0, len(sr.data))
	for _, m := range sr.data {
		w := &Word{
			Word:  m.Word,
			Count: m.Freq,
			Score: tfidf.GetTF(float64(m.Freq), float64(sr.freq)) * tfidf.GetIDF(float64(m.Docs), float64(totalDocs)),
		}
		sr.words = append(sr.words, w)
		// fmt.Println(w.Word, w.Count, w.Score)
	}

	sort.Slice(sr.words, func(i, j int) bool {
		return sr.words[i].Score > sr.words[j].Score
	})
}

func (sr *ScanResult) TopN(n int) []*Word {

	if n > len(sr.words) {
		n = len(sr.words)
	}

	return sr.words[:n]
}

// Scan 通过 trie 树扫描文档
func (g *NGram) Scan(s string) *ScanResult {

	sr := &ScanResult{
		data: make(map[string]*Meta),
	}

	text := utf8string.NewString(s)
	if text.RuneCount() < 1 {
		return sr
	}

	var (
		left  int   // 窗口左界
		right int   // 窗口右界
		meta  *Meta // 窗口内匹配到信息
	)

	L := text.RuneCount()
	for left < L {

		right++

		if right > L {
			if meta != nil {
				sr.add(meta)
				meta = nil
			}
			break
		}

		str := text.Slice(left, right)
		// t.Logf("DEBUG\t%s", s)

		if !g.tr.HasKeysWithPrefix(str) {

			// t.Logf("DEBUG\tNo prefix")
			left = right

			if meta != nil {
				sr.add(meta)

				// t.Logf("DEBUG\tADD %s", meta.Word)
				meta = nil
				left--
				right--
			}

			continue
		}

		if v, ok := g.find(str); ok {
			if !v.blocked {
				meta = v
				// fmt.Println("find ", str, ": ", v.Word)
			}
		}
	}

	sr.dataToWords(g.docs)

	return sr
}

func (g *NGram) Match(text string) *ScanResult {

	var r ScanResult
	r.data = g.match(text)
	for _, m := range r.data {
		r.freq += m.Freq
	}

	r.dataToWords(g.docs)

	return &r
}
