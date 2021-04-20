// The cli package contains the code required to run redis-proxy as CLI
package cli

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/evilmartians/redis-proxy/internal/confita/prefixed_env"
	"github.com/evilmartians/redis-proxy/pkg/config"
	"github.com/evilmartians/redis-proxy/pkg/multiproxy"
	"github.com/evilmartians/redis-proxy/pkg/server"
	"github.com/evilmartians/redis-proxy/pkg/version"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/flags"
	"github.com/syossan27/tebata"

	log "github.com/sirupsen/logrus"
)

// Run configures and initializes the application,
// starts a server and setups signal traps.
func Run() error {
	var err error

	conf, err := loadConfig()

	if err != nil {
		return err
	}

	if err = initLogger(conf); err != nil {
		return err
	}

	logger := log.WithField("context", "main")

	logger.Infof("Starting Redis Proxy v%s (pid: %d)", version.Version(), os.Getpid())

	var proxy *multiproxy.Proxy

	if proxy, err = multiproxy.New(); err != nil {
		return err
	}

	if err = proxy.Boot(); err != nil {
		return err
	}

	s, err := server.New(conf.Addr, func(c net.Conn) {
		defer c.Close()

		sessionLogger := log.WithField("context", "server")

		var perr error

		sessionLogger.Debugf("Client connection %s", c.RemoteAddr())

		session := proxy.NewSession(c)

		if perr = session.HandleCommands(); perr != nil {
			if perr == io.EOF {
				sessionLogger.Infof("Client disconnected")
			} else {
				sessionLogger.Errorf("Failed to handle client commnands: %v", perr)
			}
			return
		}
	})

	if err != nil {
		return err
	}

	go s.Run()

	logger.Infof("Handle Redis connections at %s://%s", s.Type(), s.Addr())

	t := tebata.New(syscall.SIGINT, syscall.SIGTERM)

	t.Reserve(func() { // nolint:errcheck
		logger.Infof("Shutting down... (hit Ctrl-C to stop immediately)")
		go func() {
			termSig := make(chan os.Signal, 1)
			signal.Notify(termSig, syscall.SIGINT, syscall.SIGTERM)
			<-termSig
			logger.Warnf("Immediate termination requested. Stopped")
			os.Exit(0)
		}()
	})
	t.Reserve(s.Shutdown)     // nolint:errcheck
	t.Reserve(proxy.Shutdown) // nolint:errcheck

	t.Reserve(os.Exit, 0) // nolint:errcheck

	select {}
}

func loadConfig() (config.Config, error) {
	cfg := config.New()

	loader := confita.NewLoader(
		prefixed_env.NewBackend("REDIS_PROXY"),
		flags.NewBackend(),
	)

	err := loader.Load(context.Background(), &cfg)

	return cfg, err
}

func initLogger(conf config.Config) error {
	level, err := log.ParseLevel(conf.LogLevel)

	if err != nil {
		return err
	}

	log.SetLevel(level)

	if conf.LogFile != "" {
		f, err := os.OpenFile(conf.LogFile, os.O_WRONLY|os.O_CREATE, 0755) // #nosec
		if err != nil {
			return err
		}
		log.SetOutput(f)
	} else {
		log.SetOutput(os.Stdout)
	}

	if conf.LogFormat == "text" {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	} else if conf.LogFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		return fmt.Errorf("Unknown log formatter: %s", conf.LogFormat)
	}

	return nil
}
