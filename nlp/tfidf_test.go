package nlp_test

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/jack0829/letsgo/nlp/tfidf"
	"golang.org/x/exp/utf8string"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {

	tfidf.Debug = true

	m.Run()
}

func openFile(ctx context.Context, path string, flag int) (fd *os.File, err error) {

	if fd, err = os.OpenFile(path, flag, 0755); err != nil {
		return
	}

	go func() {
		<-ctx.Done()
		fd.Close()
	}()

	return
}

func loadFile(path string) (*utf8string.String, error) {

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return utf8string.NewString(string(b)), nil
}

func TestIDFMatch(t *testing.T) {

	m := tfidf.New()
	m.AddWord("ab")
	m.AddWord("abc")

	txt := utf8string.NewString("abcabcab0acabcabc")
	for _, w := range m.Match(txt) {
		t.Logf("%s\t%.0f", w.Word, w.Count)
	}
}

func TestIDFSave(t *testing.T) {

	ctx := context.Background()

	dict, err := openFile(ctx, "./samples/keywords.input", os.O_RDONLY)
	if err != nil {
		t.Error(err)
		return
	}

	m := tfidf.New(bufio.NewScanner(dict))

	path, err := filepath.Glob("./samples/weixin/*.txt")
	if err != nil {
		t.Error(err)
		return
	}

	// TF-IDF
	for _, p := range path {

		if d, err := loadFile(p); err != nil {
			t.Error(err)
			return
		} else {
			m.AddSample(d)
		}
	}

	// Save Data
	if err = m.SaveToFile("./samples/dict.csv", false); err != nil {
		t.Error(err)
		return
	}
}

func TestIDFLoad(t *testing.T) {

	m := tfidf.New()
	if err := m.LoadFromFile("./samples/dict.csv", false); err != nil {
		t.Error(err)
		return
	}

	path, err := filepath.Glob("./samples/weixin/*.txt")
	if err != nil {
		t.Error(err)
		return
	}

	out, _ := os.OpenFile("./samples/keywords.output", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	defer out.Close()
	for _, p := range path {

		d, err := loadFile(p)
		if err != nil {
			t.Error(err)
			return
		}

		// t.Log(p)
		fmt.Fprintln(out, p)
		for _, w := range m.TopN(3, d) {
			// t.Logf("%s\t%.0f\t%f", w.Word, w.Count, w.Score)
			fmt.Fprintf(out, "%s\t%.0f\t%.6f\n", w.Word, w.Count, w.Score)
		}
		fmt.Fprintln(out, "")

	}

}

func TestExtractCSV(t *testing.T) {

	matches, err := filepath.Glob(`./samples/csv/*.csv`)
	if err != nil {
		t.Error(err)
		return
	}

	writeFile := func(id, title, content string) error {
		fd, er := os.OpenFile(fmt.Sprintf("./samples/weixin/%s.txt", id), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
		if er != nil {
			return er
		}
		defer fd.Close()
		io.WriteString(fd, title)
		io.WriteString(fd, "\n")
		io.WriteString(fd, content)
		return nil
	}

	extract := func(path string) error {

		fd, er := os.OpenFile(path, os.O_RDONLY, 0755)
		if er != nil {
			return er
		}
		defer fd.Close()

		r := csv.NewReader(fd)
		r.LazyQuotes = true

		r.Read() // 跳过首行
		for {

			row, er := r.Read()
			if er == io.EOF {
				break
			}

			if er != nil {
				return er
			}

			id := strings.Trim(row[10], " \n\r\t")
			author := strings.Trim(row[4], " \n\r\t")
			title := strings.Trim(row[11], " \n\r\t")
			content := strings.Trim(row[14], " \n\r\t")

			if author != "首都广电" {
				continue
			}

			if id != "" && content != "" {
				if er = writeFile(id, title, content); er != nil {
					return er
				}
			}
		}
		return nil
	}

	for _, path := range matches {
		if err = extract(path); err != nil {
			t.Error(err)
			return
		}
		t.Log(path)
	}
}
