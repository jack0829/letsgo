package wechat

import (
	"crypto/sha1"
	"fmt"
	"github.com/google/uuid"
	"io"
	"time"
)

type JsApiWxConfig struct {
	AppID     string `json:"app_id"`
	Timestamp int64  `json:"timestamp"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
}

// GetJsApiWxConfig 获取 JsApi 接口授权签名
func (w *Wechat) GetJsApiWxConfig(url string) (*JsApiWxConfig, error) {

	ticket, err := w.GetJsApiTicket()
	if err != nil {
		return nil, err
	}

	wc := &JsApiWxConfig{
		AppID:     w.AppID,
		Timestamp: time.Now().Unix(),
		Nonce:     uuid.New().String(),
	}

	h := sha1.New()
	io.WriteString(h, fmt.Sprintf(
		"jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		ticket.Ticket,
		wc.Nonce,
		wc.Timestamp,
		url,
	))
	wc.Signature = fmt.Sprintf("%x", h.Sum(nil))

	return wc, nil
}
