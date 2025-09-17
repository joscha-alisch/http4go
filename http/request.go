package http

import (
	"io"

	"github.com/joscha-alisch/http4go/http/body"
	"github.com/joscha-alisch/http4go/http/method"
	"github.com/joscha-alisch/http4go/http/uri"
)

type Request interface {
	Method(method method.Method) Request
	Uri(uri uri.Uri) Request
	Query(name, value string) Request
	// GetMethod returns the HTTP method of this request.
	GetMethod() method.Method

	// GetUri returns the URI of this request.
	GetUri() uri.Uri

	ToMessage(includeStream bool) string

	// Version sets the HTTP version of this message.
	Version(version string) Request

	// Header adds a header to this message.
	Header(name, value string) Request

	// Headers adds multiple headers to this message.
	Headers(headers Headers) Request

	// ReplaceHeader replaces a header in this message.
	ReplaceHeader(name string, value string) Request

	// ReplaceHeaders replaces all headers in this message with the given source headers.
	ReplaceHeaders(source Headers) Request

	// RemoveHeader removes a header from this message.
	RemoveHeader(name string) Request

	// RemoveHeaders removes headers with the given prefix from this message.
	RemoveHeaders(prefix string) Request

	// Body sets the body of this message.
	Body(body body.Body) Request

	BodyReader(r io.ReadCloser) Request

	// BodyString sets the body of this message from a string.
	BodyString(body string) Request

	BodyJson(v any) (Request, error)

	// GetHeaders returns all headers of this message.
	GetHeaders() Headers

	// GetHeader returns the first value of the given header name.
	GetHeader(name string) string

	// GetHeaderValues returns all values of the given header name.
	GetHeaderValues(name string) []string

	// GetBody returns the body of this message as an io.ReadCloser.
	GetBody() body.Body

	// Close closes the body of this message.
	Close() error
}

type MemoryRequest struct {
	memoryMessage
	method method.Method
	uri    uri.Uri
}

func NewRequest() Request {
	return MemoryRequest{
		memoryMessage: memoryMessage{
			version: "HTTP/1.1",
		},
	}
}

func (r MemoryRequest) Version(version string) Request {
	r.memoryMessage = r.memoryMessage.Version(version)
	return r
}

func (r MemoryRequest) Header(name, value string) Request {
	r.memoryMessage = r.memoryMessage.Header(name, value)
	return r
}

func (r MemoryRequest) Headers(headers Headers) Request {
	r.memoryMessage = r.memoryMessage.Headers(headers)
	return r
}

func (r MemoryRequest) ReplaceHeader(name string, value string) Request {
	r.memoryMessage = r.memoryMessage.ReplaceHeader(name, value)
	return r
}

func (r MemoryRequest) ReplaceHeaders(source Headers) Request {
	r.memoryMessage = r.memoryMessage.ReplaceHeaders(source)
	return r
}

func (r MemoryRequest) RemoveHeader(name string) Request {
	r.memoryMessage = r.memoryMessage.RemoveHeader(name)
	return r
}

func (r MemoryRequest) RemoveHeaders(prefix string) Request {
	r.memoryMessage = r.memoryMessage.RemoveHeaders(prefix)
	return r
}

func (r MemoryRequest) Body(body body.Body) Request {
	r.memoryMessage = r.memoryMessage.Body(body)
	return r
}

func (r MemoryRequest) BodyReader(reader io.ReadCloser) Request {
	r.memoryMessage = r.memoryMessage.BodyReader(reader)
	return r
}

func (r MemoryRequest) BodyString(body string) Request {
	r.memoryMessage = r.memoryMessage.BodyString(body)
	return r
}

func (r MemoryRequest) BodyJson(v any) (Request, error) {
	mm, err := r.memoryMessage.BodyJson(v)
	if err != nil {
		return r, err
	}
	r.memoryMessage = mm
	return r, nil
}

func (r MemoryRequest) GetMethod() method.Method {
	return r.method
}

func (r MemoryRequest) GetUri() uri.Uri {
	return r.uri
}

func (r MemoryRequest) Method(method method.Method) Request {
	r.method = method
	return r
}

func (r MemoryRequest) Uri(uri uri.Uri) Request {
	r.uri = uri
	return r
}

func (r MemoryRequest) Query(name, value string) Request {
	r.uri = r.uri.Query(r.uri.GetQuery() + "&" + name + "=" + value)
	return r
}

func (r MemoryRequest) ToMessage(includeStream bool) string {
	return ">>>>>>>> REQUEST\n" + r.method + " " + r.uri.GetFullPath() + " " + r.version + "\n" + r.memoryMessage.ToMessage(includeStream)
}
