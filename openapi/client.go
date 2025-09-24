package openapi

import (
	"bytes"
	"fmt"
	"github.com/jack0829/letsgo/config"
	"github.com/jack0829/letsgo/restful"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	ClientUA = "OpenAPI-SDK-Client/v1.0.0"
)

type Client struct {
	cfg config.OpenAPI
	ats AccessTokenStorager
	cl  *http.Client
}

func NewClient(cfg config.OpenAPI, ops ...ClientOption) *Client {

	cl := &Client{
		cfg: cfg,
		ats: nil,
		cl: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   time.Second * 10,
		},
	}

	for _, fn := range ops {
		fn(cl)
	}

	return cl
}

func (c *Client) getAccessToken() (at *AccessToken, err error) {

	// 加载 token
	at, ok := c.loadAccessToken()
	if ok {
		return
	}

	qs := url.Values{}
	qs.Set("grant_type", "client_credentials")
	qs.Set("client_id", c.cfg.Client.ID)
	qs.Set("client_secret", c.cfg.Client.Secret)

	uri := fmt.Sprintf("%s?%s", OAuthTokenAPI, qs.Encode())
	req, err := c.newRequest(http.MethodGet, uri, nil)
	if err != nil {
		return
	}

	c.debugRequest(os.Stdout, req)
	resp, err := c.cl.Do(req)
	c.debugResponse(os.Stdout, resp, err)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	var r restful.TypedResponse[AccessTokenDTO]
	if err = jsoniter.NewDecoder(resp.Body).Decode(&r); err != nil {
		return
	}

	if r.Code != 200 {
		err = fmt.Errorf("%d | %s", r.Code, r.Msg)
		return
	}

	at = r.Data.FixExpireAt()
	fmt.Fprintf(os.Stdout, "Access Token 过期时间：%s\n", at.ExpiresAt)

	// 保存 token
	if err = c.saveAccessToken(at); err != nil {
		return
	}

	return
}

func (c *Client) saveAccessToken(tk *AccessToken) error {

	if c.ats == nil {
		return nil
	}

	if c.cfg.Debug {
		fmt.Fprintf(os.Stdout, "保存 AccessToken：%s\n", tk.AccessToken)
	}
	return c.ats.Set(tk)
}

func (c *Client) loadAccessToken() (tk *AccessToken, ok bool) {

	if c.ats == nil {
		return
	}

	if tk = c.ats.Get(); tk != nil {
		if tk.AccessToken == "" {
			return
		}
		if exp, _ := tk.Expired(); exp {
			return
		}
		ok = true
	}

	return
}

func (c *Client) newRequest(method, api string, body io.Reader) (req *http.Request, err error) {

	if req, err = http.NewRequest(method, c.cfg.Addr+api, body); err != nil {
		return
	}

	req.Header.Set("User-Agent", ClientUA)
	return
}

func (c *Client) Get(api string, qs url.Values) (*http.Response, error) {

	uri := api
	if qs != nil && len(qs) > 0 {
		uri += "?" + qs.Encode()
	}

	req, err := c.newRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

func (c *Client) PostReader(api string, body io.Reader) (*http.Response, error) {

	req, err := c.newRequest(http.MethodPost, api, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return c.do(req)
}

func (c *Client) Post(api string, body any) (resp *http.Response, err error) {
	var buf bytes.Buffer
	if err = jsoniter.NewEncoder(&buf).Encode(body); err != nil {
		return
	}
	return c.PostReader(api, &buf)
}

func (c *Client) do(req *http.Request) (resp *http.Response, err error) {

	token, err := c.getAccessToken()
	if err != nil {
		return
	}

	k, v := token.Header()
	req.Header.Set(k, v)

	c.debugRequest(os.Stdout, req)
	resp, err = c.cl.Do(req)
	c.debugResponse(os.Stdout, resp, err)

	return
}

func (c *Client) debugRequest(out io.Writer, req *http.Request) {

	if !c.cfg.Debug {
		return
	}

	if req == nil {
		return
	}

	if req.Body == nil {
		req.Write(out)
		return
	}

	rc, err := req.GetBody()
	if err != nil || rc == nil {
		return
	}
	fmt.Fprintf(out, "\n")
	req.Write(out)
	req.Body.Close()
	req.Body = rc
}

func (c *Client) debugResponse(out io.Writer, resp *http.Response, err error) {

	if !c.cfg.Debug {
		return
	}

	if resp == nil {
		return
	}

	if err != nil {
		fmt.Fprintf(out, "Request Error: %v\n", err)
		return
	}

	buf := bytes.NewBuffer(nil)
	buf.ReadFrom(resp.Body)
	resp.Body.Close()
	br, bw := io.Pipe()
	resp.Body = br
	mw := io.MultiWriter(bw, out)
	go func() {
		defer bw.Close()
		fmt.Fprintf(out, "\n%s %s\n", resp.Proto, resp.Status)
		resp.Header.Write(out)
		fmt.Fprintf(out, "\n")
		buf.WriteTo(mw)
		fmt.Fprintf(out, "\n")
	}()

}
