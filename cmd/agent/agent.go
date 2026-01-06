// Package main boots the Paradigm Agent gRPC server.
//
// Responsibilities:
//   - Parse runtime flags (IP, port, diagnostics)
//   - Bind a TCP listener with automatic port fallback
//   - Initialize and register all gRPC services
//   - Expose health and reflection endpoints
//
// Design notes:
//   - The server must be able to start even if the preferred port is occupied.
//   - Reflection is enabled by default for debugging and introspection.
//   - Diagnostics are intentionally decoupled from server startup to avoid
//     blocking gRPC reflection and request handling.
package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"paradigm-ehb/agent/gen/greet"
	"paradigm-ehb/agent/gen/journal/v1"
	"paradigm-ehb/agent/gen/resources/v1"
	"paradigm-ehb/agent/gen/services/v1"
	luke "paradigm-ehb/agent/gen/actions/v1"
	// "paradigm-ehb/agent/internal/deviceactions"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	// "paradigm-ehb/agent/internal/platform"
	"paradigm-ehb/agent/pkg/grpc_handler"
	"syscall"
	// "time"
)

var (
	/**
		diagnostics enables periodic runtime diagnostics such as
		resource usage, process health, and connectivity checks.
		NOTE: When enabled, diagnostics must never block the gRPC server.
	 **/

	diagnostics = flag.Bool("diagnostics", false, "run runtime diagnostics")

	/**
	portFlag is the preferred TCP port to bind the gRPC server to.
	If unavailable, the server will increment the port until a free
	one is found.
	*/

	portFlag = flag.Int("port", 5000, "port to listen on")

	/**
	  ipFlag defines the IP address diagnostics may use when reporting
	  or exposing runtime information.
	*/
	ipFlag = flag.String("ip", "0.0.0.0", "ip addr")
)

func main() {
	/**
	Parse command-line flags before any runtime behavior.
	*/

	flag.Parse()

	/**
	Attempt to bind a TCP listener.

	Strategy:
	  - Start with the requested port
	  - If EADDRINUSE is encountered, increment the port and retry
	  - Fail hard on any other error

	This guarantees the agent can always start, even in constrained
	or multi-agent environments.

	*/
	var lis net.Listener
	var err error
	p := *portFlag

	for {
		addr := fmt.Sprintf("0.0.0.0:%d", p)
		lis, err = net.Listen("tcp4", addr)
		if err != nil {
			var opErr *net.OpError
			if errors.As(err, &opErr) {
				if sysErr, ok := opErr.Err.(*os.SyscallError); ok {
					if sysErr.Err == syscall.EADDRINUSE {
						p++
						continue
					}
				}
			}

			/**
			Any error other than "address already in use" is fatal.
			*/
			fmt.Println("failed to listen:", err)
			return
		}
		break
	}

	/**
	Create the gRPC server instance.
	*/
	server := grpc.NewServer()

	/**
	Health server is used by orchestration systems (systemd, Kubernetes,
	external monitors) to determine liveness and readiness.
	*/
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	/**

	  Register all application services.

	  Each service implements a distinct responsibility:
	    - Greeter: connectivity / handshake testing
	    - HandlerService: service lifecycle and orchestration
	    - JournalService: event and state journaling
	    - ResourcesService: system resource inspection and reporting

	*/
	greet.RegisterGreeterServer(server, &grpc_handler.GreeterServer{})
	services.RegisterHandlerServiceServer(server, &grpc_handler.HandlerService{})
	journal.RegisterJournalServiceServer(server, &grpc_handler.JournalService{})
	resourcespb.RegisterResourcesServiceServer(server, &grpc_handler.ResourcesService{})
	luke.RegisterActionServiceServer(server, &grpc_handler.DeviceActionsService{})
	

	/*
		Diagnostics mode (disabled for now)

		Problem:
		  - Running diagnostics synchronously introduces an infinite loop
		    that blocks gRPC reflection and request handling.

		Important invariant:
		  - Diagnostics must run asynchronously and must never interfere
		    with server startup, reflection, or request processing.

		TODO(nasr):
		  - Move diagnostics into a separate goroutine
		  - Introduce a proper shutdown context
		  - Ensure diagnostics respect server lifecycle
	*/

	/**
	Enable gRPC reflection unconditionally.

	This allows tools such as grpcurl and gRPC UI to inspect services
	and message schemas at runtime.
	*/
	reflection.Register(server)

	fmt.Printf("\nserver listening at %v\n", lis.Addr())

	/**
	Start serving requests.
	This call blocks until the server is stopped or encounters a fatal error.
	*/
	if err := server.Serve(lis); err != nil {
		fmt.Println("failed to serve:", err)
	}
}
