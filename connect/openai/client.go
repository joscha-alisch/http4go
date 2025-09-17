package openai

import (
	"github.com/joscha-alisch/http4go/http"
)

type Client struct {
	transport http.Handler
}

func NewClient(transport http.Handler) Client {
	return Client{transport: transport}
}
