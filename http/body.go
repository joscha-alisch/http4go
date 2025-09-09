package http

type Body interface {
	Close()
}

type MemoryBody []byte

func (b MemoryBody) Close() {}