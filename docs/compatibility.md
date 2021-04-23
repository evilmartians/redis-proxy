# Redis compatibility

This document contains the information on commands which are currently not supported or require special attention.

## PING

Ping doesn't call Redis, it returns PONG immediately is the corresponding Redis connection is alive.

## MGET/MSET

Only the first key is used to determine the cluster slot.
See [go-redis#527](https://github.com/go-redis/redis/issues/527) (ioredis [behaves similarly](https://github.com/luin/ioredis/issues/1128)).

## ⚠️ Unsupported commands ⚠️

🚧 _To be completed_

- All blocking commands (including pub/sub and streams)
- `CLUSTER`
- `DISCARD`
- `WATCH`
