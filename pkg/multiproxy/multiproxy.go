// The multiproxy (multiplexer + proxy) package contains the actual proxying logic: mapping
// sessions to Redis pools, handling "special" commands (SCRIPT LOAD, MULTI, etc.).
package multiproxy

import (
	"io"

	log "github.com/sirupsen/logrus"

	"github.com/evilmartians/redis-proxy/pkg/redis"
)

type Proxy struct {
	logger *log.Entry

	// TEMP
	rdb redis.RedisClient
}

func New() (*Proxy, error) {
	logger := log.WithField("context", "proxy")

	return &Proxy{logger: logger}, nil
}

// Boot initializes Redis clients for all databases
// (except those marked as lazy)
func (p *Proxy) Boot() error {
	// TODO: to be implemented
	rdb, err := redis.Connect("temp")

	if err != nil {
		return err
	}

	p.rdb = rdb

	p.logger.Info("Successfully connected to databases")
	return nil
}

// Shutdown disconnects all the Redis pool (essentially causing clients to close their connections, too).
// Each pool MUST wait for active calls to complete before closing.
func (p *Proxy) Shutdown() {
	// TODO: to be implemented
}

// NewSession creates a new session struct.
func (p *Proxy) NewSession(io io.ReadWriter) *Session {
	return NewSession(io, p)
}

// LookupClient returns a Redis client corresponding to the specified name
func (p *Proxy) LookupClient(dbname string) redis.RedisClient {
	return p.rdb
}
