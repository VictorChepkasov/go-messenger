package main

import (
	"fmt"
	"log"
	"net"

	proto "messenger/internal/gen"

	broadcast "maessenger/internal/services/broadcast"

	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()

	var conn []*broadcast.Connection
	var port string = ":8080"

	pool := &broadcast.Pool{
		COnnection: conn,
	}

	// Register gRPC server
	proto.RegisterBroadcastServer(grpcServer, pool)

	// Create a TCP listener
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("error creating as port %v", port)
	}

	fmt.Printf("Server started as port %v", port)

	// Start serving requests
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Error creating the server %v", err)
	}
}
