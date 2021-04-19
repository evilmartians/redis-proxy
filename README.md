# Moon Active Redis Proxy

## Installation

- Make sure Go 1.6 installed.
- Clone repository.
- Run `make install`.

## Usage

Run proxy with the configuration passed as a JSON file:

```shell
$ ma-redis-proxy -c path/to/config.json

=> INFO 2020-02-05T08:44:57.684Z context=main Starting Moon Active Redis Proxy v0.1.0
```

You can also provide configuration parameters through the corresponding environment variables (i.e. `MA_REDIS_PROXY_CONFIG`, etc).

For more information about available options run `ma-redis-proxy -h`.

## Development

**NOTE:** Make sure Go 1.6 installed.

First, install the required dev tools:

```shell
# Installs golangci-lint, vet shadow, gosec
make dev-prepare
```

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
