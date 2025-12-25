package main

import (
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pbgreeter "paradigm-ehb/agent/gen/greet"
	"paradigm-ehb/agent/gen/journal/v1"
	"paradigm-ehb/agent/gen/services/v1"
	"paradigm-ehb/agent/pkg/service"

	"paradigm-ehb/agent/tools"

	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	resource "paradigm-ehb/agent/internal/resmanager"
)

var (
	port = flag.Int("port", 5000, "The server port")
)

func main() {

	resource.CreateAgentCpu()
	resource.CreateAgentRam()
	resource.CreateAgentDevice()
	resource.CreateAgentDisk()

	fmt.Println("started")
	err := tools.CheckOSUser()
	if err != nil {
		fmt.Println("Operating system is currently not supported. Come back in .... never! Imagine not using Linux. Not worthy.")
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

	// Register Health Checking Service
	healthServer := health.NewServer()
	healthgrpc.RegisterHealthServer(server, healthServer)

	// Register Greeter Service
	greeterServer := &service.GreeterServer{}
	actionServer := &service.HandlerService{}
	journalServer := &service.JournalService{}

	pbgreeter.RegisterGreeterServer(server, greeterServer)
	services.RegisterHandlerServiceServer(server, actionServer)
	journal.RegisterJournalServiceServer(server, journalServer)

	reflection.Register(server)

	fmt.Printf("\nserver listening at %v\n", lis.Addr())
	if err := server.Serve(lis); err != nil {
		fmt.Println("failed to serve: ", err)
	}
}
