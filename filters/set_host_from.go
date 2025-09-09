package filters

import (
	"github.com/joscha-alisch/http4go/http"
	"github.com/joscha-alisch/http4go/http/uri"
)

func SetHostFrom(u uri.Uri) http.Filter {

}

/* fun SetHostFrom(uri: Uri): Filter = Filter { next ->
{
next(it.uri(it.uri.scheme(uri.scheme).host(uri.host).port(uri.port))
.replaceHeader("Host", "${uri.host}${uri.port?.let { port -> ":$port" } ?: ""}"))
}
}*/
