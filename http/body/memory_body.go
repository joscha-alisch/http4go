package body

import (
	"bytes"
	"io"
)

type memoryBody struct {
	b []byte
}

func (m *memoryBody) Into(t any) error {
	return Into(m, t)
}

func (m *memoryBody) IsStream() bool {
	return false
}

func (m *memoryBody) Next() Chunk {
	if m.b == nil {
		return &chunk{
			ReadCloser: io.NopCloser(bytes.NewReader(nil)),
			done:       true,
		}
	}
	b := m.b
	m.b = nil
	return &chunk{
		ReadCloser: io.NopCloser(bytes.NewReader(b)),
		done:       false,
	}
}

func (m *memoryBody) Peek() Chunk {
	if m.b == nil {
		return &chunk{
			ReadCloser: io.NopCloser(bytes.NewReader(nil)),
			done:       true,
		}
	}
	return &chunk{
		ReadCloser: io.NopCloser(bytes.NewReader(m.b)),
		done:       false,
	}
}
