#!/bin/bash

# gRPCurl test script for ResourcesServiceV2
# Make sure your gRPC server is running before executing these commands

# Color codes for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default server address
SERVER="localhost:5000"

# Check if custom server address provided
if [ ! -z "$1" ]; then
    SERVER="$1"
fi

echo -e "${BLUE}Testing ResourcesServiceV2 on ${SERVER}${NC}\n"

# List available services
echo -e "${GREEN}=== Listing available services ===${NC}"
grpcurl -plaintext ${SERVER} list
echo ""

# List methods for ResourcesServiceV2
echo -e "${GREEN}=== Listing methods for resources.v2.ResourcesService ===${NC}"
grpcurl -plaintext ${SERVER} list resources.v2.ResourcesService
echo ""

# Describe the service
echo -e "${GREEN}=== Describing ResourcesService ===${NC}"
grpcurl -plaintext ${SERVER} describe resources.v2.ResourcesService
echo ""

# Test GetSystemResources
echo -e "${GREEN}=== Getting System Resources ===${NC}"
grpcurl -plaintext \
    -d '{}' \
    ${SERVER} \
    resources.v2.ResourcesService/GetSystemResources
echo ""

# Test KillProcess (example with PID 12345)
# WARNING: Uncomment and modify PID carefully - this will kill a process!
# echo -e "${GREEN}=== Killing Process (PID: 12345) ===${NC}"
# grpcurl -plaintext \
#     -d '{"pid": 12345}' \
#     ${SERVER} \
#     resources.v2.ResourcesService/KillProcess
# echo ""

# Get system resources with formatted output (using jq if available)
echo -e "${GREEN}=== Getting System Resources (formatted) ===${NC}"
if command -v jq &> /dev/null; then
    grpcurl -plaintext \
        -d '{}' \
        ${SERVER} \
        resources.v2.ResourcesService/GetSystemResources | jq '.'
else
    echo "jq not installed - skipping formatted output"
    echo "Install jq for pretty JSON: sudo apt-get install jq"
fi
echo ""

# Show CPU information only
echo -e "${GREEN}=== CPU Information ===${NC}"
if command -v jq &> /dev/null; then
    grpcurl -plaintext \
        -d '{}' \
        ${SERVER} \
        resources.v2.ResourcesService/GetSystemResources | jq '.resources.cpu'
else
    echo "jq not installed - install it to filter output"
fi
echo ""

# Show Memory information only
echo -e "${GREEN}=== Memory Information ===${NC}"
if command -v jq &> /dev/null; then
    grpcurl -plaintext \
        -d '{}' \
        ${SERVER} \
        resources.v2.ResourcesService/GetSystemResources | jq '.resources.memory'
else
    echo "jq not installed - install it to filter output"
fi
echo ""

# Show first 5 processes
echo -e "${GREEN}=== First 5 Processes ===${NC}"
if command -v jq &> /dev/null; then
    grpcurl -plaintext \
        -d '{}' \
        ${SERVER} \
        resources.v2.ResourcesService/GetSystemResources | jq '.resources.processes[:5]'
else
    echo "jq not installed - install it to filter output"
fi
echo ""

# Count total processes
echo -e "${GREEN}=== Total Process Count ===${NC}"
if command -v jq &> /dev/null; then
    grpcurl -plaintext \
        -d '{}' \
        ${SERVER} \
        resources.v2.ResourcesService/GetSystemResources | jq '.resources.processes | length'
else
    echo "jq not installed - install it to filter output"
fi
echo ""

echo -e "${BLUE}Tests completed!${NC}"
