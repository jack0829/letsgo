package wechat

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
)

type GetWxaCodeUnLimitOption func(p map[string]any)

// WxaCodeUnLimitPage 默认是主页，页面 page，例如 pages/index/index，根路径前不要填加 /，不能携带参数（参数请放在scene字段里），如果不填写这个字段，默认跳主页面
func WxaCodeUnLimitPage(v string) GetWxaCodeUnLimitOption {
	return func(p map[string]any) {
		p["page"] = v
	}
}

// WxaCodeUnLimitEnvVersion 正式版为 "release"，体验版为 "trial"，开发版为 "develop"。默认是正式版
func WxaCodeUnLimitEnvVersion(v string) GetWxaCodeUnLimitOption {
	return func(p map[string]any) {
		switch v {
		case "trial", "develop", "release":
			p["check_path"] = false
		default:
			return
		}
		p["env_version"] = v
	}
}

// WxaCodeUnLimitSize 默认430，二维码的宽度，单位 px，最小 280px，最大 1280px
func WxaCodeUnLimitSize(v int) GetWxaCodeUnLimitOption {
	return func(p map[string]any) {
		if v >= 280 && v <= 1280 {
			p["width"] = v
		}
	}
}

// GetWxaCodeUnLimit 获取不限制的小程序码
func (w *Wechat) GetWxaCodeUnLimit(
	scene string, // 最大32个可见字符，只支持数字，大小写英文以及部分特殊字符：!#$&'()*+,/:;=?@-._~，其它字符请自行编码为合法字符
	output io.Writer,
	ops ...GetWxaCodeUnLimitOption,
) error {

	at, err := w.GetAccessToken()
	if err != nil {
		return err
	}

	u := fmt.Sprintf(
		"https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s",
		at.AccessToken,
	)

	buf := bytes.NewBuffer(nil)

	p := map[string]any{
		"scene": scene,
	}

	for _, op := range ops {
		op(p)
	}

	if err = jsoniter.NewEncoder(buf).Encode(p); err != nil {
		return err
	}

	resp, err := w.c.Post(u, "application/json", buf)
	if err != nil {
		return err
	}
	defer respClose(resp)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(resp.Status)
	}

	if _, err = io.Copy(output, resp.Body); err != nil {
		return err
	}

	return nil
}
