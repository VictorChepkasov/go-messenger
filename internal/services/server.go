package main

import (
	"context"
	"log"
	"net"

	"github.com/VictorChepkasov/protos/messenger"
	pb "github.com/VictorChepkasov/protos/messenger"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMessengerServer
	messages []*pb.Message
}

func (s server) SendMessage(ctx context.Context, msg *pb.Message) (*pb.Message, error) {
	s.messages = append(s.messages, msg)
	return msg, nil
}

func (s server) ReceiveMessages(empty *messenger.Empty, stream pb.Messenger_ReceiveMessagesServer) error {
	for _, msg := range s.messages {
		if err := stream.Send(msg); err != nil {
			log.Fatalf("could not receive messages: %v", err)
			return err
		}
	}
	return nil
}

func (s server) mustEmbedUnimplementedMessengerServer() {}

func main() {
	lis, err := net.Listen("tcp", ":5051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMessengerServer(grpcServer, &server{})

	log.Println("Server is running on port 5051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
