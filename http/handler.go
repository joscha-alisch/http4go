package http

type Handler func(r Request) (Response, error)
type Filter func(next Handler) Handler

func (f Handler) AsServer(cfg ServerConfig) Server {
	return cfg.ToServer(f)
}

func (f Filter) Then(next Filter) Filter {
	return func(handler Handler) Handler {
		return f(next(handler))
	}
}

func (f Filter) Apply(handler Handler) Handler {
	return f(handler)
}
