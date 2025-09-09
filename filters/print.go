package filters

import (
	"fmt"
	"io"

	"github.com/joscha-alisch/http4go/http"
)

func PrintRequestAndResponse(to io.Writer) http.Filter {
	return PrintRequest(to).Then(PrintResponse(to))
}

func PrintRequest(to io.Writer) http.Filter {
	return func(next http.Handler) http.Handler {
		return func(req http.Request) (http.Response, error) {
			_, err := fmt.Fprintln(to, req.ToMessage())
			if err != nil {
				return nil, err
			}
			return next(req)
		}
	}
}

func PrintResponse(to io.Writer) http.Filter {
	return func(next http.Handler) http.Handler {
		return func(req http.Request) (http.Response, error) {
			resp, err := next(req)
			if err != nil {
				return nil, err
			}

			_, err = fmt.Fprintln(to, resp.ToMessage())
			if err != nil {
				return nil, err
			}

			return resp, nil
		}
	}
}
