package http

import (
	"github.com/joscha-alisch/http4go/http/method"
)

type Header struct {
	Name  string
	Value string
}
type Headers []Header

type StringBody string

type Route map[method.Method]Handler

type Routes map[string]Route

func (r Routes) AsHandler() Handler {
	return func(request Request) (Response, error) {
		route, ok := r[request.GetUri()]
		if !ok {
			return NewResponse(404).BodyString("Not Found"), nil
		}
		handler, ok := route[request.GetMethod()]
		if !ok {
			return NewResponse(405).BodyString("Method Not Allowed"), nil
		}
		return handler(request)
	}
}