#!/bin/sh
# echo -------------------------------------------------------------------------
# echo Description: build go binary with flags 
# echo -------------------------------------------------------------------------
VERSION=1.0

prod() {
    go build \
        -ldflags="-s -w -X main.version=${VERSION}" \
        -o agent \
        ./cmd/agent/agent.go
}

debug() {
    go build \
        -race \
        -gcflags="all=-N -l" \
        -o debug \
        ./cmd/agent/agent.go
}

valgrind() {
    echo "Building with debug symbols for valgrind..."
    CGO_CFLAGS="-g -O0" \
    CGO_LDFLAGS="-g" \
    go build \
        -gcflags="all=-N -l" \
        -o debug_valgrind \
        ./cmd/agent/agent.go
    
    echo "Running valgrind..."
    valgrind \
        --leak-check=full \
        --show-leak-kinds=all \
        --track-origins=yes \
        --verbose \
        --log-file=valgrind-out.txt \
        ./debug_valgrind
    
    echo "Valgrind output saved to valgrind-out.txt"
}

gdb() {
    echo "Building with debug symbols for gdb..."
    CGO_CFLAGS="-g -O0" \
    CGO_LDFLAGS="-g" \
    go build \
        -gcflags="all=-N -l" \
        -o debug_gdb \
        ./cmd/agent/agent.go
    
    echo "Starting gdb..."
    gdb ./debug_gdb
}

core() {
    echo "Building with debug symbols..."
    CGO_CFLAGS="-g -O0" \
    CGO_LDFLAGS="-g" \
    go build \
        -gcflags="all=-N -l" \
        -o debug_core \
        ./cmd/agent/agent.go
    
    echo "Enabling core dumps..."
    ulimit -c unlimited
    
    echo "Running program (will generate core dump on crash)..."
    ./debug_core
    
    # Check if core dump was created
    if [ -f core ]; then
        echo "Core dump generated. Starting gdb..."
        gdb ./debug_core core
    elif [ -f core.* ]; then
        CORE_FILE=$(ls -t core.* | head -n1)
        echo "Core dump generated: $CORE_FILE. Starting gdb..."
        gdb ./debug_core "$CORE_FILE"
    else
        echo "No core dump found. Program may have exited normally."
        echo "Core dumps might be in: /var/lib/systemd/coredump/ or /var/crash/"
        echo "Check with: coredumpctl list"
    fi
}

sanitize() {
    echo "Building with address sanitizer..."
    CGO_CFLAGS="-fsanitize=address -g -O0" \
    CGO_LDFLAGS="-fsanitize=address" \
    go build \
        -gcflags="all=-N -l" \
        -o debug_asan \
        ./cmd/agent/agent.go
    
    echo "Running with address sanitizer..."
    ASAN_OPTIONS=detect_leaks=1:halt_on_error=0 ./debug_asan
}

case "$1" in
    prod)
        prod
        ;;
    debug)
        debug
        ;;
    valgrind)
        valgrind
        ;;
    gdb)
        gdb
        ;;
    core)
        core
        ;;
    sanitize|asan)
        sanitize
        ;;
    *)
        echo "Usage: $0 [prod|debug|valgrind|gdb|core|sanitize]"
        echo ""
        echo "  prod      - Production build (stripped, optimized)"
        echo "  debug     - Debug build with race detector"
        echo "  valgrind  - Build and run with valgrind memory checker"
        echo "  gdb       - Build and start gdb debugger"
        echo "  core      - Build, enable core dumps, run and analyze crash"
        echo "  sanitize  - Build and run with address sanitizer (ASAN)"
        exit 1
        ;;
esac
