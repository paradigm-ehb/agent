package main

import (
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pbgreeter "paradigm-ehb/agent/gen/greet"
	"paradigm-ehb/agent/gen/services/v1"
	"paradigm-ehb/agent/pkg/service"

	"paradigm-ehb/agent/tools"
)

var (
	port = flag.Int("port", 50051, "The server port")
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

	greeterServer := &service.GreeterServer{}
	actionServer := &service.HandlerService{}

	pbgreeter.RegisterGreeterServer(server, greeterServer)
	services.RegisterHandlerServiceServer(server, actionServer)

	reflection.Register(server)

	fmt.Printf("\nserver listening at %v\n", lis.Addr())
	if err := server.Serve(lis); err != nil {
		fmt.Println("failed to serve: ", err)
	}
}
