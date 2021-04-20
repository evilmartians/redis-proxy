// The cli package contains the code required to run ma-redis-proxy as CLI
package cli

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/flags"
	"github.com/moonactive/ma-redis-proxy/internal/confita/prefixed_env"
	"github.com/moonactive/ma-redis-proxy/pkg/config"
	"github.com/moonactive/ma-redis-proxy/pkg/multiproxy"
	"github.com/moonactive/ma-redis-proxy/pkg/server"
	"github.com/moonactive/ma-redis-proxy/pkg/version"

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

	logger.Infof("Starting MoonActive Redis Proxy v%s (pid: %d)", version.Version(), os.Getpid())

	var proxy *multiproxy.Proxy

	if proxy, err = multiproxy.New(); err != nil {
		return err
	}

	if err = proxy.Boot(); err != nil {
		return err
	}

	// TODO: Get protocol and hostname from config
	logger.Info("Handle Redis connections at tcp://127.0.0.1:4321")

	s, err := server.New("tcp4", "127.0.0.1:4321", func(c net.Conn) {
		defer c.Close()

		sessionLogger := log.WithField("context", "server")

		var perr error

		sessionLogger.Debugf("Client connection %s", c.RemoteAddr())

		session := proxy.NewSession(c)

		if perr = session.HandleCommands(); perr != nil {
			sessionLogger.Errorf("Failed to handle client commnands: %v", perr)
			return
		}
	})

	if err != nil {
		return err
	}

	go s.Run()

	select {}
}

func loadConfig() (config.Config, error) {
	cfg := config.New()

	loader := confita.NewLoader(
		prefixed_env.NewBackend("MA_REDIS_PROXY"),
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
