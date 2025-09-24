package resolver

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

type Resolver struct {
	mutex   sync.RWMutex
	lookups map[string][]*Host // key: 域名（不含端口）
}

func New(
	ops ...Option,
) *Resolver {

	r := &Resolver{
		lookups: make(map[string][]*Host),
	}

	for _, op := range ops {
		op(r)
	}

	return r
}

func (r *Resolver) PrintLookups(w io.Writer) {

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for domain, hosts := range r.lookups {
		fmt.Fprintln(w, domain)
		for _, h := range hosts {
			fmt.Fprintln(w, "\t", h.String())
		}
	}
}

func (r *Resolver) find(domain string) ([]*Host, bool) {

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	h, ok := r.lookups[domain]
	return h, ok
}

func (r *Resolver) addIP(domain string, ips ...net.IP) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, ip := range ips {
		r.lookups[domain] = append(r.lookups[domain], &Host{
			Network: "tcp",
			IP:      ip,
		})
	}
}

func (r *Resolver) addUnix(domain, path string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.lookups[domain] = append(r.lookups[domain], &Host{
		Network: "unix",
		Path:    path,
	})
}

func (r *Resolver) DialContext(
	ctx context.Context,
	network,
	addr string,
) (net.Conn, error) {

	sep := strings.LastIndex(addr, ":")
	domain, port := addr[:sep], addr[sep:]
	if hosts, ok := r.find(domain); ok {
		for _, h := range hosts {
			if c, err := h.Conn(ctx, port); err == nil {
				return c, nil
			}
		}
	}

	var d net.Dialer
	return d.DialContext(ctx, network, addr)
}
