package ngram

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
)

func ReadDict(
	ctx context.Context,
	path string,
) (<-chan *Meta, int, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}

	var totalDocs int
	if _, err = fmt.Fscanf(f, "%d\n", &totalDocs); err != nil {
		return nil, 0, err
	}

	ch := make(chan *Meta, 2)
	go func() {

		defer close(ch)
		defer f.Close()

		for {

			var (
				word                 string
				freq, docs, tokenLen int
			)

			if _, err = fmt.Fscanf(f, "%s\t%d\t%d\t%d\n", &word, &freq, &docs, &tokenLen); err != nil {
				if err != io.EOF {
					log.Printf("read ngram dict %s err: %v\n", path, err)
				}
				break
			}

			select {
			case <-ctx.Done():
				return
			case ch <- &Meta{
				Word: word,
				Freq: freq,
				Docs: docs,
			}:
				// fmt.Println(word, freq, docs)
			}
		}
	}()

	return ch, totalDocs, nil
}
