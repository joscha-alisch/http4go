package filters

import (
	"fmt"

	"github.com/joscha-alisch/http4go/http"
	"github.com/joscha-alisch/http4go/http/uri"
)

func SetHostFrom(u uri.Uri) http.Filter {
	return func(next http.Handler) http.Handler {
		return func(r http.Request) (http.Response, error) {
			host := u.GetHost()
			port := u.GetPort()
			hostPort := host
			if port != 0 {
				hostPort = fmt.Sprintf("%s:%d", host, port)
			}
			return next(r.Uri(r.GetUri().Scheme(u.GetScheme()).Host(u.GetHost()).Port(u.GetPort())).
				ReplaceHeader("Host", hostPort))
		}
	}
}
