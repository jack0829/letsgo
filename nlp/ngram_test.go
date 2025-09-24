package nlp_test

import (
	"context"
	"crypto/md5"
	"encoding/csv"
	"fmt"
	"github.com/jack0829/letsgo/nlp/ngram"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testExtractCsv(csvPath, outputDir string) func(t *testing.T) {

	writeFile := func(title, content string) error {

		h := md5.New()
		io.WriteString(h, title)

		f, err := os.Create(filepath.Join(outputDir, fmt.Sprintf("%x", h.Sum(nil))))
		if err != nil {
			return err
		}
		defer f.Close()

		io.WriteString(f, title)
		io.WriteString(f, "\n\n")
		io.WriteString(f, content)

		return nil
	}

	return func(t *testing.T) {

		f, err := os.Open(csvPath)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		r := csv.NewReader(f)
		r.LazyQuotes = true

		if _, err = r.Read(); err != nil {
			t.Fatal(err)
		}

		for {

			row, err := r.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				t.Logf("read csv err: %v", err)
				continue
			}

			title := strings.Trim(row[0], "\n\t\r\b \x00")
			content := strings.Trim(row[1], "\n\t\r\b \x00")
			if err = writeFile(title, content); err != nil {
				t.Errorf("write csv err: %v", err)
			}
		}
	}
}

func TestExtractCsv(t *testing.T) {

	const dir = "samples/articles"

	t.Run("微信", testExtractCsv("samples/weixin.csv", dir))
	t.Run("头条", testExtractCsv("samples/toutiao.csv", dir))
	t.Run("网易号", testExtractCsv("samples/wangyi.csv", dir))
	t.Run("学习强国", testExtractCsv("samples/qiangguo.csv", dir))
	t.Run("百家号", testExtractCsv("samples/baijia.csv", dir))
}

func testNGramSample(ng *ngram.NGram, path string) func(*testing.T) {
	return func(t *testing.T) {
		f, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		r := csv.NewReader(f)
		r.LazyQuotes = true

		_, _ = r.Read()
		for {

			row, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Logf("read csv error: %v", err)
				continue
			}

			ng.Sample(row[0], row[1])
		}
	}
}

func testNGramDump(ng *ngram.NGram, path string) func(*testing.T) {
	return func(t *testing.T) {
		if err := ng.Dump(path, ngram.DefaultDumpThreshold); err != nil {
			t.Fatal(err)
		}
	}
}

func testNGram(n int, outputDir string) func(*testing.T) {
	return func(t *testing.T) {
		ng, err := ngram.New(n,
			"samples/dict/dictionary.txt",
			"samples/dict/stopWords.txt",
		)
		if err != nil {
			t.Fatal(err)
		}

		t.Run("微信", testNGramSample(ng, "samples/weixin.csv"))
		t.Run("头条", testNGramSample(ng, "samples/toutiao.csv"))
		t.Run("网易号", testNGramSample(ng, "samples/wangyi.csv"))
		t.Run("百家号", testNGramSample(ng, "samples/baijia.csv"))
		t.Run("学习强国", testNGramSample(ng, "samples/qiangguo.csv"))
		t.Run("dump", testNGramDump(ng, filepath.Join(outputDir, fmt.Sprintf("dict_%dgram.txt", n))))
	}
}

func TestNGram(t *testing.T) {

	// t.Run("2-Gram", testNGram(2, "samples"))
	// t.Run("3-Gram", testNGram(3, "samples"))
	t.Run("4-Gram", testNGram(4, "samples"))
}

func TestNGramDict(t *testing.T) {

	ng, err := ngram.New(3,
		"samples/dict/dictionary.txt",
		"samples/dict/stopWords.txt",
	)
	if err != nil {
		t.Fatal(err)
	}

	if err = ng.LoadDict(context.TODO(), "samples/dict_4gram.txt"); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile("samples/article.txt")
	if err != nil {
		t.Fatal(err)
	}

	r := ng.Scan(string(data))
	for _, w := range r.TopN(5) {
		fmt.Printf("%s\t%.4f\n", w.Word, w.Score)
	}
}
