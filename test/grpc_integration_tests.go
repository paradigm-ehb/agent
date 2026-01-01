package test

import (
	"context"
	// "net"
	// "os"
	"testing"
	"time"

	journalpb "paradigm-ehb/agent/gen/journal/v1"
	resourcespb "paradigm-ehb/agent/gen/resources/v1"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func dialGRPC(t *testing.T) *grpc.ClientConn {
	t.Helper()

	/**
	host :=  "localhost"
	port :=  "5000"
	addr := net.JoinHostPort(host, port)
	*/

	/**
	TODO(nasr): use grpc.NewClient
	*/

	return nil
}

func TestGRPCIntegration(t *testing.T) {
	conn := dialGRPC(t)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	/**
	Health check (== grpcurl list)
	*/
	t.Run("health", func(t *testing.T) {
		health := healthpb.NewHealthClient(conn)

		resp, err := health.Check(ctx, &healthpb.HealthCheckRequest{})
		if err != nil {
			t.Fatalf("health check failed: %v", err)
		}
		if resp.Status != healthpb.HealthCheckResponse_SERVING {
			t.Fatalf("server not serving: %v", resp.Status)
		}
	})

	/**
	JournalService.Action (GID=1000)
	*/
	t.Run("journal_action_gid", func(t *testing.T) {
		journal := journalpb.NewJournalServiceClient(conn)

		_, err := journal.Action(ctx, &journalpb.JournalRequest{
			NumFromTail: 5,
			Field:       journalpb.JournalRequest_GID,
			Value:       "1000",
		})
		if err != nil {
			t.Fatalf("JournalService.Action (GID) failed: %v", err)
		}
	})

	/**
	JournalService.Action (systemd unit)
	*/
	t.Run("journal_action_systemd", func(t *testing.T) {
		journal := journalpb.NewJournalServiceClient(conn)

		_, err := journal.Action(ctx, &journalpb.JournalRequest{
			NumFromTail: 5,
			Field:       journalpb.JournalRequest_Systemd,
			Value:       "systemd-journald.service",
		})
		if err != nil {
			t.Fatalf("JournalService.Action (systemd) failed: %v", err)
		}
	})

	/**
	  ResourcesService.GetSystemResources
	*/
	t.Run("resources", func(t *testing.T) {
		res := resourcespb.NewResourcesServiceClient(conn)

		resp, err := res.GetSystemResources(ctx, &resourcespb.GetSystemResourcesRequest{})
		if err != nil {
			t.Fatalf("GetSystemResources failed: %v", err)
		}

		if resp.Resources == nil {
			t.Fatal("GetSystemResources returned nil Resources")
		}
	})
}
