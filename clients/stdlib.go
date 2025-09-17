package clients

import (
	nethttp "net/http"
	"net/url"

	"github.com/joscha-alisch/http4go/http"
	"github.com/joscha-alisch/http4go/http/status"
	"github.com/joscha-alisch/http4go/http/uri"
)

var StdLib = http.Handler(func(r http.Request) (http.Response, error) {
	client := nethttp.Client{}

	resp, err := client.Do(&nethttp.Request{
		Method: r.GetMethod(),
		URL:    urlFromUri(r.GetUri()),
		Header: r.GetHeaders().AsMap(),
		Body:   r.GetBody().Next(),
	})
	if err != nil {
		return nil, err
	}

	var h http.Headers

	for name, values := range resp.Header {
		for _, value := range values {
			h = append(h, http.Header{Name: name, Value: value})
		}
	}

	return http.NewResponse(status.Status{
		Code: resp.StatusCode,
		Text: resp.Status,
	}).BodyReader(resp.Body).Headers(h), nil
})

func urlFromUri(u uri.Uri) *url.URL {
	return &url.URL{
		Scheme:   u.GetScheme(),
		User:     nil,
		Host:     u.GetHostPort(),
		Path:     u.GetPath(),
		RawQuery: u.GetQuery(),
		Fragment: u.GetFragment(),
	}
}
