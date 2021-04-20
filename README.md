# Redis Proxy

## Installation

- Make sure Go 1.6+ installed.
- Clone repository.
- Run `make install`.

## Usage

Run proxy with the configuration passed as a JSON file:

```shell
$ redis-proxy -c path/to/config.json

=> INFO[2021-04-20T14:25:56+03:00] context=main Starting Redis Proxy v0.1.0 (pid: 68626)
```

You can also provide configuration parameters through the corresponding environment variables (i.e. `REDIS_PROXY_LOG_LEVEL=debug`, etc).

For more information about available options run `redis-proxy -h`.

## Development

**NOTE:** Make sure Go 1.6+ installed.

First, install the required dev tools:

```shell
# Installs golangci-lint, vet shadow, gosec
make dev-prepare
```

If this command fails on MacBook with Apple silicon, try running it with `arch -x86_64` prefix.

The following commands are available:

```shell
# Build the Go binary (will be available in dist/redis-proxy)
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
