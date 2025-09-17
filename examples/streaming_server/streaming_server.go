package main

import (
	"bytes"
	"io"
	"time"

	"github.com/joscha-alisch/http4go/http"
	"github.com/joscha-alisch/http4go/http/status"
	"github.com/joscha-alisch/http4go/servers"
)

func streamingHandler() http.Handler {
	return func(r http.Request) (http.Response, error) {
		return http.NewResponse(status.Ok).BodyStream(func() (io.ReadCloser, error) {
			time.Sleep(1 * time.Second)
			return io.NopCloser(bytes.NewReader([]byte("Hello World!\n"))), nil
		}), nil
	}
}

func main() {
	panic(streamingHandler().AsServer(servers.StdLib(8081)).StartBlocking())
}
