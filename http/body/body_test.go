package body

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBody(t *testing.T) {
	tests := []struct {
		name     string
		body     Body
		expected []string
	}{
		{
			name:     "string body",
			body:     FromString("string body"),
			expected: []string{"string body"},
		},
		{
			name:     "byte body",
			body:     FromBytes([]byte("byte body")),
			expected: []string{"byte body"},
		},
		{
			name: "json body",
			body: func() Body {
				b, err := FromJson(map[string]string{"key": "value"})
				if err != nil {
					t.Fatalf("failed to create json body: %v", err)
				}
				return b
			}(),
			expected: []string{`{"key":"value"}`},
		},
		{
			name: "reader body",
			body: FromReader(bytes.NewReader([]byte("reader body"))),
			expected: []string{
				"reader body",
			},
		},
		{
			name:     "nil reader",
			body:     FromReader(nil),
			expected: []string{},
		},
		{
			name: "reader stream",
			body: FromStream(func() func() (io.Reader, error) {
				i := 0
				return func() (io.Reader, error) {
					if i > 0 {
						return nil, nil
					}
					i++
					return bytes.NewReader([]byte("stream body")), nil
				}
			}()),
			expected: []string{
				"stream body",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testBody(t, test.body, test.expected)
		})
	}
}

func testBody(t *testing.T, b Body, expected []string) {
	for i, s := range expected {
		shouldBeLast := i == len(expected)-1
		t.Run(fmt.Sprintf("[%d] First Peek", i), func(t *testing.T) {
			chunk := b.Peek()
			assert.Equal(t, shouldBeLast, chunk.IsLast())
			assert.Equal(t, s, readAll(t, chunk))
		})

		t.Run(fmt.Sprintf("[%d] Second Peek", i), func(t *testing.T) {
			chunk := b.Peek()
			assert.Equal(t, shouldBeLast, chunk.IsLast())
			assert.Equal(t, s, readAll(t, chunk))
		})

		t.Run(fmt.Sprintf("[%d] Next", i), func(t *testing.T) {
			chunk := b.Next()
			assert.Equal(t, shouldBeLast, chunk.IsLast())
			assert.Equal(t, s, readAll(t, chunk))
		})

	}

	t.Run("After last it's nil", func(t *testing.T) {
		chunk := b.Peek()
		assert.Nil(t, chunk, "expected no more chunks")
		chunk = b.Next()
		assert.Nil(t, chunk, "expected no more chunks")
	})
}

func readAll(t *testing.T, r Chunk) string {
	b, err := io.ReadAll(r)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	return string(b)
}
