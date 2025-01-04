package main

import (
	"context"
	"log"
	"time"

	pb "github.com/VictorChepkasov/protos/messenger"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMessengerClient(conn)

	// Sending a message
	_, err = client.SendMessage(context.Background(), &pb.Message{Sender: "Alice", Content: "Hello, Bob!"})
	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}

	// Receiving messages
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := client.ReceiveMessages(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("could not receive messages: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			break
		}
		log.Printf("Received message from %s: %s", msg.Sender, msg.Content)
	}
}
