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
    "unitName": "tailscaled.service"
  }' \
  "${ADDR}" \
  "${SERVICE}/GetUnitStatus"
echo

echo "== PerformAction: restart ssh.service =="
grpcurl -plaintext \
  -d '{
    "serviceName": "tailscaled.service",
    "unitAction": "UNIT_ACTION_STOP"
  }' \
  "${ADDR}" \
  "${SERVICE}/PerformAction"
echo

echo "== PerformAction: enable ssh.service (runtime=true, force=true) =="
grpcurl -plaintext \
  -d '{
    "serviceName": "tailscaled.service",
    "unitFileAction": "UNIT_FILE_ACTION_ENABLE",
    "runtime": true,
    "force": true
  }' \
  "${ADDR}" \
  "${SERVICE}/PerformAction"
echo

echo "== All gRPC tests passed =="
