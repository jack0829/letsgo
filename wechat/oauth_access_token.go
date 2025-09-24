package wechat

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"net/url"
	"time"
)

type OAuthAccessToken struct {
	AppID          string    `json:"app_id"`
	AccessToken    string    `json:"access_token"`
	ExpireAt       time.Time `json:"expire_at"`
	RefreshToken   string    `json:"refresh_token"`
	OpenID         string    `json:"open_id"`
	Scope          string    `json:"scope"`
	IsSnapshotUser bool      `json:"is_snapshot_user,omitempty"`
	UnionID        string    `json:"union_id"`
}

func (tk *OAuthAccessToken) Expired() bool {
	if tk.ExpireAt.IsZero() {
		return true
	}
	return time.Until(tk.ExpireAt) < time.Minute // 留1分钟冗余无缝更换
}

// GetOAuthAccessToken 获取 OAuth AccessToken
func (w *Wechat) GetOAuthAccessToken(openID string) (tk *OAuthAccessToken, err error) {

	if s := w.storage.oauthAccessToken; s != nil {

		if tk = s.GetOAuthAccessToken(w.AppID, openID); tk == nil || !tk.Expired() {
			return
		}

	}

	return w.RefreshOAuthAccessToken(tk)
}

// RefreshOAuthAccessToken 刷新 OAuth AccessToken
func (w *Wechat) RefreshOAuthAccessToken(tk *OAuthAccessToken) (*OAuthAccessToken, error) {

	if tk == nil {
		return nil, nil
	}

	qs := url.Values{}
	qs.Set("appid", w.AppID)
	qs.Set("refresh_token", tk.RefreshToken)
	qs.Set("grant_type", "refresh_token")

	u := "https://api.weixin.qq.com/sns/oauth2/refresh_token?" + qs.Encode()

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
		AccessToken  string `json:"access_token"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		OpenID       string `json:"openid"`
		Scope        string `json:"scope"`
	}

	if err = jsoniter.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	if err = r.Error(); err != nil {
		return nil, err
	}

	tk.AccessToken = r.AccessToken
	tk.ExpireAt = time.Now().Add(time.Duration(r.ExpiresIn))
	tk.RefreshToken = r.RefreshToken
	// tk.OpenID = r.OpenID
	tk.Scope = r.Scope

	err = w.storage.oauthAccessToken.SetOAuthAccessToken(tk)
	return tk, err
}

func (w *Wechat) Code2OAuthAccessToken(code string) (*OAuthAccessToken, error) {

	qs := url.Values{}
	qs.Set("appid", w.AppID)
	qs.Set("secret", w.AppSecret)
	qs.Set("code", code)
	qs.Set("grant_type", "authorization_code")

	u := "https://api.weixin.qq.com/sns/oauth2/access_token?" + qs.Encode()

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
		AccessToken    string `json:"access_token"`
		ExpiresIn      int64  `json:"expires_in"`
		RefreshToken   string `json:"refresh_token"`
		OpenID         string `json:"openid"`
		Scope          string `json:"scope"`
		IsSnapshotUser int    `json:"is_snapshotuser"`
		UnionID        string `json:"unionid"`
	}

	if err = jsoniter.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	if err = r.Error(); err != nil {
		return nil, err
	}

	tk := &OAuthAccessToken{
		AppID:          w.AppID,
		AccessToken:    r.AccessToken,
		ExpireAt:       time.Now().Add(time.Duration(r.ExpiresIn)),
		RefreshToken:   r.RefreshToken,
		OpenID:         r.OpenID,
		Scope:          r.Scope,
		IsSnapshotUser: r.IsSnapshotUser > 0,
		UnionID:        r.UnionID,
	}

	err = w.storage.oauthAccessToken.SetOAuthAccessToken(tk)
	return tk, err
}
