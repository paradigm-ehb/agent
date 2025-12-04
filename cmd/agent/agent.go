package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb_greeter "paradigm-ehb/agent/gen/greet"
	"paradigm-ehb/agent/pkg/service"

	manager "paradigm-ehb/agent/internal/svcmanager"
	tools "paradigm-ehb/agent/tools"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

func main() {
	fmt.Println(`
                               .___.__                                                       __   
___________ ____________     __| _/|__| ____   _____           _____     ____   ____   _____/  |_ 
\____ \__  \\_  __ \__  \   / __ | |  |/ ___\ /     \   ______ \__  \   / ___\_/ __ \ /    \   __\
|  |_> > __ \|  | \// __ \_/ /_/ | |  / /_/  >  Y Y  \ /_____/  / __ \_/ /_/  >  ___/|   |  \  |  
|   __(____  /__|  (____  /\____ | |__\___  /|__|_|  /         (____  /\___  / \___  >___|  /__|  
|__|       \/           \/      \/   /_____/       \/               \//_____/      \/     \/      


	`)

	err := tools.CheckOSUser()
	if err != nil {
		fmt.Println("Operating system is currently not supported. Come back in .... never! Imagine not using Linux. Not worthy.")
		os.Exit(4)
	}

	// TODO: replace with actual services
	fmt.Println(manager.ListDbusObject())
	fmt.Println(manager.GetDisplayManager())

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	greeter_server := &service.GreeterServer{}
	pb_greeter.RegisterGreeterServer(server, greeter_server)

	reflection.Register(server)

	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
