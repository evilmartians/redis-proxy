// The server package contains servers implemenations
package multiproxy

import (
	"bytes"
	"os"
	"testing"

	"github.com/moonactive/ma-redis-proxy/pkg/redis"
	"github.com/moonactive/ma-redis-proxy/test/mocks"
	"github.com/stretchr/testify/assert"
)

var (
	redisMock *mocks.RedisClient
	proxy     *Proxy
)

func init() {
	redisMock = &mocks.RedisClient{}
	var err error
	proxy, err = New()

	if err != nil {
		panic(err)
	}

	if _, ok := os.LookupEnv("REAL_REDIS"); ok {
		client, _ := redis.Connect("42")
		proxy.rdb = client
	} else {
		proxy.rdb = redisMock
	}
}

func TestHandleCommandFastlane(t *testing.T) {
	var b bytes.Buffer
	session := proxy.NewSession(&b)

	t.Run("PING", func(t *testing.T) {
		b.Write([]byte("PING\r\n"))

		err := session.HandleCommand()
		assert.NoError(t, err)

		assert.Equal(t, "+PONG\r\n", b.String())
	})
}

func TestHandleCommandWhenHandshakeState(t *testing.T) {
	t.Run("SELECT", func(t *testing.T) {
		var b bytes.Buffer
		session := proxy.NewSession(&b)

		b.Write([]byte("SELECT 42\r\n"))

		err := session.HandleCommand()
		assert.NoError(t, err)

		assert.Equal(t, "+OK\r\n", b.String())
	})

	t.Run("Other commands", func(t *testing.T) {
		var b bytes.Buffer
		session := proxy.NewSession(&b)

		b.Write([]byte("GET x\r\n"))

		err := session.HandleCommand()
		assert.Error(t, err)

		assert.Equal(t, "-ERR command is called before `select`\r\n", b.String())
	})
}

func TestHandleCommandWhenRegularState(t *testing.T) {
	var b bytes.Buffer
	session := proxy.NewSession(&b)
	session.rdb = proxy.rdb

	session.state = regularState{}

	t.Run("SET a 1", func(t *testing.T) {
		b.Reset()

		redisMock.On("Execute",
			&redis.Command{Name: "SET", Args: []interface{}{"a", "1"}, Last: true},
		).Return([]byte("+OK\r\n"), nil)

		b.Write([]byte("SET a 1\r\n"))

		err := session.HandleCommand()
		assert.NoError(t, err)

		assert.Equal(t, "+OK\r\n", b.String())
	})

	t.Run("GET", func(t *testing.T) {
		b.Reset()

		redisMock.On("Execute",
			&redis.Command{Name: "GET", Args: []interface{}{"a"}, Last: true},
		).Return([]byte("$1\r\n1\r\n"), nil)

		b.Write([]byte("GET a\r\n"))

		err := session.HandleCommand()
		assert.NoError(t, err)

		assert.Equal(t, "$1\r\n1\r\n", b.String())
	})
}
