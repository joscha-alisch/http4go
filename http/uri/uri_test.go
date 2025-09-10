package uri

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestUriOf(t *testing.T) {
	tests := []struct {
		name        string
		uriStr      string
		expected    Uri
		expectedErr error
	}{
		{"scheme and host", "http://example.com", Uri{scheme: "http", host: "example.com"}, nil},
		{"port", "https://example.com:8080", Uri{scheme: "https", host: "example.com", port: 8080}, nil},
		{"path", "http://localhost/some/path", Uri{scheme: "http", host: "localhost", path: "/some/path"}, nil},
		{"query", "http://localhost?query=1", Uri{scheme: "http", host: "localhost", query: "query=1"}, nil},
		{"fragment", "http://localhost#section", Uri{scheme: "http", host: "localhost", fragment: "section"}, nil},
		{"userinfo", "http://user:password@localhost", Uri{scheme: "http", host: "localhost", userinfo: "user:password"}, nil},
		{"full uri", "https://user:password@localhost:8443/some/path?query=1#section", Uri{scheme: "https", host: "localhost", port: 8443, path: "/some/path", query: "query=1", fragment: "section", userinfo: "user:password"}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uri, err := Of(test.uriStr)
			if !errors.Is(err, test.expectedErr) {
				t.Errorf("expected error %v, got %v", test.expectedErr, err)
			}

			if uri != test.expected {
				t.Errorf("expected uri %s, got %s\n%s", test.expected.String(), uri.String(), cmp.Diff(test.expected, uri, cmp.AllowUnexported(Uri{})))
			}
		})
	}

}
