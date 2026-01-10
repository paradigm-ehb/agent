#!/bin/sh

go clean -cache

if CGO_ENABLED=1 go build -race -gcflags="all=-N -l" -o debug ./cmd/agent/agent.go; then
    exec ./debug
else
    exit 1
fi
