package filters

import (
	"fmt"
	"io"

	"github.com/joscha-alisch/http4go/http"
)

func PrintRequestAndResponse(to io.Writer, includeStream bool) http.Filter {
	return PrintRequest(to, includeStream).Then(PrintResponse(to, includeStream))
}

func PrintRequest(to io.Writer, includeStream bool) http.Filter {
	return func(next http.Handler) http.Handler {
		return func(req http.Request) (http.Response, error) {
			_, err := fmt.Fprintln(to, req.ToMessage(includeStream))
			if err != nil {
				return nil, err
			}
			return next(req)
		}
	}
}

func PrintResponse(to io.Writer, includeStream bool) http.Filter {
	return func(next http.Handler) http.Handler {
		return func(req http.Request) (http.Response, error) {
			resp, err := next(req)
			if err != nil {
				return nil, err
			}

			_, err = fmt.Fprintln(to, resp.ToMessage(includeStream))
			if err != nil {
				return nil, err
			}

			return resp, nil
		}
	}
}
