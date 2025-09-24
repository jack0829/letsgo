package wechat

import "net/http"

type Wechat struct {
	AppID     string
	AppSecret string
	e         *Event
	c         *http.Client
	storage   storage
	ops       struct {
		stableAccessToken bool
	}
}

type Option func(w *Wechat)

func New(
	id, secret string,
	ops ...Option,
) *Wechat {

	w := &Wechat{
		AppID:     id,
		AppSecret: secret,
		c:         http.DefaultClient,
	}

	for _, op := range ops {
		op(w)
	}

	return w
}

func WithClient(c *http.Client) Option {
	return func(w *Wechat) {
		if c != nil {
			w.c = c
		}
	}
}

func WithStorage(s Storage) Option {
	return func(w *Wechat) {
		w.storage.accessToken = s
		w.storage.jsApiTicket = s
		w.storage.session = s
		w.storage.oauthAccessToken = s
	}
}

func WithAccessTokenStorage(s AccessTokenStorage) Option {
	return func(w *Wechat) {
		w.storage.accessToken = s
	}
}

func WithJsApiTicketStorage(s JsApiTicketStorage) Option {
	return func(w *Wechat) {
		w.storage.jsApiTicket = s
	}
}

func WithSessionStorage(s SessionStorage) Option {
	return func(w *Wechat) {
		w.storage.session = s
	}
}

func WithOAuthAccessTokenStorage(s OAuthAccessTokenStorage) Option {
	return func(w *Wechat) {
		w.storage.oauthAccessToken = s
	}
}

func UseStableAccessToken(w *Wechat) {
	w.ops.stableAccessToken = true
}

func respClose(resp *http.Response) {
	if resp == nil {
		return
	}
	if b := resp.Body; b != nil {
		b.Close()
	}
}
