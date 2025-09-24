package wechat

import (
	"crypto/sha1"
	"encoding/hex"
	"golang.org/x/exp/slices"
	"io"
)

func (e *Event) Hello(
	signature,
	timestamp,
	nonce string,
) bool {

	s := []string{
		e.Token,
		timestamp,
		nonce,
	}

	slices.Sort(s)

	h := sha1.New()
	for _, v := range s {
		io.WriteString(h, v)
	}

	return hex.EncodeToString(h.Sum(nil)) == signature
}

func (w *Wechat) Event() *Event {
	return w.e
}
