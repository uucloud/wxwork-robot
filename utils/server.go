package utils

import "net/http"

type Server struct {
	sc     *ServerConfig
	Server *http.Server
}

type ServerConfig struct {
	Addr string
}

func NewServer(config *ServerConfig, router http.Handler) *Server {
	eng := &Server{
		sc: config,
		Server: &http.Server{
			Addr:    config.Addr,
			Handler: router,
		},
	}

	return eng
}
