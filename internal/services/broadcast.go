package main

import (
	"context"
	"fmt"
	proto "messenger/internal/gen"
	"sync"
)

type Connection struct {
	proto.UnimplementedBroadcastServer
	stream proto.Broadcast_CreateStreamServer
	id     string
	active bool
	error  chan error
}

type Pool struct {
	proto.UnimplementedBroadcastServer
	Connection []*Connection
}

func (p *Pool) CreateStream(pconn *proto.Connect, stream proto.Broadcast_CreateStreamServer) error {
	conn := &Connection{
		stream: stream,
		id:     pconn.User.Id,
		active: true,
		error:  make(chan error),
	}

	p.Connection = append(p.Connection, conn)

	return <-conn.error
}

func (p *Pool) BroadcastMessage(ctx context.Context, msg *proto.Message) (*proto.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range p.Connection {
		wait.Add(1)

		go func(msg *proto.Message, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(msg)
				fmt.Printf("Sending message to: %v from %v", conn.id, msg.Id)

				if err != nil {
					fmt.Printf("Error withstream: %v - Error: %v\n", conn.stream, err)
					conn.active = false
					conn.error <- err
				}
			}
		}(msg, conn)
	}

	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
	return &proto.Close{}, nil
}
