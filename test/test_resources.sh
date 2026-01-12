#!/usr/bin/env bash
set -euo pipefail

# ============================================================
# Configuration
# ============================================================

ADDR="127.0.0.1:5000"
SERVICE="resources.v3.ResourcesService"

# ============================================================
# Helpers
# ============================================================

section() {
  echo
  echo "============================================================"
  echo "$1"
  echo "============================================================"
}

# ============================================================
# Tests
# ============================================================

section "List all services"
grpcurl -plaintext "$ADDR" list

section "List ResourcesService methods"
grpcurl -plaintext "$ADDR" list "$SERVICE"

section "Get full system snapshot"
grpcurl -plaintext \
  "$ADDR" \
  "$SERVICE/GetSystemResources"

section "Get CPU"
grpcurl -plaintext \
  "$ADDR" \
  "$SERVICE/GetCpu"

section "Get Memory"
grpcurl -plaintext \
  "$ADDR" \
  "$SERVICE/GetMemory"

section "Get Device"
grpcurl -plaintext \
  "$ADDR" \
  "$SERVICE/GetDevice"

section "Get Disks"
grpcurl -plaintext \
  "$ADDR" \
  "$SERVICE/GetDisks"

section "Get Filesystems (all)"
grpcurl -plaintext \
  -d '{}' \
  "$ADDR" \
  "$SERVICE/GetFileSystems"

section "Get Filesystem for /"
grpcurl -plaintext \
  -d '{ "path": "/" }' \
  "$ADDR" \
  "$SERVICE/GetFileSystems"

section "Get Processes"
grpcurl -plaintext \
  "$ADDR" \
  "$SERVICE/GetProcesses"

# ⚠️ Dangerous — sends SIGTERM to PID 1234
# Uncomment only if you know what you're doing
#
# section "Send SIGTERM to process 1234"
# grpcurl -plaintext \
#   -d '{ "pid": 1234, "signal": 15 }' \
#   "$ADDR" \
#   "$SERVICE/ProcessAction"

echo
echo "All grpcurl tests completed successfully."
