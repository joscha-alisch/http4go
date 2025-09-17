package sse

import (
	"io"
	"strings"
	"testing"

	"github.com/joscha-alisch/http4go/http/body"
	"github.com/stretchr/testify/assert"
)

func TestStreamFromBody(t *testing.T) {
	tests := []struct {
		name       string
		bodyStream []string
		expected   []*Message
	}{
		{"single message", []string{
			"data: Hello World!\n\n",
		}, []*Message{
			{Data: []byte("Hello World!")},
		}},
		{"message with id and event", []string{
			"id: 1\nevent: greeting\ndata: Hello World!\n\n",
		}, []*Message{
			{Id: "1", Event: "greeting", Data: []byte("Hello World!")},
		}},
		{
			"multiple messages",
			[]string{
				"data: First Message\n\n",
				"data: Second Message\n\n",
			},
			[]*Message{
				{Data: []byte("First Message")},
				{Data: []byte("Second Message")},
			},
		},
		{
			"message with multiline data",
			[]string{
				"data: Line 1\ndata: Line 2\ndata: Line 3\n\n",
			},
			[]*Message{
				{Data: []byte("Line 1\nLine 2\nLine 3")},
			},
		},
		{
			"multiple messages in one chunk",
			[]string{
				"data: First Message\n\ndata: Second Message\n\n",
			},
			[]*Message{
				{Data: []byte("First Message")},
				{Data: []byte("Second Message")},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i := 0
			b := body.FromStream(func() (io.ReadCloser, error) {
				if i >= len(test.bodyStream) {
					return nil, nil
				}
				s := test.bodyStream[i]
				i++
				return io.NopCloser(strings.NewReader(s)), nil
			})

			nextFn := StreamFromBody(b)

			var results []*Message
			for {
				msg, err := nextFn()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if msg == nil {
					break
				}
				results = append(results, msg)
			}

			assert.Len(t, results, len(test.expected), "number of messages")
			for j, expectedMsg := range test.expected {
				assert.Equal(t, expectedMsg.Id, results[j].Id, "message id")
				assert.Equal(t, expectedMsg.Event, results[j].Event, "message event")
				assert.Equal(t, string(expectedMsg.Data), string(results[j].Data), "message data")
			}
		})
	}
}
