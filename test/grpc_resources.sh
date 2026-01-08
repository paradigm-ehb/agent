#!/usr/bin/env sh
set -eu

ADDR="${ADDR:-localhost:5000}"

echo "== ResourcesService integration test =="
echo "Target: $ADDR"
echo

command -v grpcurl >/dev/null 2>&1 || {
  echo "ERROR: grpcurl not found"
  exit 1
}

echo "== Checking service availability =="
grpcurl -plaintext "$ADDR" list resources.v1.ResourcesService >/dev/null
echo "OK"
echo

echo "== Calling GetSystemResources =="
grpcurl -plaintext \
  -format text \
  "$ADDR" \
  resources.v1.ResourcesService/GetSystemResources
echo

echo "== Test completed =="
