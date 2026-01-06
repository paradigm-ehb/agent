#!/usr/bin/env sh
set -eu

HOST="${HOST:-localhost}"
PORT="${PORT:-5000}"
ADDR="${HOST}:${PORT}"

echo "== gRPC integration test =="
echo "Target: ${ADDR}"
echo

command -v grpcurl >/dev/null 2>&1 || {
  echo "ERROR: grpcurl not found"
  exit 1
}

echo "== Checking server availability =="
grpcurl -plaintext "${ADDR}" list >/dev/null
echo "OK"
echo

echo "== Listing services =="
grpcurl -plaintext "${ADDR}" list
echo

echo "== Checking JournalService =="
grpcurl -plaintext "${ADDR}" list journal.v1.JournalService >/dev/null
echo "OK"
echo

echo "== JournalService.Action (GID=1000) =="
grpcurl -plaintext \
  -d '{
    "numFromTail": 5,
    "field": "GID",
    "value": "1000"
  }' \
  "${ADDR}" \
  journal.v1.JournalService/Action
echo

echo "== JournalService.Action (systemd unit) =="
grpcurl -plaintext \
  -d '{
    "numFromTail": 5,
    "field": "Systemd",
    "value": "systemd-journald.service"
  }' \
  "${ADDR}" \
  journal.v1.JournalService/Action
echo

echo "== Checking ResourcesService =="
grpcurl -plaintext "${ADDR}" list resources.v1.ResourcesService >/dev/null
echo "OK"
echo

echo "== ResourcesService.GetSystemResources =="
grpcurl -plaintext \
  -d '{}' \
  "${ADDR}" \
  resources.v1.ResourcesService/GetSystemResources
echo

echo "== All gRPC tests passed =="
