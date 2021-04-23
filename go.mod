module github.com/evilmartians/redis-proxy

go 1.16

replace github.com/go-redis/redis/v8 => github.com/skryukov/redis/v8 v8.0.0-20210421140411-f2a2650293dd

require (
	github.com/go-redis/redis/v8 v8.0.0-00010101000000-000000000000
	github.com/heetch/confita v0.10.0
	github.com/secmask/go-redisproto v0.1.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/syossan27/tebata v0.0.0-20180602121909-b283fe4bc5ba
)
