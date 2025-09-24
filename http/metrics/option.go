package metrics

import "github.com/prometheus/client_golang/prometheus"

type Option func(*Metrics)

func WithLabels(kv map[string]string) Option {
	return func(m *Metrics) {
		if kv != nil {
			if m.opts.ConstLabels == nil {
				m.opts.ConstLabels = make(prometheus.Labels)
			}
			for k, v := range kv {
				m.opts.ConstLabels[k] = v
			}
		}
	}
}
