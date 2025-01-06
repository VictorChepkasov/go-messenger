package main

import (
	"context"
	"log"
	"time"

	"github.com/VictorChepkasov/protos/messenger"
	pb "github.com/VictorChepkasov/protos/messenger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:5051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMessengerClient(conn)

	// Sending a message example
	msg := &pb.Message{Sender: "Alice", Content: "Hello, Bob!"}
	sentMsg, err := client.SendMessage(context.Background(), msg)
	if err != nil {
		log.Fatalf("Could not send message: %v", err)
	}
	log.Printf("Sent message: %v", sentMsg)

	// Receiving messages example
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := client.ReceiveMessages(ctx, &messenger.Empty{})
	if err != nil {
		log.Fatalf("Could not receive messages: %v", err)
	}
	stream.SendMsg(msg)

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("Finished receiving messages: %v", err)
			break
		}
		log.Printf("Received message from %s: %s", msg.Sender, msg.Content)
	}
}
