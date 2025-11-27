#!/bin/sh

srcdir="."
outdir="internal/connection/pb"

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative "$srcdir/internal/connection/greet.proto"

