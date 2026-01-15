#!/bin/bash
set -eu

# Compile C files with clang (C compiler)
clang -c arena.c -o arena.o
clang -c resources.c -o resources.o

# Compile C++ files and link with C objects
clang configure.c arena.o resources.o -ldnf5 -o nob

./nob
