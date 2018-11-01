package main

import (
	"log"
	"net"
	"os/exec"

	pb "github.com/travis-ci/worker-agent/agent"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = "127.0.0.1:50051"
)

// server is used to implement agent.Agent.
type server struct {
	// TODO: add mutex around this?
	logOutput []byte
}

func (s *server) GetJobStatus(ctx context.Context, wr *pb.WorkerRequest) (*pb.JobStatus, error) {
	return &pb.JobStatus{}, nil
}

func (s *server) GetLogParts(wr *pb.WorkerRequest, stream pb.Agent_GetLogPartsServer) error {
	err := stream.Send(&pb.LogPart{
		Content: string(s.logOutput),
		Number:  0,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *server) RunJob(ctx context.Context, wr *pb.RunJobRequest) (*pb.RunJobResponse, error) {
	cmd := exec.Command("bash", "example/build.sh")
	out, err := cmd.CombinedOutput()
	s.logOutput = append(s.logOutput, out...)
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return &pb.RunJobResponse{Ok: false}, err
	}
	return &pb.RunJobResponse{Ok: true}, nil
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
