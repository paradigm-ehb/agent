package server_test


import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"paradigm-ehb/agent/pkg/grpc_handler"
	respb "paradigm-ehb/agent/gen/resources/v1"
	serpb_v1 "paradigm-ehb/agent/gen/services/v1"
	serpb_v2 "paradigm-ehb/agent/gen/services/v2"

	greetpb "paradigm-ehb/agent/gen/greet"
	journalpb "paradigm-ehb/agent/gen/journal/v1"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"


)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func BufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func init() {
	lis = bufconn.Listen(bufSize)

	server := grpc.NewServer()
	healthServer := health.NewServer()

	grpc_health_v1.RegisterHealthServer(server, healthServer)

	respb.RegisterResourcesServiceServer(server, &grpc_handler.ResourcesService{})
	serpb_v1.RegisterHandlerServiceServer(server, &grpc_handler.HandlerService{})
	serpb_v2.RegisterHandlerServiceServer(server, &grpc_handler.HandlerServiceV2{})
	greetpb.RegisterGreeterServer(server, &grpc_handler.GreeterServer{})
	journalpb.RegisterJournalServiceServer(server, &grpc_handler.JournalService{})

	go func() {
		if err := server.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

