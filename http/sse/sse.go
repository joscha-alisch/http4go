package sse

import (
	"io"
	"strings"
)

type SseMessage struct {
	Id    string
	Event string
	Data  []byte
}

func Stream(f func() (*SseMessage, error)) func() (io.ReadCloser, error) {
	return func() (io.ReadCloser, error) {
		msg, err := f()
		if err != nil {
			return nil, err
		}

		if msg == nil {
			return nil, nil
		}

		var b strings.Builder
		if msg.Id != "" {
			b.WriteString("id: ")
			b.WriteString(msg.Id)
			b.WriteString("\n")
		}
		if msg.Event != "" {
			b.WriteString("event: ")
			b.WriteString(msg.Event)
			b.WriteString("\n")
		}
		if msg.Data != nil {
			for line := range strings.Lines(string(msg.Data)) {
				b.WriteString("data: ")
				b.WriteString(line)
				b.WriteString("\n")
			}
		}
		b.WriteString("\n")
		return io.NopCloser(strings.NewReader(b.String())), nil
	}
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
