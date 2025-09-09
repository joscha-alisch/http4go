package uri

import "regexp"

var authorityRegex = regexp.MustCompile("(?:([^@]+)@)?([^:]+)(?::([\\d]+))?")
var rfc3986Regex = regexp.MustCompile("^(?:([^:/?#]+):)?(?://([^/?#]*))?([^?#]*)(?:\\?([^#]*))?(?:#(.*))?")

type Uri struct {
	Scheme   string
	Host     string
	Port     int
	Path     string
	Query    string
	Fragment string
}

func Of(uri string) (Uri, error) {
	groups := rfc3986Regex.FindStringSubmatch(uri)
	authority := groups[2]
	authorityGroups := authorityRegex.FindStringSubmatch(authority)

	return Uri{
		Scheme: groups[1],
		Host:   authorityGroups[2],
	}, nil
}
