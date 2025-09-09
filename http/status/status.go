package status

type Status struct {
	Code int
	Text string
}

var (
	Continue           = Status{100, "Continue"}
	SwitchingProtocols = Status{101, "Switching Protocols"}

	Ok                          = Status{200, "OK"}
	Created                     = Status{201, "Created"}
	Accepted                    = Status{202, "Accepted"}
	NonAuthoritativeInformation = Status{203, "Non-Authoritative Information"}
	NoContent                   = Status{204, "No Content"}
	ResetContent                = Status{205, "Reset Content"}
	PartialContent              = Status{206, "Partial Content"}

	MultipleChoices   = Status{300, "Multiple Choices"}
	MovedPermanently  = Status{301, "Moved Permanently"}
	Found             = Status{302, "Found"}
	SeeOther          = Status{303, "See Other"}
	NotModified       = Status{304, "Not Modified"}
	UseProxy          = Status{305, "Use Proxy"}
	TemporaryRedirect = Status{307, "Temporary Redirect"}
	PermanentRedirect = Status{308, "Permanent Redirect"}

	BadRequest                   = Status{400, "Bad Request"}
	Unauthorized                 = Status{401, "Unauthorized"}
	PaymentRequired              = Status{402, "Payment Required"}
	Forbidden                    = Status{403, "Forbidden"}
	NotFound                     = Status{404, "Not Found"}
	MethodNotAllowed             = Status{405, "Method Not Allowed"}
	NotAcceptable                = Status{406, "Not Acceptable"}
	ProxyAuthRequired            = Status{407, "Proxy Authentication Required"}
	RequestTimeout               = Status{408, "Request Timeout"}
	Conflict                     = Status{409, "Conflict"}
	Gone                         = Status{410, "Gone"}
	LengthRequired               = Status{411, "Length Required"}
	PreconditionFailed           = Status{412, "Precondition Failed"}
	RequestEntityTooLarge        = Status{413, "Request Entity Too Large"}
	RequestUriTooLong            = Status{414, "Request-URI Too Long"}
	UnsupportedMediaType         = Status{415, "Unsupported Media Type"}
	RequestedRangeNotSatisfiable = Status{416, "Requested Range Not Satisfiable"}
	ExpectationFailed            = Status{417, "Expectation Failed"}
	ImATeapot                    = Status{418, "I'm a teapot"} //RFC2324
	MisdirectedRequest           = Status{421, "Misdirected Request"}
	UnprocessableEntity          = Status{422, "Unprocessable Entity"}
	Locked                       = Status{423, "Locked"}
	FailedDependency             = Status{424, "Failed Dependency"}
	TooEarly                     = Status{425, "Too Early"}
	UpgradeRequired              = Status{426, "Upgrade Required"}
	PreconditionRequired         = Status{428, "Precondition Required"}
	TooManyRequests              = Status{429, "Too many requests"}
	RequestHeaderFieldsTooLarge  = Status{431, "Request Header Fields Too Large"}
	UnavailableForLegalReasons   = Status{451, "Unavailable For Legal Reasons"}

	InternalServerError     = Status{500, "Internal Server Error"}
	NotImplemented          = Status{501, "Not Implemented"}
	BadGateway              = Status{502, "Bad Gateway"}
	ServiceUnavailable      = Status{503, "Service Unavailable"}
	ConnectionRefused       = Status{503, "Connection Refused"}
	UnknownHost             = Status{503, "Unknown Host"}
	GatewayTimeout          = Status{504, "Gateway Timeout"}
	HttpVersionNotSupported = Status{505, "HTTP Version Not Supported"}
)
