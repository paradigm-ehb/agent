package test 

import (
	"context"
	pb "paradigm-ehb/agent/internal/connection/pb"
	"reflect"
	"testing"
	
)

func Test_server_SayHello(t *testing.T) {
	type fields struct {
		UnimplementedGreeterServer pb.UnimplementedGreeterServer
	}
	type args struct {
		in0 context.Context
		in  *pb.HelloRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pb.HelloReply
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server{
				UnimplementedGreeterServer: tt.fields.UnimplementedGreeterServer,
			}
			got, err := s.SayHello(tt.args.in0, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("server.SayHello() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.SayHello() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
