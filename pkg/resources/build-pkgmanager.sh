#!/bin/sh
set -eu

CXX=c++
SRC=libdnf_handler.cpp
OUT=build/pkgmanager

CXXFLAGS="
-std=c++20
-Wall -Wextra -Wpedantic
-Wshadow -Wconversion
-Wundef
-Wpointer-arith -Wcast-align -Wcast-qual
-Wformat=2 -Wformat-security
-Wnull-dereference
-Wmisleading-indentation
-Wunused -Wuninitialized
-fexceptions
-fno-rtti
-ldnf5
-lfmt
"

mkdir -p build

echo "Compiling..."
$CXX $CXXFLAGS "$SRC" -o "$OUT"

echo "Running..."
"$OUT"
