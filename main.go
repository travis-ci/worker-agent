package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"

	pb "github.com/travis-ci/worker-agent/agent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = "127.0.0.1:50051"
)

// server is used to implement agent.Agent.
type server struct {
	// TODO: add mutex around this?
	logOutput  []byte
	sentOffset int64
	outChan    chan *pb.LogPart
}

func (s *server) GetJobStatus(ctx context.Context, wr *pb.WorkerRequest) (*pb.JobStatus, error) {
	return &pb.JobStatus{}, nil
}

func (s *server) GetLogParts(wr *pb.LogPartsRequest, stream pb.Agent_GetLogPartsServer) error {
	s.sentOffset = int64(len(s.logOutput))
	err := stream.Send(&pb.LogPart{
		Content: string(s.logOutput[wr.Offset:]),
		Number:  s.sentOffset,
	})
	if err != nil {
		return err
	}

	for part := range s.outChan {
		if part.Number < s.sentOffset {
			continue
		}

		err := stream.Send(part)
		if err != nil {
			return err
		}
		s.sentOffset = part.Number
	}

	return nil
}

func (s *server) RunJob(ctx context.Context, wr *pb.RunJobRequest) (*pb.RunJobResponse, error) {
	cmd := exec.Command("bash", "example/build.sh")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return &pb.RunJobResponse{Ok: false}, err
	}

	err = cmd.Start()
	if err != nil {
		return &pb.RunJobResponse{Ok: false}, err
	}

	reader := bufio.NewReader(stdout)
	go func() {
		offset := 0
		for {
			fmt.Println("reading from stdout")
			out := make([]byte, 512)
			n, err := reader.Read(out)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("failed to read from stdout: %v\n", err)
			}

			offset += n

			s.logOutput = append(s.logOutput, out[:n]...)
			s.outChan <- &pb.LogPart{
				Content: string(out[:n]),
				Number:  int64(offset),
			}
		}
		close(s.outChan)
	}()

	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return &pb.RunJobResponse{Ok: false}, err
	}
	return &pb.RunJobResponse{Ok: true}, nil
}

func main() {
	fmt.Println("starting up...")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAgentServer(s, &server{
		outChan: make(chan *pb.LogPart),
	})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
