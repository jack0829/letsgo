package wechat

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"time"
)

type AccessToken struct {
	AppID       string    `json:"app_id,omitempty"`
	AccessToken string    `json:"access_token"`
	ExpireAt    time.Time `json:"expire_at"`
}

func (tk *AccessToken) Expired() bool {
	if tk.ExpireAt.IsZero() {
		return true
	}
	return time.Until(tk.ExpireAt) < time.Minute // 留1分钟冗余无缝更换
}

// GetAccessToken 获取应用 AccessToken
func (w *Wechat) GetAccessToken() (tk *AccessToken, err error) {

	if s := w.storage.accessToken; s != nil {

		if tk = s.GetAccessToken(w.AppID); tk != nil && !tk.Expired() {
			return
		}

		defer func() {
			if err == nil {
				err = s.SetAccessToken(tk)
			}
		}()
	}

	if w.ops.stableAccessToken {
		return w.getStableAccessToken()
	}

	return w.getAccessToken()
}

func (w *Wechat) getAccessToken() (*AccessToken, error) {

	u := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=%s&appid=%s&secret=%s",
		"client_credential",
		w.AppID,
		w.AppSecret,
	)

	now := time.Now()
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
		AccessToken string `json:"access_token,omitempty"`
		ExpiresIn   int    `json:"expires_in,omitempty"`
	}
	if err = jsoniter.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	if err = r.Error(); err != nil {
		return nil, err
	}

	return &AccessToken{
		AppID:       w.AppID,
		AccessToken: r.AccessToken,
		ExpireAt:    now.Add(time.Second * time.Duration(r.ExpiresIn)),
	}, nil
}

type stableAccessTokenRequest struct {
	GrantType    string `json:"grant_type"`
	AppID        string `json:"appid"`
	Secret       string `json:"secret"`
	ForceRefresh bool   `json:"force_refresh"`
}

type getStableAccessTokenOption func(r *stableAccessTokenRequest)

func ForceRefresh(r *stableAccessTokenRequest) {
	r.ForceRefresh = true
}

func (w *Wechat) getStableAccessToken(ops ...getStableAccessTokenOption) (tk *AccessToken, err error) {

	req := &stableAccessTokenRequest{
		GrantType: "client_credential",
		AppID:     w.AppID,
		Secret:    w.AppSecret,
	}

	for _, op := range ops {
		op(req)
	}

	body := bytes.NewBuffer(nil)
	if err = jsoniter.NewEncoder(body).Encode(req); err != nil {
		return
	}

	now := time.Now()
	resp, err := w.c.Post(
		"https://api.weixin.qq.com/cgi-bin/stable_token",
		"application/json",
		body,
	)
	if err != nil {
		return nil, err
	}
	defer respClose(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}

	var r struct {
		response
		AccessToken string `json:"access_token,omitempty"`
		ExpiresIn   int    `json:"expires_in,omitempty"`
	}
	if err = jsoniter.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	if err = r.Error(); err != nil {
		return nil, err
	}

	return &AccessToken{
		AppID:       w.AppID,
		AccessToken: r.AccessToken,
		ExpireAt:    now.Add(time.Second * time.Duration(r.ExpiresIn)),
	}, nil
	return
}
