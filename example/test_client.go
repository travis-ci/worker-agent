package main

import (
	"io"
	"log"
	"time"

	pb "github.com/travis-ci/worker-agent/agent"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	r, err := c.RunJob(ctx, &pb.RunJobRequest{
		JobId:        "123",
		LogTimeoutS:  10,
		HardTimeoutS: 10,
		MaxLogLength: 10,
	})
	if err != nil {
		log.Fatalf("could not run job: %v", err)
	}
	log.Printf("Received: %t", r.Ok)

	stream, err := c.GetLogParts(ctx, &pb.WorkerRequest{})
	if err != nil {
		log.Fatalf("could not get log parts: %v", err)
	}
	for {
		part, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetLogParts(_) = _, %v", c, err)
		}
		log.Println("got log part:")
		log.Println(part.Content)
	}
}
