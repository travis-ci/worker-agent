package main

import (
	"log"
	"net"

	pb "github.com/travis-ci/worker-agent/agent"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement agent.Agent.
type server struct{}

func (s *server) GetJobStatus(ctx context.Context, wr *pb.WorkerRequest) (*pb.JobStatus, error) {
	return &pb.JobStatus{}, nil
}

func (s *server) GetLogPart(ctx context.Context,  wr *pb.WorkerRequest) (*pb.LogPart, error) {
	return &pb.LogPart{}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAgentServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
