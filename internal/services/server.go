package main

import (
	"context"
	"log"
	"net"

	pb "github.com/VictorChepkasov/protos/messenger"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedMessengerServer
	messages []pb.Message
}

func (s server) SendMessage(ctx context.Context, msg pb.Message) (pb.Message, error) {
	s.messages = append(s.messages, msg)
	return msg, nil
}

func (s server) ReceiveMessages(empty emptypb.Empty, stream pb.Messenger_ReceiveMessagesServer) error {
	for _, msg := range s.messages {
		if err := stream.Send(&msg); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMessengerServer(grpcServer, &server{})

	log.Println("Server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
