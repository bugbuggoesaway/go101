package service

import (
	"context"
	"fmt"
	helloworld "go101/rpc/grpc"
)

type HelloWorldService struct {
	helloworld.UnimplementedGreeterServer
}

func (s *HelloWorldService) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	fmt.Printf("Receive. name=[%s]\n", in.GetName())
	return &helloworld.HelloReply{Message: "Hello " + in.GetName()}, nil
}
