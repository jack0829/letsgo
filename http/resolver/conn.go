package resolver

import (
	"context"
	"net"
	"strings"
)

type Host struct {
	Network string
	IP      net.IP
	Path    string
}

func (h *Host) String() string {
	if h.Path != "" {
		return h.Network + ":" + h.Path
	}
	return h.Network + ":" + h.IP.String()
}

func (h *Host) Conn(
	ctx context.Context,
	port string,
) (net.Conn, error) {
	var d net.Dialer
	if strings.ToLower(h.Network) == "unix" {
		return d.DialContext(ctx, "unix", h.Path)
	}
	return d.DialContext(ctx, h.Network, h.IP.String()+":"+strings.Trim(port, ":"))
}
