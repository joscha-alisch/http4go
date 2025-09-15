package body

import (
	"bytes"
	"io"
)

type streamingBody struct {
	f      func() (io.Reader, error)
	peeked []byte
}

func (m *streamingBody) Next() Chunk {
	if m.peeked != nil {
		c := &chunk{
			Reader: bytes.NewReader(m.peeked),
			last:   true,
		}
		m.peeked = nil
		return c
	}

	c, err := m.f()
	if err != nil || c == nil {
		return nil
	}

	return &chunk{
		Reader: c,
		last:   true,
	}
}

func (m *streamingBody) Peek() Chunk {
	if m.peeked != nil {
		return &chunk{
			Reader: bytes.NewReader(m.peeked),
			last:   true,
		}
	}

	c, err := m.f()
	if err != nil || c == nil {
		return nil
	}

	m.peeked, err = io.ReadAll(c)
	if err != nil {
		return nil
	}

	return &chunk{
		Reader: bytes.NewReader(m.peeked),
		last:   true,
	}
}
