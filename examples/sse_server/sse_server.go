package main

import (
	"time"

	"github.com/joscha-alisch/http4go/http"
	"github.com/joscha-alisch/http4go/servers"
)

func sseHandler() http.SseHandler {
	return func(r http.Request) (http.SseResponse, error) {
		return http.NewSseResponse(func(sse http.Sse) {
			for {
				sse.Send(http.SseMessage{
					Data: "Hello, World!",
				})

				time.Sleep(1 * time.Second)
			}
		}), nil
	}
}

func main() {
	panic(sseHandler().AsServer(servers.StdLib(8081)).StartBlocking())
}
