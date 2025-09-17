package body

import "io"

type Chunk interface {
	io.ReadCloser
	Done() bool
	Into(v any) error
}

type chunk struct {
	io.ReadCloser
	done bool
}

func (c *chunk) Done() bool {
	return c.done
}

func (c *chunk) Into(v any) error {
	panic("implement me")
}
