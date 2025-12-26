package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pbgreeter "paradigm-ehb/agent/gen/greet"
	"paradigm-ehb/agent/gen/journal/v1"
	resourcespb "paradigm-ehb/agent/gen/resources/v1"
	"paradigm-ehb/agent/gen/services/v1"
	"paradigm-ehb/agent/pkg/service"

	resources "paradigm-ehb/agent/internal/resmanager"
	"paradigm-ehb/agent/tools"

	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	port = flag.Int("port", 5000, "The server port")
)

func main() {

	resources.Make()

	fmt.Println("started")

	if err := tools.CheckOSUser(); err != nil {
		fmt.Println("Operating system is currently not supported.")
		return
	}

	flag.Parse()

	var lis net.Listener
	var err error
	p := *port

	for {
		addr := fmt.Sprintf("0.0.0.0:%d", p)
		lis, err = net.Listen("tcp4", addr)
		if err != nil {
			// Retry only if port is already in use
			var opErr *net.OpError
			if errors.As(err, &opErr) {
				if sysErr, ok := opErr.Err.(*os.SyscallError); ok {
					if sysErr.Err == syscall.EADDRINUSE {
						p++
						continue
					}
				}
			}

			fmt.Println("failed to listen:", err)
			return
		}

		break
	}

	server := grpc.NewServer()

	// Health service
	healthServer := health.NewServer()
	healthgrpc.RegisterHealthServer(server, healthServer)

	// Application services
	pbgreeter.RegisterGreeterServer(server, &service.GreeterServer{})
	services.RegisterHandlerServiceServer(server, &service.HandlerService{})
	journal.RegisterJournalServiceServer(server, &service.JournalService{})
	resourcespb.RegisterResourcesServiceServer(server, &service.ResourcesService{})

	reflection.Register(server)

	fmt.Printf("\nserver listening at %v\n", lis.Addr())

	if err := server.Serve(lis); err != nil {
		fmt.Println("failed to serve:", err)
	}
}
