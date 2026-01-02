package server_test


import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)


const bufSize = 1024 * 1024

var lis *bufconn.Listener

func BufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func init() {
	lis = bufconn.Listen(bufSize)

	srv := grpc.NewServer()


	go func() {
		if err := srv.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

