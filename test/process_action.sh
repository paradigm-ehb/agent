#!/bin/sh

set -eu

grpcurl -plaintext \
  -d '{
    "pid": 2497,
    "signal": 15
  }' \
  localhost:5000 \
  resources.v2.ResourcesService/ProcessAction
