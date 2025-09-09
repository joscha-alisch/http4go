package servers

import (
	"fmt"
	"io"
	nethttp "net/http"

	"github.com/joscha-alisch/http4go/http"
)

func StdLib(port int) *StdlibServerConfig {
	return &StdlibServerConfig{
		port: port,
	}
}

type StdlibServerConfig struct {
	port int
}

func (s *StdlibServerConfig) ToServer(handler http.Handler) http.Server {
	return &stdLibServer{
		port:    s.port,
		handler: handler,
	}
}

type stdLibServer struct {
	port    int
	handler http.Handler
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
	return nethttp.ListenAndServe(fmt.Sprintf(":%d", s.port), nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		req := http.NewRequest().Method(r.Method).Uri(r.URL.String())

		res, err := s.handler(req)
		if err != nil {
			panic(err)
		}

		w.WriteHeader(res.GetStatus())
		_, _ = io.Copy(w, res.GetBody())
	}))
}

func (s *stdLibServer) Stop() http.Server {
	return s
}