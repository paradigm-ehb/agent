/*
Package main boots the Paradigm Agent gRPC server.

Responsibilities:
- Parse runtime flags (IP, port, diagnostics)
- Bind a TCP listener with automatic port fallback
- Initialize and register all gRPC services
- Expose health and reflection endpoints

Design notes:
  - The server must be able to start even if the preferred port is occupied.
  - Reflection is enabled by default for debugging and introspection.
  - Diagnostics are intentionally decoupled from server startup to avoid
    blocking gRPC reflection and request handling.
*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	devacpb "paradigm-ehb/agent/gen/actions/v1"
	"paradigm-ehb/agent/gen/greet"
	"paradigm-ehb/agent/gen/journal/v1"

	"log"

	servicesV1 "paradigm-ehb/agent/gen/services/v1"
	servicesV2 "paradigm-ehb/agent/gen/services/v2"
	servicesV3 "paradigm-ehb/agent/gen/services/v3"

	resourcesv1 "paradigm-ehb/agent/gen/resources/v1"
	resourcesv2 "paradigm-ehb/agent/gen/resources/v2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"paradigm-ehb/agent/internal/platform"

	"paradigm-ehb/agent/pkg/grpchandler"

	resourcesHandlerV1 "paradigm-ehb/agent/pkg/grpchandler/resources/v1"
	resourcesHandlerV2 "paradigm-ehb/agent/pkg/grpchandler/resources/v2"

	servicesHandlerV1 "paradigm-ehb/agent/pkg/grpchandler/services/v1"
	servicesHandlerV2 "paradigm-ehb/agent/pkg/grpchandler/services/v2"
	servicesHandlerV3 "paradigm-ehb/agent/pkg/grpchandler/services/v3"

	"syscall"
	"time"
)

var (
	/*
		diagnostics enables periodic runtime diagnostics such as
		resource usage, process health, and connectivity checks.

		TODO:
		- Add structured configuration for diagnostics intervals.
		- Allow diagnostics to be toggled or reconfigured at runtime.
		- Ensure diagnostics respect context cancellation on shutdown.

		NOTE:
		When enabled, diagnostics must never block the gRPC server.
	*/
	diagnostics = flag.Bool("diagnostics", true, "run runtime diagnostics")

	/*
		portFlag is the preferred TCP port to bind the gRPC server to.

		If unavailable:
		- The server increments the port until a free one is found.

		TODO:
		- Log the final selected port explicitly.
		- Optionally expose the selected port via diagnostics or metadata.
	*/
	portFlag = flag.Int("port", 5000, "port to listen on")

	/*
		ipFlag defines the IP address diagnostics may use when reporting
		or exposing runtime information.

		TODO:
		- Validate IP format early.
		- Clarify distinction between bind address vs diagnostics address.
	*/
	ipFlag = flag.String("ip", "0.0.0.0", "ip addr")

	/**
	- implement mTLS
	- user generates own certificates
	*/

	mTlsFlag = flag.Bool("mtls", false, "mTLS authentication")

	/**
	  - generate token token token
	  - generates a short lived token that is used to generate the ssl certificates
	*/

	generate = flag.Int64("generate", 10, "generate a short lived token that will live for x amount of time")
)

func startGrpc() {

	var lis net.Listener
	/*
		Create the gRPC server instance.

		TODO:
		- Configure server options (timeouts, interceptors, limits).
		- Add graceful shutdown handling.
	*/
	server := grpc.NewServer()

	/*
		Health server is used by orchestration systems (systemd, Kubernetes,
		external monitors) to determine liveness and readiness.

		TODO:
		- Set explicit serving statuses per service.
	*/
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	/*
		Register all application services.

		Each service implements a distinct responsibility:
		- Greeter: connectivity / handshake testing
		- HandlerService: service lifecycle and orchestration
		- JournalService: event and state journaling
		- ResourcesService: system resource inspection and reporting

		TODO:
		- Centralize service registration.
		- Version-gate deprecated service versions.
	*/

	greet.RegisterGreeterServer(
		server,
		&grpc_handler.GreeterServer{},
	)

	servicesV1.RegisterHandlerServiceServer(
		server,
		&servicesHandlerV1.HandlerService{},
	)

	servicesV2.RegisterHandlerServiceServer(
		server,
		&servicesHandlerV2.HandlerServiceV2{},
	)

	servicesV3.RegisterHandlerServiceServer(
		server,
		&servicesHandlerV3.HandlerServicev3{},
	)

	journal.RegisterJournalServiceServer(
		server,
		&grpc_handler.JournalService{},
	)

	resourcesv1.RegisterResourcesServiceServer(
		server,
		&resourcesHandlerV1.ResourcesService{},
	)

	resourcesv2.RegisterResourcesServiceServer(
		server,
		&resourcesHandlerV2.ResourcesServiceV2{},
	)

	devacpb.RegisterActionServiceServer(
		server,
		&grpc_handler.DeviceActionsService{},
	)

	/*
		Enable gRPC reflection unconditionally.

		This allows tools such as grpcurl and gRPC UI to inspect services
		and message schemas at runtime.

		TODO:
		- Make reflection configurable for production environments.
	*/
	reflection.Register(server)

	log.Printf("\nport: %v\n", lis.Addr())

	if *diagnostics {

		/*
			TODO:
			- Tie diagnostics lifecycle to server context.
			- Ensure diagnostics terminate on server shutdown.
		*/
		go platform.RunRuntimeDiagnostics(time.Second*2, *ipFlag, *portFlag)
	}

	/*
		Start serving requests.

		This call blocks until the server is stopped or encounters a fatal error.

		TODO:
		- Implement graceful shutdown (signals, context).
		- Flush diagnostics and logs before exit.
	*/
	if err := server.Serve(lis); err != nil {
		log.Printf("failed to serve: %v", err)
	}
}

func openPort() {
	/*
		Attempt to bind a TCP listener.

		Strategy:
		- Start with the requested port
		- If EADDRINUSE is encountered, increment the port and retry
		- Fail hard on any other error

		This guarantees the agent can always start, even in constrained
		or multi-agent environments.

		TODO:
		- Add an upper bound to port scanning.
		- Support IPv6 or configurable network protocols.
	*/
	p := *portFlag

	for {

		addr := fmt.Sprintf("0.0.0.0:%d", p)
		_, err := net.Listen("tcp4", addr)
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

			/*
				Any error other than "address already in use" is fatal.

				TODO:
				- Emit structured logs.
				- Exit with non-zero status code.
			*/
			log.Printf("failed to listen: %v", err)
			return
		}
		break
	}

}

func main() {

	/*
		Parse command-line flags before any runtime behavior.

		TODO:
		- Add validation for flag combinations.
		- Print effective configuration at startup.
	*/
	flag.Parse()

	startGrpc()

}
