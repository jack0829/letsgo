package openapi

import (
	"fmt"
	"time"
)

type AccessToken struct {
	Scope       string    `json:"scope"`
	Type        string    `json:"type"`
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
}

func (t *AccessToken) Header() (key, value string) {
	return "Authorization", fmt.Sprintf("%s %s", t.Type, t.AccessToken)
}

func (t *AccessToken) Expired() (b bool, err error) {

	if t.ExpiresAt.IsZero() {
		err = fmt.Errorf("expires_at 未同步修正")
		return
	}

	return t.ExpiresAt.Before(time.Now()), nil
}

type AccessTokenDTO struct {
	AccessToken
	ExpiresIn int64 `json:"expires_in"` // 秒
}

func (t *AccessTokenDTO) FixExpireAt() *AccessToken {
	t.ExpiresAt = time.Now().Add(time.Second * time.Duration(t.ExpiresIn-10)) // 提前 10 秒逻辑过期
	return &t.AccessToken
}
