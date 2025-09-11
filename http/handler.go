package http

import "github.com/joscha-alisch/http4go/http/status"

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

type SseHandler func(r Request) (SseResponse, error)

func (h SseHandler) AsServer(cfg PolyServerConfig) Server {
	return cfg.ToSseServer(h)
}

type SseResponse interface {
	Status(s status.Status) SseResponse
	Headers(h Headers) SseResponse
	Handled(b bool) SseResponse

	GetStatus() int
	GetHeaders() Headers
	GetHandled() bool

	Consume(Sse)
}

type SseConsumer func(Sse)

type Sse interface {
	Send(msg SseMessage) Sse
	OnClose(fn func()) Sse
	Close()
}

type SseMessage struct {
	Data string
}

func NewSseResponse(consume SseConsumer) SseResponse {
	return MemorySseResponse{
		status:  status.Ok,
		headers: Headers{},
		handled: true,
		consume: consume,
	}
}

type MemorySseResponse struct {
	status  status.Status
	headers Headers
	handled bool
	consume SseConsumer
}

func (m MemorySseResponse) Consume(sse Sse) {
	m.consume(sse)
}

func (m MemorySseResponse) Status(s status.Status) SseResponse {
	m.status = s
	return m
}

func (m MemorySseResponse) Headers(h Headers) SseResponse {
	m.headers = h
	return m
}

func (m MemorySseResponse) Handled(b bool) SseResponse {
	m.handled = b
	return m
}

func (m MemorySseResponse) GetStatus() int {
	return m.status.Code
}

func (m MemorySseResponse) GetHeaders() Headers {
	return m.headers
}

func (m MemorySseResponse) GetHandled() bool {
	return m.handled
}

var myHandler = SseHandler(func(r Request) (SseResponse, error) {
	return MemorySseResponse{
		status: status.Ok,
		consume: func(s Sse) {
			s.Send(SseMessage{})
		},
	}, nil
})
