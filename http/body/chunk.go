package body

import "io"

type Chunk interface {
	io.Reader
	IsLast() bool
	Into(v any) error
}

type chunk struct {
	io.Reader
	last bool
}

func (c *chunk) IsLast() bool {
	return c.last
}

func (c *chunk) Into(v any) error {
	panic("implement me")
}
