package servers

import (
	"fmt"
	"io"
	"net"
	nethttp "net/http"

	"github.com/joscha-alisch/http4go/http"
	"github.com/joscha-alisch/http4go/http/uri"
)

func StdLib(port int, opts ...ServerOption) *StdlibServerConfig {
	cfg := StdlibServerConfig{
		port: port,
	}
	for _, opt := range opts {
		cfg = opt(cfg)
	}
	return &cfg
}

type ServerOption func(StdlibServerConfig) StdlibServerConfig

func WithListener(l net.Listener) ServerOption {
	return func(cfg StdlibServerConfig) StdlibServerConfig {
		cfg.listener = l
		return cfg
	}
}

type StdlibServerConfig struct {
	port     int
	listener net.Listener
}

func (s *StdlibServerConfig) ToServer(handler http.Handler) http.Server {
	return &stdLibServer{
		port:     s.port,
		handler:  handler,
		listener: s.listener,
	}
}

type stdLibServer struct {
	port     int
	handler  http.Handler
	listener net.Listener
}

func (s *stdLibServer) Start() http.Server {
	go func() {
		err := s.StartBlocking()
		if err != nil {
			panic(err)
		}
	}()
	return s
}

func (s *stdLibServer) StartBlocking() error {
	if s.handler != nil {
		return s.startHttpServer()
	}

	return fmt.Errorf("no handler configured")
}

func (s *stdLibServer) startHttpServer() error {
	l := s.listener
	if l == nil {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
		if err != nil {
			return err
		}
		s.listener = l
	}

	return nethttp.Serve(s.listener, nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		req := http.NewRequest().Method(r.Method).Uri(
			uri.NewUri().
				Scheme("http").
				Host(r.Host).
				Path(r.URL.Path)).BodyReader(r.Body)

		headers := make(http.Headers, 0)
		for key, values := range r.Header {
			for _, value := range values {
				headers = append(headers, http.Header{Name: key, Value: value})
			}
		}

		req = req.Headers(headers)

		res, err := s.handler(req)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(res.GetStatus().Code)

		for {
			chunk := res.GetBody().Next()
			_, err = io.Copy(w, chunk)
			if err != nil {
				panic(err)
			}
			if f, ok := w.(nethttp.Flusher); ok {
				f.Flush()
			}

			if chunk.IsDone() {
				return
			}
		}
	}))
}

func (s *stdLibServer) Stop() http.Server {
	return s
}
