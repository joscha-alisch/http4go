package connect

import (
	"github.com/joscha-alisch/http4go/http"
)

type APIAction[T any] interface {
	ToRequest() http.Request
	ToResult(response http.Response, err error) (T, error)
}

func Do[A any](transport http.Handler, action APIAction[A]) (A, error) {
	return action.ToResult(transport(action.ToRequest()))
}
