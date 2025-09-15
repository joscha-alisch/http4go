package body

import "bytes"

type memoryBody struct {
	b []byte
}

func (m *memoryBody) Next() Chunk {
	if m.b == nil {
		return nil
	}
	b := m.b
	m.b = nil
	return &chunk{
		Reader: bytes.NewReader(b),
		last:   true,
	}
}

func (m *memoryBody) Peek() Chunk {
	if m.b == nil {
		return nil
	}
	return &chunk{
		Reader: bytes.NewReader(m.b),
		last:   true,
	}
}
