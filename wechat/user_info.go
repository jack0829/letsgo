package wechat

import (
	jsoniter "github.com/json-iterator/go"
	"net/url"
	"time"
)

const (
	SubscribeSceneSearch      = "ADD_SCENE_SEARCH"               // 公众号搜索
	SubscribeSceneMigration   = "ADD_SCENE_ACCOUNT_MIGRATION"    // 公众号迁移
	SubscribeSceneProfileCard = "ADD_SCENE_PROFILE_CARD"         // 名片分享
	SubscribeSceneQRCode      = "ADD_SCENE_QR_CODE"              // 扫描二维码
	SubscribeSceneProfileLink = "ADD_SCENE_PROFILE_LINK"         // 图文页内名称点击
	SubscribeSceneProfileItem = "ADD_SCENE_PROFILE_ITEM"         // 图文页右上角菜单
	SubscribeScenePaid        = "ADD_SCENE_PAID"                 // 支付后关注
	SubscribeSceneWechatAD    = "ADD_SCENE_WECHAT_ADVERTISEMENT" // 微信广告
	SubscribeSceneReprint     = "ADD_SCENE_REPRINT"              // 他人转载
	SubscribeSceneLiveStream  = "ADD_SCENE_LIVESTREAM"           // 视频号直播
	SubscribeSceneChannel     = "ADD_SCENE_CHANNELS"             // 视频号
	SubscribeSceneMPAttention = "ADD_SCENE_WXA"                  // 小程序关注
	SubscribeSceneOther       = "ADD_SCENE_OTHERS"               // 其他
)

type UserInfo struct {
	Subscribe      bool      `json:"subscribe"`                 // 用户是否订阅该公众号，用户没有关注该公众号时拉取不到其余信息。
	OpenID         string    `json:"openid"`                    //  OpenID
	Language       string    `json:"language"`                  // 用户的语言，简体中文为zh_CN
	SubscribeTime  time.Time `json:"subscribe_time"`            // 用户关注时间。如果用户曾多次关注，则取最后关注时间
	UnionID        string    `json:"unionid"`                   // UnionID
	Remark         string    `json:"remark,omitempty"`          // 公众号运营者对粉丝的备注，公众号运营者可在微信公众平台用户管理界面对粉丝添加备注
	GroupID        int64     `json:"groupid,omitempty"`         // 用户所在的分组ID（兼容旧的用户分组接口）
	TagIDList      []int64   `json:"tagid_list,omitempty"`      // 用户被打上的标签ID列表
	SubscribeScene string    `json:"subscribe_scene,omitempty"` // 返回用户关注的渠道来源
	QrScene        int64     `json:"qr_scene,omitempty"`        // 二维码扫码场景（开发者自定义）
	QrSceneStr     string    `json:"qr_scene_str,omitempty"`    // 二维码扫码场景描述（开发者自定义）
}

// GetUserInfo 拉取用户信息
func (w *Wechat) GetUserInfo(openID string) (*UserInfo, error) {

	tk, err := w.GetAccessToken()
	if err != nil {
		return nil, err
	}

	qs := url.Values{}
	qs.Set("access_token", tk.AccessToken)
	qs.Set("openid", openID)
	qs.Set("lang", "zh_CN")

	resp, err := w.c.Get("https://api.weixin.qq.com/cgi-bin/user/info?" + qs.Encode())
	if err != nil {
		return nil, err
	}
	defer respClose(resp)

	var r struct {
		response
		Subscribe      int     `json:"subscribe"`       // 用户是否订阅该公众号标识，值为0时，代表此用户没有关注该公众号，拉取不到其余信息。
		OpenID         string  `json:"openid"`          // 用户的标识，对当前公众号唯一
		Language       string  `json:"language"`        // 用户的语言，简体中文为zh_CN
		SubscribeTime  int64   `json:"subscribe_time"`  // 用户关注时间，为时间戳。如果用户曾多次关注，则取最后关注时间
		UnionID        string  `json:"unionid"`         // 只有在用户将公众号绑定到微信开放平台账号后，才会出现该字段。
		Remark         string  `json:"remark"`          // 公众号运营者对粉丝的备注，公众号运营者可在微信公众平台用户管理界面对粉丝添加备注
		GroupID        int64   `json:"groupid"`         // 用户所在的分组ID（兼容旧的用户分组接口）
		TagIDList      []int64 `json:"tagid_list"`      // 用户被打上的标签ID列表
		SubscribeScene string  `json:"subscribe_scene"` // 返回用户关注的渠道来源
		QrScene        int64   `json:"qr_scene"`        // 二维码扫码场景（开发者自定义）
		QrSceneStr     string  `json:"qr_scene_str"`    // 二维码扫码场景描述（开发者自定义）
	}

	if err = jsoniter.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &UserInfo{
		Subscribe:      r.Subscribe != 0,
		OpenID:         r.OpenID,
		Language:       r.Language,
		SubscribeTime:  time.Unix(r.SubscribeTime, 0),
		UnionID:        r.UnionID,
		Remark:         r.Remark,
		GroupID:        r.GroupID,
		TagIDList:      r.TagIDList,
		SubscribeScene: r.SubscribeScene,
		QrScene:        r.QrScene,
		QrSceneStr:     r.QrSceneStr,
	}, nil
}
