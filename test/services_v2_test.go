package server_test

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	pb "paradigm-ehb/agent/gen/services/v2"
)

func TestServiceV2_PerformAction_Start(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resolver.SetDefaultScheme("passthrough")
	clientConn, err := grpc.NewClient(
		"bufnet",
		grpc.WithContextDialer(BufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer clientConn.Close()

	client := pb.NewHandlerServiceClient(clientConn)
	resp, err := client.PerformAction(ctx, &pb.ServiceActionRequest{
		ServiceName:    "tailscaled.service",
		UnitAction:     pb.ServiceActionRequest_UNIT_ACTION_START.Enum(),
		UnitFileAction: pb.ServiceActionRequest_UNIT_FILE_ACTION_UNSPECIFIED.Enum(),
	})

	if err != nil {
		t.Fatalf("rpc failed: %v", err)
	}

	if resp == nil {
		t.Fatalf("response is nil")
	}

	if len(resp.Status) == 0 {
		t.Fatalf("response status is empty")
	}

	statusStr := string(resp.Status)
	t.Logf("received status: %s", statusStr)
	t.Logf("success: %v", resp.Success)

	if resp.ErrorMessage != "" {
		t.Logf("error message: %s", resp.ErrorMessage)
	}
}

func TestServiceV2_PerformAction_Stop(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resolver.SetDefaultScheme("passthrough")
	clientConn, err := grpc.NewClient(
		"bufnet",
		grpc.WithContextDialer(BufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer clientConn.Close()

	client := pb.NewHandlerServiceClient(clientConn)
	resp, err := client.PerformAction(ctx, &pb.ServiceActionRequest{
		ServiceName:    "tailscaled.service",
		UnitAction:     pb.ServiceActionRequest_UNIT_ACTION_STOP.Enum(),
		UnitFileAction: pb.ServiceActionRequest_UNIT_FILE_ACTION_UNSPECIFIED.Enum(),
	})
	if err != nil {
		t.Fatalf("rpc failed: %v", err)
	}
	if resp == nil {
		t.Fatalf("response is nil")
	}

	statusStr := string(resp.Status)
	t.Logf("received status: %s", statusStr)
	t.Logf("success: %v", resp.Success)
}
