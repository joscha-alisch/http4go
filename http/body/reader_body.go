package body

import (
	"bytes"
	"io"
)

type readerBody struct {
	r io.Reader
}

func (m *readerBody) Next() Chunk {
	if m.r == nil {
		return nil
	}

	r := m.r
	m.r = nil

	return &chunk{
		Reader: r,
		last:   true,
	}
}

func (m *readerBody) Peek() Chunk {
	if m.r == nil {
		return nil
	}
	peeked, err := io.ReadAll(m.r)
	if err != nil {
		return nil
	}

	m.r = bytes.NewReader(peeked)
	return &chunk{
		Reader: bytes.NewReader(peeked),
		last:   true,
	}
}
