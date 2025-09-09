package http

import (
	"fmt"
	"io"

	"github.com/joscha-alisch/http4go/http/body"
)

type memoryMessage struct {
	version string
	headers Headers
	body    body.Body
}

func (m memoryMessage) ToMessage(includeStream bool) string {
	if m.body.IsStream() && !includeStream {
		return fmt.Sprintf("%s\n<stream>", m.headers.String())
	}

	return fmt.Sprintf("%s\n%s", m.headers.String(), m.body.String())
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

func (m memoryMessage) Body(r io.ReadCloser) memoryMessage {
	m.body = body.FromStream(r)
	return m
}

func (m memoryMessage) BodyString(s string) memoryMessage {
	m.body = body.FromString(s)
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

func (m memoryMessage) GetBody() body.Body {
	return m.body
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
