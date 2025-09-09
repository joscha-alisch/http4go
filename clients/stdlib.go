package clients

import (
	nethttp "net/http"
	"net/url"

	"github.com/joscha-alisch/http4go/http"
)

var StdLib = http.Handler(func(r http.Request) (http.Response, error) {
	client := nethttp.Client{}
	parsedUrl, err := url.Parse(r.GetUri())
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(&nethttp.Request{
		Method: r.GetMethod(),
		URL:    parsedUrl,
		Header: nil,
		Body:   nil,
	})
	if err != nil {
		return nil, err
	}

	return http.NewResponse(resp.StatusCode).
		Body(resp.Body), nil
})