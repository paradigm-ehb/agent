#!/bin/sh

echo -------------------------------------------------------------------------
echo Description: Script to download protoc compiler and compile protoc files
echo -------------------------------------------------------------------------

# Install protocol buffer compiler
# source : https://protobuf.dev/installation/
PB_REL="https://github.com/protocolbuffers/protobuf/releases"
curl -LO $PB_REL/download/v30.2/protoc-30.2-linux-x86_64.zip
unzip protoc-30.2-linux-x86_64.zip -d $HOME/.local
export PATH="$PATH:$HOME/.local/bin"

# Install protoc compiler plugins for Go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Add protoc compiler to binary
export PATH="$PATH:$(go env GOPATH)/bin"

# compile protoc files
srcdir="."
outdir="internal/connection/pb"

if ls $srcdir/internal/connection/*proto; then
    protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        $srcdir/internal/connection/*.proto
else
    echo "Geen proto file gevonden :("
    exit 1
fi

