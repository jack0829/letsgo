package wechat

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/url"
)

// CheckSession 验证 微信登录态
func (w *Wechat) CheckSession(openID string) (session *Session, err error) {

	session = w.GetSession(openID)
	if session == nil {
		err = fmt.Errorf("no session")
		return
	}

	tk, err := w.GetAccessToken()
	if err != nil {
		return
	}

	qs := url.Values{}
	qs.Set("access_token", tk.AccessToken)
	qs.Set("signature", session.Signature())
	qs.Set("openid", openID)
	qs.Set("sig_method", "hmac_sha256")

	u := "https://api.weixin.qq.com/wxa/checksession?" + qs.Encode()
	resp, err := w.c.Get(u)
	if err != nil {
		return
	}
	defer respClose(resp)

	var r response
	if err = jsoniter.NewDecoder(resp.Body).Decode(&r); err != nil {
		return
	}

	if err = r.Error(); err != nil {
		if s := w.storage.session; s != nil {
			s.DeleteSession(w.AppID, openID)
		}
	}

	return
}
