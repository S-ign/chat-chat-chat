package client

import (
	"context"
	"io"
	"log"

	"github.com/S-ign/chat-chat-chat/src/api/chatpb"
	"google.golang.org/grpc"
)

// Connect connects client to gRPC
func Connect() (chatpb.ChatServiceClient, *grpc.ClientConn) {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	c := chatpb.NewChatServiceClient(cc)

	return c, cc
}

// Chat sends messages to the ChatServiceServer
func Chat(message string, c chatpb.ChatServiceClient) {
	stream, err := c.Chat(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream: %v\n", err)
		return
	}

	waitc := make(chan struct{})

	// generate requests from names
	req := &chatpb.ChatRequest{
		Chatting: &chatpb.Chatting{
			ChatMessage: message,
		},
	}
	// send messsage
	go func() {
		stream.Send(req)
	}()

	// recv messages
	go func() {
		for {
			_, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while receiving: %v\n", err)
				break
			}
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}
