package tfidf

// Meta Trie-Tree 节点数据
type Meta struct {
	Word string
	Docs float64 // 包含本词的文档数量
}

func (m *Meta) addDoc() {
	m.Docs++
}
