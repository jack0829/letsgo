package wechat

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/exp/maps"
	"net/url"
)

type TplMessage struct {
	TemplateID  string                     `json:"template_id"`             // 所需下发的订阅模板id
	ToUser      string                     `json:"touser"`                  // 接收者（用户）的 openid
	Data        map[string]*tplMessageData `json:"data"`                    // 模板内容
	URL         string                     `json:"url,omitempty"`           // 模板跳转链接（海外账号没有跳转能力）
	MiniProgram *miniProgram               `json:"miniprogram,omitempty"`   // 跳小程序所需数据，不需跳小程序可不用传该数据
	ClientMsgID string                     `json:"client_msg_id,omitempty"` // 防重入id。对于同一个openid + client_msg_id, 只发送一条消息,10分钟有效,超过10分钟不保证效果
}

type tplMessageData struct {
	Value any `json:"value"`
}

type miniProgram struct {
	AppID    string `json:"appid"`    // 所需跳转到的小程序appid（该小程序appid必须与发模板消息的公众号是绑定关联关系，暂不支持小游戏）
	PagePath string `json:"pagepath"` // 所需跳转到小程序的具体页面路径，支持带参数,（示例index?foo=bar），要求该小程序已发布，暂不支持小游戏
}

func NewTplMessage(
	id,
	to string,
) *TplMessage {
	return &TplMessage{
		TemplateID: id,
		ToUser:     to,
		Data:       make(map[string]*tplMessageData),
	}
}

// SetField 设置模板字段
func (tm *TplMessage) SetField(k string, v any) *TplMessage {
	if f, ok := tm.Data[k]; ok {
		f.Value = v
	} else {
		tm.Data[k] = &tplMessageData{
			Value: v,
		}
	}
	return tm
}

// ClearField 删除模板字段，不传为全部
func (tm *TplMessage) ClearField(key ...string) *TplMessage {
	if len(key) < 1 {
		key = maps.Keys(tm.Data)
	}
	for _, k := range key {
		delete(tm.Data, k)
	}
	return tm
}

// SetMiniProgramPage 指定跳小程序的页面，比 URL 字段优先，不跳小程序不传，不支持小游戏，小程序必须已发布，且与发模板消息的公众号是绑定关联关系
func (tm *TplMessage) SetMiniProgramPage(
	appID, // 小程序 AppID
	page string, // 页面路径（示例index?foo=bar）
) *TplMessage {
	if tm.MiniProgram == nil {
		tm.MiniProgram = &miniProgram{}
	}
	tm.MiniProgram.AppID = appID
	tm.MiniProgram.PagePath = page
	return tm
}

// ClearMiniProgramPage 取消跳转小程序
func (tm *TplMessage) ClearMiniProgramPage() *TplMessage {
	tm.MiniProgram = nil
	return tm
}

// SetURL 指定跳转 URL 地址（https://）
func (tm *TplMessage) SetURL(v string) *TplMessage {
	tm.URL = v
	return tm
}

func (tm *TplMessage) WithID(v string) *TplMessage {
	tm.ClientMsgID = v
	return tm
}

// SendTplMessage 发送模板消息（必须是公众号）
func (w *Wechat) SendTplMessage(mt *TplMessage) (msgID int64, err error) {

	tk, err := w.GetAccessToken()
	if err != nil {
		return
	}

	qs := url.Values{}
	qs.Set("access_token", tk.AccessToken)

	body := bytes.NewBuffer(nil)
	if err = jsoniter.NewEncoder(body).Encode(mt); err != nil {
		return
	}

	resp, err := w.c.Post(
		"https://api.weixin.qq.com/cgi-bin/message/template/send?"+qs.Encode(),
		"application/json",
		body,
	)
	if err != nil {
		return
	}
	defer respClose(resp)

	var r struct {
		response
		MsgID int64 `json:"msgid"`
	}
	if err = jsoniter.NewDecoder(resp.Body).Decode(&r); err != nil {
		return
	}

	if err = r.Error(); err != nil {
		return
	}

	msgID = r.MsgID
	return
}
