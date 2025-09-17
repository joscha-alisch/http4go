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
			body: FromStream(func() func() (io.ReadCloser, error) {
				i := 0
				return func() (io.ReadCloser, error) {
					if i > 1 {
						return nil, nil
					}
					i++
					return io.NopCloser(bytes.NewReader([]byte(fmt.Sprintf("stream body %d", i)))), nil
				}
			}()),
			expected: []string{
				"stream body 1",
				"stream body 2",
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
		t.Run(fmt.Sprintf("[%d] First Peek", i), func(t *testing.T) {
			chunk := b.Peek()
			assert.Equal(t, false, chunk.Done())
			assert.Equal(t, s, readAll(t, chunk))
		})

		t.Run(fmt.Sprintf("[%d] Second Peek", i), func(t *testing.T) {
			chunk := b.Peek()
			assert.Equal(t, false, chunk.Done())
			assert.Equal(t, s, readAll(t, chunk))
		})

		t.Run(fmt.Sprintf("[%d] Next", i), func(t *testing.T) {
			chunk := b.Next()
			assert.Equal(t, false, chunk.Done())
			assert.Equal(t, s, readAll(t, chunk))
		})

	}

	t.Run("After done it's nil", func(t *testing.T) {
		chunk := b.Peek()
		assert.True(t, chunk.Done())
		chunk = b.Next()
		assert.True(t, chunk.Done())
	})
}

func readAll(t *testing.T, r Chunk) string {
	b, err := io.ReadAll(r)
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	return string(b)
}
