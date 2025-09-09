package body

import (
	"bytes"
	"io"
)

type Body interface {
	io.ReadCloser
	String() string
	IsStream() bool
}

type memoryBackedBody struct {
	io.ReadCloser
	cached   []byte
	isStream bool
}

func FromBytes(b []byte) Body {
	return &memoryBackedBody{
		ReadCloser: io.NopCloser(bytes.NewBuffer(b)),
		cached:     b,
		isStream:   false,
	}
}

func FromString(s string) Body {
	return FromBytes([]byte(s))
}

func FromStream(r io.ReadCloser) Body {
	return &memoryBackedBody{
		ReadCloser: r,
		cached:     nil,
		isStream:   true,
	}
}

func (b *memoryBackedBody) String() string {
	if b.cached != nil {
		return string(b.cached)
	}

	if b.ReadCloser == nil {
		return ""
	}

	data, err := io.ReadAll(b.ReadCloser)
	if err != nil {
		return ""
	}
	b.cached = data
	b.ReadCloser = io.NopCloser(bytes.NewBuffer(b.cached))
	return string(data)
}

func (b *memoryBackedBody) IsStream() bool {
	return b.isStream
}
