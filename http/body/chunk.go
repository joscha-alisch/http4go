package body

import "io"

type Chunk interface {
	io.ReadCloser
	IsDone() bool
	Into(v any) error
}

type chunk struct {
	io.ReadCloser
	done bool
}

func (c *chunk) IsDone() bool {
	return c.done
}

func (c *chunk) Into(v any) error {
	panic("implement me")
}
