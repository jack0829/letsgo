package wechat

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"time"
)

type JsApiTicket struct {
	AppID    string    `json:"app_id,omitempty"`
	Ticket   string    `json:"ticket"`
	ExpireAt time.Time `json:"expire_at"`
}

func (tk *JsApiTicket) Expired() bool {
	if tk.ExpireAt.IsZero() {
		return true
	}
	return time.Until(tk.ExpireAt) < time.Minute // 留1分钟冗余无缝更换
}

// GetJsApiTicket 获取 JsApi 调用凭证
func (w *Wechat) GetJsApiTicket() (tk *JsApiTicket, err error) {

	if s := w.storage.jsApiTicket; s != nil {

		if tk = s.GetJsApiTicket(w.AppID); tk != nil && !tk.Expired() {
			return
		}

		defer func() {
			if err == nil {
				err = s.SetJsApiTicket(tk)
			}
		}()
	}

	return w.getJsApiTicket()
}

func (w *Wechat) getJsApiTicket() (*JsApiTicket, error) {

	at, err := w.GetAccessToken()
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=%s",
		at.AccessToken,
		"jsapi",
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

	// {
	//  "errcode": 0,
	//  "errmsg": "ok",
	//  "ticket": "LIKLckvwlJT9cWIhEQTwfN0Fi5UDTH6qpjEe28eQV2A51VjDIPx1defI3qqtajDNQaVvPH8-fjs3FBblZ9PNsw",
	//  "expires_in": 7200
	// }

	var r struct {
		response
		Ticket    string `json:"ticket,omitempty"`
		ExpiresIn int    `json:"expires_in,omitempty"`
	}

	if err = jsoniter.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	if err = r.Error(); err != nil {
		return nil, err
	}

	return &JsApiTicket{
		AppID:    w.AppID,
		Ticket:   r.Ticket,
		ExpireAt: now.Add(time.Second * time.Duration(r.ExpiresIn)),
	}, nil
}
