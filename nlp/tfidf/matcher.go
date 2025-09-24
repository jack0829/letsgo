package tfidf

import (
	"bufio"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"github.com/derekparker/trie"
	"github.com/jack0829/letsgo/common/fs"
	"golang.org/x/exp/utf8string"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 头信息列
const (
	headerTotalDocs int = iota // 总样本文档数 (D)
	headerColumns              // 头信息列数
)

// 主体信息列
const (
	bodyWord    int = iota // 词
	bodyDocs               // 全部样本文档中出现该词的文档数 (d)
	bodyColumns            // 主体信息列数
)

// Matcher 生成器
type Matcher struct {
	docs float64
	tr   *trie.Trie
	sc   *Scanner
}

func New(dict ...*bufio.Scanner) *Matcher {

	m := &Matcher{
		tr: trie.New(),
		sc: NewScanner(),
	}

	for _, d := range dict {
		for d.Scan() {
			if w := strings.Trim(d.Text(), " \r\n\t"); w != "" {
				m.AddWord(w)
			}
		}
	}

	return m
}

// AddWord 添加词语
func (m *Matcher) AddWord(w string) {
	m.tr.Add(w, &Meta{
		Word: w,
	})
}

// RemoveWord 删除词语
func (m *Matcher) RemoveWord(w string) {
	m.tr.Remove(w)
}

// AddSample 添加样本，会改变词语 IDF
func (m *Matcher) AddSample(text *utf8string.String) {

	// 文档中的词（去重）
	words := make(map[string]*Meta)
	m.sc.Scan(m.tr, text, func(meta *Meta) {
		// fmt.Printf("Meta %s = %f\n", meta.Word, meta.Docs)
		words[meta.Word] = meta
	})

	// 包含某词的文档数
	for _, meta := range words {
		meta.addDoc()
	}

	// 总文档数
	m.docs++
}

func (m *Matcher) traversal(n *trie.Node, w *csv.Writer) error {

	// 当前节点写入
	if meta, ok := n.Meta().(*Meta); ok {
		if er := w.Write([]string{
			meta.Word,
			fmt.Sprintf("%.0f", meta.Docs),
		}); er != nil {
			return er
		}
	}

	// 子节点（递归）
	for _, c := range n.Children() {
		if er := m.traversal(c, w); er != nil {
			return er
		}
	}

	return nil
}

// Save 保存当前词典数据
func (m *Matcher) Save(w io.Writer) error {

	cw := csv.NewWriter(w)
	defer cw.Flush()

	// 写入头信息
	head := make([]string, headerColumns)
	head[headerTotalDocs] = fmt.Sprintf("%.0f", m.docs) // 总文档数
	if err := cw.Write(head); err != nil {
		return err
	}

	return m.traversal(m.tr.Root(), cw)

}

// SaveToFile 保存当前词典数据到文件
func (m *Matcher) SaveToFile(path string, withCompress bool) error {

	if err := fs.MustDir(filepath.Dir(path)); err != nil {
		return err
	}

	fd, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer fd.Close()

	if !withCompress {
		return m.Save(fd)
	}

	gz := gzip.NewWriter(fd)
	defer gz.Close()

	return m.Save(gz)
}

// Load 加载词典数据
func (m *Matcher) Load(r io.Reader) error {

	cr := csv.NewReader(r)
	cr.LazyQuotes = true

	// 读头信息
	cr.FieldsPerRecord = headerColumns
	head, err := cr.Read()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	if docs, e := strconv.ParseFloat(head[headerTotalDocs], 64); e == nil && docs > 0 {
		m.docs += docs // 累加总文档数 (D)
	}

	cr.FieldsPerRecord = bodyColumns
	for {

		row, er := cr.Read()
		if er == io.EOF {
			er = nil
			break
		}
		if er != nil {
			return er
		}

		word := strings.Trim(row[bodyWord], " \r\n\t")
		docs, _ := strconv.ParseFloat(row[bodyDocs], 64)

		if word == "" {
			break
		}

		// 已有词
		if n, ok := m.tr.Find(word); ok {
			if meta, ok := n.Meta().(*Meta); ok && meta != nil {
				meta.Docs += docs // 累加该词文档数 (d)
				continue
			}
		}

		// 新词
		m.tr.Add(word, &Meta{
			Word: word,
			Docs: docs,
		})

	}

	return nil
}

// LoadFromFile 从文件加载词典数据
func (m *Matcher) LoadFromFile(path string, withCompress bool) error {

	fd, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer fd.Close()

	if !withCompress {
		return m.Load(fd)
	}

	gz, err := gzip.NewReader(fd)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}
	defer gz.Close()

	return m.Load(gz)
}

// Match 匹配文本
func (m *Matcher) Match(text *utf8string.String) (words []*Word) {

	ts := NewScanner()
	ts.Scan(m.tr, text, nil)
	all := ts.GetCount() // 文档中关键词数量

	if Debug {
		if l := text.RuneCount(); l > 50 {
			fmt.Println(text.Slice(0, 25), "......", text.Slice(l-25, l))
		} else {
			fmt.Println(text.String())
		}
	}

	for _, w := range ts.GetWords() {
		words = append(words, w.SetScore(m.docs, all))
	}

	if Debug {
		fmt.Println("")
	}

	return
}

// TopN 获取权重最高的 N 个关键词
func (m *Matcher) TopN(n int, text *utf8string.String) []*Word {
	words := m.Match(text)
	// desc(words)
	return topN(n, words)
}
