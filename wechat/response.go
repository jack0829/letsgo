package wechat

import "fmt"

type response struct {
	Errcode int    `json:"errcode,omitempty"`
	Errmsg  string `json:"errmsg,omitempty"`
}

func (r *response) Error() error {
	if r.Errcode == 0 {
		return nil
	}
	return fmt.Errorf("%d | %s", r.Errcode, r.Errmsg)
}
