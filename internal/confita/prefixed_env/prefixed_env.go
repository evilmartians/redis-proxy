package prefixed_env // nolint:stylecheck

import (
	"context"
	"os"
	"strings"

	"github.com/heetch/confita/backend"
)

// NewBackend creates a configuration loader that loads from the environment and
// allows specifying the env keys prefix.
// Based on https://github.com/heetch/confita/blob/master/backend/env/env.go.
func NewBackend(prefix string) backend.Backend {
	return backend.Func("prefixed_env", func(ctx context.Context, key string) ([]byte, error) {
		key = strings.Join([]string{prefix, key}, "_")

		if val := os.Getenv(key); val != "" {
			return []byte(val), nil
		}
		key = strings.ReplaceAll(strings.ToUpper(key), "-", "_")
		if val := os.Getenv(key); val != "" {
			return []byte(val), nil
		}
		return nil, backend.ErrNotFound
	})
}
