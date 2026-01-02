package server_test

import (
	"context"
	"net"
	"testing"
	"time"

	pb "paradigm-ehb/agent/gen/resources/v1"
	"paradigm-ehb/agent/pkg/grpc_handler"

	"google.golang.org/grpc/resolver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)

	srv := grpc.NewServer()


	pb.RegisterResourcesServiceServer(srv, &grpc_handler.ResourcesService{})

	go func() {
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestResources_All(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()


	resolver.SetDefaultScheme("passthrough")

	clientConn, err := grpc.NewClient(
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer clientConn.Close()

	client := pb.NewResourcesServiceClient(clientConn)

	resp, err := client.GetSystemResources(ctx, &pb.GetSystemResourcesRequest{})
	if err != nil {
		t.Fatalf("rpc failed: %v", err)
	}

	if resp != nil {
		t.Fatalf("unexpected response: %q", resp)
	}
}
