#!/bin/sh

set -eu

echo "-------------------------------------------------------------------------"
echo "Description: Script to download protoc compiler (x86_64 / arm64) and compile proto files"
echo "-------------------------------------------------------------------------"

PROTOC_VERSION="30.2"
PB_REL="https://github.com/protocolbuffers/protobuf/releases"
INSTALL_DIR="$HOME/.local"
BIN_DIR="$INSTALL_DIR/bin"

ARCH="$(uname -m)"

case "$ARCH" in
    x86_64)
        PROTOC_ARCH="linux-x86_64"
        ;;
    aarch64)
        PROTOC_ARCH="linux-aarch_64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

ZIP_NAME="protoc-${PROTOC_VERSION}-${PROTOC_ARCH}.zip"
DOWNLOAD_URL="$PB_REL/download/v${PROTOC_VERSION}/${ZIP_NAME}"

echo "Detected architecture: $ARCH"
echo "Downloading: $DOWNLOAD_URL"

mkdir -p "$INSTALL_DIR"
curl -L -o "$ZIP_NAME" "$DOWNLOAD_URL"
unzip -o "$ZIP_NAME" -d "$INSTALL_DIR"
rm -f "$ZIP_NAME"

export PATH="$PATH:$BIN_DIR"

# Install Go plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

export PATH="$PATH:$(go env GOPATH)/bin"

# Compile proto files
SRCDIR="."
OUTDIR="gen"

cd "$SRCDIR/proto"

for DIR in */; do
    DIR="${DIR%/}"
    PROTOS="$(find "$DIR" -type f -name '*.proto')"

    if [ -z "$PROTOS" ]; then
        echo "Skipping $DIR (no proto files)"
        continue
    fi

    protoc \
        --go_out="../$OUTDIR" --go_opt=paths=source_relative \
        --go-grpc_out="../$OUTDIR" --go-grpc_opt=paths=source_relative \
        $PROTOS

    echo "Compiled: $PROTOS"
done

