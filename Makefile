OUTPUT ?= dist/redis-proxy

export GO111MODULE=on

ifdef VERSION
	LD_FLAGS="-s -w -X github.com/evilmartians/redis-proxy/version.number=$(VERSION)"
else
	COMMIT := $(shell sh -c 'git log --pretty=format:"%h" -n 1 ')
	VERSION := $(shell sh -c 'git tag -l --sort=-version:refname "v*" | head -n1')
	LD_FLAGS="-s -w -X github.com/evilmartians/redis-proxy/version.sha=$(COMMIT) -X github.com/evilmartians/redis-proxy/version.number=$(VERSION)"
endif

GOBUILD=go build -ldflags $(LD_FLAGS) -a

# Standard build
default: build

# Install current version
install:
	go install ./...

build:
	go build -ldflags $(LD_FLAGS) -o $(OUTPUT) cmd/redis-proxy/main.go

build-clean:
	rm -rf ./dist

# Run server
run:
	@go run ./cmd/redis-proxy/main.go

test:
	go test -count=1 -timeout=30s -race ./...

bin/shadow:
	@which shadow &> /dev/null || \
		env GO111MODULE=off go get golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow

bin/golangci-lint:
	@test -x $$(go env GOPATH)/bin/golangci-lint || \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.39.0

bin/gosec:
	@test -x $$(go env GOPATH)/bin/gosec || \
		curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v2.7.0

vet: bin/shadow
	go vet ./...
	go vet -vettool=$$(which shadow) ./...

sec: bin/gosec
	$$(go env GOPATH)/bin/gosec -quiet -confidence=medium -severity=medium  ./...

fmt:
	go fmt ./...

lint: bin/golangci-lint
	$$(go env GOPATH)/bin/golangci-lint run

.PHONY: test

# Docker compose related commands

redis-cluster-up:
	docker-compose -f ./etc/docker-compose.redis.yml up redis-cluster

redis-up:
	docker-compose -f ./etc/docker-compose.redis.yml up redis

SHARDN ?= 1

stop-shard:
	docker-compose -f ./etc/docker-compose.redis.yml exec redis-cluster supervisorctl stop redis-$(SHARDN)

start-shard:
	docker-compose -f ./etc/docker-compose.redis.yml exec redis-cluster supervisorctl start redis-$(SHARDN)

HOST ?= host.docker.internal
PORT ?= 6379
DBNUM ?= 42

redis-benchmark:
	docker-compose -f ./etc/docker-compose.redis.yml run redis redis-benchmark -q -n 1000 -c 50 -r 50 -k 1 -h $(HOST) -p $(PORT) --dbnum $(DBNUM)
