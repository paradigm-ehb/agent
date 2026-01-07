#!/usr/bin/env sh
set -euo pipefail

HOST="${HOST:-localhost}"
PORT="${PORT:-5000}"
ADDR="${HOST}:${PORT}"

SERVICE="services.v2.HandlerService"

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

echo "== GetUnitStatus (tailscaled.service) =="
grpcurl -plaintext \
  -d '{
    "unitName": "tailscaled.service"
  }' \
  "${ADDR}" \
  "${SERVICE}/GetUnitStatus"
echo

echo "== PerformUnitAction: START tailscaled.service =="
grpcurl -plaintext \
  -d '{
    "unitName": "tailscaled.service",
    "action": "UNIT_ACTION_START"
  }' \
  "${ADDR}" \
  "${SERVICE}/PerformUnitAction"
echo

echo "== PerformUnitFileAction: ENABLE tailscaled.service (runtime=true, force=true) =="
grpcurl -plaintext \
  -d '{
    "unitName": "tailscaled.service",
    "action": "UNIT_FILE_ACTION_ENABLE",
    "runtime": true,
    "force": true
  }' \
  "${ADDR}" \
  "${SERVICE}/PerformUnitFileAction"
echo

echo "== PerformUnitAction: STOP tailscaled.service =="
grpcurl -plaintext \
  -d '{
    "unitName": "tailscaled.service",
    "action": "UNIT_ACTION_STOP",
    "force": true
  }' \
  "${ADDR}" \
  "${SERVICE}/PerformUnitAction"
echo


echo "== PerformUnitFileAction: DISABLE tailscaled.service =="
grpcurl -plaintext \
  -d '{
    "unitName": "tailscaled.service",
    "action": "UNIT_FILE_ACTION_DISABLE",
    "runtime": true,
    "force": true
  }' \
  "${ADDR}" \
  "${SERVICE}/PerformUnitFileAction"
echo

echo "== All gRPC v2 tests passed =="
