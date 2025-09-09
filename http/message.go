package http

import (
	"fmt"
	"io"
	"strings"
)

type memoryMessage struct {
	version string
	headers Headers
	body    io.ReadCloser
}

func (m memoryMessage) ToMessage() string {
	return fmt.Sprintf("HTTP/%s\nHeaders: %v\nBody: %v", m.version, m.headers, m.body)
}

func (m memoryMessage) Version(version string) memoryMessage {
	m.version = version
	return m
}

func (m memoryMessage) Header(name, value string) memoryMessage {
	m.headers = append(m.headers, Header{Name: name, Value: value})
	return m
}

func (m memoryMessage) Headers(headers Headers) memoryMessage {
	m.headers = append(m.headers, headers...)
	return m
}

func (m memoryMessage) ReplaceHeader(name string, value string) memoryMessage {
	m.RemoveHeader(name)
	m.headers = append(m.headers, Header{Name: name, Value: value})
	return m
}

func (m memoryMessage) ReplaceHeaders(source Headers) memoryMessage {
	m.headers = make(Headers, len(source))
	copy(m.headers, source)
	return m
}

func (m memoryMessage) RemoveHeader(name string) memoryMessage {
	headers := make(Headers, 0, len(m.headers))
	for i, h := range m.headers {
		if h.Name != name {
			headers = append(headers, m.headers[i])
		}
	}
	m.headers = headers
	return m
}

func (m memoryMessage) RemoveHeaders(prefix string) memoryMessage {
	headers := make(Headers, 0, len(m.headers))
	for i, h := range m.headers {
		if len(h.Name) < len(prefix) || h.Name[:len(prefix)] != prefix {
			headers = append(headers, m.headers[i])
		}
	}
	m.headers = headers
	return m
}

func (m memoryMessage) Body(body io.ReadCloser) memoryMessage {
	m.body = body
	return m
}

func (m memoryMessage) BodyString(body string) memoryMessage {
	m.body = io.NopCloser(io.Reader(strings.NewReader(body)))
	return m
}

func (m memoryMessage) GetHeaders() Headers {
	return m.headers
}

func (m memoryMessage) GetHeaderValues(name string) []string {
	values := make([]string, 0)
	for _, h := range m.headers {
		if h.Name == name {
			values = append(values, h.Value)
		}
	}
	return values
}

func (m memoryMessage) GetBody() io.Reader {
	return m.body
}

func (m memoryMessage) GetBodyString() string {
	if m.body == nil {
		return ""
	}
	b, err := io.ReadAll(m.body)
	if err != nil {
		return ""
	}
	return string(b)
}

func (m memoryMessage) Close() error {
	if m.body != nil {
		return m.body.Close()
	}
	return nil
}

func (m memoryMessage) GetHeader(name string) string {
	for _, h := range m.headers {
		if h.Name == name {
			return h.Value
		}
	}
	return ""
}