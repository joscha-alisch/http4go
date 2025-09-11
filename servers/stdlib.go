package servers

import (
	"fmt"
	"io"
	nethttp "net/http"

	"github.com/joscha-alisch/http4go/http"
	"github.com/joscha-alisch/http4go/http/uri"
)

func StdLib(port int) *StdlibServerConfig {
	return &StdlibServerConfig{
		port: port,
	}
}

type StdlibServerConfig struct {
	port int
}

func (s *StdlibServerConfig) ToPolyServer(handler http.Handler, sseHandler http.SseHandler) http.Server {
	return &stdLibServer{
		port:       s.port,
		handler:    handler,
		sseHandler: sseHandler,
	}
}

func (s *StdlibServerConfig) ToServer(handler http.Handler) http.Server {
	return s.ToPolyServer(handler, nil)
}

func (s *StdlibServerConfig) ToSseServer(handler http.SseHandler) http.Server {
	return s.ToPolyServer(nil, handler)
}

type stdLibServer struct {
	port       int
	handler    http.Handler
	sseHandler http.SseHandler
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

	if s.sseHandler != nil {
		return s.startSseServer()
	}

	return fmt.Errorf("no handler configured")
}

func (s *stdLibServer) startHttpServer() error {
	return nethttp.ListenAndServe(fmt.Sprintf(":%d", s.port), nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		req := http.NewRequest().Method(r.Method).Uri(
			uri.NewUri().
				Scheme("http").
				Host(r.Host).
				Path(r.URL.Path)).Body(r.Body)

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
		_, _ = io.Copy(w, res.GetBody())
	}))
}

func (s *stdLibServer) startSseServer() error {
	return nethttp.ListenAndServe(fmt.Sprintf(":%d", s.port), nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		req := http.NewRequest().Method(r.Method).Uri(
			uri.NewUri().
				Scheme("http").
				Host(r.Host).
				Path(r.URL.Path)).Body(r.Body)

		headers := make(http.Headers, 0)
		for key, values := range r.Header {
			for _, value := range values {
				headers = append(headers, http.Header{Name: key, Value: value})
			}
		}

		req = req.Headers(headers)

		resp, err := s.sseHandler(req)
		if err != nil {
			panic(err)
		}

		if !resp.GetHandled() {
			w.WriteHeader(404)
			return
		}

		sse := &stdLibSse{c: make(chan http.SseMessage)}
		w.WriteHeader(resp.GetStatus())
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")

		f, ok := w.(nethttp.Flusher)
		if !ok {
			panic("streaming unsupported")
		}

		go resp.Consume(sse)

		for {
			msg, ok := <-sse.c
			if !ok {
				sse.Close()
			}
			_, _ = w.Write([]byte(fmt.Sprintf("%s\n\n", msg.Data)))
			f.Flush()
		}
	}))
}

type stdLibSse struct {
	c       chan http.SseMessage
	onClose func()
}

func (s *stdLibSse) Send(msg http.SseMessage) http.Sse {
	s.c <- msg
	return s
}

func (s *stdLibSse) OnClose(fn func()) http.Sse {
	s.onClose = fn
	return s
}

func (s *stdLibSse) Close() {
	close(s.c)
	if s.onClose != nil {
		s.onClose()
	}
}

func (s *stdLibServer) Stop() http.Server {
	return s
}
