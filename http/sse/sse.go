package sse

import (
	"io"
	"strings"

	"github.com/joscha-alisch/http4go/http/body"
)

type Message struct {
	Id    string
	Event string
	Data  []byte
}

func StreamFromBody(body body.Body) func() (*Message, error) {
	return func() (*Message, error) {
		bodyChunk := body.Next()
		if bodyChunk == nil || bodyChunk.IsDone() {
			return nil, nil
		}

		msg, err := readMessage(bodyChunk)
		if err != nil {
			return nil, err
		}
		return msg, nil
	}
}

func Stream(f func() (*Message, error)) func() (io.ReadCloser, error) {
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

func readMessage(chunk io.Reader) (*Message, error) {
	msg := &Message{}

	b, err := io.ReadAll(chunk)
	if err != nil {
		return nil, err
	}

	for line := range strings.Lines(string(b)) {
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
	return msg, nil
}
