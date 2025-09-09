package main

import (
	"os"

	"github.com/joscha-alisch/http4go/filters"
	"github.com/joscha-alisch/http4go/http"
	"github.com/joscha-alisch/http4go/http/method"
	"github.com/joscha-alisch/http4go/http/status"
	"github.com/joscha-alisch/http4go/servers"
)

func helloWorld(r http.Request) (http.Response, error) {
	return http.NewResponse(status.Ok).BodyString("Hello, World!"), nil
}

func App() http.Handler {
	return http.Routes{
		"/": {
			method.GET: filters.PrintRequestAndResponse(os.Stdout, true).Apply(helloWorld),
		},
	}.AsHandler()
}

func main() {
	panic(App().AsServer(servers.StdLib(8080)).StartBlocking())
}
