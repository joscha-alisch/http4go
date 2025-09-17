package body

import (
	"bytes"
	"io"
)

type readerBody struct {
	r      io.Reader
	peeked []byte
}

func (m *readerBody) Into(t any) error {
	return Into(m, t)
}

func (m *readerBody) IsStream() bool {
	return false
}

func (m *readerBody) Next() Chunk {
	if m.r == nil {
		return &chunk{
			ReadCloser: io.NopCloser(nil),
			done:       true,
		}
	}

	r := m.r
	m.r = nil
	m.peeked = nil

	return &chunk{
		ReadCloser: io.NopCloser(r),
		done:       false,
	}
}

func (m *readerBody) Peek() Chunk {
	if m.peeked != nil {
		return &chunk{
			ReadCloser: io.NopCloser(bytes.NewReader(m.peeked)),
			done:       false,
		}
	}

	if m.r == nil {
		return &chunk{
			ReadCloser: io.NopCloser(nil),
			done:       true,
		}
	}

	peeked, err := io.ReadAll(m.r)
	if err != nil {
		return nil
	}

	m.peeked = peeked

	m.r = bytes.NewReader(peeked)
	return &chunk{
		ReadCloser: io.NopCloser(bytes.NewReader(peeked)),
		done:       false,
	}
}
