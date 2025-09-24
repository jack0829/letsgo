package resolver

import "net"

type IP string

func (v IP) Value() (ip net.IP, err error) {
	err = ip.UnmarshalText([]byte(v))
	return
}

func (v IP) String() string {
	if ip, err := v.Value(); err == nil {
		return ip.String()
	}
	return string(v)
}
