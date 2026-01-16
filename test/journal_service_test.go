package server_test
//
// import (
// 	"context"
// 	"testing"
// 	"time"
//
// 	pb "paradigm-ehb/agent/gen/journal/v1"
//
// 	"google.golang.org/grpc/resolver"
//
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )
//
//
// func TestJournal_all(t *testing.T) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()
//
// 	resolver.SetDefaultScheme("passthrough")
//
// 	clientConn, err := grpc.NewClient(
// 		"bufnet",
// 		grpc.WithContextDialer(BufDialer),
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	)
//
// 	if err != nil {
// 		t.Fatalf("failed to create client: %v", err)
// 	}
// 	defer clientConn.Close()
//
// 	client := pb.NewJournalServiceClient(clientConn)
//
// 	resp, err := client.Action(ctx, &pb.JournalRequest{})
// 	if err != nil {
// 		t.Fatalf("rpc failed: %v", err)
// 	}
//
// 	if resp == nil {
// 		t.Fatalf("unexpected response: %q", resp)
// 	}
// }
