package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb_greeter "paradigm-ehb/agent/gen/greet"
	"paradigm-ehb/agent/gen/services/v1"
	"paradigm-ehb/agent/pkg/service"

	tools "paradigm-ehb/agent/tools"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	fmt.Println("started")

	err := tools.CheckOSUser()
	if err != nil {
		fmt.Println("Operating system is currently not supported. Come back in .... never! Imagine not using Linux. Not worthy.")
		os.Exit(4)
	}

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		fmt.Println("failed to listen: ", err)
	}

	server := grpc.NewServer()

	greeter_server := &service.GreeterServer{}
	action_server := &service.HandlerService{}

	pb_greeter.RegisterGreeterServer(server, greeter_server)
	services.RegisterHandlerServiceServer(server, action_server)

	reflection.Register(server)

	fmt.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		fmt.Println("failed to serve: ", err)
	}

}
