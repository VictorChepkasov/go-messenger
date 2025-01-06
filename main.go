package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/VictorChepkasov/go-messenger/broadcast"
	proto "github.com/VictorChepkasov/go-messenger/gen"
)

func main() {
	// Create a new gRPC server
	grpcServer := grpc.NewServer()

	// Create a new connection pool
	var conn []*broadcast.Connection

	pool := &broadcast.Pool{
		Connection: conn,
	}

	// Register the pool with the gRPC server
	proto.RegisterBroadcastServer(grpcServer, pool)

	// Create a TCP listener at port 8080
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatalf("Error creating the server %v", err)
	}

	fmt.Println("Server started at port :8080")

	// Start serving requests at port 8080
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error creating the server %v", err)
	}
}
