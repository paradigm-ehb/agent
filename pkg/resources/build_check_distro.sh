#!/bin/sh

set -xeu

clang configure.cpp -o -Iarena -Ibase check
