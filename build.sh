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
SRCDIR="."
OUTDIR="gen/"

cd "$SRCDIR/proto"

for DIR in */; do
    DIR=${DIR%/}   

	PROTOS=$(find "$DIR" -type f -name "*.proto") 

# check length of protos
   if [ -z "$PROTOS" ]; then
        echo "tsjoe nothing in here ,  $DIR :("
        exit 1
    fi

    protoc --go_out=../$OUTDIR --go_opt=paths=source_relative \
           --go-grpc_out=../$OUTDIR --go-grpc_opt=paths=source_relative \
           $PROTOS

	echo "yahoooo yippie yay $PROTOS"

done
