package multiproxy

import "github.com/moonactive/ma-redis-proxy/pkg/redis"

var (
	fastlane = map[string]func(s *Session, cmd *redis.Command) error{
		"PING": func(s *Session, cmd *redis.Command) error {
			err := s.writer.WriteSimpleString("PONG")
			if err != nil {
				return nil
			}

			err = s.writer.Flush()
			if err != nil {
				return nil
			}

			return nil
		},
	}
)
