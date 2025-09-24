package tfidf

import (
	"github.com/derekparker/trie"
	"golang.org/x/exp/utf8string"
)

type Scanner struct {
	words map[string]*Word
	count float64 // 所有词出现次数总和
	list  []*Word
}

type MatchHandler func(*Meta)

func NewScanner() *Scanner {
	return &Scanner{
		words: make(map[string]*Word),
	}
}

func (ts *Scanner) add(meta *Meta) {

	ts.count++

	if w, ok := ts.words[meta.Word]; ok {
		w.Count++
		return
	}

	w := &Word{
		Word:  meta.Word,
		Count: 1,
		meta:  meta,
	}
	ts.words[meta.Word] = w
	ts.list = append(ts.list, w)
}

func (ts *Scanner) GetWords() []*Word {
	return ts.list
}

func (ts *Scanner) GetCount() float64 {
	return ts.count
}

// Scan 通过 trie 树扫描文档
func (ts *Scanner) Scan(tr *trie.Trie, text *utf8string.String, handler MatchHandler) {

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
				ts.add(meta)
				if handler != nil {
					handler(meta)
				}
				meta = nil
			}
			break
		}

		s := text.Slice(left, right)
		// t.Logf("DEBUG\t%s", s)

		if !tr.HasKeysWithPrefix(s) {

			// t.Logf("DEBUG\tNo prefix")
			left = right

			if meta != nil {
				ts.add(meta)
				if handler != nil {
					handler(meta)
				}
				// t.Logf("DEBUG\tADD %s", meta.Word)
				meta = nil
				left--
				right--
			}

			continue
		}

		if node, ok := tr.Find(s); ok {
			if v, ok := node.Meta().(*Meta); ok {
				// t.Logf("%s Matched", v.Word)
				meta = v
			}
		}

	}

	return
}
