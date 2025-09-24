package signature

type Option func(o *Signature)

func HeaderNameSignature(v string) Option {
	return func(o *Signature) {
		o.headerName = v
	}
}

func HeaderNameClock(v string) Option {
	return func(o *Signature) {
		o.clockName = v
	}
}

func HeaderNameURL(v string) Option {
	return func(o *Signature) {
		o.originUrlName = v
	}
}
