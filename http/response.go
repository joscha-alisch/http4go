package http

import (
	"fmt"
	"io"

	"github.com/joscha-alisch/http4go/http/status"
)

type Response interface {
	Status(s status.Status) Response
	GetStatus() status.Status

	ToMessage() string

	// Version sets the HTTP version of this message.
	Version(version string) Response

	// Header adds a header to this message.
	Header(name, value string) Response

	// Headers adds multiple headers to this message.
	Headers(headers Headers) Response

	// ReplaceHeader replaces a header in this message.
	ReplaceHeader(name string, value string) Response

	// ReplaceHeaders replaces all headers in this message with the given source headers.
	ReplaceHeaders(source Headers) Response

	// RemoveHeader removes a header from this message.
	RemoveHeader(name string) Response

	// RemoveHeaders removes headers with the given prefix from this message.
	RemoveHeaders(prefix string) Response

	// Body sets the body of this message.
	Body(body io.ReadCloser) Response

	// BodyString sets the body of this message from a string.
	BodyString(body string) Response

	// GetHeaders returns all headers of this message.
	GetHeaders() Headers

	// GetHeader returns the first value of the given header name.
	GetHeader(name string) string

	// GetHeaderValues returns all values of the given header name.
	GetHeaderValues(name string) []string

	// GetBody returns the body of this message as an io.Reader.
	GetBody() io.Reader

	// GetBodyString returns the body of this message as a string.
	GetBodyString() string

	// Close closes the body of this message.
	Close() error
}

type MemoryResponse struct {
	memoryMessage
	status status.Status
}

func NewResponse(status status.Status) Response {
	return MemoryResponse{
		memoryMessage: memoryMessage{
			version: "HTTP/1.1",
		},
	}.Status(status)
}

func (r MemoryResponse) Status(status status.Status) Response {
	r.status = status
	return r
}

func (r MemoryResponse) GetStatus() status.Status {
	return r.status
}

func (r MemoryResponse) Version(version string) Response {
	r.memoryMessage = r.memoryMessage.Version(version)
	return r
}

func (r MemoryResponse) Header(name, value string) Response {
	r.memoryMessage = r.memoryMessage.Header(name, value)
	return r
}

func (r MemoryResponse) Headers(headers Headers) Response {
	r.memoryMessage = r.memoryMessage.Headers(headers)
	return r
}

func (r MemoryResponse) ReplaceHeader(name string, value string) Response {
	r.memoryMessage = r.memoryMessage.ReplaceHeader(name, value)
	return r
}

func (r MemoryResponse) ReplaceHeaders(source Headers) Response {
	r.memoryMessage = r.memoryMessage.ReplaceHeaders(source)
	return r
}

func (r MemoryResponse) RemoveHeader(name string) Response {
	r.memoryMessage = r.memoryMessage.RemoveHeader(name)
	return r
}

func (r MemoryResponse) RemoveHeaders(prefix string) Response {
	r.memoryMessage = r.memoryMessage.RemoveHeaders(prefix)
	return r
}

func (r MemoryResponse) Body(body io.ReadCloser) Response {
	r.memoryMessage = r.memoryMessage.Body(body)
	return r
}

func (r MemoryResponse) BodyString(body string) Response {
	r.memoryMessage = r.memoryMessage.BodyString(body)
	return r
}

func (r MemoryResponse) ToMessage() string {
	return fmt.Sprintf("<<<<<<<< RESPONSE\n%s %d %s %s", r.version, r.status.Code, r.status.Text, r.memoryMessage.ToMessage())
}
