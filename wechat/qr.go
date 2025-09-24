package wechat

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/url"
	"time"
)

const (
	QRActionNameScene         qrActionName = "QR_SCENE"           // 带场景ID的临时二维码
	QRActionNameStrScene      qrActionName = "QR_STR_SCENE"       // 带场景str的临时二维码
	QRActionNameLimitScene    qrActionName = "QR_LIMIT_SCENE"     // 带场景ID的永久二维码
	QRActionNameLimitStrScene qrActionName = "QR_LIMIT_STR_SCENE" // 带场景str的永久二维码
)

type qrActionName string

type QRTicket struct {
	Ticket   string    `json:"ticket"`    // 获取的二维码ticket，凭借此ticket可以在有效时间内换取二维码。
	ExpireAt time.Time `json:"expire_at"` // 二维码有效时间
	URL      string    `json:"url"`       // 二维码图片解析后的地址，开发者可根据该地址自行生成需要的二维码图片
}

type createQRRequest struct {
	ActionName qrActionName `json:"action_name"`
	ActionInfo struct {
		Scene struct {
			SceneID  int64  `json:"scene_id,omitempty"`  // 1-100000
			SceneStr string `json:"scene_str,omitempty"` // 长度 1-64
		} `json:"scene"`
	} `json:"action_info"`
	ExpireSeconds int64 `json:"expire_seconds,omitempty"` // 最大 30 天，即 86400*30
}

func (r *createQRRequest) buffer() *bytes.Buffer {
	b := bytes.NewBuffer(nil)
	jsoniter.NewEncoder(b).Encode(r)
	return b
}

type CreateQROption func(r *createQRRequest)

// TemporaryQR 临时二维码（最长 30 天）
func TemporaryQR(scene string, exp time.Duration) CreateQROption {
	return func(r *createQRRequest) {
		if exp > time.Hour*24*30 {
			exp = time.Hour * 24 * 30
		}
		r.ActionName = QRActionNameStrScene
		r.ActionInfo.Scene.SceneStr = scene
		r.ExpireSeconds = int64(exp.Seconds())
	}
}

// PermanentQR 永久二维码（上限 100000 个）
func PermanentQR(scene string) CreateQROption {
	return func(r *createQRRequest) {
		r.ActionName = QRActionNameLimitStrScene
		r.ActionInfo.Scene.SceneStr = scene
	}
}

// CreateQRTicket 生成场景二维码 ticket
func (w *Wechat) CreateQRTicket(opt CreateQROption) (tick *QRTicket, err error) {

	tk, err := w.GetAccessToken()
	if err != nil {
		return
	}

	var req createQRRequest
	opt(&req)

	u := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%s", url.QueryEscape(tk.AccessToken))
	resp, err := w.c.Post(u, "application/json", req.buffer())
	if err != nil {
		return
	}
	defer respClose(resp)

	var r struct {
		response
		Ticket        string `json:"ticket"`         // 获取的二维码ticket，凭借此ticket可以在有效时间内换取二维码。
		ExpireSeconds int64  `json:"expire_seconds"` // 该二维码有效时间，以秒为单位。 最大不超过2592000（即30天）。
		URL           string `json:"url"`            // 二维码图片解析后的地址，开发者可根据该地址自行生成需要的二维码图片
	}

	if err = jsoniter.NewDecoder(resp.Body).Decode(&r); err != nil {
		return
	}

	return &QRTicket{
		Ticket:   r.Ticket,
		ExpireAt: time.Now().Add(time.Second * time.Duration(r.ExpireSeconds)),
		URL:      r.URL,
	}, nil
}
