package http

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/jack0829/letsgo/http/signature"
	"io"
	"net/http"
)

type Transport struct {
	http.Transport
	debugger   io.Writer
	setHeaders map[string]string
	signature  *signature.Signature
}

func (t *Transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {

	if t.setHeaders != nil {
		for k, v := range t.setHeaders {
			req.Header.Set(k, v)
		}
	}

	t.Sign(req)

	if w := t.debugger; w != nil {

		// debug request
		fmt.Fprintln(w, req.Method, req.URL.RequestURI())
		if req.GetBody != nil {
			if body, err := req.GetBody(); err == nil {
				bufio.NewReader(body).WriteTo(w)
				fmt.Fprintln(w)
				body.Close()
			}
		}

		// debug response
		defer func() {

			if err != nil {
				fmt.Fprintln(w, err.Error())
				return
			}

			fmt.Fprintln(w, resp.Proto, resp.Status)
			if body := resp.Body; body != nil {
				buf := bytes.NewBuffer(nil)
				bufio.NewReader(body).WriteTo(io.MultiWriter(buf, w))
				fmt.Fprintln(w)
				body.Close()
				resp.Body = io.NopCloser(buf)
			}
		}()
	}

	return t.Transport.RoundTrip(req)
}

func (t *Transport) Debug(w io.Writer) *Transport {
	t.debugger = w
	return t
}

func (t *Transport) SetHeader(k, v string) {

	if t.setHeaders == nil {
		t.setHeaders = make(map[string]string)
	}

	t.setHeaders[k] = v
}

func (t *Transport) SetSignature(secret string) {
	t.signature = signature.New(secret)
}

func (t *Transport) Sign(req *http.Request) {
	if s := t.signature; s != nil {
		s.Sign(req)
	}
}
