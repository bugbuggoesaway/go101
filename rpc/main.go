package main

import (
	"go101/rpc/server"
	"go101/rpc/service"
	"log"
)

func main() {
	srv := &service.HelloWorldService{}
	server := server.NewServer(srv, ":8080", server.GRPC)
	if err := server.Start(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
