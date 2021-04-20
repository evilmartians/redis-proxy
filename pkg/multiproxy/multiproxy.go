// The multiproxy (multiplexer + proxy) package contains the actual proxying logic: mapping
// sessions to Redis pools, handling "special" commands (SCRIPT LOAD, MULTI, etc.).
package multiproxy

import (
	"io"
)

type Proxy struct {
}

func New() (*Proxy, error) {
	return &Proxy{}, nil
}

// Boot initializes Redis clients for all databases
// (except those marked as lazy)
func (*Proxy) Boot() error {
	// TODO: to be implemented
	return nil
}

// NewSession creates a new session struct.
func (*Proxy) NewSession(io io.ReadWriter) *Session {
	return &Session{io: io}
}
