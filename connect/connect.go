package connect

import (
	"fmt"

	"github.com/joscha-alisch/http4go/http"
)

type APIAction[T any] interface {
	ToRequest() http.Request
	ToResult(response http.Response, err error) (T, error)
}

type SSEAPIAction[T any] interface {
	ToRequest() http.Request
	ToEvent(message http.SseMessage) (T, error)
}

func Do[A any](transport http.Handler, action APIAction[A]) (A, error) {
	return action.ToResult(transport(action.ToRequest()))
}

func DoSse[A any](transport http.Handler, action SSEAPIAction[A]) (next func() (A, error), err error) {
	resp, err := transport(action.ToRequest())
	if err != nil {
		return nil, err
	}

	if resp.GetStatus().Code < 200 || resp.GetStatus().Code >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.GetStatus().Code)
	}

	return func() (A, error) {

		return action.ToEvent(http.SseMessage{})
	}, nil
}
