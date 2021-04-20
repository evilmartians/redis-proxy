package prefixed_env

import (
	"context"
	"os"
	"testing"

	"github.com/heetch/confita/backend"
	"github.com/stretchr/testify/require"
)

func TestPrefixedEnvBackend(t *testing.T) {
	t.Run("NotFoundBecauseUnset", func(t *testing.T) {
		b := NewBackend("TEST")

		_, err := b.Get(context.Background(), "something that doesn't exist")
		require.Equal(t, backend.ErrNotFound, err)
	})

	t.Run("NotFoundBecauseEmpty", func(t *testing.T) {
		b := NewBackend("TEST")
		os.Setenv("TEST_CONFIG", "")
		_, err := b.Get(context.Background(), "CONFIG")
		require.Equal(t, backend.ErrNotFound, err)
	})

	t.Run("ExactMatch", func(t *testing.T) {
		b := NewBackend("TEST")

		os.Setenv("TEST_CONFIG1", "ok")
		val, err := b.Get(context.Background(), "CONFIG1")
		require.NoError(t, err)
		require.Equal(t, "ok", string(val))
	})

	t.Run("DifferentCase", func(t *testing.T) {
		b := NewBackend("TEST")

		os.Setenv("TEST_CONFIG_2", "ok")
		val, err := b.Get(context.Background(), "config-2")
		require.NoError(t, err)
		require.Equal(t, "ok", string(val))
	})
}
