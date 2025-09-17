package http

import (
	"strings"
)

type SseMessage struct {
	Id    string
	Event string
	Data  []byte
}

func MessageFromChunk(chunk []byte) *SseMessage {
	msg := &SseMessage{}
	for line := range strings.Lines(string(chunk)) {
		if strings.HasPrefix(line, "id: ") {
			msg.Id = strings.TrimPrefix(line, "id: ")
		} else if strings.HasPrefix(line, "event: ") {
			msg.Event = strings.TrimPrefix(line, "event: ")
		} else if strings.HasPrefix(line, "data: ") {
			if msg.Data != nil {
				msg.Data = append(msg.Data, '\n')
			}
			msg.Data = append(msg.Data, []byte(strings.TrimPrefix(line, "data: "))...)
		}
	}
	return msg
}

type SseHandler func(Request) (Response, SseStream, error)

func (h SseHandler) AsServer(cfg SseServerConfig) Server {
	return cfg.ToSseServer(h)
}

type SseStream func() *SseMessage

type SseFilter func(next SseHandler) SseHandler

func (f SseFilter) Then(next SseFilter) SseFilter {
	return func(handler SseHandler) SseHandler {
		return f(next(handler))
	}
}

func (f SseFilter) Apply(next SseHandler) SseHandler {
	return f(next)
}
