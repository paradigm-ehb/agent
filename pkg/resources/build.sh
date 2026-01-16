#!/bin/sh
set -eu

CC=cc
AR=ar

SRC=resources.c
OUT_DIR=build
OUT_OBJ=$OUT_DIR/resources.o
OUT_LIB=$OUT_DIR/libagent_resources.a

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

# TODO(nasr): compile the package manager libs, static bin
mkdir -p "$OUT_DIR"

echo "Compiling object..."
$CC $CFLAGS -c "$SRC" -o "$OUT_OBJ"

echo "Creating static library..."
$AR rcs "$OUT_LIB" "$OUT_OBJ"

echo "Done:"
echo "  $OUT_LIB"
