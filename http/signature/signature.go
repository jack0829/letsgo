package signature

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
)

const (
	DefaultHeaderName    = "Signature"
	DefaultClockName     = "Clock"
	DefaultOriginUrlName = "OriginURL"
)

var (
	请求过期 = fmt.Errorf("请求已过期")
	签名错误 = fmt.Errorf("签名错误")
)

type Signature struct {
	secret        string // 密钥
	headerName    string // HTTP 头部签名字段名
	clockName     string // HTTP 头部时钟字段名
	originUrlName string // HTTP 头部原始 URL 字段名
}

func New(secret string, ops ...Option) *Signature {

	o := &Signature{
		secret:        secret,
		headerName:    DefaultHeaderName,
		clockName:     DefaultClockName,
		originUrlName: DefaultOriginUrlName,
	}

	for _, op := range ops {
		op(o)
	}

	return o
}

// Sign 给请求添加签名
func (s *Signature) Sign(req *http.Request) {

	uri := req.URL.String()
	req.Header.Set(s.originUrlName, uri)

	clock := fmt.Sprintf("%d", time.Now().Unix())
	req.Header.Set(s.clockName, clock)

	sign, body := s.sum(uri, clock, req.Body)

	req.Header.Set(s.headerName, sign)
	req.Body = body
}

// Check 验证请求
func (s *Signature) Check(req *http.Request) error {

	sign := req.Header.Get(s.headerName)
	uri := req.Header.Get(s.originUrlName)
	clock := req.Header.Get(s.clockName)

	if clock == "" {
		return 请求过期
	}

	ts, err := strconv.ParseInt(clock, 10, 64)
	if err != nil {
		return 请求过期
	}

	if d := float64(time.Now().Unix() - ts); math.Abs(d) > 10 {
		return 请求过期
	}

	if uri == "" {
		uri = req.URL.String()
	}

	var sum string
	if sum, req.Body = s.sum(uri, clock, req.Body); sum != sign {
		return 签名错误
	}

	return nil
}

func (s *Signature) sum(uri, clock string, body io.ReadCloser) (string, io.ReadCloser) {

	h := md5.New()
	io.WriteString(h, s.secret+"\n")
	io.WriteString(h, uri+"\n")
	io.WriteString(h, clock+"\n")

	var newBody io.ReadCloser
	if body != nil {
		buf := bytes.NewBuffer(nil)
		io.Copy(io.MultiWriter(h, buf), body)
		body.Close()
		newBody = io.NopCloser(buf)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), newBody
}
