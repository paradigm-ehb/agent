#!/bin/sh
set -eu

echo "========================================================================="
echo "Build Script: Protobuf compilation + Agent Resources C library"
echo "========================================================================="

# ============================================================================
echo ""
echo "[1/2] Protobuf Compilation"
echo "-------------------------------------------------------------------------"

PROTOC_VERSION="30.2"
PB_REL="https://github.com/protocolbuffers/protobuf/releases"
INSTALL_DIR="$HOME/.local"
BIN_DIR="$INSTALL_DIR/bin"
PROTOC_BIN="$BIN_DIR/protoc"

# Check if protoc is already installed
if [ ! -f "$PROTOC_BIN" ]; then
    echo "protoc not found, downloading..."
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
    
    echo "Installing Go protobuf plugins..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
else
    echo "protoc already installed, skipping download"
fi

export PATH="$PATH:$BIN_DIR:$(go env GOPATH)/bin"

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

cd ..

# TOOD(nasr): call the external resources buildscript instead of doing this
echo ""
echo "[2/2] Agent Resources C Library"
echo "-------------------------------------------------------------------------"

CC=cc
AR=ar
AGENT_RES_DIR="pkg/agent-resources"
SRC="$AGENT_RES_DIR/resources.c"
OUT_DIR="$AGENT_RES_DIR/build"
OUT_OBJ="$OUT_DIR/resources.o"
OUT_LIB="$OUT_DIR/libagent_resources.a"

CFLAGS="
-std=c99
-Wall
-Wextra
-Wpedantic
-Wshadow
-Wconversion
-Wundef
-Wpointer-arith
-Wcast-align
-Wcast-qual
-Wwrite-strings
-Wformat=2
-Wformat-security
-Wnull-dereference
-Wmisleading-indentation
-Wunused
-Wuninitialized
-Werror
-Wdouble-promotion
-Wstrict-overflow=2
-D_POSIX_C_SOURCE=200809L
"

mkdir -p "$OUT_DIR"

echo "Compiling object..."
$CC $CFLAGS -c "$SRC" -o "$OUT_OBJ"

echo "Creating static library..."
$AR rcs "$OUT_LIB" "$OUT_OBJ"

echo "Done: $OUT_LIB"

echo ""
echo "========================================================================="
echo "Build Complete!"
echo "========================================================================="
echo "✓ Protobuf files compiled to: $OUTDIR"
echo "✓ C library built at: $OUT_LIB"
echo "========================================================================="
