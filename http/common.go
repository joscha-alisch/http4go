package http

import (
	"github.com/joscha-alisch/http4go/http/method"
	"github.com/joscha-alisch/http4go/http/status"
)

type Header struct {
	Name  string
	Value string
}
type Headers []Header

func (h Headers) String() string {
	result := ""
	for _, header := range h {
		result += header.Name + ": " + header.Value + "\r\n"
	}
	return result
}

func (h Headers) AsMap() map[string][]string {
	result := make(map[string][]string)
	for _, header := range h {
		result[header.Name] = append(result[header.Name], header.Value)
	}
	return result
}

type StringBody string

type Route map[method.Method]Handler

type Routes map[string]Route

func (r Routes) AsHandler() Handler {
	return func(request Request) (Response, error) {
		route, ok := r[request.GetUri().GetPath()]
		if !ok {
			return NewResponse(status.NotFound).BodyString("Not Found"), nil
		}
		handler, ok := route[request.GetMethod()]
		if !ok {
			return NewResponse(status.NotFound).BodyString("Method Not Allowed"), nil
		}
		return handler(request)
	}
}
