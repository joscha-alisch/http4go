package body

import (
	"encoding/json"
	"io"
)

type Body interface {
	IsStream() bool

	// Next returns the next chunk of the body, consuming it.
	// If the body is not a stream, Next always returns the full body.
	Next() Chunk

	// Peek returns a reader that can be used to read the next chunk of the body without consuming it -
	// useful for debugging and logging. If the body is not a stream, Peek always returns the full body.
	Peek() Chunk
}

func FromBytes(b []byte) Body {
	return &memoryBody{b}
}

func FromString(s string) Body {
	return FromBytes([]byte(s))
}

func FromJson(v any) (Body, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return FromBytes(b), nil
}

func FromReader(r io.Reader) Body {
	return &readerBody{r: r}
}

func FromStream(f func() (io.ReadCloser, error)) Body {
	return &streamingBody{f: f}
}
