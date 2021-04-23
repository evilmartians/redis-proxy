// The redis package includes the Redis clients and pools code
package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// Command is an application-specific Redis command representation
type Command struct {
	// Name is the first argument of the command
	Name string
	// Args is the optional list of arguments (1..)
	Args []interface{}
	// Last is true iff this command is the last in the bulk payload
	Last bool
}

type RedisClient interface {
	Execute(cmd *Command) ([]byte, error)
}

type Client struct {
	id  string
	rdb *redis.Client
}

func Connect(id string) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &Client{id: id, rdb: rdb}, nil
}

func (c *Client) Execute(cmd *Command) ([]byte, error) {
	args := append([]interface{}{cmd.Name}, cmd.Args...)

	ctx := context.Background()
	redisCmd := redis.NewCmdWithRawResponse(ctx, args...)
	_ = c.rdb.Process(ctx, redisCmd)

	return redisCmd.RawResponse(), nil
}
