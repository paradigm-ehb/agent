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

	"paradigm-ehb/agent/tools"

	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

var (
	port = flag.Int("port", 5000, "The server port")
)

func main() {
	flag.Parse()

	var lis net.Listener
	var err error
	p := *port

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

			fmt.Println("Failed to find open port", err)
			return
		}
		break
	}

	server := grpc.NewServer()

	healthServer := health.NewServer()
	healthgrpc.RegisterHealthServer(server, healthServer)

	pbgreeter.RegisterGreeterServer(server, &service.GreeterServer{})
	services.RegisterHandlerServiceServer(server, &service.HandlerService{})
	journal.RegisterJournalServiceServer(server, &service.JournalService{})
	resourcespb.RegisterResourcesServiceServer(server, &service.ResourcesService{})

	reflection.Register(server)

	fmt.Printf("Agent initialized\nListening:\nport:%v\n", p)

	if err := tools.AssertLinux(); err != nil {
		fmt.Println(err)
		return
	}

	go tools.RunRuntimeDiagnostics(3 * time.Second)

	if err := server.Serve(lis); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
