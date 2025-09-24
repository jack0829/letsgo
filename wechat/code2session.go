package wechat

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Session struct {
	OpenID  string `json:"open_id"`
	AppID   string `json:"app_id"`
	Key     string `json:"key"`
	UnionID string `json:"union_id,omitempty"`
}

func (s *Session) Signature() string {
	h := hmac.New(sha256.New, []byte(s.Key))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Code2Session wx.login code 换取 微信登录态 session_key
func (w *Wechat) Code2Session(code string) (session *Session, err error) {

	if s := w.storage.session; s != nil {

		defer func() {
			if err == nil {
				err = s.SetSession(session)
			}
		}()
	}

	return w.code2Session(code)
}

func (w *Wechat) code2Session(code string) (session *Session, err error) {

	qs := url.Values{}
	qs.Set("appid", w.AppID)
	qs.Set("secret", w.AppSecret)
	qs.Set("js_code", code)
	qs.Set("grant_type", "authorization_code")

	u := "https://api.weixin.qq.com/sns/jscode2session?" + qs.Encode()

	resp, err := w.c.Get(u)
	if err != nil {
		return nil, err
	}
	defer respClose(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}

	var r struct {
		response
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionid,omitempty"`
	}

	tr := io.TeeReader(resp.Body, log.Writer())
	if err = jsoniter.NewDecoder(tr).Decode(&r); err != nil {
		log.Println("jscode2session decode err: ", err)
		return nil, err
	}

	if err = r.Error(); err != nil {
		return nil, err
	}

	return &Session{
		OpenID:  r.OpenID,
		AppID:   w.AppID,
		Key:     r.SessionKey,
		UnionID: r.UnionID,
	}, nil
}

// GetSession 获取用户微信登录态
func (w *Wechat) GetSession(openID string) *Session {
	if s := w.storage.session; s != nil {
		return s.GetSession(w.AppID, openID)
	}
	return nil
}
