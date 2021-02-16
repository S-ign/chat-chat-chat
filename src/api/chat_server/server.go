package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/S-ign/chat-chat-chat/src/api/chatpb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

type server struct {
	chatpb.UnimplementedChatServiceServer
}

type chatMessage struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
}

func (*server) Chat(stream chatpb.ChatService_ChatServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return status.Errorf(
				codes.Unknown,
				fmt.Sprintf("error receiving message from client: %v\n", err),
			)
		}
		stream.Send(&chatpb.ChatResponse{
			Chatting: &chatpb.Chatting{
				ChatMessage: req.GetChatting().GetChatMessage(),
			},
		})
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Chat Service Started")

	// Connect to Mongodb
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://MongoAdmin:F00dForThought!@192.81.208.69:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("chat-db").Collection("messages")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	chatpb.RegisterChatServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting gRPC Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the Server")
	s.Stop()
	fmt.Println("Closing the Listener")
	lis.Close()
	fmt.Println("Closing MongoDB connection")
	client.Disconnect(context.TODO())
	fmt.Println("End of Program")
}
