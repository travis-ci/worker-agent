package main

import (
	"log"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/travis-ci/worker-agent/agent"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewAgentClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.RunJob(ctx, &pb.RunJobRequest{
		JobId: "123",
		LogTimeoutS: 10,
		HardTimeoutS: 10,
		MaxLogLength: 10,	
	})
	if err != nil {
		log.Fatalf("could not run job: %v", err)
	}
	log.Printf("Received: %t", r.Ok)
}
