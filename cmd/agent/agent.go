package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb_greeter "paradigm-ehb/agent/gen/greet"
	"paradigm-ehb/agent/gen/services/v1"
	"paradigm-ehb/agent/pkg/service"

	tools "paradigm-ehb/agent/tools"

	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
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
		//os.Exit(4)
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
	greeter_server := &service.GreeterServer{}
	action_server := &service.HandlerService{}

	pb_greeter.RegisterGreeterServer(server, greeter_server)
	services.RegisterHandlerServiceServer(server, action_server)

	reflection.Register(server)

	go func(service string) {
		// asynchronously inspect dependencies and toggle serving status as needed
		next := healthpb.HealthCheckResponse_SERVING

		for {
			healthServer.SetServingStatus(service, next)

			if next == healthpb.HealthCheckResponse_SERVING {
				next = healthpb.HealthCheckResponse_NOT_SERVING
			} else {
				next = healthpb.HealthCheckResponse_SERVING
			}

			time.Sleep(*sleep)
		}
	}("Greeter")

	go func(service string) {
		// asynchronously inspect dependencies and toggle serving status as needed
		next := healthpb.HealthCheckResponse_SERVING

		for {
			healthServer.SetServingStatus(service, next)

			if next == healthpb.HealthCheckResponse_SERVING {
				next = healthpb.HealthCheckResponse_NOT_SERVING
			} else {
				next = healthpb.HealthCheckResponse_SERVING
			}

			time.Sleep(*sleep)
		}
	}("Services")

	healthServer.SetServingStatus("Greeter", healthpb.HealthCheckResponse_SERVING)

	fmt.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		fmt.Println("failed to serve: ", err)
	}
}
