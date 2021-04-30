ARG K6_VERSION='v0.31.1'
ARG GO_VERSION='1.16'

FROM golang:${GO_VERSION}-alpine as builder
RUN go get -u github.com/k6io/xk6/cmd/xk6
WORKDIR /go/bin
RUN xk6 build ${K6_VERSION} --with github.com/skryukov/xk6-redis

FROM alpine:3.13
RUN apk add --no-cache ca-certificates && \
    adduser -D -u 12345 -g 12345 k6
COPY --from=builder /go/bin/k6 /usr/bin/k6

USER 12345
WORKDIR /home/k6
ENTRYPOINT ["k6"]
