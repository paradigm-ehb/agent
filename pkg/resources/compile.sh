#!/bin/bash
set -eu

CC=clang
SRC=resources_main.c
OUT_DIR=build
OUT_BIN=$OUT_DIR/agent_resources

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
-D_POSIX_C_SOURCE=200809L
"

sep() {
  printf '%s\n' "============================================================"
}

header() {
  sep
  printf '  %s\n' "$1"
  sep
}

info() {
  printf '  • %s\n' "$1"
}

ok() {
  printf '  ✓ %s\n' "$1"
}

mkdir -p "$OUT_DIR"

header "BUILD"
info "Compiler : $CC"
info "Source   : $SRC"
info "Output   : $OUT_BIN"

printf '\n'
info "Compiling…"
$CC $CFLAGS "$SRC" -o "$OUT_BIN"
ok "Compilation finished"

printf '\n'
header "RUN"
info "Executing binary"
sep
"$OUT_BIN"
sep
ok "Execution finished"
