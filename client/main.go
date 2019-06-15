package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	pb "github.com/emilyruby/cling/api"
)

const (
	address     = "localhost:50051"
)

func login(c pb.ClingClient, ctx context.Context, username, password string) (context.Context, error) {
	var header metadata.MD
	_, err := c.Login(ctx, &pb.LoginRequest{Username: username, Password: password}, grpc.Header(&header))

	if err != nil {
		return nil, err
	}

	md, _ := metadata.FromOutgoingContext(ctx) //TODO: should return some error if not ok
	return metadata.NewOutgoingContext(ctx, metadata.Join(md, header)), nil
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewClingClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	ctx, err = login(c, ctx, "test", "this")

	post, err := c.NewPost(ctx, &pb.Post{Content: "hello", Title: "HAHA"})

	log.Println("Greeting: ", post, err)
}