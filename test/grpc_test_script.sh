#!/usr/bin/env sh
set -euo pipefail

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

echo "== Listing all services =="
grpcurl -plaintext "${ADDR}" list
echo

SERVICE="services.v2.HandlerService"

echo "== Checking ${SERVICE} existence =="
grpcurl -plaintext "${ADDR}" list "${SERVICE}" >/dev/null
echo "OK"
echo

echo "== Listing ${SERVICE} methods =="
grpcurl -plaintext "${ADDR}" list "${SERVICE}"
echo

echo "== GetAllUnits =="
grpcurl -plaintext \
  -d '{}' \
  "${ADDR}" \
  "${SERVICE}/GetAllUnits"
echo

echo "== GetLoadedUnits =="
grpcurl -plaintext \
  -d '{}' \
  "${ADDR}" \
  "${SERVICE}/GetLoadedUnits"
echo

echo "== GetUnitStatus (example: ssh.service) =="
grpcurl -plaintext \
  -d '{
    "unitName": "ssh.service"
  }' \
  "${ADDR}" \
  "${SERVICE}/GetUnitStatus"
echo

echo "== PerformAction: restart ssh.service =="
grpcurl -plaintext \
  -d '{
    "serviceName": "ssh.service",
    "unitAction": "UNIT_ACTION_RESTART"
  }' \
  "${ADDR}" \
  "${SERVICE}/PerformAction"
echo

echo "== PerformAction: enable ssh.service (runtime=false, force=true) =="
grpcurl -plaintext \
  -d '{
    "serviceName": "ssh.service",
    "unitFileAction": "UNIT_FILE_ACTION_ENABLE",
    "runtime": false,
    "force": true
  }' \
  "${ADDR}" \
  "${SERVICE}/PerformAction"
echo

echo "== All gRPC tests passed =="
