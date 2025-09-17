package openai

import (
	"github.com/joscha-alisch/http4go/http"
)

type Client struct {
	transport    http.Handler
	sseTransport http.SseHandler
}

func NewClient(transport http.Handler, sseTransport http.SseHandler) Client {
	return Client{transport: transport, sseTransport: sseTransport}
}
