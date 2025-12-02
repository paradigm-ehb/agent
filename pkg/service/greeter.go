package service

import (
	"context"
	"log"

	"paradigm-ehb/agent/gen/greet"
)

type GreeterServer struct {
	greet.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *GreeterServer) SayHello(_ context.Context, in *greet.HelloRequest) (*greet.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &greet.HelloReply{Message: "Hello " + in.GetName()}, nil
}
