package openapi

import (
	"net/http"
	"time"
)

type ClientOption func(c *Client)

func ClientTimeOut(d time.Duration) ClientOption {
	return func(c *Client) {
		c.cl.Timeout = d
	}
}

func WithTransport(tr http.RoundTripper) ClientOption {
	return func(c *Client) {
		c.cl.Transport = tr
	}
}

func WithAccessTokenStorager(ats AccessTokenStorager) ClientOption {
	return func(c *Client) {
		c.ats = ats
	}
}
