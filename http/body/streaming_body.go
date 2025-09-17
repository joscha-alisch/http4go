package body

import (
	"bytes"
	"io"
)

type streamingBody struct {
	f            func() (io.ReadCloser, error)
	peeked       []byte
	peekedIsLast bool
}

func (m *streamingBody) IsStream() bool {
	return true
}

func (m *streamingBody) Next() Chunk {
	if m.peeked != nil {
		c := &chunk{
			ReadCloser: io.NopCloser(bytes.NewReader(m.peeked)),
			done:       m.peekedIsLast,
		}
		m.peeked = nil
		m.peekedIsLast = false
		return c
	}

	c, err := m.f()
	if err != nil || c == nil {
		return &chunk{
			ReadCloser: io.NopCloser(bytes.NewReader(nil)),
			done:       true,
		}
	}

	return &chunk{
		ReadCloser: c,
		done:       false,
	}
}

func (m *streamingBody) Peek() Chunk {
	if m.peeked != nil {
		return &chunk{
			ReadCloser: io.NopCloser(bytes.NewReader(m.peeked)),
			done:       m.peekedIsLast,
		}
	}

	c, err := m.f()
	if err != nil || c == nil {
		return &chunk{
			ReadCloser: io.NopCloser(bytes.NewReader(nil)),
			done:       true,
		}
	}

	m.peeked, err = io.ReadAll(c)
	if err != nil {
		return nil
	}

	return &chunk{
		ReadCloser: io.NopCloser(bytes.NewReader(m.peeked)),
		done:       false,
	}
}
