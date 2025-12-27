package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"paradigm-ehb/agent/internal/platform"
	"strconv"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"paradigm-ehb/agent/gen/greet"
	"paradigm-ehb/agent/gen/journal/v1"

	"paradigm-ehb/agent/gen/resources/v1"

	"paradigm-ehb/agent/gen/services/v1"
	"paradigm-ehb/agent/pkg/service"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var diagnostics = flag.Bool("diagnostics", false, "run runtime diagnostics")
var portFlag = flag.Int("port", 5000, "port to listen on")
var ipFlag = flag.String("ip", "0.0.0.0", "ip addr")

func main() {
	if err := platform.AssertLinux(); err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	addr := net.JoinHostPort(*ipFlag, strconv.Itoa(*portFlag))

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()

	server := grpc.NewServer()

	healthServer := health.NewServer()

	grpc_health_v1.RegisterHealthServer(server, healthServer)
	greet.RegisterGreeterServer(server, &service.GreeterServer{})
	services.RegisterHandlerServiceServer(server, &service.HandlerService{})
	journal.RegisterJournalServiceServer(server, &service.JournalService{})
	resourcespb.RegisterResourcesServiceServer(server, &service.ResourcesService{})

	reflection.Register(server)

	if *diagnostics == true {

		platform.RunRuntimeDiagnostics(time.Second*3, *ipFlag, *portFlag)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sig
		server.GracefulStop()
	}()

	log.Printf("listening on %s\n", addr)

	if err := server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
