package resolver

import "net"

type Option func(r *Resolver)

func WithDNS(domain string, ip ...net.IP) Option {
	return func(r *Resolver) {
		r.addIP(domain, ip...)
	}
}

func WithUnix(domain, path string) Option {
	return func(r *Resolver) {
		r.addUnix(domain, path)
	}
}

func WithDNSMap(dns map[string][]string) Option {

	return func(r *Resolver) {

		if dns == nil {
			return
		}

		for domain, ips := range dns {
			var ip []net.IP
			for _, s := range ips {
				if v, err := IP(s).Value(); err == nil {
					ip = append(ip, v)
				}
			}
			if len(ip) > 0 {
				WithDNS(domain, ip...)(r)
			}
		}
		return
	}
}

func WithUnixMap(unix map[string]string) Option {
	return func(r *Resolver) {

		if unix == nil {
			return
		}

		for domain, path := range unix {
			WithUnix(domain, path)(r)
		}
	}
}
