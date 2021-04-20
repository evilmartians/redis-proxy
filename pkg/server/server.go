// The server package contains servers implemenations
package server

import (
	"fmt"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Handler func(net.Conn)

type Server struct {
	protocol string
	addr     string
	listener net.Listener
	handler  Handler

	quit   chan (struct{})
	logger *log.Entry
}

var (
	protocolToSocketType = map[string]string{"tcp": "tcp4", "unix": "unix"}
)

func New(addr string, handler Handler) (*Server, error) {
	parts := strings.SplitN(addr, "://", 2)

	var protocol string
	var hostname string

	if len(parts) < 2 {
		protocol = "tcp"
		hostname = parts[0]
	} else {
		protocol = parts[0]
		hostname = parts[1]
	}

	socketType, ok := protocolToSocketType[protocol]

	if !ok {
		return nil, fmt.Errorf("Unsupported protocol: %s", protocol)
	}

	l, err := net.Listen(socketType, hostname)

	if err != nil {
		return nil, err
	}

	logger := log.WithField("context", "server")

	return &Server{
		protocol: protocol,
		addr:     addr,
		listener: l,
		handler:  handler,
		quit:     make(chan struct{}),
		logger:   logger,
	}, nil
}

// Run starts accepting and handling the connections
func (s *Server) Run() {
	for {
		c, err := s.listener.Accept()
		if err != nil {
			select {
			// Handle shutdown
			case <-s.quit:
				return
			default:
				s.logger.Errorf("Failed to accept connection: %v", err)
			}
			continue
		}
		go s.handler(c)
	}
}

// Shutdown stops listening for new connections
func (s *Server) Shutdown() {
	// Stop accepting new connections
	close(s.quit)
	s.listener.Close()
}

// Addr returns the actual server address (useful in case of 0 port)
func (s *Server) Addr() string {
	addr := s.listener.Addr()
	return addr.String()
}

// Type returns the server type (tcp or unix)
func (s *Server) Type() string {
	return s.protocol
}
