package wechat

import "fmt"

type Phone struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
	Watermark       struct {
		Timestamp int64  `json:"timestamp"`
		AppID     string `json:"appid"`
	} `json:"watermark"`
}

func (p Phone) Check(appID string) error {
	if p.Watermark.AppID != appID {
		return fmt.Errorf("AppID 错误")
	}
	return nil
}
