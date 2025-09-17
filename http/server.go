package http

type Server interface {
	Start() Server
	Stop() Server
	StartBlocking() error
}

type ServerConfig interface {
	ToServer(Handler) Server
}

type PolyServerConfig interface {
	ServerConfig
	ToPolyServer(Handler) Server
}
