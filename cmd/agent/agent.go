package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pbgreeter "paradigm-ehb/agent/gen/greet"
	"paradigm-ehb/agent/gen/journal/v1"
	"paradigm-ehb/agent/gen/services/v1"
	"paradigm-ehb/agent/pkg/service"

	"paradigm-ehb/agent/tools"

	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	port  = flag.Int("port", 50051, "The server port")
	sleep = flag.Duration("sleep", time.Second*5, "duration between changes in health")

	system = "Greeter" // empty string represents the health of the system
)

func main() {
	fmt.Println("started")

	err := tools.CheckOSUser()
	if err != nil {
		fmt.Println("Operating system is currently not supported. Come back in .... never! Imagine not using Linux. Not worthy.")
	}

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		fmt.Println("failed to listen: ", err)
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
