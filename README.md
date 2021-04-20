# Moon Active Redis Proxy

Redis Proxy aims to reduce the number for concurrent connections to Redis databases from the monolith. Instead of initializing connections to Redis directly, clients connect to the proxy, which multiplex commands (using a small, limited number of _real_ Redis connections).

The proxy is designed to run a sidecar service on a EC2 machine (which runs dozens of app instances).

## Installation

- Make sure Go 1.6 installed.
- Clone repository.
- Run `make install`.

## Usage

Run proxy with the configuration passed as a JSON file:

```shell
$ ma-redis-proxy -c path/to/config.json

=> INFO[2021-04-20T14:25:56+03:00] Starting MoonActive Redis Proxy v0.1.0 (pid: 68626)
```

You can also provide configuration parameters through the corresponding environment variables (i.e. `MA_REDIS_PROXY_LOG_LEVEL=debug`, etc).

For more information about available options run `ma-redis-proxy -h`.

### Configuration format

Configuration file describes how to connect to Redis databases (credentials, custom settings):

```json
// TODO: Example
```

**IMPORTANT:** Each database is identified by the uniq **number**. To connect to a database via proxy, a client MUST used this number as a `db` part of the url, e.g., `redis://localhost:6379/42` means connecting to the database with the identifier `"42"`.

See more details in the [docs](docs/architecture.md).

## Development

**NOTE:** Make sure Go 1.6 installed.

The following commands are available:

```shell
# Build the Go binary (will be available in dist/ma-redis-proxy)
make

# Run Golang tests
make test
```

We use [golangci-lint](https://golangci-lint.run) to lint Go source code:

```sh
make lint
```

We use [gosec](https://github.com/securego/gosec) to scan for possible security issues:

```sh
make sec
```

Other commands:

```sh
# Run vet and vet shadow
make vet

# Run fmt
make fmt
```

**NOTE**: If the commands above fails on MacBook with Apple silicon, try running with `arch -x86_64` prefix.

### Links & Resources

Here is the list of relevant resources and tools used to build this app:

- [Redis Protocol spec](https://redis.io/topics/protocol).
- [Redis Cluster spec](https://redis.io/topics/cluster-spec).
- [logrus](https://github.com/sirupsen/logrus): structured, pluggable logging.
- [confita](https://github.com/heetch/confita): configuration management.
- [mockery](https://github.com/vektra/mockery): test mocks.
