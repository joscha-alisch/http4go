package uri

import (
	"fmt"
	"regexp"
)

var authorityRegex = regexp.MustCompile("(?:([^@]+)@)?([^:]+)(?::([\\d]+))?")
var rfc3986Regex = regexp.MustCompile("^(?:([^:/?#]+):)?(?://([^/?#]*))?([^?#]*)(?:\\?([^#]*))?(?:#(.*))?")

type Uri struct {
	scheme   string
	host     string
	port     int
	path     string
	query    string
	fragment string
}

func NewUri() Uri {
	return Uri{}
}

func (u Uri) Scheme(scheme string) Uri {
	u.scheme = scheme
	return u
}

func (u Uri) Host(host string) Uri {
	u.host = host
	return u
}

func (u Uri) Port(port int) Uri {
	u.port = port
	return u
}

func (u Uri) Path(path string) Uri {
	u.path = path
	return u
}

func (u Uri) Query(query string) Uri {
	u.query = query
	return u
}

func (u Uri) Fragment(fragment string) Uri {
	u.fragment = fragment
	return u
}

func (u Uri) GetScheme() string {
	return u.scheme
}

func (u Uri) GetHost() string {
	return u.host
}

func (u Uri) GetPort() int {
	return u.port
}

func (u Uri) GetPath() string {
	return u.path
}

func (u Uri) GetQuery() string {
	return u.query
}

func (u Uri) GetFragment() string {
	return u.fragment
}

func (u Uri) GetHostPort() string {
	if u.port != 0 {
		return fmt.Sprintf("%s:%d", u.host, u.port)
	}
	return u.host
}

func (u Uri) GetFullPath() string {
	result := u.path
	if u.query != "" {
		result += "?" + u.query
	}

	if u.fragment != "" {
		result += "#" + u.fragment
	}
	return result
}
func (u Uri) String() string {
	return u.scheme + "://" + u.GetHostPort() + u.GetFullPath()
}
