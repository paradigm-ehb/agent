#!/bin/bash
# gRPCurl test commands for all services
# Make sure your gRPC server is running on localhost:5000 (adjust port as needed)
# Install grpcurl: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

SERVER="localhost:5000"

echo "======================================"
echo "ResourcesServiceV2 Tests"
echo "======================================"

# Get System Resources
echo -e "\n[1] GetSystemResources"
grpcurl -plaintext \
  -d '{}' \
  $SERVER \
  resources.v2.ResourcesService/GetSystemResources

# Kill Process (replace with actual PID)
echo -e "\n[2] KillProcess (SIGTERM - signal 15)"
grpcurl -plaintext \
  -d '{
    "pid": 12345,
    "signal": 15
  }' \
  $SERVER \
  resources.v2.ResourcesService/KillProcess

echo -e "\n[3] KillProcess (SIGKILL - signal 9)"
grpcurl -plaintext \
  -d '{
    "pid": 12345,
    "signal": 9
  }' \
  $SERVER \
  resources.v2.ResourcesService/KillProcess

echo "======================================"
echo "HandlerServicev3 Tests (Systemd)"
echo "======================================"

# Get All Units
echo -e "\n[4] GetAllUnits"
grpcurl -plaintext \
  -d '{}' \
  $SERVER \
  services.v3.HandlerService/GetAllUnits

# Get Loaded Units
echo -e "\n[5] GetLoadedUnits"
grpcurl -plaintext \
  -d '{}' \
  $SERVER \
  services.v3.HandlerService/GetLoadedUnits

# Get Filtered Units
echo -e "\n[6] GetFilteredUnits - LOADED"
grpcurl -plaintext \
  -d '{
    "filters": ["LOADED"]
  }' \
  $SERVER \
  services.v3.HandlerService/GetFilteredUnits

echo -e "\n[7] GetFilteredUnits - ERROR and MASKED"
grpcurl -plaintext \
  -d '{
    "filters": ["ERROR", "MASKED"]
  }' \
  $SERVER \
  services.v3.HandlerService/GetFilteredUnits

# Get Unit Status
echo -e "\n[8] GetUnitStatus - nginx.service"
grpcurl -plaintext \
  -d '{
    "unit_name": "nginx.service"
  }' \
  $SERVER \
  services.v3.HandlerService/GetUnitStatus

echo -e "\n[9] GetUnitStatus - sshd.service"
grpcurl -plaintext \
  -d '{
    "unit_name": "sshd.service"
  }' \
  $SERVER \
  services.v3.HandlerService/GetUnitStatus

# Perform Unit Action - Start
echo -e "\n[10] PerformUnitAction - START nginx.service"
grpcurl -plaintext \
  -d '{
    "unit_name": "nginx.service",
    "action": "UNIT_ACTION_START"
  }' \
  $SERVER \
  services.v3.HandlerService/PerformUnitAction

# Perform Unit Action - Stop
echo -e "\n[11] PerformUnitAction - STOP nginx.service"
grpcurl -plaintext \
  -d '{
    "unit_name": "nginx.service",
    "action": "UNIT_ACTION_STOP"
  }' \
  $SERVER \
  services.v3.HandlerService/PerformUnitAction

# Perform Unit Action - Restart
echo -e "\n[12] PerformUnitAction - RESTART nginx.service"
grpcurl -plaintext \
  -d '{
    "unit_name": "nginx.service",
    "action": "UNIT_ACTION_RESTART"
  }' \
  $SERVER \
  services.v3.HandlerService/PerformUnitAction

# Perform Unit File Action - Enable
echo -e "\n[13] PerformUnitFileAction - ENABLE nginx.service"
grpcurl -plaintext \
  -d '{
    "unit_name": "nginx.service",
    "action": "UNIT_FILE_ACTION_ENABLE",
    "runtime": false,
    "force": false
  }' \
  $SERVER \
  services.v3.HandlerService/PerformUnitFileAction

# Perform Unit File Action - Disable
echo -e "\n[14] PerformUnitFileAction - DISABLE nginx.service"
grpcurl -plaintext \
  -d '{
    "unit_name": "nginx.service",
    "action": "UNIT_FILE_ACTION_DISABLE",
    "runtime": false,
    "force": true
  }' \
  $SERVER \
  services.v3.HandlerService/PerformUnitFileAction

echo "======================================"
echo "DeviceActionsService Tests"
echo "======================================"

# Shutdown (WARNING: Will shutdown the system!)
echo -e "\n[15] Action - SHUTDOWN (commented out for safety)"
# grpcurl -plaintext \
#   -d '{
#     "device_action": "DEVICE_ACTION_SHUTDOWN"
#   }' \
#   $SERVER \
#   actions.v1.ActionService/Action

# Reboot (WARNING: Will reboot the system!)
echo -e "\n[16] Action - REBOOT (commented out for safety)"
# grpcurl -plaintext \
#   -d '{
#     "device_action": "DEVICE_ACTION_REBOOT"
#   }' \
#   $SERVER \
#   actions.v1.ActionService/Action

# Suspend
echo -e "\n[17] Action - SUSPEND"
grpcurl -plaintext \
  -d '{
    "device_action": "DEVICE_ACTION_SUSPEND"
  }' \
  $SERVER \
  actions.v1.ActionService/Action

# Hibernate
echo -e "\n[18] Action - HIBERNATE"
grpcurl -plaintext \
  -d '{
    "device_action": "DEVICE_ACTION_HIBERNATE"
  }' \
  $SERVER \
  actions.v1.ActionService/Action

echo "======================================"
echo "JournalService Tests (Server Streaming)"
echo "======================================"

# Journal by Systemd Unit
echo -e "\n[19] Action - Get journal for nginx.service"
grpcurl -plaintext \
  -d '{
    "field": 0,
    "value": "nginx.service",
    "num_from_tail": 100,
    "cursor": "",
    "path": ""
  }' \
  $SERVER \
  journal.v1.JournalService/Action

# Journal by PID
echo -e "\n[20] Action - Get journal for PID 1"
grpcurl -plaintext \
  -d '{
    "field": 1,
    "value": "1",
    "num_from_tail": 50,
    "cursor": "",
    "path": ""
  }' \
  $SERVER \
  journal.v1.JournalService/Action

# Journal by UID
echo -e "\n[21] Action - Get journal for UID 0 (root)"
grpcurl -plaintext \
  -d '{
    "field": 2,
    "value": "0",
    "num_from_tail": 100,
    "cursor": "",
    "path": ""
  }' \
  $SERVER \
  journal.v1.JournalService/Action

# Journal by GID
echo -e "\n[22] Action - Get journal for GID 0 (root)"
grpcurl -plaintext \
  -d '{
    "field": 3,
    "value": "0",
    "num_from_tail": 100,
    "cursor": "",
    "path": ""
  }' \
  $SERVER \
  journal.v1.JournalService/Action

echo -e "\n======================================"
echo "List Available Services"
echo "======================================"

grpcurl -plaintext $SERVER list

echo -e "\n======================================"
echo "Describe a Service"
echo "======================================"

grpcurl -plaintext $SERVER describe resources.v2.ResourcesService

echo -e "\n======================================"
echo "Tests Complete!"
echo "======================================"
