package servers

import (
	"bytes"
	"io"
	"net"
	nethttp "net/http"
	"testing"

	"github.com/joscha-alisch/http4go/http"
	"github.com/joscha-alisch/http4go/http/status"
	"github.com/stretchr/testify/assert"
)

func TestHttpRequest(t *testing.T) {
	url := setupServer(t, func(r http.Request) (http.Response, error) {
		return http.NewResponse(status.Ok).BodyString("Hello, World!"), nil
	})

	resp, err := nethttp.Get(url)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "Hello, World!", readBody(t, resp))
}

func TestStreamingResponse(t *testing.T) {
	url := setupServer(t, func(request http.Request) (http.Response, error) {
		data := []string{"Hello ", "World ", "from ", "http4go"}
		i := 0
		return http.NewResponse(status.Ok).BodyStream(func() (io.ReadCloser, error) {
			if i >= len(data) {
				return nil, nil
			}

			msg := io.NopCloser(bytes.NewReader([]byte(data[i])))
			i++
			return msg, nil
		}), nil
	})

	resp, err := nethttp.Get(url)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.Equal(t, 200, resp.StatusCode)
	body := readBody(t, resp)
	expected := `Hello World from http4go`

	assert.Equal(t, expected, body)
}

func setupServer(t *testing.T, handler http.Handler) (url string) {
	ln, err := net.Listen("tcp", ":0")
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	server := StdLib(0, WithListener(ln)).ToServer(handler).Start()
	t.Cleanup(func() { server.Stop() })

	return "http://" + ln.Addr().String()
}

func readBody(t *testing.T, r *nethttp.Response) string {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
