#!/bin/sh
# echo -------------------------------------------------------------------------
# echo Description: build go binary with flags 
# echo -------------------------------------------------------------------------

VERSION=1.0

prod() {
    go build \
        -ldflags="-s -w -X main.version=${VERSION}" \
        -o agent \
        ./cmd/agent/agent.go
}

debug() {
    go build \
        -race \
        -o debug \
        ./cmd/agent/agent.go
}

case "$1" in
    prod)
        prod
        ;;
    debug)
        debug
        ;;
    *)
        echo "Usage: [debug|prod]"
        exit 1
        ;;
esac
