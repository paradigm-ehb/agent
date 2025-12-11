package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb_greeter "paradigm-ehb/agent/gen/greet"
	manager "paradigm-ehb/agent/internal/svcmanager"
	"paradigm-ehb/agent/internal/svcmanager/servicecontrol"
	"paradigm-ehb/agent/pkg/service"

	tools "paradigm-ehb/agent/tools"

	dh "paradigm-ehb/agent/internal/svcmanager/dbushandler"
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

	conn, err := dh.CreateSystemBus()
	if err != nil {
		fmt.Println("failed to create shared systembus")
	}

	defer conn.Close()

	err = manager.RunAction(conn, servicecontrol.Stop, "mariadb.service")
	if err != nil {
		fmt.Println("failed to perform action on mariadb")
	}

	err = manager.RunSymlinkAction(conn, servicecontrol.Disable, true, true, []string{"mariadb.service"})
	if err != nil {

		fmt.Println("failed to perform symlink action on mariadb")
	}

	err = manager.RunRetrieval(conn, true)
	if err != nil {
		fmt.Println("failed to do this")
	}

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		fmt.Println("failed to listen: ", err)
	}

	server := grpc.NewServer()

	greeter_server := &service.GreeterServer{}
	pb_greeter.RegisterGreeterServer(server, greeter_server)

	reflection.Register(server)

	fmt.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		fmt.Println("failed to serve: ", err)
	}

}
