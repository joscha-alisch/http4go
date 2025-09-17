package main

import (
	"fmt"
	"time"

	"github.com/joscha-alisch/http4go/http"
	"github.com/joscha-alisch/http4go/http/sse"
	"github.com/joscha-alisch/http4go/http/status"
	"github.com/joscha-alisch/http4go/servers"
)

func streamingHandler() http.Handler {
	return func(r http.Request) (http.Response, error) {
		i := 0
		return http.NewResponse(status.Ok).BodyStream(sse.Stream(func() (*sse.SseMessage, error) {
			time.Sleep(time.Second)
			i++
			if i > 5 {
				return nil, nil
			}

			return &sse.SseMessage{
				Id:   fmt.Sprintf("%d", i),
				Data: []byte("Hello World!"),
			}, nil
		})), nil
	}
}

func main() {
	panic(streamingHandler().AsServer(servers.StdLib(8082)).StartBlocking())
}
