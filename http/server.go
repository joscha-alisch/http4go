package http

type Server interface {
	Start() Server
	Stop() Server
	StartBlocking() error
}

type ServerConfig interface {
	ToServer(Handler) Server
}

type SseServerConfig interface {
	ToSseServer(SseHandler) Server
}

type PolyServerConfig interface {
	ServerConfig
	SseServerConfig
	ToPolyServer(Handler, SseHandler) Server
}
