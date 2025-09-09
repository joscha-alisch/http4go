package main

import (
	"github.com/joscha-alisch/http4go/http"
	"github.com/joscha-alisch/http4go/http/method"
	"github.com/joscha-alisch/http4go/servers"
)

func HelloWorld() http.Handler {
	return http.Routes{
		"/": {
			method.GET: func(r http.Request) (http.Response, error) {
				return http.NewResponse(200).BodyString("Hello, World!"), nil
			},
		},
	}.AsHandler()
}

func main() {
	panic(HelloWorld().AsServer(servers.StdLib(8080)).StartBlocking())
}