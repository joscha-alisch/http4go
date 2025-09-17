package openai

import (
	"encoding/json"
	"fmt"

	"github.com/joscha-alisch/http4go/connect"
	"github.com/joscha-alisch/http4go/http"
	"github.com/joscha-alisch/http4go/http/body"
	"github.com/joscha-alisch/http4go/http/sse"
	"github.com/joscha-alisch/http4go/http/uri"
)

func (c Client) ChatCompletions(request ChatCompletionsRequest) (next func() (ChatCompletionsChunk, error), err error) {
	action, err := NewChatCompletionsAction(request)
	if err != nil {
		return nil, err
	}
	return connect.DoSse(c.transport, action)
}

type ChatCompletionsAction struct {
	Body body.Body
}

func NewChatCompletionsAction(request ChatCompletionsRequest) (ChatCompletionsAction, error) {
	b, err := body.FromJson(request)
	if err != nil {
		return ChatCompletionsAction{}, err
	}

	return ChatCompletionsAction{
		Body: b,
	}, nil
}

func (c ChatCompletionsAction) ToRequest() http.Request {
	return http.NewRequest().
		Method("POST").
		Uri(uri.NewUri().Path("/v1/chat/completions")).
		Header("Content-Type", "application/json").
		Header("Accept", "text/event-stream").
		Body(c.Body)
}

func (c ChatCompletionsAction) ToEvent(message sse.Message) (ChatCompletionsChunk, error) {
	var chunk ChatCompletionsChunk
	return chunk, json.Unmarshal(message.Data, &chunk)
}

type ChatCompletionsRequest struct {
	Stream   bool                    `json:"stream,omitempty"`
	Model    string                  `json:"model"`
	Messages []ChatCompletionMessage `json:"messages"`
}

type ChatCompletionMessage struct {
	Role    string                        `json:"role"`
	Content ChatCompletionsMessageContent `json:"content"`
	Name    string                        `json:"name,omitempty"`
}

type ChatCompletionsMessageContent struct {
	Text  *string
	Parts []ChatCompletionsMessageContentPart
}

func (c *ChatCompletionsMessageContent) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return fmt.Errorf("content: empty json")
	}
	switch b[0] {
	case '"': // string
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		c.Text = &s
		c.Parts = nil
		return nil
	case '[': // array
		var parts []ChatCompletionsMessageContentPart
		if err := json.Unmarshal(b, &parts); err != nil {
			return err
		}
		c.Text = nil
		c.Parts = parts
		return nil
	case 'n': // null
		*c = ChatCompletionsMessageContent{}
		return nil
	default:
		return fmt.Errorf("content: expected string or array, got %q...", b[0])
	}
}

func (c *ChatCompletionsMessageContent) MarshalJSON() ([]byte, error) {
	switch {
	case c.Text != nil && len(c.Parts) == 0:
		return json.Marshal(*c.Text)
	case c.Text == nil && len(c.Parts) > 0:
		return json.Marshal(c.Parts)
	case c.Text == nil && len(c.Parts) == 0:
		// choose your policy: null or empty string/array
		return []byte("null"), nil
	default:
		return nil, fmt.Errorf("content: both Text and Parts set")
	}
}

type ChatCompletionsMessageContentPart struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

type ChatCompletionsChunk struct {
	Model             string                       `json:"model"`
	Choices           []ChatCompletionsChunkChoice `json:"choices"`
	Created           int64                        `json:"created"`
	ID                string                       `json:"id"`
	Object            string                       `json:"object"`
	SystemFingerprint string                       `json:"system_fingerprint,omitempty"`
	Usage             map[string]any               `json:"usage,omitempty"`
}

type ChatCompletionsChunkChoice struct {
	Delta        ChatCompletionsChunkChoiceDelta `json:"delta"`
	Index        int                             `json:"index"`
	FinishReason string                          `json:"finish_reason,omitempty"`
}

type ChatCompletionsChunkChoiceDelta struct {
	Content string `json:"content,omitempty"`
	Role    string `json:"role,omitempty"`
}
