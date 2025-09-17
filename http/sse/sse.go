package sse

import (
	"bufio"
	"errors"
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
	r := bufio.NewReader(&chunkedReader{body: body})

	return func() (*Message, error) {
		return readMessage(r)
	}
}

func readMessage(r *bufio.Reader) (*Message, error) {
	var msg *Message
	for {
		line, err := r.ReadString('\n')
		if len(line) > 0 {
			s := strings.TrimRight(line, "\r\n")
			if s == "" {
				if msg != nil {
					return msg, nil
				}
				continue
			}

			if msg == nil {
				msg = &Message{}
			}

			switch {
			case strings.HasPrefix(s, "id: "):
				msg.Id = strings.TrimPrefix(s, "id: ")
			case strings.HasPrefix(s, "event: "):
				msg.Event = strings.TrimPrefix(s, "event: ")
			case strings.HasPrefix(s, "data: "):
				if len(msg.Data) > 0 {
					msg.Data = append(msg.Data, '\n')
				}
				msg.Data = append(msg.Data, s[len("data: "):]...)
				// You can add support for "retry:" etc. here if needed.
			}
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				if msg != nil {
					out := msg
					msg = nil
					return out, nil
				}
				return nil, nil
			}
			return nil, err
		}
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

type chunkedReader struct {
	body body.Body
	cur  io.ReadCloser
	done bool
}

func (cr *chunkedReader) Read(p []byte) (n int, err error) {
	for {
		if cr.done {
			return 0, io.EOF
		}
		if cr.cur == nil {
			nxt := cr.body.Next()
			if nxt == nil || nxt.IsDone() {
				cr.done = true
				return 0, io.EOF
			}
			cr.cur = nxt
		}

		n, err := cr.cur.Read(p)
		if err == io.EOF {
			// Exhausted this chunk — try the next one on the next loop.
			cr.cur = nil
			if n > 0 {
				// Return what we read; caller will come back for more.
				return n, nil
			}
			continue
		}
		return n, err
	}
}
