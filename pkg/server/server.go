// The server package contains servers implemenations
package server

import (
	"net"

	log "github.com/sirupsen/logrus"
)

type Handler func(net.Conn)

type Server struct {
	socketType string
	addr       string
	listener   net.Listener
	handler    Handler

	logger *log.Entry
}

func New(serverType string, addr string, handler Handler) (*Server, error) {
	l, err := net.Listen(serverType, addr)

	if err != nil {
		return nil, err
	}

	logger := log.WithField("context", "server")

	return &Server{serverType, addr, l, handler, logger}, nil
}

// Run starts accepting and handling the connections
func (s *Server) Run() {
	for {
		c, err := s.listener.Accept()
		if err != nil {
			s.logger.Errorf("Failed to accept connection: %v", err)
			continue
		}
		go s.handler(c)
	}
}
