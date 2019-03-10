package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/executor/proto"
)

// Execute executes a command
func (s *Server) Execute(ctx context.Context, req *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {
	sTime := time.Now()
	output, err := s.scheduler.schedule(req.Command)

	return &pb.ExecuteResponse{
		TimeTakenInMillis: time.Now().Sub(sTime).Nanoseconds() / 100000,
		CommandOutput:     output,
	}, err
}

// Execute executes a command
func (s *Server) QueueExecute(ctx context.Context, req *pb.ExecuteRequest) (*pb.ExecuteResponse, error) {
	return nil, fmt.Errorf("Not implemented yet")
}
