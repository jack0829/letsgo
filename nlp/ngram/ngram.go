package ngram

import (
	"bufio"
	"context"
	"fmt"
	"github.com/derekparker/trie"
	"github.com/fumiama/jieba"
	"io"
	"os"
	"regexp"
	"unicode/utf8"
)

var (
	onlySymbolNumerals = regexp.MustCompile(`^[\d!@#$%^&*()、=_+\[\]{}\\|;:'",.<>/?]+$`)
)

type (
	NGram struct {
		jb        *jieba.Segmenter
		tr        *trie.Trie
		n         int
		threshold int
		docs      int
		freq      int
	}
	Meta struct {
		Word    string   `json:"word"`      // 词
		Freq    int      `json:"frequency"` // 词频
		Docs    int      `json:"docs"`      // 出现的文章数
		Token   []string `json:"-"`         // 子词元
		blocked bool     // 是否为屏蔽词
	}

	ThresholdFunc func(totalDocs, totalFreq int, meta *Meta) bool
)

func New(
	n int,
	baseDictPath, stopWordsPath string,
) (*NGram, error) {

	if n < 2 || n > 4 {
		return nil, fmt.Errorf("n 值范围必须为 2 至 4")
	}

	var (
		g   NGram
		err error
	)
	g.n = n
	if g.jb, err = jieba.LoadDictionaryAt(baseDictPath); err != nil {
		return nil, err
	} // , "", "", "", stopWordsPath)

	_ = g.jb.LoadUserDictionaryAt(stopWordsPath)

	g.tr = trie.New()
	if stopWordsPath != "" {
		if err := scanFile(stopWordsPath, func(word string) {
			g.tr.Add(word, &Meta{
				Word:    word,
				blocked: true,
			})
		}); err != nil {
			return nil, err
		}
	}

	return &g, nil
}

func (g *NGram) GetDocs() int {
	return g.docs
}

func (g *NGram) LoadDict(bgCtx context.Context, path string) error {

	ctx, stop := context.WithCancel(bgCtx)
	defer stop()

	ch, total, err := ReadDict(ctx, path)
	if err != nil {
		return err
	}
	g.docs += total
	for v := range ch {
		g.merge(v)
	}

	return nil
}

func (g *NGram) find(word string) (*Meta, bool) {

	if n, ok := g.tr.Find(word); ok {
		if m := n.Meta(); m != nil {
			v, ok := m.(*Meta)
			return v, ok
		}
	}

	return nil, false
}

func (g *NGram) merge(v *Meta) {

	if v == nil {
		return
	}

	if m, ok := g.find(v.Word); ok {
		m.merge(v)
		g.freq += v.Freq
		return
	}

	g.tr.Add(v.Word, v)
	g.freq += v.Freq
}

func (m *Meta) merge(v *Meta) {

	if v == nil {
		return
	}

	m.Docs += v.Docs
	m.Freq += v.Freq
	if v.blocked {
		m.blocked = true
	}
	if len(v.Token) < len(m.Token) {
		m.Token = v.Token
	}
}

func (g *NGram) match(text string) map[string]*Meta {

	words := g.jb.Cut(text, true)
	L := len(words)
	if L < 1 {
		return nil
	}

	tf := make(map[string]*Meta)

	for i := 0; i < L; i++ {

		j := i + g.n
		if j >= L {
			j = L
		}

		var s string
		window := words[i:j]
	Window:
		for n, w := range window {

			// if utf8.RuneCountInString(w) < 2 {
			// 	break
			// }

			if v, ok := g.find(w); ok && v.blocked {
				break Window
			}

			s += w

			if onlySymbolNumerals.MatchString(s) {
				continue
			}

			token := window[:n+1]
			if meta, ok := tf[s]; !ok {
				tf[s] = &Meta{
					Word:  s,
					Freq:  1,
					Token: token,
				}
			} else {
				meta.Freq++
				if len(meta.Token) < len(token) {
					meta.Token = token
				}
			}
		}
	}

	return tf
}

func (g *NGram) Sample(title, content string) {

	mt := g.match(title)
	mc := g.match(content)

	merged := make(map[string]struct{})

	threshold := g.threshold

	// 正文中的
	for s, meta := range mc {

		// 词频 >= 阈值
		if meta.Freq >= threshold {
			if _, ok := merged[s]; !ok {
				meta.Docs = 1
			}
			g.merge(meta)
			merged[s] = struct{}{}
			if _, ok := mt[s]; ok {
				delete(mt, s)
			}
		}
	}

	// 标题中的
	for s, meta := range mt {

		// 正文有出现 && 词组
		if _, ok := mc[s]; ok && len(meta.Token) > 1 {
			if _, ok := merged[s]; !ok {
				meta.Docs = 1
			}
			g.merge(meta)
			merged[s] = struct{}{}
		}
	}

	g.docs++
}

func (g *NGram) dumpTrieNode(n *trie.Node, out io.Writer, threshold ThresholdFunc) {

	if n == nil {
		return
	}

	if m := n.Meta(); m != nil {
		if w, ok := m.(*Meta); ok && w != nil {
			if !w.blocked && utf8.RuneCountInString(w.Word) >= 2 {

				var (
					isPrefix bool
					tf, df   int
				)
				for _, tail := range g.tr.PrefixSearch(w.Word) {

					v, ok := g.find(tail)
					if !ok || v.Word == w.Word {
						continue
					}

					tf += v.Freq
					df += v.Docs

					if v.Docs == w.Docs && v.Freq == w.Freq {
						isPrefix = true
						break
					}
				}

				if tf == w.Freq {
					isPrefix = true
				}

				if !isPrefix {
					if threshold == nil || threshold(g.docs, g.freq, w) {
						fmt.Fprintf(out, "%s\t%d\t%d\t%d\n", w.Word, w.Freq, w.Docs, len(w.Token))
					}
				}

			}
		}
	}

	for _, c := range n.Children() {
		g.dumpTrieNode(c, out, threshold)
	}
}

func (g *NGram) Dump(path string, threshold ThresholdFunc) error {

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "%d\n", g.docs)
	g.dumpTrieNode(g.tr.Root(), f, threshold)
	return nil
}

func scanFile(path string, fn func(line string)) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		if v := s.Text(); v != "" {
			fn(v)
		}
	}
	return nil
}

func DefaultDumpThreshold(totalDocs, totalFreq int, meta *Meta) bool {

	if meta == nil {
		return false
	}

	if totalDocs < 1 {
		return false
	}

	D, F, d, f := float64(totalDocs), float64(totalFreq), float64(meta.Docs), float64(meta.Freq)

	if f/F > 1e-1 {
		return false
	}

	idf := d / D
	fr := f / d
	// n := math.Log(D); d >= n &&
	if idf <= 1e-1 && fr > 2 && fr < 10 && (d > 2 || f > 15) {
		return true
	}

	if idf < 1e-3 && fr > 1 && d > 2 && len(meta.Token) < 2 {
		return true
	}

	return false
}
